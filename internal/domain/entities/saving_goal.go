package entities

import (
	"my-finance-hub-api/pkg/errors"
	"time"
)

type SavingGoal struct {
	ID            uint
	Name          string
	TargetAmount  float64
	CurrentAmount float64
	Description   string
	UserID        uint
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// NewSavingGoal creates a new SavingGoal entity
func NewSavingGoal(name string, targetAmount float64, userID uint, currentAmount float64, description string) *SavingGoal {
	return &SavingGoal{
		Name:          name,
		TargetAmount:  targetAmount,
		CurrentAmount: currentAmount,
		Description:   description,
		UserID:        userID,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
}

// Deposit adiciona um valor ao cofrinho
func (sg *SavingGoal) Deposit(amount float64) {
	sg.CurrentAmount += amount
	sg.UpdatedAt = time.Now()
}

// Withdraw remove um valor do cofrinho
func (sg *SavingGoal) Withdraw(amount float64) error {
	if sg.CurrentAmount < amount {
		return errors.NewDomainError("insufficient_funds", "Saldo insuficiente no cofrinho")
	}
	sg.CurrentAmount -= amount
	sg.UpdatedAt = time.Now()
	return nil
}

// GetProgress retorna o progresso em porcentagem
func (sg *SavingGoal) GetProgress() float64 {
	if sg.TargetAmount == 0 {
		return 0
	}
	progress := (sg.CurrentAmount / sg.TargetAmount) * 100
	if progress > 100 {
		return 100
	}
	return progress
}

// IsCompleted verifica se a meta foi atingida
func (sg *SavingGoal) IsCompleted() bool {
	return sg.CurrentAmount >= sg.TargetAmount
}

// Update atualiza os dados da meta de economia
func (sg *SavingGoal) Update(name string, targetAmount float64, currentAmount float64, description string) {
	sg.Name = name
	sg.TargetAmount = targetAmount
	sg.CurrentAmount = currentAmount
	sg.Description = description
	sg.UpdatedAt = time.Now()
}

// BelongsToUser verifica se a meta de economia pertence ao usuário
func (sg *SavingGoal) BelongsToUser(userID uint) bool {
	return sg.UserID == userID
}
