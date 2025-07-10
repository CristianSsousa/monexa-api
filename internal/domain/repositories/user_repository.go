package repositories

import (
	"context"
	"my-finance-hub-api/internal/domain/entities"
)

type UserRepository interface {
	Create(ctx context.Context, user *entities.User) error
	GetByID(ctx context.Context, id uint) (*entities.User, error)
	GetByEmail(ctx context.Context, email string) (*entities.User, error)
	Update(ctx context.Context, user *entities.User) error
	Delete(ctx context.Context, id uint) error
	Exists(ctx context.Context, id uint) (bool, error)
	EmailExists(ctx context.Context, email string) (bool, error)
	UpdateUser(ctx context.Context, userID uint, updateData *entities.User) error
}
