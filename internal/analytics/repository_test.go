package analytics

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"testing"

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
