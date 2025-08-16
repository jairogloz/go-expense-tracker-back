package infra

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jairogloz/go-expense-tracker-back/config"
	"github.com/jairogloz/go-expense-tracker-back/internal/domain"
)

// SupabaseAuthService implements the AuthService interface
type SupabaseAuthService struct {
	jwtSecret string
}

// NewSupabaseAuthService creates a new Supabase authentication service
func NewSupabaseAuthService(cfg *config.Config) *SupabaseAuthService {
	return &SupabaseAuthService{
		jwtSecret: cfg.Supabase.JWTSecret,
	}
}

// SupabaseClaims represents the custom claims in a Supabase JWT token
type SupabaseClaims struct {
	jwt.RegisteredClaims
	Email string `json:"email"`
	Role  string `json:"role"`
}

// ValidateToken validates a Supabase JWT token and returns user information
func (s *SupabaseAuthService) ValidateToken(ctx context.Context, tokenString string) (*domain.AuthUser, error) {
	// Remove Bearer prefix if present
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	// Parse and validate the token
	token, err := jwt.ParseWithClaims(tokenString, &SupabaseClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	// Check if token is valid
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	// Extract claims
	claims, ok := token.Claims.(*SupabaseClaims)
	if !ok {
		return nil, fmt.Errorf("failed to parse token claims")
	}

	// Check if token is expired using time comparison
	if claims.ExpiresAt != nil && claims.ExpiresAt.Time.Before(time.Now()) {
		return nil, fmt.Errorf("token is expired")
	}

	// Check if token is issued in the future
	if claims.IssuedAt != nil && claims.IssuedAt.Time.After(time.Now()) {
		return nil, fmt.Errorf("token used before issued")
	}

	// Check if token is not valid yet
	if claims.NotBefore != nil && claims.NotBefore.Time.After(time.Now()) {
		return nil, fmt.Errorf("token used before valid")
	}

	// Create and return user
	user := &domain.AuthUser{
		ID:    claims.Subject,
		Email: claims.Email,
	}

	return user, nil
}
