package analytics

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"testing"
	"time"

	"realtek-connect/internal/leads"
)

func TestOpenInitializesSQLiteSchema(t *testing.T) {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "analytics.db")

	repository, err := Open(context.Background(), func() Config {
		cfg := DefaultConfig()
		cfg.DatabasePath = dbPath
		return cfg
	}())
	if err != nil {
		t.Fatalf("open analytics store: %v", err)
	}
	defer repository.Close()

	if _, err := os.Stat(dbPath); err != nil {
		t.Fatalf("analytics database file was not created: %v", err)
	}

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	defer db.Close()

	for _, name := range []string{
		"analytics_events",
		"idx_analytics_events_ts",
		"idx_analytics_events_event_page",
	} {
		var count int
		if err := db.QueryRowContext(context.Background(), `
SELECT COUNT(*)
FROM sqlite_master
WHERE name = ?`, name).Scan(&count); err != nil {
			t.Fatalf("lookup %s: %v", name, err)
		}
		if count != 1 {
			t.Fatalf("sqlite object %q not initialized", name)
		}
	}

	if got := repository.RetentionDays(); got != DefaultRetentionDays {
		t.Fatalf("retention days = %d, want %d", got, DefaultRetentionDays)
	}
}

func TestOpenDisabledSkipsStorageInitialization(t *testing.T) {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "analytics.db")

	repository, err := Open(context.Background(), Config{
		Enabled:      false,
		DatabasePath: dbPath,
	})
	if err != nil {
		t.Fatalf("open disabled analytics store: %v", err)
	}
	if repository != nil {
		t.Fatal("repository = non-nil, want nil when analytics is disabled")
	}

	if _, err := os.Stat(dbPath); !os.IsNotExist(err) {
		t.Fatalf("analytics database file exists when disabled: %v", err)
	}
}

func TestAnalyticsStorageIsSeparateFromLeadStorage(t *testing.T) {
	dir := t.TempDir()
	leadPath := filepath.Join(dir, "connectplus.db")
	analyticsPath := filepath.Join(dir, "analytics.db")

	leadDB, err := sql.Open("sqlite", leadPath)
	if err != nil {
		t.Fatalf("open lead sqlite: %v", err)
	}
	defer leadDB.Close()

	leadRepository := leads.NewRepository(leadDB)
	if err := leadRepository.Init(); err != nil {
		t.Fatalf("initialize lead schema: %v", err)
	}
	if err := leadRepository.Insert(context.Background(), leads.Lead{
		Name:     "Ada",
		Company:  "Example",
		Email:    "ada@example.com",
		Interest: "evaluation-access",
		Message:  "Please follow up.",
	}); err != nil {
		t.Fatalf("insert lead: %v", err)
	}

	analyticsRepository, err := Open(context.Background(), Config{
		Enabled:      true,
		DatabasePath: analyticsPath,
	})
	if err != nil {
		t.Fatalf("open analytics store: %v", err)
	}
	defer analyticsRepository.Close()

	assertSQLiteObjectExists(t, leadDB, "leads")
	assertSQLiteObjectMissing(t, leadDB, "analytics_events")

	analyticsDB, err := sql.Open("sqlite", analyticsPath)
	if err != nil {
		t.Fatalf("open analytics sqlite: %v", err)
	}
	defer analyticsDB.Close()

	assertSQLiteObjectExists(t, analyticsDB, "analytics_events")
	assertSQLiteObjectMissing(t, analyticsDB, "leads")
}

func TestCleanupExpiredEventsUsesDefaultRetention(t *testing.T) {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "analytics.db")

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	defer db.Close()

	repository := NewRepository(db, 0)
	if err := repository.Init(); err != nil {
		t.Fatalf("initialize analytics schema: %v", err)
	}

	asOf := time.Date(2026, 5, 5, 12, 0, 0, 0, time.UTC)
	cutoff := asOf.Add(-DefaultRetentionDays * 24 * time.Hour).Unix()

	seedAnalyticsEvent(t, db, cutoff-1)
	seedAnalyticsEvent(t, db, cutoff)
	seedAnalyticsEvent(t, db, cutoff+1)

	deleted, err := repository.CleanupExpiredEvents(context.Background(), asOf)
	if err != nil {
		t.Fatalf("cleanup expired events: %v", err)
	}
	if deleted != 1 {
		t.Fatalf("deleted rows = %d, want 1", deleted)
	}

	assertAnalyticsEventTimestamps(t, db, cutoff, cutoff+1)
}

func TestCleanupExpiredEventsUsesCustomRetention(t *testing.T) {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "analytics.db")

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	defer db.Close()

	repository := NewRepository(db, 14)
	if err := repository.Init(); err != nil {
		t.Fatalf("initialize analytics schema: %v", err)
	}

	asOf := time.Date(2026, 5, 5, 12, 0, 0, 0, time.UTC)
	cutoff := asOf.Add(-14 * 24 * time.Hour).Unix()

	seedAnalyticsEvent(t, db, cutoff-86400)
	seedAnalyticsEvent(t, db, cutoff)
	seedAnalyticsEvent(t, db, cutoff+86400)

	deleted, err := repository.CleanupExpiredEvents(context.Background(), asOf)
	if err != nil {
		t.Fatalf("cleanup expired events: %v", err)
	}
	if deleted != 1 {
		t.Fatalf("deleted rows = %d, want 1", deleted)
	}

	assertAnalyticsEventTimestamps(t, db, cutoff, cutoff+86400)
}

