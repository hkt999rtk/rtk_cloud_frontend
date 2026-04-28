package web

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"realtek-connect/internal/features"
	"realtek-connect/internal/leads"
)

type memoryLeadStore struct {
	leads []leads.Lead
}

func (s *memoryLeadStore) Insert(_ context.Context, lead leads.Lead) error {
	s.leads = append(s.leads, lead)
	return nil
}

func (s *memoryLeadStore) List(_ context.Context, _ int) ([]leads.LeadRecord, error) {
	records := make([]leads.LeadRecord, 0, len(s.leads))
	for index, lead := range s.leads {
		records = append(records, leads.LeadRecord{
			ID:       int64(index + 1),
			Name:     lead.Name,
			Company:  lead.Company,
			Email:    lead.Email,
			Interest: lead.Interest,
			Message:  lead.Message,
		})
	}
	return records, nil
}

func testServer(t *testing.T, store LeadStore) http.Handler {
	t.Helper()
	return testServerWithAdminToken(t, store, "")
}

func testServerWithAdminToken(t *testing.T, store LeadStore, adminToken string) http.Handler {
	t.Helper()
	server, err := NewServer(Config{
		TemplatesDir: "../../templates",
		StaticDir:    "../../static",
		LeadStore:    store,
		AdminToken:   adminToken,
	})
	if err != nil {
		t.Fatalf("new server: %v", err)
	}
	return server.Routes()
}

func TestRoutesReturnOK(t *testing.T) {
	handler := testServer(t, &memoryLeadStore{})
	paths := []string{"/", "/features", "/contact", "/healthz"}
	for _, feature := range features.All() {
		paths = append(paths, "/features/"+feature.Slug)
	}

	for _, path := range paths {
		req := httptest.NewRequest(http.MethodGet, path, nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Fatalf("%s returned %d, want 200", path, rec.Code)
		}
	}
}

func TestUnknownFeatureReturnsNotFound(t *testing.T) {
	handler := testServer(t, &memoryLeadStore{})

	req := httptest.NewRequest(http.MethodGet, "/features/unknown", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", rec.Code)
	}
}

func TestValidContactPostStoresLead(t *testing.T) {
	store := &memoryLeadStore{}
	handler := testServer(t, store)

	form := url.Values{
		"name":     {"Kevin Huang"},
		"company":  {"Realtek"},
		"email":    {"kevin@example.com"},
		"interest": {"OTA"},
		"message":  {"Need scheduled rollout support."},
	}
	req := httptest.NewRequest(http.MethodPost, "/contact", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	if len(store.leads) != 1 {
		t.Fatalf("stored leads = %d, want 1", len(store.leads))
	}
	if !strings.Contains(rec.Body.String(), "Thanks, Kevin Huang") {
		t.Fatalf("response does not contain success message: %s", rec.Body.String())
	}
}

func TestInvalidContactPostDoesNotStoreLead(t *testing.T) {
	store := &memoryLeadStore{}
	handler := testServer(t, store)

	form := url.Values{
		"name":     {""},
		"company":  {"Realtek"},
		"email":    {"not-an-email"},
		"interest": {""},
	}
	req := httptest.NewRequest(http.MethodPost, "/contact", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", rec.Code)
	}
	if len(store.leads) != 0 {
		t.Fatalf("stored leads = %d, want 0", len(store.leads))
	}
	if !strings.Contains(rec.Body.String(), "Name is required") {
		t.Fatalf("response does not contain validation error: %s", rec.Body.String())
	}
}

func TestAdminLeadsRequiresToken(t *testing.T) {
	handler := testServerWithAdminToken(t, &memoryLeadStore{}, "secret")

	req := httptest.NewRequest(http.MethodGet, "/admin/leads", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want 401", rec.Code)
	}
}

func TestAdminLeadsDisabledWithoutToken(t *testing.T) {
	handler := testServer(t, &memoryLeadStore{})

	req := httptest.NewRequest(http.MethodGet, "/admin/leads", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", rec.Code)
	}
}

func TestAdminLeadsWithToken(t *testing.T) {
	store := &memoryLeadStore{leads: []leads.Lead{{
		Name:     "Kevin Huang",
		Company:  "Realtek",
		Email:    "kevin@example.com",
		Interest: "Provision",
		Message:  "Evaluate onboarding.",
	}}}
	handler := testServerWithAdminToken(t, store, "secret")

	req := httptest.NewRequest(http.MethodGet, "/admin/leads?token=secret", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "kevin@example.com") {
		t.Fatalf("response does not contain lead email: %s", rec.Body.String())
	}
}

func TestAdminLeadsCSVWithHeaderToken(t *testing.T) {
	store := &memoryLeadStore{leads: []leads.Lead{{
		Name:     "Kevin Huang",
		Company:  "Realtek",
		Email:    "kevin@example.com",
		Interest: "OTA",
		Message:  "CSV export.",
	}}}
	handler := testServerWithAdminToken(t, store, "secret")

	req := httptest.NewRequest(http.MethodGet, "/admin/leads.csv", nil)
	req.Header.Set("X-Admin-Token", "secret")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	if got := rec.Header().Get("Content-Type"); !strings.Contains(got, "text/csv") {
		t.Fatalf("content type = %q, want text/csv", got)
	}
	if !strings.Contains(rec.Body.String(), "kevin@example.com") {
		t.Fatalf("csv does not contain lead email: %s", rec.Body.String())
	}
}
