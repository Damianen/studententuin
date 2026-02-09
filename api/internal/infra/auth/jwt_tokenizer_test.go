package auth

import (
	"api/internal/infra/utils"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func newTestTokenizer(secret string) *JwtTokenizer {
	return &JwtTokenizer{
		Clock:     utils.SystemClock{},
		SecretKey: secret,
	}
}

func TestJwtTokenizer_CreateToken(t *testing.T) {
	tok := newTestTokenizer("test-secret")

	tokenStr, err := tok.CreateToken("user-123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tokenStr == "" {
		t.Fatal("expected non-empty token")
	}

	// Verify it's a parseable JWT
	parsed, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (any, error) {
		return []byte("test-secret"), nil
	})
	if err != nil {
		t.Fatalf("failed to parse token: %v", err)
	}
	claims, ok := parsed.Claims.(*Claims)
	if !ok {
		t.Fatal("failed to cast claims")
	}
	if claims.UserID != "user-123" {
		t.Fatalf("expected UserID user-123, got %s", claims.UserID)
	}
}

func TestJwtTokenizer_VerifyToken(t *testing.T) {
	tok := newTestTokenizer("test-secret")

	tokenStr, err := tok.CreateToken("user-456")
	if err != nil {
		t.Fatalf("unexpected error creating token: %v", err)
	}

	userID, err := tok.VerifyToken(tokenStr)
	if err != nil {
		t.Fatalf("unexpected error verifying token: %v", err)
	}
	if userID != "user-456" {
		t.Fatalf("expected user-456, got %s", userID)
	}
}

func TestJwtTokenizer_VerifyToken_TamperedToken(t *testing.T) {
	tok := newTestTokenizer("test-secret")

	tokenStr, err := tok.CreateToken("user-789")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Tamper with the token by modifying the last character
	tampered := tokenStr[:len(tokenStr)-1] + "X"

	_, err = tok.VerifyToken(tampered)
	if err == nil {
		t.Fatal("expected error for tampered token, got nil")
	}
}

func TestJwtTokenizer_VerifyToken_ExpiredToken(t *testing.T) {
	secret := "test-secret"

	// Manually create an expired token
	claims := Claims{
		UserID: "user-expired",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-25 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	tok := newTestTokenizer(secret)
	_, err = tok.VerifyToken(tokenStr)
	if err == nil {
		t.Fatal("expected error for expired token, got nil")
	}
}
