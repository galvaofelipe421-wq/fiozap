package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"fiozap/internal/config"
	"fiozap/internal/database"
	"fiozap/internal/database/migration"
	"fiozap/internal/logger"
	"fiozap/internal/router"

	_ "fiozap/docs"
)

// @title FioZap API
// @version 1.0
// @description WhatsApp API using whatsmeow
// @termsOfService http://swagger.io/terms/

// @contact.name FioZap Support
// @contact.email support@fiozap.io

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Token
// @description User API token for session operations

// @securityDefinitions.apikey AdminKeyAuth
// @in header
// @name Authorization
// @description Admin token for user management

func main() {
	ctx := context.Background()

	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	logger.Init(cfg.LogLevel, cfg.LogType == "console")
	logger.Component("server").Str("version", "1.0").Msg("starting")

	dbUtil, err := database.ConnectDBUtil(cfg)
	if err != nil {
		logger.WithError(err).Msg("failed to connect to database")
		os.Exit(1)
	}
	defer dbUtil.Close()

	if err := migration.Run(ctx, dbUtil); err != nil {
		logger.WithError(err).Msg("failed to run migrations")
		os.Exit(1)
	}

	db, err := database.Connect(cfg)
	if err != nil {
		logger.WithError(err).Msg("failed to connect to database")
		os.Exit(1)
	}
	defer db.Close()

	r := router.New(cfg, db)

	r.StartDispatcher()
	defer r.StopDispatcher()

	go func() {
		time.Sleep(2 * time.Second)
		r.GetSessionService().ReconnectAll(ctx)
	}()

	srv := &http.Server{
		Addr:         cfg.Address + ":" + cfg.Port,
		Handler:      r,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 120 * time.Second,
		IdleTimeout:  180 * time.Second,
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		logger.Component("server").
			Str("addr", cfg.Address).
			Str("port", cfg.Port).
			Msg("listening")
		logger.Component("server").Str("admin_token", cfg.AdminToken).Msg("config")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.WithError(err).Msg("server error")
			os.Exit(1)
		}
	}()

	<-done
	logger.Component("server").Msg("shutting down")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.WithError(err).Msg("shutdown error")
	}

	logger.Component("server").Msg("stopped")
}
