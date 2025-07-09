package repositories

import (
	"context"
	"errors"
	"my-finance-hub-api/internal/domain/entities"
	"my-finance-hub-api/internal/domain/repositories"
	"my-finance-hub-api/internal/infrastructure/database/models"
	pkgErrors "my-finance-hub-api/pkg/errors"

	"gorm.io/gorm"
)

type userRepositoryImpl struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) repositories.UserRepository {
	return &userRepositoryImpl{
		db: db,
	}
}

func (r *userRepositoryImpl) Create(ctx context.Context, user *entities.User) error {
	model := &models.User{}
	model.FromEntity(user)

	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return pkgErrors.ErrEmailAlreadyExists
		}
		return err
	}

	// Atualiza a entidade com o ID gerado
	user.ID = model.ID
	user.CreatedAt = model.CreatedAt
	user.UpdatedAt = model.UpdatedAt

	return nil
}

func (r *userRepositoryImpl) GetByID(ctx context.Context, id uint) (*entities.User, error) {
	var model models.User

	if err := r.db.WithContext(ctx).First(&model, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, pkgErrors.ErrUserNotFound
		}
		return nil, err
	}

	return model.ToEntity(), nil
}

func (r *userRepositoryImpl) GetByEmail(ctx context.Context, email string) (*entities.User, error) {
	var model models.User

	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, pkgErrors.ErrUserNotFound
		}
		return nil, err
	}

	return model.ToEntity(), nil
}

func (r *userRepositoryImpl) Update(ctx context.Context, user *entities.User) error {
	model := &models.User{}
	model.FromEntity(user)

	if err := r.db.WithContext(ctx).Save(model).Error; err != nil {
		return err
	}

	// Atualiza a entidade com o timestamp
	user.UpdatedAt = model.UpdatedAt

	return nil
}

func (r *userRepositoryImpl) Delete(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Delete(&models.User{}, id)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return pkgErrors.ErrUserNotFound
	}

	return nil
}

func (r *userRepositoryImpl) Exists(ctx context.Context, id uint) (bool, error) {
	var count int64

	if err := r.db.WithContext(ctx).Model(&models.User{}).Where("id = ?", id).Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *userRepositoryImpl) EmailExists(ctx context.Context, email string) (bool, error) {
	var count int64

	if err := r.db.WithContext(ctx).Model(&models.User{}).Where("email = ?", email).Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}
