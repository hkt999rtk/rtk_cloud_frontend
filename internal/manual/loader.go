package manual

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync"

	"gopkg.in/yaml.v3"

	"realtek-connect/internal/content"
)

type Loader struct {
	root string

	mu      sync.RWMutex
	indexes map[string]ManualIndex
	pages   map[string]map[string]ManualPage
}

type indexFile struct {
	Title       string          `yaml:"title"`
	Description string          `yaml:"description"`
	Sections    []ManualSection `yaml:"sections"`
}

type pageFile struct {
	Title       string `yaml:"title"`
	Description string `yaml:"description"`
}

var manualPagePattern = regexp.MustCompile(`^(.+)\.(en|zh-TW|zh-CN)\.md$`)

func NewLoader(root string) *Loader {
	if root == "" {
		root = filepath.Join("content", "manual")
	}
	l := &Loader{root: root}
	_ = l.Reload()
	return l
}

func (l *Loader) Reload() error {
	if l == nil {
		return nil
	}

	state := struct {
		indexes map[string]ManualIndex
		pages   map[string]map[string]ManualPage
	}{
		indexes: map[string]ManualIndex{},
		pages:   map[string]map[string]ManualPage{},
	}

	errs := make([]error, 0, 4)

	indexes, indexErrs := l.loadIndexes()
	errs = append(errs, indexErrs...)
	for locale, index := range indexes {
		state.indexes[locale] = index
	}

	pages, pageErrs := l.loadPages()
	errs = append(errs, pageErrs...)
	for locale, localePages := range pages {
		state.pages[locale] = localePages
	}

	if enIndex, ok := state.indexes["en"]; ok {
		for _, locale := range []string{"zh-TW", "zh-CN"} {
			if _, exists := state.indexes[locale]; !exists {
				state.indexes[locale] = enIndex
			}
		}
	}
	if enPages, ok := state.pages["en"]; ok {
		for _, locale := range []string{"zh-TW", "zh-CN"} {
			if _, exists := state.pages[locale]; !exists {
				state.pages[locale] = clonePageMap(enPages)
				continue
			}
			for slug, page := range enPages {
				if _, exists := state.pages[locale][slug]; !exists {
					state.pages[locale][slug] = page
				}
			}
		}
	}

	l.mu.Lock()
	l.indexes = state.indexes
	l.pages = state.pages
	l.mu.Unlock()

	return joinErrors(errs)
}

func (l *Loader) Index(locale content.Locale) (ManualIndex, bool) {
	if l == nil {
		return ManualIndex{}, false
	}

	l.mu.RLock()
	defer l.mu.RUnlock()
	for _, code := range localeFallbacks(locale.Code) {
		if index, ok := l.indexes[code]; ok {
			return index, true
		}
	}
	return ManualIndex{}, false
}

func (l *Loader) Page(locale content.Locale, slug string) (ManualPage, bool) {
	if l == nil {
		return ManualPage{}, false
	}

	l.mu.RLock()
	defer l.mu.RUnlock()
	for _, code := range localeFallbacks(locale.Code) {
		if localePages, ok := l.pages[code]; ok {
			if page, ok := localePages[slug]; ok {
				return page, true
			}
		}
	}
	return ManualPage{}, false
}

func (l *Loader) loadIndexes() (map[string]ManualIndex, []error) {
	result := map[string]ManualIndex{}
	errs := make([]error, 0, 3)
	for _, locale := range []string{"en", "zh-TW", "zh-CN"} {
		index, err := l.loadIndex(locale)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			errs = append(errs, err)
			continue
		}
		result[locale] = index
	}
	return result, errs
}

