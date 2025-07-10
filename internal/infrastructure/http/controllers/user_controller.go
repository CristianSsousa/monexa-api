package controllers

import (
	"net/http"

	"my-finance-hub-api/internal/application/interfaces"
	"my-finance-hub-api/internal/domain/entities"
	"my-finance-hub-api/internal/infrastructure/http/dto"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService interfaces.UserService
}

func NewUserController(userService interfaces.UserService) *UserController {
	return &UserController{userService: userService}
}

func (c *UserController) UpdateProfile(ctx *gin.Context) {
	// Obter ID do usuário do middleware de autenticação
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não autenticado"})
		return
	}

	var updateDTO dto.UpdateUserDTO
	if err := ctx.ShouldBindJSON(&updateDTO); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Converter DTO para entidade
	updateData := &entities.User{
		Name:     updateDTO.Name,
		Email:    updateDTO.Email,
		Password: updateDTO.Password,
	}

	// Chamar serviço para atualizar perfil
	err := c.userService.UpdateUserProfile(ctx, userID.(uint), updateData)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Perfil atualizado com sucesso"})
}
