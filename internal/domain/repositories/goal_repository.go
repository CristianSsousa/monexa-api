package repositories

import (
	"context"
	"my-finance-hub-api/internal/domain/entities"
)

type GoalRepository interface {
	Create(ctx context.Context, goal *entities.Goal) error
	GetByID(ctx context.Context, id uint) (*entities.Goal, error)
	GetByUserID(ctx context.Context, userID uint) ([]*entities.Goal, error)
	Update(ctx context.Context, goal *entities.Goal) error
	Delete(ctx context.Context, id uint) error
}
