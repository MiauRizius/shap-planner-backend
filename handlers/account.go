package handlers

import (
	"encoding/json"
	"net/http"
	"shap-planner-backend/auth"
	"shap-planner-backend/models"
	"shap-planner-backend/storage"
	"shap-planner-backend/utils"
)

func Register(w http.ResponseWriter, r *http.Request) {
	var user models.User
	_ = json.NewDecoder(r.Body).Decode(&user)
	hashed, _ := auth.HashPassword(user.Password)
	user.Password = hashed
	user.ID = utils.GenerateUUID()

	err := storage.AddUser(user)
	if err != nil {
		http.Error(w, "User exists", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func Login(w http.ResponseWriter, r *http.Request) {
	var creds struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	_ = json.NewDecoder(r.Body).Decode(&creds)

	user, err := storage.GetUserByUsername(creds.Username)
	if err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	if !auth.CheckPasswordHash(creds.Password, user.Password) {
		http.Error(w, "Wrong password", http.StatusUnauthorized)
		return
	}

	// TODO: JWT oder Session-Token erzeugen
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		return
	}
}
