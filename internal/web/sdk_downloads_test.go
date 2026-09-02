package web

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"realtek-connect/internal/sdkdownloads"
)

type webSDKStore struct{ objects map[string][]byte }

func (s webSDKStore) Get(_ context.Context, key string) ([]byte, error) { return s.objects[key], nil }
func (s webSDKStore) PresignGet(key string, _ time.Duration) (string, error) {
	return "https://objects.example/" + key + "?X-Amz-Signature=test", nil
}

func sdkDownloadTestService(t *testing.T) *sdkdownloads.Service {
	t.Helper()
	packages := []sdkdownloads.Artifact{}
	for _, slug := range []string{"native", "android", "javascript", "ios", "freertos-pro2"} {
		packages = append(packages, sdkdownloads.Artifact{Slug: slug, Title: strings.ToUpper(slug), Description: slug + " SDK", Filename: slug + ".tgz", ObjectKey: "sdk/releases/0.1.0-rc.2/artifacts/" + slug + ".tgz", SHA256: strings.Repeat("a", 64), SizeBytes: 1024, ValidationStatus: "PASS", Capabilities: []string{"WebRTC signaling"}, Limitations: []string{"No media runtime"}})
	}
	catalog := sdkdownloads.Catalog{Schema: sdkdownloads.CatalogSchema, Version: "0.1.0-rc.2", ReleaseTrain: "rtk-cloud-client-0.1.0-rc.2", CreatedAt: "2026-08-29T00:00:00Z", Distribution: "public-evaluation", SigningStatus: "not_configured", TermsVersion: "eval-v1", TermsObjectKey: "sdk/releases/0.1.0-rc.2/EVALUATION_LICENSE.txt", Packages: packages, CompleteBundle: sdkdownloads.Artifact{Slug: "all", Title: "Complete", Description: "All SDKs", Filename: "all.tgz", ObjectKey: "sdk/releases/0.1.0-rc.2/artifacts/all.tgz", SHA256: strings.Repeat("b", 64), SizeBytes: 2048, ValidationStatus: "PASS", Capabilities: []string{"All packages"}, Limitations: []string{"No media runtime"}}}
	catalogJSON, _ := json.Marshal(catalog)
	latestJSON := []byte(`{"schema":"rtk-portal-sdk-latest/v1","version":"0.1.0-rc.2","catalog_object_key":"sdk/releases/0.1.0-rc.2/catalog.json","terms_version":"eval-v1"}`)
	return sdkdownloads.NewService(webSDKStore{objects: map[string][]byte{"sdk/latest.json": latestJSON, "sdk/releases/0.1.0-rc.2/catalog.json": catalogJSON, catalog.TermsObjectKey: []byte("Approved evaluation terms")}}, "", time.Minute)
}

func TestSDKCatalogAPIIsPublicAndRedacted(t *testing.T) {
	handler := testServerWithConfig(t, Config{SDKDownloads: sdkDownloadTestService(t)})
	request := httptest.NewRequest(http.MethodGet, "/api/sdk/catalog", nil)
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %q", response.Code, response.Body.String())
	}
	if got := response.Header().Get("Cache-Control"); got != "public, max-age=60" {
		t.Fatalf("Cache-Control = %q", got)
	}
	for _, forbidden := range []string{"object_key", "terms_object_key", "X-Amz-", "sdk/releases/"} {
		if strings.Contains(response.Body.String(), forbidden) {
			t.Fatalf("catalog leaked %q: %s", forbidden, response.Body.String())
		}
	}
	var catalog sdkdownloads.PublicCatalog
	if err := json.Unmarshal(response.Body.Bytes(), &catalog); err != nil {
		t.Fatal(err)
	}
	if catalog.Schema != sdkdownloads.PublicCatalogSchema || len(catalog.Packages) != 5 || catalog.CompleteBundle.Slug != "all" {
		t.Fatalf("unexpected catalog: %#v", catalog)
	}
}

func TestSDKCatalogAPIUnavailableAndMethod(t *testing.T) {
	handler := testServerWithConfig(t, Config{})
	request := httptest.NewRequest(http.MethodGet, "/api/sdk/catalog", nil)
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	if response.Code != http.StatusServiceUnavailable || response.Header().Get("Cache-Control") != "no-store" {
		t.Fatalf("unavailable response = %d %q", response.Code, response.Header().Get("Cache-Control"))
	}
	request = httptest.NewRequest(http.MethodPost, "/api/sdk/catalog", nil)
	response = httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	if response.Code != http.StatusMethodNotAllowed || response.Header().Get("Allow") != http.MethodGet {
		t.Fatalf("method response = %d Allow=%q", response.Code, response.Header().Get("Allow"))
	}
}

