package content

import "testing"

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
