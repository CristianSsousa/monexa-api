package interfaces

import (
	"context"
	"my-finance-hub-api/internal/domain/entities"
)

type GoalService interface {
	CreateGoal(ctx context.Context, userID uint, goal *entities.Goal) (*entities.Goal, error)
	GetGoalByID(ctx context.Context, userID, goalID uint) (*entities.Goal, error)
	GetGoalsByUser(ctx context.Context, userID uint) ([]*entities.Goal, error)
	UpdateGoal(ctx context.Context, userID, goalID uint, updates *entities.Goal) (*entities.Goal, error)
	DeleteGoal(ctx context.Context, userID, goalID uint) error
}
