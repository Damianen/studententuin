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
	UserID string `json:"string"`
	jwt.RegisteredClaims
}

func (t *JwtTokenizer) CreateToken(userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"id": userID,
			"exp": t.Clock.Now().Add(time.Hour * 24).Unix(),
		})

	tokenString, err := token.SignedString([]byte(t.SecretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (t *JwtTokenizer) VerifyToken(tokenStr string) (string, error) {
	token, err := jwt.ParseWithClaims(
		tokenStr,
		&Claims{},
		func(token *jwt.Token) (interface{}, error) {
			return t.SecretKey, nil
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
