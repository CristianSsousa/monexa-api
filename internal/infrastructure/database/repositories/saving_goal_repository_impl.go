package repositories

import (
	"context"
	"errors"
	"log"
	"my-finance-hub-api/internal/domain/entities"
	"my-finance-hub-api/internal/domain/repositories"
	"my-finance-hub-api/internal/infrastructure/database/models"
	pkgErrors "my-finance-hub-api/pkg/errors"
	"time"

	"gorm.io/gorm"
)

type savingGoalRepositoryImpl struct {
	db *gorm.DB
}

func NewSavingGoalRepository(db *gorm.DB) repositories.SavingGoalRepository {
	return &savingGoalRepositoryImpl{
		db: db,
	}
}

func (r *savingGoalRepositoryImpl) Create(ctx context.Context, savingGoal *entities.SavingGoal) error {
	model := &models.SavingGoal{}
	model.FromEntity(savingGoal)

	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return err
	}

	// Atualiza a entidade com o ID gerado
	savingGoal.ID = model.ID
	savingGoal.CreatedAt = model.CreatedAt
	savingGoal.UpdatedAt = model.UpdatedAt

	return nil
}

func (r *savingGoalRepositoryImpl) GetByID(ctx context.Context, id uint) (*entities.SavingGoal, error) {
	var model models.SavingGoal

	if err := r.db.WithContext(ctx).First(&model, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, pkgErrors.ErrSavingGoalNotFound
		}
		return nil, err
	}

	return model.ToEntity(), nil
}

func (r *savingGoalRepositoryImpl) GetByUserID(ctx context.Context, userID uint) ([]*entities.SavingGoal, error) {
	var models []models.SavingGoal

	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&models).Error; err != nil {
		return nil, err
	}

	savingGoals := make([]*entities.SavingGoal, len(models))
	for i, model := range models {
		savingGoals[i] = model.ToEntity()
	}

	return savingGoals, nil
}

func (r *savingGoalRepositoryImpl) Update(ctx context.Context, savingGoal *entities.SavingGoal) error {
	model := &models.SavingGoal{}
	model.FromEntity(savingGoal)

	if err := r.db.WithContext(ctx).Save(model).Error; err != nil {
		return err
	}

	// Atualiza a entidade com o timestamp
	savingGoal.UpdatedAt = model.UpdatedAt

	return nil
}

func (r *savingGoalRepositoryImpl) Delete(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Delete(&models.SavingGoal{}, id)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return pkgErrors.ErrSavingGoalNotFound
	}

	return nil
}

func (r *savingGoalRepositoryImpl) UpdateAmount(ctx context.Context, id uint, newAmount float64) error {
	// Primeiro, verificar se o registro existe
	var existingSavingGoal models.SavingGoal
	if err := r.db.WithContext(ctx).First(&existingSavingGoal, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return pkgErrors.ErrSavingGoalNotFound
		}
		return err
	}

	// Atualizar o valor
	result := r.db.WithContext(ctx).Model(&models.SavingGoal{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"current_amount": newAmount,
			"updated_at":     time.Now(),
		})

	if result.Error != nil {
		log.Printf("Erro ao atualizar valor do cofrinho: %v", result.Error)
		return result.Error
	}

	if result.RowsAffected == 0 {
		log.Printf("Nenhuma linha atualizada para o cofrinho ID %d", id)
		return pkgErrors.ErrSavingGoalNotFound
	}

	log.Printf("Cofrinho ID %d atualizado com sucesso. Novo valor: %.2f", id, newAmount)
	return nil
}
