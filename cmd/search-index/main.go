package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"realtek-connect/internal/search"

	cloudlogger "github.com/hkt999rtk/rtk_cloud_logger"
	"go.uber.org/zap"
	_ "modernc.org/sqlite"
)

func main() {
	logger, err := newSearchIndexRootLogger()
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = logger.Sync()
	}()

	if err := run(logger); err != nil {
		logSearchIndexError(logger, "run", "search_index_failed", err)
		os.Exit(1)
	}
}

func run(logger *zap.Logger) error {
	if logger == nil {
		logger = zap.NewNop()
	}

	repoRootDefault, err := os.Getwd()
	if err != nil {
		return err
	}

	repoRoot := flag.String("repo-root", repoRootDefault, "repository root to index")
	databasePath := flag.String("database", envOrDefault("SEARCH_DATABASE_PATH", "data/search.db"), "SQLite search database path")
	contentRoot := flag.String("content-root", "", "content root override; defaults to <repo-root>/content")
	flag.Parse()

	apiKey := strings.TrimSpace(os.Getenv("OPENAI_API_KEY"))
	if apiKey == "" {
		return fmt.Errorf("OPENAI_API_KEY is required")
	}

	root, err := filepath.Abs(*repoRoot)
	if err != nil {
		return err
	}
	contentDir := strings.TrimSpace(*contentRoot)
	if contentDir == "" {
		contentDir = filepath.Join(root, "content")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	embedder := search.OpenAIClient{
		APIKey:         apiKey,
		EmbeddingModel: envOrDefault("SEARCH_EMBEDDING_MODEL", "text-embedding-3-small"),
	}

	documents, err := search.CollectWebsiteDocuments(search.CollectionConfig{
		RepoRoot:    root,
		ContentRoot: contentDir,
	})
	if err != nil {
		logSearchIndexError(logger, "collect_content", "content_load_failed", err)
		return err
	}
	chunks, err := search.BuildIndexChunks(ctx, documents, embedder)
	if err != nil {
		logSearchIndexError(logger, "build_chunks", "embedding_failed", err)
		return err
	}

	dbPath := *databasePath
	if !filepath.IsAbs(dbPath) {
		dbPath = filepath.Join(root, dbPath)
	}
	if err := os.MkdirAll(filepath.Dir(dbPath), 0o755); err != nil {
		logSearchIndexError(logger, "prepare_database", "storage_failed", err)
		return err
	}

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		logSearchIndexError(logger, "open_database", "storage_failed", err)
		return err
	}
	defer db.Close()

	repository := search.NewRepository(db)
	if err := repository.Init(ctx); err != nil {
		logSearchIndexError(logger, "init_database", "storage_failed", err)
		return err
	}
	if err := repository.Replace(ctx, chunks); err != nil {
		logSearchIndexError(logger, "replace_index", "storage_failed", err)
		return err
	}

	seenDocs := map[string]bool{}
	for _, chunk := range chunks {
		seenDocs[chunk.Document.ID] = true
	}
	logger.Info(
		"search index complete",
		zap.Int("documents", len(seenDocs)),
		zap.Int("chunks", len(chunks)),
		zap.String("database_path", dbPath),
	)
	return nil
}

func newSearchIndexRootLogger() (*zap.Logger, error) {
	logger, err := cloudlogger.New(cloudlogger.Config{
		Service: "realtek-connect",
		Env:     envOrDefault("REALTEK_CONNECT_ENV", envOrDefault("APP_ENV", envOrDefault("ENVIRONMENT", "unknown"))),
		Version: serviceVersion(),
		Unit:    "realtek-connect-search-index.service",
		Level:   envOrDefault("LOG_LEVEL", "info"),
	})
	if err != nil {
		return nil, err
	}
	return logger.With(zap.String("component", "search-index")), nil
}

func logSearchIndexError(logger *zap.Logger, operation, category string, err error) {
	if logger == nil {
		logger = zap.NewNop()
	}
	logger.Error(
		"search index failed",
		zap.String("operation", operation),
		zap.String("error_category", category),
		zap.Error(err),
	)
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
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	return value
}
