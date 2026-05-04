package web

import (
	"context"
	"io"
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
	cfg := testConfig(store)
	cfg.AdminToken = adminToken
	server := newTestServer(t, cfg)
	return server.Routes()
}

func testConfig(store LeadStore) Config {
	return Config{
		TemplatesDir: "../../templates",
		StaticDir:    "../../static",
		LeadStore:    store,
	}
}

func testServerWithConfig(t *testing.T, cfg Config) http.Handler {
	t.Helper()
	if cfg.TemplatesDir == "" {
		cfg.TemplatesDir = "../../templates"
	}
	if cfg.StaticDir == "" {
		cfg.StaticDir = "../../static"
	}
	if cfg.LeadStore == nil {
		cfg.LeadStore = &memoryLeadStore{}
	}
	server := newTestServer(t, cfg)
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
	paths := []string{"/", "/docs", "/features", "/contact", "/privacy", "/healthz", "/robots.txt", "/sitemap.xml"}
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

func TestLocalizedPublicRoutesReturnOK(t *testing.T) {
	handler := testServer(t, &memoryLeadStore{})
	paths := []string{
		"/zh-tw",
		"/zh-tw/docs",
		"/zh-tw/docs/apis",
		"/zh-tw/features",
		"/zh-tw/features/provision",
		"/zh-tw/contact",
		"/zh-tw/privacy",
		"/zh-cn",
		"/zh-cn/docs",
		"/zh-cn/docs/apis",
		"/zh-cn/features",
		"/zh-cn/features/provision",
		"/zh-cn/contact",
		"/zh-cn/privacy",
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

func TestLocalizedHomeIncludesLangSwitcherAndAlternates(t *testing.T) {
	handler := testServer(t, &memoryLeadStore{})

	req := httptest.NewRequest(http.MethodGet, "/zh-tw/features/provision", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}

	body := rec.Body.String()
	for _, want := range []string{
		`<html lang="zh-Hant">`,
		`<link rel="canonical" href="http://example.com/zh-tw/features/provision">`,
		`hreflang="en" href="http://example.com/features/provision"`,
		`hreflang="zh-Hant" href="http://example.com/zh-tw/features/provision"`,
		`hreflang="zh-Hans" href="http://example.com/zh-cn/features/provision"`,
		`hreflang="x-default" href="http://example.com/features/provision"`,
		`href="http://example.com/zh-tw/features/provision" aria-current="true">繁體中文</a>`,
		"Provision 配網",
		"以合約支撐的基礎來描述裝置導入。",
	} {
		if !strings.Contains(body, want) {
			t.Fatalf("response does not contain %q: %s", want, body)
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

func TestHomeIncludesLocalizedBrandFilmEmbed(t *testing.T) {
	handler := testServer(t, &memoryLeadStore{})

	tests := []struct {
		path       string
		title      string
		body       string
		point      string
		videoTitle string
		fallback   string
	}{
		{
			path:       "/",
			title:      "Built on Realtek&#39;s connected intelligence.",
			body:       "Realtek Connect&#43; extends a semiconductor and connectivity foundation into a cloud platform story",
			point:      "Semiconductor foundation",
			videoTitle: `title="Realtek corporate brand film"`,
			fallback:   "Your browser does not support the video tag.",
		},
		{
			path:       "/zh-tw/",
			title:      "建立在 Realtek 的連網智慧之上。",
			body:       "Realtek Connect&#43; 將半導體與連線技術基礎延伸為雲端平台敘事",
			point:      "半導體技術基礎",
			videoTitle: `title="Realtek 企業形象影片"`,
			fallback:   "你的瀏覽器不支援 video 標籤。",
		},
		{
			path:       "/zh-cn/",
			title:      "建立在 Realtek 的連网智慧之上。",
			body:       "Realtek Connect&#43; 将半导体与連线技术基础延伸为云端平台敘事",
			point:      "半导体技术基础",
			videoTitle: `title="Realtek 企业形象影片"`,
			fallback:   "你的浏览器不支援 video 标签。",
		},
	}

	for _, tc := range tests {
		req := httptest.NewRequest(http.MethodGet, tc.path, nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("%s status = %d, want 200", tc.path, rec.Code)
		}

		body := rec.Body.String()
		for _, want := range []string{
			`<section class="section brand-film">`,
			`<video controls preload="metadata" poster="/static/assets/realtek-brand-film-poster.jpg"`,
			`<source src="/static/assets/realtek-brand-film.mp4" type="video/mp4">`,
			tc.videoTitle,
			tc.title,
			tc.body,
			tc.point,
			tc.fallback,
		} {
			if !strings.Contains(body, want) {
				t.Fatalf("%s response does not contain %q: %s", tc.path, want, body)
			}
		}
		for _, unwanted := range []string{
			`<iframe`,
			`youtube-nocookie.com`,
			`data-video-id="QqC06634wcI"`,
		} {
			if strings.Contains(body, unwanted) {
				t.Fatalf("%s initial response should not contain %q: %s", tc.path, unwanted, body)
			}
		}

		architectureIndex := strings.Index(body, `<section class="section architecture"`)
		filmIndex := strings.Index(body, `<section class="section brand-film">`)
		deploymentIndex := strings.Index(body, `<section class="section deployment-section">`)
		if architectureIndex == -1 || filmIndex == -1 || deploymentIndex == -1 {
			t.Fatalf("%s missing expected home sections", tc.path)
		}
		if !(architectureIndex < filmIndex && filmIndex < deploymentIndex) {
			t.Fatalf("%s brand film section should be between architecture and deployment", tc.path)
		}
	}
}

func TestPrivacyPagesIncludeLocalizedNoticeAndMetadata(t *testing.T) {
	handler := testServer(t, &memoryLeadStore{})

	tests := []struct {
		path      string
		lang      string
		canonical string
		title     string
		body      string
		snippets  []string
	}{
		{
			path:      "/privacy",
			lang:      "en",
			canonical: `http://example.com/privacy`,
			title:     "Privacy notice for Realtek Connect&#43; website inquiries.",
			body:      "Website leads are intended to be retained for up to 24 months",
			snippets: []string{
				"first-party SQLite analytics",
				"page_view",
				"click_cta",
				"scroll",
				"engaged",
				"referrer origin only",
				"ephemeral session id",
				"Raw analytics event rows are retained for 90 days",
				"third-party analytics services",
				"advertising pixels",
				"fingerprinting scripts",
			},
		},
		{
			path:      "/zh-tw/privacy",
			lang:      "zh-Hant",
			canonical: `http://example.com/zh-tw/privacy`,
			title:     "Realtek Connect&#43; 網站詢問隱私權聲明。",
			body:      "網站 leads 預期最多保存 24 個月",
			snippets: []string{
				"第一方 SQLite analytics",
				"page_view",
				"click_cta",
				"scroll",
				"engaged events",
				"referrer origin only",
				"ephemeral session id",
				"保存 90 天",
				"第三方 analytics services",
				"advertising pixels",
				"fingerprinting scripts",
			},
		},
		{
			path:      "/zh-cn/privacy",
			lang:      "zh-Hans",
			canonical: `http://example.com/zh-cn/privacy`,
			title:     "Realtek Connect&#43; 网站询问隐私权声明。",
			body:      "网站 leads 预期最多保存 24 个月",
			snippets: []string{
				"第一方 SQLite analytics",
				"page_view",
				"click_cta",
				"scroll",
				"engaged events",
				"referrer origin only",
				"ephemeral session id",
				"保存 90 天",
				"第三方 analytics services",
				"advertising pixels",
				"fingerprinting scripts",
			},
		},
	}

	for _, tc := range tests {
		req := httptest.NewRequest(http.MethodGet, tc.path, nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Fatalf("%s status = %d, want 200", tc.path, rec.Code)
		}
		body := rec.Body.String()
		for _, want := range []string{
			`<html lang="` + tc.lang + `">`,
			`<link rel="canonical" href="` + tc.canonical + `">`,
			`hreflang="en" href="http://example.com/privacy"`,
			`hreflang="zh-Hant" href="http://example.com/zh-tw/privacy"`,
			`hreflang="zh-Hans" href="http://example.com/zh-cn/privacy"`,
			tc.title,
			tc.body,
			"privacy@example.com",
			"local MP4",
		} {
			if !strings.Contains(body, want) {
				t.Fatalf("%s response does not contain %q: %s", tc.path, want, body)
			}
		}
		for _, want := range tc.snippets {
			if !strings.Contains(body, want) {
				t.Fatalf("%s response does not contain analytics privacy snippet %q: %s", tc.path, want, body)
			}
		}
	}
}

func TestPublicBaseURLOverridesGeneratedAbsoluteURLs(t *testing.T) {
	handler := testServerWithConfig(t, Config{
		LeadStore:     &memoryLeadStore{},
		PublicBaseURL: "https://webtest.mgmeet.io/",
	})

	req := httptest.NewRequest(http.MethodGet, "/zh-tw/features/provision", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}

	body := rec.Body.String()
	for _, want := range []string{
		`<link rel="canonical" href="https://webtest.mgmeet.io/zh-tw/features/provision">`,
		`hreflang="en" href="https://webtest.mgmeet.io/features/provision"`,
		`hreflang="zh-Hant" href="https://webtest.mgmeet.io/zh-tw/features/provision"`,
		`hreflang="zh-Hans" href="https://webtest.mgmeet.io/zh-cn/features/provision"`,
		`<meta property="og:url" content="https://webtest.mgmeet.io/zh-tw/features/provision">`,
		`<meta property="og:image" content="https://webtest.mgmeet.io/static/assets/connectplus-hero.png">`,
	} {
		if !strings.Contains(body, want) {
			t.Fatalf("response does not contain %q: %s", want, body)
		}
	}

	req = httptest.NewRequest(http.MethodGet, "/robots.txt", nil)
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("robots status = %d, want 200", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "Sitemap: https://webtest.mgmeet.io/sitemap.xml") {
		t.Fatalf("robots does not use public base URL: %s", rec.Body.String())
	}

	req = httptest.NewRequest(http.MethodGet, "/sitemap.xml", nil)
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("sitemap status = %d, want 200", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), `<loc>https://webtest.mgmeet.io/zh-cn/contact</loc>`) {
		t.Fatalf("sitemap does not use public base URL: %s", rec.Body.String())
	}
}

func TestAssetFingerprintsAreOptional(t *testing.T) {
	handler := testServer(t, &memoryLeadStore{})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	body := rec.Body.String()
	if !strings.Contains(body, `href="/static/styles.css"`) {
		t.Fatalf("default stylesheet path changed: %s", body)
	}
	if strings.Contains(body, `/static/styles.css?v=`) {
		t.Fatalf("default response should not fingerprint assets: %s", body)
	}

	handler = testServerWithConfig(t, Config{
		LeadStore:               &memoryLeadStore{},
		EnableAssetFingerprints: true,
		EnableCDNCacheHeaders:   false,
		DisableSearchIndexing:   false,
		PublicBaseURL:           "https://webtest.mgmeet.io",
	})

	req = httptest.NewRequest(http.MethodGet, "/", nil)
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("fingerprinted status = %d, want 200", rec.Code)
	}
	body = rec.Body.String()
	for _, want := range []string{
		`href="/static/styles.css?v=`,
		`src="/static/assets/realtek-logo.png?v=`,
		`src="/static/assets/connectplus-hero-v2.jpg?v=`,
		`<meta property="og:image" content="https://webtest.mgmeet.io/static/assets/connectplus-hero.png?v=`,
	} {
		if !strings.Contains(body, want) {
			t.Fatalf("fingerprinted response does not contain %q: %s", want, body)
		}
	}
}

func TestCDNCacheHeadersAreOptional(t *testing.T) {
	handler := testServer(t, &memoryLeadStore{})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if got := rec.Header().Get("Cache-Control"); got != "" {
		t.Fatalf("default home Cache-Control = %q, want empty", got)
	}

	req = httptest.NewRequest(http.MethodGet, "/static/styles.css", nil)
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if got := rec.Header().Get("Cache-Control"); got != "" {
		t.Fatalf("default static Cache-Control = %q, want empty", got)
	}

	handler = testServerWithConfig(t, Config{
		LeadStore:               &memoryLeadStore{},
		AdminToken:              "secret",
		EnableCDNCacheHeaders:   true,
		DisableSearchIndexing:   false,
		EnableAssetFingerprints: false,
	})

	tests := []struct {
		method string
		path   string
		want   string
	}{
		{method: http.MethodGet, path: "/", want: "no-store"},
		{method: http.MethodGet, path: "/zh-tw/contact", want: "no-store"},
		{method: http.MethodPost, path: "/contact", want: "no-store"},
		{method: http.MethodGet, path: "/admin/leads", want: "no-store"},
		{method: http.MethodGet, path: "/healthz", want: "no-store"},
		{method: http.MethodGet, path: "/robots.txt", want: "public, max-age=300"},
		{method: http.MethodGet, path: "/sitemap.xml", want: "public, max-age=300"},
		{method: http.MethodGet, path: "/static/styles.css", want: "public, max-age=31536000, immutable"},
	}

	for _, tc := range tests {
		var body io.Reader
		if tc.method == http.MethodPost {
			form := url.Values{
				"name":     {"Kevin Huang"},
				"email":    {"kevin@example.com"},
				"interest": {"ota"},
			}
			body = strings.NewReader(form.Encode())
		}
		req := httptest.NewRequest(tc.method, tc.path, body)
		if tc.method == http.MethodPost {
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		}
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
		if got := rec.Header().Get("Cache-Control"); got != tc.want {
			t.Fatalf("%s %s Cache-Control = %q, want %q", tc.method, tc.path, got, tc.want)
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
		`<a href="/privacy">Privacy</a>`,
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
		`<meta name="description" content="Firmware upload, catalog, target enablement, rollout status, report, cancel, and download are available foundations; advanced campaign policy remains contract-defined follow-up work.">`,
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
		"Firmware upload, catalog, target enablement, rollout status, report, cancel, and download are available foundations; advanced campaign policy remains contract-defined follow-up work.",
		"Firmware campaign interface contract",
		`href="https://github.com/hkt999rtk/rtk_cloud_contracts_doc/blob/main/FIRMWARE_CAMPAIGN.md"`,
		"Use the current firmware lifecycle as the implementation boundary",
		"Describe publish, enablement, whitelist, rollout query/report, cancel, and download behavior as the available firmware lifecycle foundation.",
		"Scheduled and time-window OTA are contract-defined policy concepts until backend enforcement and SDK handling are documented as available.",
		"User-consent-required OTA is a policy flag in phase one, not a shipped mobile UX or app-side consent flow.",
		"Approval workflow, operator dashboards, analytics, and success-rate reporting are roadmap capabilities, not phase-one availability claims.",
		"Map each OTA concept to the right implementation status",
		"<th scope=\"col\">Status</th>",
		"Available foundation",
		"Integration-ready contract",
		"Roadmap campaign management",
		"Staged percentage rollout and automatic cohort ramping stay out of the available feature list until a campaign engine implements them.",
	} {
		if !strings.Contains(body, want) {
			t.Fatalf("response does not contain %q: %s", want, body)
		}
	}
}

func TestOTAFeaturePageDoesNotPromoteCampaignPolicyScope(t *testing.T) {
	handler := testServer(t, &memoryLeadStore{})

	req := httptest.NewRequest(http.MethodGet, "/features/ota", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}

	body := rec.Body.String()
	assertTableRowStatus(t, body, "Scheduled policy", "Integration-ready contract")
	assertTableRowStatus(t, body, "Time-window policy", "Integration-ready contract")
	assertTableRowStatus(t, body, "User-consent policy", "Integration-ready contract")
	assertTableRowStatus(t, body, "Archive", "Roadmap campaign management")

	for _, tc := range []struct {
		concept string
		status  string
	}{
		{concept: "Scheduled policy", status: "Available foundation"},
		{concept: "Time-window policy", status: "Available foundation"},
		{concept: "User-consent policy", status: "Available foundation"},
		{concept: "Archive", status: "Available foundation"},
	} {
		assertTableRowDoesNotHaveStatus(t, body, tc.concept, tc.status)
	}
}

func TestProvisionFeatureAlignsPublicCopyWithContractStatus(t *testing.T) {
	handler := testServer(t, &memoryLeadStore{})

	req := httptest.NewRequest(http.MethodGet, "/features/provision", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}

	body := rec.Body.String()
	for _, want := range []string{
		"Cloud registry and activation foundations are contract-backed; local Wi-Fi/BLE onboarding, claim UX, transfer/reset policy, and product readiness remain integration or roadmap scope.",
		"Product onboarding interface contract",
		`href="https://github.com/hkt999rtk/rtk_cloud_contracts_doc/blob/main/PRODUCT_ONBOARDING.md"`,
		"Cloud-side provisioning is the implemented contract boundary",
		"Account-side device registration, cross-service provisioning requests, video activation results, scoped token issuance, and owner transport readiness are the public cloud-side behaviors to discuss today.",
		"Claim material has a defined interface, not final ownership policy",
		"BLE provisioning, SoftAP provisioning, local Wi-Fi credential transport, QR onboarding UX, ECDH or challenge-response handshakes, and manufacturing CA policy are not yet stable website-available implementation claims.",
		"Separate what is available, integration-ready, and roadmap",
		"<th scope=\"col\">Public status</th>",
		"Available foundation",
		"Integration-ready",
		"Roadmap",
		"Transfer, reset, and product readiness",
	} {
		if !strings.Contains(body, want) {
			t.Fatalf("response does not contain %q: %s", want, body)
		}
	}
}

func assertTableRowStatus(t *testing.T, body, concept, status string) {
	t.Helper()

	row := tableRowForConcept(t, body, concept)
	want := "<td>" + status + "</td>"
	if !strings.Contains(row, want) {
		t.Fatalf("table row for %q does not contain status %q: %s", concept, status, row)
	}
}

func assertTableRowDoesNotHaveStatus(t *testing.T, body, concept, status string) {
	t.Helper()

	row := tableRowForConcept(t, body, concept)
	forbidden := "<td>" + status + "</td>"
	if strings.Contains(row, forbidden) {
		t.Fatalf("table row for %q unexpectedly contains status %q: %s", concept, status, row)
	}
}

func tableRowForConcept(t *testing.T, body, concept string) string {
	t.Helper()

	cell := "<td>" + concept + "</td>"
	index := strings.Index(body, cell)
	if index < 0 {
		t.Fatalf("response does not contain table concept %q: %s", concept, body)
	}
	start := strings.LastIndex(body[:index], "<tr>")
	end := strings.Index(body[index:], "</tr>")
	if start < 0 || end < 0 {
		t.Fatalf("response does not contain a complete table row for %q: %s", concept, body)
	}
	return body[start : index+end+len("</tr>")]
}

func TestFeaturePagesUseLocalVisualAssets(t *testing.T) {
	handler := testServer(t, &memoryLeadStore{})

	tests := []struct {
		path string
		src  string
		alt  string
	}{
		{
			path: "/features/provision",
			src:  `/static/assets/feature-provision-flow.jpg`,
			alt:  `alt="Provisioning dashboard concept with mobile pairing steps, QR onboarding, and device activation status cards."`,
		},
		{
			path: "/features/ota",
			src:  `/static/assets/feature-ota-control-center.jpg`,
			alt:  `alt="Firmware rollout control center with staged release timeline, device cohorts, and OTA job status cards."`,
		},
		{
			path: "/features/fleet-management",
			src:  `/static/assets/feature-fleet-management.png`,
			alt:  `alt="Fleet management dashboard with connected device groups, health status tiles, tags, and batch operation queue."`,
		},
		{
			path: "/features/smart-home",
			src:  `/static/assets/feature-smart-home-experience.png`,
			alt:  `alt="Smart home app control surface with connected home devices, scenes, schedules, grouping, and notification cards."`,
		},
		{
			path: "/features/user-management",
			src:  `/static/assets/feature-user-management.png`,
			alt:  `alt="User management console with profile cards, security verification, sharing permissions, and account lifecycle controls."`,
		},
		{
			path: "/features/app-sdk",
			src:  `/static/assets/feature-app-sdk.png`,
			alt:  `alt="Mobile app SDK workspace with app screens, code modules, push notification blocks, and publishing checklist."`,
		},
		{
			path: "/features/insights",
			src:  `/static/assets/feature-insights-dashboard.jpg`,
			alt:  `alt="Operations insights dashboard with fleet health charts, alert cards, and device telemetry summaries."`,
		},
		{
			path: "/features/private-cloud",
			src:  `/static/assets/feature-private-cloud-architecture.jpg`,
			alt:  `alt="Private cloud architecture showing container and VM workloads running across multiple cloud providers and on-premises data centers."`,
		},
		{
			path: "/features/integrations",
			src:  `/static/assets/feature-integrations.png`,
			alt:  `alt="Integration hub connecting generic Matter, voice assistant, REST API, MQTT over TLS, webhook, app, and enterprise system endpoints."`,
		},
	}

	for _, tc := range tests {
		req := httptest.NewRequest(http.MethodGet, tc.path, nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("%s status = %d, want 200", tc.path, rec.Code)
		}

		body := rec.Body.String()
		if !strings.Contains(body, `class="feature-visual"`) {
			t.Fatalf("%s missing feature visual wrapper: %s", tc.path, body)
		}
		if !strings.Contains(body, `src="`+tc.src+`"`) {
			t.Fatalf("%s missing local asset %q: %s", tc.path, tc.src, body)
		}
		if !strings.Contains(body, tc.alt) {
			t.Fatalf("%s missing alt text %q: %s", tc.path, tc.alt, body)
		}

		assetReq := httptest.NewRequest(http.MethodGet, tc.src, nil)
		assetRec := httptest.NewRecorder()
		handler.ServeHTTP(assetRec, assetReq)
		if assetRec.Code != http.StatusOK {
			t.Fatalf("%s asset %s status = %d, want 200", tc.path, tc.src, assetRec.Code)
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

func TestHomeDeploySectionDisclosesEvaluationLimits(t *testing.T) {
	handler := testServer(t, &memoryLeadStore{})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}

	body := rec.Body.String()
	for _, want := range []string{
		"5 devices by default",
		"raise up to 200 on request",
		"No expiry",
		"See plans &amp; limits",
		`href="/features/private-cloud"`,
	} {
		if !strings.Contains(body, want) {
			t.Fatalf("home page missing %q: %s", want, body)
		}
	}
}

func TestContactFormRendersCanonicalInquiryOptions(t *testing.T) {
	handler := testServer(t, &memoryLeadStore{})

	req := httptest.NewRequest(http.MethodGet, "/contact", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}

	body := rec.Body.String()
	for _, want := range []string{
		`value="evaluation-access"`,
		`value="commercial-deployment"`,
		`value="partnership"`,
		`value="technical-question"`,
		`value="other"`,
		"Evaluation access",
		"Commercial deployment",
	} {
		if !strings.Contains(body, want) {
			t.Fatalf("contact form missing %q: %s", want, body)
		}
	}

	// The dropdown must no longer offer feature slugs as inquiry types.
	for _, mustNot := range []string{
		`value="provision"`,
		`value="ota"`,
		`value="private-cloud"`,
	} {
		if strings.Contains(body, mustNot) {
			t.Fatalf("contact form unexpectedly still offers feature slug %q", mustNot)
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
		"VM/container deployment on GCP, Azure, AWS, or on-premises",
		"Standard container and VM workloads",
		"no serverless runtime requirement",
		"Transition to a dedicated private deployment once product teams need tenant isolation, formal support processes, and customer-specific change windows.",
		"Custom domains and branded entry points let the deployment align with the customer&#39;s DNS, certificate, and support model.",
		"Choose regional placement around residency, latency, and operational coverage requirements",
		"Use release promotion, maintenance windows, and rollback checkpoints to move from pilot tenants into production operations safely.",
		"Is there a cloud vendor requirement? No.",
		// Plans & Limits disclosure
		"5-device default quota",
		"up to 200 devices on request",
		"Evaluation access does not expire",
		"no minimum scale for the commercial tier",
		// Pricing Factors disclosure (no price list, just the inputs)
		"How commercial pricing is shaped",
		"Fleet size",
		"Deployment topology",
		"Support coverage",
		"Customization scope",
		"Term length",
		// SDK Licensing (split from Support)
		"What you can build with",
		"open-source SDK release is planned at general availability",
		"platform backend stays a proprietary commercial product",
		// Support tier (split from SDK)
		"What support looks like at each tier",
		"Evaluation support is community-tier",
		"Commercial support is contract-defined",
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

func TestSearchIndexingDisabledAddsNoIndexSignals(t *testing.T) {
	server := newTestServer(t, Config{
		TemplatesDir:          "../../templates",
		StaticDir:             "../../static",
		LeadStore:             &memoryLeadStore{},
		DisableSearchIndexing: true,
	})
	handler := server.Routes()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	if got := rec.Header().Get("X-Robots-Tag"); got != "noindex, nofollow, noarchive" {
		t.Fatalf("X-Robots-Tag = %q, want noindex", got)
	}
	if !strings.Contains(rec.Body.String(), `<meta name="robots" content="noindex, nofollow, noarchive">`) {
		t.Fatalf("homepage does not contain noindex meta: %s", rec.Body.String())
	}

	req = httptest.NewRequest(http.MethodGet, "/robots.txt", nil)
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("robots status = %d, want 200", rec.Code)
	}
	body := rec.Body.String()
	for _, want := range []string{"User-agent: *", "Disallow: /"} {
		if !strings.Contains(body, want) {
			t.Fatalf("robots.txt does not contain %q: %s", want, body)
		}
	}
	if strings.Contains(body, "Sitemap:") {
		t.Fatalf("disabled robots.txt should not advertise sitemap: %s", body)
	}

	req = httptest.NewRequest(http.MethodGet, "/sitemap.xml", nil)
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("sitemap status = %d, want 404", rec.Code)
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
		`<loc>http://example.com/privacy</loc>`,
		`<loc>http://example.com/zh-tw/features/ota</loc>`,
		`<loc>http://example.com/zh-tw/privacy</loc>`,
		`<loc>http://example.com/zh-cn/docs/apis</loc>`,
		`<loc>http://example.com/zh-cn/contact</loc>`,
		`<loc>http://example.com/zh-cn/privacy</loc>`,
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

func TestUnknownLocalizedRoutesReturnNotFound(t *testing.T) {
	handler := testServer(t, &memoryLeadStore{})
	paths := []string{"/fr/features", "/zh-tw/features/not-found", "/zh-cn/docs/not-found"}

	for _, path := range paths {
		req := httptest.NewRequest(http.MethodGet, path, nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
		if rec.Code != http.StatusNotFound {
			t.Fatalf("%s status = %d, want 404", path, rec.Code)
		}
	}
}

func TestValidContactPostStoresLead(t *testing.T) {
	store := &memoryLeadStore{}
	handler := testServer(t, store)

	form := url.Values{
		"name":     {"Kevin Huang"},
		"company":  {"Realtek"},
		"email":    {"kevin@example.com"},
		"interest": {"evaluation-access"},
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

func TestLocalizedContactPostStoresCanonicalInterest(t *testing.T) {
	store := &memoryLeadStore{}
	handler := testServer(t, store)

	form := url.Values{
		"name":     {"Kevin Huang"},
		"company":  {"Realtek"},
		"email":    {"kevin@example.com"},
		"interest": {"commercial-deployment"},
		"message":  {"需要排程更新支援。"},
	}
	req := httptest.NewRequest(http.MethodPost, "/zh-tw/contact", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	if len(store.leads) != 1 {
		t.Fatalf("stored leads = %d, want 1", len(store.leads))
	}
	if got := store.leads[0].Interest; got != "commercial-deployment" {
		t.Fatalf("interest = %q, want canonical slug", got)
	}
	if !strings.Contains(rec.Body.String(), "你的 Realtek Connect&#43; 請求已記錄。") {
		t.Fatalf("response does not contain localized success message: %s", rec.Body.String())
	}
}

func TestLocalizedInvalidContactPostUsesLocalizedErrors(t *testing.T) {
	store := &memoryLeadStore{}
	handler := testServer(t, store)

	form := url.Values{
		"name":     {""},
		"email":    {"not-an-email"},
		"interest": {""},
	}
	req := httptest.NewRequest(http.MethodPost, "/zh-cn/contact", strings.NewReader(form.Encode()))
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
	for _, want := range []string{"姓名为必填栏位。", "请输入有效的 Email。", "请选择关注服务。"} {
		if !strings.Contains(body, want) {
			t.Fatalf("response does not contain %q: %s", want, body)
		}
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

func TestContactFormIncludesLocalizedPrivacyNotice(t *testing.T) {
	handler := testServer(t, &memoryLeadStore{})

	tests := []struct {
		path string
		href string
		text string
		link string
	}{
		{
			path: "/contact",
			href: `/privacy`,
			text: "By submitting this form, you understand that your inquiry will be handled according to the Realtek Connect&#43; privacy notice.",
			link: "Privacy notice",
		},
		{
			path: "/zh-tw/contact",
			href: `/zh-tw/privacy`,
			text: "送出此表單即表示你理解我們會依 Realtek Connect&#43; 隱私權聲明處理你的詢問資料。",
			link: "隱私權聲明",
		},
		{
			path: "/zh-cn/contact",
			href: `/zh-cn/privacy`,
			text: "送出此表单即表示你理解我们会依 Realtek Connect&#43; 隐私权声明处理你的询问资料。",
			link: "隐私权声明",
		},
	}

	for _, tc := range tests {
		req := httptest.NewRequest(http.MethodGet, tc.path, nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Fatalf("%s status = %d, want 200", tc.path, rec.Code)
		}
		body := rec.Body.String()
		for _, want := range []string{
			`<p class="privacy-note">` + tc.text,
			`<a href="` + tc.href + `">` + tc.link + `</a>`,
		} {
			if !strings.Contains(body, want) {
				t.Fatalf("%s response does not contain %q: %s", tc.path, want, body)
			}
		}
	}
}

func TestOversizedContactPostDoesNotStoreLead(t *testing.T) {
	store := &memoryLeadStore{}
	handler := testServer(t, store)

	form := url.Values{
		"name":     {strings.Repeat("N", leads.NameMaxLength+1)},
		"company":  {strings.Repeat("C", leads.CompanyMaxLength+1)},
		"email":    {"kevin@example.com"},
		"interest": {"evaluation-access"},
		"message":  {strings.Repeat("M", leads.MessageMaxLength+1)},
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
		"interest": {"evaluation-access"},
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
		"interest": {"evaluation-access"},
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
		"interest": {"evaluation-access"},
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
