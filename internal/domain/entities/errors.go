package entities

import "my-finance-hub-api/pkg/errors"

// Alias para facilitar o uso dos erros de domínio
type DomainError = errors.DomainError

var (
	NewDomainError            = errors.NewDomainError
	NewDomainErrorWithDetails = errors.NewDomainErrorWithDetails

	ErrUserNotFound       = errors.ErrUserNotFound
	ErrEmailAlreadyExists = errors.ErrEmailAlreadyExists
	ErrInvalidCredentials = errors.ErrInvalidCredentials
	ErrUnauthorized       = errors.ErrUnauthorized
	ErrForbidden          = errors.ErrForbidden

	ErrTransactionNotFound = errors.ErrTransactionNotFound
	ErrCategoryNotFound    = errors.ErrCategoryNotFound
	ErrGoalNotFound        = errors.ErrGoalNotFound
	ErrSavingGoalNotFound  = errors.ErrSavingGoalNotFound

	ErrInsufficientFunds = errors.ErrInsufficientFunds
	ErrInvalidAmount     = errors.ErrInvalidAmount
	ErrInvalidDate       = errors.ErrInvalidDate
)
