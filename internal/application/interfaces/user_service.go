package interfaces

import (
	"context"
	"my-finance-hub-api/internal/domain/entities"
)

type UserService interface {
	UpdateUserProfile(ctx context.Context, userID uint, updateData *entities.User) error
}
