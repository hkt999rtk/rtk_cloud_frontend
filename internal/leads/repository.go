package leads

import (
	"context"
	"database/sql"
	"strings"
	"time"
)

type Lead struct {
	Name     string
	Company  string
	Email    string
	Interest string
	Message  string
}

type LeadRecord struct {
	ID        int64
	Name      string
	Company   string
	Email     string
	Interest  string
	Message   string
	CreatedAt time.Time
}

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Init() error {
	_, err := r.db.Exec(`
CREATE TABLE IF NOT EXISTS leads (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  name TEXT NOT NULL,
  company TEXT,
  email TEXT NOT NULL,
  interest TEXT NOT NULL,
  message TEXT,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);`)
	return err
}

func (r *Repository) Insert(ctx context.Context, lead Lead) error {
	_, err := r.db.ExecContext(ctx, `
INSERT INTO leads (name, company, email, interest, message)
VALUES (?, ?, ?, ?, ?)`,
		strings.TrimSpace(lead.Name),
		strings.TrimSpace(lead.Company),
		strings.TrimSpace(lead.Email),
		strings.TrimSpace(lead.Interest),
		strings.TrimSpace(lead.Message),
	)
	return err
}

func (r *Repository) Count(ctx context.Context) (int, error) {
	var count int
	err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM leads`).Scan(&count)
	return count, err
}

func (r *Repository) List(ctx context.Context, limit int) ([]LeadRecord, error) {
	if limit <= 0 || limit > 500 {
		limit = 100
	}

	rows, err := r.db.QueryContext(ctx, `
SELECT id, name, company, email, interest, message, created_at
FROM leads
ORDER BY id DESC
LIMIT ?`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []LeadRecord
	for rows.Next() {
		var record LeadRecord
		var createdAt string
		if err := rows.Scan(
			&record.ID,
			&record.Name,
			&record.Company,
			&record.Email,
			&record.Interest,
			&record.Message,
			&createdAt,
		); err != nil {
			return nil, err
		}
		record.CreatedAt = parseSQLiteTime(createdAt)
		records = append(records, record)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return records, nil
}

func parseSQLiteTime(value string) time.Time {
	layouts := []string{
		"2006-01-02 15:04:05",
		time.RFC3339Nano,
		time.RFC3339,
	}
	for _, layout := range layouts {
		parsed, err := time.Parse(layout, value)
		if err == nil {
			return parsed
		}
	}
	return time.Time{}
}
