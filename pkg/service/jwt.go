package service

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/renlin-code/mock-shop-api/pkg/errors_handler"
)

const confirmationTokenTTL = time.Hour
const signInTokenTTL = 24 * time.Hour

func generateToken(payload map[string]interface{}, signKey string, tokenTTL time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"payload": payload,
		"exp":     time.Now().Add(tokenTTL).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(signKey))
}

func parseToken(tokenString, signKey string) (map[string]interface{}, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(signKey), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		payload := claims["payload"].(map[string]interface{})
		return payload, nil
	}

	return nil, errors_handler.Forbidden("invalid token")
}
