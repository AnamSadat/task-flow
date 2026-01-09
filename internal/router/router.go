package router

import (
	"net/http"

	"task-flow/internal/handler"
	"task-flow/internal/middleware"
)

type Deps struct {
	AuthHandler *handler.AuthHandler
	TaskHandler *handler.TaskHandler
	UserHandler *handler.UserHandler
	AuthMid     *middleware.AuthMiddleware
}

func New(d Deps) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"message": "Server Golang is running..."}`))
	})

	// Health check
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"OK"}`))
	})

	// Auth routes (public)
	mux.HandleFunc("POST /auth/register", d.AuthHandler.Register)
	mux.HandleFunc("POST /auth/login", d.AuthHandler.Login)
	mux.HandleFunc("POST /auth/refresh", d.AuthHandler.Refresh)
	mux.HandleFunc("POST /auth/logout", d.AuthHandler.Logout)

	// User routes (protected)
	mux.Handle("GET /users/me", middleware.RequireAccessJWT(d.AuthMid)(http.HandlerFunc(d.UserHandler.Me)))

	// Task routes
	mux.HandleFunc("GET /tasks", d.TaskHandler.GetTasks)
	mux.HandleFunc("POST /tasks", d.TaskHandler.AddTask)

	return mux
}
