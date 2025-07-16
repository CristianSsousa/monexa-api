package dto

import (
	"my-finance-hub-api/internal/domain/entities"
	"time"
)

// Request DTOs
type CreateCategoryRequest struct {
	Name  string `json:"name" binding:"required,min=2,max=100"`
	Color string `json:"color" binding:"required,hexcolor"`
	Types string `json:"type" binding:"required,oneof=expense income investment"`
}

type UpdateCategoryRequest struct {
	Name  string `json:"name" binding:"required,min=2,max=100"`
	Color string `json:"color" binding:"required,hexcolor"`
	Types string `json:"type" binding:"required,oneof=expense income investment"`
}

// Response DTOs
type CategoryResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Color     string    `json:"color"`
	UserID    uint      `json:"user_id"`
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Mappers
func ToCategoryResponse(category *entities.Category) CategoryResponse {
	return CategoryResponse{
		ID:        category.ID,
		Name:      category.Name,
		Color:     category.Color,
		UserID:    category.UserID,
		Type:      category.Type,
		CreatedAt: category.CreatedAt,
		UpdatedAt: category.UpdatedAt,
	}
}

func ToCategoryResponseList(categories []*entities.Category) []CategoryResponse {
	result := make([]CategoryResponse, len(categories))
	for i, category := range categories {
		result[i] = ToCategoryResponse(category)
	}
	return result
}

func (req *CreateCategoryRequest) ToEntity(userID uint) *entities.Category {
	return entities.NewCategory(req.Name, req.Color, userID, req.Types)
}

func (req *UpdateCategoryRequest) ToEntity(userID uint) *entities.Category {
	return entities.NewCategory(req.Name, req.Color, userID, req.Types)
}
