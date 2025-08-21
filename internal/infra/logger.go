package infra

import (
	"go.uber.org/zap"
)

// NewLogger creates a new zap logger instance
func NewLogger() (*zap.Logger, error) {
	config := zap.NewProductionConfig()

	// Customize configuration
	config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	config.Development = false
	config.DisableCaller = false
	config.DisableStacktrace = false
	config.Sampling = nil // Disable sampling for now

	// Build the logger
	logger, err := config.Build()
	if err != nil {
		return nil, err
	}

	return logger, nil
}

// NewDevelopmentLogger creates a logger optimized for development
func NewDevelopmentLogger() (*zap.Logger, error) {
	config := zap.NewDevelopmentConfig()
	config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)

	logger, err := config.Build()
	if err != nil {
		return nil, err
	}

	return logger, nil
}

// NewLoggerForEnvironment creates a logger based on environment
func NewLoggerForEnvironment(isDevelopment bool) (*zap.Logger, error) {
	if isDevelopment {
		return NewDevelopmentLogger()
	}
	return NewLogger()
}
