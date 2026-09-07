package sdkdownloads

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"path"
	"strings"
	"time"
)

const ExamplesPrefix = "pro2-examples/"

type ExampleArtifact struct {
	ID        string `json:"id"`
	Filename  string `json:"filename"`
	SHA256    string `json:"sha256"`
	SizeBytes int64  `json:"size_bytes"`
	Kind      string `json:"kind"`
}
type CloudExample struct {
	ID          string            `json:"id"`
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Board       string            `json:"board"`
	Sensor      string            `json:"sensor"`
	FirmwareID  string            `json:"firmware_id"`
	ChecksumID  string            `json:"checksum_id"`
	FlashOffset uint32            `json:"flash_offset"`
	ImageType   string            `json:"image_type"`
	Validation  map[string]string `json:"validation"`
}
type ExamplesCatalog struct {
	Schema       string            `json:"schema"`
	Version      string            `json:"version"`
	SourceCommit string            `json:"source_commit"`
	CreatedAt    string            `json:"created_at"`
	TermsVersion string            `json:"terms_version"`
	Terms        string            `json:"terms"`
	TestOnly     bool              `json:"test_only"`
	Dependencies map[string]string `json:"dependencies"`
	Examples     []CloudExample    `json:"examples"`
	Artifacts    []ExampleArtifact `json:"artifacts"`
}

// Example objects live separately from the existing SDK release train. No object
// keys or credentials are returned; keys are derived from validated identities.
func (s *Service) Examples(ctx context.Context, version string) (ExamplesCatalog, error) {
	var c ExamplesCatalog
	if s == nil || s.store == nil {
		return c, errors.New("examples unavailable")
	}
	if version == "" {
		b, e := s.store.Get(ctx, s.examplesPrefix+"latest.json")
		if e != nil {
			return c, e
		}
		var p struct {
			Version string `json:"version"`
		}
		if e = json.Unmarshal(b, &p); e != nil {
			return c, e
		}
		version = p.Version
	}
	if !versionPattern.MatchString(version) {
		return c, errors.New("invalid version")
	}
	b, e := s.store.Get(ctx, s.examplesPrefix+"releases/"+version+"/manifest.json")
	if e != nil {
		return c, e
	}
	if e = json.Unmarshal(b, &c); e != nil {
		return c, e
	}
	if c.Schema != "rtk-pro2-examples/v1" || c.Version != version || !versionPattern.MatchString(c.TermsVersion) || strings.TrimSpace(c.Terms) == "" || !c.TestOnly || len(c.Examples) != 3 || len(c.SourceCommit) != 40 {
		return c, errors.New("invalid examples release")
	}
	if _, e = time.Parse(time.RFC3339, c.CreatedAt); e != nil {
		return c, e
	}
	artifacts := map[string]ExampleArtifact{}
	for _, a := range c.Artifacts {
		if !versionPattern.MatchString(a.ID) || a.Filename == "" || path.Base(a.Filename) != a.Filename || !versionPattern.MatchString(a.Filename) || a.SizeBytes <= 0 || a.SizeBytes > 64<<20 || !sha256Pattern.MatchString(a.SHA256) {
			return c, errors.New("invalid example artifact")
		}
		if _, ok := artifacts[a.ID]; ok {
			return c, errors.New("duplicate artifact")
		}
		artifacts[a.ID] = a
	}
	if artifacts["source"].Kind != "source" {
		return c, errors.New("missing source")
	}
	seen := map[string]bool{}
	for _, x := range c.Examples {
		if (x.ID != "mqtt" && x.ID != "webrtc_test_video" && x.ID != "webrtc_camera") || seen[x.ID] || x.Title == "" || x.ImageType != "full-flash-ntz" || x.FlashOffset != 0 || artifacts[x.FirmwareID].Kind != "firmware" || artifacts[x.ChecksumID].Kind != "checksum" {
			return c, errors.New("invalid example")
		}
		seen[x.ID] = true
	}
	return c, nil
}

func (s *Service) ExampleDownloadURL(ctx context.Context, version, id, terms string, ttl time.Duration) (ExampleArtifact, string, error) {
	c, e := s.Examples(ctx, version)
	if e != nil {
		return ExampleArtifact{}, "", e
	}
	if version != c.Version || terms != c.TermsVersion {
		return ExampleArtifact{}, "", errors.New("version or terms mismatch")
	}
	for _, a := range c.Artifacts {
		if a.ID == id {
			if ttl <= 0 || ttl > 15*time.Minute {
				ttl = 10 * time.Minute
			}
			u, e := s.store.PresignGet(fmt.Sprintf("%sreleases/%s/%s", s.examplesPrefix, c.Version, a.Filename), ttl)
			return a, u, e
		}
	}
	return ExampleArtifact{}, "", errors.New("unknown artifact")
}

// SetExamplesPrefix is called once at startup to isolate environment releases.
func (s *Service) SetExamplesPrefix(prefix string) error {
	if prefix == "" {
		return nil
	}
	if !strings.HasSuffix(prefix, "/") || strings.HasPrefix(prefix, "/") || strings.Contains(prefix, "..") || path.Clean(prefix)+"/" != prefix {
		return errors.New("invalid examples object prefix")
	}
	s.examplesPrefix = prefix
	return nil
}
