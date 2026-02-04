package db

import "database/sql"

func Migrate(db *sql.DB) error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			email TEXT NOT NULL UNIQUE,
			password_hash TEXT NOT NULL,
			is_admin INTEGER NOT NULL DEFAULT 0,
			created_at TEXT NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS refresh_tokens (
			token TEXT PRIMARY KEY,
			user_id INTEGER NOT NULL,
			expires_at TEXT NOT NULL,
			FOREIGN KEY(user_id) REFERENCES users(id)
		)`,
		`CREATE TABLE IF NOT EXISTS courses (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			level TEXT NOT NULL,
			title TEXT NOT NULL,
			description TEXT NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS lessons (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			course_id INTEGER NOT NULL,
			title TEXT NOT NULL,
			content TEXT NOT NULL,
			order_index INTEGER NOT NULL,
			FOREIGN KEY(course_id) REFERENCES courses(id)
		)`,
		`CREATE TABLE IF NOT EXISTS placement_questions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			question TEXT NOT NULL,
			options_json TEXT NOT NULL,
			correct_index INTEGER NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS lesson_quiz_questions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			lesson_id INTEGER NOT NULL,
			question TEXT NOT NULL,
			options_json TEXT NOT NULL,
			correct_index INTEGER NOT NULL,
			explanation TEXT NOT NULL,
			FOREIGN KEY(lesson_id) REFERENCES lessons(id)
		)`,
		`CREATE TABLE IF NOT EXISTS user_lesson_progress (
			user_id INTEGER NOT NULL,
			lesson_id INTEGER NOT NULL,
			passed INTEGER NOT NULL,
			score INTEGER NOT NULL,
			updated_at TEXT NOT NULL,
			PRIMARY KEY(user_id, lesson_id),
			FOREIGN KEY(user_id) REFERENCES users(id),
			FOREIGN KEY(lesson_id) REFERENCES lessons(id)
		)`,
	}

	for _, q := range queries {
		if _, err := db.Exec(q); err != nil {
			return err
		}
	}
	// Backward compatible: add is_admin if table existed before
	_, _ = db.Exec("ALTER TABLE users ADD COLUMN is_admin INTEGER NOT NULL DEFAULT 0")
	return nil
}
