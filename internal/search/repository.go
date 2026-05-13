package search

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"math"
	"sort"
	"strings"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Init(ctx context.Context) error {
	if r == nil || r.db == nil {
		return fmt.Errorf("search repository database is required")
	}
	statements := []string{
		`CREATE TABLE IF NOT EXISTS search_documents (
			id TEXT PRIMARY KEY,
			locale TEXT NOT NULL,
			source_type TEXT NOT NULL,
			title TEXT NOT NULL,
			url TEXT NOT NULL,
			body TEXT NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS search_chunks (
			id TEXT PRIMARY KEY,
			document_id TEXT NOT NULL REFERENCES search_documents(id) ON DELETE CASCADE,
			text TEXT NOT NULL,
			embedding_json TEXT NOT NULL
		)`,
		`CREATE INDEX IF NOT EXISTS search_chunks_document_idx ON search_chunks(document_id)`,
	}
	for _, statement := range statements {
		if _, err := r.db.ExecContext(ctx, statement); err != nil {
			return err
		}
	}
	return nil
}

func (r *Repository) Replace(ctx context.Context, chunks []IndexedChunk) error {
	if r == nil || r.db == nil {
		return fmt.Errorf("search repository database is required")
	}
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	if _, err := tx.ExecContext(ctx, `DELETE FROM search_chunks`); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, `DELETE FROM search_documents`); err != nil {
		return err
	}

	seenDocs := map[string]bool{}
	for _, chunk := range chunks {
		doc := normalizeDocument(chunk.Document)
		if doc.ID == "" || chunk.ChunkID == "" || strings.TrimSpace(chunk.Text) == "" || len(chunk.Embedding) == 0 {
			return fmt.Errorf("incomplete indexed chunk")
		}
		if !seenDocs[doc.ID] {
			if _, err := tx.ExecContext(ctx, `
				INSERT INTO search_documents (id, locale, source_type, title, url, body)
				VALUES (?, ?, ?, ?, ?, ?)`,
				doc.ID, doc.Locale, doc.SourceType, doc.Title, doc.URL, doc.Body,
			); err != nil {
				return err
			}
			seenDocs[doc.ID] = true
		}
		embeddingJSON, err := json.Marshal(chunk.Embedding)
		if err != nil {
			return err
		}
		if _, err := tx.ExecContext(ctx, `
			INSERT INTO search_chunks (id, document_id, text, embedding_json)
			VALUES (?, ?, ?, ?)`,
			chunk.ChunkID, doc.ID, strings.TrimSpace(chunk.Text), string(embeddingJSON),
		); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (r *Repository) Search(ctx context.Context, query []float64, limit int) ([]Source, error) {
	if r == nil || r.db == nil {
		return nil, fmt.Errorf("search repository database is required")
	}
	if len(query) == 0 {
		return nil, fmt.Errorf("query embedding is required")
	}
	if limit <= 0 {
		limit = 5
	}
	rows, err := r.db.QueryContext(ctx, `
		SELECT d.locale, d.source_type, d.title, d.url, c.text, c.embedding_json
		FROM search_chunks c
		JOIN search_documents d ON d.id = c.document_id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := make([]Source, 0)
	for rows.Next() {
		var source Source
		var embeddingJSON string
		if err := rows.Scan(&source.Locale, &source.SourceType, &source.Title, &source.URL, &source.Snippet, &embeddingJSON); err != nil {
			return nil, err
		}
		var embedding []float64
		if err := json.Unmarshal([]byte(embeddingJSON), &embedding); err != nil {
			return nil, err
		}
		source.Score = cosine(query, embedding)
		results = append(results, source)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	sort.SliceStable(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})
	if len(results) > limit {
		results = results[:limit]
	}
	return results, nil
}

func normalizeDocument(doc Document) Document {
	doc.ID = strings.TrimSpace(doc.ID)
	doc.Locale = strings.TrimSpace(doc.Locale)
	if doc.Locale == "" {
		doc.Locale = "en"
	}
	doc.SourceType = strings.TrimSpace(doc.SourceType)
	doc.Title = strings.TrimSpace(doc.Title)
	doc.URL = strings.TrimSpace(doc.URL)
	doc.Body = strings.TrimSpace(doc.Body)
	return doc
}

func cosine(a, b []float64) float64 {
	if len(a) == 0 || len(a) != len(b) {
		return 0
	}
	var dot, normA, normB float64
	for i := range a {
		dot += a[i] * b[i]
		normA += a[i] * a[i]
		normB += b[i] * b[i]
	}
	if normA == 0 || normB == 0 {
		return 0
	}
	return dot / (math.Sqrt(normA) * math.Sqrt(normB))
}
