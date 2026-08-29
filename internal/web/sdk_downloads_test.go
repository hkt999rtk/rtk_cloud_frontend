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
		packages = append(packages, sdkdownloads.Artifact{Slug: slug, Title: strings.ToUpper(slug), Description: slug + " SDK", Filename: slug + ".tgz", ObjectKey: "sdk/releases/0.1.0-rc.2/artifacts/" + slug + ".tgz", SHA256: strings.Repeat("a", 64), SizeBytes: 1024, ValidationStatus: "PASS"})
	}
	catalog := sdkdownloads.Catalog{Schema: sdkdownloads.CatalogSchema, Version: "0.1.0-rc.2", Distribution: "public-evaluation", TermsVersion: "eval-v1", TermsObjectKey: "sdk/releases/0.1.0-rc.2/EVALUATION_LICENSE.txt", Packages: packages, CompleteBundle: sdkdownloads.Artifact{Slug: "all", Title: "Complete", Description: "All SDKs", Filename: "all.tgz", ObjectKey: "sdk/releases/0.1.0-rc.2/artifacts/all.tgz", SHA256: strings.Repeat("b", 64), SizeBytes: 2048, ValidationStatus: "PASS"}}
	catalogJSON, _ := json.Marshal(catalog)
	latestJSON := []byte(`{"schema":"rtk-portal-sdk-latest/v1","version":"0.1.0-rc.2","catalog_object_key":"sdk/releases/0.1.0-rc.2/catalog.json","terms_version":"eval-v1"}`)
	return sdkdownloads.NewService(webSDKStore{objects: map[string][]byte{"sdk/latest.json": latestJSON, "sdk/releases/0.1.0-rc.2/catalog.json": catalogJSON, catalog.TermsObjectKey: []byte("Approved evaluation terms")}}, "", time.Minute)
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
