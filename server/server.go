package server

import (
	"log"
	"net/http"
	"os"
	"shap-planner-backend/auth"
	"shap-planner-backend/config"
	"shap-planner-backend/handlers"
)

type Server struct {
	Port         string
	JWTSecret    []byte
	DatabasePath string
}

func InitServer() *Server {

	err := config.CheckIfExists()
	if err != nil {
		log.Fatal(err)
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	jwtSecret := os.Getenv("SHAP_JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("SHAP_JWT_SECRET environment variable not set.")
	}
	if len(jwtSecret) < 32 {
		log.Fatal("SHAP_JWT_SECRET must be at least 32 characters long.")
	}

	return &Server{
		Port:         cfg.Port,
		JWTSecret:    []byte(jwtSecret),
		DatabasePath: cfg.DatabasePath,
	}
}

func (server *Server) Run() {
	mux := http.NewServeMux()

	mux.HandleFunc("/login", handlers.Login)

	protected := auth.AuthMiddleware(server.JWTSecret)(http.HandlerFunc(handlers.GetExpenses))
	mux.Handle("/expenses", protected)

	adminOnly := auth.AuthMiddleware(server.JWTSecret)(auth.RequireRole("admin")(http.HandlerFunc(handlers.AdminPanel)))
	mux.Handle("/admin", adminOnly)

	log.Printf("Listening on port %s", server.Port)
	log.Fatal(http.ListenAndServe(":"+server.Port, mux))
}
