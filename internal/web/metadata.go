package web

import (
	"encoding/xml"
	"net/http"
	"strings"

	"realtek-connect/internal/docs"
	"realtek-connect/internal/features"
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

func (s *Server) basePageData(r *http.Request, title, description string) pageData {
	return pageData{
		Title:           title,
		MetaDescription: description,
		CanonicalURL:    absoluteURL(r, r.URL.Path),
		SocialImageURL:  absoluteURL(r, heroImagePath),
		SocialImageAlt:  heroImageAlt,
		CurrentPath:     r.URL.Path,
		Docs:            docs.All(),
		Features:        features.All(),
	}
}

func (s *Server) adminPageData(r *http.Request, title, description string) pageData {
	data := s.basePageData(r, title, description)
	data.MetaRobots = "noindex, nofollow"
	return data
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

	paths := []string{"/", "/docs", "/features", "/contact"}
	for _, section := range docs.All() {
		paths = append(paths, "/docs/"+section.Slug)
	}
	for _, feature := range features.All() {
		paths = append(paths, "/features/"+feature.Slug)
	}

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
