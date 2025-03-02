package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"customize_crm/controller"
	_ "customize_crm/docs"
	"customize_crm/middleware"
	"customize_crm/service"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title CRM API
// @version 1.0
// @description API for Customize CRM
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.example.com/support
// @contact.email support@example.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and the JWT token.

func main() {
	// 1. Load environment variables
	loadEnvFile()

	// 2. Connect to database
	dbPool := connectToDatabase()
	defer dbPool.Close()

	// 3. Configure and start server
	startServer(dbPool)
}

// loadEnvFile loads variables from .env file
func loadEnvFile() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found or could not be loaded: %v\n", err)
	}
}

func connectToDatabase() *pgxpool.Pool {
	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_SSL_MODE"),
	)

	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		log.Fatalf("Unable to parse connection string: %v", err)
	}

	maxConns, err := strconv.Atoi(os.Getenv("DB_MAX_CONNECTIONS"))
	if err != nil {
		log.Fatalf("Invalid DB_MAX_CONNECTIONS value: %v", err)
	}

	config.MaxConns = int32(maxConns)

	dbPool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}

	if err := dbPool.Ping(context.Background()); err != nil {
		log.Fatalf("Unable to ping database: %v", err)
	}

	log.Println("Successfully connected to database")
	return dbPool
}

func startServer(dbPool *pgxpool.Pool) {
	//  services
	userService := service.NewUserService(dbPool)
	authService := service.NewAuthService(userService)

	// controllers
	authController := controller.NewAuthController(authService)
	userController := controller.NewUserController(userService)

	router := setupRouter()

	router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))

	setupAuthRoutes(router, authController, userService)
	setupUserRoutes(router, userController, userService)

	port := getEnv("SERVER_PORT", "8080")
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		log.Printf("Starting server on port %s", port)
		log.Printf("Swagger documentation available at http://localhost:%s/swagger/index.html", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	waitForShutdownSignal(server)
}

func setupRouter() *chi.Mux {
	router := chi.NewRouter()

	// Global middleware
	router.Use(chimiddleware.Logger)
	router.Use(chimiddleware.Recoverer)
	router.Use(chimiddleware.RequestID)
	router.Use(chimiddleware.RealIP)
	router.Use(chimiddleware.Timeout(60 * time.Second))
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	return router
}

func setupAuthRoutes(router *chi.Mux, controller *controller.AuthController, userService *service.UserService) {
	// Public auth
	router.Route("/api/v1/auth", func(r chi.Router) {
		r.Post("/login", controller.Login)
		r.Post("/refresh-token", controller.RefreshToken)
		r.Post("/forgot-password", controller.ForgotPassword)
		r.Post("/reset-password", controller.ResetPassword)

		// Protected auth
		r.Group(func(r chi.Router) {
			authMiddleware := middleware.NewAuthMiddleware(userService)
			r.Use(authMiddleware.Authenticate)
			r.Post("/logout", controller.Logout)
		})
	})
}

func setupUserRoutes(router *chi.Mux, controller *controller.UserController, userService *service.UserService) {
	router.Route("/api/v1/users", func(r chi.Router) {
		authMiddleware := middleware.NewAuthMiddleware(userService)
		r.Use(authMiddleware.Authenticate)

		r.Get("/me", controller.GetCurrentUser)
		r.Patch("/me", controller.UpdateCurrentUser)

		r.Group(func(r chi.Router) {
			r.Use(authMiddleware.RequireAdmin)
			r.Get("/", controller.GetAllUsers)
			r.Post("/", controller.CreateUser)
			r.Get("/{id}", controller.GetUserByID)
			r.Patch("/{id}", controller.UpdateUser)
			r.Delete("/", controller.DeleteUsers)
		})
	})
}

func waitForShutdownSignal(server *http.Server) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Error during server shutdown: %v\n", err)
		if err := server.Close(); err != nil {
			log.Printf("Could not close server: %v\n", err)
		}
	}

	log.Println("Server stopped successfully")
}

func getEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if exists {
		return value
	}
	return defaultValue
}
