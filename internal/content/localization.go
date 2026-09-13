package content

import (
	"crypto/sha256"
	"encoding/hex"
	"regexp"
	"sort"
	"strings"
)

var placeholderPattern = regexp.MustCompile(`{{\s*[-\w.]+\s*}}`)

// TranslationSource is the stable interchange format used by the workspace
// localization tool. English remains the only source language.
type TranslationSource struct {
	Key          string   `json:"key"`
	Source       string   `json:"source"`
	Context      string   `json:"context"`
	Placeholders []string `json:"placeholders"`
}

// TranslationArtifact is the Git-reviewed localized output consumed by Portal.
type TranslationArtifact struct {
	SchemaVersion int                                 `json:"schemaVersion"`
	Locale        string                              `json:"locale"`
	Entries       map[string]TranslationArtifactEntry `json:"entries"`
}

type TranslationArtifactEntry struct {
	SourceHash  string `json:"sourceHash"`
	Fingerprint string `json:"fingerprint"`
	Text        string `json:"text"`
	Status      string `json:"status"`
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
	return hash(payload)
}

// TranslationFingerprint matches the shared Node tool. Portal currently has no
// glossary, so its artifacts use GlossaryHash for an empty object.
func TranslationFingerprint(entry TranslationSource, locale string, policyVersion int, glossaryHash string) string {
	payload := `{"key":` + jsonString(entry.Key) + `,"sourceHash":` + jsonString(TranslationSourceHash(entry)) + `,"locale":` + jsonString(locale) + `,"policyVersion":` + integerJSON(policyVersion) + `,"glossaryHash":` + jsonString(glossaryHash) + `}`
	return hash(payload)
}

func emptyGlossaryHash() string { return hash(`{}`) }

func hash(payload string) string {
	digest := sha256.Sum256([]byte(payload))
	return hex.EncodeToString(digest[:])
}

func placeholders(value string) []string {
	matches := placeholderPattern.FindAllString(value, -1)
	seen := make(map[string]struct{}, len(matches))
	for _, match := range matches {
		seen[match] = struct{}{}
	}
	result := make([]string, 0, len(seen))
	for match := range seen {
		result = append(result, match)
	}
	sort.Strings(result)
	return result
}

func jsonString(value string) string {
	var builder strings.Builder
	builder.Grow(len(value) + 2)
	builder.WriteByte('"')
	for _, runeValue := range value {
		switch runeValue {
		case '\\':
			builder.WriteString(`\\`)
		case '"':
			builder.WriteString(`\"`)
		case '\b':
			builder.WriteString(`\b`)
		case '\f':
			builder.WriteString(`\f`)
		case '\n':
			builder.WriteString(`\n`)
		case '\r':
			builder.WriteString(`\r`)
		case '\t':
			builder.WriteString(`\t`)
		default:
			if runeValue < 0x20 {
				builder.WriteString(`\u00`)
				const hexadecimal = "0123456789abcdef"
				builder.WriteByte(hexadecimal[(runeValue>>4)&0x0f])
				builder.WriteByte(hexadecimal[runeValue&0x0f])
			} else {
				builder.WriteRune(runeValue)
			}
		}
	}
	builder.WriteByte('"')
	return builder.String()
}

func jsonStrings(values []string) string {
	if len(values) == 0 {
		return "[]"
	}
	normalized := append([]string(nil), values...)
	sort.Strings(normalized)
	parts := make([]string, 0, len(normalized))
	for _, value := range normalized {
		parts = append(parts, jsonString(value))
	}
	return "[" + strings.Join(parts, ",") + "]"
}

func integerJSON(value int) string {
	if value == 0 {
		return "0"
	}
	negative := value < 0
	if negative {
		value = -value
	}
	const digits = "0123456789"
	var reversed [20]byte
	index := len(reversed)
	for value > 0 {
		index--
		reversed[index] = digits[value%10]
		value /= 10
	}
	if negative {
		return "-" + string(reversed[index:])
	}
	return string(reversed[index:])
}
