package leads

import (
	"context"
	"database/sql"
	"strings"
)

type Lead struct {
	Name     string
	Company  string
	Email    string
	Interest string
	Message  string
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