func (l *Loader) loadIndex(locale string) (ManualIndex, error) {
	path := filepath.Join(l.root, "index."+locale+".yaml")
	body, err := os.ReadFile(path)
	if err != nil {
		return ManualIndex{}, err
	}

	var file indexFile
	if err := yaml.Unmarshal(body, &file); err != nil {
		return ManualIndex{}, fmt.Errorf("parse manual index %s: %w", locale, err)
	}
	if strings.TrimSpace(file.Title) == "" {
		return ManualIndex{}, fmt.Errorf("manual index %s title is required", locale)
	}
	if strings.TrimSpace(file.Description) == "" {
		return ManualIndex{}, fmt.Errorf("manual index %s description is required", locale)
	}

	sections := make([]ManualSection, 0, len(file.Sections))
	for _, section := range file.Sections {
		section.Slug = strings.TrimSpace(section.Slug)
		section.Title = strings.TrimSpace(section.Title)
		section.Summary = strings.TrimSpace(section.Summary)
		if section.Slug == "" || section.Title == "" || section.Summary == "" {
			return ManualIndex{}, fmt.Errorf("manual index %s contains an incomplete section entry", locale)
		}
		sections = append(sections, section)
	}

	return ManualIndex{
		Title:       strings.TrimSpace(file.Title),
		Description: strings.TrimSpace(file.Description),
		Sections:    sections,
	}, nil
}

func (l *Loader) loadPages() (map[string]map[string]ManualPage, []error) {
	result := map[string]map[string]ManualPage{}
	errs := make([]error, 0, 3)

	entries, err := os.ReadDir(l.root)
	if err != nil {
		if os.IsNotExist(err) {
			return result, nil
		}
		return result, []error{fmt.Errorf("read manual directory: %w", err)}
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		matches := manualPagePattern.FindStringSubmatch(entry.Name())
		if matches == nil {
			continue
		}
		slug := matches[1]
		locale := matches[2]
		page, err := l.loadPage(locale, slug)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		if _, ok := result[locale]; !ok {
			result[locale] = map[string]ManualPage{}
		}
		result[locale][slug] = page
	}

	return result, errs
}

func (l *Loader) loadPage(locale, slug string) (ManualPage, error) {
	path := filepath.Join(l.root, slug+"."+locale+".md")
	body, err := os.ReadFile(path)
	if err != nil {
		return ManualPage{}, err
	}

	frontmatter, markdown, err := splitFrontmatter(body)
	if err != nil {
		return ManualPage{}, fmt.Errorf("parse manual page %s %s: %w", locale, slug, err)
	}

	var file pageFile
	if err := yaml.Unmarshal(frontmatter, &file); err != nil {
		return ManualPage{}, fmt.Errorf("parse manual page %s %s frontmatter: %w", locale, slug, err)
	}
	if strings.TrimSpace(file.Title) == "" {
		return ManualPage{}, fmt.Errorf("manual page %s %s title is required", locale, slug)
	}
	if strings.TrimSpace(file.Description) == "" {
		return ManualPage{}, fmt.Errorf("manual page %s %s description is required", locale, slug)
	}

	return ManualPage{
		Slug:        slug,
		Title:       strings.TrimSpace(file.Title),
		Description: strings.TrimSpace(file.Description),
		BodyHTML:    RenderMarkdown(markdown),
	}, nil
}

func splitFrontmatter(input []byte) ([]byte, []byte, error) {
	normalized := bytes.ReplaceAll(input, []byte("\r\n"), []byte("\n"))
	trimmed := bytes.TrimSpace(normalized)
	if !bytes.HasPrefix(trimmed, []byte("---\n")) {
		return nil, nil, fmt.Errorf("manual page must start with YAML frontmatter")
	}
	rest := trimmed[len("---\n"):]
	end := bytes.Index(rest, []byte("\n---\n"))
	if end < 0 {
		return nil, nil, fmt.Errorf("manual page must close YAML frontmatter")
	}
	return rest[:end], rest[end+len("\n---\n"):], nil
}

func clonePageMap(src map[string]ManualPage) map[string]ManualPage {
	dst := make(map[string]ManualPage, len(src))
	for slug, page := range src {
		dst[slug] = page
	}
	return dst
}

func localeFallbacks(locale string) []string {
	if locale == "en" || locale == "" {
		return []string{"en"}
	}
	return []string{locale, "en"}
}

func joinErrors(errs []error) error {
	filtered := make([]error, 0, len(errs))
	for _, err := range errs {
		if err != nil {
			filtered = append(filtered, err)
		}
	}
	if len(filtered) == 0 {
		return nil
	}
	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].Error() < filtered[j].Error()
	})
	return errors.Join(filtered...)
}
