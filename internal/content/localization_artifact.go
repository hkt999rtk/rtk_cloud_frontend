package content

import (
	"embed"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
)

//go:embed localization/translations/zh-TW.json localization/translations/zh-CN.json
var portalTranslationFiles embed.FS

var portalArtifacts = map[string]*struct {
	once     sync.Once
	artifact TranslationArtifact
	err      error
}{
	"zh-TW": {},
	"zh-CN": {},
}

func approvedPortalArtifact(locale string) TranslationArtifact {
	state, ok := portalArtifacts[locale]
	if !ok {
		panic(fmt.Sprintf("unsupported Portal translation locale %q", locale))
	}
	state.once.Do(func() {
		payload, err := portalTranslationFiles.ReadFile("localization/translations/" + locale + ".json")
		if err != nil {
			state.err = err
			return
		}
		if err := json.Unmarshal(payload, &state.artifact); err != nil {
			state.err = fmt.Errorf("parse %s translation artifact: %w", locale, err)
			return
		}
		state.err = validatePortalArtifact(locale, state.artifact)
	})
	if state.err != nil {
		panic(state.err)
	}
	return state.artifact
}

func validatePortalArtifact(locale string, artifact TranslationArtifact) error {
	if artifact.SchemaVersion != 1 || artifact.Locale != locale {
		return fmt.Errorf("invalid %s translation artifact header", locale)
	}
	manifest := TranslationManifest()
	if len(artifact.Entries) != len(manifest) {
		return fmt.Errorf("%s translation artifact has %d entries, want %d", locale, len(artifact.Entries), len(manifest))
	}
	for _, entry := range manifest {
		translation, ok := artifact.Entries[entry.Key]
		if !ok {
			return fmt.Errorf("%s translation artifact is missing %s", locale, entry.Key)
		}
		if translation.Status != "approved" || strings.TrimSpace(translation.Text) == "" {
			return fmt.Errorf("%s translation artifact has unapproved or blank %s", locale, entry.Key)
		}
		if !sameStrings(placeholders(translation.Text), placeholders(entry.Source)) {
			return fmt.Errorf("%s translation artifact has mismatched placeholders for %s", locale, entry.Key)
		}
		if translation.SourceHash != TranslationSourceHash(entry) {
			return fmt.Errorf("%s translation artifact has stale source hash for %s", locale, entry.Key)
		}
		if translation.Fingerprint != TranslationFingerprint(entry, locale, 1, emptyGlossaryHash()) {
			return fmt.Errorf("%s translation artifact has stale fingerprint for %s", locale, entry.Key)
		}
	}
	return nil
}

func sameStrings(left, right []string) bool {
	if len(left) != len(right) {
		return false
	}
	for index := range left {
		if left[index] != right[index] {
			return false
		}
	}
	return true
}

func portalApprovedText(locale string) map[string]string {
	artifact := approvedPortalArtifact(locale)
	result := make(map[string]string)
	for key, entry := range artifact.Entries {
		if strings.HasPrefix(key, "text.") {
			result[strings.TrimPrefix(key, "text.")] = entry.Text
		}
	}
	return result
}

func portalApprovedPages(locale string) map[string]PageMeta {
	artifact := approvedPortalArtifact(locale)
	result := make(map[string]PageMeta)
	for key, entry := range artifact.Entries {
		if !strings.HasPrefix(key, "page.") {
			continue
		}
		pageKey := strings.TrimPrefix(key, "page.")
		switch {
		case strings.HasSuffix(pageKey, ".title"):
			page := result[strings.TrimSuffix(pageKey, ".title")]
			page.Title = entry.Text
			result[strings.TrimSuffix(pageKey, ".title")] = page
		case strings.HasSuffix(pageKey, ".description"):
			page := result[strings.TrimSuffix(pageKey, ".description")]
			page.Description = entry.Text
			result[strings.TrimSuffix(pageKey, ".description")] = page
		}
	}
	return result
}
