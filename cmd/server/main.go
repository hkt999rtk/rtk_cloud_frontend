package main

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"realtek-connect/internal/leads"
	"realtek-connect/internal/web"

	_ "modernc.org/sqlite"
)

func main() {
	logger := log.New(os.Stdout, "", log.LstdFlags)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := run(ctx, logger); err != nil {
		logger.Printf("server exited with error: %v", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, logger *log.Logger) error {
	if logger == nil {
		logger = log.New(os.Stdout, "", log.LstdFlags)
	}

	databasePath := envOrDefault("DATABASE_PATH", "data/connectplus.db")
	if err := os.MkdirAll(filepath.Dir(databasePath), 0o755); err != nil {
		return err
	}

	db, err := sql.Open("sqlite", databasePath)
	if err != nil {
		return err
	}
	defer db.Close()

	if err := db.PingContext(ctx); err != nil {
		return err
	}

	repository := leads.NewRepository(db)
	if err := repository.Init(); err != nil {
		return err
	}

	application, err := web.NewServer(web.Config{
		LeadStore:  repository,
		AdminToken: os.Getenv("ADMIN_TOKEN"),
	})
	if err != nil {
		return err
	}

	address := ":" + envOrDefault("PORT", "8080")
	server := web.NewHTTPServer(web.HTTPServerConfig{
		Addr:    address,
		Handler: web.LoggingMiddleware(logger)(application.Routes()),
	})

	logger.Printf("listening on %s", address)

	return serveWithGracefulShutdown(ctx, server, logger, web.DefaultShutdownTimeout)
}

type httpServerLifecycle interface {
	ListenAndServe() error
	Shutdown(context.Context) error
}

func serveWithGracefulShutdown(ctx context.Context, server httpServerLifecycle, logger *log.Logger, shutdownTimeout time.Duration) error {
	if logger == nil {
		logger = log.New(os.Stdout, "", log.LstdFlags)
	}

	errCh := make(chan error, 1)
	go func() {
		errCh <- server.ListenAndServe()
	}()

	select {
	case err := <-errCh:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}
		return nil
	case <-ctx.Done():
		logger.Printf("shutdown requested")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		return err
	}

	if err := <-errCh; err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	logger.Printf("shutdown complete")
	return nil
}

func envOrDefault(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
