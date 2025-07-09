package models

import (
	"my-finance-hub-api/internal/domain/entities"
	"time"

	"gorm.io/gorm"
)

type SavingGoal struct {
	ID            uint   `gorm:"primaryKey"`
	Name          string `gorm:"not null"`
	Description   string
	TargetAmount  float64 `gorm:"not null"`
	CurrentAmount float64 `gorm:"default:0"`
	UserID        uint    `gorm:"not null"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt `gorm:"index"`
}

// FromEntity converte uma entidade de domínio para o modelo GORM
func (sg *SavingGoal) FromEntity(entity *entities.SavingGoal) {
	sg.ID = entity.ID
	sg.Name = entity.Name
	sg.TargetAmount = entity.TargetAmount
	sg.CurrentAmount = entity.CurrentAmount
	sg.UserID = entity.UserID
	sg.CreatedAt = entity.CreatedAt
	sg.UpdatedAt = entity.UpdatedAt
}

// ToEntity converte o modelo GORM para uma entidade de domínio
func (sg *SavingGoal) ToEntity() *entities.SavingGoal {
	return &entities.SavingGoal{
		ID:            sg.ID,
		Name:          sg.Name,
		TargetAmount:  sg.TargetAmount,
		CurrentAmount: sg.CurrentAmount,
		UserID:        sg.UserID,
		CreatedAt:     sg.CreatedAt,
		UpdatedAt:     sg.UpdatedAt,
	}
}

// TableName especifica o nome da tabela
func (SavingGoal) TableName() string {
	return "saving_goals"
}
