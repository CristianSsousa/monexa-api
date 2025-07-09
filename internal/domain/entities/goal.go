package entities

import "time"

type Goal struct {
	ID          uint
	Name        string
	Description string
	Amount      float64
	UserID      uint
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// NewGoal creates a new Goal entity
func NewGoal(name, description string, amount float64, userID uint) *Goal {
	return &Goal{
		Name:        name,
		Description: description,
		Amount:      amount,
		UserID:      userID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

// Update atualiza os dados da meta
func (g *Goal) Update(name, description string, amount float64) {
	g.Name = name
	g.Description = description
	g.Amount = amount
	g.UpdatedAt = time.Now()
}

// BelongsToUser verifica se a meta pertence ao usuário
func (g *Goal) BelongsToUser(userID uint) bool {
	return g.UserID == userID
}
