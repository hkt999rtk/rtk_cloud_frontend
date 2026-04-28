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

func TestLayoutIncludesSkipLinkToMainContent(t *testing.T) {
	handler := testServer(t, &memoryLeadStore{})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}

	body := rec.Body.String()
	for _, want := range []string{
		`<a class="skip-link" href="#main-content">Skip to main content</a>`,
		`<main id="main-content" tabindex="-1">`,
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
		`<meta name="description" content="Upload firmware, extract release metadata, target staged rollouts, and manage dynamic OTA jobs with force, normal, and user-controlled policies.">`,
		`<meta property="og:url" content="http://example.com/features/ota">`,
		`<meta name="twitter:title" content="OTA | Realtek Connect&#43;">`,
	} {
		if !strings.Contains(body, want) {
			t.Fatalf("response does not contain %q: %s", want, body)
		}
	}
}

func TestOTAFeaturePageIncludesProductionDetail(t *testing.T) {
	handler := testServer(t, &memoryLeadStore{})

	req := httptest.NewRequest(http.MethodGet, "/features/ota", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}

	body := rec.Body.String()
	for _, want := range []string{
		"Upload signed firmware images and extract embedded project, version, model, checksum, and release-note metadata.",
		"Attach rollout notes, force, normal, or user-controlled install policy, and maintenance-window guidance before approval.",
		"Target by product family, hardware model, current firmware version, customer tier, region, or support cohort.",
		"Validate project and version compatibility before devices accept a package.",
		"Choose the delivery mode that fits the release",
		"<th scope=\"col\">Strategy</th>",
		"Force, normal, scheduled, user-controlled, and time-window rollout modes",
		"Dynamic OTA keeps device eligibility aligned with the latest approved campaign even when endpoints reconnect later.",
		"Force",
		"Normal",
		"Cancel active waves and archive completed campaigns without losing audit history.",
	} {
		if !strings.Contains(body, want) {
			t.Fatalf("response does not contain %q: %s", want, body)
		}
	}
}

func TestFleetManagementFeatureCoversAdminOperationsScope(t *testing.T) {
	handler := testServer(t, &memoryLeadStore{})

	req := httptest.NewRequest(http.MethodGet, "/features/fleet-management", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}

	body := rec.Body.String()
	for _, want := range []string{
		"Node registration, certificate provisioning, device registry, OTA job coordination, batch operations, and operator widgets for commercial fleets.",
		"Record serial number, model, MAC, factory lot, and claim state when a node is first registered into the platform catalog.",
		"Issue bootstrap certificates or device credentials, then support rotation or revocation workflows when products are repaired, replaced, or reworked.",
		"Search the device registry by region, firmware, product family, installer, or customer account and save groups for repeat operations.",
		"Coordinate firmware images and OTA jobs from the same operations surface so release managers can move from device search to rollout action without spreadsheet handoffs.",
		"Summarize activation counts, firmware mix, online-versus-offline ratios, alert backlogs, and support escalations in operator-facing statistics widgets.",
		"The existing /admin/leads page only covers website sales leads; the future IoT platform admin console described here remains a public product narrative rather than a shipped control plane in this repo.",
		"Map each admin workflow to the right platform boundary",
		"<th scope=\"col\">Workflow</th>",
		"Node registration",
		"Device registry",
		"Release operations",
		"Statistics widgets",
		"the shipped admin endpoint remains /admin/leads for website sales workflow only.",
	} {
		if !strings.Contains(body, want) {
			t.Fatalf("response does not contain %q: %s", want, body)
		}
	}
}

func TestHomeAndSmartHomeFeatureCoverEndUserWorkflow(t *testing.T) {
	handler := testServer(t, &memoryLeadStore{})

	homeReq := httptest.NewRequest(http.MethodGet, "/", nil)
	homeRec := httptest.NewRecorder()
	handler.ServeHTTP(homeRec, homeReq)

	if homeRec.Code != http.StatusOK {
		t.Fatalf("home status = %d, want 200", homeRec.Code)
	}
	if !strings.Contains(homeRec.Body.String(), `/features/smart-home`) {
		t.Fatalf("home page does not link to smart-home feature: %s", homeRec.Body.String())
	}

	req := httptest.NewRequest(http.MethodGet, "/features/smart-home", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}

	body := rec.Body.String()
	for _, want := range []string{
		"Smart Home Experience",
		"Remote control, local control fallback, schedules, scenes, grouping, device sharing, push notifications, and alerts for connected home products.",
		"Use remote control for away-from-home power, mode, and status changes when devices stay connected through the Realtek Connect&#43; cloud path.",
		"Keep local control available on the home network so core actions can stay responsive during WAN degradation or when products intentionally prioritize nearby control.",
		"Create recurring schedules around daily routines, quiet hours, occupancy assumptions, or energy-saving windows.",
		"Bundle scenes so users can trigger coordinated actions across lights, climate, appliances, or custom device categories from one tap.",
		"Group devices by room, home, or product set so the app can present household-level control instead of one-node-at-a-time management.",
		"Support node sharing so primary owners can invite family members, installers, or temporary guests with bounded access expectations.",
		"Use push notifications for onboarding completion, automation results, offline alerts, abnormal events, and OTA prompts that need the user back in the app.",
		"Map the home experience to the right control pattern",
		"<th scope=\"col\">Workflow</th>",
		"Remote control",
		"Local control",
		"Schedules and scenes",
		"Grouping and sharing",
		"Push notifications and alerts",
		"this repo does not ship the native control client",
	} {
		if !strings.Contains(body, want) {
			t.Fatalf("response does not contain %q: %s", want, body)
		}
	}
}

