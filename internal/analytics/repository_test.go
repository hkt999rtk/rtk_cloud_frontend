package analytics

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"testing"

	_ "modernc.org/sqlite"
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
