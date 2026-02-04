package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"golearn/internal/auth"
	"golearn/internal/middleware"
	"golearn/internal/models"
	"golearn/internal/utils"
)

const accessTokenTTL = 15 * time.Minute
const refreshTokenTTL = 30 * 24 * time.Hour

func (a *App) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	var req struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid json")
		return
	}
	if req.Name == "" || req.Email == "" || req.Password == "" {
		utils.WriteError(w, http.StatusBadRequest, "name, email and password are required")
		return
	}
	hash, err := auth.HashPassword(req.Password)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "hash error")
		return
	}
	res, err := a.DB.Exec(
		"INSERT INTO users(name, email, password_hash, is_admin, created_at) VALUES(?,?,?,?,?)",
		req.Name,
		strings.ToLower(req.Email),
		hash,
		0,
		time.Now().Format(time.RFC3339),
	)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE") {
			utils.WriteError(w, http.StatusBadRequest, "email already exists")
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, "db error")
		return
	}
	id, _ := res.LastInsertId()
	utils.WriteJSON(w, http.StatusCreated, map[string]any{
		"id":    id,
		"name":  req.Name,
		"email": strings.ToLower(req.Email),
	})
}

func (a *App) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid json")
		return
	}
	var id int
	var name string
	var email string
	var hash string
	var isAdmin int
	err := a.DB.QueryRow(
		"SELECT id, name, email, password_hash, is_admin FROM users WHERE email = ?",
		strings.ToLower(req.Email),
	).Scan(&id, &name, &email, &hash, &isAdmin)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, "invalid credentials")
		return
	}
	if !auth.CheckPassword(hash, req.Password) {
		utils.WriteError(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	user := models.User{ID: id, Name: name, Email: email, IsAdmin: isAdmin == 1}
	accessToken, err := auth.GenerateAccessToken(user, a.Secret, accessTokenTTL)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "token error")
		return
	}
	refreshToken, err := auth.GenerateRefreshToken()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "token error")
		return
	}
	if err := auth.StoreRefreshToken(a.DB, refreshToken, id, time.Now().Add(refreshTokenTTL)); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "db error")
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]any{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"user":          user,
	})
}

func (a *App) Refresh(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid json")
		return
	}
	if req.RefreshToken == "" {
		utils.WriteError(w, http.StatusBadRequest, "refresh token required")
		return
	}
	userID, expires, err := auth.GetRefreshToken(a.DB, req.RefreshToken)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.WriteError(w, http.StatusUnauthorized, "invalid refresh token")
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, "db error")
		return
	}
	if time.Now().After(expires) {
		_ = auth.DeleteRefreshToken(a.DB, req.RefreshToken)
		utils.WriteError(w, http.StatusUnauthorized, "refresh token expired")
		return
	}
	var user models.User
	var isAdmin int
	err = a.DB.QueryRow("SELECT id, name, email, is_admin FROM users WHERE id = ?", userID).Scan(&user.ID, &user.Name, &user.Email, &isAdmin)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "db error")
		return
	}
	user.IsAdmin = isAdmin == 1
	accessToken, err := auth.GenerateAccessToken(user, a.Secret, accessTokenTTL)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "token error")
		return
	}
	utils.WriteJSON(w, http.StatusOK, map[string]any{"access_token": accessToken})
}

func (a *App) Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	_ = json.NewDecoder(r.Body).Decode(&req)
	if req.RefreshToken != "" {
		_ = auth.DeleteRefreshToken(a.DB, req.RefreshToken)
	}
	utils.WriteJSON(w, http.StatusOK, map[string]any{"ok": true})
}

func (a *App) Me(w http.ResponseWriter, r *http.Request) {
	user, _ := middleware.UserFromContext(r.Context())
	var isAdmin int
	_ = a.DB.QueryRow("SELECT is_admin FROM users WHERE id = ?", user.ID).Scan(&isAdmin)
	user.IsAdmin = isAdmin == 1
	utils.WriteJSON(w, http.StatusOK, user)
}

func (a *App) Auth(next http.HandlerFunc) http.HandlerFunc {
	return middleware.Auth(a.Secret, next)
}
