package docs

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParseContentPageRendersMarkdownImage(t *testing.T) {
	page, err := ParseContentPage([]byte(`---
title: "Docs"
subtitle: "Source-backed docs"
hero_image: "/static/assets/example.png"
hero_image_alt: "Example"
sections:
  - title: "One"
    icon: "document"
    body: "Body"
seo:
  meta_title: "Docs | Realtek Connect+"
  meta_description: "Docs description"
  social_image: "/static/assets/example.png"
---
Intro paragraph.

![Preview](/static/assets/example.png)
`))
	if err != nil {
		t.Fatalf("parse content page: %v", err)
	}
	if page.Title != "Docs" || page.Subtitle != "Source-backed docs" {
		t.Fatalf("page title/subtitle = %q/%q", page.Title, page.Subtitle)
	}
	if got := string(page.BodyHTML); !strings.Contains(got, `<img class="docs-content-image" src="/static/assets/example.png" alt="Preview">`) {
		t.Fatalf("body html missing rendered image: %s", got)
	}
	if len(page.Sections) != 1 || page.Sections[0].Title != "One" {
		t.Fatalf("sections = %#v", page.Sections)
	}
}

func TestContentSourceFallsBackMissingLocaleToEnglish(t *testing.T) {
	root := t.TempDir()
	enDir := filepath.Join(root, "en")
	if err := os.MkdirAll(enDir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(enDir, "docs.yaml"), []byte(`---
title: "English Docs"
subtitle: "Fallback"
---
Fallback body.
`), 0o644); err != nil {
		t.Fatalf("write content: %v", err)
	}

	pages, err := NewContentSource(root).Load()
	if err != nil {
		t.Fatalf("load content: %v", err)
	}
	if pages["zh-TW"].Title != "English Docs" {
		t.Fatalf("zh-TW fallback title = %q", pages["zh-TW"].Title)
	}
}
