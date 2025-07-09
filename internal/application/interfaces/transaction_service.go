package interfaces

import (
	"context"
	"my-finance-hub-api/internal/domain/entities"
	"my-finance-hub-api/internal/domain/repositories"
	"time"
)

type TransactionService interface {
	CreateTransaction(ctx context.Context, userID uint, transaction *entities.Transaction) (*entities.Transaction, error)
	GetTransactionByID(ctx context.Context, userID, transactionID uint) (*entities.Transaction, error)
	GetTransactionsByUser(ctx context.Context, userID uint, filters *repositories.TransactionFilters) ([]*entities.Transaction, error)
	UpdateTransaction(ctx context.Context, userID, transactionID uint, updates *entities.Transaction) (*entities.Transaction, error)
	DeleteTransaction(ctx context.Context, userID, transactionID uint) error
	TogglePaidStatus(ctx context.Context, userID, transactionID uint) (*entities.Transaction, error)
	// GetTransactionStats obtém estatísticas de transações com suporte a filtros opcionais
	GetTransactionStats(ctx context.Context, userID uint, filters *repositories.TransactionFilters) (map[string]interface{}, error)
	GetReports(ctx context.Context, userID uint, filters *repositories.TransactionFilters) (map[string]interface{}, error)
	GetDashboardReports(ctx context.Context, userID uint) (map[string]interface{}, error)
	CreateRecurringTransactions(ctx context.Context, transaction *entities.Transaction) error
	GetMonthlyStats(ctx context.Context, userID uint, year int, startDate, endDate *time.Time) ([]repositories.MonthlyStats, error)
}
