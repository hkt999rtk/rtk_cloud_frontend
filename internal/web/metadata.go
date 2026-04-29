package web

import (
	"encoding/xml"
	"net/http"
	"strings"

	"realtek-connect/internal/content"
)

const (
	heroImagePath = "/static/assets/connectplus-hero.png"
	heroImageAlt  = "Realtek Connect+ device, cloud, mobile app, and dashboard platform flow"
)

type sitemapURLSet struct {
	XMLName xml.Name     `xml:"urlset"`
	XMLNS   string       `xml:"xmlns,attr"`
	URLs    []sitemapURL `xml:"url"`
}

type sitemapURL struct {
	Loc string `xml:"loc"`
}

func (s *Server) basePageData(r *http.Request, locale content.Locale, publicPath, title, description string) pageData {
	catalog := content.CatalogFor(locale)
	data := pageData{
		Title:           title,
		MetaDescription: description,
		CanonicalURL:    absoluteURL(r, content.PathForLocale(locale, publicPath)),
		SocialImageURL:  absoluteURL(r, heroImagePath),
		SocialImageAlt:  heroImageAlt,
		CurrentPath:     r.URL.Path,
		PublicPath:      publicPath,
		Lang:            locale.Lang,
		Locale:          locale,
		LocalePrefix:    locale.Prefix,
		Text:            catalog.Text,
		AlternateLinks:  alternateLinks(r, publicPath, locale),
		Docs:            catalog.Docs,
		Features:        catalog.Features,
	}
	if s.disableSearchIndexing {
		data.MetaRobots = "noindex, nofollow, noarchive"
	}
	return data
}

func (s *Server) adminPageData(r *http.Request, title, description string) pageData {
	data := s.basePageData(r, content.DefaultLocale(), r.URL.Path, title, description)
	data.MetaRobots = "noindex, nofollow"
	data.AlternateLinks = nil
	return data
}

func alternateLinks(r *http.Request, publicPath string, current content.Locale) []content.AlternateLink {
	locales := content.SupportedLocales()
	links := make([]content.AlternateLink, 0, len(locales)+1)
	for _, locale := range locales {
		links = append(links, content.AlternateLink{
			HrefLang: locale.Lang,
			Label:    locale.Label,
			Href:     absoluteURL(r, content.PathForLocale(locale, publicPath)),
			Current:  locale.Code == current.Code,
		})
	}
	links = append(links, content.AlternateLink{
		HrefLang: "x-default",
		Label:    "Default",
		Href:     absoluteURL(r, content.PathForLocale(content.DefaultLocale(), publicPath)),
		Current:  false,
	})
	return links
}

func (s *Server) handleRobotsTxt(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/robots.txt" {
		http.NotFound(w, r)
		return
	}
	if r.Method != http.MethodGet {
		methodNotAllowed(w)
		return
	}

	if s.disableSearchIndexing {
		body := strings.Join([]string{
			"User-agent: *",
			"Disallow: /",
			"",
		}, "\n")

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(body))
		return
	}

	body := strings.Join([]string{
		"User-agent: *",
		"Allow: /",
		"Disallow: /admin/",
		"Disallow: /healthz",
		"Sitemap: " + absoluteURL(r, "/sitemap.xml"),
		"",
	}, "\n")

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(body))
}

func (s *Server) handleSitemapXML(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/sitemap.xml" {
		http.NotFound(w, r)
		return
	}
	if r.Method != http.MethodGet {
		methodNotAllowed(w)
		return
	}
	if s.disableSearchIndexing {
		http.NotFound(w, r)
		return
	}

	paths := publicSitemapPaths()

	payload := sitemapURLSet{
		XMLNS: "http://www.sitemaps.org/schemas/sitemap/0.9",
		URLs:  make([]sitemapURL, 0, len(paths)),
	}
	for _, path := range paths {
		payload.URLs = append(payload.URLs, sitemapURL{Loc: absoluteURL(r, path)})
	}

	body, err := xml.MarshalIndent(payload, "", "  ")
	if err != nil {
		http.Error(w, "could not build sitemap", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(xml.Header))
	_, _ = w.Write(body)
	_, _ = w.Write([]byte("\n"))
}

func publicSitemapPaths() []string {
	catalog := content.CatalogFor(content.DefaultLocale())
	basePaths := []string{"/", "/docs", "/features", "/contact"}
	for _, section := range catalog.Docs {
		basePaths = append(basePaths, "/docs/"+section.Slug)
	}
	for _, feature := range catalog.Features {
		basePaths = append(basePaths, "/features/"+feature.Slug)
	}

	paths := make([]string, 0, len(basePaths)*len(content.SupportedLocales()))
	for _, locale := range content.SupportedLocales() {
		for _, path := range basePaths {
			paths = append(paths, content.PathForLocale(locale, path))
		}
	}
	return paths
}

func absoluteURL(r *http.Request, path string) string {
	host := strings.TrimSpace(r.Header.Get("X-Forwarded-Host"))
	if host == "" {
		host = r.Host
	}
	if host == "" {
		host = "localhost"
	}

	return requestScheme(r) + "://" + host + path
}

func requestScheme(r *http.Request) string {
	forwarded := strings.TrimSpace(r.Header.Get("X-Forwarded-Proto"))
	if forwarded != "" {
		return strings.TrimSpace(strings.Split(forwarded, ",")[0])
	}
	if r.TLS != nil {
		return "https"
	}
	return "http"
}
