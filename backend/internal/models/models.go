package models

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	IsAdmin bool `json:"is_admin"`
}

type Course struct {
	ID          int    `json:"id"`
	Level       string `json:"level"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type Lesson struct {
	ID         int    `json:"id"`
	CourseID   int    `json:"course_id"`
	Title      string `json:"title"`
	Content    string `json:"content"`
	OrderIndex int    `json:"order_index"`
}

type QuizQuestion struct {
	ID           int      `json:"id"`
	Question     string   `json:"question"`
	Options      []string `json:"options"`
	CorrectIndex int      `json:"-"`
	Explanation  string   `json:"explanation"`
}
