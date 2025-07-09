package services

import (
	"context"
	"log"
	"my-finance-hub-api/internal/application/interfaces"
	"my-finance-hub-api/internal/domain/entities"
	"my-finance-hub-api/internal/domain/repositories"
	pkgErrors "my-finance-hub-api/pkg/errors"
	"time"
)

type savingGoalServiceImpl struct {
	savingGoalRepo repositories.SavingGoalRepository
}

func NewSavingGoalService(savingGoalRepo repositories.SavingGoalRepository) interfaces.SavingGoalService {
	return &savingGoalServiceImpl{
		savingGoalRepo: savingGoalRepo,
	}
}

func (s *savingGoalServiceImpl) CreateSavingGoal(ctx context.Context, userID uint, savingGoal *entities.SavingGoal) (*entities.SavingGoal, error) {
	// Validações
	if savingGoal.Name == "" {
		return nil, pkgErrors.NewDomainError("validation_error", "Nome da meta de economia é obrigatório")
	}

	if savingGoal.TargetAmount <= 0 {
		return nil, pkgErrors.NewDomainError("validation_error", "Valor da meta deve ser maior que zero")
	}

	// Criar nova meta de economia
	newSavingGoal := entities.NewSavingGoal(savingGoal.Name, savingGoal.TargetAmount, userID, savingGoal.CurrentAmount, savingGoal.Description)

	if err := s.savingGoalRepo.Create(ctx, newSavingGoal); err != nil {
		return nil, err
	}

	return newSavingGoal, nil
}

func (s *savingGoalServiceImpl) GetSavingGoalByID(ctx context.Context, userID, savingGoalID uint) (*entities.SavingGoal, error) {
	savingGoal, err := s.savingGoalRepo.GetByID(ctx, savingGoalID)
	if err != nil {
		return nil, err
	}

	// Verificar se a meta de economia pertence ao usuário
	if !savingGoal.BelongsToUser(userID) {
		return nil, pkgErrors.ErrForbidden
	}

	return savingGoal, nil
}

func (s *savingGoalServiceImpl) GetSavingGoalsByUser(ctx context.Context, userID uint) ([]*entities.SavingGoal, error) {
	return s.savingGoalRepo.GetByUserID(ctx, userID)
}

func (s *savingGoalServiceImpl) UpdateSavingGoal(ctx context.Context, userID, savingGoalID uint, updates *entities.SavingGoal) (*entities.SavingGoal, error) {
	// Buscar meta de economia existente
	savingGoal, err := s.GetSavingGoalByID(ctx, userID, savingGoalID)
	if err != nil {
		return nil, err
	}

	// Validações
	if updates.Name == "" {
		return nil, pkgErrors.NewDomainError("validation_error", "Nome da meta de economia é obrigatório")
	}

	if updates.TargetAmount <= 0 {
		return nil, pkgErrors.NewDomainError("validation_error", "Valor da meta deve ser maior que zero")
	}

	if updates.CurrentAmount <= 0 {
		return nil, pkgErrors.NewDomainError("validation_error", "Valor atual da meta deve ser maior que zero")
	}

	// Atualizar meta de economia
	savingGoal.Update(updates.Name, updates.TargetAmount, updates.CurrentAmount, updates.Description)

	if err := s.savingGoalRepo.Update(ctx, savingGoal); err != nil {
		return nil, err
	}

	return savingGoal, nil
}

func (s *savingGoalServiceImpl) DeleteSavingGoal(ctx context.Context, userID, savingGoalID uint) error {
	// Verificar se a meta de economia existe e pertence ao usuário
	_, err := s.GetSavingGoalByID(ctx, userID, savingGoalID)
	if err != nil {
		return err
	}

	// Excluir meta de economia
	return s.savingGoalRepo.Delete(ctx, savingGoalID)
}

func (s *savingGoalServiceImpl) Deposit(ctx context.Context, userID, savingGoalID uint, amount float64) (*entities.SavingGoal, error) {
	// Validações
	if amount <= 0 {
		return nil, pkgErrors.NewDomainError("validation_error", "Valor do depósito deve ser maior que zero")
	}

	// Buscar meta de economia
	savingGoal, err := s.GetSavingGoalByID(ctx, userID, savingGoalID)
	if err != nil {
		return nil, err
	}

	// Log de depuração
	log.Printf("Depósito - ID: %d, Valor atual: %.2f, Depósito: %.2f",
		savingGoalID, savingGoal.CurrentAmount, amount)

	// Calcular novo valor
	newAmount := savingGoal.CurrentAmount + amount

	// Atualizar valor no banco de dados
	if err := s.savingGoalRepo.(repositories.SavingGoalRepository).UpdateAmount(ctx, savingGoalID, newAmount); err != nil {
		log.Printf("Erro ao atualizar valor: %v", err)
		return nil, err
	}

	// Log após atualização
	log.Printf("Depósito concluído - ID: %d, Novo valor: %.2f", savingGoalID, newAmount)

	// Atualizar valor na entidade
	savingGoal.CurrentAmount = newAmount
	savingGoal.UpdatedAt = time.Now()

	return savingGoal, nil
}
