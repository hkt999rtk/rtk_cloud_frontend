package analytics

import "testing"

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if !cfg.Enabled {
		t.Fatal("enabled = false, want true")
	}
	if cfg.DatabasePath != DefaultDatabasePath {
		t.Fatalf("database path = %q, want %q", cfg.DatabasePath, DefaultDatabasePath)
	}
	if cfg.RetentionDays != DefaultRetentionDays {
		t.Fatalf("retention days = %d, want %d", cfg.RetentionDays, DefaultRetentionDays)
	}
}

func TestConfigFromEnvAppliesOverrides(t *testing.T) {
	t.Setenv("ANALYTICS_ENABLED", "off")
	t.Setenv("ANALYTICS_DATABASE_PATH", "tmp/analytics.db")
	t.Setenv("ANALYTICS_RETENTION_DAYS", "14")

	cfg := ConfigFromEnv()

	if cfg.Enabled {
		t.Fatal("enabled = true, want false")
	}
	if cfg.DatabasePath != "/var/lib/realtek-connect/tmp/analytics.db" {
		t.Fatalf("database path = %q, want /var/lib/realtek-connect/tmp/analytics.db", cfg.DatabasePath)
	}
	if cfg.RetentionDays != 14 {
		t.Fatalf("retention days = %d, want 14", cfg.RetentionDays)
	}
}
