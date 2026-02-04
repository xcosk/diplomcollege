package middleware

import (
	"database/sql"
	"net/http"

	"golearn/internal/utils"
)

func Admin(db *sql.DB, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := UserFromContext(r.Context())
		if !ok {
			utils.WriteError(w, http.StatusUnauthorized, "unauthorized")
			return
		}
		var isAdmin int
		err := db.QueryRow("SELECT is_admin FROM users WHERE id = ?", user.ID).Scan(&isAdmin)
		if err != nil {
			utils.WriteError(w, http.StatusForbidden, "admin only")
			return
		}
		if isAdmin != 1 {
			utils.WriteError(w, http.StatusForbidden, "admin only")
			return
		}
		next(w, r)
	}
}
