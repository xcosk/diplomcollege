package handlers

import (
	"database/sql"
	"net/http"

	"golearn/internal/middleware"
)

type App struct {
	DB     *sql.DB
	Secret []byte
}

func NewApp(db *sql.DB, secret []byte) *App {
	return &App{DB: db, Secret: secret}
}

func (a *App) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/register", a.Register)
	mux.HandleFunc("/api/login", a.Login)
	mux.HandleFunc("/api/refresh", a.Refresh)
	mux.HandleFunc("/api/logout", a.Auth(a.Logout))
	mux.HandleFunc("/api/me", a.Auth(a.Me))

	mux.HandleFunc("/api/placement-test", a.PlacementQuestions)
	mux.HandleFunc("/api/placement-test/submit", a.PlacementSubmit)

	mux.HandleFunc("/api/courses", a.Auth(a.Courses))
	mux.HandleFunc("/api/courses/", a.Auth(a.CourseDetail))
	mux.HandleFunc("/api/lessons/", a.Auth(a.Lesson))
	mux.HandleFunc("/api/lessons-quiz/", a.Auth(a.LessonQuizSubmit))
	mux.HandleFunc("/api/progress", a.Auth(a.Progress))

	// Admin
	mux.HandleFunc("/api/admin/courses", a.Admin(a.AdminCourses))
	mux.HandleFunc("/api/admin/courses/", a.Admin(a.AdminCoursesSub))
	mux.HandleFunc("/api/admin/lessons/", a.Admin(a.AdminLesson))
	mux.HandleFunc("/api/admin/quiz/", a.Admin(a.AdminQuiz))
	mux.HandleFunc("/api/admin/placement", a.Admin(a.AdminPlacement))
	mux.HandleFunc("/api/admin/placement/", a.Admin(a.AdminPlacementItem))
}

func (a *App) Admin(next http.HandlerFunc) http.HandlerFunc {
	return a.Auth(middleware.Admin(a.DB, next))
}
