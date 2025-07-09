package controllers

import (
	"net/http"

	"my-finance-hub-api/internal/application/interfaces"
	"my-finance-hub-api/internal/infrastructure/http/dto"
	pkgErrors "my-finance-hub-api/pkg/errors"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authService interfaces.AuthService
}

func NewAuthController(authService interfaces.AuthService) *AuthController {
	return &AuthController{
		authService: authService,
	}
}

func (c *AuthController) Register(ctx *gin.Context) {
	var req dto.RegisterUserRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := c.authService.Register(ctx.Request.Context(), req.Name, req.Email, req.Password)
	if err != nil {
		c.handleError(ctx, err)
		return
	}

	response := dto.ToUserResponse(user)
	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Usuário registrado com sucesso",
		"user":    response,
	})
}

func (c *AuthController) Login(ctx *gin.Context) {
	var req dto.LoginRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, token, err := c.authService.Login(ctx.Request.Context(), req.Email, req.Password)
	if err != nil {
		c.handleError(ctx, err)
		return
	}

	response := dto.AuthResponse{
		Token: token,
		User:  dto.ToUserResponse(user),
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *AuthController) GetCurrentUser(ctx *gin.Context) {
	userID := c.getUserIDFromContext(ctx)

	user, err := c.authService.GetUserByID(ctx.Request.Context(), userID)
	if err != nil {
		c.handleError(ctx, err)
		return
	}

	response := dto.ToUserResponse(user)
	ctx.JSON(http.StatusOK, response)
}

func (c *AuthController) UpdateUser(ctx *gin.Context) {
	userID := c.getUserIDFromContext(ctx)

	var req dto.UpdateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := c.authService.UpdateUser(ctx.Request.Context(), userID, req.Name, req.Email)
	if err != nil {
		c.handleError(ctx, err)
		return
	}

	response := dto.ToUserResponse(user)
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Usuário atualizado com sucesso",
		"user":    response,
	})
}

func (c *AuthController) ChangePassword(ctx *gin.Context) {
	userID := c.getUserIDFromContext(ctx)

	var req dto.ChangePasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := c.authService.ChangePassword(ctx.Request.Context(), userID, req.CurrentPassword, req.NewPassword)
	if err != nil {
		c.handleError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Senha alterada com sucesso"})
}

func (c *AuthController) Logout(ctx *gin.Context) {
	// Em uma implementação real, você poderia invalidar o token aqui
	// Por enquanto, apenas retorna sucesso
	ctx.JSON(http.StatusOK, gin.H{"message": "Logout realizado com sucesso"})
}

// Helpers
func (c *AuthController) getUserIDFromContext(ctx *gin.Context) uint {
	userID, exists := ctx.Get("user_id")
	if !exists {
		// Isso não deveria acontecer se o middleware estiver funcionando corretamente
		panic("user_id não encontrado no contexto")
	}
	return userID.(uint)
}

func (c *AuthController) handleError(ctx *gin.Context, err error) {
	if domainErr, ok := err.(pkgErrors.DomainError); ok {
		ctx.JSON(domainErr.HTTPStatus(), gin.H{"error": domainErr.Message})
		return
	}

	// Erro genérico
	ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Erro interno do servidor"})
}
