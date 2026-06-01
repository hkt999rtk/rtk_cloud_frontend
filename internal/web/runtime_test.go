package web

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
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
	logger := jsonTestLogger(t, &output)

	handler := LoggingMiddleware(logger)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/healthz" {
			t.Fatalf("path = %s, want /healthz", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	}))

	req := httptest.NewRequest(http.MethodGet, "/healthz?check=1", nil)
	req.RemoteAddr = "203.0.113.10:54321"
	req.Header.Set("X-Request-Id", "req-123")
	req.Header.Set("X-Trace-Id", "trace-456")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want 204", rec.Code)
	}

	event := decodeRuntimeLogEvent(t, output.Bytes())
	for key, want := range map[string]any{
		"msg":         "http request",
		"method":      http.MethodGet,
		"path":        "/healthz?check=1",
		"status":      float64(http.StatusNoContent),
		"remote_addr": "203.0.113.10",
		"request_id":  "req-123",
		"trace_id":    "trace-456",
	} {
		if event[key] != want {
			t.Fatalf("%s = %#v, want %#v in %#v", key, event[key], want, event)
		}
	}
	if _, ok := event["duration_ms"].(float64); !ok {
		t.Fatalf("duration_ms = %#v, want numeric field in %#v", event["duration_ms"], event)
	}
}

func TestLoggingMiddlewareRedactsSensitiveQueryParameters(t *testing.T) {
	var output bytes.Buffer
	logger := jsonTestLogger(t, &output)

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
	event := decodeRuntimeLogEvent(t, output.Bytes())
	if event["path"] != "/admin/leads?token=[REDACTED]&view=full" {
		t.Fatalf("path = %#v, want redacted token and preserved request details in %#v", event["path"], event)
	}
}

func jsonTestLogger(t *testing.T, output *bytes.Buffer) *zap.Logger {
	t.Helper()

	encoderCfg := zap.NewProductionEncoderConfig()
	core := zapcore.NewCore(zapcore.NewJSONEncoder(encoderCfg), zapcore.AddSync(output), zapcore.DebugLevel)
	return zap.New(core)
}

func decodeRuntimeLogEvent(t *testing.T, line []byte) map[string]any {
	t.Helper()

	var event map[string]any
	if err := json.Unmarshal(bytes.TrimSpace(line), &event); err != nil {
		t.Fatalf("decode log event %q: %v", string(line), err)
	}
	return event
}
