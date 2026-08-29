package sdkdownloads

import (
	"context"
	"encoding/json"
	"strings"
	"testing"
	"time"
)

type fakeStore struct{ objects map[string][]byte }

func (f fakeStore) Get(_ context.Context, key string) ([]byte, error) { return f.objects[key], nil }
func (f fakeStore) PresignGet(key string, _ time.Duration) (string, error) {
	return "https://objects.example/" + key + "?signed=1", nil
}

func testCatalog() Catalog {
	packages := []Artifact{}
	for _, slug := range []string{"native", "android", "javascript", "ios", "freertos-pro2"} {
		packages = append(packages, Artifact{Slug: slug, Title: slug, Filename: slug + ".tgz", ObjectKey: "sdk/releases/0.1.0-rc.2/artifacts/" + slug + ".tgz", SHA256: strings.Repeat("a", 64), SizeBytes: 10, ValidationStatus: "PASS"})
	}
	return Catalog{Schema: CatalogSchema, Version: "0.1.0-rc.2", Distribution: "public-evaluation", TermsVersion: "eval-v1", TermsObjectKey: "sdk/releases/0.1.0-rc.2/EVALUATION_LICENSE.txt", Packages: packages, CompleteBundle: Artifact{Slug: "all", Title: "All", Filename: "all.tgz", ObjectKey: "sdk/releases/0.1.0-rc.2/artifacts/all.tgz", SHA256: strings.Repeat("b", 64), SizeBytes: 50, ValidationStatus: "PASS"}}
}

func TestServiceLoadsAndSignsAllowlistedArtifact(t *testing.T) {
	catalog := testCatalog()
	catalogJSON, _ := json.Marshal(catalog)
	latestJSON, _ := json.Marshal(latestPointer{Schema: LatestSchema, Version: catalog.Version, CatalogObjectKey: "sdk/releases/0.1.0-rc.2/catalog.json", TermsVersion: catalog.TermsVersion})
	service := NewService(fakeStore{objects: map[string][]byte{"sdk/latest.json": latestJSON, "sdk/releases/0.1.0-rc.2/catalog.json": catalogJSON, catalog.TermsObjectKey: []byte("terms")}}, "", time.Minute)
	artifact, downloadURL, err := service.DownloadURL(context.Background(), catalog.Version, "native", catalog.TermsVersion, 10*time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	if artifact.Slug != "native" || !strings.Contains(downloadURL, "signed=1") {
		t.Fatalf("unexpected download: %#v %s", artifact, downloadURL)
	}
}

func TestServiceRejectsStaleTermsAndUnknownPackage(t *testing.T) {
	catalog := testCatalog()
	catalogJSON, _ := json.Marshal(catalog)
	latestJSON, _ := json.Marshal(latestPointer{Schema: LatestSchema, Version: catalog.Version, CatalogObjectKey: "catalog.json", TermsVersion: catalog.TermsVersion})
	service := NewService(fakeStore{objects: map[string][]byte{"sdk/latest.json": latestJSON, "catalog.json": catalogJSON, catalog.TermsObjectKey: []byte("terms")}}, "", time.Minute)
	if _, _, err := service.DownloadURL(context.Background(), catalog.Version, "native", "old", time.Minute); err == nil {
		t.Fatal("expected stale terms rejection")
	}
	if _, _, err := service.DownloadURL(context.Background(), catalog.Version, "go", catalog.TermsVersion, time.Minute); err == nil {
		t.Fatal("expected package rejection")
	}
}

func TestPresignedURLIsBoundedAndSigned(t *testing.T) {
	store, err := NewS3Store(StoreConfig{Bucket: "bucket", Endpoint: "https://objects.example", Region: "us-sea", AccessKey: "access", SecretKey: "secret"})
	if err != nil {
		t.Fatal(err)
	}
	store.now = func() time.Time { return time.Date(2026, 8, 29, 0, 0, 0, 0, time.UTC) }
	u, err := store.PresignGet("sdk/releases/v/artifact.tgz", 10*time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"X-Amz-Expires=600", "X-Amz-Signature=", "X-Amz-Credential="} {
		if !strings.Contains(u, want) {
			t.Fatalf("URL missing %s: %s", want, u)
		}
	}
}
