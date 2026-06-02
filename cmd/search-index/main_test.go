package main

import (
	"errors"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

func TestLogSearchIndexErrorUsesTypedFields(t *testing.T) {
	core, recorded := observer.New(zapcore.DebugLevel)
	logger := newSearchIndexTestLogger(core, "ci-abc123")

	logSearchIndexError(logger, "collect_content", "content_load_failed", errors.New("missing docs frontmatter"))

	entries := recorded.All()
	if len(entries) != 1 {
		t.Fatalf("entries = %d, want 1", len(entries))
	}
	entry := entries[0]
	if entry.Message != "search index failed" {
		t.Fatalf("message = %q, want search index failed", entry.Message)
	}

	fields := entry.ContextMap()
	for key, want := range map[string]any{
		"service":        "realtek-connect",
		"component":      "search-index",
		"version":        "ci-abc123",
		"operation":      "collect_content",
		"error_category": "content_load_failed",
	} {
		if fields[key] != want {
			t.Fatalf("%s = %#v, want %#v in %#v", key, fields[key], want, fields)
		}
	}
	if fields["error"] == "" {
		t.Fatalf("error field missing in %#v", fields)
	}
}

func newSearchIndexTestLogger(core zapcore.Core, version string) *zap.Logger {
	return zap.New(core).With(
		zap.String("service", "realtek-connect"),
		zap.String("component", "search-index"),
		zap.String("version", version),
	)
}
