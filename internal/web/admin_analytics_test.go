package web

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"realtek-connect/internal/analytics"
)

func TestAdminAnalyticsPageRequiresToken(t *testing.T) {
	repo, _ := openAnalyticsTestStore(t)
	defer repo.Close()

	handler := testServerWithConfig(t, Config{
		LeadStore:      &memoryLeadStore{},
		AnalyticsStore: repo,
		AdminToken:     "secret",
	})

	req := httptest.NewRequest(http.MethodGet, "/admin/analytics", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusUnauthorized)
	}
}

func TestAdminAnalyticsPageRendersAggregateMetrics(t *testing.T) {
	repo, _ := openAnalyticsTestStore(t)
	defer repo.Close()

	seedAdminAnalyticsData(t, repo)

	handler := testServerWithConfig(t, Config{
		LeadStore:      &memoryLeadStore{},
		AnalyticsStore: repo,
		AdminToken:     "secret",
	})

	req := httptest.NewRequest(http.MethodGet, "/admin/analytics?token=secret", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200: %s", rec.Code, rec.Body.String())
	}

	body := rec.Body.String()
	for _, want := range []string{
		"Analytics summary",
		"Page views",
		"CTA clicks",
		"Conversion rate",
		"75.0%",
		"Top referrer origins",
		"https://example.com",
		"(direct)",
		"Scroll distribution",
		"25%",
		"100%",
		"CTA clicks by page",
		"home",
		"contact_us",
		"features",
		"talk_to_sales",
	} {
		if !strings.Contains(body, want) {
			t.Fatalf("analytics page missing %q: %s", want, body)
		}
	}

	if !strings.Contains(body, `href="/admin/leads?token=secret"`) {
		t.Fatalf("analytics page should link back to leads with token: %s", body)
	}
}

func seedAdminAnalyticsData(t *testing.T, repo *analytics.Repository) {
	t.Helper()

	now := time.Date(2026, 5, 5, 12, 0, 0, 0, time.UTC)
	for i := 0; i < 4; i++ {
		if err := repo.InsertEvent(context.Background(), analytics.Event{
			TS:        now.Add(time.Duration(i) * time.Minute).Unix(),
			Type:      "page_view",
			Page:      "home",
			SessionID: "session-page-" + strconv.Itoa(i),
			CreatedAt: now,
		}); err != nil {
			t.Fatalf("seed page_view: %v", err)
		}
	}

	if err := repo.InsertEvent(context.Background(), analytics.Event{
		TS:             now.Add(10 * time.Minute).Unix(),
		Type:           "click_cta",
		Page:           "home",
		CTA:            "contact_us",
		ReferrerOrigin: "https://example.com",
		SessionID:      "session-click-1",
		CreatedAt:      now,
	}); err != nil {
		t.Fatalf("seed click_cta 1: %v", err)
	}
	if err := repo.InsertEvent(context.Background(), analytics.Event{
		TS:             now.Add(11 * time.Minute).Unix(),
		Type:           "click_cta",
		Page:           "home",
		CTA:            "contact_us",
		ReferrerOrigin: "https://example.com",
		SessionID:      "session-click-2",
		CreatedAt:      now,
	}); err != nil {
		t.Fatalf("seed click_cta 2: %v", err)
	}
	if err := repo.InsertEvent(context.Background(), analytics.Event{
		TS:        now.Add(12 * time.Minute).Unix(),
		Type:      "click_cta",
		Page:      "features",
		CTA:       "talk_to_sales",
		SessionID: "session-click-3",
		CreatedAt: now,
	}); err != nil {
		t.Fatalf("seed click_cta 3: %v", err)
	}

	for _, percent := range []int{25, 100} {
		p := percent
		if err := repo.InsertEvent(context.Background(), analytics.Event{
			TS:        now.Add(20 * time.Minute).Unix(),
			Type:      "scroll",
			Page:      "home",
			Percent:   &p,
			SessionID: "session-scroll-" + strconv.Itoa(percent),
			CreatedAt: now,
		}); err != nil {
			t.Fatalf("seed scroll: %v", err)
		}
	}

	if err := repo.InsertEvent(context.Background(), analytics.Event{
		TS:        now.Add(30 * time.Minute).Unix(),
		Type:      "engaged",
		Page:      "home",
		Duration:  durationPtr(10),
		SessionID: "session-engaged",
		CreatedAt: now,
	}); err != nil {
		t.Fatalf("seed engaged: %v", err)
	}
}

func durationPtr(value int) *int {
	return &value
}
