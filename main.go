package main

import (
	"log"
	"shap-planner-backend/server"
	"shap-planner-backend/storage"
)

func main() {
	var SERVER = server.InitServer()

	err := storage.InitDB(SERVER.DatabasePath)
	if err != nil {
		log.Fatal(err)
		return
	}

	SERVER.Run()
}

func Setup() {
	//TODO: first configuration
}
