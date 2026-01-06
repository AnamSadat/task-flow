package main

import (
	"log"
	"net/http"

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

	// Connect database
	db := config.ConnectDB()
	defer db.Close()

	// TODO: Initialize repositories
	// userRepo := mysql.NewUserRepo(db)
	// refreshRepo := mysql.NewRefreshTokenRepo(db)

	// TODO: Initialize services
	// jwtInstance := jwt.New([]byte(os.Getenv("JWT_SECRET")))
	// authService := auth.NewService(userRepo, refreshRepo, jwtInstance, 15*time.Minute, 7*24*time.Hour)

	// Initialize handlers
	// authHandler := handler.NewAuthHandler(authService)
	authHandler := &handler.AuthHandler{} // Placeholder until repos are implemented
	taskHandler := handler.NewTaskHandler()

	// Setup router
	mux := router.New(authHandler, taskHandler)

	log.Println("server running on :3000")
	log.Fatal(http.ListenAndServe(":3000", middleware.Logger(mux)))
}
