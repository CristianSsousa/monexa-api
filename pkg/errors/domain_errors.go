package errors

import (
	"fmt"
	"net/http"
)

// DomainError representa um erro específico do domínio
type DomainError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

func (e DomainError) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("%s: %s - %s", e.Code, e.Message, e.Details)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// NewDomainError cria um novo erro de domínio
func NewDomainError(code, message string) DomainError {
	return DomainError{
		Code:    code,
		Message: message,
	}
}

// NewDomainErrorWithDetails cria um novo erro de domínio com detalhes
func NewDomainErrorWithDetails(code, message, details string) DomainError {
	return DomainError{
		Code:    code,
		Message: message,
		Details: details,
	}
}

// HTTP Status mapping para diferentes tipos de erro
func (e DomainError) HTTPStatus() int {
	switch e.Code {
	case "not_found":
		return http.StatusNotFound
	case "unauthorized":
		return http.StatusUnauthorized
	case "forbidden":
		return http.StatusForbidden
	case "validation_error":
		return http.StatusBadRequest
	case "insufficient_funds":
		return http.StatusBadRequest
	case "already_exists":
		return http.StatusConflict
	case "invalid_credentials":
		return http.StatusUnauthorized
	default:
		return http.StatusInternalServerError
	}
}

// Erros comuns do domínio
var (
	ErrUserNotFound       = NewDomainError("not_found", "Usuário não encontrado")
	ErrEmailAlreadyExists = NewDomainError("already_exists", "Email já está em uso")
	ErrInvalidCredentials = NewDomainError("invalid_credentials", "Email ou senha inválidos")
	ErrUnauthorized       = NewDomainError("unauthorized", "Não autorizado")
	ErrForbidden          = NewDomainError("forbidden", "Acesso negado")

	ErrTransactionNotFound = NewDomainError("not_found", "Transação não encontrada")
	ErrCategoryNotFound    = NewDomainError("not_found", "Categoria não encontrada")
	ErrGoalNotFound        = NewDomainError("not_found", "Meta não encontrada")
	ErrSavingGoalNotFound  = NewDomainError("not_found", "Meta de economia não encontrada")

	ErrInsufficientFunds = NewDomainError("insufficient_funds", "Saldo insuficiente")
	ErrInvalidAmount     = NewDomainError("validation_error", "Valor inválido")
	ErrInvalidDate       = NewDomainError("validation_error", "Data inválida")
)
