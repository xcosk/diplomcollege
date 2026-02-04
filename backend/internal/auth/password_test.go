package auth

import "testing"

func TestHashAndCheckPassword(t *testing.T) {
	password := "super-secret"
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("hash error: %v", err)
	}
	if hash == "" {
		t.Fatal("expected non-empty hash")
	}
	if !CheckPassword(hash, password) {
		t.Fatal("expected password to match hash")
	}
	if CheckPassword(hash, "wrong") {
		t.Fatal("expected wrong password to fail")
	}
}
