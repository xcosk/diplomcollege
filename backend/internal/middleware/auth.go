package middleware

import (
	"context"
	"net/http"
	"strings"

	"golearn/internal/auth"
	"golearn/internal/models"
	"golearn/internal/utils"
)

type ctxKey string

const userKey ctxKey = "user"

func UserFromContext(ctx context.Context) (models.User, bool) {
	u, ok := ctx.Value(userKey).(models.User)
	return u, ok
}

func Auth(secret []byte, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
		if token == "" {
			utils.WriteError(w, http.StatusUnauthorized, "missing token")
			return
		}
		claims, err := auth.ParseAccessToken(token, secret)
		if err != nil {
			utils.WriteError(w, http.StatusUnauthorized, "invalid token")
			return
		}
		ctx := context.WithValue(r.Context(), userKey, models.User{
			ID:    claims.UserID,
			Name:  claims.Name,
			Email: claims.Email,
		})
		next(w, r.WithContext(ctx))
	}
}
