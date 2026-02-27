package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"shap-planner-backend/auth"
	"shap-planner-backend/models"
	"shap-planner-backend/storage"
	"shap-planner-backend/utils"
)

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
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	secret := []byte(os.Getenv("SHAP_JWT_SECRET"))
	if len(secret) == 0 {
		http.Error(w, "Server misconfiguration", http.StatusInternalServerError)
		return
	}

	token, err := auth.GenerateJWT(user.ID, user.Role, secret)
	if err != nil {
		http.Error(w, "Could not generate token", http.StatusInternalServerError)
		return
	}

	type userResp struct {
		ID       string `json:"id"`
		Username string `json:"username"`
		Role     string `json:"role"`
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"token": token,
		"user": userResp{
			ID:       user.ID,
			Username: user.Username,
			Role:     user.Role,
		},
	})
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
