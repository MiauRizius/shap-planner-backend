package utils

import (
	"os"
	"time"

	"crypto/rand"
	"encoding/base64"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var secret = []byte(os.Getenv("SHAP_JWT_SECRET"))

// Password Validation
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// JWT Validation
func GenerateJWT(userID string) (string, error) {
	claims := jwt.MapClaims{
		"userId": userID,
		"exp":    time.Now().Add(time.Hour * 1).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}
func ValidateJWT(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})
}

func GenerateUUID() string {
	return uuid.New().String()
}
func GenerateSecret() string {
	b := make([]byte, 64)
	_, err := rand.Read(b)
	if err != nil {
		return err.Error()
	}
	return base64.StdEncoding.EncodeToString(b)
}