func TestUserManagementFeatureClarifiesPlatformScope(t *testing.T) {
	handler := testServer(t, &memoryLeadStore{})

	req := httptest.NewRequest(http.MethodGet, "/features/user-management", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}

	body := rec.Body.String()
	for _, want := range []string{
		"User Management",
		"sign up and sign in",
		"One-time password verification",
		"Third-party login and account-linking paths",
		"Forgot-password, change-password, and session-management controls",
		"Account deletion and retention workflows",
		"This website does not expose end-user sign-in or account management flows today.",
	} {
		if !strings.Contains(body, want) {
			t.Fatalf("response does not contain %q: %s", want, body)
		}
	}
}

func TestPrivateCloudFeatureCoversCommercialDeploymentPaths(t *testing.T) {
	handler := testServer(t, &memoryLeadStore{})

	req := httptest.NewRequest(http.MethodGet, "/features/private-cloud", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}

	body := rec.Body.String()
	for _, want := range []string{
		"Compare evaluation and private commercial operating models",
		"Public evaluation versus dedicated private commercial deployment",
		"Transition to a dedicated deployment once product teams need tenant isolation, formal support processes, and customer-specific change windows.",
		"Offer custom domains and branded entry points so the deployment can align with the customer&#39;s DNS, certificate, and support model.",
		"Choose regional placement around residency, latency, and operational coverage requirements instead of forcing every product through one public region.",
		"Use release promotion, maintenance windows, and rollback checkpoints to move from pilot tenants into production operations safely.",
		"Production TLS still terminates at a reverse proxy, ingress, or deployment platform in front of the Go website runtime.",
		"<th scope=\"col\">Model</th>",
		"Managed private deployment",
		"Customer-operated private region",
	} {
		if !strings.Contains(body, want) {
			t.Fatalf("response does not contain %q: %s", want, body)
		}
	}
}

func TestAppSDKFeatureCoversMobileDeliveryPaths(t *testing.T) {
	handler := testServer(t, &memoryLeadStore{})

	req := httptest.NewRequest(http.MethodGet, "/features/app-sdk", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}

	body := rec.Body.String()
	for _, want := range []string{
		"Deliver branded mobile apps without rebuilding the connected product stack",
		"Cover shared onboarding, authentication, device control, and account-linking primitives through iOS and Android SDK layers instead of promising a full client framework in this repo.",
		"Use a sample app to accelerate white-label or branded launches while preserving room for custom navigation, design systems, and product-specific device flows.",
		"Plan push notifications around onboarding completion, shared-device events, OTA prompts, alerts, and support workflows that need deep links back into the branded app.",
		"Coordinate bundle identifiers, signing assets, store metadata, review checklists, and staged rollout plans for both the App Store and Google Play.",
		"Discuss App SDK",
		"Choose the mobile delivery path that fits launch speed and brand control",
		"<th scope=\"col\">Delivery path</th>",
		"Rebranded starter app",
		"Custom app on shared SDK",
		"not a shipped mobile framework",
	} {
		if !strings.Contains(body, want) {
			t.Fatalf("response does not contain %q: %s", want, body)
		}
	}
}

func TestIntegrationsFeatureCoversMatterAndEcosystemPaths(t *testing.T) {
	handler := testServer(t, &memoryLeadStore{})

	req := httptest.NewRequest(http.MethodGet, "/features/integrations", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}

	body := rec.Body.String()
	for _, want := range []string{
		"Position Realtek products inside the customer&#39;s chosen ecosystem",
		"Describe how devices can participate in a Matter Fabric while still keeping Realtek app, cloud, and support flows in scope where products need them.",
		"Cover Alexa and Google Assistant paths for products that need voice control, routine support, and ecosystem discovery alongside the branded app experience.",
		"Document REST APIs for authenticated product, support, and operations workflows that need request-response access to platform state.",
		"Position MQTT over TLS for policy-scoped telemetry, near-real-time command paths, and event fan-out into downstream infrastructure.",
		"Use webhooks for signed lifecycle, alert, and workflow events so CRM, ticketing, analytics, and fulfillment systems can react without polling.",
		"Choose the ecosystem contract that fits the product",
		"<th scope=\"col\">Path</th>",
		"Matter Fabric",
		"Voice assistants",
		"MQTT over TLS",
		"Webhooks",
		"without claiming every protocol surface is already live in this repository",
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

func TestInvalidContactPostUsesAccessibleErrorMarkup(t *testing.T) {
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

	body := rec.Body.String()
	for _, want := range []string{
		`<div class="form-error-summary" role="alert" aria-labelledby="contact-errors-title" tabindex="-1">`,
		`<a href="#contact-name">Name is required.</a>`,
		`<a href="#contact-email">Enter a valid email address.</a>`,
		`<a href="#contact-interest">Select an area of interest.</a>`,
		`id="contact-name"`,
		`aria-describedby="contact-name-error"`,
		`id="contact-email"`,
		`aria-describedby="contact-email-error"`,
		`id="contact-interest"`,
		`aria-describedby="contact-interest-error"`,
	} {
		if !strings.Contains(body, want) {
			t.Fatalf("response does not contain %q: %s", want, body)
		}
	}
	if strings.Count(body, `aria-invalid="true"`) != 3 {
		t.Fatalf("response should flag three invalid fields: %s", body)
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
