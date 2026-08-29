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

type SDKDownloadAcceptance struct {
	AcceptedAt   time.Time
	SessionID    string
	TermsVersion string
	SDKVersion   string
	Package      string
	RequestID    string
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
	if cfg.UnsafeDisableSync {
		if _, err := db.ExecContext(ctx, `PRAGMA synchronous = OFF; PRAGMA journal_mode = MEMORY;`); err != nil {
			_ = db.Close()
			return nil, err
		}
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

CREATE TABLE IF NOT EXISTS sdk_download_acceptances (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  accepted_at TEXT NOT NULL,
  session_id TEXT NOT NULL,
  terms_version TEXT NOT NULL,
  sdk_version TEXT NOT NULL,
  package TEXT NOT NULL,
  request_id TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_sdk_download_acceptances_time
  ON sdk_download_acceptances(accepted_at);
`)
	return err
}

func (r *Repository) InsertSDKDownloadAcceptance(ctx context.Context, acceptance SDKDownloadAcceptance) error {
	if r == nil || r.db == nil {
		return nil
	}
	_, err := r.db.ExecContext(ctx, `
INSERT INTO sdk_download_acceptances (
  accepted_at, session_id, terms_version, sdk_version, package, request_id
) VALUES (?, ?, ?, ?, ?, ?)`,
		acceptance.AcceptedAt.UTC().Format(time.RFC3339),
		acceptance.SessionID,
		acceptance.TermsVersion,
		acceptance.SDKVersion,
		acceptance.Package,
		acceptance.RequestID,
	)
	return err
}

func (r *Repository) InsertEvent(ctx context.Context, event Event) error {
	if r == nil || r.db == nil {
		return nil
	}

	_, err := r.db.ExecContext(ctx, `
INSERT INTO analytics_events (
  ts,
  event,
  page,
  cta,
  percent,
  duration,
  variant,
  referrer_origin,
  session_id,
  created_at
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		event.TS,
		event.Type,
		event.Page,
		nullString(event.CTA),
		nullInt(event.Percent),
		nullInt(event.Duration),
		nullString(event.Variant),
		nullString(event.ReferrerOrigin),
		event.SessionID,
		createdAtValue(event.CreatedAt),
	)
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
	acceptanceResult, err := r.db.ExecContext(ctx, `DELETE FROM sdk_download_acceptances WHERE accepted_at < ?`, time.Unix(cutoff, 0).UTC().Format(time.RFC3339))
	if err != nil {
		return 0, err
	}
	acceptanceRows, err := acceptanceResult.RowsAffected()
	if err != nil {
		return 0, err
	}
	return rowsAffected + acceptanceRows, nil
}

func nullString(value string) any {
	if value == "" {
		return nil
	}
	return value
}

func nullInt(value *int) any {
	if value == nil {
		return nil
	}
	return int64(*value)
}

func createdAtValue(value time.Time) string {
	if value.IsZero() {
		return ""
	}
	return value.UTC().Format(time.RFC3339)
}
