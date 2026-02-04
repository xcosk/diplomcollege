package handlers

import (
	"net/http"

	"golearn/internal/middleware"
	"golearn/internal/utils"
)

func (a *App) Progress(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	u, _ := middleware.UserFromContext(r.Context())
	rows, err := a.DB.Query(`
		SELECT c.id, c.title,
			COUNT(l.id) AS total_lessons,
			SUM(CASE WHEN p.lesson_id IS NOT NULL THEN 1 ELSE 0 END) AS attempts,
			SUM(CASE WHEN p.passed = 1 THEN 1 ELSE 0 END) AS passed_lessons
		FROM courses c
		JOIN lessons l ON l.course_id = c.id
		LEFT JOIN user_lesson_progress p ON p.lesson_id = l.id AND p.user_id = ?
		GROUP BY c.id, c.title
		ORDER BY c.id ASC
	`, u.ID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "db error")
		return
	}
	defer rows.Close()
	var result []map[string]any
	for rows.Next() {
		var id int
		var title string
		var total int
		var attempts int
		var passed int
		if err := rows.Scan(&id, &title, &total, &attempts, &passed); err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "db error")
			return
		}
		progress := 0
		if total > 0 {
			progress = int(float64(passed) / float64(total) * 100.0)
		}
		result = append(result, map[string]any{
			"course_id": id,
			"title":     title,
			"total":     total,
			"attempts":  attempts,
			"passed":    passed,
			"progress":  progress,
		})
	}
	utils.WriteJSON(w, http.StatusOK, result)
}
