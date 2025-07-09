package entities

import (
	"time"
)

type TransactionType string
type RecurrenceType string

const (
	INCOME     TransactionType = "income"
	EXPENSE    TransactionType = "expense"
	INVESTMENT TransactionType = "investment"
)

const (
	NONE    RecurrenceType = "none"
	DAILY   RecurrenceType = "daily"
	WEEKLY  RecurrenceType = "weekly"
	MONTHLY RecurrenceType = "monthly"
	YEARLY  RecurrenceType = "yearly"
)

type Transaction struct {
	ID             uint
	Description    string
	Amount         float64
	Type           TransactionType
	Date           time.Time
	CategoryID     *uint
	PiggyBankID    *uint
	UserID         uint
	ParentID       *uint
	Paid           bool
	IsRecurrent    bool
	RecurrenceType RecurrenceType
	RecurrenceEnd  *time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// NewTransaction creates a new Transaction entity
func NewTransaction(description string, amount float64, transactionType TransactionType, date time.Time, userID uint) *Transaction {
	return &Transaction{
		Description:    description,
		Amount:         amount,
		Type:           transactionType,
		Date:           date,
		UserID:         userID,
		Paid:           false,
		IsRecurrent:    false,
		RecurrenceType: NONE,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
}

// SetCategory define a categoria da transação
func (t *Transaction) SetCategory(categoryID uint) {
	t.CategoryID = &categoryID
	t.UpdatedAt = time.Now()
}

// SetPiggyBank define o cofrinho da transação
func (t *Transaction) SetPiggyBank(piggyBankID uint) {
	t.PiggyBankID = &piggyBankID
	t.UpdatedAt = time.Now()
}

// TogglePaid alterna o status de pagamento
func (t *Transaction) TogglePaid() {
	t.Paid = !t.Paid
	t.UpdatedAt = time.Now()
}

// SetRecurrence define a recorrência da transação
func (t *Transaction) SetRecurrence(recurrenceType RecurrenceType, endDate *time.Time) {
	t.IsRecurrent = true
	t.RecurrenceType = recurrenceType
	t.RecurrenceEnd = endDate
	t.UpdatedAt = time.Now()
}

// Update atualiza os dados da transação
func (t *Transaction) Update(description string, amount float64, transactionType TransactionType, date time.Time) {
	t.Description = description
	t.Amount = amount
	t.Type = transactionType
	t.Date = date
	t.UpdatedAt = time.Now()
}

// IsInvestment verifica se a transação é um investimento
func (t *Transaction) IsInvestment() bool {
	return t.Type == INVESTMENT
}

// BelongsToUser verifica se a transação pertence ao usuário
func (t *Transaction) BelongsToUser(userID uint) bool {
	return t.UserID == userID
}
