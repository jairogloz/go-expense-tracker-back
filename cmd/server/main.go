package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jairogloz/go-expense-tracker-back/config"
	"github.com/jairogloz/go-expense-tracker-back/internal/app"
	"github.com/jairogloz/go-expense-tracker-back/internal/handlers"
	"github.com/jairogloz/go-expense-tracker-back/internal/infra"
	"github.com/jairogloz/go-expense-tracker-back/internal/services"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Create context with timeout for db connection
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// Initialize database connection
	db, err := infra.NewDatabaseConnection(ctx, cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize repository
	transactionRepo := infra.NewPostgreSQLTransactionRepository(db)

	// Use background context for the rest of the operations
	ctx = context.Background()

	// Create tables if they don't exist
	if err := transactionRepo.CreateTransactionsTable(ctx); err != nil {
		log.Fatalf("Failed to create database tables: %v", err)
	}

	// Initialize services
	aiService := infra.NewOpenAIService(cfg.OpenAI.APIKey)
	transactionService := services.NewTransactionService(transactionRepo)

	// Initialize use cases
	parseInputUseCase := app.NewParseInputUseCase(aiService, transactionService)

	// Initialize auth service
	authService := infra.NewSupabaseAuthService(cfg)

	// Initialize handlers
	transactionHandler := handlers.NewTransactionHandler(parseInputUseCase, transactionService)
	authMiddleware := handlers.NewAuthMiddleware(authService)

	// Setup routes
	r := gin.Default()

	// Middleware
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Health check endpoint (no auth required)
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
			"time":   time.Now().UTC(),
		})
	})

	// Protected routes group
	protected := r.Group("/")
	protected.Use(authMiddleware.Authenticate())

	// Setup routes with authentication
	transactionHandler.SetupRoutes(protected)

	// Create HTTP server
	srv := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: r,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Starting server on port %s", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Give outstanding requests 30 seconds to complete
	ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")
}
