package services

import (
	"context"
	"os"
	"time"

	"my-finance-hub-api/internal/application/interfaces"
	"my-finance-hub-api/internal/domain/entities"
	"my-finance-hub-api/internal/domain/repositories"
	pkgErrors "my-finance-hub-api/pkg/errors"

	"github.com/golang-jwt/jwt/v5"
)

type authServiceImpl struct {
	userRepo repositories.UserRepository
}

func NewAuthService(userRepo repositories.UserRepository) interfaces.AuthService {
	return &authServiceImpl{
		userRepo: userRepo,
	}
}

func (s *authServiceImpl) Register(ctx context.Context, name, email, password string) (*entities.User, error) {
	// Verificar se email já existe
	exists, err := s.userRepo.EmailExists(ctx, email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, pkgErrors.ErrEmailAlreadyExists
	}

	// Criar usuário
	user, err := entities.NewUser(name, email, password)
	if err != nil {
		return nil, err
	}

	// Salvar no repositório
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *authServiceImpl) Login(ctx context.Context, email, password string) (*entities.User, string, error) {
	// Buscar usuário por email
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, "", pkgErrors.ErrInvalidCredentials
	}

	// Verificar senha
	if err := user.CheckPassword(password); err != nil {
		return nil, "", pkgErrors.ErrInvalidCredentials
	}

	// Gerar token JWT
	token, err := s.generateJWT(user.ID)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

func (s *authServiceImpl) GetUserByID(ctx context.Context, userID uint) (*entities.User, error) {
	return s.userRepo.GetByID(ctx, userID)
}

func (s *authServiceImpl) ChangePassword(ctx context.Context, userID uint, currentPassword, newPassword string) error {
	// Buscar usuário
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	// Verificar senha atual
	if err := user.CheckPassword(currentPassword); err != nil {
		return pkgErrors.ErrInvalidCredentials
	}

	// Atualizar senha
	if err := user.UpdatePassword(newPassword); err != nil {
		return err
	}

	// Salvar no repositório
	return s.userRepo.Update(ctx, user)
}

func (s *authServiceImpl) UpdateUser(ctx context.Context, userID uint, name, email string) (*entities.User, error) {
	// Buscar usuário
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Verificar se o email não está sendo usado por outro usuário
	if email != user.Email {
		exists, err := s.userRepo.EmailExists(ctx, email)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, pkgErrors.ErrEmailAlreadyExists
		}
	}

	// Atualizar dados
	user.Update(name, email)

	// Salvar no repositório
	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *authServiceImpl) generateJWT(userID uint) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
		"iat":     time.Now().Unix(),
	})

	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		secretKey = "your-256-bit-secret" // Chave padrão para desenvolvimento
	}

	return token.SignedString([]byte(secretKey))
}
