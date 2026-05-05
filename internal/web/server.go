package web

import (
	"context"
	"crypto/sha256"
	"encoding/csv"
	"encoding/hex"
	"fmt"
	"html/template"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"realtek-connect/internal/content"
	"realtek-connect/internal/docs"
	"realtek-connect/internal/features"
	"realtek-connect/internal/leads"
)

type LeadStore interface {
	Insert(context.Context, leads.Lead) error
	Count(context.Context, leads.ListFilter) (int, error)
	List(context.Context, leads.ListOptions) ([]leads.LeadRecord, error)
}

type Config struct {
	TemplatesDir            string
	StaticDir               string
	ContentDir              string
	LeadStore               LeadStore
	AdminToken              string
	DisableSearchIndexing   bool
	PublicBaseURL           string
	EnableAssetFingerprints bool
	EnableCDNCacheHeaders   bool
}

type Server struct {
	templatesDir            string
	staticDir               string
	contentDir              string
	leadStore               LeadStore
	adminToken              string
	disableSearchIndexing   bool
	publicBaseURL           string
	enableAssetFingerprints bool
	enableCDNCacheHeaders   bool
	assetVersions           map[string]string
	contactLimit            *submissionRateLimiter
	docsContentMu           sync.RWMutex
	docsContent             map[string]docs.ContentPage
}

type pageData struct {
	Title           string
	MetaDescription string
	CanonicalURL    string
	SocialImageURL  string
	SocialImageAlt  string
	MetaRobots      string
	CurrentPath     string
	PublicPath      string
	Lang            string
	Locale          content.Locale
	LocalePrefix    string
	Text            map[string]string
	AlternateLinks  []content.AlternateLink
	Docs            []docs.Section
	DocsPage        docs.ContentPage
	Doc             docs.Section
	Features        []features.Feature
	Feature         features.Feature
	InterestOptions []content.ContactInterestOption
	Form            contactForm
	Errors          map[string]string
	Success         bool
	SubmittedFor    string
	Leads           []leads.LeadRecord
	AdminEnabled    bool
	AdminCSVHref    string
	LeadFilters     adminLeadFilters
	LeadPagination  adminLeadPagination
}

type contactForm struct {
	Name     string
	Company  string
	Email    string
	Interest string
	Message  string
	Website  string
}

const adminLeadPageSize = 25

type adminLeadFilters struct {
	Email     string
	Company   string
	Interest  string
	Token     string
	HasActive bool
	ClearHref string
}

type adminLeadPagination struct {
	Page         int
	PageSize     int
	TotalCount   int
	TotalPages   int
	Start        int
	End          int
	PreviousHref string
	NextHref     string
}

func NewServer(cfg Config) (*Server, error) {
	if cfg.TemplatesDir == "" {
		cfg.TemplatesDir = "templates"
	}
	if cfg.StaticDir == "" {
		cfg.StaticDir = "static"
	}
	if cfg.ContentDir == "" {
		cfg.ContentDir = "content/docs"
	}
	docsContent, err := docs.NewContentSource(cfg.ContentDir).Load()
	if err != nil {
		return nil, err
	}
	assetVersions := map[string]string{}
	if cfg.EnableAssetFingerprints {
		assetVersions = buildAssetVersions(cfg.StaticDir)
	}
	return &Server{
		templatesDir:            cfg.TemplatesDir,
		staticDir:               cfg.StaticDir,
		contentDir:              cfg.ContentDir,
		leadStore:               cfg.LeadStore,
		adminToken:              cfg.AdminToken,
		disableSearchIndexing:   cfg.DisableSearchIndexing,
		publicBaseURL:           normalizePublicBaseURL(cfg.PublicBaseURL),
		enableAssetFingerprints: cfg.EnableAssetFingerprints,
		enableCDNCacheHeaders:   cfg.EnableCDNCacheHeaders,
		assetVersions:           assetVersions,
		contactLimit:            newSubmissionRateLimiter(5, 10*time.Minute),
		docsContent:             docsContent,
	}, nil
}

