package database

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.mau.fi/util/dbutil"

	"fiozap/internal/config"
	"fiozap/internal/database/migration"
	"fiozap/internal/logger"
)

const (
	driverPostgres = "postgres"
	dbOwner        = "fiozap"
)

func Connect(cfg *config.Config) (*sqlx.DB, error) {
	db, err := sqlx.Open(driverPostgres, cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("failed to open postgres: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping postgres: %w", err)
	}

	logConnection(cfg)
	return db, nil
}

func ConnectDBUtil(cfg *config.Config) (*dbutil.Database, error) {
	db, err := dbutil.NewWithDialect(cfg.PostgresURL(), driverPostgres)
	if err != nil {
		return nil, fmt.Errorf("failed to open postgres: %w", err)
	}

	db.Owner = dbOwner
	db.VersionTable = migration.VersionTable

	logConnection(cfg)
	return db, nil
}

func logConnection(cfg *config.Config) {
	logger.Component("database").
		Str("host", cfg.DBHost).
		Str("port", cfg.DBPort).
		Str("name", cfg.DBName).
		Msg("connected")
}
