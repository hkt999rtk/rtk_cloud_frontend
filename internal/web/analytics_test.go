package web

import (
	"context"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"net/url"
	"path/filepath"
	"strings"
	"testing"

	"realtek-connect/internal/analytics"
)

func TestAnalyticsEventEndpointStoresValidEvent(t *testing.T) {
	repo, dbPath := openAnalyticsTestStore(t)
	defer repo.Close()

	handler := testServerWithConfig(t, Config{
		LeadStore:      &memoryLeadStore{},
		AnalyticsStore: repo,
	})

	req := httptest.NewRequest(http.MethodPost, "/api/event", strings.NewReader(`{"event":"click_cta","page":"home","cta":"contact_us","variant":"control","session_id":"session-123","extra":"ignored"}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Referer", "https://example.com/path?q=1#frag")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusAccepted {
		t.Fatalf("status = %d, want %d: %s", rec.Code, http.StatusAccepted, rec.Body.String())
	}
	if got := rec.Header().Get("Cache-Control"); got != "no-store" {
		t.Fatalf("Cache-Control = %q, want no-store", got)
	}
	if got := strings.TrimSpace(rec.Body.String()); got != `{"status":"ok"}` {
		t.Fatalf("body = %q, want %q", got, `{"status":"ok"}`)
	}

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	defer db.Close()

	var (
		ts             int64
		eventType      string
		page           string
		cta            sql.NullString
		percent        sql.NullInt64
		duration       sql.NullInt64
		variant        sql.NullString
		referrerOrigin sql.NullString
		sessionID      string
		createdAt      string
	)
	if err := db.QueryRowContext(context.Background(), `
SELECT ts, event, page, cta, percent, duration, variant, referrer_origin, session_id, created_at
FROM analytics_events
LIMIT 1`).Scan(&ts, &eventType, &page, &cta, &percent, &duration, &variant, &referrerOrigin, &sessionID, &createdAt); err != nil {
		t.Fatalf("query event: %v", err)
	}

	if ts == 0 {
		t.Fatal("ts = 0, want generated timestamp")
	}
	if eventType != "click_cta" || page != "home" {
		t.Fatalf("event/page = %q/%q, want click_cta/home", eventType, page)
	}
	if !cta.Valid || cta.String != "contact_us" {
		t.Fatalf("cta = %+v, want contact_us", cta)
	}
	if percent.Valid || duration.Valid {
		t.Fatalf("unexpected optional fields stored: percent=%+v duration=%+v", percent, duration)
	}
	if !variant.Valid || variant.String != "control" {
		t.Fatalf("variant = %+v, want control", variant)
	}
	if !referrerOrigin.Valid || referrerOrigin.String != "https://example.com" {
		t.Fatalf("referrer origin = %+v, want https://example.com", referrerOrigin)
	}
	if sessionID != "session-123" {
		t.Fatalf("session id = %q, want session-123", sessionID)
	}
	if createdAt == "" {
		t.Fatal("created_at = empty, want server timestamp")
	}
}

func TestPublicPageRendersAnalyticsConfigWhenStoreExists(t *testing.T) {
	repo, _ := openAnalyticsTestStore(t)
	defer repo.Close()

	handler := testServerWithConfig(t, Config{
		LeadStore:      &memoryLeadStore{},
		AnalyticsStore: repo,
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	body := rec.Body.String()
	for _, want := range []string{
		`endpoint: "\/api\/event"`,
		`page: "home"`,
		`event: "page_view"`,
	} {
		if !strings.Contains(body, want) {
			t.Fatalf("analytics config missing %q: %s", want, body)
		}
	}
	if strings.Contains(body, "template render error") {
		t.Fatalf("page rendered template error: %s", body)
	}
}

func TestAnalyticsEventEndpointRejectsInvalidEvents(t *testing.T) {
	cases := []struct {
		name    string
		body    string
		want    int
		headers map[string]string
	}{
		{
			name: "invalid event",
			body: `{"event":"hover","page":"home","session_id":"session-123"}`,
			want: http.StatusBadRequest,
		},
		{
			name: "invalid page",
			body: `{"event":"page_view","page":"/home","session_id":"session-123"}`,
			want: http.StatusBadRequest,
		},
		{
			name: "invalid cta",
			body: `{"event":"click_cta","page":"home","cta":"bad","session_id":"session-123"}`,
			want: http.StatusBadRequest,
		},
		{
			name: "invalid percent",
			body: `{"event":"scroll","page":"home","percent":33,"session_id":"session-123"}`,
			want: http.StatusBadRequest,
		},
		{
			name: "invalid duration",
			body: `{"event":"engaged","page":"home","duration":5,"session_id":"session-123"}`,
			want: http.StatusBadRequest,
		},
		{
			name: "invalid session id",
			body: `{"event":"page_view","page":"home","session_id":"has spaces"}`,
			want: http.StatusBadRequest,
		},
		{
			name: "invalid variant",
			body: `{"event":"page_view","page":"home","variant":"unknown","session_id":"session-123"}`,
			want: http.StatusBadRequest,
		},
		{
			name: "invalid content type",
			body: `{"event":"page_view","page":"home","session_id":"session-123"}`,
			want: http.StatusUnsupportedMediaType,
			headers: map[string]string{
				"Content-Type": "text/plain",
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			repo, dbPath := openAnalyticsTestStore(t)
			defer repo.Close()

			handler := testServerWithConfig(t, Config{
				LeadStore:      &memoryLeadStore{},
				AnalyticsStore: repo,
			})

			req := httptest.NewRequest(http.MethodPost, "/api/event", strings.NewReader(tc.body))
			req.Header.Set("Content-Type", "application/json")
			if tc.headers != nil {
				for key, value := range tc.headers {
					req.Header.Set(key, value)
				}
			}
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			if rec.Code != tc.want {
				t.Fatalf("status = %d, want %d: %s", rec.Code, tc.want, rec.Body.String())
			}
			if rec.Header().Get("Cache-Control") != "no-store" {
				t.Fatalf("Cache-Control = %q, want no-store", rec.Header().Get("Cache-Control"))
			}
			if got := analyticsRowCount(t, dbPath); got != 0 {
				t.Fatalf("analytics rows = %d, want 0", got)
			}
		})
	}
}

func TestAnalyticsEventEndpointReturnsNoContentWhenDisabled(t *testing.T) {
	handler := testServerWithConfig(t, Config{
		LeadStore:      &memoryLeadStore{},
		AnalyticsStore: nil,
	})

	req := httptest.NewRequest(http.MethodPost, "/api/event", strings.NewReader(`{"event":"page_view","page":"home","session_id":"session-123"}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusNoContent)
	}
	if rec.Body.Len() != 0 {
		t.Fatalf("disabled response body = %q, want empty", rec.Body.String())
	}
}

