package web

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestPortalDesignNavigationAndManualMarkup(t *testing.T) {
	handler := testServer(t, &memoryLeadStore{})
	for _, prefix := range []string{"", "/zh-tw", "/zh-cn"} {
		for _, path := range []string{"/docs", "/features/provision", "/manual/getting-started", "/contact"} {
			t.Run(prefix+path, func(t *testing.T) {
				rec := httptest.NewRecorder()
				handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, prefix+path, nil))
				body := rec.Body.String()
				if rec.Code != http.StatusOK {
					t.Fatalf("status = %d", rec.Code)
				}
				for _, want := range []string{`class="portal-ui"`, `/static/portal-ui.css`, `aria-current="page"`, `href="` + prefix + `/manual"`} {
					if !strings.Contains(body, want) {
						t.Errorf("missing %s", want)
					}
				}
				if path == "/docs" && !strings.Contains(body, `data-analytics-cta="docs_cta_manual"`) {
					t.Error("docs must offer the manual as its primary developer entry")
				}
				if path == "/manual/getting-started" {
					for _, want := range []string{`class="portal-breadcrumb"`, `data-manual-article`, `data-manual-toc hidden`} {
						if !strings.Contains(body, want) {
							t.Errorf("missing %s", want)
						}
					}
				}
			})
		}
	}
}

func TestPortalDesignDoesNotStylePrivateAdmin(t *testing.T) {
	handler := testServerWithAdminToken(t, &memoryLeadStore{}, "test-token")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/admin/leads?token=test-token", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d", rec.Code)
	}
	if strings.Contains(rec.Body.String(), "portal-ui") {
		t.Error("private admin must retain its existing styles")
	}
}

func TestPortalManualCTAAnalyticsIsAccepted(t *testing.T) {
	repo, _ := openAnalyticsTestStore(t)
	defer repo.Close()
	handler := testServerWithConfig(t, Config{LeadStore: &memoryLeadStore{}, AnalyticsStore: repo})
	req := httptest.NewRequest(http.MethodPost, "/api/event", strings.NewReader(`{"event":"click_cta","page":"docs","cta":"docs_cta_manual","session_id":"portal-design-test"}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusAccepted {
		t.Fatalf("manual CTA status = %d: %s", rec.Code, rec.Body.String())
	}
}
