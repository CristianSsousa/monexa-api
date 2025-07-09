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

type transactionServiceImpl struct {
	transactionRepo repositories.TransactionRepository
}

func NewTransactionService(transactionRepo repositories.TransactionRepository) interfaces.TransactionService {
	return &transactionServiceImpl{
		transactionRepo: transactionRepo,
	}
}

func (s *transactionServiceImpl) CreateTransaction(ctx context.Context, userID uint, transaction *entities.Transaction) (*entities.Transaction, error) {
	// Validações
	if transaction.Description == "" {
		return nil, pkgErrors.NewDomainError("validation_error", "Descrição da transação é obrigatória")
	}

	if transaction.Amount <= 0 {
		return nil, pkgErrors.NewDomainError("validation_error", "Valor da transação deve ser maior que zero")
	}

	if transaction.Type == "" {
		return nil, pkgErrors.NewDomainError("validation_error", "Tipo da transação é obrigatório")
	}

	// Criar nova transação
	newTransaction := entities.NewTransaction(
		transaction.Description,
		transaction.Amount,
		transaction.Type,
		transaction.Date,
		userID,
	)

	// Configurar campos opcionais
	if transaction.CategoryID != nil {
		newTransaction.SetCategory(*transaction.CategoryID)
	}
	if transaction.PiggyBankID != nil {
		newTransaction.SetPiggyBank(*transaction.PiggyBankID)
	}
	if transaction.IsRecurrent {
		newTransaction.SetRecurrence(transaction.RecurrenceType, transaction.RecurrenceEnd)
	}

	if err := s.transactionRepo.Create(ctx, newTransaction); err != nil {
		return nil, err
	}

	return newTransaction, nil
}

func (s *transactionServiceImpl) GetTransactionByID(ctx context.Context, userID, transactionID uint) (*entities.Transaction, error) {
	transaction, err := s.transactionRepo.GetByID(ctx, transactionID)
	if err != nil {
		return nil, err
	}

	// Verificar se a transação pertence ao usuário
	if !transaction.BelongsToUser(userID) {
		return nil, pkgErrors.ErrForbidden
	}

	return transaction, nil
}

func (s *transactionServiceImpl) GetTransactionsByUser(ctx context.Context, userID uint, filters *repositories.TransactionFilters) ([]*entities.Transaction, error) {
	if filters == nil {
		filters = &repositories.TransactionFilters{}
	}
	filters.UserID = userID

	return s.transactionRepo.GetByUserID(ctx, userID, filters)
}

func (s *transactionServiceImpl) UpdateTransaction(ctx context.Context, userID, transactionID uint, updates *entities.Transaction) (*entities.Transaction, error) {
	// Buscar transação existente
	transaction, err := s.GetTransactionByID(ctx, userID, transactionID)
	if err != nil {
		return nil, err
	}

	// Validações
	if updates.Description == "" {
		return nil, pkgErrors.NewDomainError("validation_error", "Descrição da transação é obrigatória")
	}

	if updates.Amount <= 0 {
		return nil, pkgErrors.NewDomainError("validation_error", "Valor da transação deve ser maior que zero")
	}

	if updates.Type == "" {
		return nil, pkgErrors.NewDomainError("validation_error", "Tipo da transação é obrigatório")
	}

	// Atualizar transação
	transaction.Update(updates.Description, updates.Amount, updates.Type, updates.Date)

	// Atualizar campos opcionais
	if updates.CategoryID != nil {
		transaction.SetCategory(*updates.CategoryID)
	}
	if updates.PiggyBankID != nil {
		transaction.SetPiggyBank(*updates.PiggyBankID)
	}

	if err := s.transactionRepo.Update(ctx, transaction); err != nil {
		return nil, err
	}

	return transaction, nil
}

func (s *transactionServiceImpl) DeleteTransaction(ctx context.Context, userID, transactionID uint) error {
	// Verificar se a transação existe e pertence ao usuário
	_, err := s.GetTransactionByID(ctx, userID, transactionID)
	if err != nil {
		return err
	}

	// Excluir transação
	return s.transactionRepo.Delete(ctx, transactionID)
}

