package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"golearn/internal/models"
	"golearn/internal/utils"
)

func (a *App) PlacementQuestions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	rows, err := a.DB.Query("SELECT id, question, options_json FROM placement_questions ORDER BY id ASC")
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "db error")
		return
	}
	defer rows.Close()
	var result []models.QuizQuestion
	for rows.Next() {
		var q models.QuizQuestion
		var options string
		if err := rows.Scan(&q.ID, &q.Question, &options); err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "db error")
			return
		}
		_ = json.Unmarshal([]byte(options), &q.Options)
		result = append(result, q)
	}
	utils.WriteJSON(w, http.StatusOK, result)
}

func (a *App) PlacementSubmit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	var req struct {
		Answers map[string]int `json:"answers"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid json")
		return
	}
	if len(req.Answers) == 0 {
		utils.WriteError(w, http.StatusBadRequest, "answers required")
		return
	}
	rows, err := a.DB.Query("SELECT id, correct_index FROM placement_questions")
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "db error")
		return
	}
	defer rows.Close()
	correct := 0
	total := 0
	for rows.Next() {
		var id int
		var correctIndex int
		if err := rows.Scan(&id, &correctIndex); err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "db error")
			return
		}
		total++
		if ans, ok := req.Answers[strconv.Itoa(id)]; ok && ans == correctIndex {
			correct++
		}
	}
	level := placementLevel(correct)
	utils.WriteJSON(w, http.StatusOK, map[string]any{
		"score":   correct,
		"total":   total,
		"level":   level,
		"message": placementMessage(level),
	})
}

func placementLevel(correct int) string {
	if correct >= 8 {
		return "pro"
	}
	if correct >= 4 {
		return "mid"
	}
	return "base"
}

func placementMessage(level string) string {
	switch level {
	case "pro":
		return "Вам подойдет профессиональный курс: сложные темы, архитектура, производительность и production-практики."
	case "mid":
		return "Вам подойдет средний курс: интерфейсы, конкурентность, тесты и работа с данными."
	default:
		return "Вам подойдет базовый курс: основы синтаксиса, типы, функции и базовые структуры данных."
	}
}
