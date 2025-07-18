package repositories

import (
	"context"
	"errors"

	"my-finance-hub-api/internal/domain/entities"
	"my-finance-hub-api/internal/domain/repositories"
	"my-finance-hub-api/internal/infrastructure/database/models"

	"gorm.io/gorm"
)

type UserRepositoryImpl struct {
	DB *gorm.DB
}

var _ repositories.UserRepository = &UserRepositoryImpl{}

func NewUserRepository(db *gorm.DB) repositories.UserRepository {
	return &UserRepositoryImpl{
		DB: db,
	}
}

func (r *UserRepositoryImpl) Create(ctx context.Context, user *entities.User) error {
	model := &models.User{}
	model.FromEntity(user)

	if err := r.DB.WithContext(ctx).Create(model).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return errors.New("email já está em uso") // Changed from pkgErrors.ErrEmailAlreadyExists
		}
		return err
	}

	// Atualiza a entidade com o ID gerado
	user.ID = model.ID
	user.CreatedAt = model.CreatedAt
	user.UpdatedAt = model.UpdatedAt

	return nil
}

func (r *UserRepositoryImpl) GetByID(ctx context.Context, id uint) (*entities.User, error) {
	var model models.User

	if err := r.DB.WithContext(ctx).First(&model, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("usuário não encontrado") // Changed from pkgErrors.ErrUserNotFound
		}
		return nil, err
	}

	return model.ToEntity(), nil
}

func (r *UserRepositoryImpl) GetByEmail(ctx context.Context, email string) (*entities.User, error) {
	var model models.User

	if err := r.DB.WithContext(ctx).Where("email = ?", email).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("usuário não encontrado") // Changed from pkgErrors.ErrUserNotFound
		}
		return nil, err
	}

	return model.ToEntity(), nil
}

func (r *UserRepositoryImpl) Update(ctx context.Context, user *entities.User) error {
	model := &models.User{}
	model.FromEntity(user)

	if err := r.DB.WithContext(ctx).Save(model).Error; err != nil {
		return err
	}

	// Atualiza a entidade com o timestamp
	user.UpdatedAt = model.UpdatedAt

	return nil
}

func (r *UserRepositoryImpl) Delete(ctx context.Context, id uint) error {
	result := r.DB.WithContext(ctx).Delete(&models.User{}, id)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("usuário não encontrado") // Changed from pkgErrors.ErrUserNotFound
	}

	return nil
}

func (r *UserRepositoryImpl) Exists(ctx context.Context, id uint) (bool, error) {
	var count int64

	if err := r.DB.WithContext(ctx).Model(&models.User{}).Where("id = ?", id).Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *UserRepositoryImpl) EmailExists(ctx context.Context, email string) (bool, error) {
	var count int64

	if err := r.DB.WithContext(ctx).Model(&models.User{}).Where("email = ?", email).Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *UserRepositoryImpl) UpdateUser(ctx context.Context, userID uint, updateData *entities.User) error {
	result := r.DB.WithContext(ctx).Model(&models.User{}).Where("id = ?", userID).Updates(updateData)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("usuário não encontrado")
	}
	return nil
}
