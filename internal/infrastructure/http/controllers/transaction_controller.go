package controllers

import (
	"net/http"
	"strconv"
	"time"

	"my-finance-hub-api/internal/application/interfaces"
	"my-finance-hub-api/internal/domain/repositories"
	"my-finance-hub-api/internal/infrastructure/http/dto"
	pkgErrors "my-finance-hub-api/pkg/errors"

	"github.com/gin-gonic/gin"
)

type TransactionController struct {
	transactionService interfaces.TransactionService
}

func NewTransactionController(transactionService interfaces.TransactionService) *TransactionController {
	return &TransactionController{
		transactionService: transactionService,
	}
}

func (c *TransactionController) CreateTransaction(ctx *gin.Context) {
	userID := ctx.GetUint("user_id")

	var req dto.CreateTransactionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	transactionEntity := req.ToEntity(userID)
	transaction, err := c.transactionService.CreateTransaction(ctx.Request.Context(), userID, transactionEntity)
	if err != nil {
		c.handleError(ctx, err)
		return
	}

	response := dto.ToTransactionResponse(transaction)
	ctx.JSON(http.StatusCreated, response)
}

func (c *TransactionController) GetTransactions(ctx *gin.Context) {
	userID := ctx.GetUint("user_id")

	// Parse filters from query parameters
	var filters dto.TransactionFiltersRequest
	if err := ctx.ShouldBindQuery(&filters); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convert DTO filters to repository filters
	repoFilters := &repositories.TransactionFilters{
		UserID:     userID,
		Paid:       filters.Paid,
		Month:      filters.Month,
		Year:       filters.Year,
		Type:       filters.Type,
		CategoryID: filters.CategoryID,
	}

	transactions, err := c.transactionService.GetTransactionsByUser(ctx.Request.Context(), userID, repoFilters)
	if err != nil {
		c.handleError(ctx, err)
		return
	}

	response := dto.ToTransactionResponseList(transactions)
	ctx.JSON(http.StatusOK, response)
}

func (c *TransactionController) GetTransaction(ctx *gin.Context) {
	userID := ctx.GetUint("user_id")

	transactionID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	transaction, err := c.transactionService.GetTransactionByID(ctx.Request.Context(), userID, uint(transactionID))
	if err != nil {
		c.handleError(ctx, err)
		return
	}

	response := dto.ToTransactionResponse(transaction)
	ctx.JSON(http.StatusOK, response)
}

func (c *TransactionController) UpdateTransaction(ctx *gin.Context) {
	userID := ctx.GetUint("user_id")

	transactionID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var req dto.UpdateTransactionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updates := req.ToEntity(userID)
	transaction, err := c.transactionService.UpdateTransaction(ctx.Request.Context(), userID, uint(transactionID), updates)
	if err != nil {
		c.handleError(ctx, err)
		return
	}

	response := dto.ToTransactionResponse(transaction)
	ctx.JSON(http.StatusOK, response)
}

func (c *TransactionController) DeleteTransaction(ctx *gin.Context) {
	userID := ctx.GetUint("user_id")

	transactionID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	err = c.transactionService.DeleteTransaction(ctx.Request.Context(), userID, uint(transactionID))
	if err != nil {
		c.handleError(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (c *TransactionController) TogglePaid(ctx *gin.Context) {
	userID := ctx.GetUint("user_id")

	transactionID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	transaction, err := c.transactionService.TogglePaidStatus(ctx.Request.Context(), userID, uint(transactionID))
	if err != nil {
		c.handleError(ctx, err)
		return
	}

	response := dto.ToTransactionResponse(transaction)
	ctx.JSON(http.StatusOK, response)
}

func (c *TransactionController) GetReports(ctx *gin.Context) {
	userID := ctx.GetUint("user_id")

	// Obter datas do frontend
	startDateStr := ctx.Query("startDate")
	endDateStr := ctx.Query("endDate")

	// Converter datas para o formato correto
	var startDate, endDate time.Time
	var err error

	if startDateStr != "" {
		startDate, err = time.Parse("2006-01-02", startDateStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Data inicial inválida"})
			return
		}
	}

	if endDateStr != "" {
		endDate, err = time.Parse("2006-01-02", endDateStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Data final inválida"})
			return
		}
	}

	// Criar filtros com as datas
	repoFilters := &repositories.TransactionFilters{
		UserID:    userID,
		StartDate: startDate,
		EndDate:   endDate,
	}

	reports, err := c.transactionService.GetReports(ctx.Request.Context(), userID, repoFilters)
	if err != nil {
		c.handleError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, reports)
}

func (c *TransactionController) GetDashboardReports(ctx *gin.Context) {
	userID := ctx.GetUint("user_id")

	reports, err := c.transactionService.GetDashboardReports(ctx.Request.Context(), userID)
	if err != nil {
		c.handleError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, reports)
}

func (c *TransactionController) handleError(ctx *gin.Context, err error) {
	if domainErr, ok := err.(pkgErrors.DomainError); ok {
		ctx.JSON(domainErr.HTTPStatus(), gin.H{"error": domainErr.Message})
		return
	}
	ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Erro interno do servidor"})
}
