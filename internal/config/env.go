package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

var Env *Config = LoadConfig()

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Fatal("unable to load .env file")
	}

	return &Config{
		App: AppConfig{
			Port:      fmt.Sprintf(":%s", getEnv("APP_PORT", "4000")),
			Version:   getEnv("APP_VERSION", "v1"),
			RateLimit: getEnvInt("RATE_LIMIT", 5),
		},
		DB: DBConfig{
			URI:      getEnv("MONGODB_URI", "mongodb://localhost:27017/auth_db"),
			Name:     getEnv("MONGODB_NAME", "auth_db"),
			Username: getEnv("MONGODB_USERNAME", "admin"),
			Password: getEnv("MONGODB_PASSWORD", "secret_password"),
		},
		JWT: JwtConfig{
			Secret:     getEnv("JWT_SECRET", "wHRi8StvZJdGcPW9OmhpadcqLXQx5xFKQocoh+7rSIc"),
			Expiration: getEnvDuration("JWT_EXIPIRATION", 10*time.Second),
		},
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	value := getEnv(key, "")
	if value == "" {
		return fallback
	}

	valueAsInt, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}

	return valueAsInt
}

func getEnvDuration(key string, fallback time.Duration) time.Duration {
	valueAsString := getEnv(key, "")
	if valueAsString == "" {
		return fallback
	}

	valueAsDuration, err := time.ParseDuration(valueAsString)
	if err != nil {
		return fallback
	}

	return valueAsDuration
}
