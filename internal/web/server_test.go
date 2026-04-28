package web

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"
	"time"

	"realtek-connect/internal/docs"
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

func (s *memoryLeadStore) Count(_ context.Context, filter leads.ListFilter) (int, error) {
	return len(s.filteredRecords(filter)), nil
}

func (s *memoryLeadStore) List(_ context.Context, opts leads.ListOptions) ([]leads.LeadRecord, error) {
	records := s.filteredRecords(opts.Filter)
	if opts.Offset >= len(records) {
		return []leads.LeadRecord{}, nil
	}
	if opts.Offset > 0 {
		records = records[opts.Offset:]
	}
	if opts.Limit > 0 && opts.Limit < len(records) {
		records = records[:opts.Limit]
	}
	return records, nil
}

func (s *memoryLeadStore) filteredRecords(filter leads.ListFilter) []leads.LeadRecord {
	records := make([]leads.LeadRecord, 0, len(s.leads))
	email := strings.ToLower(strings.TrimSpace(filter.Email))
	company := strings.ToLower(strings.TrimSpace(filter.Company))
	interest := strings.ToLower(strings.TrimSpace(filter.Interest))
	for index := len(s.leads) - 1; index >= 0; index-- {
		lead := s.leads[index]
		if email != "" && !strings.Contains(strings.ToLower(lead.Email), email) {
			continue
		}
		if company != "" && !strings.Contains(strings.ToLower(lead.Company), company) {
			continue
		}
		if interest != "" && !strings.Contains(strings.ToLower(lead.Interest), interest) {
			continue
		}
		records = append(records, leads.LeadRecord{
			ID:       int64(index + 1),
			Name:     lead.Name,
			Company:  lead.Company,
			Email:    lead.Email,
			Interest: lead.Interest,
			Message:  lead.Message,
		})
	}
	return records
}

func testServer(t *testing.T, store LeadStore) http.Handler {
	t.Helper()
	return testServerWithAdminToken(t, store, "")
}

func testServerWithAdminToken(t *testing.T, store LeadStore, adminToken string) http.Handler {
	t.Helper()
	server := newTestServer(t, Config{
		TemplatesDir: "../../templates",
		StaticDir:    "../../static",
		LeadStore:    store,
		AdminToken:   adminToken,
	})
	return server.Routes()
}

func newTestServer(t *testing.T, cfg Config) *Server {
	t.Helper()
	server, err := NewServer(cfg)
	if err != nil {
		t.Fatalf("new server: %v", err)
	}
	return server
}

func TestRoutesReturnOK(t *testing.T) {
	handler := testServer(t, &memoryLeadStore{})
	paths := []string{"/", "/docs", "/features", "/contact", "/healthz", "/robots.txt", "/sitemap.xml"}
	for _, section := range docs.All() {
		paths = append(paths, "/docs/"+section.Slug)
	}
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

func TestHomeMetadataIncludesSocialTags(t *testing.T) {
	handler := testServer(t, &memoryLeadStore{})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}

	body := rec.Body.String()
	for _, want := range []string{
		"<title>Realtek Connect&#43; | IoT Cloud Platform</title>",
		`<meta name="description" content="Realtek Connect&#43; is an IoT cloud platform for provisioning, OTA, fleet management, app SDKs, insights, private cloud, and integrations.">`,
		`<link rel="canonical" href="http://example.com/">`,
		`<meta property="og:title" content="Realtek Connect&#43; | IoT Cloud Platform">`,
		`<meta property="og:url" content="http://example.com/">`,
		`<meta property="og:image" content="http://example.com/static/assets/connectplus-hero.png">`,
		`<meta name="twitter:card" content="summary_large_image">`,
	} {
		if !strings.Contains(body, want) {
			t.Fatalf("response does not contain %q: %s", want, body)
		}
	}
}

