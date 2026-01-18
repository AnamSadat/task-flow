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
	"task-flow/internal/pkg/jwt"
	"task-flow/internal/repository/mysql"
	"task-flow/internal/router"
	"task-flow/internal/service"
	authservice "task-flow/internal/service/auth"

	"github.com/joho/godotenv"
	"github.com/rs/cors"
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

	// Initialize repositories
	userRepo := mysql.NewUserRepo(db)
	refreshRepo := mysql.NewRefreshTokenRepo(db)
	taskRepo := mysql.NewTaskRepo(db)

	// Initialize JWT
	jwtInstance := jwt.New([]byte(cfg.JWTSecret))

	// Initialize middleware
	authMid := middleware.NewAuthMiddleware(jwtInstance)

	// Initialize services
	authSvc := authservice.NewService(userRepo, refreshRepo, jwtInstance, cfg.AccessTTL, cfg.RefreshTTL)
	taskSvc := service.NewServiceTask(taskRepo)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authSvc)
	userHandler := handler.NewUserHandler(userRepo)
	taskHandler := handler.NewTaskHandler(taskSvc)

	// Setup router
	mux := router.New(router.Deps{
		AuthHandler: authHandler,
		TaskHandler: taskHandler,
		UserHandler: userHandler,
		AuthMid:     authMid,
	})

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	srv := &http.Server{
		Addr: cfg.Addr,
		// Handler: middleware.CORS(middleware.Logger(mux)), // If using manual CORS
		Handler:           c.Handler(middleware.Logger(mux)),
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
