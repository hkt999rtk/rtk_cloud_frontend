package search

import (
	"context"
	"fmt"
	"html"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"realtek-connect/internal/content"
	"realtek-connect/internal/docs"
	"realtek-connect/internal/features"
	"realtek-connect/internal/manual"
)

type CollectionConfig struct {
	RepoRoot    string
	ContentRoot string
}

func CollectWebsiteDocuments(cfg CollectionConfig) ([]Document, error) {
	repoRoot := cfg.RepoRoot
	if repoRoot == "" {
		repoRoot = "."
	}
	contentRoot := cfg.ContentRoot
	if contentRoot == "" {
		contentRoot = filepath.Join(repoRoot, "content")
	}
	documents := make([]Document, 0, 64)
	for _, locale := range content.SupportedLocales() {
		catalog := content.CatalogFor(locale)
		for _, feature := range catalog.Features {
			body := featureBody(feature)
			documents = append(documents, Document{
				ID:         fmt.Sprintf("feature:%s:%s", feature.Slug, locale.Code),
				Locale:     locale.Code,
				SourceType: "feature",
				Title:      feature.Title,
				URL:        content.PathForLocale(locale, "/features/"+feature.Slug),
				Body:       body,
			})
		}
		for _, section := range catalog.Docs {
			body := docBody(section)
			documents = append(documents, Document{
				ID:         fmt.Sprintf("doc:%s:%s", section.Slug, locale.Code),
				Locale:     locale.Code,
				SourceType: "docs",
				Title:      section.Title,
				URL:        content.PathForLocale(locale, "/docs/"+section.Slug),
				Body:       body,
			})
		}
	}
	docPages, err := docs.NewContentSource(filepath.Join(contentRoot, "docs")).Load()
	if err == nil {
		for localeCode, page := range docPages {
			locale := localeByCode(localeCode)
			documents = append(documents, Document{
				ID:         fmt.Sprintf("docs-index:%s", localeCode),
				Locale:     localeCode,
				SourceType: "docs",
				Title:      page.Title,
				URL:        content.PathForLocale(locale, "/docs"),
				Body:       docsContentBody(page),
			})
		}
	}
	loader := manual.NewLoader(filepath.Join(contentRoot, "manual"))
	for _, locale := range content.SupportedLocales() {
		if index, ok := loader.Index(locale); ok {
			documents = append(documents, Document{
				ID:         fmt.Sprintf("manual-index:%s", locale.Code),
				Locale:     locale.Code,
				SourceType: "manual",
				Title:      index.Title,
				URL:        content.PathForLocale(locale, "/manual"),
				Body:       manualIndexBody(index),
			})
			for _, section := range index.Sections {
				if page, ok := loader.Page(locale, section.Slug); ok {
					documents = append(documents, Document{
						ID:         fmt.Sprintf("manual:%s:%s", section.Slug, locale.Code),
						Locale:     locale.Code,
						SourceType: "manual",
						Title:      page.Title,
						URL:        content.PathForLocale(locale, "/manual/"+section.Slug),
						Body:       page.Description + "\n" + htmlToText(string(page.BodyHTML)),
					})
				}
			}
		}
		if index, ok := loader.CollectionIndex(locale, "sdk"); ok {
			documents = append(documents, Document{
				ID:         fmt.Sprintf("manual-sdk-index:%s", locale.Code),
				Locale:     locale.Code,
				SourceType: "manual",
				Title:      index.Title,
				URL:        content.PathForLocale(locale, "/manual/sdk"),
				Body:       manualIndexBody(index),
			})
			for _, section := range index.Sections {
				if page, ok := loader.Page(locale, section.Slug); ok {
					documents = append(documents, Document{
						ID:         fmt.Sprintf("manual:%s:%s", section.Slug, locale.Code),
						Locale:     locale.Code,
						SourceType: "manual",
						Title:      page.Title,
						URL:        content.PathForLocale(locale, "/manual/"+section.Slug),
						Body:       page.Description + "\n" + htmlToText(string(page.BodyHTML)),
					})
				}
			}
		}
	}
	for _, rel := range append([]string{"README.md"}, markdownFiles(filepath.Join(repoRoot, "docs"))...) {
		path := filepath.Join(repoRoot, rel)
		body, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		title := rel
		if heading := firstMarkdownHeading(string(body)); heading != "" {
			title = heading
		}
		documents = append(documents, Document{
			ID:         "file:" + filepath.ToSlash(rel) + ":en",
			Locale:     "en",
			SourceType: "file",
			Title:      title,
			URL:        "",
			Body:       markdownToText(string(body)),
		})
	}
	return compactDocuments(documents), nil
}

