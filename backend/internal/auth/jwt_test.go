package auth

import (
	"testing"
	"time"

	"golearn/internal/models"
)

func TestGenerateAndParseAccessToken(t *testing.T) {
	secret := []byte("test-secret")
	user := models.User{ID: 10, Name: "Test", Email: "t@example.com"}
	ok, err := GenerateAccessToken(user, secret, time.Minute)
	if err != nil {
		t.Fatalf("generate token error: %v", err)
	}
	claims, err := ParseAccessToken(ok, secret)
	if err != nil {
		t.Fatalf("parse token error: %v", err)
	}
	if claims.UserID != user.ID || claims.Email != user.Email {
		t.Fatalf("claims mismatch")
	}
}
