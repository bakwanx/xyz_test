package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type JWTConfig struct {
	Secret        string
	ExpiryMinutes int
}

type Config struct {
	MySQLDSN      string
	ServerAddress string
	JWT           JWTConfig
	UploadDir     string
	GinMode       string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}
	jwtExpiryStr := getEnv("JWT_EXPIRY_MINUTES", "60")
	expiry, err := strconv.Atoi(jwtExpiryStr)
	if err != nil {
		log.Fatalf("invalid JWT_EXPIRY_MINUTES: %v", err)
	}

	cfg := &Config{
		MySQLDSN:      getEnv("MYSQL_DSN", "user:password@tcp(localhost:3306)/xyz?parseTime=true"),
		ServerAddress: getEnv("SERVER_ADDR", ":8080"),
		JWT: JWTConfig{
			Secret:        getEnv("JWT_SECRET", "test key"),
			ExpiryMinutes: expiry,
		},
		UploadDir: getEnv("UPLOAD_DIR", "./uploads"),
		GinMode:   getEnv("GIN_MODE", "debug"),
	}

	return cfg
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
