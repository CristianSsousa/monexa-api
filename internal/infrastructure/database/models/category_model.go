package models

import (
	"my-finance-hub-api/internal/domain/entities"
	"time"

	"gorm.io/gorm"
)

type Category struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"not null"`
	Color     string `gorm:"not null"`
	UserID    uint   `gorm:"column:user_id"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// FromEntity converte uma entidade de domínio para o modelo GORM
func (c *Category) FromEntity(entity *entities.Category) {
	c.ID = entity.ID
	c.Name = entity.Name
	c.Color = entity.Color
	c.UserID = entity.UserID
	c.CreatedAt = entity.CreatedAt
	c.UpdatedAt = entity.UpdatedAt
}

// ToEntity converte o modelo GORM para uma entidade de domínio
func (c *Category) ToEntity() *entities.Category {
	return &entities.Category{
		ID:        c.ID,
		Name:      c.Name,
		Color:     c.Color,
		UserID:    c.UserID,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}

// TableName especifica o nome da tabela
func (Category) TableName() string {
	return "categories"
}
