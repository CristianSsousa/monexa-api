package repositories

import (
	"context"
	"errors"
	"fmt"
	"log"
	"my-finance-hub-api/internal/domain/entities"
	"my-finance-hub-api/internal/domain/repositories"
	"my-finance-hub-api/internal/infrastructure/database/models"
	pkgErrors "my-finance-hub-api/pkg/errors"
	"time"

	"gorm.io/gorm"
)

type transactionRepositoryImpl struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) repositories.TransactionRepository {
	return &transactionRepositoryImpl{
		db: db,
	}
}

func (r *transactionRepositoryImpl) Create(ctx context.Context, transaction *entities.Transaction) error {
	model := &models.Transaction{}
	model.FromEntity(transaction)

	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return err
	}

	// Atualiza a entidade com o ID gerado
	transaction.ID = model.ID
	transaction.CreatedAt = model.CreatedAt
	transaction.UpdatedAt = model.UpdatedAt

	return nil
}

func (r *transactionRepositoryImpl) GetByID(ctx context.Context, id uint) (*entities.Transaction, error) {
	var model models.Transaction

	if err := r.db.WithContext(ctx).First(&model, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, pkgErrors.ErrTransactionNotFound
		}
		return nil, err
	}

	return model.ToEntity(), nil
}

func (r *transactionRepositoryImpl) GetByUserID(ctx context.Context, userID uint, filters *repositories.TransactionFilters) ([]*entities.Transaction, error) {
	query := r.db.WithContext(ctx).Where("user_id = ?", userID)

	if filters != nil {
		if filters.Paid != nil {
			query = query.Where("paid = ?", *filters.Paid)
		}
		if filters.Type != nil {
			query = query.Where("type = ?", *filters.Type)
		}
		if filters.CategoryID != nil {
			query = query.Where("category_id = ?", *filters.CategoryID)
		}
		if filters.Month != nil && filters.Year != nil {
			query = query.Where("EXTRACT(MONTH FROM date) = ? AND EXTRACT(YEAR FROM date) = ?", *filters.Month, *filters.Year)
		}

		// Novo filtro de data
		if !filters.StartDate.IsZero() && !filters.EndDate.IsZero() {
			query = query.Where("date BETWEEN ? AND ?", filters.StartDate, filters.EndDate)
		}
	}

	var models []models.Transaction
	if err := query.Order("date DESC").Find(&models).Error; err != nil {
		return nil, err
	}

	transactions := make([]*entities.Transaction, len(models))
	for i, model := range models {
		transactions[i] = model.ToEntity()
	}

	return transactions, nil
}

func (r *transactionRepositoryImpl) Update(ctx context.Context, transaction *entities.Transaction) error {
	model := &models.Transaction{}
	model.FromEntity(transaction)

	if err := r.db.WithContext(ctx).Save(model).Error; err != nil {
		return err
	}

	// Atualiza a entidade com o timestamp
	transaction.UpdatedAt = model.UpdatedAt

	return nil
}

func (r *transactionRepositoryImpl) Delete(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Delete(&models.Transaction{}, id)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return pkgErrors.ErrTransactionNotFound
	}

	return nil
}

func (r *transactionRepositoryImpl) GetByDateRange(ctx context.Context, userID uint, startDate, endDate time.Time) ([]*entities.Transaction, error) {
	var models []models.Transaction

	if err := r.db.WithContext(ctx).
		Where("user_id = ? AND date BETWEEN ? AND ?", userID, startDate, endDate).
		Order("date DESC").
		Find(&models).Error; err != nil {
		return nil, err
	}

	transactions := make([]*entities.Transaction, len(models))
	for i, model := range models {
		transactions[i] = model.ToEntity()
	}

	return transactions, nil
}

func (r *transactionRepositoryImpl) GetRecurringTransactions(ctx context.Context, userID uint) ([]*entities.Transaction, error) {
	var models []models.Transaction

	if err := r.db.WithContext(ctx).
		Where("user_id = ? AND is_recurrent = ?", userID, true).
		Find(&models).Error; err != nil {
		return nil, err
	}

	transactions := make([]*entities.Transaction, len(models))
	for i, model := range models {
		transactions[i] = model.ToEntity()
	}

	return transactions, nil
}

func (r *transactionRepositoryImpl) GetInvestmentsByPiggyBank(ctx context.Context, piggyBankID uint) ([]*entities.Transaction, error) {
	var models []models.Transaction

	if err := r.db.WithContext(ctx).
		Where("piggy_bank_id = ? AND type = ?", piggyBankID, entities.INVESTMENT).
		Find(&models).Error; err != nil {
		return nil, err
	}

	transactions := make([]*entities.Transaction, len(models))
	for i, model := range models {
		transactions[i] = model.ToEntity()
	}

	return transactions, nil
}

func (r *transactionRepositoryImpl) GetTotalAmountByType(ctx context.Context, userID uint, transactionType entities.TransactionType, startDate, endDate *time.Time) (float64, error) {
	query := r.db.WithContext(ctx).Model(&models.Transaction{}).
		Where("user_id = ? AND type = ?", userID, transactionType)

	// Adicionar filtros de data se existirem
	if startDate != nil && !startDate.IsZero() {
		query = query.Where("date >= ?", *startDate)
	}
	if endDate != nil && !endDate.IsZero() {
		query = query.Where("date <= ?", *endDate)
	}

	var total float64
	if err := query.Select("COALESCE(SUM(amount), 0)").Scan(&total).Error; err != nil {
		return 0, err
	}

	return total, nil
}

