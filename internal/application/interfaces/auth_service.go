package interfaces

import (
	"context"
	"my-finance-hub-api/internal/domain/entities"
)

type AuthService interface {
	Register(ctx context.Context, name, email, password string) (*entities.User, error)
	Login(ctx context.Context, email, password string) (*entities.User, string, error)
	GetUserByID(ctx context.Context, userID uint) (*entities.User, error)
	ChangePassword(ctx context.Context, userID uint, currentPassword, newPassword string) error
	UpdateUser(ctx context.Context, userID uint, name, email string) (*entities.User, error)
}
