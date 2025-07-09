package config

import (
	"os"
	"strconv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
}

type ServerConfig struct {
	Port            string
	AllowedOrigins  []string
	ReadTimeout     int
	WriteTimeout    int
	ShutdownTimeout int
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
	TimeZone string
}

type JWTConfig struct {
	Secret          string
	ExpirationHours int
	RefreshHours    int
	Issuer          string
}

func Load() *Config {
	return &Config{
		Server: ServerConfig{
			Port:            getEnv("PORT", "8080"),
			AllowedOrigins:  []string{getEnv("ALLOWED_ORIGINS", "http://localhost:5174")},
			ReadTimeout:     getEnvAsInt("READ_TIMEOUT", 30),
			WriteTimeout:    getEnvAsInt("WRITE_TIMEOUT", 30),
			ShutdownTimeout: getEnvAsInt("SHUTDOWN_TIMEOUT", 5),
		},
		Database: DatabaseConfig{
			Host:     getEnv("SUPABASE_HOST", "localhost"),
			Port:     getEnv("SUPABASE_PORT", "5432"),
			User:     getEnv("SUPABASE_USER", "postgres"),
			Password: getEnv("SUPABASE_PASSWORD", "postgres"),
			DBName:   getEnv("SUPABASE_DB_NAME", "my_finance_hub"),
			SSLMode:  getEnv("SSL_MODE", "disable"),
			TimeZone: getEnv("TIME_ZONE", "America/Sao_Paulo"),
		},
		JWT: JWTConfig{
			Secret:          getEnv("JWT_SECRET", "your-256-bit-secret"),
			ExpirationHours: getEnvAsInt("JWT_EXPIRATION_HOURS", 24),
			RefreshHours:    getEnvAsInt("JWT_REFRESH_HOURS", 168), // 7 dias
			Issuer:          getEnv("JWT_ISSUER", "my-finance-hub"),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
