package main

import (
	"encoding/json"
	"log"
	"net/http"

	"task-flow/internal/middleware"
)

func main() {
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
