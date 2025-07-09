package repositories

import (
	"context"
	"my-finance-hub-api/internal/domain/entities"
	"time"
)

type TransactionFilters struct {
	UserID     uint
	Paid       *bool
	Month      *int
	Year       *int
	Type       *string
	CategoryID *uint
	StartDate  time.Time
	EndDate    time.Time
}

// MonthlyStats representa as estatísticas de transações por mês
type MonthlyStats struct {
	Month   string  `json:"month"`
	Income  float64 `json:"income"`
	Expense float64 `json:"expense"`
	Balance float64 `json:"balance"`
}

type TransactionRepository interface {
	Create(ctx context.Context, transaction *entities.Transaction) error
	GetByID(ctx context.Context, id uint) (*entities.Transaction, error)
	GetByUserID(ctx context.Context, userID uint, filters *TransactionFilters) ([]*entities.Transaction, error)
	Update(ctx context.Context, transaction *entities.Transaction) error
	Delete(ctx context.Context, id uint) error
	GetByDateRange(ctx context.Context, userID uint, startDate, endDate time.Time) ([]*entities.Transaction, error)
	GetRecurringTransactions(ctx context.Context, userID uint) ([]*entities.Transaction, error)
	GetInvestmentsByPiggyBank(ctx context.Context, piggyBankID uint) ([]*entities.Transaction, error)
	// GetTotalAmountByType busca o total de transações por tipo, com suporte a filtros de data
	GetTotalAmountByType(ctx context.Context, userID uint, transactionType entities.TransactionType, startDate, endDate *time.Time) (float64, error)
	GetTotalAmountByCategory(ctx context.Context, userID uint, categoryID uint) (float64, error)
	// GetMonthlyStats busca as estatísticas mensais de transações
	GetMonthlyStats(ctx context.Context, userID uint, year int, startDate, endDate *time.Time) ([]MonthlyStats, error)
	// GetCategoryTotals busca os totais de transações agrupados por categoria
	GetCategoryTotals(ctx context.Context, userID uint, filters *TransactionFilters) ([]CategoryTotal, error)
}

// Estrutura para representar totais por categoria
type CategoryTotal struct {
	CategoryName string
	Total        float64
	Type         string
}
