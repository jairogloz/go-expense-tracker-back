package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jairogloz/go-expense-tracker-back/internal/domain"
)

// AccountHandler handles HTTP requests related to accounts
type AccountHandler struct {
	accountService domain.AccountService
}

// NewAccountHandler creates a new account handler
func NewAccountHandler(accountService domain.AccountService) *AccountHandler {
	return &AccountHandler{
		accountService: accountService,
	}
}

// CreateAccount handles POST /accounts
func (h *AccountHandler) CreateAccount(c *gin.Context) {
	// Get user ID from context
	userID, err := getUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "User authentication required",
			"details": err.Error(),
		})
		return
	}

	var request domain.CreateAccountRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	account, err := h.accountService.CreateAccount(c.Request.Context(), userID, request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Failed to create account",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, account)
}

// GetAccounts handles GET /accounts
func (h *AccountHandler) GetAccounts(c *gin.Context) {
	// Get user ID from context
	userID, err := getUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "User authentication required",
			"details": err.Error(),
		})
		return
	}

	accounts, err := h.accountService.GetUserAccounts(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get accounts",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"accounts": accounts,
	})
}

// GetAccount handles GET /accounts/:id
func (h *AccountHandler) GetAccount(c *gin.Context) {
	// Get user ID from context
	userID, err := getUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "User authentication required",
			"details": err.Error(),
		})
		return
	}

	accountID := c.Param("id")
	if accountID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Account ID is required",
		})
		return
	}

	account, err := h.accountService.GetAccount(c.Request.Context(), userID, accountID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Account not found",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, account)
}

// GetAccountBalance handles GET /accounts/:id/balance
func (h *AccountHandler) GetAccountBalance(c *gin.Context) {
	// Get user ID from context
	userID, err := getUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "User authentication required",
			"details": err.Error(),
		})
		return
	}

	accountID := c.Param("id")
	if accountID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Account ID is required",
		})
		return
	}

	summary, err := h.accountService.GetAccountBalance(c.Request.Context(), userID, accountID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Failed to get account balance",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, summary)
}

// UpdateAccount handles PUT /accounts/:id
func (h *AccountHandler) UpdateAccount(c *gin.Context) {
	// Get user ID from context
	userID, err := getUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "User authentication required",
			"details": err.Error(),
		})
		return
	}

	accountID := c.Param("id")
	if accountID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Account ID is required",
		})
		return
	}

	var request domain.UpdateAccountRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	account, err := h.accountService.UpdateAccount(c.Request.Context(), userID, accountID, request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Failed to update account",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, account)
}

// DeleteAccount handles DELETE /accounts/:id
func (h *AccountHandler) DeleteAccount(c *gin.Context) {
	// Get user ID from context
	userID, err := getUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "User authentication required",
			"details": err.Error(),
		})
		return
	}

	accountID := c.Param("id")
	if accountID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Account ID is required",
		})
		return
	}

	if err := h.accountService.DeleteAccount(c.Request.Context(), userID, accountID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Failed to delete account",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Account deleted successfully",
	})
}

// SetDefaultAccount handles PUT /accounts/:id/set-default
func (h *AccountHandler) SetDefaultAccount(c *gin.Context) {
	// Get user ID from context
	userID, err := getUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "User authentication required",
			"details": err.Error(),
		})
		return
	}

	accountID := c.Param("id")
	if accountID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Account ID is required",
		})
		return
	}

	if err := h.accountService.SetDefaultAccount(c.Request.Context(), userID, accountID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Failed to set default account",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Default account set successfully",
	})
}

// SetupAccountRoutes sets up the HTTP routes for accounts
func (h *AccountHandler) SetupAccountRoutes(router gin.IRouter) {
	router.POST("/accounts", h.CreateAccount)
	router.GET("/accounts", h.GetAccounts)
	router.GET("/accounts/:id", h.GetAccount)
	router.GET("/accounts/:id/balance", h.GetAccountBalance)
	router.PUT("/accounts/:id", h.UpdateAccount)
	router.DELETE("/accounts/:id", h.DeleteAccount)
	router.PUT("/accounts/:id/set-default", h.SetDefaultAccount)
}