func (s *transactionServiceImpl) TogglePaidStatus(ctx context.Context, userID, transactionID uint) (*entities.Transaction, error) {
	// Buscar transação
	transaction, err := s.GetTransactionByID(ctx, userID, transactionID)
	if err != nil {
		return nil, err
	}

	// Alternar status de pagamento
	transaction.TogglePaid()

	// Salvar alteração
	if err := s.transactionRepo.Update(ctx, transaction); err != nil {
		return nil, err
	}

	return transaction, nil
}

func (s *transactionServiceImpl) GetTransactionStats(ctx context.Context, userID uint, filters *repositories.TransactionFilters) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// Preparar ponteiros para datas
	var startDate, endDate *time.Time
	if filters != nil {
		startDate = &filters.StartDate
		endDate = &filters.EndDate
	}

	// Total de receitas
	totalIncome, err := s.transactionRepo.GetTotalAmountByType(ctx, userID, entities.INCOME, startDate, endDate)
	if err != nil {
		return nil, err
	}

	// Total de gastos
	totalExpense, err := s.transactionRepo.GetTotalAmountByType(ctx, userID, entities.EXPENSE, startDate, endDate)
	if err != nil {
		return nil, err
	}

	// Total de investimentos
	totalInvestment, err := s.transactionRepo.GetTotalAmountByType(ctx, userID, entities.INVESTMENT, startDate, endDate)
	if err != nil {
		return nil, err
	}

	stats["total_income"] = totalIncome
	stats["total_expense"] = totalExpense
	stats["total_investment"] = totalInvestment
	stats["balance"] = totalIncome - totalExpense

	return stats, nil
}

func (s *transactionServiceImpl) GetReports(ctx context.Context, userID uint, filters *repositories.TransactionFilters) (map[string]interface{}, error) {
	// Estrutura de relatório esperada pelo frontend
	reportData := map[string]interface{}{
		"totalIncome":    0,
		"totalExpense":   0,
		"balance":        0,
		"categoryTotals": []map[string]interface{}{},
		"monthlyTotals":  []map[string]interface{}{},
	}

	// Log de depuração para filtros recebidos
	log.Printf("Filtros recebidos - UserID: %d, Filtros: %+v", userID, filters)

	// Validar e ajustar filtros se necessário
	if filters == nil {
		filters = &repositories.TransactionFilters{}
	}

	// Se ano não foi especificado, usar o ano atual
	currentYear := time.Now().Year()
	if filters.Year == nil {
		filters.Year = &currentYear
	}

	// Se mês não foi especificado, usar o mês atual
	currentMonth := int(time.Now().Month())
	if filters.Month == nil {
		filters.Month = &currentMonth
	}

	// Definir datas de início e fim se não existirem
	if filters.StartDate.IsZero() || filters.EndDate.IsZero() {
		startDate := time.Date(*filters.Year, time.Month(*filters.Month), 1, 0, 0, 0, 0, time.UTC)
		endDate := startDate.AddDate(0, 1, -1)

		filters.StartDate = startDate
		filters.EndDate = endDate
	}

	// Log de filtros ajustados
	log.Printf("Filtros ajustados - Ano: %d, Mês: %d, Início: %v, Fim: %v",
		*filters.Year, *filters.Month, filters.StartDate, filters.EndDate)

	// Estatísticas básicas
	stats, err := s.GetTransactionStats(ctx, userID, filters)
	if err != nil {
		log.Printf("Erro ao buscar estatísticas: %v", err)
		return nil, err
	}

	// Mapear estatísticas para o formato esperado
	reportData["totalIncome"] = stats["total_income"]
	reportData["totalExpense"] = stats["total_expense"]
	reportData["balance"] = stats["balance"]

	// Buscar totais por categoria
	categoryTotals, err := s.transactionRepo.GetCategoryTotals(ctx, userID, filters)
	if err != nil {
		log.Printf("Erro ao buscar totais por categoria: %v", err)
		return nil, err
	}

	// Log de totais por categoria
	log.Printf("Número de categorias encontradas: %d", len(categoryTotals))

	// Mapear totais por categoria
	for _, cat := range categoryTotals {
		reportData["categoryTotals"] = append(
			reportData["categoryTotals"].([]map[string]interface{}),
			map[string]interface{}{
				"name":  cat.CategoryName,
				"total": cat.Total,
				"type":  cat.Type,
			},
		)
	}

	// Estatísticas mensais do ano atual
	monthlyStats, err := s.transactionRepo.GetMonthlyStats(ctx, userID, currentYear, nil, nil)
	if err != nil {
		log.Printf("Erro ao buscar estatísticas mensais: %v", err)
		return nil, err
	}

	// Mapear estatísticas mensais
	for _, monthly := range monthlyStats {
		reportData["monthlyTotals"] = append(
			reportData["monthlyTotals"].([]map[string]interface{}),
			map[string]interface{}{
				"month":   monthly.Month,
				"income":  monthly.Income,
				"expense": monthly.Expense,
				"balance": monthly.Balance,
			},
		)
	}

	return reportData, nil
}

