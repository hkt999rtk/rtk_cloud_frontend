package analytics

import (
	"os"
	"strconv"
	"strings"
)

const (
	DefaultDatabasePath  = "data/analytics.db"
	DefaultRetentionDays = 90
)

type Config struct {
	Enabled       bool
	DatabasePath  string
	RetentionDays int
}

func DefaultConfig() Config {
	return Config{
		Enabled:       true,
		DatabasePath:  DefaultDatabasePath,
		RetentionDays: DefaultRetentionDays,
	}
}

func ConfigFromEnv() Config {
	cfg := DefaultConfig()

	if value, ok := os.LookupEnv("ANALYTICS_ENABLED"); ok {
		cfg.Enabled = parseEnabled(value)
	}

	if value := strings.TrimSpace(os.Getenv("ANALYTICS_DATABASE_PATH")); value != "" {
		cfg.DatabasePath = value
	}

	if value := strings.TrimSpace(os.Getenv("ANALYTICS_RETENTION_DAYS")); value != "" {
		if days, err := strconv.Atoi(value); err == nil && days > 0 {
			cfg.RetentionDays = days
		}
	}

	return cfg
}

func parseEnabled(value string) bool {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "0", "false", "no", "off":
		return false
	default:
		return true
	}
}
