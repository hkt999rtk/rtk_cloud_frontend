package manual

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"realtek-connect/internal/content"
)

func TestLoaderFallsBackToEnglishForMissingLocales(t *testing.T) {
	root := t.TempDir()
	writeManualIndex(t, root, "en", "User Manual", "English description")
	writeManualPage(t, root, "en", "getting-started", "Getting Started", "Get started quickly", "Hello manual.")

	loader := NewLoader(root)
	index, ok := loader.Index(content.Locale{Code: "zh-TW"})
	if !ok {
		t.Fatalf("zh-TW index missing")
	}
	if index.Title != "User Manual" {
		t.Fatalf("zh-TW title = %q, want English fallback", index.Title)
	}

	page, ok := loader.Page(content.Locale{Code: "zh-CN"}, "getting-started")
	if !ok {
		t.Fatalf("zh-CN page missing")
	}
	if page.Title != "Getting Started" {
		t.Fatalf("zh-CN page title = %q, want English fallback", page.Title)
	}
}

func TestRenderMarkdownSanitizesUnsafeLinks(t *testing.T) {
	html := RenderMarkdown([]byte(`[safe](/manual)
[unsafe](javascript:alert(1))

![image](https://example.com/image.png)
`))
	body := string(html)
	if !strings.Contains(body, `href="/manual"`) {
		t.Fatalf("safe link not rendered: %s", body)
	}
	if strings.Contains(body, `javascript:alert(1)`) {
		t.Fatalf("unsafe url was not sanitized: %s", body)
	}
	if !strings.Contains(body, `src="https://example.com/image.png"`) {
		t.Fatalf("image url missing: %s", body)
	}
}

func TestReloadKeepsLoadedSectionsAndPages(t *testing.T) {
	root := t.TempDir()
	writeManualIndex(t, root, "en", "User Manual", "Manual description")
	writeManualPage(t, root, "en", "getting-started", "Getting Started", "Get started quickly", "Manual body.")

	loader := NewLoader(root)
	index, ok := loader.Index(content.DefaultLocale())
	if !ok {
		t.Fatalf("index missing")
	}
	if len(index.Sections) != 1 {
		t.Fatalf("sections = %d, want 1", len(index.Sections))
	}

	page, ok := loader.Page(content.DefaultLocale(), "getting-started")
	if !ok {
		t.Fatalf("page missing")
	}
	if !strings.Contains(string(page.BodyHTML), "Manual body.") {
		t.Fatalf("page body missing: %s", string(page.BodyHTML))
	}

	writeManualIndex(t, root, "en", "Updated Manual", "Manual description")
	if err := loader.Reload(); err != nil {
		t.Fatalf("reload: %v", err)
	}
	index, ok = loader.Index(content.DefaultLocale())
	if !ok || index.Title != "Updated Manual" {
		t.Fatalf("reloaded index = %#v", index)
	}
}

func writeManualIndex(t *testing.T, root, locale, title, description string) {
	t.Helper()
	body := "---\n" +
		"title: \"" + title + "\"\n" +
		"description: \"" + description + "\"\n" +
		"sections:\n" +
		"  - slug: getting-started\n" +
		"    title: \"Getting Started\"\n" +
		"    summary: \"Start here\"\n" +
		"---\n"
	if err := os.WriteFile(filepath.Join(root, "index."+locale+".yaml"), []byte(body), 0o644); err != nil {
		t.Fatalf("write manual index: %v", err)
	}
}

func writeManualPage(t *testing.T, root, locale, slug, title, description, body string) {
	t.Helper()
	content := "---\n" +
		"title: \"" + title + "\"\n" +
		"description: \"" + description + "\"\n" +
		"---\n" +
		body + "\n"
	if err := os.WriteFile(filepath.Join(root, slug+"."+locale+".md"), []byte(content), 0o644); err != nil {
		t.Fatalf("write manual page: %v", err)
	}
}