func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/static/", http.StripPrefix("/static/", s.staticHandler()))
	mux.HandleFunc("/robots.txt", s.handleRobotsTxt)
	mux.HandleFunc("/sitemap.xml", s.handleSitemapXML)
	mux.HandleFunc("/admin/leads", s.handleAdminLeads)
	mux.HandleFunc("/admin/leads.csv", s.handleAdminLeadsCSV)
	mux.HandleFunc("/admin/reload-content", s.handleAdminReloadContent)
	mux.HandleFunc("/healthz", s.handleHealthz)
	mux.HandleFunc("/", s.handlePublic)
	return securityHeaders(s.searchIndexingHeaders(s.cacheHeaders(mux)))
}

func (s *Server) staticHandler() http.Handler {
	fileServer := http.FileServer(http.Dir(s.staticDir))
	if !s.enableCDNCacheHeaders {
		return fileServer
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if staticPathExists(s.staticDir, r.URL.Path) {
			w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
		}
		fileServer.ServeHTTP(w, r)
	})
}

func (s *Server) handlePublic(w http.ResponseWriter, r *http.Request) {
	locale, publicPath, ok := content.LocaleFromPath(r.URL.Path)
	if !ok {
		http.NotFound(w, r)
		return
	}
	switch {
	case publicPath == "/":
		s.handleHome(w, r, locale, publicPath)
	case publicPath == "/docs":
		s.handleDocs(w, r, locale, publicPath)
	case strings.HasPrefix(publicPath, "/docs/"):
		s.handleDocDetail(w, r, locale, publicPath)
	case publicPath == "/features":
		s.handleFeatures(w, r, locale, publicPath)
	case strings.HasPrefix(publicPath, "/features/"):
		s.handleFeatureDetail(w, r, locale, publicPath)
	case publicPath == "/contact":
		s.handleContact(w, r, locale, publicPath)
	case publicPath == "/privacy":
		s.handlePrivacy(w, r, locale, publicPath)
	default:
		http.NotFound(w, r)
	}
}

func (s *Server) handleHome(w http.ResponseWriter, r *http.Request, locale content.Locale, publicPath string) {
	if r.Method != http.MethodGet {
		methodNotAllowed(w)
		return
	}
	catalog := content.CatalogFor(locale)
	page := catalog.Page("home")
	s.render(w, http.StatusOK, "home.html", s.basePageData(r, locale, publicPath, page.Title, page.Description))
}

func (s *Server) handleDocs(w http.ResponseWriter, r *http.Request, locale content.Locale, publicPath string) {
	if r.Method != http.MethodGet {
		methodNotAllowed(w)
		return
	}
	catalog := content.CatalogFor(locale)
	docsPage := s.docsPageFor(locale)
	title := docsPage.Title + " | Realtek Connect+"
	if docsPage.SEO.MetaTitle != "" {
		title = docsPage.SEO.MetaTitle
	}
	description := docsPage.Subtitle
	if docsPage.SEO.MetaDescription != "" {
		description = docsPage.SEO.MetaDescription
	}
	data := s.basePageData(r, locale, publicPath, title, description)
	data.DocsPage = docsPage
	if docsPage.SEO.SocialImage != "" {
		data.SocialImageURL = s.absoluteURL(r, s.assetPath(docsPage.SEO.SocialImage))
	}
	if docsPage.HeroImageAlt != "" {
		data.SocialImageAlt = docsPage.HeroImageAlt
	}
	if docsPage.Title == "" {
		page := catalog.Page("docs")
		data = s.basePageData(r, locale, publicPath, page.Title, page.Description)
	}
	s.render(w, http.StatusOK, "docs.html", data)
}

func (s *Server) handlePrivacy(w http.ResponseWriter, r *http.Request, locale content.Locale, publicPath string) {
	if r.Method != http.MethodGet {
		methodNotAllowed(w)
		return
	}
	catalog := content.CatalogFor(locale)
	page := catalog.Page("privacy")
	s.render(w, http.StatusOK, "privacy.html", s.basePageData(r, locale, publicPath, page.Title, page.Description))
}

