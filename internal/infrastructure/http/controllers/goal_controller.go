package controllers

import (
	"net/http"
	"strconv"

	"my-finance-hub-api/internal/application/interfaces"
	"my-finance-hub-api/internal/infrastructure/http/dto"
	pkgErrors "my-finance-hub-api/pkg/errors"

	"github.com/gin-gonic/gin"
)

type GoalController struct {
	goalService interfaces.GoalService
}

func NewGoalController(goalService interfaces.GoalService) *GoalController {
	return &GoalController{
		goalService: goalService,
	}
}

func (c *GoalController) CreateGoal(ctx *gin.Context) {
	userID := ctx.GetUint("user_id")

	var req dto.CreateGoalRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	goalEntity := req.ToEntity(userID)
	goal, err := c.goalService.CreateGoal(ctx.Request.Context(), userID, goalEntity)
	if err != nil {
		c.handleError(ctx, err)
		return
	}

	response := dto.ToGoalResponse(goal)
	ctx.JSON(http.StatusCreated, response)
}

func (c *GoalController) GetGoals(ctx *gin.Context) {
	userID := ctx.GetUint("user_id")

	goals, err := c.goalService.GetGoalsByUser(ctx.Request.Context(), userID)
	if err != nil {
		c.handleError(ctx, err)
		return
	}

	response := dto.ToGoalResponseList(goals)
	ctx.JSON(http.StatusOK, response)
}

func (c *GoalController) GetGoal(ctx *gin.Context) {
	userID := ctx.GetUint("user_id")

	goalID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	goal, err := c.goalService.GetGoalByID(ctx.Request.Context(), userID, uint(goalID))
	if err != nil {
		c.handleError(ctx, err)
		return
	}

	response := dto.ToGoalResponse(goal)
	ctx.JSON(http.StatusOK, response)
}

func (c *GoalController) UpdateGoal(ctx *gin.Context) {
	userID := ctx.GetUint("user_id")

	goalID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var req dto.UpdateGoalRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updates := req.ToEntity(userID)
	goal, err := c.goalService.UpdateGoal(ctx.Request.Context(), userID, uint(goalID), updates)
	if err != nil {
		c.handleError(ctx, err)
		return
	}

	response := dto.ToGoalResponse(goal)
	ctx.JSON(http.StatusOK, response)
}

func (c *GoalController) DeleteGoal(ctx *gin.Context) {
	userID := ctx.GetUint("user_id")

	goalID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	err = c.goalService.DeleteGoal(ctx.Request.Context(), userID, uint(goalID))
	if err != nil {
		c.handleError(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (c *GoalController) handleError(ctx *gin.Context, err error) {
	if domainErr, ok := err.(pkgErrors.DomainError); ok {
		ctx.JSON(domainErr.HTTPStatus(), gin.H{"error": domainErr.Message})
		return
	}
	ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Erro interno do servidor"})
}
