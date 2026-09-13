package content

import (
	"crypto/sha256"
	"encoding/hex"
	"sort"
)

// TranslationSource is the stable interchange format used by the workspace
// localization tool. English remains the only source language.
type TranslationSource struct {
	Key          string   `json:"key"`
	Source       string   `json:"source"`
	Context      string   `json:"context"`
	Placeholders []string `json:"placeholders"`
}

// TranslationManifest exports Portal UI and SEO strings. Feature and document
// bodies keep their existing authored review workflow and are intentionally not
// sent to machine translation by this first integration.
func TranslationManifest() []TranslationSource {
	entries := make([]TranslationSource, 0, len(enText())+len(enPages())*2)
	for key, value := range enText() {
		entries = append(entries, TranslationSource{Key: "text." + key, Source: value, Context: "Portal shared interface text", Placeholders: placeholders(value)})
	}
	for key, value := range enPages() {
		entries = append(entries,
			TranslationSource{Key: "page." + key + ".title", Source: value.Title, Context: "Portal page title", Placeholders: placeholders(value.Title)},
			TranslationSource{Key: "page." + key + ".description", Source: value.Description, Context: "Portal page SEO description", Placeholders: placeholders(value.Description)},
		)
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].Key < entries[j].Key })
	return entries
}

// TranslationSourceHash matches the canonical field order used by
// tools/localization/localization.mjs.
func TranslationSourceHash(entry TranslationSource) string {
	payload := `{"key":` + jsonString(entry.Key) + `,"source":` + jsonString(entry.Source) + `,"context":` + jsonString(entry.Context) + `,"placeholders":` + jsonStrings(entry.Placeholders) + `}`
	digest := sha256.Sum256([]byte(payload))
	return hex.EncodeToString(digest[:])
}

func placeholders(value string) []string {
	// Portal source text currently has no interpolation tokens. Keeping this
	// helper makes the manifest compatible with future {{name}} strings.
	return nil
}

func jsonString(value string) string {
	// The standard encoder is deliberately used here so escaping stays aligned
	// with JSON artifacts consumed by the Node tool.
	encoded := make([]byte, 0, len(value)+2)
	encoded = append(encoded, '"')
	for _, r := range value {
		switch r {
		case '\\':
			encoded = append(encoded, '\\', '\\')
		case '"':
			encoded = append(encoded, '\\', '"')
		case '\n':
			encoded = append(encoded, '\\', 'n')
		case '\r':
			encoded = append(encoded, '\\', 'r')
		case '\t':
			encoded = append(encoded, '\\', 't')
		default:
			encoded = append(encoded, string(r)...)
		}
	}
	return string(append(encoded, '"'))
}

func jsonStrings(values []string) string {
	if len(values) == 0 {
		return "[]"
	}
	out := "["
	for index, value := range values {
		if index > 0 {
			out += ","
		}
		out += jsonString(value)
	}
	return out + "]"
}
