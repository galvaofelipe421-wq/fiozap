package main

import (
	"context"
	"errors"
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

const (
	version              = "1.0"
	reconnectDelay       = 2 * time.Second
	shutdownTimeout      = 5 * time.Second
	serverReadTimeout    = 60 * time.Second
	serverWriteTimeout   = 120 * time.Second
	serverIdleTimeout    = 180 * time.Second
)

func main() {
	if err := run(); err != nil {
		logger.WithError(err).Msg("application error")
		os.Exit(1)
	}
}

func run() error {
	ctx := context.Background()

	cfg, err := config.Load()
	if err != nil {
		return err
	}

	logger.Init(cfg.LogLevel, cfg.LogType == "console")
	logger.Component("server").Str("version", version).Msg("starting")

	if err := initDatabase(ctx, cfg); err != nil {
		return err
	}

	db, err := database.Connect(cfg)
	if err != nil {
		return err
	}
	defer func() { _ = db.Close() }()

	r := router.New(cfg, db)
	r.StartDispatcher()
	defer r.StopDispatcher()

	go scheduleReconnect(ctx, r)

	srv := newServer(cfg, r)

	return startServer(srv)
}

func initDatabase(ctx context.Context, cfg *config.Config) error {
	dbUtil, err := database.ConnectDBUtil(cfg)
	if err != nil {
		return err
	}
	defer func() { _ = dbUtil.Close() }()

	return migration.Run(ctx, dbUtil)
}

func newServer(cfg *config.Config, handler http.Handler) *http.Server {
	return &http.Server{
		Addr:         cfg.Address + ":" + cfg.Port,
		Handler:      handler,
		ReadTimeout:  serverReadTimeout,
		WriteTimeout: serverWriteTimeout,
		IdleTimeout:  serverIdleTimeout,
	}
}

func scheduleReconnect(ctx context.Context, r *router.Router) {
	time.Sleep(reconnectDelay)
	r.GetSessionService().ReconnectAll(ctx)
}

func startServer(srv *http.Server) error {
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	errCh := make(chan error, 1)
	go func() {
		logger.Component("server").Str("addr", srv.Addr).Msg("listening")
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	select {
	case err := <-errCh:
		return err
	case <-done:
		logger.Component("server").Msg("shutting down")
	}

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.WithError(err).Msg("shutdown error")
	}

	logger.Component("server").Msg("stopped")
	return nil
}
