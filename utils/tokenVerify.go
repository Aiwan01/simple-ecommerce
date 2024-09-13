package utils

import (
	"errors"
	"os"

	"github.com/golang-jwt/jwt"
)

func VerifyUserToken(token string) (string, string, string, error) {
	if token == "" {
		return "", "", "", errors.New("token is empty")
	}
	claims := &Claims{}

	parsedToken, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return "", "", "", errors.New("signature is invalid")
		}
		return "", "", "", errors.New("token is invalid")
	}
	if !parsedToken.Valid {
		return "", "", "", errors.New("token is not valid")
	}
	if claims == nil {
		return "", "", "", errors.New("token claim is nil")
	}

	return claims.Id.Hex(), claims.Email, claims.UserType, nil
}
