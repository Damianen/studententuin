package auth

import (
	"api/internal/infra/utils"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JwtTokenizer struct {
	Clock utils.SystemClock
	SecretKey string
}

type Claims struct {
	UserID string `json:"id"`
	jwt.RegisteredClaims
}

func (t *JwtTokenizer) CreateToken(userID string) (string, error) {
	now := t.Clock.Now()

	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(24 * time.Hour)),
			IssuedAt: jwt.NewNumericDate(now),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(t.SecretKey))
}

func (t *JwtTokenizer) VerifyToken(tokenStr string) (string, error) {
	token, err := jwt.ParseWithClaims(
		tokenStr,
		&Claims{},
		func(token *jwt.Token) (any, error) {
			return []byte(t.SecretKey), nil
		},
	)
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return "", errors.New("user id missing in token")
	}

	return claims.UserID, nil
}
