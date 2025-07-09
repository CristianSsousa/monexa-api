package interfaces

import (
	"context"
	"my-finance-hub-api/internal/domain/entities"
)

type SavingGoalService interface {
	CreateSavingGoal(ctx context.Context, userID uint, savingGoal *entities.SavingGoal) (*entities.SavingGoal, error)
	GetSavingGoalByID(ctx context.Context, userID, savingGoalID uint) (*entities.SavingGoal, error)
	GetSavingGoalsByUser(ctx context.Context, userID uint) ([]*entities.SavingGoal, error)
	UpdateSavingGoal(ctx context.Context, userID, savingGoalID uint, updates *entities.SavingGoal) (*entities.SavingGoal, error)
	DeleteSavingGoal(ctx context.Context, userID, savingGoalID uint) error
	Deposit(ctx context.Context, userID, savingGoalID uint, amount float64) (*entities.SavingGoal, error)
}
