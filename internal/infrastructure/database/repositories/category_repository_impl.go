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

type categoryRepositoryImpl struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) repositories.CategoryRepository {
	return &categoryRepositoryImpl{
		db: db,
	}
}

func (r *categoryRepositoryImpl) Create(ctx context.Context, category *entities.Category) error {
	model := &models.Category{}
	model.FromEntity(category)

	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return err
	}

	// Atualiza a entidade com o ID gerado
	category.ID = model.ID
	category.CreatedAt = model.CreatedAt
	category.UpdatedAt = model.UpdatedAt

	return nil
}

func (r *categoryRepositoryImpl) GetByID(ctx context.Context, id uint) (*entities.Category, error) {
	var model models.Category

	if err := r.db.WithContext(ctx).First(&model, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, pkgErrors.ErrCategoryNotFound
		}
		return nil, err
	}

	return model.ToEntity(), nil
}

func (r *categoryRepositoryImpl) GetByUserID(ctx context.Context, userID uint) ([]*entities.Category, error) {
	var models []models.Category

	// Buscar categorias do usuário e categorias globais (sem UserID)
	if err := r.db.WithContext(ctx).
		Where("user_id = ? OR user_id IS NULL", userID).
		Find(&models).Error; err != nil {
		return nil, err
	}

	categories := make([]*entities.Category, len(models))
	for i, model := range models {
		categories[i] = model.ToEntity()
	}

	return categories, nil
}

func (r *categoryRepositoryImpl) Update(ctx context.Context, category *entities.Category) error {
	model := &models.Category{}
	model.FromEntity(category)

	if err := r.db.WithContext(ctx).Save(model).Error; err != nil {
		return err
	}

	// Atualiza a entidade com o timestamp
	category.UpdatedAt = model.UpdatedAt

	return nil
}

func (r *categoryRepositoryImpl) Delete(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Delete(&models.Category{}, id)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return pkgErrors.ErrCategoryNotFound
	}

	return nil
}

func (r *categoryRepositoryImpl) ExistsByName(ctx context.Context, userID uint, name string) (bool, error) {
	var count int64

	if err := r.db.WithContext(ctx).Model(&models.Category{}).
		Where("user_id = ? AND name = ?", userID, name).
		Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}
