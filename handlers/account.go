package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"shap-planner-backend/auth"
	"shap-planner-backend/config"
	"shap-planner-backend/models"
	"shap-planner-backend/storage"
	"shap-planner-backend/utils"
	"time"
)

var cfg, _ = config.LoadConfig()

func Register(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if user.Username == "" || user.Password == "" {
		http.Error(w, "username and password required", http.StatusBadRequest)
		return
	}

	hashed, err := auth.HashPassword(user.Password)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	user.Password = hashed
	user.ID = utils.GenerateUUID()
	user.Role = "user"

	if err := storage.AddUser(user); err != nil {
		http.Error(w, "user already exists", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
func Login(w http.ResponseWriter, r *http.Request) {
	var creds struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	user, err := storage.GetUserByUsername(creds.Username)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	if !auth.CheckPasswordHash(creds.Password, user.Password) {
		println("invalid password")
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	secret := []byte(os.Getenv("SHAP_JWT_SECRET"))
	if len(secret) == 0 {
		http.Error(w, "Server misconfiguration", http.StatusInternalServerError)
		return
	}

	accessToken, err := auth.GenerateJWT(user.ID, user.Role, secret)
	if err != nil {
		http.Error(w, "Could not generate token", http.StatusInternalServerError)
		return
	}

	refreshTokenPlain, err := utils.GenerateRefreshToken()
	if err != nil {
		http.Error(w, "could not generate refresh token", http.StatusInternalServerError)
		return
	}
	refreshHash := utils.HashToken(refreshTokenPlain)
	refreshID := utils.GenerateUUID()
	refreshExpires := time.Now().Add(7 * 24 * time.Hour).Unix() // expiry: 7 days

	deviceInfo := r.Header.Get("User-Agent")

	if err := storage.AddRefreshToken(models.RefreshToken{
		ID:         refreshID,
		UserID:     user.ID,
		Token:      refreshHash,
		ExpiresAt:  refreshExpires,
		DeviceInfo: deviceInfo,
		CreatedAt:  time.Now().Unix(),
		Revoked:    false,
	}); err != nil {
		http.Error(w, "could not save refresh token", http.StatusInternalServerError)
		return
	}

	// Return access + refresh token (refresh in plain for client to store securely)
	resp := map[string]interface{}{
		"access_token":  accessToken,
		"refresh_token": refreshTokenPlain,
		"user": map[string]interface{}{
			"id":       user.ID,
			"username": user.Username,
			"role":     user.Role,
		},
		"wgName": cfg.HouseholdName,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
func Logout(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value(auth.UserContextKey).(*auth.Claims)
	storage.RevokeAllRefreshTokensForUser(claims.UserID)
	w.WriteHeader(204)
}
func TestHandler(w http.ResponseWriter, r *http.Request) {
	claimsRaw := r.Context().Value(auth.UserContextKey)
	if claimsRaw == nil {
		http.Error(w, "No claims in context", http.StatusUnauthorized)
		return
	}

	claims, ok := claimsRaw.(*auth.Claims)
	if !ok {
		http.Error(w, "Invalid claims", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user_id": claims.UserID,
		"role":    claims.Role,
		"msg":     "access granted to protected endpoint",
	})
}
func RefreshToken(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	hashed := utils.HashToken(req.RefreshToken)

	tokenRow, err := storage.GetRefreshToken(hashed)
	if err != nil || tokenRow.Revoked || tokenRow.ExpiresAt < time.Now().Unix() {
		http.Error(w, "Invalid refresh token", http.StatusUnauthorized)
		return
	}

	if err := storage.RevokeRefreshToken(tokenRow.ID); err != nil {
		log.Println(err)
	}

	newToken, _ := utils.GenerateRefreshToken()
	newHash := utils.HashToken(newToken)
	newExpires := time.Now().Add(7 * 24 * time.Hour).Unix() //7 days
	newID := utils.GenerateUUID()
	deviceInfo := r.Header.Get("User-Agent")
	if err = storage.AddRefreshToken(models.RefreshToken{
		ID:         newID,
		UserID:     tokenRow.UserID,
		Token:      newHash,
		ExpiresAt:  newExpires,
		CreatedAt:  time.Now().Unix(),
		Revoked:    false,
		DeviceInfo: deviceInfo,
	}); err != nil {
		return
	}

	accessToken, _ := auth.GenerateJWT(tokenRow.UserID, "", []byte(os.Getenv("SHAP_JWT_SECRET")))

	if err = json.NewEncoder(w).Encode(map[string]string{
		"access_token":  accessToken,
		"refresh_token": newToken,
	}); err != nil {
		return
	}
}
