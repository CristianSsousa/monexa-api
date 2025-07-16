package services

import (
	"context"
	"my-finance-hub-api/internal/application/interfaces"
	"my-finance-hub-api/internal/domain/entities"
	"my-finance-hub-api/internal/domain/repositories"
	pkgErrors "my-finance-hub-api/pkg/errors"
	"time"
)

type categoryServiceImpl struct {
	categoryRepo repositories.CategoryRepository
}

func NewCategoryService(categoryRepo repositories.CategoryRepository) interfaces.CategoryService {
	return &categoryServiceImpl{
		categoryRepo: categoryRepo,
	}
}

func (s *categoryServiceImpl) CreateCategory(ctx context.Context, userID uint, category *entities.Category) (*entities.Category, error) {
	// Validações
	if category.Name == "" {
		return nil, pkgErrors.NewDomainError("validation_error", "Nome da categoria é obrigatório")
	}

	if category.Color == "" {
		return nil, pkgErrors.NewDomainError("validation_error", "Cor da categoria é obrigatória")
	}

	// Verificar se já existe uma categoria com o mesmo nome para o usuário
	exists, err := s.categoryRepo.ExistsByName(ctx, userID, category.Name)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, pkgErrors.NewDomainError("already_exists", "Já existe uma categoria com este nome")
	}

	// Criar nova categoria
	// Se nenhum UserID for fornecido, cria como categoria global (UserID = 0)
	newCategory := &entities.Category{
		Name:      category.Name,
		Color:     category.Color,
		UserID:    userID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.categoryRepo.Create(ctx, newCategory); err != nil {
		return nil, err
	}

	return newCategory, nil
}

func (s *categoryServiceImpl) GetCategoryByID(ctx context.Context, userID, categoryID uint) (*entities.Category, error) {
	category, err := s.categoryRepo.GetByID(ctx, categoryID)
	if err != nil {
		return nil, err
	}

	// Permite acesso a categorias globais (UserID = 0) ou categorias do usuário
	if category.UserID != 0 && category.UserID != userID {
		return nil, pkgErrors.ErrForbidden
	}

	return category, nil
}

func (s *categoryServiceImpl) GetCategoriesByUser(ctx context.Context, userID uint) ([]*entities.Category, error) {
	return s.categoryRepo.GetByUserID(ctx, userID)
}

func (s *categoryServiceImpl) UpdateCategory(ctx context.Context, userID, categoryID uint, updates *entities.Category) (*entities.Category, error) {
	// Buscar categoria existente
	category, err := s.GetCategoryByID(ctx, userID, categoryID)
	if err != nil {
		return nil, err
	}

	// Validações
	if updates.Name == "" {
		return nil, pkgErrors.NewDomainError("validation_error", "Nome da categoria é obrigatório")
	}

	if updates.Color == "" {
		return nil, pkgErrors.NewDomainError("validation_error", "Cor da categoria é obrigatória")
	}

	// Verificar se o novo nome já existe (se foi alterado)
	if updates.Name != category.Name {
		exists, err := s.categoryRepo.ExistsByName(ctx, userID, updates.Name)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, pkgErrors.NewDomainError("already_exists", "Já existe uma categoria com este nome")
		}
	}

	// Atualizar categoria
	category.Update(updates.Name, updates.Color, updates.Type)

	if err := s.categoryRepo.Update(ctx, category); err != nil {
		return nil, err
	}

	return category, nil
}

func (s *categoryServiceImpl) DeleteCategory(ctx context.Context, userID, categoryID uint) error {
	// Verificar se a categoria existe e pertence ao usuário
	_, err := s.GetCategoryByID(ctx, userID, categoryID)
	if err != nil {
		return err
	}

	// Excluir categoria
	return s.categoryRepo.Delete(ctx, categoryID)
}
