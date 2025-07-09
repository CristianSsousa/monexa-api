package repositories

import (
	"context"
	"my-finance-hub-api/internal/domain/entities"
)

type SavingGoalRepository interface {
	Create(ctx context.Context, savingGoal *entities.SavingGoal) error
	GetByID(ctx context.Context, id uint) (*entities.SavingGoal, error)
	GetByUserID(ctx context.Context, userID uint) ([]*entities.SavingGoal, error)
	Update(ctx context.Context, savingGoal *entities.SavingGoal) error
	Delete(ctx context.Context, id uint) error
	UpdateAmount(ctx context.Context, id uint, newAmount float64) error
}
