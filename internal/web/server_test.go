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

func testServer(t *testing.T, store LeadStore) http.Handler {
	t.Helper()
	server, err := NewServer(Config{
		TemplatesDir: "../../templates",
		StaticDir:    "../../static",
		LeadStore:    store,
	})
	if err != nil {
		t.Fatalf("new server: %v", err)
	}
	return server.Routes()
}

func TestRoutesReturnOK(t *testing.T) {
	handler := testServer(t, &memoryLeadStore{})
	paths := []string{"/", "/features", "/contact"}
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
