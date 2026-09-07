package sdkdownloads

import (
	"context"
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func exampleFixture() ExamplesCatalog {
	c := ExamplesCatalog{Schema: "rtk-pro2-examples/v1", Version: "v1", SourceCommit: strings.Repeat("a", 40), CreatedAt: "2026-09-07T00:00:00Z", TermsVersion: "eval-v1", Terms: "Evaluation only", TestOnly: true}
	for _, id := range []string{"source", "mqtt", "webrtc_test_video", "webrtc_camera"} {
		kind := "firmware"
		if id == "source" {
			kind = "source"
		}
		c.Artifacts = append(c.Artifacts, ExampleArtifact{ID: id, Filename: id + ".bin", Kind: kind, SHA256: strings.Repeat("b", 64), SizeBytes: 10})
		if id != "source" {
			c.Artifacts = append(c.Artifacts, ExampleArtifact{ID: id + "-sha", Filename: id + ".sha256", Kind: "checksum", SHA256: strings.Repeat("c", 64), SizeBytes: 80})
			c.Examples = append(c.Examples, CloudExample{ID: id, Title: id, ImageType: "full-flash-ntz", FirmwareID: id, ChecksumID: id + "-sha"})
		}
	}
	return c
}
func TestExamplesVersionedDownloads(t *testing.T) {
	ctx := context.Background()
	c := exampleFixture()
	b, _ := json.Marshal(c)
	store := fakeStore{objects: map[string][]byte{ExamplesPrefix + "latest.json": []byte(`{"version":"v1"}`), ExamplesPrefix + "releases/v1/manifest.json": b}}
	s := NewService(store, "", time.Minute)
	got, e := s.Examples(ctx, "")
	if e != nil || got.Version != "v1" {
		t.Fatal(got, e)
	}
	a, u, e := s.ExampleDownloadURL(ctx, "v1", "mqtt", "eval-v1", time.Minute)
	if e != nil || a.ID != "mqtt" || !strings.Contains(u, "pro2-examples/releases/v1/mqtt.bin") {
		t.Fatal(a, u, e)
	}
	for _, v := range [][3]string{{"missing", "mqtt", "eval-v1"}, {"../v1", "mqtt", "eval-v1"}, {"v1", "unknown", "eval-v1"}, {"v1", "mqtt", "old"}} {
		if _, _, e := s.ExampleDownloadURL(ctx, v[0], v[1], v[2], time.Minute); e == nil {
			t.Fatal(v)
		}
	}
	// A saved release remains usable after latest moves.
	store.objects[ExamplesPrefix+"latest.json"] = []byte(`{"version":"v2"}`)
	if _, _, e := s.ExampleDownloadURL(ctx, "v1", "mqtt", "eval-v1", time.Minute); e != nil {
		t.Fatal(e)
	}
}
func TestExamplesRejectUnsafeCatalog(t *testing.T) {
	for _, change := range []func(*ExamplesCatalog){func(c *ExamplesCatalog) { c.TestOnly = false }, func(c *ExamplesCatalog) { c.Artifacts[0].Filename = "../secret" }, func(c *ExamplesCatalog) { c.Artifacts[1].SHA256 = "bad" }, func(c *ExamplesCatalog) { c.Examples[0].FlashOffset = 4096 }, func(c *ExamplesCatalog) { c.Examples[0].FirmwareID = "missing" }} {
		c := exampleFixture()
		change(&c)
		b, _ := json.Marshal(c)
		s := NewService(fakeStore{objects: map[string][]byte{ExamplesPrefix + "releases/v1/manifest.json": b}}, "", time.Minute)
		if _, e := s.Examples(context.Background(), "v1"); e == nil {
			t.Fatal("accepted invalid catalog")
		}
	}
}