func TestCleanupExpiredEventsPreservesBoundaryTimestamp(t *testing.T) {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "analytics.db")

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	defer db.Close()

	repository := NewRepository(db, DefaultRetentionDays)
	if err := repository.Init(); err != nil {
		t.Fatalf("initialize analytics schema: %v", err)
	}

	asOf := time.Date(2026, 5, 5, 12, 0, 0, 0, time.UTC)
	boundary := asOf.Add(-DefaultRetentionDays * 24 * time.Hour).Unix()

	seedAnalyticsEvent(t, db, boundary)

	deleted, err := repository.CleanupExpiredEvents(context.Background(), asOf)
	if err != nil {
		t.Fatalf("cleanup expired events: %v", err)
	}
	if deleted != 0 {
		t.Fatalf("deleted rows = %d, want 0", deleted)
	}

	assertAnalyticsEventTimestamps(t, db, boundary)
}

func TestOpenCleansExpiredEventsOnStartup(t *testing.T) {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "analytics.db")

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}

	repository := NewRepository(db, DefaultRetentionDays)
	if err := repository.Init(); err != nil {
		_ = db.Close()
		t.Fatalf("initialize analytics schema: %v", err)
	}

	now := time.Now().UTC()
	seedAnalyticsEvent(t, db, now.Add(-91*24*time.Hour).Unix())
	seedAnalyticsEvent(t, db, now.Add(-89*24*time.Hour).Unix())

	if err := db.Close(); err != nil {
		t.Fatalf("close seeded sqlite: %v", err)
	}

	reopened, err := Open(context.Background(), Config{
		Enabled:       true,
		DatabasePath:  dbPath,
		RetentionDays: DefaultRetentionDays,
	})
	if err != nil {
		t.Fatalf("reopen analytics store: %v", err)
	}
	defer reopened.Close()

	verifyDB, err := sql.Open("sqlite", dbPath)
	if err != nil {
		t.Fatalf("open sqlite for verification: %v", err)
	}
	defer verifyDB.Close()

	if got := countRows(t, verifyDB); got != 1 {
		t.Fatalf("analytics rows = %d, want 1", got)
	}
}

func assertSQLiteObjectExists(t *testing.T, db *sql.DB, name string) {
	t.Helper()
	if countSQLiteObjects(t, db, name) != 1 {
		t.Fatalf("sqlite object %q not initialized", name)
	}
}

func assertSQLiteObjectMissing(t *testing.T, db *sql.DB, name string) {
	t.Helper()
	if countSQLiteObjects(t, db, name) != 0 {
		t.Fatalf("sqlite object %q should not exist", name)
	}
}

func countSQLiteObjects(t *testing.T, db *sql.DB, name string) int {
	t.Helper()

	var count int
	if err := db.QueryRowContext(context.Background(), `
SELECT COUNT(*)
FROM sqlite_master
WHERE name = ?`, name).Scan(&count); err != nil {
		t.Fatalf("lookup %s: %v", name, err)
	}
	return count
}

func seedAnalyticsEvent(t *testing.T, db *sql.DB, ts int64) {
	t.Helper()

	if _, err := db.ExecContext(context.Background(), `
INSERT INTO analytics_events (
  ts,
  event,
  page,
  created_at
) VALUES (?, ?, ?, ?)`,
		ts,
		"page_view",
		"home",
		time.Unix(ts, 0).UTC().Format(time.RFC3339),
	); err != nil {
		t.Fatalf("seed analytics event: %v", err)
	}
}

func assertAnalyticsEventTimestamps(t *testing.T, db *sql.DB, want ...int64) {
	t.Helper()

	rows, err := db.QueryContext(context.Background(), `
SELECT ts
FROM analytics_events
ORDER BY ts ASC`)
	if err != nil {
		t.Fatalf("query analytics timestamps: %v", err)
	}
	defer rows.Close()

	var got []int64
	for rows.Next() {
		var ts int64
		if err := rows.Scan(&ts); err != nil {
			t.Fatalf("scan analytics timestamp: %v", err)
		}
		got = append(got, ts)
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("iterate analytics timestamps: %v", err)
	}

	if len(got) != len(want) {
		t.Fatalf("timestamps = %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("timestamps = %v, want %v", got, want)
		}
	}
}

func countRows(t *testing.T, db *sql.DB) int {
	t.Helper()

	var count int
	if err := db.QueryRowContext(context.Background(), `SELECT COUNT(*) FROM analytics_events`).Scan(&count); err != nil {
		t.Fatalf("count analytics rows: %v", err)
	}
	return count
}
