package utils

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/argon2"
)

func HashPassword(password string) (string, error) {
	salt := make([]byte, 32)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	hashPassword := argon2.IDKey([]byte(password), []byte(salt), 1, 64*1024, 4, 32)
	encodedPassword := fmt.Sprintf("%s.%s", base64.RawStdEncoding.EncodeToString(salt),
		base64.RawStdEncoding.EncodeToString(hashPassword))

	return encodedPassword, nil
}

type Claims struct {
	Id    primitive.ObjectID `json:"id"`
	Email string             `json:"email"`
	jwt.StandardClaims
	UserType string `json:"userType"`
}

func CreateToken(id primitive.ObjectID, email string, userType string) (string, error) {
	tokenKeys := &Claims{
		Id:       id,
		Email:    email,
		UserType: userType,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodES256, tokenKeys)
	if signedToken, err := token.SignedString([]byte(os.Getenv("JWT_SECRET"))); err != nil {
		return "", nil
	} else {
		return signedToken, nil
	}
}
