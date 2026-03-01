package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"net/http"
	"shap-planner-backend/auth"

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
func GenerateRefreshToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}
func HashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}
func IsLoggedIn(w http.ResponseWriter, r *http.Request) (*auth.Claims, bool) {
	claimsRaw := r.Context().Value(auth.UserContextKey)
	if claimsRaw == nil {
		http.Error(w, "No claims in context", http.StatusUnauthorized)
		return nil, false
	}

	claims, ok := claimsRaw.(*auth.Claims)
	if !ok {
		http.Error(w, "Invalid claims", http.StatusUnauthorized)
		return nil, false
	}
	return claims, true
}
