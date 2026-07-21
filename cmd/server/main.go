package main

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"realtek-connect/internal/analytics"
	"realtek-connect/internal/leads"
	"realtek-connect/internal/search"
	"realtek-connect/internal/web"

	cloudlogger "github.com/hkt999rtk/rtk_cloud_logger"
	"go.uber.org/zap"
	_ "modernc.org/sqlite"
)

func main() {
	logger, err := newServerLogger()
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = logger.Sync()
	}()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := run(ctx, logger); err != nil {
		logger.Error("server exited with error", zap.Error(err))
		os.Exit(1)
	}
}

func run(ctx context.Context, logger *zap.Logger) error {
	if logger == nil {
		logger = zap.NewNop()
	}

	databasePath := envOrDefault("DATABASE_PATH", "/var/lib/realtek-connect/connectplus.db")
	if !filepath.IsAbs(databasePath) {
		databasePath = filepath.Clean(filepath.Join("/var/lib/realtek-connect", databasePath))
	}
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

	analyticsStore, err := analytics.Open(ctx, analytics.ConfigFromEnv())
	if err != nil {
		return err
	}
	if analyticsStore != nil {
		defer analyticsStore.Close()
	}

	var searchDB *sql.DB
	var searchService web.SearchService
	searchEnabled := truthyEnv("SEARCH_ENABLED")
	if searchEnabled {
		openAIAPIKey := os.Getenv("OPENAI_API_KEY")
		if strings.TrimSpace(openAIAPIKey) == "" {
			logger.Warn(
				"documentation search disabled",
				zap.String("component", "search"),
				zap.String("reason", "missing_openai_api_key"),
			)
			searchEnabled = false
		}
	}
	if searchEnabled {
		searchDatabasePath := envOrDefault("SEARCH_DATABASE_PATH", "data/search.db")
		if !filepath.IsAbs(searchDatabasePath) {
			searchDatabasePath = filepath.Clean(searchDatabasePath)
		}
		if err := os.MkdirAll(filepath.Dir(searchDatabasePath), 0o755); err != nil {
			return err
		}
		searchDB, err = sql.Open("sqlite", searchDatabasePath)
		if err != nil {
			return err
		}
		defer searchDB.Close()
		if err := searchDB.PingContext(ctx); err != nil {
			return err
		}
		searchRepository := search.NewRepository(searchDB)
		if err := searchRepository.Init(ctx); err != nil {
			return err
		}
		openAIClient := search.OpenAIClient{
			APIKey:         os.Getenv("OPENAI_API_KEY"),
			EmbeddingModel: envOrDefault("SEARCH_EMBEDDING_MODEL", "text-embedding-3-small"),
			AnswerModel:    envOrDefault("SEARCH_ANSWER_MODEL", "gpt-4.1-mini"),
		}
		searchService = search.NewService(searchRepository, openAIClient, openAIClient, search.ServiceConfig{
			MaxSources: 5,
			MinScore:   0.35,
		})
	}

	application, err := web.NewServer(web.Config{
		LeadStore:               repository,
		AnalyticsStore:          analyticsStore,
		AdminToken:              os.Getenv("ADMIN_TOKEN"),
		DisableSearchIndexing:   truthyEnv("DISABLE_SEARCH_INDEXING"),
		PublicBaseURL:           os.Getenv("PUBLIC_BASE_URL"),
		EnableAssetFingerprints: truthyEnv("ENABLE_ASSET_FINGERPRINTS"),
		EnableCDNCacheHeaders:   truthyEnv("ENABLE_CDN_CACHE_HEADERS"),
		SearchEnabled:           searchEnabled,
		SearchService:           searchService,
		SDKDocsDir:              envOrDefault("SDK_DOCS_DIR", filepath.Join("dist", "sdk-docs", "current")),
	})
	if err != nil {
		logger.Error(
			"content load failed",
			zap.String("component", "content-load"),
			zap.String("error_category", "content_load_failed"),
			zap.Error(err),
		)
		return err
	}

	address := ":" + envOrDefault("PORT", "8080")
	server := web.NewHTTPServer(web.HTTPServerConfig{
		Addr:    address,
		Handler: web.LoggingMiddleware(logger)(application.Routes()),
	})

	logger.Info("listening", zap.String("addr", address))

	return serveWithGracefulShutdown(ctx, server, logger, web.DefaultShutdownTimeout)
}

type httpServerLifecycle interface {
	ListenAndServe() error
	Shutdown(context.Context) error
}

func serveWithGracefulShutdown(ctx context.Context, server httpServerLifecycle, logger *zap.Logger, shutdownTimeout time.Duration) error {
	if logger == nil {
		logger = zap.NewNop()
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
		logger.Info("shutdown requested")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		return err
	}

	if err := <-errCh; err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	logger.Info("shutdown complete")
	return nil
}

func newServerLogger() (*zap.Logger, error) {
	logger, err := cloudlogger.New(cloudlogger.Config{
		Service: "realtek-connect",
		Env:     envFirst("REALTEK_CONNECT_ENV", "APP_ENV", "ENVIRONMENT"),
		Version: serviceVersion(),
		Level:   envOrDefault("LOG_LEVEL", "info"),
	})
	if err != nil {
		return nil, err
	}
	return logger.With(zap.String("component", "server")), nil
}

func serviceVersion() string {
	if version := envFirst("REALTEK_CONNECT_VERSION", "SERVICE_VERSION", "VERSION"); version != "" {
		return version
	}
	if contents, err := os.ReadFile("VERSION"); err == nil {
		if version := strings.TrimSpace(string(contents)); version != "" {
			return version
		}
	}
	return "dev"
}

func envFirst(keys ...string) string {
	for _, key := range keys {
		if value := strings.TrimSpace(os.Getenv(key)); value != "" {
			return value
		}
	}
	return ""
}

func envOrDefault(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func truthyEnv(key string) bool {
	switch strings.ToLower(strings.TrimSpace(os.Getenv(key))) {
	case "1", "true", "yes", "on":
		return true
	default:
		return false
	}
}
