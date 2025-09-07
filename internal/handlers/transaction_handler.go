package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jairogloz/go-expense-tracker-back/internal/app"
	"github.com/jairogloz/go-expense-tracker-back/internal/domain"
	"go.uber.org/zap"
)

const RequestIDKey domain.ContextKey = "request_id"

// getUserID extracts the user ID from the Gin context
func getUserID(c *gin.Context) (string, error) {
	userID, exists := c.Get(string(domain.UserIDKey))
	if !exists {
		return "", fmt.Errorf("user ID not found in context")
	}

	id, ok := userID.(string)
	if !ok {
		return "", fmt.Errorf("user ID is not a string")
	}

	return id, nil
}

// TransactionHandler handles HTTP requests related to transactions
type TransactionHandler struct {
	parseInputUseCase  *app.ParseInputUseCase
	transactionService domain.TransactionService
	logger             *zap.Logger
}

// NewTransactionHandler creates a new transaction handler
func NewTransactionHandler(parseInputUseCase *app.ParseInputUseCase, transactionService domain.TransactionService, logger *zap.Logger) *TransactionHandler {
	return &TransactionHandler{
		parseInputUseCase:  parseInputUseCase,
		transactionService: transactionService,
		logger:             logger,
	}
}

// ParseInput handles the POST /parse endpoint
func (h *TransactionHandler) ParseInput(c *gin.Context) {
	startTime := time.Now()
	requestID := c.GetHeader("X-Request-ID")
	if requestID == "" {
		requestID = fmt.Sprintf("req_%d", time.Now().UnixNano())
	}

	h.logger.Info("Parse input request started",
		zap.String("request_id", requestID),
		zap.String("method", c.Request.Method),
		zap.String("path", c.Request.URL.Path),
		zap.String("user_agent", c.GetHeader("User-Agent")),
	)

	// Get user ID from context
	userID, err := getUserID(c)
	if err != nil {
		h.logger.Error("User authentication failed",
			zap.String("request_id", requestID),
			zap.Error(err),
		)
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "User authentication required",
			"details": err.Error(),
		})
		return
	}

	var request domain.ParseInputRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		h.logger.Error("Invalid request body",
			zap.String("request_id", requestID),
			zap.String("user_id", userID),
			zap.Error(err),
		)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	h.logger.Info("Parse input request validated",
		zap.String("request_id", requestID),
		zap.String("user_id", userID),
		zap.Int("input_text_length", len(request.Text)),
	)

	// Add user ID and request ID to context for the use case
	ctx := context.WithValue(c.Request.Context(), domain.UserIDKey, userID)
	ctx = context.WithValue(ctx, RequestIDKey, requestID)

	response, err := h.parseInputUseCase.Execute(ctx, request)
	if err != nil {
		h.logger.Error("Parse input use case failed",
			zap.String("request_id", requestID),
			zap.String("user_id", userID),
			zap.Error(err),
			zap.Duration("duration", time.Since(startTime)),
		)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to parse input",
			"details": err.Error(),
		})
		return
	}

	h.logger.Info("Parse input request completed successfully",
		zap.String("request_id", requestID),
		zap.String("user_id", userID),
		zap.Int("transactions_parsed", len(response.Transactions)),
		zap.Duration("duration", time.Since(startTime)),
	)

	c.JSON(http.StatusOK, response)
}

// GetTransaction handles GET /transactions/:id
func (h *TransactionHandler) GetTransaction(c *gin.Context) {
	// Get user ID from context
	userID, err := getUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "User authentication required",
			"details": err.Error(),
		})
		return
	}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid transaction ID",
		})
		return
	}

	transaction, err := h.transactionService.GetTransactionByID(c.Request.Context(), userID, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get transaction",
			"details": err.Error(),
		})
		return
	}

	if transaction == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Transaction not found",
		})
		return
	}

	c.JSON(http.StatusOK, transaction)
}

// GetTransactions handles GET /transactions
func (h *TransactionHandler) GetTransactions(c *gin.Context) {
	// Get user ID from context
	userID, err := getUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "User authentication required",
			"details": err.Error(),
		})
		return
	}

	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	transactions, err := h.transactionService.GetTransactions(c.Request.Context(), userID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get transactions",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"transactions": transactions,
		"limit":        limit,
		"offset":       offset,
	})
}

// UpdateTransaction handles PUT /transactions/:id
// UpdateTransaction handles PUT /transactions/:id
func (h *TransactionHandler) UpdateTransaction(c *gin.Context) {
	// Get user ID from context
	userID, err := getUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "User authentication required",
			"details": err.Error(),
		})
		return
	}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid transaction ID",
		})
		return
	}

	var request domain.UpdateTransactionRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	// Check if transaction exists and belongs to user
	existing, err := h.transactionService.GetTransactionByID(c.Request.Context(), userID, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get transaction",
			"details": err.Error(),
		})
		return
	}

	if existing == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Transaction not found",
		})
		return
	}

	// Create updated transaction
	transaction := &domain.Transaction{
		ID:          id,
		UserID:      userID, // Ensure user ID is set
		Amount:      request.Amount,
		Currency:    request.Currency,
		Category:    request.Category,
		Type:        request.Type,
		Date:        request.Date,
		Description: request.Description,
		SubCategory: request.SubCategory,
	}

	if err := h.transactionService.UpdateTransaction(c.Request.Context(), userID, transaction); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update transaction",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, transaction)
}

// DeleteTransaction handles DELETE /transactions/:id
func (h *TransactionHandler) DeleteTransaction(c *gin.Context) {
	// Get user ID from context
	userID, err := getUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "User authentication required",
			"details": err.Error(),
		})
		return
	}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid transaction ID",
		})
		return
	}

	if err := h.transactionService.DeleteTransaction(c.Request.Context(), userID, id); err != nil {
		if err.Error() == fmt.Sprintf("transaction with id %d not found or doesn't belong to user", id) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Transaction not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to delete transaction",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Transaction deleted successfully",
	})
}

// SetupRoutes sets up the HTTP routes
func (h *TransactionHandler) SetupRoutes(router gin.IRouter) {
	router.POST("/parse", h.ParseInput)
	router.GET("/transactions/:id", h.GetTransaction)
	router.GET("/transactions", h.GetTransactions)
	router.PUT("/transactions/:id", h.UpdateTransaction)
	router.DELETE("/transactions/:id", h.DeleteTransaction)
}
