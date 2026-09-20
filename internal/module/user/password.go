package user

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

const (
	// minPasswordLen is the minimum plain password length.
	minPasswordLen = 8
	// maxPasswordLen is the maximum plain password length.
	maxPasswordLen = 20
	// maxPasswordBytes is the maximum plain password length in bytes.
	maxPasswordBytes = 72
)

// HashPassword hashes the plain password.
func HashPassword(plain string) (string, error) {
	if len([]rune(plain)) < minPasswordLen || len([]rune(plain)) > maxPasswordLen {
		return "", fmt.Errorf("[user] password length must be between %d and %d", minPasswordLen, maxPasswordLen)
	}
	// Reject multibyte passwords above the bcrypt input limit before hashing.
	if len(plain) > maxPasswordBytes {
		return "", fmt.Errorf("[user] password must not exceed %d bytes", maxPasswordBytes)
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
