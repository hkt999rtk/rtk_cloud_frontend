package web

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sort"
	"strings"
	"sync"
	"testing"
	"time"

	"realtek-connect/internal/analytics"
)

func TestHandleAnalyticsEventDisabledReturnsNotFound(t *testing.T) {
	handler := testServerWithConfig(t, Config{
		TemplatesDir: "../../templates",
		StaticDir:    "../../static",
		ContentDir:   "../../content/docs",
		LeadStore:    &memoryLeadStore{},
	})

	req := httptest.NewRequest(http.MethodPost, "/api/event", strings.NewReader(`{"event":"page_view","page":"home","session_id":"sid_1"}`))
	req.Header.Set("Content-Type", "application/json")
	req.RemoteAddr = "198.51.100.10:1234"

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusNotFound)
	}
}

func TestHandleAnalyticsEventSuccess(t *testing.T) {
	store := &memoryAnalyticsStore{}
	handler := testServerWithConfig(t, Config{
		TemplatesDir:   "../../templates",
		StaticDir:      "../../static",
		ContentDir:     "../../content/docs",
		LeadStore:      &memoryLeadStore{},
		AnalyticsStore: store,
	})

	req := httptest.NewRequest(http.MethodPost, "/api/event", strings.NewReader(`{"event":"page_view","page":"home","session_id":"sid_01"}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Origin", "https://example.com")
	req.Host = "example.com"
	req.RemoteAddr = "198.51.100.10:1234"

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusAccepted {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusAccepted)
	}
	if got := rec.Header().Get("Cache-Control"); got != "no-store" {
		t.Fatalf("cache-control = %q, want no-store", got)
	}
	if got := strings.TrimSpace(rec.Body.String()); got != `{"status":"ok"}` {
		t.Fatalf("body = %q, want %q", got, `{"status":"ok"}`)
	}
	if len(store.events) != 1 {
		t.Fatalf("events stored = %d, want 1", len(store.events))
	}
	if store.events[0].Page != "home" {
		t.Fatalf("stored page = %q, want home", store.events[0].Page)
	}
}

func TestHandleAnalyticsEventRejectsDisallowedRequests(t *testing.T) {
	store := &memoryAnalyticsStore{}
	handler := testServerWithConfig(t, Config{
		TemplatesDir:   "../../templates",
		StaticDir:      "../../static",
		ContentDir:     "../../content/docs",
		LeadStore:      &memoryLeadStore{},
		AnalyticsStore: store,
	})

	req := httptest.NewRequest(http.MethodGet, "/api/event", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("method status = %d, want %d", rec.Code, http.StatusMethodNotAllowed)
	}

	testCases := []struct {
		name  string
		body  string
		ctype string
	}{
		{
			name:  "invalid content type",
			body:  `{"event":"page_view","page":"home","session_id":"sid_01"}`,
			ctype: "text/plain",
		},
		{
			name:  "unknown event",
			body:  `{"event":"bad","page":"home","session_id":"sid_01"}`,
			ctype: "application/json",
		},
		{
			name:  "invalid page",
			body:  `{"event":"page_view","page":"../home","session_id":"sid_01"}`,
			ctype: "application/json",
		},
		{
			name:  "invalid cta",
			body:  `{"event":"click_cta","page":"home","cta":"bad_key","session_id":"sid_01"}`,
			ctype: "application/json",
		},
		{
			name:  "scroll missing percent",
			body:  `{"event":"scroll","page":"home","session_id":"sid_01"}`,
			ctype: "application/json",
		},
		{
			name:  "engaged invalid duration",
			body:  `{"event":"engaged","page":"home","duration":15,"session_id":"sid_01"}`,
			ctype: "application/json",
		},
		{
			name:  "bad session",
			body:  `{"event":"page_view","page":"home","session_id":""}`,
			ctype: "application/json",
		},
		{
			name:  "unknown field",
			body:  `{"event":"page_view","page":"home","session_id":"sid_01","unexpected":"x"}`,
			ctype: "application/json",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/api/event", strings.NewReader(tc.body))
			req.Header.Set("Content-Type", tc.ctype)
			req.Header.Set("Origin", "https://example.com")
			req.Host = "example.com"
			req.RemoteAddr = "198.51.100.10:1234"
			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)
			if rec.Code != http.StatusBadRequest {
				t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
			}
		})
	}

	if len(store.events) != 1 {
		t.Fatalf("events stored = %d, want 1", len(store.events))
	}
}

func TestHandleAnalyticsEventRejectsCrossOrigin(t *testing.T) {
	store := &memoryAnalyticsStore{}
	handler := testServerWithConfig(t, Config{
		TemplatesDir:   "../../templates",
		StaticDir:      "../../static",
		ContentDir:     "../../content/docs",
		LeadStore:      &memoryLeadStore{},
		AnalyticsStore: store,
	})

	req := httptest.NewRequest(http.MethodPost, "/api/event", strings.NewReader(`{"event":"page_view","page":"home","session_id":"sid_01"}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Origin", "https://evil.example")
	req.Header.Set("Referer", "https://example.com")
	req.Host = "example.com"
	req.RemoteAddr = "198.51.100.10:1234"

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusForbidden)
	}
	if len(store.events) != 0 {
		t.Fatalf("events stored = %d, want 0", len(store.events))
	}
}