func (s *Server) handleDocDetail(w http.ResponseWriter, r *http.Request, locale content.Locale, publicPath string) {
	if r.Method != http.MethodGet {
		methodNotAllowed(w)
		return
	}

	slug := strings.TrimPrefix(publicPath, "/docs/")
	slug = strings.Trim(slug, "/")
	if slug == "" {
		http.Redirect(w, r, content.PathForLocale(locale, "/docs"), http.StatusSeeOther)
		return
	}

	catalog := content.CatalogFor(locale)
	doc, ok := catalog.DocBySlug(slug)
	if !ok {
		http.NotFound(w, r)
		return
	}

	data := s.basePageData(
		r,
		locale,
		publicPath,
		doc.Title+" | Realtek Connect+ Docs",
		doc.Summary,
	)
	data.Doc = doc
	s.render(w, http.StatusOK, "doc.html", data)
}

func (s *Server) handleHealthz(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/healthz" {
		http.NotFound(w, r)
		return
	}
	if r.Method != http.MethodGet {
		methodNotAllowed(w)
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok\n"))
}

func (s *Server) handleFeatures(w http.ResponseWriter, r *http.Request, locale content.Locale, publicPath string) {
	if r.Method != http.MethodGet {
		methodNotAllowed(w)
		return
	}
	catalog := content.CatalogFor(locale)
	page := catalog.Page("features")
	s.render(w, http.StatusOK, "features.html", s.basePageData(r, locale, publicPath, page.Title, page.Description))
}

func (s *Server) handleAdminLeads(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/admin/leads" {
		http.NotFound(w, r)
		return
	}
	if r.Method != http.MethodGet {
		methodNotAllowed(w)
		return
	}
	if !s.authorized(r) {
		s.unauthorized(w)
		return
	}

	filters := parseAdminLeadFilters(r)
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	totalCount, err := s.leadStore.Count(ctx, filters.listFilter())
	if err != nil {
		http.Error(w, "could not load leads", http.StatusInternalServerError)
		return
	}

	page := parseAdminLeadPage(r.URL.Query().Get("page"))
	totalPages := adminLeadTotalPages(totalCount, adminLeadPageSize)
	if totalPages > 0 && page > totalPages {
		page = totalPages
	}

	records, err := s.leadStore.List(ctx, leads.ListOptions{
		Filter: filters.listFilter(),
		Limit:  adminLeadPageSize,
		Offset: (page - 1) * adminLeadPageSize,
	})
	if err != nil {
		http.Error(w, "could not load leads", http.StatusInternalServerError)
		return
	}

	data := s.adminPageData(
		r,
		"Leads | Realtek Connect+",
		"Protected Realtek Connect+ lead review interface.",
	)
	data.Leads = records
	data.AdminEnabled = s.adminToken != ""
	data.AdminCSVHref = s.adminCSVHref(filters)
	data.LeadFilters = filters
	data.LeadFilters.ClearHref = s.adminLeadsHref(filters.Token, adminLeadFilters{})
	data.LeadPagination = s.adminLeadPagination(filters, page, totalCount)
	s.render(w, http.StatusOK, "admin_leads.html", data)
}

