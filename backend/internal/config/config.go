package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Chain    ChainConfig
	Indexer  IndexerConfig
	Auth     AuthConfig
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

type ChainConfig struct {
	PRCURL          string
	WSURL           string
	ChainID         int64
	AuctionContract string
	StartBlock      uint64
}

type IndexerConfig struct {
	Enabled       bool
	Confirmations uint64
	BatchSize     uint64
	PollInterval  time.Duration
}

type AuthConfig struct {
	Domain         string
	URI            string
	JWTSecret      string
	AccessTokenTTL time.Duration
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
		Chain: ChainConfig{
			PRCURL:          getEnv("CHAIN_PRC_URL", "http://127.0.0.1:8545"),
			WSURL:           getEnv("CHAIN_WS_URL", "ws://127.0.0.1:8545"),
			ChainID:         getInt64Env("CHAIN_ID", 31337),
			AuctionContract: getEnv("AUCTION_CONTRACT", ""),
			StartBlock:      getUint64Env("START_BLOCK", 0),
		},
		Indexer: IndexerConfig{
			Enabled:       getBoolEnv("INDEXER_ENABLED", false),
			Confirmations: getUint64Env("INDEXER_CONFIRMATIONS", 1),
			BatchSize:     getUint64Env("INDEXER_BATCH_SIZE", 500),
			PollInterval:  getDurationEnv("INDEXER_POLL_INTERVAL", 3*time.Second),
		},
		Auth: AuthConfig{
			Domain:         getEnv("AUTH_DOMAIN", "localhost:8080"),
			URI:            getEnv("AUTH_URI", "http://localhost:8080"),
			JWTSecret:      getEnv("AUTH_JWT_SECRET", "dev-secret-change-me"),
			AccessTokenTTL: getDurationEnv("AUTH_ACCESS_TOKEN_TTL", 24*time.Hour),
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

func getBoolEnv(key string, fallback bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	v, err := strconv.ParseBool(value)
	if err != nil {
		return fallback
	}
	return v
}

func getInt64Env(key string, fallback int64) int64 {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	v, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return fallback
	}

	return v
}

func getUint64Env(key string, fallback uint64) uint64 {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	v, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return fallback
	}

	return v
}
