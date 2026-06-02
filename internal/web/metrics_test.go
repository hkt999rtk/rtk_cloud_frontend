package web

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"realtek-connect/internal/leads"
)

func TestPrometheusMetrics(t *testing.T) {
	store := &memoryLeadStore{leads: []leads.Lead{{
		Name:     "Ada",
		Company:  "Acme",
		Email:    "ada@example.com",
		Interest: "ota",
		Message:  "hello",
	}}}
	handler := testServer(t, store)

	req := httptest.NewRequest(http.MethodGet, "/metrics/prometheus", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("metrics status = %d, want %d; body=%s", rec.Code, http.StatusOK, rec.Body.String())
	}
	if got := rec.Header().Get("Content-Type"); !strings.Contains(got, "text/plain") {
		t.Fatalf("content type = %q, want text/plain", got)
	}
	for _, want := range []string{
		"rtk_cloud_frontend_up 1",
		"rtk_cloud_frontend_leads_total 1",
	} {
		if !strings.Contains(rec.Body.String(), want) {
			t.Fatalf("metrics body missing %q:\n%s", want, rec.Body.String())
		}
	}
}