func (s *Server) handleAdminLeadsCSV(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/admin/leads.csv" {
		http.NotFound(w, r)
		return
	}
	if r.Method != http.MethodGet {
		methodNotAllowed(w)
		return
	}
	if !s.authorized(r) {
		s.unauthorized(w)
		return
	}

	filters := parseAdminLeadFilters(r)
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()
	records, err := s.leadStore.List(ctx, leads.ListOptions{
		Filter: filters.listFilter(),
	})
	if err != nil {
		http.Error(w, "could not load leads", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/csv; charset=utf-8")
	w.Header().Set("Content-Disposition", `attachment; filename="realtek-connect-leads.csv"`)
	w.WriteHeader(http.StatusOK)

	writer := csv.NewWriter(w)
	_ = writer.Write([]string{"id", "name", "company", "email", "interest", "message", "created_at"})
	for _, record := range records {
		_ = writer.Write([]string{
			strconv.FormatInt(record.ID, 10),
			record.Name,
			record.Company,
			record.Email,
			record.Interest,
			record.Message,
			formatTime(record.CreatedAt),
		})
	}
	writer.Flush()
}

func (s *Server) handleAdminReloadContent(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/admin/reload-content" {
		http.NotFound(w, r)
		return
	}
	if r.Method != http.MethodPost {
		methodNotAllowed(w)
		return
	}
	if !s.authorized(r) {
		s.unauthorized(w)
		return
	}

	loaded, err := docs.NewContentSource(s.contentDir).Load()
	if err != nil {
		http.Error(w, fmt.Sprintf("could not reload content: %v", err), http.StatusInternalServerError)
		return
	}

	s.docsContentMu.Lock()
	s.docsContent = loaded
	s.docsContentMu.Unlock()

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("content reloaded\n"))
}

func (s *Server) docsPageFor(locale content.Locale) docs.ContentPage {
	s.docsContentMu.RLock()
	defer s.docsContentMu.RUnlock()
	if page, ok := s.docsContent[locale.Code]; ok {
		return page
	}
	return s.docsContent["en"]
}

func (s *Server) authorized(r *http.Request) bool {
	if s.adminToken == "" {
		return false
	}
	token := r.Header.Get("X-Admin-Token")
	if token == "" {
		token = r.URL.Query().Get("token")
	}
	return token == s.adminToken
}

func (s *Server) adminCSVHref(filters adminLeadFilters) string {
	values := adminLeadQueryValues(filters, 0)
	if len(values) == 0 {
		return "/admin/leads.csv"
	}
	return "/admin/leads.csv?" + values.Encode()
}

func (s *Server) unauthorized(w http.ResponseWriter) {
	if s.adminToken == "" {
		http.Error(w, "404 page not found", http.StatusNotFound)
		return
	}
	http.Error(w, "unauthorized", http.StatusUnauthorized)
}

func (s *Server) handleFeatureDetail(w http.ResponseWriter, r *http.Request, locale content.Locale, publicPath string) {
	if r.Method != http.MethodGet {
		methodNotAllowed(w)
		return
	}

	slug := strings.TrimPrefix(publicPath, "/features/")
	slug = strings.Trim(slug, "/")
	if slug == "" {
		http.Redirect(w, r, content.PathForLocale(locale, "/features"), http.StatusSeeOther)
		return
	}

	catalog := content.CatalogFor(locale)
	feature, ok := catalog.FeatureBySlug(slug)
	if !ok {
		http.NotFound(w, r)
		return
	}

	data := s.basePageData(
		r,
		locale,
		publicPath,
		feature.Title+" | Realtek Connect+",
		feature.Summary,
	)
	data.Feature = feature
	s.render(w, http.StatusOK, "feature.html", data)
}

func (s *Server) handleContact(w http.ResponseWriter, r *http.Request, locale content.Locale, publicPath string) {
	switch r.Method {
	case http.MethodGet:
		catalog := content.CatalogFor(locale)
		page := catalog.Page("contact")
		s.render(w, http.StatusOK, "contact.html", s.basePageData(r, locale, publicPath, page.Title, page.Description))
	case http.MethodPost:
		s.submitContact(w, r, locale, publicPath)
	default:
		methodNotAllowed(w)
	}
}

