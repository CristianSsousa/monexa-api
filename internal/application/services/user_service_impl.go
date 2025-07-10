package services

import (
	"context"
	"errors"

	"my-finance-hub-api/internal/application/interfaces"
	"my-finance-hub-api/internal/domain/entities"
	"my-finance-hub-api/internal/domain/repositories"

	"golang.org/x/crypto/bcrypt"
)

type UserServiceImpl struct {
	userRepo repositories.UserRepository
}

func NewUserService(userRepo repositories.UserRepository) interfaces.UserService {
	return &UserServiceImpl{userRepo: userRepo}
}

func (s *UserServiceImpl) UpdateUserProfile(ctx context.Context, userID uint, updateData *entities.User) error {
	// Validar dados de entrada
	if updateData.Name == "" {
		return errors.New("nome não pode ser vazio")
	}

	// Se a senha for fornecida, hash da senha
	if updateData.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(updateData.Password), bcrypt.DefaultCost)
		if err != nil {
			return errors.New("erro ao criptografar senha")
		}
		updateData.Password = string(hashedPassword)
	}

	// Atualizar usuário no repositório
	return s.userRepo.UpdateUser(ctx, userID, updateData)
}
