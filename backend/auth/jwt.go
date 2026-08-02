package auth

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jnie1/MTGViewer-V2/users"
)

type Claims struct {
	Role string `json:"role"`
	jwt.RegisteredClaims
}

func getTokenKey(_ *jwt.Token) (any, error) {
	tokenKey := os.Getenv("TOKEN_KEY")
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

func GererateToken(user users.UserInfo, expiresAt time.Time) (string, error) {
	userClaims := Claims{
		Role: user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			Subject:   user.Email,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, userClaims)
	tokenKey, err := getTokenKey(token)

	if err != nil {
		return "", err
	}

	return token.SignedString(tokenKey)
}
