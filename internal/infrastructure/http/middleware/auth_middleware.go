package middleware

import (
	"errors"
	"net/http"
	"os"
	"strings"

	"my-finance-hub-api/internal/application/interfaces"
	pkgErrors "my-finance-hub-api/pkg/errors"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type AuthMiddleware struct {
	authService interfaces.AuthService
}

func NewAuthMiddleware(authService interfaces.AuthService) *AuthMiddleware {
	return &AuthMiddleware{
		authService: authService,
	}
}

func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Permitir requisições de pré-vôo (OPTIONS) sem autenticação
		if c.Request.Method == http.MethodOptions {
			c.Next()
			return
		}

		// Obter o token do cabeçalho Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token não fornecido"})
			c.Abort()
			return
		}

		// Remover o prefixo "Bearer " do token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Formato de token inválido"})
			c.Abort()
			return
		}

		// Validar o token
		userID, err := m.validateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido"})
			c.Abort()
			return
		}

		// Verificar se o usuário ainda existe
		user, err := m.authService.GetUserByID(c.Request.Context(), userID)
		if err != nil {
			if _, ok := err.(pkgErrors.DomainError); ok {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não encontrado"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro interno do servidor"})
			}
			c.Abort()
			return
		}

		// Adicionar informações do usuário ao contexto
		c.Set("user_id", user.ID)
		c.Set("user", user)

		c.Next()
	}
}

func (m *AuthMiddleware) validateToken(tokenString string) (uint, error) {
	// Parse do token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Verificar o método de assinatura
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("método de assinatura inválido")
		}

		// Retornar a chave secreta
		secretKey := os.Getenv("JWT_SECRET")
		if secretKey == "" {
			secretKey = "your-256-bit-secret" // Chave padrão para desenvolvimento
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return 0, err
	}

	// Verificar se o token é válido e extrair claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Extrair user_id do token
		userIDFloat, ok := claims["user_id"].(float64)
		if !ok {
			return 0, errors.New("user_id inválido no token")
		}

		userID := uint(userIDFloat)
		if userID == 0 {
			return 0, errors.New("user_id deve ser maior que zero")
		}

		return userID, nil
	}

	return 0, errors.New("token inválido")
}

// OptionalAuth middleware que permite requests autenticados e não autenticados
func (m *AuthMiddleware) OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.Next()
			return
		}

		userID, err := m.validateToken(tokenString)
		if err != nil {
			c.Next()
			return
		}

		user, err := m.authService.GetUserByID(c.Request.Context(), userID)
		if err != nil {
			c.Next()
			return
		}

		c.Set("user_id", user.ID)
		c.Set("user", user)
		c.Next()
	}
}
