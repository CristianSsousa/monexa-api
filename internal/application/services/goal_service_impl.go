package services

import (
	"context"
	"my-finance-hub-api/internal/application/interfaces"
	"my-finance-hub-api/internal/domain/entities"
	"my-finance-hub-api/internal/domain/repositories"
	pkgErrors "my-finance-hub-api/pkg/errors"
)

type goalServiceImpl struct {
	goalRepo repositories.GoalRepository
}

func NewGoalService(goalRepo repositories.GoalRepository) interfaces.GoalService {
	return &goalServiceImpl{
		goalRepo: goalRepo,
	}
}

func (s *goalServiceImpl) CreateGoal(ctx context.Context, userID uint, goal *entities.Goal) (*entities.Goal, error) {
	// Validações
	if goal.Name == "" {
		return nil, pkgErrors.NewDomainError("validation_error", "Nome da meta é obrigatório")
	}

	if goal.Amount <= 0 {
		return nil, pkgErrors.NewDomainError("validation_error", "Valor da meta deve ser maior que zero")
	}

	// Criar nova meta
	newGoal := entities.NewGoal(goal.Name, goal.Description, goal.Amount, userID)

	if err := s.goalRepo.Create(ctx, newGoal); err != nil {
		return nil, err
	}

	return newGoal, nil
}

func (s *goalServiceImpl) GetGoalByID(ctx context.Context, userID, goalID uint) (*entities.Goal, error) {
	goal, err := s.goalRepo.GetByID(ctx, goalID)
	if err != nil {
		return nil, err
	}

	// Verificar se a meta pertence ao usuário
	if !goal.BelongsToUser(userID) {
		return nil, pkgErrors.ErrForbidden
	}

	return goal, nil
}

func (s *goalServiceImpl) GetGoalsByUser(ctx context.Context, userID uint) ([]*entities.Goal, error) {
	return s.goalRepo.GetByUserID(ctx, userID)
}

func (s *goalServiceImpl) UpdateGoal(ctx context.Context, userID, goalID uint, updates *entities.Goal) (*entities.Goal, error) {
	// Buscar meta existente
	goal, err := s.GetGoalByID(ctx, userID, goalID)
	if err != nil {
		return nil, err
	}

	// Validações
	if updates.Name == "" {
		return nil, pkgErrors.NewDomainError("validation_error", "Nome da meta é obrigatório")
	}

	if updates.Amount <= 0 {
		return nil, pkgErrors.NewDomainError("validation_error", "Valor da meta deve ser maior que zero")
	}

	// Atualizar meta
	goal.Update(updates.Name, updates.Description, updates.Amount)

	if err := s.goalRepo.Update(ctx, goal); err != nil {
		return nil, err
	}

	return goal, nil
}

func (s *goalServiceImpl) DeleteGoal(ctx context.Context, userID, goalID uint) error {
	// Verificar se a meta existe e pertence ao usuário
	_, err := s.GetGoalByID(ctx, userID, goalID)
	if err != nil {
		return err
	}

	// Excluir meta
	return s.goalRepo.Delete(ctx, goalID)
}
