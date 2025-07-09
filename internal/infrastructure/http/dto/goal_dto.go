package dto

import (
	"my-finance-hub-api/internal/domain/entities"
	"time"
)

// Request DTOs
type CreateGoalRequest struct {
	Name        string  `json:"name" binding:"required,min=2,max=100"`
	Description string  `json:"description" binding:"required,max=500"`
	Amount      float64 `json:"amount" binding:"required,gt=0"`
}

type UpdateGoalRequest struct {
	Name        string  `json:"name" binding:"required,min=2,max=100"`
	Description string  `json:"description" binding:"required,max=500"`
	Amount      float64 `json:"amount" binding:"required,gt=0"`
}

// Response DTOs
type GoalResponse struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Amount      float64   `json:"amount"`
	UserID      uint      `json:"user_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Mappers
func ToGoalResponse(goal *entities.Goal) GoalResponse {
	return GoalResponse{
		ID:          goal.ID,
		Name:        goal.Name,
		Description: goal.Description,
		Amount:      goal.Amount,
		UserID:      goal.UserID,
		CreatedAt:   goal.CreatedAt,
		UpdatedAt:   goal.UpdatedAt,
	}
}

func ToGoalResponseList(goals []*entities.Goal) []GoalResponse {
	result := make([]GoalResponse, len(goals))
	for i, goal := range goals {
		result[i] = ToGoalResponse(goal)
	}
	return result
}

func (req *CreateGoalRequest) ToEntity(userID uint) *entities.Goal {
	return entities.NewGoal(req.Name, req.Description, req.Amount, userID)
}

func (req *UpdateGoalRequest) ToEntity(userID uint) *entities.Goal {
	return entities.NewGoal(req.Name, req.Description, req.Amount, userID)
}
