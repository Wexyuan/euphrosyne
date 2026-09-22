package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/golang-jwt/jwt/v5"
)

// refreshTokenBytes is the random byte length of a refresh token.
const refreshTokenBytes = 32

// refreshTokenSeparator is the separator between the user ID and the token secret.
const refreshTokenSeparator = "."

// JWTClaims holds the access token claims.
type JWTClaims struct {
	jwt.RegisteredClaims
	UserID   int64  `json:"uid"`   // token owner ID
	Username string `json:"uname"` // token owner name
}

// TokenManager generates and verifies tokens.
type TokenManager struct {
	secret     []byte        // HS256 signing key
	issuer     string        // issuer claim
	accessTTL  time.Duration // access token lifetime
	refreshTTL time.Duration // refresh token lifetime
}

func NewTokenManager(secret, issuer string, accessTTL, refreshTTL time.Duration) (*TokenManager, error) {
	// Refuse to start with an empty secret so tokens are never signed weakly.
	if secret == "" {
		return nil, fmt.Errorf("[auth] jwt secret is required")
	}
	if accessTTL <= 0 {
		accessTTL = 15 * time.Minute
	}
	if refreshTTL <= 0 {
		refreshTTL = 168 * time.Hour
	}
	return &TokenManager{
		secret:     []byte(secret),
		issuer:     issuer,
		accessTTL:  accessTTL,
		refreshTTL: refreshTTL,
	}, nil
}

// GenerateAccessToken generates an access token.
func (m *TokenManager) GenerateAccessToken(userID int64, username string) (string, time.Duration, error) {
	now := time.Now()
	// Use a unique jti so tokens issued within the same second stay distinguishable.
	id, err := uuid.NewV4()
	if err != nil {
		return "", 0, fmt.Errorf("[auth] generate token id error: %w", err)
	}

	claims := JWTClaims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    m.issuer,
			Subject:   fmt.Sprintf("%d", userID),
			ID:        id.String(),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(m.accessTTL)),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(m.secret)
	if err != nil {
		return "", 0, fmt.Errorf("[auth] sign access token error: %w", err)
	}
	return token, m.accessTTL, nil
}

// VerifyAccessToken verifies the access token.
func (m *TokenManager) VerifyAccessToken(token string) (int64, error) {
	claims := &JWTClaims{}
	parsed, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("[auth] unexpected signing method %q", t.Header["alg"])
		}
		return m.secret, nil
	},
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		jwt.WithIssuer(m.issuer),
		jwt.WithExpirationRequired(),
		jwt.WithLeeway(30*time.Second),
	)
	if err != nil {
		return 0, fmt.Errorf("[auth] parse access token error: %w", err)
	}
	if !parsed.Valid || claims.UserID <= 0 {
		return 0, fmt.Errorf("[auth] invalid access token")
	}
	return claims.UserID, nil
}

// GenerateRefreshToken generates a refresh token.
func (m *TokenManager) GenerateRefreshToken(userID int64) (string, error) {
	buf := make([]byte, refreshTokenBytes)
	if _, err := rand.Read(buf); err != nil {
		return "", fmt.Errorf("[auth] generate refresh token error: %w", err)
	}
	return strconv.FormatInt(userID, 10) + refreshTokenSeparator + base64.RawURLEncoding.EncodeToString(buf), nil
}

// ParseRefreshToken parses the refresh token.
func (m *TokenManager) ParseRefreshToken(token string) (int64, bool) {
	id, secret, ok := strings.Cut(token, refreshTokenSeparator)
	// Reject a token without a secret so it can never match a stored digest.
	if !ok || secret == "" {
		return 0, false
	}

	userID, err := strconv.ParseInt(id, 10, 64)
	if err != nil || userID <= 0 {
		return 0, false
	}
	return userID, true
}

// HashRefreshToken hashes the refresh token.
func (m *TokenManager) HashRefreshToken(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}

// RefreshTTL returns the refresh token lifetime.
func (m *TokenManager) RefreshTTL() time.Duration {
	return m.refreshTTL
}
