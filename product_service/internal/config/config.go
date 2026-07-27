package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	BDName     string
	APIPort    string
	GRPCPort   string
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("Aviso: no se encontró el archivo .env, se usarán las variables del sistema")
	}

	return &Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "password"),
		BDName:     getEnv("BD_NAME", "bd_product_service"),
		APIPort:    getEnv("API_PORT", "8081"),
		GRPCPort:   getEnv("GRPC_PORT", "50052"),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
