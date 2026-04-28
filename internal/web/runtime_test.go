package web

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestNewHTTPServerAppliesDefaultTimeouts(t *testing.T) {
	server := NewHTTPServer(HTTPServerConfig{
		Addr:    ":8080",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {}),
	})

	if server.ReadTimeout != DefaultReadTimeout {
		t.Fatalf("read timeout = %s, want %s", server.ReadTimeout, DefaultReadTimeout)
	}
	if server.WriteTimeout != DefaultWriteTimeout {
		t.Fatalf("write timeout = %s, want %s", server.WriteTimeout, DefaultWriteTimeout)
	}
	if server.IdleTimeout != DefaultIdleTimeout {
		t.Fatalf("idle timeout = %s, want %s", server.IdleTimeout, DefaultIdleTimeout)
	}
}

func TestNewHTTPServerUsesExplicitTimeouts(t *testing.T) {
	server := NewHTTPServer(HTTPServerConfig{
		Addr:         ":8080",
		Handler:      http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {}),
		ReadTimeout:  2 * time.Second,
		WriteTimeout: 3 * time.Second,
		IdleTimeout:  4 * time.Second,
	})

	if server.ReadTimeout != 2*time.Second {
		t.Fatalf("read timeout = %s, want 2s", server.ReadTimeout)
	}
	if server.WriteTimeout != 3*time.Second {
		t.Fatalf("write timeout = %s, want 3s", server.WriteTimeout)
	}
	if server.IdleTimeout != 4*time.Second {
		t.Fatalf("idle timeout = %s, want 4s", server.IdleTimeout)
	}
}

func TestLoggingMiddlewareLogsRequestOutcome(t *testing.T) {
	var output bytes.Buffer
	logger := log.New(&output, "", 0)

	handler := LoggingMiddleware(logger)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/healthz" {
			t.Fatalf("path = %s, want /healthz", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	}))

	req := httptest.NewRequest(http.MethodGet, "/healthz?check=1", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want 204", rec.Code)
	}

	logLine := output.String()
	if !strings.Contains(logLine, "GET /healthz?check=1 204") {
		t.Fatalf("log line = %q, want method, path, and status", logLine)
	}
}

func TestLoggingMiddlewareRedactsSensitiveQueryParameters(t *testing.T) {
	var output bytes.Buffer
	logger := log.New(&output, "", 0)

	handler := LoggingMiddleware(logger)(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/admin/leads?token=secret&view=full", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	logLine := output.String()
	if strings.Contains(logLine, "secret") {
		t.Fatalf("log line = %q, secret token leaked", logLine)
	}
	if !strings.Contains(logLine, "/admin/leads?token=REDACTED&view=full 200") {
		t.Fatalf("log line = %q, want redacted token and preserved request details", logLine)
	}
}