func (s *transactionServiceImpl) GetDashboardReports(ctx context.Context, userID uint) (map[string]interface{}, error) {
	reports := make(map[string]interface{})

	// Estatísticas básicas
	stats, err := s.GetTransactionStats(ctx, userID, nil)
	if err != nil {
		return nil, err
	}
	reports["stats"] = stats

	// Transações recentes (últimas 10)
	recentTransactions, err := s.GetTransactionsByUser(ctx, userID, nil)
	if err != nil {
		return nil, err
	}

	// Limitar a 10 transações mais recentes
	if len(recentTransactions) > 10 {
		recentTransactions = recentTransactions[:10]
	}
	reports["recent_transactions"] = recentTransactions

	// Estatísticas mensais do ano atual
	currentYear := time.Now().Year()
	monthlyStats, err := s.transactionRepo.GetMonthlyStats(ctx, userID, currentYear, nil, nil)
	if err != nil {
		return nil, err
	}
	reports["monthly_stats"] = monthlyStats

	return reports, nil
}

func (s *transactionServiceImpl) CreateRecurringTransactions(ctx context.Context, transaction *entities.Transaction) error {
	// TODO: Implementar lógica para criar transações recorrentes
	// Esta funcionalidade pode ser implementada posteriormente
	return pkgErrors.NewDomainError("not_implemented", "Funcionalidade de transações recorrentes ainda não implementada")
}

func (s *transactionServiceImpl) GetMonthlyStats(ctx context.Context, userID uint, year int, startDate, endDate *time.Time) ([]repositories.MonthlyStats, error) {
	// Validar entrada
	if year == 0 {
		year = time.Now().Year()
	}

	// Buscar estatísticas mensais
	monthlyStats, err := s.transactionRepo.GetMonthlyStats(ctx, userID, year, startDate, endDate)
	if err != nil {
		return nil, err
	}

	// Garantir que todos os meses estejam presentes
	monthMap := make(map[string]repositories.MonthlyStats)
	monthOrder := []string{
		"January", "February", "March", "April",
		"May", "June", "July", "August",
		"September", "October", "November", "December",
	}

	// Preencher mapa com meses existentes
	for _, stat := range monthlyStats {
		monthMap[stat.Month] = stat
	}

	// Criar lista final com todos os meses
	var finalStats []repositories.MonthlyStats
	for _, monthName := range monthOrder {
		if stat, exists := monthMap[monthName]; exists {
			finalStats = append(finalStats, stat)
		} else {
			// Adicionar mês sem transações
			finalStats = append(finalStats, repositories.MonthlyStats{
				Month:   monthName,
				Income:  0,
				Expense: 0,
			})
		}
	}

	return finalStats, nil
}
