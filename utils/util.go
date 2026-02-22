package utils

import (
	"crypto/rand"
	"encoding/base64"

	"github.com/google/uuid"
)

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