func TestHandleAnalyticsEventRateLimitsAndSanitizesReferrer(t *testing.T) {
	store := &memoryAnalyticsStore{}
	server := newTestServer(t, Config{
		TemplatesDir:   "../../templates",
		StaticDir:      "../../static",
		ContentDir:     "../../content/docs",
		LeadStore:      &memoryLeadStore{},
		AnalyticsStore: store,
	})
	server.analyticsLimiter = &submissionRateLimiter{
		limit:  1,
		window: time.Minute,
		now:    func() time.Time { return time.Date(2026, 5, 6, 0, 0, 0, 0, time.UTC) },
		hits:   make(map[string][]time.Time),
	}
	handler := server.Routes()

	req := httptest.NewRequest(http.MethodPost, "/api/event", strings.NewReader(`{"event":"click_cta","page":"home","cta":"home_cta_primary","session_id":"sid_01"}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Origin", "https://example.com")
	req.Header.Set("Referer", "https://example.com/docs/overview?x=1#foo")
	req.Host = "example.com"
	req.RemoteAddr = "203.0.113.10:1234"
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusAccepted {
		t.Fatalf("first status = %d, want %d", rec.Code, http.StatusAccepted)
	}

	req = httptest.NewRequest(http.MethodPost, "/api/event", strings.NewReader(`{"event":"click_cta","page":"home","cta":"home_cta_primary","session_id":"sid_01"}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Origin", "https://example.com")
	req.Header.Set("Referer", "https://example.com/docs/overview?x=1#foo")
	req.Host = "example.com"
	req.RemoteAddr = "203.0.113.10:1234"
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusTooManyRequests {
		t.Fatalf("second status = %d, want %d", rec.Code, http.StatusTooManyRequests)
	}

	if len(store.events) != 1 {
		t.Fatalf("events stored = %d, want 1", len(store.events))
	}
	if store.events[0].ReferrerOrigin != "https://example.com" {
		t.Fatalf("referrer origin = %q, want %q", store.events[0].ReferrerOrigin, "https://example.com")
	}
}

func TestAnalyticsScriptRenderedBasedOnAnalyticsStore(t *testing.T) {
	t.Run("enabled", func(t *testing.T) {
		handler := testServerWithConfig(t, Config{
			TemplatesDir:   "../../templates",
			StaticDir:      "../../static",
			ContentDir:     "../../content/docs",
			LeadStore:      &memoryLeadStore{},
			AnalyticsStore: &memoryAnalyticsStore{},
		})

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
		}
		body := rec.Body.String()
		for _, want := range []string{"navigator.sendBeacon", "analyticsConfig", "trackScroll()", "sendEvent({", "analyticsConfig.page"} {
			if !strings.Contains(body, want) {
				t.Fatalf("expected body to contain %q", want)
			}
		}
		for _, banned := range []string{"localstorage", "sessionstorage", "document.cookie", "ga(", "gtag(", "google-analytics"} {
			if strings.Contains(strings.ToLower(body), banned) {
				t.Fatalf("body contains banned string %q", banned)
			}
		}
	})

	t.Run("disabled", func(t *testing.T) {
		handler := testServerWithConfig(t, Config{
			TemplatesDir: "../../templates",
			StaticDir:    "../../static",
			ContentDir:   "../../content/docs",
			LeadStore:    &memoryLeadStore{},
		})

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
		}
		if strings.Contains(rec.Body.String(), "analyticsConfig") {
			t.Fatalf("disabled analytics should not include analyticsConfig")
		}
	})
}