func (s *Server) submitContact(w http.ResponseWriter, r *http.Request, locale content.Locale, publicPath string) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form payload", http.StatusBadRequest)
		return
	}

	form := contactForm{
		Name:     strings.TrimSpace(r.FormValue("name")),
		Company:  strings.TrimSpace(r.FormValue("company")),
		Email:    strings.TrimSpace(r.FormValue("email")),
		Interest: strings.TrimSpace(r.FormValue("interest")),
		Message:  strings.TrimSpace(r.FormValue("message")),
		Website:  strings.TrimSpace(r.FormValue("website")),
	}

	if isSpamContact(form) {
		catalog := content.CatalogFor(locale)
		page := catalog.Page("contact")
		data := s.basePageData(r, locale, publicPath, page.Title, page.Description)
		data.Form = form
		data.Errors = map[string]string{
			"form": localizeError("Request could not be processed.", catalog),
		}
		s.render(w, http.StatusBadRequest, "contact.html", data)
		return
	}

	catalog := content.CatalogFor(locale)
	errors := validateContact(form, catalog)
	if len(errors) > 0 {
		page := catalog.Page("contact")
		data := s.basePageData(r, locale, publicPath, page.Title, page.Description)
		data.Form = form
		data.Errors = errors
		s.render(w, http.StatusBadRequest, "contact.html", data)
		return
	}

	if s.contactLimit != nil && !s.contactLimit.Allow(contactSubmissionKey(r)) {
		page := catalog.Page("contact")
		data := s.basePageData(r, locale, publicPath, page.Title, page.Description)
		data.Form = form
		data.Errors = map[string]string{
			"form": localizeError("Too many requests from this address. Please wait a few minutes and try again.", catalog),
		}
		s.render(w, http.StatusTooManyRequests, "contact.html", data)
		return
	}

	if s.leadStore != nil {
		ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
		defer cancel()
		err := s.leadStore.Insert(ctx, leads.Lead{
			Name:     form.Name,
			Company:  form.Company,
			Email:    form.Email,
			Interest: form.Interest,
			Message:  form.Message,
		})
		if err != nil {
			http.Error(w, "could not save contact request", http.StatusInternalServerError)
			return
		}
	}

	page := catalog.Page("contact")
	data := s.basePageData(r, locale, publicPath, page.Title, page.Description)
	data.Success = true
	data.SubmittedFor = form.Name
	s.render(w, http.StatusOK, "contact.html", data)
}

func validateContact(form contactForm, catalog content.Catalog) map[string]string {
	errors := make(map[string]string)
	for field, message := range leads.Validate(leads.Lead{
		Name:     form.Name,
		Company:  form.Company,
		Email:    form.Email,
		Interest: form.Interest,
		Message:  form.Message,
	}) {
		errors[field] = localizeError(message, catalog)
	}

	if form.Email != "" && !emailPattern.MatchString(form.Email) {
		errors["email"] = localizeError("Enter a valid email address.", catalog)
	}

	if len(errors) == 0 {
		return nil
	}
	return errors
}

var emailPattern = regexp.MustCompile(`^[^@\s]+@[^@\s]+\.[^@\s]+$`)

func isSpamContact(form contactForm) bool {
	return form.Website != ""
}

func contactSubmissionKey(r *http.Request) string {
	// The app is not behind a trusted proxy chain, so client-supplied forwarding
	// headers are ignored for abuse controls.
	host, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr))
	if err == nil && host != "" {
		return host
	}
	if value := strings.TrimSpace(r.RemoteAddr); value != "" {
		return value
	}
	return "unknown"
}

