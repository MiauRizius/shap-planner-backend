package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"shap-planner-backend/storage"
)

func GetBalance(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	userParam := query.Get("user")

	if userParam == "all" {
		// TODO: add later
	} else {
		balance, err := storage.ComputeBalance(userParam)
		if err != nil {
			log.Println("GET [api/balance] " + r.RemoteAddr + ": " + err.Error())
			http.Error(w, "Invalid request query", http.StatusBadRequest)
			return
		}
		err = json.NewEncoder(w).Encode(map[string]interface{}{
			"balance": balance,
		})
		if err != nil {
			log.Println("GET [api/balance] " + r.RemoteAddr + ": " + err.Error())
			return
		}
		log.Println("GET [api/balance] " + r.RemoteAddr + ": Successfully retrieved balance")
	}
}
