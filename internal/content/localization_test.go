package content

import (
	"strings"
	"testing"
)

func TestTranslationManifestIsStableEnglishSource(t *testing.T) {
	entries := TranslationManifest()
	if len(entries) == 0 {
		t.Fatal("translation manifest is empty")
	}
	seen := map[string]bool{}
	for _, entry := range entries {
		if entry.Key == "" || entry.Source == "" || entry.Context == "" {
			t.Fatalf("invalid source entry: %#v", entry)
		}
		if seen[entry.Key] {
			t.Fatalf("duplicate source key %q", entry.Key)
		}
		seen[entry.Key] = true
		if TranslationSourceHash(entry) != TranslationSourceHash(entry) {
			t.Fatalf("unstable hash for %q", entry.Key)
		}
	}
	if !seen["text.nav.language"] {
		t.Fatal("manifest does not include shared Portal UI text")
	}
}

func TestTranslationHashMatchesSharedNodeFixtures(t *testing.T) {
	tests := []struct {
		name        string
		entry       TranslationSource
		sourceHash  string
		fingerprint string
	}{
		{
			name:        "sorted placeholders",
			entry:       TranslationSource{Key: "sorted", Source: "Hello {{a}} {{z}}", Context: "test", Placeholders: []string{"{{z}}", "{{a}}"}},
			sourceHash:  "2a137a20c17eeb9abd530b214598f3eb25ff6b266eec309f3a419eda975de85c",
			fingerprint: "a36381abfdcfb2631692f6591c3872843fbdc47eb10a1980bd531850dc219766",
		},
		{
			name:        "JSON control and Unicode characters",
			entry:       TranslationSource{Key: "control", Source: "Line\b\f<&>\u2028\u2029", Context: "test", Placeholders: []string{}},
			sourceHash:  "70960de109c8cc951fc3ba51fce150656eda33da7f233c44a60b4f9fba0ebbda",
			fingerprint: "d15858bbe37eb46d82ab8aacf5ed2acc35e356864e32ff0c73016eb10509ab02",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := TranslationSourceHash(tc.entry); got != tc.sourceHash {
				t.Fatalf("TranslationSourceHash() = %s, want %s", got, tc.sourceHash)
			}
			if got := TranslationFingerprint(tc.entry, "zh-TW", 1, emptyGlossaryHash()); got != tc.fingerprint {
				t.Fatalf("TranslationFingerprint() = %s, want %s", got, tc.fingerprint)
			}
		})
	}
	if got := placeholders("Hello {{name}} {{name}} {{date}}"); len(got) != 2 || got[0] != "{{date}}" || got[1] != "{{name}}" {
		t.Fatalf("placeholders() = %#v", got)
	}
}

func TestPortalCatalogUsesApprovedTranslationArtifacts(t *testing.T) {
	manifest := TranslationManifest()
	for _, locale := range []string{"zh-TW", "zh-CN"} {
		artifact := approvedPortalArtifact(locale)
		if len(artifact.Entries) != len(manifest) {
			t.Fatalf("%s artifact has %d entries, want %d", locale, len(artifact.Entries), len(manifest))
		}
		catalog := CatalogFor(localeByCode(t, locale))
		for _, entry := range manifest {
			translation := artifact.Entries[entry.Key]
			switch {
			case strings.HasPrefix(entry.Key, "text."):
				key := strings.TrimPrefix(entry.Key, "text.")
				if got := catalog.T(key); got != translation.Text {
					t.Fatalf("%s %s runtime text = %q, want artifact %q", locale, entry.Key, got, translation.Text)
				}
			case strings.HasPrefix(entry.Key, "page."):
				pageKey := strings.TrimPrefix(entry.Key, "page.")
				switch {
				case strings.HasSuffix(pageKey, ".title"):
					if got := catalog.Page(strings.TrimSuffix(pageKey, ".title")).Title; got != translation.Text {
						t.Fatalf("%s %s runtime title = %q, want artifact %q", locale, entry.Key, got, translation.Text)
					}
				case strings.HasSuffix(pageKey, ".description"):
					if got := catalog.Page(strings.TrimSuffix(pageKey, ".description")).Description; got != translation.Text {
						t.Fatalf("%s %s runtime description = %q, want artifact %q", locale, entry.Key, got, translation.Text)
					}
				default:
					t.Fatalf("unexpected Portal page key %q", entry.Key)
				}
			}
		}
	}
}

func TestPortalArtifactRejectsMismatchedPlaceholders(t *testing.T) {
	original := approvedPortalArtifact("zh-TW")
	artifact := original
	artifact.Entries = make(map[string]TranslationArtifactEntry, len(original.Entries))
	for key, entry := range original.Entries {
		artifact.Entries[key] = entry
	}
	entry := artifact.Entries["text.nav.docs"]
	entry.Text += " {{name}}"
	artifact.Entries["text.nav.docs"] = entry
	if err := validatePortalArtifact("zh-TW", artifact); err == nil {
		t.Fatal("validatePortalArtifact accepted an unexpected placeholder")
	}
}