func (s *Server) render(w http.ResponseWriter, status int, name string, data pageData) {
	files := []string{
		filepath.Join(s.templatesDir, "layout.html"),
		filepath.Join(s.templatesDir, name),
	}
	tmpl, err := template.New("layout.html").Funcs(template.FuncMap{
		"formatTime":    formatTime,
		"icon":          icon,
		"t":             templateText,
		"localizedPath": localizedPath,
		"asset":         s.assetPath,
	}).ParseFiles(files...)
	if err != nil {
		http.Error(w, "template parse error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	if err := tmpl.ExecuteTemplate(w, "layout", data); err != nil {
		http.Error(w, "template render error", http.StatusInternalServerError)
	}
}

func templateText(data pageData, key string) string {
	if value, ok := data.Text[key]; ok {
		return value
	}
	return key
}

func localizedPath(data pageData, publicPath string) string {
	return content.PathForLocale(data.Locale, publicPath)
}

func (s *Server) assetPath(rawPath string) string {
	cleanPath := cleanAssetPath(rawPath)
	if !s.enableAssetFingerprints {
		return cleanPath
	}
	version := s.assetVersions[cleanPath]
	if version == "" {
		return cleanPath
	}
	values := url.Values{}
	values.Set("v", version)
	return cleanPath + "?" + values.Encode()
}

func localizeError(message string, catalog content.Catalog) string {
	if catalog.Locale.Code == "en" {
		return message
	}
	localized := message
	switch {
	case message == "Name is required.":
		localized = catalog.T("contact.name") + "為必填欄位。"
	case message == "Email is required.":
		localized = catalog.T("contact.email") + "為必填欄位。"
	case message == "Enter a valid email address.":
		localized = "請輸入有效的 Email。"
	case message == "Select an area of interest.":
		localized = "請選擇關注服務。"
	case message == "Request could not be processed.":
		localized = "無法處理此請求。"
	case strings.HasPrefix(message, "Too many requests"):
		localized = "此來源送出太多請求，請稍候再試。"
	case strings.HasPrefix(message, "Name must be"):
		localized = catalog.T("contact.name") + "最多 120 個字元。"
	case strings.HasPrefix(message, "Company must be"):
		localized = catalog.T("contact.company") + "最多 160 個字元。"
	case strings.HasPrefix(message, "Message must be"):
		localized = catalog.T("contact.message") + "最多 2000 個字元。"
	}
	if catalog.Locale.Code == "zh-CN" {
		return content.ToSimplified(localized)
	}
	return localized
}

func formatTime(value time.Time) string {
	if value.IsZero() {
		return ""
	}
	return value.UTC().Format("2006-01-02 15:04:05 UTC")
}

func methodNotAllowed(w http.ResponseWriter) {
	w.Header().Set("Allow", "GET, POST")
	http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
}

func parseAdminLeadFilters(r *http.Request) adminLeadFilters {
	email := strings.TrimSpace(r.URL.Query().Get("email"))
	company := strings.TrimSpace(r.URL.Query().Get("company"))
	interest := strings.TrimSpace(r.URL.Query().Get("interest"))
	return adminLeadFilters{
		Email:     email,
		Company:   company,
		Interest:  interest,
		Token:     strings.TrimSpace(r.URL.Query().Get("token")),
		HasActive: email != "" || company != "" || interest != "",
	}
}

func (filters adminLeadFilters) listFilter() leads.ListFilter {
	return leads.ListFilter{
		Email:    filters.Email,
		Company:  filters.Company,
		Interest: filters.Interest,
	}
}

func parseAdminLeadPage(raw string) int {
	page, err := strconv.Atoi(strings.TrimSpace(raw))
	if err != nil || page < 1 {
		return 1
	}
	return page
}

func adminLeadTotalPages(totalCount, pageSize int) int {
	if totalCount == 0 || pageSize <= 0 {
		return 0
	}
	return (totalCount + pageSize - 1) / pageSize
}

func adminLeadQueryValues(filters adminLeadFilters, page int) url.Values {
	values := url.Values{}
	if filters.Token != "" {
		values.Set("token", filters.Token)
	}
	if filters.Email != "" {
		values.Set("email", filters.Email)
	}
	if filters.Company != "" {
		values.Set("company", filters.Company)
	}
	if filters.Interest != "" {
		values.Set("interest", filters.Interest)
	}
	if page > 1 {
		values.Set("page", strconv.Itoa(page))
	}
	return values
}

func (s *Server) adminLeadsHref(token string, filters adminLeadFilters) string {
	values := adminLeadQueryValues(adminLeadFilters{
		Email:    filters.Email,
		Company:  filters.Company,
		Interest: filters.Interest,
		Token:    token,
	}, 0)
	if len(values) == 0 {
		return "/admin/leads"
	}
	return "/admin/leads?" + values.Encode()
}

func (s *Server) adminLeadPagination(filters adminLeadFilters, page, totalCount int) adminLeadPagination {
	pagination := adminLeadPagination{
		Page:       page,
		PageSize:   adminLeadPageSize,
		TotalCount: totalCount,
		TotalPages: adminLeadTotalPages(totalCount, adminLeadPageSize),
	}
	if totalCount == 0 {
		return pagination
	}

	pagination.Start = (page-1)*adminLeadPageSize + 1
	pagination.End = pagination.Start + adminLeadPageSize - 1
	if pagination.End > totalCount {
		pagination.End = totalCount
	}
	if page > 1 {
		pagination.PreviousHref = "/admin/leads?" + adminLeadQueryValues(filters, page-1).Encode()
	}
	if page < pagination.TotalPages {
		pagination.NextHref = "/admin/leads?" + adminLeadQueryValues(filters, page+1).Encode()
	}
	return pagination
}

func securityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		next.ServeHTTP(w, r)
	})
}

