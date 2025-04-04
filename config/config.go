package config

import (
	"ApiSmart/pkg/database"
	"os"
)

type Config struct {
	ServerPort string
	DBConfig   database.DBConfig
	JWTSecret  string
	CORSConfig CORSConfig
}

// CORSConfig contiene la configuración para CORS
type CORSConfig struct {
	AllowedOrigins   []string
	AllowedMethods   []string
	AllowedHeaders   []string
	AllowCredentials bool
	MaxAge           int
}

func LoadConfig() *Config {
	return &Config{
		ServerPort: getEnv("SERVER_PORT", "8000"),
		DBConfig: database.DBConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "3306"),
			User:     getEnv("DB_USER", "root"),
			Password: getEnv("DB_PASSWORD", "manuel"),
			DBName:   getEnv("DB_NAME", "sensores_db"),
		},
		JWTSecret: getEnv("JWT_SECRET", "secret_key_cambiar_en_produccion"),
		CORSConfig: CORSConfig{
			AllowedOrigins:   []string{"*"},                                                // Permite cualquier origen
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"}, // Métodos HTTP permitidos
			AllowedHeaders:   []string{"Content-Type", "Authorization", "X-Requested-With"},
			AllowCredentials: true,  // Permite enviar cookies en solicitudes cross-origin
			MaxAge:           86400, // Tiempo en segundos para cachear preflight requests (24 horas)
		},
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
