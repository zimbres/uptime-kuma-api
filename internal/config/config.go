package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string
	ServerPort string
	AuthToken  string
	EnableAuth bool
	Version    string
}

func LoadConfig(version string) *Config {
	godotenv.Load()

	dbPort, _ := strconv.Atoi(getEnv("DB_PORT", "3306"))

	return &Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     dbPort,
		DBUser:     getEnv("DB_USER", "root"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:     getEnv("DB_NAME", "uptime"),
		ServerPort: getEnv("SERVER_PORT", "8080"),
		AuthToken:  getEnv("AUTH_TOKEN", ""),
		EnableAuth: getEnv("ENABLE_AUTH", "false") == "true",
		Version:    version,
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}