func TestSDKCatalogAndTermsArePublic(t *testing.T) {
	handler := testServerWithConfig(t, Config{SDKDownloads: sdkDownloadTestService(t)})
	for _, test := range []struct{ path, want string }{{"/manual/sdk", "0.1.0-rc.2"}, {"/zh-tw/manual/sdk", "下載套件"}, {"/legal/sdk-evaluation-terms", "Approved evaluation terms"}} {
		request := httptest.NewRequest(http.MethodGet, test.path, nil)
		response := httptest.NewRecorder()
		handler.ServeHTTP(response, request)
		body, _ := io.ReadAll(response.Result().Body)
		if response.Code != http.StatusOK || !strings.Contains(string(body), test.want) {
			t.Fatalf("GET %s = %d, body missing %q", test.path, response.Code, test.want)
		}
	}
}

func TestSDKDownloadRequiresCurrentTermsAndRedirects(t *testing.T) {
	handler := testServerWithConfig(t, Config{SDKDownloads: sdkDownloadTestService(t)})
	form := url.Values{"package": {"native"}, "version": {"0.1.0-rc.2"}, "terms_version": {"eval-v1"}, "accepted": {"true"}}
	request := httptest.NewRequest(http.MethodPost, "/manual/sdk/download", strings.NewReader(form.Encode()))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	if response.Code != http.StatusSeeOther {
		t.Fatalf("status = %d, want 303", response.Code)
	}
	if location := response.Header().Get("Location"); !strings.Contains(location, "X-Amz-Signature=test") {
		t.Fatalf("unexpected location %q", location)
	}
	if response.Header().Get("X-Request-ID") == "" {
		t.Fatal("missing request id")
	}

	form.Set("terms_version", "old")
	request = httptest.NewRequest(http.MethodPost, "/manual/sdk/download", strings.NewReader(form.Encode()))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	response = httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	if response.Code != http.StatusBadRequest {
		t.Fatalf("stale terms status = %d, want 400", response.Code)
	}
}

func TestSDKDownloadRequiresAcceptance(t *testing.T) {
	handler := testServerWithConfig(t, Config{SDKDownloads: sdkDownloadTestService(t)})
	request := httptest.NewRequest(http.MethodPost, "/manual/sdk/download", strings.NewReader("package=native"))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	if response.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", response.Code)
	}
}

func TestSDKDownloadUnavailableAndMalformedRequests(t *testing.T) {
	tests := []struct {
		name        string
		method      string
		body        string
		service     *sdkdownloads.Service
		wantStatus  int
		wantMessage string
	}{
		{name: "method", method: http.MethodGet, wantStatus: http.StatusMethodNotAllowed},
		{name: "unavailable", method: http.MethodPost, body: "accepted=true", wantStatus: http.StatusServiceUnavailable, wantMessage: "SDK downloads are unavailable"},
		{name: "malformed form", method: http.MethodPost, body: "%zz", service: sdkDownloadTestService(t), wantStatus: http.StatusBadRequest, wantMessage: "invalid download request"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			handler := testServerWithConfig(t, Config{SDKDownloads: test.service})
			request := httptest.NewRequest(test.method, "/manual/sdk/download", strings.NewReader(test.body))
			request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			response := httptest.NewRecorder()
			handler.ServeHTTP(response, request)
			if response.Code != test.wantStatus {
				t.Fatalf("status = %d, want %d", response.Code, test.wantStatus)
			}
			if test.wantMessage != "" && !strings.Contains(response.Body.String(), test.wantMessage) {
				t.Fatalf("body = %q, want %q", response.Body.String(), test.wantMessage)
			}
		})
	}
}

