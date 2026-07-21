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

	mu                sync.RWMutex
	indexes           map[string]ManualIndex
	collectionIndexes map[string]map[string]ManualIndex
	pages             map[string]map[string]ManualPage
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
		indexes           map[string]ManualIndex
		collectionIndexes map[string]map[string]ManualIndex
		pages             map[string]map[string]ManualPage
	}{
		indexes:           map[string]ManualIndex{},
		collectionIndexes: map[string]map[string]ManualIndex{},
		pages:             map[string]map[string]ManualPage{},
	}

	errs := make([]error, 0, 4)

	indexes, indexErrs := l.loadIndexes()
	errs = append(errs, indexErrs...)
	for locale, index := range indexes {
		state.indexes[locale] = index
	}
	collectionIndexes, collectionErrs := l.loadCollectionIndexes()
	errs = append(errs, collectionErrs...)
	for collection, localized := range collectionIndexes {
		state.collectionIndexes[collection] = localized
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
				state.pages[locale] = clonePageMap(enPages, manualLocalePrefix(locale))
				continue
			}
			for slug, page := range enPages {
				if _, exists := state.pages[locale][slug]; !exists {
					page.BodyHTML = LocalizeRenderedHTML(page.BodyHTML, manualLocalePrefix(locale))
					state.pages[locale][slug] = page
				}
			}
		}
	}
	for _, localized := range state.collectionIndexes {
		if enIndex, ok := localized["en"]; ok {
			for _, locale := range []string{"zh-TW", "zh-CN"} {
				if _, exists := localized[locale]; !exists {
					localized[locale] = enIndex
				}
			}
		}
	}

	l.mu.Lock()
	l.indexes = state.indexes
	l.collectionIndexes = state.collectionIndexes
	l.pages = state.pages
	l.mu.Unlock()

	return joinErrors(errs)
}

func (l *Loader) CollectionIndex(locale content.Locale, collection string) (ManualIndex, bool) {
	if l == nil {
		return ManualIndex{}, false
	}
	collection = strings.Trim(strings.TrimSpace(collection), "/")
	if collection == "" {
		return l.Index(locale)
	}

	l.mu.RLock()
	defer l.mu.RUnlock()
	localized, ok := l.collectionIndexes[collection]
	if !ok {
		return ManualIndex{}, false
	}
	for _, code := range localeFallbacks(locale.Code) {
		if index, ok := localized[code]; ok {
			return index, true
		}
	}
	return ManualIndex{}, false
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
	return loadIndexFile(path, locale, "")
}

func loadIndexFile(path, locale, collection string) (ManualIndex, error) {
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
		if collection != "" && !strings.HasPrefix(section.Slug, collection+"/") {
			section.Slug = collection + "/" + section.Slug
		}
		sections = append(sections, section)
	}

	return ManualIndex{
		Title:       strings.TrimSpace(file.Title),
		Description: strings.TrimSpace(file.Description),
		Sections:    sections,
	}, nil
}

func (l *Loader) loadCollectionIndexes() (map[string]map[string]ManualIndex, []error) {
	result := map[string]map[string]ManualIndex{}
	errs := []error{}
	if _, err := os.Stat(l.root); os.IsNotExist(err) {
		return result, nil
	}
	err := filepath.WalkDir(l.root, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			errs = append(errs, walkErr)
			return nil
		}
		if entry.IsDir() || filepath.Dir(path) == l.root {
			return nil
		}
		matches := regexp.MustCompile(`^index\.(en|zh-TW|zh-CN)\.yaml$`).FindStringSubmatch(entry.Name())
		if matches == nil {
			return nil
		}
		collection, relErr := filepath.Rel(l.root, filepath.Dir(path))
		if relErr != nil {
			errs = append(errs, relErr)
			return nil
		}
		collection = filepath.ToSlash(collection)
		index, loadErr := loadIndexFile(path, matches[1], collection)
		if loadErr != nil {
			errs = append(errs, loadErr)
			return nil
		}
		if result[collection] == nil {
			result[collection] = map[string]ManualIndex{}
		}
		result[collection][matches[1]] = index
		return nil
	})
	if err != nil {
		errs = append(errs, err)
	}
	return result, errs
}

func (l *Loader) loadPages() (map[string]map[string]ManualPage, []error) {
	result := map[string]map[string]ManualPage{}
	errs := make([]error, 0, 3)
	if _, err := os.Stat(l.root); os.IsNotExist(err) {
		return result, nil
	}

	err := filepath.WalkDir(l.root, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			errs = append(errs, walkErr)
			return nil
		}
		if entry.IsDir() {
			return nil
		}
		matches := manualPagePattern.FindStringSubmatch(entry.Name())
		if matches == nil {
			return nil
		}
		relative, relErr := filepath.Rel(l.root, path)
		if relErr != nil {
			errs = append(errs, relErr)
			return nil
		}
		slug := strings.TrimSuffix(filepath.ToSlash(relative), "."+matches[2]+".md")
		locale := matches[2]
		page, err := l.loadPage(locale, slug)
		if err != nil {
			errs = append(errs, err)
			return nil
		}
		if _, ok := result[locale]; !ok {
			result[locale] = map[string]ManualPage{}
		}
		result[locale][slug] = page
		return nil
	})
	if err != nil && !os.IsNotExist(err) {
		errs = append(errs, fmt.Errorf("read manual directory: %w", err))
	}

	return result, errs
}

func (l *Loader) loadPage(locale, slug string) (ManualPage, error) {
	path := filepath.Join(l.root, filepath.FromSlash(slug)+"."+locale+".md")
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
		BodyHTML:    RenderMarkdownWithPrefix(markdown, manualLocalePrefix(locale)),
	}, nil
}

func manualLocalePrefix(locale string) string {
	switch locale {
	case "zh-TW":
		return "/zh-tw"
	case "zh-CN":
		return "/zh-cn"
	default:
		return ""
	}
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

func clonePageMap(src map[string]ManualPage, localePrefix string) map[string]ManualPage {
	dst := make(map[string]ManualPage, len(src))
	for slug, page := range src {
		page.BodyHTML = LocalizeRenderedHTML(page.BodyHTML, localePrefix)
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