func TestAnalyticsEventEndpointRejectsOversizedBodies(t *testing.T) {
	repo, dbPath := openAnalyticsTestStore(t)
	defer repo.Close()

	handler := testServerWithConfig(t, Config{
		LeadStore:      &memoryLeadStore{},
		AnalyticsStore: repo,
	})

	hugeSessionID := strings.Repeat("a", analyticsRequestBodyLimit)
	req := httptest.NewRequest(http.MethodPost, "/api/event", strings.NewReader(`{"event":"page_view","page":"home","session_id":"`+hugeSessionID+`"}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("status = %d, want %d: %s", rec.Code, http.StatusRequestEntityTooLarge, rec.Body.String())
	}
	if got := analyticsRowCount(t, dbPath); got != 0 {
		t.Fatalf("analytics rows = %d, want 0", got)
	}
}

func TestContactSubmissionDoesNotCaptureRawFormValuesInAnalytics(t *testing.T) {
	repo, dbPath := openAnalyticsTestStore(t)
	defer repo.Close()

	store := &memoryLeadStore{}
	handler := testServerWithConfig(t, Config{
		LeadStore:      store,
		AnalyticsStore: repo,
	})

	form := url.Values{
		"name":     {"Ada Lovelace"},
		"company":  {"Example Systems"},
		"email":    {"ada@example.com"},
		"interest": {"evaluation-access"},
		"message":  {"Please do not track this text in analytics."},
	}
	req := httptest.NewRequest(http.MethodPost, "/contact", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d: %s", rec.Code, http.StatusOK, rec.Body.String())
	}
	if got := len(store.leads); got != 1 {
		t.Fatalf("lead rows = %d, want 1", got)
	}
	if got := analyticsRowCount(t, dbPath); got != 0 {
		t.Fatalf("analytics rows = %d, want 0", got)
	}
}

func openAnalyticsTestStore(t *testing.T) (*analytics.Repository, string) {
	t.Helper()

	dir := t.TempDir()
	dbPath := filepath.Join(dir, "analytics.db")
	repo, err := analytics.Open(context.Background(), analytics.Config{
		Enabled:      true,
		DatabasePath: dbPath,
	})
	if err != nil {
		t.Fatalf("open analytics store: %v", err)
	}
	return repo, dbPath
}

func analyticsRowCount(t *testing.T, dbPath string) int {
	t.Helper()

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	defer db.Close()

	var count int
	if err := db.QueryRowContext(context.Background(), `SELECT COUNT(*) FROM analytics_events`).Scan(&count); err != nil {
		t.Fatalf("count analytics events: %v", err)
	}
	return count
}
