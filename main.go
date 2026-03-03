package main

import (
	"log"
	"shap-planner-backend/server"
	"shap-planner-backend/storage"
)

func main() {
	var _server = server.InitServer()

	err := storage.InitDB(_server.DatabasePath)
	if err != nil {
		log.Fatal(err)
		return
	}

	_server.Run()
}
