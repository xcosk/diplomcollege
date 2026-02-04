package auth

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"time"
)

func GenerateRefreshToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func StoreRefreshToken(db *sql.DB, token string, userID int, expires time.Time) error {
	_, err := db.Exec(
		"INSERT INTO refresh_tokens(token, user_id, expires_at) VALUES(?,?,?)",
		token,
		userID,
		expires.Format(time.RFC3339),
	)
	return err
}

func DeleteRefreshToken(db *sql.DB, token string) error {
	_, err := db.Exec("DELETE FROM refresh_tokens WHERE token = ?", token)
	return err
}

func GetRefreshToken(db *sql.DB, token string) (int, time.Time, error) {
	var userID int
	var expires string
	if err := db.QueryRow("SELECT user_id, expires_at FROM refresh_tokens WHERE token = ?", token).Scan(&userID, &expires); err != nil {
		return 0, time.Time{}, err
	}
	parsed, err := time.Parse(time.RFC3339, expires)
	if err != nil {
		return 0, time.Time{}, err
	}
	return userID, parsed, nil
}
