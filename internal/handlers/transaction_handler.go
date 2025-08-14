package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jairogloz/go-expense-tracker-back/internal/app"
	"github.com/jairogloz/go-expense-tracker-back/internal/domain"
)

// TransactionHandler handles HTTP requests related to transactions
type TransactionHandler struct {
	parseInputUseCase  *app.ParseInputUseCase
	transactionService domain.TransactionService
}

// NewTransactionHandler creates a new transaction handler
func NewTransactionHandler(parseInputUseCase *app.ParseInputUseCase, transactionService domain.TransactionService) *TransactionHandler {
	return &TransactionHandler{
		parseInputUseCase:  parseInputUseCase,
		transactionService: transactionService,
	}
}

// ParseInput handles the POST /parse endpoint
func (h *TransactionHandler) ParseInput(c *gin.Context) {
	var request domain.ParseInputRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	response, err := h.parseInputUseCase.Execute(c.Request.Context(), request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to parse input",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetTransaction handles GET /transactions/:id
func (h *TransactionHandler) GetTransaction(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid transaction ID",
		})
		return
	}

	transaction, err := h.transactionService.GetTransactionByID(c.Request.Context(), id)
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

	transactions, err := h.transactionService.GetTransactions(c.Request.Context(), limit, offset)
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

// SetupRoutes sets up the HTTP routes
func (h *TransactionHandler) SetupRoutes(router *gin.Engine) {
	router.POST("/parse", h.ParseInput)
	router.GET("/transactions/:id", h.GetTransaction)
	router.GET("/transactions", h.GetTransactions)
}
