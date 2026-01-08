package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"task-flow/internal/config"
	"task-flow/internal/handler"
	"task-flow/internal/middleware"
	"task-flow/internal/router"

	"github.com/joho/godotenv"
)

func main() {
	// Load .env
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	cfg := config.MustLoad()

	// Connect database
	db := config.ConnectDB()
	defer db.Close()

	// TODO: Initialize repositories
	// userRepo := repository.NewUserRepo(db)
	// refreshRepo := repository.NewRefreshTokenRepo(db)

	// TODO: Initialize services
	// jwtInstance := jwt.New([]byte(cfg.JWTSecret))
	// authSvc := authservice.NewService(userRepo, refreshRepo, jwtInstance, cfg.AccessTTL, cfg.RefreshTTL)

	// Initialize handlers
	// authHandler := handler.NewAuthHandler(authSvc)
	authHandler := &handler.AuthHandler{} // Placeholder
	taskHandler := handler.NewTaskHandler()

	// Setup router
	mux := router.New(authHandler, taskHandler)

	srv := &http.Server{
		Addr:              cfg.Addr,
		Handler:           middleware.Logger(mux),
		ReadHeaderTimeout: 5 * time.Second,
	}

	// Start server
	go func() {
		log.Printf("listening on %s", cfg.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	log.Println("shutting down...")
	_ = srv.Close()
	log.Println("shutdown complete")
}
