package models

import (
	"my-finance-hub-api/internal/domain/entities"
	"time"

	"gorm.io/gorm"
)

type Transaction struct {
	ID          uint      `gorm:"primaryKey"`
	Description string    `gorm:"not null"`
	Amount      float64   `gorm:"not null"`
	Type        string    `gorm:"not null"`
	Date        time.Time `gorm:"not null"`
	Paid        bool      `gorm:"default:false"`
	UserID      uint      `gorm:"not null"`
	CategoryID  *uint     `gorm:"column:category_id"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

// FromEntity converte uma entidade de domínio para o modelo GORM
func (t *Transaction) FromEntity(entity *entities.Transaction) {
	t.ID = entity.ID
	t.Description = entity.Description
	t.Amount = entity.Amount
	t.Type = string(entity.Type)
	t.Date = entity.Date
	t.Paid = entity.Paid
	t.UserID = entity.UserID

	// Converter CategoryID para ponteiro
	if entity.CategoryID != nil {
		categoryID := *entity.CategoryID
		t.CategoryID = &categoryID
	}

	t.CreatedAt = entity.CreatedAt
	t.UpdatedAt = entity.UpdatedAt
}

// ToEntity converte o modelo GORM para uma entidade de domínio
func (t *Transaction) ToEntity() *entities.Transaction {
	return &entities.Transaction{
		ID:          t.ID,
		Description: t.Description,
		Amount:      t.Amount,
		Type:        entities.TransactionType(t.Type),
		Date:        t.Date,
		Paid:        t.Paid,
		UserID:      t.UserID,
		CategoryID:  t.CategoryID,
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
	}
}

// TableName especifica o nome da tabela
func (Transaction) TableName() string {
	return "transactions"
}
