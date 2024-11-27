package auth

import (
	"errors"
	"os"

	"github.com/golang-jwt/jwt"
)

type Claims struct {
	jwt.StandardClaims
	Role string `json:"role"`
}

var tokenKey = os.Getenv("TOKEN_KEY")

func getTokenKey(_ *jwt.Token) (interface{}, error) {
	return []byte(tokenKey), nil
}

func ParseToken(tokenString string) (claims *Claims, err error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, getTokenKey)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)

	if !ok {
		return nil, errors.New("claims failed to parse")
	}

	return claims, nil
}

func GenerateToken(claims *Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	tokenKey, err := getTokenKey(token)

	if err != nil {
		return "", err
	}

	return token.SignedString(tokenKey)
}
