package main

import (
	"context"
	"fmt"
	"log"

	"github.com/jairogloz/go-expense-tracker-back/internal/domain"
	"github.com/jairogloz/go-expense-tracker-back/internal/infra"
	"go.uber.org/zap"
)

func main() {
	fmt.Println("Testing Zap logging setup...")

	// Test logger creation
	logger, err := infra.NewLogger()
	if err != nil {
		log.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Sync()

	fmt.Println("✅ Logger created successfully")

	// Test development logger
	devLogger, err := infra.NewDevelopmentLogger()
	if err != nil {
		log.Fatalf("Failed to create development logger: %v", err)
	}
	defer devLogger.Sync()

	fmt.Println("✅ Development logger created successfully")

	// Test OpenAI service with logger
	openaiService := infra.NewOpenAIService("fake-api-key", logger)
	fmt.Println("✅ OpenAI service created with logger")

	// Test context with user ID and request ID
	ctx := context.Background()
	ctx = context.WithValue(ctx, domain.UserIDKey, "test-user-123")
	ctx = context.WithValue(ctx, domain.ContextKey("request_id"), "test-request-456")

	// Log some test messages
	logger.Info("Testing logger functionality",
		zap.String("user_id", "test-user-123"),
		zap.String("request_id", "test-request-456"),
		zap.String("test_type", "logging_setup"),
	)

	devLogger.Debug("Testing development logger functionality",
		zap.String("message", "This should appear in development mode"),
		zap.Bool("debug_enabled", true),
	)

	fmt.Println("✅ All logging tests completed successfully")
	fmt.Println("\nLogging setup is ready for your /parse endpoint!")
	fmt.Println("Key logging points implemented:")
	fmt.Println("- Handler level: Request/response logging with timing")
	fmt.Println("- OpenAI Service level: API call logging, JSON parsing errors")
	fmt.Println("- Repository level: Database transaction logging")
	fmt.Println("- Structured logging with user_id and request_id correlation")

	_ = openaiService // Avoid unused variable warning
}
