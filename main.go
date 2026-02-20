package main

import (
	"log"
	"net/http"
	"shap-planner-backend/handlers"
	"shap-planner-backend/storage"
)

func main() {
	err := storage.InitDB("database.db")
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/register", handlers.Register)
	http.HandleFunc("/login", handlers.Login)

	log.Println("Server läuft auf :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
