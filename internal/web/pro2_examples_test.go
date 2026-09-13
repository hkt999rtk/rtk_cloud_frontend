package web

import (
	"encoding/json"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"realtek-connect/internal/sdkdownloads"
	"strings"
	"testing"
	"time"
)

func TestPRO2CatalogAndTermsDownload(t *testing.T) {
	c := sdkdownloads.ExamplesCatalog{Schema: "rtk-pro2-examples/v1", Version: "v1", SourceCommit: strings.Repeat("a", 40), CreatedAt: "2026-09-07T00:00:00Z", TermsVersion: "eval-v1", Terms: "Evaluation only", TestOnly: true}
	for _, id := range []string{"source", "mqtt", "webrtc_test_video", "webrtc_camera"} {
		kind := "firmware"
		if id == "source" {
			kind = "source"
		}
		c.Artifacts = append(c.Artifacts, sdkdownloads.ExampleArtifact{ID: id, Filename: id + ".bin", SHA256: strings.Repeat("b", 64), SizeBytes: 10, Kind: kind})
		if id != "source" {
			c.Artifacts = append(c.Artifacts, sdkdownloads.ExampleArtifact{ID: id + "-sha", Filename: id + ".sha256", SHA256: strings.Repeat("c", 64), SizeBytes: 80, Kind: "checksum"})
			c.Examples = append(c.Examples, sdkdownloads.CloudExample{ID: id, Title: id, ImageType: "full-flash-ntz", FirmwareID: id, ChecksumID: id + "-sha"})
		}
	}
	b, _ := json.Marshal(c)
	store := webSDKStore{objects: map[string][]byte{"pro2-examples/latest.json": []byte(`{"version":"v1"}`), "pro2-examples/releases/v1/manifest.json": b}}
	h := testServerWithConfig(t, Config{SDKDownloads: sdkdownloads.NewService(store, "", time.Minute)})
	r := httptest.NewRecorder()
	h.ServeHTTP(r, httptest.NewRequest("GET", "/api/pro2-examples/catalog", nil))
	if r.Code != 200 || strings.Contains(r.Body.String(), "object_key") {
		t.Fatal(r.Code, r.Body.String())
	}
	for _, tc := range []struct {
		accepted, id, terms string
		status              int
	}{{"true", "mqtt", "eval-v1", 200}, {"false", "mqtt", "eval-v1", 400}, {"true", "missing", "eval-v1", 400}, {"true", "mqtt", "old", 400}} {
		form := url.Values{"accepted": {tc.accepted}, "artifact": {tc.id}, "version": {"v1"}, "terms_version": {tc.terms}}
		q := httptest.NewRequest("POST", "/api/pro2-examples/download", strings.NewReader(form.Encode()))
		q.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		r := httptest.NewRecorder()
		h.ServeHTTP(r, q)
		if r.Code != tc.status {
			t.Fatal(r.Code, r.Body.String())
		}
		if r.Code == http.StatusOK && !strings.Contains(r.Body.String(), "expires_at") {
			t.Fatal("no expiration")
		}
	}
	// A browser must retain the API acceptance session between artifact downloads.
	server := httptest.NewServer(h)
	defer server.Close()
	jar, _ := cookiejar.New(nil)
	client := &http.Client{Jar: jar}
	var previous string
	for i := 0; i < 2; i++ {
		response, err := client.PostForm(server.URL+"/api/pro2-examples/download", url.Values{"accepted": {"true"}, "artifact": {"mqtt"}, "version": {"v1"}, "terms_version": {"eval-v1"}})
		if err != nil {
			t.Fatal(err)
		}
		response.Body.Close()
		if response.StatusCode != 200 {
			t.Fatal(response.StatusCode)
		}
		endpoint, _ := url.Parse(server.URL + "/api/pro2-examples/download")
		cookies := jar.Cookies(endpoint)
		if len(cookies) != 1 || (previous != "" && cookies[0].Value != previous) {
			t.Fatal("acceptance session not retained")
		}
		previous = cookies[0].Value
		if i == 1 && len(response.Cookies()) != 0 {
			t.Fatal("existing session replaced")
		}
	}
	delete(store.objects, "pro2-examples/releases/v1/manifest.json")
	r = httptest.NewRecorder()
	h.ServeHTTP(r, httptest.NewRequest("GET", "/api/pro2-examples/catalog", nil))
	if r.Code != 503 {
		t.Fatal(r.Code)
	}
}
