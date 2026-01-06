package router

import (
	"net/http"

	"task-flow/internal/handler"
)

func New(authHandler *handler.AuthHandler, taskHandler *handler.TaskHandler) *http.ServeMux {
	mux := http.NewServeMux()

	// Health check
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"OK"}`))
	})

	// Auth routes
	mux.HandleFunc("POST /auth/login", authHandler.Login)
	mux.HandleFunc("POST /auth/refresh", authHandler.Refresh)
	mux.HandleFunc("POST /auth/logout", authHandler.Logout)

	// Task routes
	mux.HandleFunc("GET /tasks", taskHandler.List)
	mux.HandleFunc("POST /tasks", taskHandler.Create)

	return mux
}
