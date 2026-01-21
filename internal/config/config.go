package config

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

const (
	defaultPort      = "8080"
	defaultAddress   = "0.0.0.0"
	defaultDBHost    = "localhost"
	defaultDBPort    = "5432"
	defaultDBUser    = "fiozap"
	defaultDBPass    = "fiozap123"
	defaultDBName    = "fiozap"
	defaultDBSSLMode = "disable"
	defaultLogLevel  = "info"
	defaultLogType   = "console"
	tokenLength      = 16
)

type Config struct {
	Port       string
	Address    string
	AdminToken string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
	LogLevel   string
	LogType    string
	WADebug    string
	CORSOrigin string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		Port:       getEnv("PORT", defaultPort),
		Address:    getEnv("ADDRESS", defaultAddress),
		AdminToken: getEnv("ADMIN_TOKEN", ""),
		DBHost:     getEnv("DB_HOST", defaultDBHost),
		DBPort:     getEnv("DB_PORT", defaultDBPort),
		DBUser:     getEnv("DB_USER", defaultDBUser),
		DBPassword: getEnv("DB_PASSWORD", defaultDBPass),
		DBName:     getEnv("DB_NAME", defaultDBName),
		DBSSLMode:  getEnv("DB_SSLMODE", defaultDBSSLMode),
		LogLevel:   getEnv("LOG_LEVEL", defaultLogLevel),
		LogType:    getEnv("LOG_TYPE", defaultLogType),
		WADebug:    getEnv("WA_DEBUG", ""),
		CORSOrigin: getEnv("CORS_ORIGIN", "*"),
	}

	if cfg.AdminToken == "" {
		cfg.AdminToken = generateToken()
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) Validate() error {
	if c.DBHost == "" {
		return errors.New("DB_HOST is required")
	}
	if c.DBName == "" {
		return errors.New("DB_NAME is required")
	}
	if c.DBUser == "" {
		return errors.New("DB_USER is required")
	}
	return nil
}

func (c *Config) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName, c.DBSSLMode,
	)
}

func (c *Config) PostgresURL() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName, c.DBSSLMode,
	)
}

func (c *Config) ServerAddr() string {
	return c.Address + ":" + c.Port
}

func (c *Config) IsPrettyLog() bool {
	return c.LogType == "console"
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func generateToken() string {
	b := make([]byte, tokenLength)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