func TestSDKCatalogAndTermsReportUnavailableService(t *testing.T) {
	brokenService := sdkdownloads.NewService(webSDKStore{}, "", time.Minute)
	handler := testServerWithConfig(t, Config{SDKDownloads: brokenService})
	for _, path := range []string{"/manual/sdk", "/legal/sdk-evaluation-terms"} {
		request := httptest.NewRequest(http.MethodGet, path, nil)
		response := httptest.NewRecorder()
		handler.ServeHTTP(response, request)
		if response.Code != http.StatusOK || !strings.Contains(response.Body.String(), "temporarily unavailable") {
			t.Fatalf("GET %s = %d, body = %q", path, response.Code, response.Body.String())
		}
	}

	request := httptest.NewRequest(http.MethodPost, "/legal/sdk-evaluation-terms", nil)
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	if response.Code != http.StatusMethodNotAllowed {
		t.Fatalf("POST terms status = %d, want 405", response.Code)
	}
}

func TestSDKDownloadRateLimitAndMetrics(t *testing.T) {
	for _, sessionID := range []string{"", "-" + strings.Repeat("a", 23), "_" + strings.Repeat("a", 23)} {
		t.Run("session="+sessionID, func(t *testing.T) {
			testSDKDownloadRateLimitAndMetrics(t, sessionID)
		})
	}
}

func testSDKDownloadRateLimitAndMetrics(t *testing.T, sessionID string) {
	t.Helper()
	server := newTestServer(t, Config{
		TemplatesDir: "../../templates",
		StaticDir:    "../../static",
		ContentDir:   "../../content/docs",
		SDKDocsDir:   "testdata/sdk-docs",
		LeadStore:    &memoryLeadStore{},
		SDKDownloads: sdkDownloadTestService(t),
	})
	server.sdkDownloadLimit = newSubmissionRateLimiter(1, time.Hour)
	handler := server.Routes()
	form := url.Values{"package": {"native"}, "version": {"0.1.0-rc.2"}, "terms_version": {"eval-v1"}, "accepted": {"true"}}

	firstRequest := httptest.NewRequest(http.MethodPost, "/manual/sdk/download", strings.NewReader(form.Encode()))
	firstRequest.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	var sessionCookie *http.Cookie
	if sessionID != "" {
		sessionCookie = &http.Cookie{Name: sdkSessionCookie, Value: sessionID}
		firstRequest.AddCookie(sessionCookie)
	}
	firstResponse := httptest.NewRecorder()
	handler.ServeHTTP(firstResponse, firstRequest)
	if firstResponse.Code != http.StatusSeeOther {
		t.Fatalf("first status = %d, want 303", firstResponse.Code)
	}
	for _, cookie := range firstResponse.Result().Cookies() {
		if cookie.Name == sdkSessionCookie {
			if sessionID != "" {
				t.Fatal("valid SDK session cookie was replaced")
			}
			sessionCookie = cookie
			break
		}
	}
	if sessionCookie == nil {
		t.Fatal("missing SDK session cookie")
	}

	secondRequest := httptest.NewRequest(http.MethodPost, "/manual/sdk/download", strings.NewReader(form.Encode()))
	secondRequest.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	secondRequest.AddCookie(sessionCookie)
	secondResponse := httptest.NewRecorder()
	handler.ServeHTTP(secondResponse, secondRequest)
	if secondResponse.Code != http.StatusTooManyRequests {
		t.Fatalf("second status = %d, want 429", secondResponse.Code)
	}

	metricsRequest := httptest.NewRequest(http.MethodGet, "/metrics/prometheus", nil)
	metricsResponse := httptest.NewRecorder()
	handler.ServeHTTP(metricsResponse, metricsRequest)
	for _, want := range []string{
		`rtk_cloud_frontend_sdk_download_acceptances_total{version="0.1.0-rc.2",package="native"} 1`,
		`rtk_cloud_frontend_sdk_download_redirects_total{version="0.1.0-rc.2",package="native"} 1`,
		"rtk_cloud_frontend_sdk_download_errors_total 1",
	} {
		if !strings.Contains(metricsResponse.Body.String(), want) {
			t.Fatalf("metrics missing %q", want)
		}
	}
}

func TestFormatBytes(t *testing.T) {
	for _, test := range []struct {
		size int64
		want string
	}{{512, "512 B"}, {1024, "1.0 KB"}, {1024 * 1024, "1.0 MB"}, {1024 * 1024 * 1024, "1.0 GB"}} {
		if got := formatBytes(test.size); got != test.want {
			t.Fatalf("formatBytes(%d) = %q, want %q", test.size, got, test.want)
		}
	}
}
