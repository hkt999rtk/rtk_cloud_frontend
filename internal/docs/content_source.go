package docs

import (
	"bytes"
	"fmt"
	"html"
	"html/template"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

type ContentPage struct {
	Title        string
	Subtitle     string
	HeroImage    string
	HeroImageAlt string
	Sections     []ContentSection
	SEO          ContentSEO
	BodyHTML     template.HTML
}

type ContentSection struct {
	Title string `yaml:"title"`
	Body  string `yaml:"body"`
	Icon  string `yaml:"icon"`
}

type ContentSEO struct {
	MetaTitle       string `yaml:"meta_title"`
	MetaDescription string `yaml:"meta_description"`
	SocialImage     string `yaml:"social_image"`
}

type contentFrontmatter struct {
	Title        string           `yaml:"title"`
	Subtitle     string           `yaml:"subtitle"`
	HeroImage    string           `yaml:"hero_image"`
	HeroImageAlt string           `yaml:"hero_image_alt"`
	Sections     []ContentSection `yaml:"sections"`
	SEO          ContentSEO       `yaml:"seo"`
}

type ContentSource struct {
	root string
}

func NewContentSource(root string) ContentSource {
	if root == "" {
		root = "content/docs"
	}
	return ContentSource{root: root}
}

func (s ContentSource) Load() (map[string]ContentPage, error) {
	en, err := s.LoadLocale("en")
	if err != nil {
		return nil, err
	}

	pages := map[string]ContentPage{"en": en}
	for _, locale := range []string{"zh-TW", "zh-CN"} {
		page, err := s.LoadLocale(locale)
		if err != nil {
			if os.IsNotExist(err) {
				pages[locale] = en
				continue
			}
			return nil, err
		}
		pages[locale] = page
	}
	return pages, nil
}

func (s ContentSource) LoadLocale(locale string) (ContentPage, error) {
	path := filepath.Join(s.root, locale, "docs.yaml")
	body, err := os.ReadFile(path)
	if err != nil {
		return ContentPage{}, err
	}
	return ParseContentPage(body)
}

func ParseContentPage(input []byte) (ContentPage, error) {
	frontmatter, markdown, err := splitFrontmatter(input)
	if err != nil {
		return ContentPage{}, err
	}

	var meta contentFrontmatter
	if err := yaml.Unmarshal(frontmatter, &meta); err != nil {
		return ContentPage{}, fmt.Errorf("parse docs frontmatter: %w", err)
	}
	if strings.TrimSpace(meta.Title) == "" {
		return ContentPage{}, fmt.Errorf("docs frontmatter title is required")
	}
	if strings.TrimSpace(meta.Subtitle) == "" {
		return ContentPage{}, fmt.Errorf("docs frontmatter subtitle is required")
	}

	return ContentPage{
		Title:        meta.Title,
		Subtitle:     meta.Subtitle,
		HeroImage:    strings.TrimSpace(meta.HeroImage),
		HeroImageAlt: strings.TrimSpace(meta.HeroImageAlt),
		Sections:     meta.Sections,
		SEO:          meta.SEO,
		BodyHTML:     renderMarkdown(markdown),
	}, nil
}

func splitFrontmatter(input []byte) ([]byte, []byte, error) {
	normalized := bytes.ReplaceAll(input, []byte("\r\n"), []byte("\n"))
	trimmed := bytes.TrimSpace(normalized)
	if !bytes.HasPrefix(trimmed, []byte("---\n")) {
		return nil, nil, fmt.Errorf("docs content must start with YAML frontmatter")
	}
	rest := trimmed[len("---\n"):]
	end := bytes.Index(rest, []byte("\n---\n"))
	if end < 0 {
		return nil, nil, fmt.Errorf("docs content must close YAML frontmatter")
	}
	return rest[:end], rest[end+len("\n---\n"):], nil
}

var markdownImagePattern = regexp.MustCompile(`^!\[([^\]]*)\]\(([^)]+)\)$`)

func renderMarkdown(input []byte) template.HTML {
	var out strings.Builder
	paragraph := make([]string, 0, 4)
	flushParagraph := func() {
		if len(paragraph) == 0 {
			return
		}
		out.WriteString("<p>")
		out.WriteString(html.EscapeString(strings.Join(paragraph, " ")))
		out.WriteString("</p>\n")
		paragraph = paragraph[:0]
	}

	for _, rawLine := range strings.Split(string(input), "\n") {
		line := strings.TrimSpace(rawLine)
		if line == "" {
			flushParagraph()
			continue
		}
		if matches := markdownImagePattern.FindStringSubmatch(line); matches != nil {
			flushParagraph()
			out.WriteString(`<img class="docs-content-image" src="`)
			out.WriteString(html.EscapeString(matches[2]))
			out.WriteString(`" alt="`)
			out.WriteString(html.EscapeString(matches[1]))
			out.WriteString(`">` + "\n")
			continue
		}
		switch {
		case strings.HasPrefix(line, "### "):
			flushParagraph()
			out.WriteString("<h3>")
			out.WriteString(html.EscapeString(strings.TrimSpace(strings.TrimPrefix(line, "### "))))
			out.WriteString("</h3>\n")
		case strings.HasPrefix(line, "## "):
			flushParagraph()
			out.WriteString("<h2>")
			out.WriteString(html.EscapeString(strings.TrimSpace(strings.TrimPrefix(line, "## "))))
			out.WriteString("</h2>\n")
		default:
			paragraph = append(paragraph, line)
		}
	}
	flushParagraph()
	return template.HTML(out.String())
}
