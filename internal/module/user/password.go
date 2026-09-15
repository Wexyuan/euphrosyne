package user

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// minPasswordLen and maxPasswordLen bound the accepted plain password length.
const (
	minPasswordLen = 8
	maxPasswordLen = 32
)

// HashPassword hashes the plain password with bcrypt.
func HashPassword(plain string) (string, error) {
	if len([]rune(plain)) < minPasswordLen || len([]rune(plain)) > maxPasswordLen {
		return "", fmt.Errorf("[user] password length must be between %d and %d", minPasswordLen, maxPasswordLen)
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("[user] hash password error: %w", err)
	}
	return string(hashed), nil
}

// VerifyPassword reports whether the plain password matches the bcrypt hash.
func VerifyPassword(hash, plain string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain)); err != nil {
		return fmt.Errorf("[user] verify password error: %w", err)
	}
	return nil
}
