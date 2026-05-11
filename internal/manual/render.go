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
	return RenderMarkdownWithPrefix(input, "")
}

func RenderMarkdownWithPrefix(input []byte, localePrefix string) template.HTML {
	var buf bytes.Buffer
	if err := markdownRenderer.Convert(input, &buf); err != nil {
		return template.HTML(html.EscapeString(string(input)))
	}
	return template.HTML(sanitizeRenderedHTML(buf.String(), localePrefix))
}

func LocalizeRenderedHTML(input template.HTML, localePrefix string) template.HTML {
	if localePrefix == "" {
		return input
	}
	return template.HTML(sanitizeRenderedHTML(string(input), localePrefix))
}

func sanitizeRenderedHTML(input, localePrefix string) string {
	return urlAttributePattern.ReplaceAllStringFunc(input, func(match string) string {
		parts := urlAttributePattern.FindStringSubmatch(match)
		if len(parts) != 3 {
			return match
		}
		safe := sanitizeURL(parts[1], parts[2], localePrefix)
		return parts[1] + `="` + html.EscapeString(safe) + `"`
	})
}

func sanitizeURL(attr, raw, localePrefix string) string {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return ""
	}

	if strings.HasPrefix(trimmed, "#") ||
		strings.HasPrefix(trimmed, "./") ||
		strings.HasPrefix(trimmed, "../") {
		return trimmed
	}
	if strings.HasPrefix(trimmed, "/") {
		return localizeRootRelativeHref(attr, trimmed, localePrefix)
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

func localizeRootRelativeHref(attr, raw, localePrefix string) string {
	if strings.ToLower(attr) != "href" || localePrefix == "" || !strings.HasPrefix(raw, "/") {
		return raw
	}
	if strings.HasPrefix(raw, "//") ||
		strings.HasPrefix(raw, localePrefix+"/") ||
		raw == localePrefix ||
		strings.HasPrefix(raw, "/zh-tw/") ||
		strings.HasPrefix(raw, "/zh-cn/") ||
		strings.HasPrefix(raw, "/static/") ||
		strings.HasPrefix(raw, "/content-assets/") ||
		strings.HasPrefix(raw, "/admin/") ||
		strings.HasPrefix(raw, "/api/") ||
		raw == "/healthz" ||
		raw == "/robots.txt" ||
		raw == "/sitemap.xml" {
		return raw
	}
	return localePrefix + raw
}
