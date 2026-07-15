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
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("Aviso: No se encontró el archivo .env, se usarán las variables del sistema")
	}

	return &Config{
		DBHost:     getEnv("DB_HOST", "localhost"), // "localhost" es el valor por defecto
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "password"),
		BDName:     getEnv("BD_NAME", "bd_tests"),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
