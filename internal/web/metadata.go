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
		CanonicalURL:    s.absoluteURL(r, content.PathForLocale(locale, publicPath)),
		SocialImageURL:  s.absoluteURL(r, s.assetPath(heroImagePath)),
		SocialImageAlt:  heroImageAlt,
		CurrentPath:     r.URL.Path,
		PublicPath:      publicPath,
		Lang:            locale.Lang,
		Locale:          locale,
		LocalePrefix:    locale.Prefix,
		Text:            catalog.Text,
		AlternateLinks:  s.alternateLinks(r, publicPath, locale),
		FooterSitemap:   s.footerSitemap(locale, catalog),
		Docs:            catalog.Docs,
		Features:        catalog.Features,
		Analytics: pageAnalyticsView{
			Enabled: s.analyticsStore != nil && isPublicAnalyticsPage(publicPath),
		},
		AnalyticsEndpoint: "/api/event",
		AnalyticsPage:     analyticsPageKey(publicPath),
		InterestOptions:   catalog.ContactInterestOptions(),
	}
	if s.disableSearchIndexing {
		data.MetaRobots = "noindex, nofollow, noarchive"
	}
	return data
}

func (s *Server) footerSitemap(locale content.Locale, catalog content.Catalog) []footerSitemapGroup {
	manualIndex, _ := s.manualIndexFor(locale)

	groups := []footerSitemapGroup{
		{
			Title: catalog.T("footer.group.platform"),
			Links: []footerSitemapLink{
				{Label: catalog.T("footer.home"), Href: content.PathForLocale(locale, "/")},
				{Label: catalog.T("footer.features"), Href: content.PathForLocale(locale, "/features")},
				{Label: catalog.T("footer.docs"), Href: content.PathForLocale(locale, "/docs")},
				{Label: catalog.T("footer.manual"), Href: content.PathForLocale(locale, "/manual")},
			},
		},
		{
			Title: catalog.T("footer.group.features"),
			Links: make([]footerSitemapLink, 0, len(catalog.Features)),
		},
		{
			Title: catalog.T("footer.group.docs"),
			Links: make([]footerSitemapLink, 0, len(catalog.Docs)),
		},
		{
			Title: catalog.T("footer.group.manual"),
			Links: make([]footerSitemapLink, 0, len(manualIndex.Sections)),
		},
		{
			Title: catalog.T("footer.group.company"),
			Links: []footerSitemapLink{
				{Label: catalog.T("footer.contact"), Href: content.PathForLocale(locale, "/contact")},
				{Label: catalog.T("footer.privacy"), Href: content.PathForLocale(locale, "/privacy")},
			},
		},
	}

	for _, feature := range catalog.Features {
		groups[1].Links = append(groups[1].Links, footerSitemapLink{
			Label: feature.Title,
			Href:  content.PathForLocale(locale, "/features/"+feature.Slug),
		})
	}
	for _, section := range catalog.Docs {
		groups[2].Links = append(groups[2].Links, footerSitemapLink{
			Label: section.Title,
			Href:  content.PathForLocale(locale, "/docs/"+section.Slug),
		})
	}
	for _, section := range manualIndex.Sections {
		groups[3].Links = append(groups[3].Links, footerSitemapLink{
			Label: section.Title,
			Href:  content.PathForLocale(locale, "/manual/"+section.Slug),
		})
	}

	return groups
}

func (s *Server) adminPageData(r *http.Request, title, description string) pageData {
	data := s.basePageData(r, content.DefaultLocale(), r.URL.Path, title, description)
	data.MetaRobots = "noindex, nofollow"
	data.AlternateLinks = nil
	return data
}

func (s *Server) alternateLinks(r *http.Request, publicPath string, current content.Locale) []content.AlternateLink {
	locales := content.SupportedLocales()
	links := make([]content.AlternateLink, 0, len(locales)+1)
	for _, locale := range locales {
		links = append(links, content.AlternateLink{
			HrefLang: locale.Lang,
			Label:    locale.Label,
			Href:     s.absoluteURL(r, content.PathForLocale(locale, publicPath)),
			Current:  locale.Code == current.Code,
		})
	}
	links = append(links, content.AlternateLink{
		HrefLang: "x-default",
		Label:    "Default",
		Href:     s.absoluteURL(r, content.PathForLocale(content.DefaultLocale(), publicPath)),
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
		"Sitemap: " + s.absoluteURL(r, "/sitemap.xml"),
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
		payload.URLs = append(payload.URLs, sitemapURL{Loc: s.absoluteURL(r, path)})
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
	basePaths := []string{
		"/",
		"/docs",
		"/manual",
		"/features",
		"/contact",
		"/privacy",
		"/manual/getting-started",
		"/manual/deployment-notes",
		"/manual/reference",
		"/manual/sdk-samples",
	}
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

func isPublicAnalyticsPage(publicPath string) bool {
	key := analyticsPageKey(publicPath)
	_, ok := analyticsAllowedPages[key]
	return ok
}

func analyticsPageKey(publicPath string) string {
	path := "/" + strings.Trim(strings.TrimSpace(publicPath), "/")
	switch {
	case path == "/":
		return "home"
	case path == "/features":
		return "features"
	case strings.HasPrefix(path, "/features/"):
		return strings.TrimPrefix(path, "/features/")
	case path == "/docs":
		return "docs"
	case strings.HasPrefix(path, "/docs/"):
		return strings.TrimPrefix(path, "/docs/")
	case path == "/contact":
		return "contact"
	case path == "/privacy":
		return "privacy"
	default:
		return ""
	}
}

func (s *Server) absoluteURL(r *http.Request, path string) string {
	if s.publicBaseURL != "" {
		return s.publicBaseURL + path
	}
	return absoluteURL(r, path)
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
