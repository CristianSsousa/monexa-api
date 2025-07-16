package entities

import "time"

type Category struct {
	ID        uint
	Name      string
	Color     string
	UserID    uint
	Type      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewCategory creates a new Category entity
func NewCategory(name, color string, userID uint, types string) *Category {
	return &Category{
		Name:      name,
		Color:     color,
		UserID:    userID,
		Type:      types,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// Update atualiza os dados da categoria
func (c *Category) Update(name, color, types string) {
	c.Name = name
	c.Color = color
	c.Type = types
	c.UpdatedAt = time.Now()
}

// BelongsToUser verifica se a categoria pertence ao usuário
func (c *Category) BelongsToUser(userID uint) bool {
	return c.UserID == userID
}
