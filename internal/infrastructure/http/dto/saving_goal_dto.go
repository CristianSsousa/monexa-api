package dto

import (
	"my-finance-hub-api/internal/domain/entities"
	"time"
)

// Request DTOs
type CreateSavingGoalRequest struct {
	Name          string  `json:"name" binding:"required,min=2,max=100"`
	TargetAmount  float64 `json:"target_amount" binding:"required,gt=0"`
	CurrentAmount float64 `json:"current_amount" binding:"omitempty,gte=0"`
	UserID        uint    `json:"user_id" binding:"omitempty"`
	Description   string  `json:"description"`
}

type UpdateSavingGoalRequest struct {
	Name          string  `json:"name" binding:"required,min=2,max=100"`
	TargetAmount  float64 `json:"target_amount" binding:"required,gt=0"`
	CurrentAmount float64 `json:"current_amount" binding:"required,gt=0"`
	Description   string  `json:"description"`
}

type DepositRequest struct {
	Amount float64 `json:"amount" binding:"required,gt=0"`
}

// Response DTOs
type SavingGoalResponse struct {
	ID            uint      `json:"id"`
	Name          string    `json:"name"`
	TargetAmount  float64   `json:"target_amount"`
	CurrentAmount float64   `json:"current_amount"`
	Progress      float64   `json:"progress"`
	IsCompleted   bool      `json:"is_completed"`
	UserID        uint      `json:"user_id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// Mappers
func ToSavingGoalResponse(savingGoal *entities.SavingGoal) SavingGoalResponse {
	return SavingGoalResponse{
		ID:            savingGoal.ID,
		Name:          savingGoal.Name,
		TargetAmount:  savingGoal.TargetAmount,
		CurrentAmount: savingGoal.CurrentAmount,
		Progress:      savingGoal.GetProgress(),
		IsCompleted:   savingGoal.IsCompleted(),
		UserID:        savingGoal.UserID,
		CreatedAt:     savingGoal.CreatedAt,
		UpdatedAt:     savingGoal.UpdatedAt,
	}
}

func ToSavingGoalResponseList(savingGoals []*entities.SavingGoal) []SavingGoalResponse {
	result := make([]SavingGoalResponse, len(savingGoals))
	for i, savingGoal := range savingGoals {
		result[i] = ToSavingGoalResponse(savingGoal)
	}
	return result
}

func (req *CreateSavingGoalRequest) ToEntity(userID uint) *entities.SavingGoal {
	return entities.NewSavingGoal(req.Name, req.TargetAmount, userID, req.CurrentAmount, req.Description)
}

func (req *UpdateSavingGoalRequest) ToEntity(userID uint) *entities.SavingGoal {
	return entities.NewSavingGoal(req.Name, req.TargetAmount, userID, req.CurrentAmount, req.Description)
}
