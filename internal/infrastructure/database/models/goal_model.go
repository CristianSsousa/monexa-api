package models

import (
	"my-finance-hub-api/internal/domain/entities"
	"time"

	"gorm.io/gorm"
)

type Goal struct {
	ID          uint   `gorm:"primaryKey"`
	Name        string `gorm:"not null"`
	Description string
	Amount      float64 `gorm:"not null"`
	UserID      uint    `gorm:"not null"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

// FromEntity converte uma entidade de domínio para o modelo GORM
func (g *Goal) FromEntity(entity *entities.Goal) {
	g.ID = entity.ID
	g.Name = entity.Name
	g.Description = entity.Description
	g.Amount = entity.Amount
	g.UserID = entity.UserID
	g.CreatedAt = entity.CreatedAt
	g.UpdatedAt = entity.UpdatedAt
}

// ToEntity converte o modelo GORM para uma entidade de domínio
func (g *Goal) ToEntity() *entities.Goal {
	return &entities.Goal{
		ID:          g.ID,
		Name:        g.Name,
		Description: g.Description,
		Amount:      g.Amount,
		UserID:      g.UserID,
		CreatedAt:   g.CreatedAt,
		UpdatedAt:   g.UpdatedAt,
	}
}

// TableName especifica o nome da tabela
func (Goal) TableName() string {
	return "goals"
}
