package manual

import (
	"bytes"
	"html"
	"html/template"
	"net/url"
	"regexp"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
)

var (
	markdownRenderer = goldmark.New(
		goldmark.WithExtensions(extension.GFM),
	)
	urlAttributePattern = regexp.MustCompile(`\b(href|src)="([^"]*)"`)
)

func RenderMarkdown(input []byte) template.HTML {
	var buf bytes.Buffer
	if err := markdownRenderer.Convert(input, &buf); err != nil {
		return template.HTML(html.EscapeString(string(input)))
	}
	return template.HTML(sanitizeRenderedHTML(buf.String()))
}

func sanitizeRenderedHTML(input string) string {
	return urlAttributePattern.ReplaceAllStringFunc(input, func(match string) string {
		parts := urlAttributePattern.FindStringSubmatch(match)
		if len(parts) != 3 {
			return match
		}
		safe := sanitizeURL(parts[2])
		return parts[1] + `="` + html.EscapeString(safe) + `"`
	})
}

func sanitizeURL(raw string) string {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return ""
	}

	if strings.HasPrefix(trimmed, "#") ||
		strings.HasPrefix(trimmed, "/") ||
		strings.HasPrefix(trimmed, "./") ||
		strings.HasPrefix(trimmed, "../") {
		return trimmed
	}

	parsed, err := url.Parse(trimmed)
	if err != nil {
		return "#"
	}
	switch strings.ToLower(parsed.Scheme) {
	case "http", "https", "mailto", "tel":
		return trimmed
	default:
		return "#"
	}
}
