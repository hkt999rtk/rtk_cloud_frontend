package analytics

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "modernc.org/sqlite"
)

type Repository struct {
	db            *sql.DB
	retentionDays int
}

type Event struct {
	Event          string
	Page           string
	CTA            string
	Percent        int
	Duration       int
	Variant        string
	ReferrerOrigin string
	SessionID      string
}

type ReferrerMetric struct {
	Origin string
	Count  int64
}

type ScrollMilestone struct {
	Percent int
	Count   int64
}

type CTAClickMetric struct {
	Page  string
	CTA   string
	Count int64
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

func (r *Repository) InsertEvent(ctx context.Context, event Event) error {
	if r == nil || r.db == nil {
		return nil
	}
	if event.Event == "" || event.Page == "" {
		return fmt.Errorf("invalid event payload")
	}

	ts := time.Now().UTC().Unix()
	createdAt := time.Now().UTC().Format(time.RFC3339)

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
		ts,
		event.Event,
		event.Page,
		nullableString(event.CTA),
		nullableInt(event.Percent),
		nullableInt(event.Duration),
		nullableString(event.Variant),
		nullableString(event.ReferrerOrigin),
		nullableString(event.SessionID),
		createdAt,
	)
	return err
}

func (r *Repository) ConversionRate(ctx context.Context) (float64, error) {
	if r == nil || r.db == nil {
		return 0, nil
	}

	var rate sql.NullFloat64
	if err := r.db.QueryRowContext(ctx, `
SELECT
  CAST(COUNT(CASE WHEN event = 'click_cta' THEN 1 END) AS REAL) /
  NULLIF(COUNT(CASE WHEN event = 'page_view' THEN 1 END), 0)
FROM analytics_events;
`).Scan(&rate); err != nil {
		return 0, err
	}
	if !rate.Valid {
		return 0, nil
	}
	return rate.Float64, nil
}

func (r *Repository) TopReferrerOrigins(ctx context.Context, limit int) ([]ReferrerMetric, error) {
	if r == nil || r.db == nil {
		return nil, nil
	}
	if limit <= 0 {
		limit = 10
	}

	rows, err := r.db.QueryContext(ctx, `
SELECT referrer_origin, COUNT(*)
FROM analytics_events
WHERE event = 'click_cta' AND referrer_origin IS NOT NULL AND referrer_origin != ''
GROUP BY referrer_origin
ORDER BY COUNT(*) DESC
LIMIT ?;
`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []ReferrerMetric
	for rows.Next() {
		var metric ReferrerMetric
		if err := rows.Scan(&metric.Origin, &metric.Count); err != nil {
			return nil, err
		}
		items = append(items, metric)
	}
	return items, rows.Err()
}

func (r *Repository) ScrollDistribution(ctx context.Context) ([]ScrollMilestone, error) {
	if r == nil || r.db == nil {
		return nil, nil
	}

	rows, err := r.db.QueryContext(ctx, `
SELECT percent, COUNT(*)
FROM analytics_events
WHERE event = 'scroll' AND percent IS NOT NULL
GROUP BY percent
ORDER BY percent;
`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []ScrollMilestone
	for rows.Next() {
		var metric ScrollMilestone
		if err := rows.Scan(&metric.Percent, &metric.Count); err != nil {
			return nil, err
		}
		items = append(items, metric)
	}
	return items, rows.Err()
}

func (r *Repository) CTAClicksByPage(ctx context.Context, limit int) ([]CTAClickMetric, error) {
	if r == nil || r.db == nil {
		return nil, nil
	}
	if limit <= 0 {
		limit = 20
	}

	rows, err := r.db.QueryContext(ctx, `
SELECT page, cta, COUNT(*)
FROM analytics_events
WHERE event = 'click_cta' AND cta IS NOT NULL AND cta != ''
GROUP BY page, cta
ORDER BY COUNT(*) DESC
LIMIT ?;
`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []CTAClickMetric
	for rows.Next() {
		var metric CTAClickMetric
		if err := rows.Scan(&metric.Page, &metric.CTA, &metric.Count); err != nil {
			return nil, err
		}
		items = append(items, metric)
	}
	return items, rows.Err()
}

func nullableString(value string) interface{} {
	if value == "" {
		return nil
	}
	return value
}

func nullableInt(value int) interface{} {
	if value == 0 {
		return nil
	}
	return value
}