func TestAdminAnalyticsRouteRequiresAuthAndRendersAggregateView(t *testing.T) {
	store := &memoryAnalyticsStore{}
	normalizeErr := func(err error) {
		if err != nil {
			t.Fatalf("unexpected analytics store error: %v", err)
		}
	}
	normalizeErr(store.InsertEvent(context.Background(), analytics.Event{Event: "page_view", Page: "home", SessionID: "sid-1"}))
	normalizeErr(store.InsertEvent(context.Background(), analytics.Event{Event: "page_view", Page: "home", SessionID: "sid-2"}))
	normalizeErr(store.InsertEvent(context.Background(), analytics.Event{Event: "page_view", Page: "docs", SessionID: "sid-3"}))
	normalizeErr(store.InsertEvent(context.Background(), analytics.Event{Event: "click_cta", Page: "home", SessionID: "sid-1", CTA: "home_cta_primary", ReferrerOrigin: "https://example.com"}))
	normalizeErr(store.InsertEvent(context.Background(), analytics.Event{Event: "scroll", Page: "home", Percent: 25, SessionID: "sid-1"}))
	normalizeErr(store.InsertEvent(context.Background(), analytics.Event{Event: "scroll", Page: "home", Percent: 50, SessionID: "sid-2"}))

	handler := testServerWithConfig(t, Config{
		TemplatesDir:   "../../templates",
		StaticDir:      "../../static",
		ContentDir:     "../../content/docs",
		LeadStore:      &memoryLeadStore{},
		AnalyticsStore: store,
		AdminToken:     "secret",
	})

	req := httptest.NewRequest(http.MethodGet, "/admin/analytics", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status without token = %d, want %d", rec.Code, http.StatusUnauthorized)
	}

	req = httptest.NewRequest(http.MethodGet, "/admin/analytics?token=secret", nil)
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status with token = %d, want %d", rec.Code, http.StatusOK)
	}
	body := rec.Body.String()
	if !strings.Contains(body, "Conversion rate") {
		t.Fatalf("missing conversion section")
	}
	if !strings.Contains(body, "0.3333") {
		t.Fatalf("missing conversion value")
	}
	if !strings.Contains(body, "https://example.com") {
		t.Fatalf("missing referrer origin")
	}
	if !strings.Contains(body, "Top referrer origins") {
		t.Fatalf("missing referrer section")
	}
	if !strings.Contains(body, "CTA clicks by page") {
		t.Fatalf("missing CTA section")
	}
	if !strings.Contains(body, "Scroll distribution") {
		t.Fatalf("missing scroll section")
	}
}