func (s *Server) cacheHeaders(next http.Handler) http.Handler {
	if !s.enableCDNCacheHeaders {
		return next
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.HasPrefix(r.URL.Path, "/static/"):
			// Static cache headers are applied by staticHandler only when the file exists.
		case r.URL.Path == "/robots.txt" || r.URL.Path == "/sitemap.xml":
			w.Header().Set("Cache-Control", "public, max-age=300")
		case r.URL.Path == "/healthz" || strings.HasPrefix(r.URL.Path, "/admin/"):
			w.Header().Set("Cache-Control", "no-store")
		case r.Method == http.MethodPost && contentPublicPath(r.URL.Path) == "/contact":
			w.Header().Set("Cache-Control", "no-store")
		case r.Method == http.MethodGet && isPublicHTMLPath(r.URL.Path):
			w.Header().Set("Cache-Control", "no-store")
		}
		next.ServeHTTP(w, r)
	})
}

func (s *Server) searchIndexingHeaders(next http.Handler) http.Handler {
	if !s.disableSearchIndexing {
		return next
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Robots-Tag", "noindex, nofollow, noarchive")
		next.ServeHTTP(w, r)
	})
}

func contentPublicPath(requestPath string) string {
	_, publicPath, ok := content.LocaleFromPath(requestPath)
	if !ok {
		return requestPath
	}
	return publicPath
}

func isPublicHTMLPath(requestPath string) bool {
	if strings.HasPrefix(requestPath, "/static/") ||
		strings.HasPrefix(requestPath, "/admin/") ||
		requestPath == "/healthz" ||
		requestPath == "/robots.txt" ||
		requestPath == "/sitemap.xml" {
		return false
	}
	_, _, ok := content.LocaleFromPath(requestPath)
	return ok
}

func buildAssetVersions(staticDir string) map[string]string {
	versions := make(map[string]string)
	_ = filepath.WalkDir(staticDir, func(filePath string, entry os.DirEntry, err error) error {
		if err != nil || entry.IsDir() {
			return nil
		}
		relative, err := filepath.Rel(staticDir, filePath)
		if err != nil {
			return nil
		}
		version, err := fileHash(filePath)
		if err != nil {
			return nil
		}
		assetPath := "/static/" + filepath.ToSlash(relative)
		versions[assetPath] = version
		return nil
	})
	return versions
}

func fileHash(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()
	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil))[:12], nil
}

func cleanAssetPath(rawPath string) string {
	if rawPath == "" {
		return ""
	}
	parsed, err := url.Parse(rawPath)
	if err == nil && parsed.Path != "" {
		rawPath = parsed.Path
	}
	if !strings.HasPrefix(rawPath, "/static/") {
		return rawPath
	}
	return path.Clean(rawPath)
}

func staticPathExists(staticDir, requestPath string) bool {
	cleanPath := path.Clean("/" + strings.TrimLeft(requestPath, "/"))
	filePath := filepath.Join(staticDir, filepath.FromSlash(strings.TrimPrefix(cleanPath, "/")))
	info, err := os.Stat(filePath)
	return err == nil && !info.IsDir()
}

func normalizePublicBaseURL(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}
	return strings.TrimRight(value, "/")
}
