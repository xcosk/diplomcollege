package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"golearn/internal/models"
	"golearn/internal/utils"
)

func (a *App) AdminCourses(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/api/admin/courses" {
		utils.WriteError(w, http.StatusNotFound, "not found")
		return
	}
	switch r.Method {
	case http.MethodGet:
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
	case http.MethodPost:
		var req struct {
			Level       string `json:"level"`
			Title       string `json:"title"`
			Description string `json:"description"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.WriteError(w, http.StatusBadRequest, "invalid json")
			return
		}
		if req.Level == "" || req.Title == "" || req.Description == "" {
			utils.WriteError(w, http.StatusBadRequest, "level, title and description required")
			return
		}
		res, err := a.DB.Exec("INSERT INTO courses(level, title, description) VALUES(?,?,?)", req.Level, req.Title, req.Description)
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "db error")
			return
		}
		id, _ := res.LastInsertId()
		utils.WriteJSON(w, http.StatusCreated, map[string]any{"id": id})
	default:
		utils.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (a *App) AdminCoursesSub(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/admin/courses/")
	path = strings.Trim(path, "/")
	parts := strings.Split(path, "/")
	if len(parts) == 0 || parts[0] == "" {
		utils.WriteError(w, http.StatusNotFound, "not found")
		return
	}
	courseID, err := strconv.Atoi(parts[0])
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid course id")
		return
	}

	if len(parts) == 1 {
		switch r.Method {
		case http.MethodPut:
			var req struct {
				Level       string `json:"level"`
				Title       string `json:"title"`
				Description string `json:"description"`
			}
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				utils.WriteError(w, http.StatusBadRequest, "invalid json")
				return
			}
			_, err := a.DB.Exec("UPDATE courses SET level=?, title=?, description=? WHERE id = ?", req.Level, req.Title, req.Description, courseID)
			if err != nil {
				utils.WriteError(w, http.StatusInternalServerError, "db error")
				return
			}
			utils.WriteJSON(w, http.StatusOK, map[string]any{"ok": true})
		case http.MethodDelete:
			tx, err := a.DB.Begin()
			if err != nil {
				utils.WriteError(w, http.StatusInternalServerError, "db error")
				return
			}
			_, _ = tx.Exec("DELETE FROM lesson_quiz_questions WHERE lesson_id IN (SELECT id FROM lessons WHERE course_id = ?)", courseID)
			_, _ = tx.Exec("DELETE FROM user_lesson_progress WHERE lesson_id IN (SELECT id FROM lessons WHERE course_id = ?)", courseID)
			_, _ = tx.Exec("DELETE FROM lessons WHERE course_id = ?", courseID)
			_, err = tx.Exec("DELETE FROM courses WHERE id = ?", courseID)
			if err != nil {
				_ = tx.Rollback()
				utils.WriteError(w, http.StatusInternalServerError, "db error")
				return
			}
			_ = tx.Commit()
			utils.WriteJSON(w, http.StatusOK, map[string]any{"ok": true})
		default:
			utils.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		}
		return
	}

	if len(parts) == 2 && parts[1] == "lessons" {
		switch r.Method {
		case http.MethodGet:
			rows, err := a.DB.Query("SELECT id, course_id, title, content, order_index FROM lessons WHERE course_id = ? ORDER BY order_index ASC", courseID)
			if err != nil {
				utils.WriteError(w, http.StatusInternalServerError, "db error")
				return
			}
			defer rows.Close()
			var result []models.Lesson
			for rows.Next() {
				var l models.Lesson
				if err := rows.Scan(&l.ID, &l.CourseID, &l.Title, &l.Content, &l.OrderIndex); err != nil {
					utils.WriteError(w, http.StatusInternalServerError, "db error")
					return
				}
				result = append(result, l)
			}
			utils.WriteJSON(w, http.StatusOK, result)
		case http.MethodPost:
			var req struct {
				Title      string `json:"title"`
				Content    string `json:"content"`
				OrderIndex int    `json:"order_index"`
			}
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				utils.WriteError(w, http.StatusBadRequest, "invalid json")
				return
			}
			if req.Title == "" || req.Content == "" {
				utils.WriteError(w, http.StatusBadRequest, "title and content required")
				return
			}
			orderIndex := req.OrderIndex
			if orderIndex == 0 {
				_ = a.DB.QueryRow("SELECT COALESCE(MAX(order_index), 0) + 1 FROM lessons WHERE course_id = ?", courseID).Scan(&orderIndex)
			}
			res, err := a.DB.Exec("INSERT INTO lessons(course_id, title, content, order_index) VALUES(?,?,?,?)", courseID, req.Title, req.Content, orderIndex)
			if err != nil {
				utils.WriteError(w, http.StatusInternalServerError, "db error")
				return
			}
			id, _ := res.LastInsertId()
			utils.WriteJSON(w, http.StatusCreated, map[string]any{"id": id})
		default:
			utils.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		}
		return
	}

	utils.WriteError(w, http.StatusNotFound, "not found")
}

func (a *App) AdminLesson(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/admin/lessons/")
	path = strings.Trim(path, "/")
	parts := strings.Split(path, "/")
	if len(parts) == 0 || parts[0] == "" {
		utils.WriteError(w, http.StatusNotFound, "not found")
		return
	}
	lessonID, err := strconv.Atoi(parts[0])
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid lesson id")
		return
	}

	if len(parts) == 1 {
		switch r.Method {
		case http.MethodGet:
			var l models.Lesson
			err := a.DB.QueryRow("SELECT id, course_id, title, content, order_index FROM lessons WHERE id = ?", lessonID).Scan(&l.ID, &l.CourseID, &l.Title, &l.Content, &l.OrderIndex)
			if err != nil {
				utils.WriteError(w, http.StatusNotFound, "lesson not found")
				return
			}
			utils.WriteJSON(w, http.StatusOK, l)
		case http.MethodPut:
			var req struct {
				Title      string `json:"title"`
				Content    string `json:"content"`
				OrderIndex int    `json:"order_index"`
			}
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				utils.WriteError(w, http.StatusBadRequest, "invalid json")
				return
			}
			_, err := a.DB.Exec("UPDATE lessons SET title=?, content=?, order_index=? WHERE id = ?", req.Title, req.Content, req.OrderIndex, lessonID)
			if err != nil {
				utils.WriteError(w, http.StatusInternalServerError, "db error")
				return
			}
			utils.WriteJSON(w, http.StatusOK, map[string]any{"ok": true})
		case http.MethodDelete:
			_, err := a.DB.Exec("DELETE FROM lessons WHERE id = ?", lessonID)
			if err != nil {
				utils.WriteError(w, http.StatusInternalServerError, "db error")
				return
			}
			utils.WriteJSON(w, http.StatusOK, map[string]any{"ok": true})
		default:
			utils.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		}
		return
	}

	if len(parts) == 2 && parts[1] == "quiz" {
		switch r.Method {
		case http.MethodGet:
			rows, err := a.DB.Query("SELECT id, question, options_json, correct_index, explanation FROM lesson_quiz_questions WHERE lesson_id = ? ORDER BY id ASC", lessonID)
			if err != nil {
				utils.WriteError(w, http.StatusInternalServerError, "db error")
				return
			}
			defer rows.Close()
			var result []models.QuizQuestion
			for rows.Next() {
				var q models.QuizQuestion
				var options string
				if err := rows.Scan(&q.ID, &q.Question, &options, &q.CorrectIndex, &q.Explanation); err != nil {
					utils.WriteError(w, http.StatusInternalServerError, "db error")
					return
				}
				_ = json.Unmarshal([]byte(options), &q.Options)
				result = append(result, q)
			}
			utils.WriteJSON(w, http.StatusOK, result)
		case http.MethodPost:
			var req struct {
				Question     string   `json:"question"`
				Options      []string `json:"options"`
				CorrectIndex int      `json:"correct_index"`
				Explanation  string   `json:"explanation"`
			}
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				utils.WriteError(w, http.StatusBadRequest, "invalid json")
				return
			}
			if req.Question == "" || len(req.Options) < 2 {
				utils.WriteError(w, http.StatusBadRequest, "question and options required")
				return
			}
			options, _ := json.Marshal(req.Options)
			res, err := a.DB.Exec("INSERT INTO lesson_quiz_questions(lesson_id, question, options_json, correct_index, explanation) VALUES(?,?,?,?,?)", lessonID, req.Question, string(options), req.CorrectIndex, req.Explanation)
			if err != nil {
				utils.WriteError(w, http.StatusInternalServerError, "db error")
				return
			}
			id, _ := res.LastInsertId()
			utils.WriteJSON(w, http.StatusCreated, map[string]any{"id": id})
		default:
			utils.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		}
		return
	}

	utils.WriteError(w, http.StatusNotFound, "not found")
}

func (a *App) AdminQuiz(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/admin/quiz/")
	path = strings.Trim(path, "/")
	id, err := strconv.Atoi(path)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid quiz id")
		return
	}
	if r.Method == http.MethodPut {
		var req struct {
			Question     string   `json:"question"`
			Options      []string `json:"options"`
			CorrectIndex int      `json:"correct_index"`
			Explanation  string   `json:"explanation"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.WriteError(w, http.StatusBadRequest, "invalid json")
			return
		}
		options, _ := json.Marshal(req.Options)
		_, err := a.DB.Exec("UPDATE lesson_quiz_questions SET question=?, options_json=?, correct_index=?, explanation=? WHERE id = ?", req.Question, string(options), req.CorrectIndex, req.Explanation, id)
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "db error")
			return
		}
		utils.WriteJSON(w, http.StatusOK, map[string]any{"ok": true})
		return
	}
	if r.Method == http.MethodDelete {
		_, err := a.DB.Exec("DELETE FROM lesson_quiz_questions WHERE id = ?", id)
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "db error")
			return
		}
		utils.WriteJSON(w, http.StatusOK, map[string]any{"ok": true})
		return
	}
	utils.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
}