func TestFeatureMetadataUsesFeatureSummary(t *testing.T) {
	handler := testServer(t, &memoryLeadStore{})

	req := httptest.NewRequest(http.MethodGet, "/features/ota", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}

	body := rec.Body.String()
	for _, want := range []string{
		`<title>OTA | Realtek Connect&#43;</title>`,
		`<meta name="description" content="Upload firmware, create rollout campaigns, monitor jobs, and protect devices with version validation.">`,
		`<meta property="og:url" content="http://example.com/features/ota">`,
		`<meta name="twitter:title" content="OTA | Realtek Connect&#43;">`,
	} {
		if !strings.Contains(body, want) {
			t.Fatalf("response does not contain %q: %s", want, body)
		}
	}
}

func TestRobotsTxtIncludesSitemap(t *testing.T) {
	handler := testServer(t, &memoryLeadStore{})

	req := httptest.NewRequest(http.MethodGet, "/robots.txt", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	if got := rec.Header().Get("Content-Type"); !strings.Contains(got, "text/plain") {
		t.Fatalf("content type = %q, want text/plain", got)
	}

	body := rec.Body.String()
	for _, want := range []string{
		"User-agent: *",
		"Disallow: /admin/",
		"Sitemap: http://example.com/sitemap.xml",
	} {
		if !strings.Contains(body, want) {
			t.Fatalf("robots.txt does not contain %q: %s", want, body)
		}
	}
}

func TestSitemapXMLIncludesPublicRoutes(t *testing.T) {
	handler := testServer(t, &memoryLeadStore{})

	req := httptest.NewRequest(http.MethodGet, "/sitemap.xml", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	if got := rec.Header().Get("Content-Type"); !strings.Contains(got, "application/xml") {
		t.Fatalf("content type = %q, want application/xml", got)
	}

	body := rec.Body.String()
	for _, want := range []string{
		`<?xml version="1.0" encoding="UTF-8"?>`,
		`<loc>http://example.com/</loc>`,
		`<loc>http://example.com/docs/product-overview</loc>`,
		`<loc>http://example.com/features/ota</loc>`,
		`<loc>http://example.com/contact</loc>`,
	} {
		if !strings.Contains(body, want) {
			t.Fatalf("sitemap does not contain %q: %s", want, body)
		}
	}
	if strings.Contains(body, "/admin/leads") {
		t.Fatalf("sitemap should not contain admin routes: %s", body)
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

func TestUnknownDocReturnsNotFound(t *testing.T) {
	handler := testServer(t, &memoryLeadStore{})

	req := httptest.NewRequest(http.MethodGet, "/docs/unknown", nil)
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

func TestOversizedContactPostDoesNotStoreLead(t *testing.T) {
	store := &memoryLeadStore{}
	handler := testServer(t, store)

	form := url.Values{
		"name":     {strings.Repeat("N", contactNameMaxLength+1)},
		"company":  {strings.Repeat("C", contactCompanyMaxLength+1)},
		"email":    {"kevin@example.com"},
		"interest": {"OTA"},
		"message":  {strings.Repeat("M", contactMessageMaxLength+1)},
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
	body := rec.Body.String()
	if !strings.Contains(body, "Name must be 120 characters or fewer.") {
		t.Fatalf("response does not contain name limit error: %s", body)
	}
	if !strings.Contains(body, "Company must be 160 characters or fewer.") {
		t.Fatalf("response does not contain company limit error: %s", body)
	}
	if !strings.Contains(body, "Message must be 2000 characters or fewer.") {
		t.Fatalf("response does not contain message limit error: %s", body)
	}
}

func TestSpamContactPostDoesNotStoreLead(t *testing.T) {
	store := &memoryLeadStore{}
	handler := testServer(t, store)

	form := url.Values{
		"name":     {"Kevin Huang"},
		"email":    {"kevin@example.com"},
		"interest": {"OTA"},
		"website":  {"https://spam.example.com"},
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
	if !strings.Contains(rec.Body.String(), "Request could not be processed.") {
		t.Fatalf("response does not contain form error: %s", rec.Body.String())
	}
}

func TestContactPostRateLimitsRepeatedSubmissions(t *testing.T) {
	store := &memoryLeadStore{}
	server := newTestServer(t, Config{
		TemplatesDir: "../../templates",
		StaticDir:    "../../static",
		LeadStore:    store,
	})

	now := time.Date(2026, 4, 28, 0, 0, 0, 0, time.UTC)
	server.contactLimit = &submissionRateLimiter{
		limit:  1,
		window: time.Minute,
		now:    func() time.Time { return now },
		hits:   make(map[string][]time.Time),
	}
	handler := server.Routes()

	form := url.Values{
		"name":     {"Kevin Huang"},
		"email":    {"kevin@example.com"},
		"interest": {"OTA"},
	}

	firstReq := httptest.NewRequest(http.MethodPost, "/contact", strings.NewReader(form.Encode()))
	firstReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	firstReq.RemoteAddr = "198.51.100.10:1234"
	firstRec := httptest.NewRecorder()
	handler.ServeHTTP(firstRec, firstReq)

	if firstRec.Code != http.StatusOK {
		t.Fatalf("first status = %d, want 200", firstRec.Code)
	}

	secondReq := httptest.NewRequest(http.MethodPost, "/contact", strings.NewReader(form.Encode()))
	secondReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	secondReq.RemoteAddr = "198.51.100.10:1234"
	secondRec := httptest.NewRecorder()
	handler.ServeHTTP(secondRec, secondReq)

	if secondRec.Code != http.StatusTooManyRequests {
		t.Fatalf("second status = %d, want 429", secondRec.Code)
	}
	if len(store.leads) != 1 {
		t.Fatalf("stored leads = %d, want 1", len(store.leads))
	}
	if !strings.Contains(secondRec.Body.String(), "Too many requests from this address.") {
		t.Fatalf("response does not contain rate limit error: %s", secondRec.Body.String())
	}
}

func TestContactPostRateLimitIgnoresSpoofedForwardingHeaders(t *testing.T) {
	store := &memoryLeadStore{}
	server := newTestServer(t, Config{
		TemplatesDir: "../../templates",
		StaticDir:    "../../static",
		LeadStore:    store,
	})

	now := time.Date(2026, 4, 28, 0, 0, 0, 0, time.UTC)
	server.contactLimit = &submissionRateLimiter{
		limit:  1,
		window: time.Minute,
		now:    func() time.Time { return now },
		hits:   make(map[string][]time.Time),
	}
	handler := server.Routes()

	form := url.Values{
		"name":     {"Kevin Huang"},
		"email":    {"kevin@example.com"},
		"interest": {"OTA"},
	}

	firstReq := httptest.NewRequest(http.MethodPost, "/contact", strings.NewReader(form.Encode()))
	firstReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	firstReq.Header.Set("X-Forwarded-For", "203.0.113.10")
	firstReq.Header.Set("X-Real-IP", "203.0.113.20")
	firstReq.RemoteAddr = "198.51.100.10:1234"
	firstRec := httptest.NewRecorder()
	handler.ServeHTTP(firstRec, firstReq)

	if firstRec.Code != http.StatusOK {
		t.Fatalf("first status = %d, want 200", firstRec.Code)
	}

	secondReq := httptest.NewRequest(http.MethodPost, "/contact", strings.NewReader(form.Encode()))
	secondReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	secondReq.Header.Set("X-Forwarded-For", "203.0.113.11")
	secondReq.Header.Set("X-Real-IP", "203.0.113.21")
	secondReq.RemoteAddr = "198.51.100.10:1234"
	secondRec := httptest.NewRecorder()
	handler.ServeHTTP(secondRec, secondReq)

	if secondRec.Code != http.StatusTooManyRequests {
		t.Fatalf("second status = %d, want 429", secondRec.Code)
	}
	if len(store.leads) != 1 {
		t.Fatalf("stored leads = %d, want 1", len(store.leads))
	}
	if !strings.Contains(secondRec.Body.String(), "Too many requests from this address.") {
		t.Fatalf("response does not contain rate limit error: %s", secondRec.Body.String())
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

func TestAdminLeadsSupportsFilteringAndPagination(t *testing.T) {
	seed := make([]leads.Lead, 0, 27)
	for index := 1; index <= 27; index++ {
		seed = append(seed, leads.Lead{
			Name:     "Lead " + strconv.Itoa(index),
			Company:  "Company " + strconv.Itoa(index),
			Email:    "lead-" + strconv.Itoa(index) + "@example.com",
			Interest: "Provision",
			Message:  "seed",
		})
	}
	handler := testServerWithAdminToken(t, &memoryLeadStore{leads: seed}, "secret")

	req := httptest.NewRequest(http.MethodGet, "/admin/leads?token=secret&page=2", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}

	body := rec.Body.String()
	if !strings.Contains(body, "Page 2 of 2") {
		t.Fatalf("response does not contain page summary: %s", body)
	}
	if !strings.Contains(body, "Showing 26-27 of 27 leads") {
		t.Fatalf("response does not contain range summary: %s", body)
	}
	if !strings.Contains(body, "lead-2@example.com") {
		t.Fatalf("response does not contain second page lead: %s", body)
	}
	if strings.Contains(body, "lead-27@example.com") {
		t.Fatalf("response unexpectedly contains first page lead: %s", body)
	}
	if !strings.Contains(body, "/admin/leads?token=secret") {
		t.Fatalf("response does not contain previous-page link: %s", body)
	}
}

func TestAdminLeadsFiltersAndCSVExportRespectActiveFilters(t *testing.T) {
	store := &memoryLeadStore{leads: []leads.Lead{
		{
			Name:     "Alpha",
			Company:  "Acme",
			Email:    "alpha@example.com",
			Interest: "Provision",
		},
		{
			Name:     "Beta",
			Company:  "Acme Labs",
			Email:    "beta@example.com",
			Interest: "OTA",
		},
		{
			Name:     "Gamma",
			Company:  "Zenith",
			Email:    "gamma@example.com",
			Interest: "OTA",
		},
	}}
	handler := testServerWithAdminToken(t, store, "secret")

	req := httptest.NewRequest(http.MethodGet, "/admin/leads?token=secret&email=beta&company=acme&interest=ota", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}

	body := rec.Body.String()
	if !strings.Contains(body, "beta@example.com") {
		t.Fatalf("response does not contain filtered lead: %s", body)
	}
	if strings.Contains(body, "alpha@example.com") || strings.Contains(body, "gamma@example.com") {
		t.Fatalf("response contains unfiltered leads: %s", body)
	}
	if !strings.Contains(body, "/admin/leads.csv?company=acme&amp;email=beta&amp;interest=ota&amp;token=secret") {
		t.Fatalf("response does not preserve filters in csv link: %s", body)
	}

	csvReq := httptest.NewRequest(http.MethodGet, "/admin/leads.csv?token=secret&email=beta&company=acme&interest=ota", nil)
	csvRec := httptest.NewRecorder()
	handler.ServeHTTP(csvRec, csvReq)

	if csvRec.Code != http.StatusOK {
		t.Fatalf("csv status = %d, want 200", csvRec.Code)
	}
	if !strings.Contains(csvRec.Body.String(), "beta@example.com") {
		t.Fatalf("csv does not contain filtered lead: %s", csvRec.Body.String())
	}
	if strings.Contains(csvRec.Body.String(), "alpha@example.com") || strings.Contains(csvRec.Body.String(), "gamma@example.com") {
		t.Fatalf("csv contains unfiltered leads: %s", csvRec.Body.String())
	}
}

func TestAdminLeadsNoMatchEmptyStateMentionsFilters(t *testing.T) {
	handler := testServerWithAdminToken(t, &memoryLeadStore{leads: []leads.Lead{
		{
			Name:     "Alpha",
			Company:  "Acme",
			Email:    "alpha@example.com",
			Interest: "Provision",
		},
	}}, "secret")

	req := httptest.NewRequest(http.MethodGet, "/admin/leads?token=secret&email=missing", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}

	body := rec.Body.String()
	if !strings.Contains(body, "No leads match the current filters.") {
		t.Fatalf("response does not contain filtered empty-state message: %s", body)
	}
	if strings.Contains(body, "No leads yet.") {
		t.Fatalf("response unexpectedly contains generic empty-state message: %s", body)
	}
}