func TestContactSubmitDoesNotWriteAnalytics(t *testing.T) {
	leadStore := &memoryLeadStore{}
	analyticsStore := &memoryAnalyticsStore{}
	handler := testServerWithConfig(t, Config{
		TemplatesDir:   "../../templates",
		StaticDir:      "../../static",
		ContentDir:     "../../content/docs",
		LeadStore:      leadStore,
		AnalyticsStore: analyticsStore,
	})

	req := httptest.NewRequest(http.MethodPost, "/contact", strings.NewReader("name=Kevin+Huang&email=kevin%40example.com&interest=evaluation-access&message=hello"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.RemoteAddr = "198.51.100.10:1234"
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if len(leadStore.leads) != 1 {
		t.Fatalf("leads stored = %d, want 1", len(leadStore.leads))
	}
	if len(analyticsStore.events) != 0 {
		t.Fatalf("analytics events = %d, want 0", len(analyticsStore.events))
	}
}

type memoryAnalyticsStore struct {
	mu     sync.Mutex
	events []analytics.Event
	err    error
}

func (s *memoryAnalyticsStore) InsertEvent(_ context.Context, event analytics.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.err != nil {
		return s.err
	}
	s.events = append(s.events, event)
	return nil
}

func (s *memoryAnalyticsStore) ConversionRate(_ context.Context) (float64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	var views, clicks float64
	for _, event := range s.events {
		switch event.Event {
		case "page_view":
			views++
		case "click_cta":
			clicks++
		}
	}
	if views == 0 {
		return 0, nil
	}
	return clicks / views, nil
}

func (s *memoryAnalyticsStore) TopReferrerOrigins(_ context.Context, limit int) ([]analytics.ReferrerMetric, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if limit <= 0 {
		limit = 10
	}

	counts := map[string]int64{}
	for _, event := range s.events {
		if event.Event != "click_cta" || event.ReferrerOrigin == "" {
			continue
		}
		counts[event.ReferrerOrigin]++
	}
	items := make([]analytics.ReferrerMetric, 0, len(counts))
	for origin, count := range counts {
		items = append(items, analytics.ReferrerMetric{Origin: origin, Count: count})
	}
	sort.Slice(items, func(i, j int) bool {
		if items[i].Count == items[j].Count {
			return items[i].Origin < items[j].Origin
		}
		return items[i].Count > items[j].Count
	})
	if len(items) > limit {
		items = items[:limit]
	}
	return items, nil
}

func (s *memoryAnalyticsStore) ScrollDistribution(_ context.Context) ([]analytics.ScrollMilestone, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	counts := map[int]int64{}
	for _, event := range s.events {
		if event.Event != "scroll" {
			continue
		}
		counts[event.Percent]++
	}
	items := make([]analytics.ScrollMilestone, 0, len(counts))
	for percent, count := range counts {
		items = append(items, analytics.ScrollMilestone{Percent: percent, Count: count})
	}
	sort.Slice(items, func(i, j int) bool {
		if items[i].Percent == items[j].Percent {
			return items[i].Count > items[j].Count
		}
		return items[i].Percent < items[j].Percent
	})
	return items, nil
}

func (s *memoryAnalyticsStore) CTAClicksByPage(_ context.Context, limit int) ([]analytics.CTAClickMetric, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if limit <= 0 {
		limit = 20
	}

	type key struct {
		page string
		cta  string
	}
	counts := map[key]int64{}
	for _, event := range s.events {
		if event.Event != "click_cta" {
			continue
		}
		counts[key{page: event.Page, cta: event.CTA}]++
	}
	items := make([]analytics.CTAClickMetric, 0, len(counts))
	for group, count := range counts {
		items = append(items, analytics.CTAClickMetric{Page: group.page, CTA: group.cta, Count: count})
	}
	sort.Slice(items, func(i, j int) bool {
		if items[i].Count == items[j].Count {
			if items[i].Page == items[j].Page {
				return items[i].CTA < items[j].CTA
			}
			return items[i].Page < items[j].Page
		}
		return items[i].Count > items[j].Count
	})
	if len(items) > limit {
		items = items[:limit]
	}
	return items, nil
}
