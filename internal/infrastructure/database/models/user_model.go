package models

import (
	"my-finance-hub-api/internal/domain/entities"
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"not null"`
	Email     string `gorm:"unique;not null"`
	Password  string `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// TableName especifica o nome da tabela
func (User) TableName() string {
	return "users"
}

// ToEntity converte o modelo GORM para entidade de domínio
func (u *User) ToEntity() *entities.User {
	return &entities.User{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		Password:  u.Password,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

// FromEntity converte entidade de domínio para modelo GORM
func (u *User) FromEntity(entity *entities.User) {
	u.ID = entity.ID
	u.Name = entity.Name
	u.Email = entity.Email
	u.Password = entity.Password
	u.CreatedAt = entity.CreatedAt
	u.UpdatedAt = entity.UpdatedAt
}
