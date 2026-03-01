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
		log.Println("[api/register] " + r.RemoteAddr + ": " + err.Error())
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if user.Username == "" || user.Password == "" {
		log.Println("[api/register] " + r.RemoteAddr + ": Username or Password is empty")
		http.Error(w, "username and password required", http.StatusBadRequest)
		return
	}

	hashed, err := auth.HashPassword(user.Password)
	if err != nil {
		log.Println("[api/register] " + r.RemoteAddr + ": " + err.Error())
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	user.Password = hashed
	user.ID = utils.GenerateUUID()
	user.Role = "user"

	if err := storage.AddUser(&user); err != nil {
		log.Println("[api/register] " + r.RemoteAddr + ": " + err.Error())
		http.Error(w, "user already exists", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	log.Println("[api/register] " + r.RemoteAddr + ": Successfully created user")
}
func Login(w http.ResponseWriter, r *http.Request) {
	var creds struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		log.Println("[api/login] " + r.RemoteAddr + ": " + err.Error())
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	user, err := storage.GetUserByUsername(creds.Username)
	if err != nil {
		log.Println("[api/login] " + r.RemoteAddr + ": " + err.Error())
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	if !auth.CheckPasswordHash(creds.Password, user.Password) {
		log.Println("[api/login] " + r.RemoteAddr + ": Invalid credentials")
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	secret := []byte(os.Getenv("SHAP_JWT_SECRET"))
	if len(secret) == 0 {
		log.Println("[api/login] " + r.RemoteAddr + ": Server misconfiguration")
		http.Error(w, "Server misconfiguration", http.StatusInternalServerError)
		return
	}

	accessToken, err := auth.GenerateJWT(user.ID, user.Role, secret)
	if err != nil {
		log.Println("[api/login] " + r.RemoteAddr + ": " + err.Error())
		http.Error(w, "Could not generate token", http.StatusInternalServerError)
		return
	}

	refreshTokenPlain, err := utils.GenerateRefreshToken()
	if err != nil {
		log.Println("[api/login] " + r.RemoteAddr + ": " + err.Error())
		http.Error(w, "could not generate refresh token", http.StatusInternalServerError)
		return
	}
	refreshHash := utils.HashToken(refreshTokenPlain)
	refreshID := utils.GenerateUUID()
	refreshExpires := time.Now().Add(7 * 24 * time.Hour).Unix() // expiry: 7 days

	deviceInfo := r.Header.Get("User-Agent")

	if err := storage.AddRefreshToken(&models.RefreshToken{
		ID:         refreshID,
		UserID:     user.ID,
		Token:      refreshHash,
		ExpiresAt:  refreshExpires,
		DeviceInfo: deviceInfo,
		CreatedAt:  time.Now().Unix(),
		Revoked:    false,
	}); err != nil {
		log.Println("[api/login] " + r.RemoteAddr + ": " + err.Error())
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
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Println("[api/login] " + r.RemoteAddr + ": " + err.Error())
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	log.Println("[api/login] " + r.RemoteAddr + ": Successfully logged in")
}
func Logout(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value(auth.UserContextKey).(*auth.Claims)
	err := storage.RevokeAllRefreshTokensForUser(claims.UserID)
	if err != nil {
		log.Println("[api/logout] " + r.RemoteAddr + ": " + err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(204)
}
func TestHandler(w http.ResponseWriter, r *http.Request) {
	claims, _ := utils.IsLoggedIn(w, r)

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(map[string]interface{}{
		"user_id": claims.UserID,
		"role":    claims.Role,
		"msg":     "access granted to protected endpoint",
	})
	if err != nil {
		log.Println("[api/ping] " + r.RemoteAddr + ": " + err.Error())
		return
	}
	log.Println("[api/login] " + r.RemoteAddr + ": Successfully tested connection")
}
func RefreshToken(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println("[api/refresh] " + r.RemoteAddr + ": " + err.Error())
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	hashed := utils.HashToken(req.RefreshToken)

	tokenRow, err := storage.GetRefreshToken(hashed)
	if err != nil || tokenRow.Revoked || tokenRow.ExpiresAt < time.Now().Unix() {
		log.Println("[api/refresh] " + r.RemoteAddr + ": Invalid refresh token")
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
	if err = storage.AddRefreshToken(&models.RefreshToken{
		ID:         newID,
		UserID:     tokenRow.UserID,
		Token:      newHash,
		ExpiresAt:  newExpires,
		CreatedAt:  time.Now().Unix(),
		Revoked:    false,
		DeviceInfo: deviceInfo,
	}); err != nil {
		log.Println("[api/refresh] " + r.RemoteAddr + ": " + err.Error())
		http.Error(w, "Could not generate new refresh token", http.StatusInternalServerError)
		return
	}

	accessToken, _ := auth.GenerateJWT(tokenRow.UserID, "", []byte(os.Getenv("SHAP_JWT_SECRET")))

	if err = json.NewEncoder(w).Encode(map[string]string{
		"access_token":  accessToken,
		"refresh_token": newToken,
	}); err != nil {
		log.Println("[api/refresh] " + r.RemoteAddr + ": " + err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
