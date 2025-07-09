package database

import (
	"fmt"
	"log"
	"os"

	"my-finance-hub-api/internal/infrastructure/database/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Initialize() error {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=America/Sao_Paulo",
		os.Getenv("SUPABASE_HOST"),
		os.Getenv("SUPABASE_USER"),
		os.Getenv("SUPABASE_PASSWORD"),
		os.Getenv("SUPABASE_DB_NAME"),
		os.Getenv("SUPABASE_PORT"),
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		// Habilitar log de SQL para depuração
		// Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Printf("Erro ao conectar com o banco de dados: %v", err)
		return err
	}

	log.Println("Conexão com o banco de dados estabelecida com sucesso")
	return nil
}

func SetupDatabase() error {
	// Migrar todos os modelos
	err := DB.AutoMigrate(
		&models.User{},
		&models.Category{},
		&models.Goal{},
		&models.SavingGoal{},
		&models.Transaction{},
	)

	if err != nil {
		log.Printf("Erro ao realizar migração: %v", err)
		return err
	}

	log.Println("Banco de dados configurado com sucesso")
	return nil
}