func (r *transactionRepositoryImpl) GetTotalAmountByCategory(ctx context.Context, userID uint, categoryID uint) (float64, error) {
	var total float64

	if err := r.db.WithContext(ctx).Model(&models.Transaction{}).
		Where("user_id = ? AND category_id = ? AND paid = ?", userID, categoryID, true).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&total).Error; err != nil {
		return 0, err
	}

	return total, nil
}

func (r *transactionRepositoryImpl) GetMonthlyStats(ctx context.Context, userID uint, year int, startDate, endDate *time.Time) ([]repositories.MonthlyStats, error) {
	var monthlyStats []repositories.MonthlyStats

	// Consulta SQL completa para estatísticas mensais
	query := `
		WITH month_series AS (
			SELECT 
				generate_series(1, 12) AS month_num,
				TO_CHAR(make_date(?, generate_series(1, 12), 1), 'Month') AS month_name
		),
		transaction_data AS (
			SELECT 
				EXTRACT(MONTH FROM date) AS month_num,
				type,
				SUM(amount) AS total_amount
			FROM 
				transactions
			WHERE 
				user_id = ?
				AND deleted_at IS NULL
				AND EXTRACT(YEAR FROM date) = ?
			GROUP BY 
				EXTRACT(MONTH FROM date),
				type
		)
		SELECT 
			month_series.month_name AS month,
			COALESCE(
				(SELECT total_amount FROM transaction_data 
				 WHERE month_num = month_series.month_num AND type = 'income'), 
				0
			) AS income,
			COALESCE(
				(SELECT total_amount FROM transaction_data 
				 WHERE month_num = month_series.month_num AND type = 'expense'), 
				0
			) AS expense,
			COALESCE(
				(SELECT total_amount FROM transaction_data 
				 WHERE month_num = month_series.month_num AND type = 'income'), 
				0
			) - COALESCE(
				(SELECT total_amount FROM transaction_data 
				 WHERE month_num = month_series.month_num AND type = 'expense'), 
				0
			) AS balance
		FROM 
			month_series
		ORDER BY 
			month_series.month_num
	`

	// Parâmetros para a query
	params := []interface{}{year, userID, year}

	// Executar query
	err := r.db.WithContext(ctx).Raw(query, params...).Scan(&monthlyStats).Error

	if err != nil {
		fmt.Printf("Erro ao executar consulta: %v\n", err)
		return nil, err
	}

	return monthlyStats, nil
}

func (r *transactionRepositoryImpl) GetCategoryTotals(ctx context.Context, userID uint, filters *repositories.TransactionFilters) ([]repositories.CategoryTotal, error) {
	var categoryTotals []repositories.CategoryTotal

	// Log detalhado dos filtros recebidos
	log.Printf("Filtros recebidos - UserID: %d, Filtros: %+v", userID, filters)

	// Verificar se há transações para o usuário
	var transactionCount int64
	err := r.db.WithContext(ctx).Model(&models.Transaction{}).
		Where("user_id = ?", userID).
		Count(&transactionCount).Error
	if err != nil {
		log.Printf("Erro ao contar transações: %v", err)
		return nil, err
	}
	log.Printf("Total de transações do usuário: %d", transactionCount)

	// Construir query base
	query := `
		SELECT 
			c.name AS category_name,
			SUM(t.amount) AS total,
			t.type AS type
		FROM 
			transactions t
		JOIN 
			categories c ON t.category_id = c.id
		WHERE 
			t.user_id = ?
			AND t.deleted_at IS NULL
	`

	// Parâmetros para a query
	params := []interface{}{userID}

	// Adicionar filtros de data se existirem
	if filters != nil {
		// Log de verificação de datas
		log.Printf("Datas do filtro - Início: %v, Fim: %v", filters.StartDate, filters.EndDate)

		if !filters.StartDate.IsZero() {
			query += " AND t.date >= ?"
			params = append(params, filters.StartDate)
		}
		if !filters.EndDate.IsZero() {
			query += " AND t.date <= ?"
			params = append(params, filters.EndDate)
		}

		// Filtro por tipo de transação, se especificado
		if filters.Type != nil {
			query += " AND t.type = ?"
			params = append(params, *filters.Type)
		}
	}

	// Agrupar por categoria e tipo
	query += `
		GROUP BY 
			c.name, 
			t.type
		ORDER BY 
			total DESC
	`

	// Log da query final
	log.Printf("Query final: %s", query)
	log.Printf("Parâmetros da query: %+v", params)

	// Executar query
	err = r.db.WithContext(ctx).Raw(query, params...).Scan(&categoryTotals).Error

	if err != nil {
		log.Printf("Erro ao buscar totais por categoria: %v", err)
		return nil, err
	}

	// Log dos resultados
	log.Printf("Número de categorias encontradas: %d", len(categoryTotals))
	for _, cat := range categoryTotals {
		log.Printf("Categoria: %s, Total: %.2f, Tipo: %s",
			cat.CategoryName, cat.Total, cat.Type)
	}

	return categoryTotals, nil
}
