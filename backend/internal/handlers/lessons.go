package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"golearn/internal/middleware"
	"golearn/internal/models"
	"golearn/internal/utils"
)

func (a *App) Lesson(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	idStr := strings.TrimPrefix(r.URL.Path, "/api/lessons/")
	id, err := strconv.Atoi(strings.Trim(idStr, "/"))
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid lesson id")
		return
	}
	u, _ := middleware.UserFromContext(r.Context())
	allowed, err := a.lessonAllowed(u.ID, id)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "db error")
		return
	}
	if !allowed {
		utils.WriteError(w, http.StatusForbidden, "previous lesson not passed")
		return
	}
	var l models.Lesson
	err = a.DB.QueryRow("SELECT id, course_id, title, content, order_index FROM lessons WHERE id = ?", id).Scan(&l.ID, &l.CourseID, &l.Title, &l.Content, &l.OrderIndex)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, "lesson not found")
		return
	}
	questions, err := a.lessonQuizQuestions(id)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "db error")
		return
	}
	utils.WriteJSON(w, http.StatusOK, map[string]any{
		"lesson":    l,
		"questions": questions,
	})
}

func (a *App) LessonQuizSubmit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	idStr := strings.TrimPrefix(r.URL.Path, "/api/lessons-quiz/")
	id, err := strconv.Atoi(strings.Trim(idStr, "/"))
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid lesson id")
		return
	}
	var req struct {
		Answers map[string]int `json:"answers"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid json")
		return
	}
	u, _ := middleware.UserFromContext(r.Context())
	questions, err := a.lessonQuizQuestions(id)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "db error")
		return
	}
	if len(questions) == 0 {
		utils.WriteError(w, http.StatusBadRequest, "no quiz for lesson")
		return
	}
	correct := 0
	for _, q := range questions {
		key := strconv.Itoa(q.ID)
		if ans, ok := req.Answers[key]; ok && ans == q.CorrectIndex {
			correct++
		}
	}
	score := int(float64(correct) / float64(len(questions)) * 100.0)
	passed := score >= 70
	_, _ = a.DB.Exec(`
		INSERT INTO user_lesson_progress(user_id, lesson_id, passed, score, updated_at)
		VALUES(?,?,?,?,?)
		ON CONFLICT(user_id, lesson_id) DO UPDATE SET passed=excluded.passed, score=excluded.score, updated_at=excluded.updated_at
	`,
		u.ID,
		id,
		utils.BoolToInt(passed),
		score,
		time.Now().Format(time.RFC3339),
	)

	utils.WriteJSON(w, http.StatusOK, map[string]any{
		"score":  score,
		"passed": passed,
	})
}

func (a *App) lessonAllowed(userID, lessonID int) (bool, error) {
	var courseID int
	var orderIndex int
	err := a.DB.QueryRow("SELECT course_id, order_index FROM lessons WHERE id = ?", lessonID).Scan(&courseID, &orderIndex)
	if err != nil {
		return false, err
	}
	if orderIndex == 1 {
		return true, nil
	}
	var prevID int
	err = a.DB.QueryRow("SELECT id FROM lessons WHERE course_id = ? AND order_index = ?", courseID, orderIndex-1).Scan(&prevID)
	if err != nil {
		return false, err
	}
	var passed int
	err = a.DB.QueryRow("SELECT COALESCE(passed, 0) FROM user_lesson_progress WHERE user_id = ? AND lesson_id = ?", userID, prevID).Scan(&passed)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return passed == 1, nil
}

func (a *App) lessonQuizQuestions(lessonID int) ([]models.QuizQuestion, error) {
	rows, err := a.DB.Query("SELECT id, question, options_json, correct_index, explanation FROM lesson_quiz_questions WHERE lesson_id = ? ORDER BY id ASC", lessonID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []models.QuizQuestion
	for rows.Next() {
		var q models.QuizQuestion
		var options string
		if err := rows.Scan(&q.ID, &q.Question, &options, &q.CorrectIndex, &q.Explanation); err != nil {
			return nil, err
		}
		_ = json.Unmarshal([]byte(options), &q.Options)
		result = append(result, q)
	}
	return result, nil
}
