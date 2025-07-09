package dto

import (
	"my-finance-hub-api/internal/domain/entities"
	"time"
)

// Request DTOs
type CreateTransactionRequest struct {
	Description    string                   `json:"description" binding:"required,max=255"`
	Amount         float64                  `json:"amount" binding:"required,gt=0"`
	Type           entities.TransactionType `json:"type" binding:"required"`
	Date           time.Time                `json:"date" binding:"required"`
	CategoryID     *uint                    `json:"category_id"`
	PiggyBankID    *uint                    `json:"piggy_bank_id"`
	Paid           bool                     `json:"paid"`
	IsRecurrent    bool                     `json:"is_recurrent"`
	RecurrenceType entities.RecurrenceType  `json:"recurrence_type"`
	RecurrenceEnd  *time.Time               `json:"recurrence_end"`
}

type UpdateTransactionRequest struct {
	Description    string                   `json:"description" binding:"required,max=255"`
	Amount         float64                  `json:"amount" binding:"required,gt=0"`
	Type           entities.TransactionType `json:"type" binding:"required"`
	Date           time.Time                `json:"date" binding:"required"`
	CategoryID     *uint                    `json:"category_id"`
	PiggyBankID    *uint                    `json:"piggy_bank_id"`
	Paid           bool                     `json:"paid"`
	IsRecurrent    bool                     `json:"is_recurrent"`
	RecurrenceType entities.RecurrenceType  `json:"recurrence_type"`
	RecurrenceEnd  *time.Time               `json:"recurrence_end"`
}

type TransactionFiltersRequest struct {
	Paid       *bool   `form:"paid"`
	Month      *int    `form:"month"`
	Year       *int    `form:"year"`
	Type       *string `form:"type"`
	CategoryID *uint   `form:"category_id"`
}

// Response DTOs
type TransactionResponse struct {
	ID             uint                     `json:"id"`
	Description    string                   `json:"description"`
	Amount         float64                  `json:"amount"`
	Type           entities.TransactionType `json:"type"`
	Date           time.Time                `json:"date"`
	CategoryID     *uint                    `json:"category_id"`
	PiggyBankID    *uint                    `json:"piggy_bank_id"`
	UserID         uint                     `json:"user_id"`
	ParentID       *uint                    `json:"parent_id"`
	Paid           bool                     `json:"paid"`
	IsRecurrent    bool                     `json:"is_recurrent"`
	RecurrenceType entities.RecurrenceType  `json:"recurrence_type"`
	RecurrenceEnd  *time.Time               `json:"recurrence_end"`
	CreatedAt      time.Time                `json:"created_at"`
	UpdatedAt      time.Time                `json:"updated_at"`
}

type TransactionStatsResponse struct {
	TotalIncome    float64 `json:"total_income"`
	TotalExpense   float64 `json:"total_expense"`
	Balance        float64 `json:"balance"`
	MonthlyIncome  float64 `json:"monthly_income"`
	MonthlyExpense float64 `json:"monthly_expense"`
	MonthlyBalance float64 `json:"monthly_balance"`
}

type ReportResponse struct {
	TotalIncome      float64                  `json:"total_income"`
	TotalExpense     float64                  `json:"total_expense"`
	Balance          float64                  `json:"balance"`
	Transactions     []TransactionResponse    `json:"transactions"`
	CategoryStats    []map[string]interface{} `json:"category_stats"`
	MonthlyStats     []map[string]interface{} `json:"monthly_stats"`
	YearlyComparison []map[string]interface{} `json:"yearly_comparison"`
}

// Mappers
func ToTransactionResponse(transaction *entities.Transaction) TransactionResponse {
	return TransactionResponse{
		ID:             transaction.ID,
		Description:    transaction.Description,
		Amount:         transaction.Amount,
		Type:           transaction.Type,
		Date:           transaction.Date,
		CategoryID:     transaction.CategoryID,
		PiggyBankID:    transaction.PiggyBankID,
		UserID:         transaction.UserID,
		ParentID:       transaction.ParentID,
		Paid:           transaction.Paid,
		IsRecurrent:    transaction.IsRecurrent,
		RecurrenceType: transaction.RecurrenceType,
		RecurrenceEnd:  transaction.RecurrenceEnd,
		CreatedAt:      transaction.CreatedAt,
		UpdatedAt:      transaction.UpdatedAt,
	}
}

func ToTransactionResponseList(transactions []*entities.Transaction) []TransactionResponse {
	result := make([]TransactionResponse, len(transactions))
	for i, transaction := range transactions {
		result[i] = ToTransactionResponse(transaction)
	}
	return result
}

func (req *CreateTransactionRequest) ToEntity(userID uint) *entities.Transaction {
	transaction := entities.NewTransaction(req.Description, req.Amount, req.Type, req.Date, userID)

	if req.CategoryID != nil {
		transaction.SetCategory(*req.CategoryID)
	}

	if req.PiggyBankID != nil {
		transaction.SetPiggyBank(*req.PiggyBankID)
	}

	if req.IsRecurrent {
		transaction.SetRecurrence(req.RecurrenceType, req.RecurrenceEnd)
	}

	transaction.Paid = req.Paid

	return transaction
}

func (req *UpdateTransactionRequest) ToEntity(userID uint) *entities.Transaction {
	transaction := entities.NewTransaction(req.Description, req.Amount, req.Type, req.Date, userID)

	if req.CategoryID != nil {
		transaction.SetCategory(*req.CategoryID)
	}

	if req.PiggyBankID != nil {
		transaction.SetPiggyBank(*req.PiggyBankID)
	}

	if req.IsRecurrent {
		transaction.SetRecurrence(req.RecurrenceType, req.RecurrenceEnd)
	}

	transaction.Paid = req.Paid

	return transaction
}
