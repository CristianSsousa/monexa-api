package controllers

import (
	"net/http"
	"strconv"

	"my-finance-hub-api/internal/application/interfaces"
	"my-finance-hub-api/internal/infrastructure/http/dto"
	pkgErrors "my-finance-hub-api/pkg/errors"

	"github.com/gin-gonic/gin"
)

type SavingGoalController struct {
	savingGoalService interfaces.SavingGoalService
}

func NewSavingGoalController(savingGoalService interfaces.SavingGoalService) *SavingGoalController {
	return &SavingGoalController{
		savingGoalService: savingGoalService,
	}
}

func (c *SavingGoalController) CreateSavingGoal(ctx *gin.Context) {
	userID := ctx.GetUint("user_id")

	var req dto.CreateSavingGoalRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	savingGoalEntity := req.ToEntity(userID)
	savingGoal, err := c.savingGoalService.CreateSavingGoal(ctx.Request.Context(), userID, savingGoalEntity)
	if err != nil {
		c.handleError(ctx, err)
		return
	}

	response := dto.ToSavingGoalResponse(savingGoal)
	ctx.JSON(http.StatusCreated, response)
}

func (c *SavingGoalController) GetSavingGoals(ctx *gin.Context) {
	userID := ctx.GetUint("user_id")

	savingGoals, err := c.savingGoalService.GetSavingGoalsByUser(ctx.Request.Context(), userID)
	if err != nil {
		c.handleError(ctx, err)
		return
	}

	response := dto.ToSavingGoalResponseList(savingGoals)
	ctx.JSON(http.StatusOK, response)
}

func (c *SavingGoalController) GetSavingGoal(ctx *gin.Context) {
	userID := ctx.GetUint("user_id")

	savingGoalID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	savingGoal, err := c.savingGoalService.GetSavingGoalByID(ctx.Request.Context(), userID, uint(savingGoalID))
	if err != nil {
		c.handleError(ctx, err)
		return
	}

	response := dto.ToSavingGoalResponse(savingGoal)
	ctx.JSON(http.StatusOK, response)
}

func (c *SavingGoalController) UpdateSavingGoal(ctx *gin.Context) {
	userID := ctx.GetUint("user_id")

	savingGoalID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var req dto.UpdateSavingGoalRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updates := req.ToEntity(userID)
	savingGoal, err := c.savingGoalService.UpdateSavingGoal(ctx.Request.Context(), userID, uint(savingGoalID), updates)
	if err != nil {
		c.handleError(ctx, err)
		return
	}

	response := dto.ToSavingGoalResponse(savingGoal)
	ctx.JSON(http.StatusOK, response)
}

func (c *SavingGoalController) DeleteSavingGoal(ctx *gin.Context) {
	userID := ctx.GetUint("user_id")

	savingGoalID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	err = c.savingGoalService.DeleteSavingGoal(ctx.Request.Context(), userID, uint(savingGoalID))
	if err != nil {
		c.handleError(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (c *SavingGoalController) Deposit(ctx *gin.Context) {
	userID := ctx.GetUint("user_id")

	savingGoalID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var req dto.DepositRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	savingGoal, err := c.savingGoalService.Deposit(ctx.Request.Context(), userID, uint(savingGoalID), req.Amount)
	if err != nil {
		c.handleError(ctx, err)
		return
	}

	response := dto.ToSavingGoalResponse(savingGoal)
	ctx.JSON(http.StatusOK, response)
}

func (c *SavingGoalController) handleError(ctx *gin.Context, err error) {
	if domainErr, ok := err.(pkgErrors.DomainError); ok {
		ctx.JSON(domainErr.HTTPStatus(), gin.H{"error": domainErr.Message})
		return
	}
	ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Erro interno do servidor"})
}
