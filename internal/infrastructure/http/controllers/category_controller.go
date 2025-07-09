package controllers

import (
	"net/http"
	"strconv"

	"my-finance-hub-api/internal/application/interfaces"
	"my-finance-hub-api/internal/infrastructure/http/dto"
	pkgErrors "my-finance-hub-api/pkg/errors"

	"github.com/gin-gonic/gin"
)

type CategoryController struct {
	categoryService interfaces.CategoryService
}

func NewCategoryController(categoryService interfaces.CategoryService) *CategoryController {
	return &CategoryController{
		categoryService: categoryService,
	}
}

// CreateCategory godoc
// @Summary Criar nova categoria
// @Description Cria uma nova categoria para o usuário autenticado
// @Tags categorias
// @Accept json
// @Produce json
// @Param category body dto.CreateCategoryRequest true "Dados da categoria"
// @Success 201 {object} dto.CategoryResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /categories [post]
func (c *CategoryController) CreateCategory(ctx *gin.Context) {
	userID := ctx.GetUint("user_id")

	var req dto.CreateCategoryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	categoryEntity := req.ToEntity(userID)
	category, err := c.categoryService.CreateCategory(ctx.Request.Context(), userID, categoryEntity)
	if err != nil {
		c.handleError(ctx, err)
		return
	}

	response := dto.ToCategoryResponse(category)
	ctx.JSON(http.StatusCreated, response)
}

// GetCategories godoc
// @Summary Listar categorias
// @Description Lista todas as categorias do usuário autenticado
// @Tags categorias
// @Produce json
// @Success 200 {array} dto.CategoryResponse
// @Failure 401 {object} map[string]interface{}
// @Router /categories [get]
func (c *CategoryController) GetCategories(ctx *gin.Context) {
	userID := ctx.GetUint("user_id")

	categories, err := c.categoryService.GetCategoriesByUser(ctx.Request.Context(), userID)
	if err != nil {
		c.handleError(ctx, err)
		return
	}

	response := dto.ToCategoryResponseList(categories)
	ctx.JSON(http.StatusOK, response)
}

// GetCategory godoc
// @Summary Obter categoria por ID
// @Description Obtém uma categoria específica pelo ID
// @Tags categorias
// @Produce json
// @Param id path int true "ID da categoria"
// @Success 200 {object} dto.CategoryResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /categories/{id} [get]
func (c *CategoryController) GetCategory(ctx *gin.Context) {
	userID := ctx.GetUint("user_id")

	categoryID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	category, err := c.categoryService.GetCategoryByID(ctx.Request.Context(), userID, uint(categoryID))
	if err != nil {
		c.handleError(ctx, err)
		return
	}

	response := dto.ToCategoryResponse(category)
	ctx.JSON(http.StatusOK, response)
}

// UpdateCategory godoc
// @Summary Atualizar categoria
// @Description Atualiza uma categoria existente
// @Tags categorias
// @Accept json
// @Produce json
// @Param id path int true "ID da categoria"
// @Param category body dto.UpdateCategoryRequest true "Dados atualizados da categoria"
// @Success 200 {object} dto.CategoryResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /categories/{id} [put]
func (c *CategoryController) UpdateCategory(ctx *gin.Context) {
	userID := ctx.GetUint("user_id")

	categoryID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var req dto.UpdateCategoryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updates := req.ToEntity(userID)
	category, err := c.categoryService.UpdateCategory(ctx.Request.Context(), userID, uint(categoryID), updates)
	if err != nil {
		c.handleError(ctx, err)
		return
	}

	response := dto.ToCategoryResponse(category)
	ctx.JSON(http.StatusOK, response)
}

// DeleteCategory godoc
// @Summary Excluir categoria
// @Description Exclui uma categoria existente
// @Tags categorias
// @Param id path int true "ID da categoria"
// @Success 204
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /categories/{id} [delete]
func (c *CategoryController) DeleteCategory(ctx *gin.Context) {
	userID := ctx.GetUint("user_id")

	categoryID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	err = c.categoryService.DeleteCategory(ctx.Request.Context(), userID, uint(categoryID))
	if err != nil {
		c.handleError(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (c *CategoryController) handleError(ctx *gin.Context, err error) {
	if domainErr, ok := err.(pkgErrors.DomainError); ok {
		ctx.JSON(domainErr.HTTPStatus(), gin.H{"error": domainErr.Message})
		return
	}
	ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Erro interno do servidor"})
}
