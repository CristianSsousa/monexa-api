package main

import (
	"log"

	"my-finance-hub-api/config"
	"my-finance-hub-api/internal/infrastructure/container"
	"my-finance-hub-api/internal/infrastructure/database"
	"my-finance-hub-api/internal/infrastructure/http/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Carregar variáveis de ambiente
	if err := godotenv.Load(); err != nil {
		log.Println("Aviso: arquivo .env não encontrado, usando variáveis de ambiente do sistema")
	}

	// Carregar configuração
	cfg := config.Load()

	// Inicializar banco de dados
	if err := database.Initialize(); err != nil {
		log.Fatal("Falha ao conectar com o banco de dados: ", err)
	}

	// Setup do banco (migrações, índices, seeds)
	if err := database.SetupDatabase(); err != nil {
		log.Fatal("Falha no setup do banco de dados: ", err)
	}

	// Criar container de dependências
	container := container.NewContainer(database.DB)

	// Configurar router
	router := gin.Default()

	// Configurar CORS
	corsConfig := cors.Config{
		AllowOrigins:     cfg.Server.AllowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}
	router.Use(cors.New(corsConfig))

	// Configurar rotas
	routes.SetupRoutes(router, container)

	// Iniciar servidor
	port := cfg.Server.Port
	log.Printf("🚀 Servidor iniciado na porta %s", port)
	log.Printf("📍 Endpoint: http://localhost:%s", port)

	if err := router.Run(":" + port); err != nil {
		log.Fatal("Falha ao iniciar servidor: ", err)
	}
}
