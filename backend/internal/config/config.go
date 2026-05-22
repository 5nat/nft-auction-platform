package config

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
}

type ServerConfig struct {
	Addr            string
	GinMode         string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	IdleTimeout     time.Duration
	ShutdownTimeout time.Duration
}

type DatabaseConfig struct {
	MySQLDSN string
}

func Load() (Config, error) {
	_ = godotenv.Load()

	cfg := Config{
		Server: ServerConfig{
			Addr:            getEnv("APP_ADDR", ":8080"),
			GinMode:         getEnv("GIN_MODE", "debug"),
			ReadTimeout:     getDurationEnv("SERVER_READ_TIMEOUT", 5*time.Second),
			WriteTimeout:    getDurationEnv("SERVER_WRITE_TIMEOUT", 10*time.Second),
			IdleTimeout:     getDurationEnv("SERVER_IDLE_TIMEOUT", 60*time.Second),
			ShutdownTimeout: getDurationEnv("SERVER_SHUTDOWN_TIMEOUT", 5*time.Second),
		},
		Database: DatabaseConfig{
			MySQLDSN: getEnv("MYSQL_DSN", ""),
		},
	}

	if cfg.Database.MySQLDSN == "" {
		return Config{}, fmt.Errorf("MYSQL_DSN is required")
	}

	return cfg, nil
}

func getEnv(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func getDurationEnv(key string, fallback time.Duration) time.Duration {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	d, err := time.ParseDuration(value)
	if err != nil {
		return fallback
	}

	return d
}