func BuildIndexChunks(ctx context.Context, documents []Document, embedder Embedder) ([]IndexedChunk, error) {
	chunkTexts := make([]string, 0)
	chunkDocs := make([]Document, 0)
	chunkIDs := make([]string, 0)
	for _, doc := range documents {
		doc = normalizeDocument(doc)
		if doc.ID == "" || doc.Body == "" {
			continue
		}
		for i, text := range splitChunks(doc.Body, 900) {
			chunkTexts = append(chunkTexts, text)
			chunkDocs = append(chunkDocs, doc)
			chunkIDs = append(chunkIDs, fmt.Sprintf("%s:%d", doc.ID, i+1))
		}
	}
	if len(chunkTexts) == 0 {
		return nil, nil
	}
	embeddings, err := embedder.Embed(ctx, chunkTexts)
	if err != nil {
		return nil, err
	}
	if len(embeddings) != len(chunkTexts) {
		return nil, fmt.Errorf("embedding count mismatch")
	}
	chunks := make([]IndexedChunk, 0, len(chunkTexts))
	for i := range chunkTexts {
		chunks = append(chunks, IndexedChunk{
			Document:  chunkDocs[i],
			ChunkID:   chunkIDs[i],
			Text:      chunkTexts[i],
			Embedding: embeddings[i],
		})
	}
	return chunks, nil
}

func featureBody(feature features.Feature) string {
	parts := []string{feature.Title, feature.Kicker, feature.Summary, feature.Description}
	parts = append(parts, feature.Highlights...)
	parts = append(parts, feature.Capabilities...)
	parts = append(parts, feature.Outcomes...)
	for _, section := range feature.Sections {
		parts = append(parts, section.Eyebrow, section.Title, section.Intro)
		parts = append(parts, section.Items...)
	}
	if feature.Table.Title != "" {
		parts = append(parts, feature.Table.Eyebrow, feature.Table.Title, feature.Table.Intro)
		for _, row := range feature.Table.Rows {
			parts = append(parts, row.Cells...)
		}
	}
	return strings.Join(parts, "\n")
}

func docBody(section docs.Section) string {
	parts := []string{section.Title, section.Kicker, section.Summary, section.Description}
	parts = append(parts, section.Highlights...)
	parts = append(parts, section.Deliverables...)
	parts = append(parts, section.Audience...)
	if section.Table.Title != "" {
		parts = append(parts, section.Table.Eyebrow, section.Table.Title, section.Table.Intro)
		for _, row := range section.Table.Rows {
			parts = append(parts, row.Cells...)
		}
	}
	return strings.Join(parts, "\n")
}

func docsContentBody(page docs.ContentPage) string {
	parts := []string{page.Title, page.Subtitle, htmlToText(string(page.BodyHTML))}
	for _, section := range page.Sections {
		parts = append(parts, section.Title, section.Body)
	}
	return strings.Join(parts, "\n")
}

func manualIndexBody(index manual.ManualIndex) string {
	parts := []string{index.Title, index.Description}
	for _, section := range index.Sections {
		parts = append(parts, section.Title, section.Summary)
	}
	return strings.Join(parts, "\n")
}

func localeByCode(code string) content.Locale {
	for _, locale := range content.SupportedLocales() {
		if locale.Code == code {
			return locale
		}
	}
	return content.DefaultLocale()
}

func compactDocuments(input []Document) []Document {
	out := make([]Document, 0, len(input))
	for _, doc := range input {
		doc = normalizeDocument(doc)
		if doc.ID == "" || doc.Title == "" || doc.Body == "" {
			continue
		}
		out = append(out, doc)
	}
	return out
}

func splitChunks(body string, maxRunes int) []string {
	body = strings.Join(strings.Fields(body), " ")
	if body == "" {
		return nil
	}
	runes := []rune(body)
	if len(runes) <= maxRunes {
		return []string{body}
	}
	chunks := make([]string, 0, len(runes)/maxRunes+1)
	for len(runes) > 0 {
		end := maxRunes
		if len(runes) < end {
			end = len(runes)
		}
		chunks = append(chunks, strings.TrimSpace(string(runes[:end])))
		runes = runes[end:]
	}
	return chunks
}

func markdownFiles(root string) []string {
	result := []string{}
	_ = filepath.WalkDir(root, func(path string, entry os.DirEntry, err error) error {
		if err != nil || entry.IsDir() || filepath.Ext(path) != ".md" {
			return nil
		}
		rel, err := filepath.Rel(filepath.Dir(root), path)
		if err != nil {
			return nil
		}
		result = append(result, filepath.ToSlash(rel))
		return nil
	})
	return result
}

func firstMarkdownHeading(body string) string {
	for _, line := range strings.Split(body, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "# ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "# "))
		}
	}
	return ""
}

func markdownToText(body string) string {
	lines := make([]string, 0)
	for _, line := range strings.Split(body, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "```") {
			continue
		}
		lines = append(lines, strings.Trim(line, "#-*` "))
	}
	return strings.Join(lines, "\n")
}

var htmlTagPattern = regexp.MustCompile(`<[^>]+>`)

func htmlToText(body string) string {
	body = htmlTagPattern.ReplaceAllString(body, " ")
	return html.UnescapeString(strings.Join(strings.Fields(body), " "))
}
