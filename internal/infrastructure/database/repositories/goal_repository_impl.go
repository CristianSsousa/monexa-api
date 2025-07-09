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

type goalRepositoryImpl struct {
	db *gorm.DB
}

func NewGoalRepository(db *gorm.DB) repositories.GoalRepository {
	return &goalRepositoryImpl{
		db: db,
	}
}

func (r *goalRepositoryImpl) Create(ctx context.Context, goal *entities.Goal) error {
	model := &models.Goal{}
	model.FromEntity(goal)

	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return err
	}

	// Atualiza a entidade com o ID gerado
	goal.ID = model.ID
	goal.CreatedAt = model.CreatedAt
	goal.UpdatedAt = model.UpdatedAt

	return nil
}

func (r *goalRepositoryImpl) GetByID(ctx context.Context, id uint) (*entities.Goal, error) {
	var model models.Goal

	if err := r.db.WithContext(ctx).First(&model, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, pkgErrors.ErrGoalNotFound
		}
		return nil, err
	}

	return model.ToEntity(), nil
}

func (r *goalRepositoryImpl) GetByUserID(ctx context.Context, userID uint) ([]*entities.Goal, error) {
	var models []models.Goal

	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&models).Error; err != nil {
		return nil, err
	}

	goals := make([]*entities.Goal, len(models))
	for i, model := range models {
		goals[i] = model.ToEntity()
	}

	return goals, nil
}

func (r *goalRepositoryImpl) Update(ctx context.Context, goal *entities.Goal) error {
	model := &models.Goal{}
	model.FromEntity(goal)

	if err := r.db.WithContext(ctx).Save(model).Error; err != nil {
		return err
	}

	// Atualiza a entidade com o timestamp
	goal.UpdatedAt = model.UpdatedAt

	return nil
}

func (r *goalRepositoryImpl) Delete(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Delete(&models.Goal{}, id)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return pkgErrors.ErrGoalNotFound
	}

	return nil
}
