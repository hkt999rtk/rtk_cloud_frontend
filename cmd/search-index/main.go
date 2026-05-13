package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"realtek-connect/internal/search"

	_ "modernc.org/sqlite"
)

func main() {
	logger := log.New(os.Stdout, "", log.LstdFlags)
	if err := run(logger); err != nil {
		logger.Printf("search index failed: %v", err)
		os.Exit(1)
	}
}

func run(logger *log.Logger) error {
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
		return err
	}
	chunks, err := search.BuildIndexChunks(ctx, documents, embedder)
	if err != nil {
		return err
	}

	dbPath := *databasePath
	if !filepath.IsAbs(dbPath) {
		dbPath = filepath.Join(root, dbPath)
	}
	if err := os.MkdirAll(filepath.Dir(dbPath), 0o755); err != nil {
		return err
	}

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return err
	}
	defer db.Close()

	repository := search.NewRepository(db)
	if err := repository.Init(ctx); err != nil {
		return err
	}
	if err := repository.Replace(ctx, chunks); err != nil {
		return err
	}

	seenDocs := map[string]bool{}
	for _, chunk := range chunks {
		seenDocs[chunk.Document.ID] = true
	}
	logger.Printf("indexed %d documents and %d chunks into %s", len(seenDocs), len(chunks), dbPath)
	return nil
}

func envOrDefault(key, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	return value
}
