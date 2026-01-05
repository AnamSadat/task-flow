package main

import (
	"encoding/json"
	"log"
	"net/http"

	"task-flow/internal/config"
	"task-flow/internal/middleware"

	"github.com/joho/godotenv"
)

func main() {
	// Load .env
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Connect database
	db := config.ConnectDB()
	defer db.Close()

	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "OK",
			"message": "Server is running!",
		})
	})

	log.Println("server running on :3000")
	log.Fatal(http.ListenAndServe(":3000", middleware.Logger(mux)))
}
