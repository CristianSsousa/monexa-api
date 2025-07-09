package repositories

import (
	"context"
	"my-finance-hub-api/internal/domain/entities"
)

type CategoryRepository interface {
	Create(ctx context.Context, category *entities.Category) error
	GetByID(ctx context.Context, id uint) (*entities.Category, error)
	GetByUserID(ctx context.Context, userID uint) ([]*entities.Category, error)
	Update(ctx context.Context, category *entities.Category) error
	Delete(ctx context.Context, id uint) error
	ExistsByName(ctx context.Context, userID uint, name string) (bool, error)
}
