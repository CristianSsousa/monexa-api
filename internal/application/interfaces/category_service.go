package interfaces

import (
	"context"
	"my-finance-hub-api/internal/domain/entities"
)

type CategoryService interface {
	CreateCategory(ctx context.Context, userID uint, category *entities.Category) (*entities.Category, error)
	GetCategoryByID(ctx context.Context, userID, categoryID uint) (*entities.Category, error)
	GetCategoriesByUser(ctx context.Context, userID uint) ([]*entities.Category, error)
	UpdateCategory(ctx context.Context, userID, categoryID uint, updates *entities.Category) (*entities.Category, error)
	DeleteCategory(ctx context.Context, userID, categoryID uint) error
}
