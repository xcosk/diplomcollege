package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"golearn/internal/middleware"
	"golearn/internal/models"
	"golearn/internal/utils"
)

func (a *App) Courses(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	rows, err := a.DB.Query("SELECT id, level, title, description FROM courses ORDER BY id ASC")
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "db error")
		return
	}
	defer rows.Close()
	var result []models.Course
	for rows.Next() {
		var c models.Course
		if err := rows.Scan(&c.ID, &c.Level, &c.Title, &c.Description); err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "db error")
			return
		}
		result = append(result, c)
	}
	utils.WriteJSON(w, http.StatusOK, result)
}

func (a *App) CourseDetail(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	idStr := strings.TrimPrefix(r.URL.Path, "/api/courses/")
	id, err := strconv.Atoi(strings.Trim(idStr, "/"))
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid course id")
		return
	}
	var c models.Course
	err = a.DB.QueryRow("SELECT id, level, title, description FROM courses WHERE id = ?", id).Scan(&c.ID, &c.Level, &c.Title, &c.Description)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, "course not found")
		return
	}

	u, _ := middleware.UserFromContext(r.Context())
	rows, err := a.DB.Query(`
		SELECT l.id, l.title, l.order_index,
			COALESCE(p.passed, 0) AS passed
		FROM lessons l
		LEFT JOIN user_lesson_progress p
			ON p.lesson_id = l.id AND p.user_id = ?
		WHERE l.course_id = ?
		ORDER BY l.order_index ASC
	`, u.ID, id)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "db error")
		return
	}
	defer rows.Close()
	lessons := []map[string]any{}
	prevPassed := true
	for rows.Next() {
		var lessonID int
		var title string
		var orderIndex int
		var passed int
		if err := rows.Scan(&lessonID, &title, &orderIndex, &passed); err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "db error")
			return
		}
		locked := !prevPassed
		if passed == 1 {
			prevPassed = true
		} else {
			prevPassed = false
		}
		lessons = append(lessons, map[string]any{
			"id":          lessonID,
			"title":       title,
			"order_index": orderIndex,
			"passed":      passed == 1,
			"locked":      locked,
		})
	}

	utils.WriteJSON(w, http.StatusOK, map[string]any{
		"course":  c,
		"lessons": lessons,
	})
}