func (a *App) AdminPlacement(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/api/admin/placement" {
		utils.WriteError(w, http.StatusNotFound, "not found")
		return
	}
	switch r.Method {
	case http.MethodGet:
		rows, err := a.DB.Query("SELECT id, question, options_json, correct_index FROM placement_questions ORDER BY id ASC")
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "db error")
			return
		}
		defer rows.Close()
		var result []models.QuizQuestion
		for rows.Next() {
			var q models.QuizQuestion
			var options string
			if err := rows.Scan(&q.ID, &q.Question, &options, &q.CorrectIndex); err != nil {
				utils.WriteError(w, http.StatusInternalServerError, "db error")
				return
			}
			_ = json.Unmarshal([]byte(options), &q.Options)
			result = append(result, q)
		}
		utils.WriteJSON(w, http.StatusOK, result)
	case http.MethodPost:
		var req struct {
			Question     string   `json:"question"`
			Options      []string `json:"options"`
			CorrectIndex int      `json:"correct_index"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.WriteError(w, http.StatusBadRequest, "invalid json")
			return
		}
		options, _ := json.Marshal(req.Options)
		res, err := a.DB.Exec("INSERT INTO placement_questions(question, options_json, correct_index) VALUES(?,?,?)", req.Question, string(options), req.CorrectIndex)
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "db error")
			return
		}
		id, _ := res.LastInsertId()
		utils.WriteJSON(w, http.StatusCreated, map[string]any{"id": id})
	default:
		utils.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (a *App) AdminPlacementItem(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/admin/placement/")
	path = strings.Trim(path, "/")
	id, err := strconv.Atoi(path)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid placement id")
		return
	}
	if r.Method == http.MethodPut {
		var req struct {
			Question     string   `json:"question"`
			Options      []string `json:"options"`
			CorrectIndex int      `json:"correct_index"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.WriteError(w, http.StatusBadRequest, "invalid json")
			return
		}
		options, _ := json.Marshal(req.Options)
		_, err := a.DB.Exec("UPDATE placement_questions SET question=?, options_json=?, correct_index=? WHERE id = ?", req.Question, string(options), req.CorrectIndex, id)
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "db error")
			return
		}
		utils.WriteJSON(w, http.StatusOK, map[string]any{"ok": true})
		return
	}
	if r.Method == http.MethodDelete {
		_, err := a.DB.Exec("DELETE FROM placement_questions WHERE id = ?", id)
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "db error")
			return
		}
		utils.WriteJSON(w, http.StatusOK, map[string]any{"ok": true})
		return
	}
	utils.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
}
