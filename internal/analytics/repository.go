package analytics

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"time"

	_ "modernc.org/sqlite"
)

type Repository struct {
	db            *sql.DB
	retentionDays int
}

func Open(ctx context.Context, cfg Config) (*Repository, error) {
	if !cfg.Enabled {
		return nil, nil
	}

	if cfg.DatabasePath == "" {
		cfg.DatabasePath = DefaultDatabasePath
	}
	if cfg.RetentionDays <= 0 {
		cfg.RetentionDays = DefaultRetentionDays
	}

	if err := os.MkdirAll(filepath.Dir(cfg.DatabasePath), 0o755); err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite", cfg.DatabasePath)
	if err != nil {
		return nil, err
	}

	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, err
	}

	repository := &Repository{
		db:            db,
		retentionDays: cfg.RetentionDays,
	}
	if err := repository.Init(); err != nil {
		_ = db.Close()
		return nil, err
	}
	if _, err := repository.CleanupExpiredEvents(ctx, time.Now().UTC()); err != nil {
		_ = db.Close()
		return nil, err
	}

	return repository, nil
}

func NewRepository(db *sql.DB, retentionDays int) *Repository {
	if retentionDays <= 0 {
		retentionDays = DefaultRetentionDays
	}
	return &Repository{
		db:            db,
		retentionDays: retentionDays,
	}
}

func (r *Repository) Close() error {
	if r == nil || r.db == nil {
		return nil
	}
	return r.db.Close()
}

func (r *Repository) RetentionDays() int {
	if r == nil {
		return 0
	}
	return r.retentionDays
}

func (r *Repository) Init() error {
	_, err := r.db.Exec(`
CREATE TABLE IF NOT EXISTS analytics_events (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  ts INTEGER NOT NULL,
  event TEXT NOT NULL,
  page TEXT NOT NULL,
  cta TEXT,
  percent INTEGER,
  duration INTEGER,
  variant TEXT,
  referrer_origin TEXT,
  session_id TEXT,
  created_at TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_analytics_events_ts
  ON analytics_events(ts);

CREATE INDEX IF NOT EXISTS idx_analytics_events_event_page
  ON analytics_events(event, page);
`)
	return err
}

func (r *Repository) CleanupExpiredEvents(ctx context.Context, now time.Time) (int64, error) {
	if r == nil || r.db == nil {
		return 0, nil
	}
	if r.retentionDays <= 0 {
		return 0, nil
	}

	cutoff := now.UTC().Add(-time.Duration(r.retentionDays) * 24 * time.Hour).Unix()
	result, err := r.db.ExecContext(ctx, `DELETE FROM analytics_events WHERE ts < ?`, cutoff)
	if err != nil {
		return 0, err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return rowsAffected, nil
}
