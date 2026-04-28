package web

import (
	"context"
	"encoding/csv"
	"html/template"
	"net"
	"net/http"
	"net/url"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"realtek-connect/internal/docs"
	"realtek-connect/internal/features"
	"realtek-connect/internal/leads"
)

type LeadStore interface {
	Insert(context.Context, leads.Lead) error
	List(context.Context, int) ([]leads.LeadRecord, error)
}

type Config struct {
	TemplatesDir string
	StaticDir    string
	LeadStore    LeadStore
	AdminToken   string
}

type Server struct {
	templatesDir string
	staticDir    string
	leadStore    LeadStore
	adminToken   string
	contactLimit *submissionRateLimiter
}

type pageData struct {
	Title           string
	MetaDescription string
	CanonicalURL    string
	SocialImageURL  string
	SocialImageAlt  string
	MetaRobots      string
	CurrentPath     string
	Docs            []docs.Section
	Doc             docs.Section
	Features        []features.Feature
	Feature         features.Feature
	Form            contactForm
	Errors          map[string]string
	Success         bool
	SubmittedFor    string
	Leads           []leads.LeadRecord
	AdminEnabled    bool
	AdminCSVHref    string
}

type contactForm struct {
	Name     string
	Company  string
	Email    string
	Interest string
	Message  string
	Website  string
}

const (
	contactNameMaxLength     = 120
	contactCompanyMaxLength  = 160
	contactEmailMaxLength    = 254
	contactInterestMaxLength = 120
	contactMessageMaxLength  = 2000
)

func NewServer(cfg Config) (*Server, error) {
	if cfg.TemplatesDir == "" {
		cfg.TemplatesDir = "templates"
	}
	if cfg.StaticDir == "" {
		cfg.StaticDir = "static"
	}
	return &Server{
		templatesDir: cfg.TemplatesDir,
		staticDir:    cfg.StaticDir,
		leadStore:    cfg.LeadStore,
		adminToken:   cfg.AdminToken,
		contactLimit: newSubmissionRateLimiter(5, 10*time.Minute),
	}, nil
}

func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(s.staticDir))))
	mux.HandleFunc("/robots.txt", s.handleRobotsTxt)
	mux.HandleFunc("/sitemap.xml", s.handleSitemapXML)
	mux.HandleFunc("/", s.handleHome)
	mux.HandleFunc("/docs", s.handleDocs)
	mux.HandleFunc("/docs/", s.handleDocDetail)
	mux.HandleFunc("/features", s.handleFeatures)
	mux.HandleFunc("/features/", s.handleFeatureDetail)
	mux.HandleFunc("/contact", s.handleContact)
	mux.HandleFunc("/admin/leads", s.handleAdminLeads)
	mux.HandleFunc("/admin/leads.csv", s.handleAdminLeadsCSV)
	mux.HandleFunc("/healthz", s.handleHealthz)
	return securityHeaders(mux)
}

func (s *Server) handleHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	if r.Method != http.MethodGet {
		methodNotAllowed(w)
		return
	}
	s.render(w, http.StatusOK, "home.html", s.basePageData(
		r,
		"Realtek Connect+ | IoT Cloud Platform",
		"Realtek Connect+ is an IoT cloud platform for provisioning, OTA, fleet management, app SDKs, insights, private cloud, and integrations.",
	))
}

func (s *Server) handleDocs(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/docs" {
		http.NotFound(w, r)
		return
	}
	if r.Method != http.MethodGet {
		methodNotAllowed(w)
		return
	}
	s.render(w, http.StatusOK, "docs.html", s.basePageData(
		r,
		"Developer Docs | Realtek Connect+",
		"Browse Realtek Connect+ documentation entry points for product overview, development, APIs, SDKs, firmware, CLI, deployment, and release notes.",
	))
}

func (s *Server) handleDocDetail(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		methodNotAllowed(w)
		return
	}

	slug := strings.TrimPrefix(r.URL.Path, "/docs/")
	slug = strings.Trim(slug, "/")
	if slug == "" {
		http.Redirect(w, r, "/docs", http.StatusSeeOther)
		return
	}

	doc, ok := docs.BySlug(slug)
	if !ok {
		http.NotFound(w, r)
		return
	}

	data := s.basePageData(
		r,
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

func (s *Server) handleFeatures(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/features" {
		http.NotFound(w, r)
		return
	}
	if r.Method != http.MethodGet {
		methodNotAllowed(w)
		return
	}
	s.render(w, http.StatusOK, "features.html", s.basePageData(
		r,
		"Features | Realtek Connect+",
		"Explore provisioning, OTA, fleet management, app SDK, insights, private cloud, and ecosystem integrations for Realtek-based IoT products.",
	))
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

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()
	records, err := s.leadStore.List(ctx, 100)
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
	data.AdminCSVHref = s.adminCSVHref(r)
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

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()
	records, err := s.leadStore.List(ctx, 500)
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

func (s *Server) adminCSVHref(r *http.Request) string {
	token := r.URL.Query().Get("token")
	if token == "" {
		return "/admin/leads.csv"
	}
	return "/admin/leads.csv?token=" + url.QueryEscape(token)
}

func (s *Server) unauthorized(w http.ResponseWriter) {
	if s.adminToken == "" {
		http.Error(w, "404 page not found", http.StatusNotFound)
		return
	}
	http.Error(w, "unauthorized", http.StatusUnauthorized)
}

func (s *Server) handleFeatureDetail(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		methodNotAllowed(w)
		return
	}

	slug := strings.TrimPrefix(r.URL.Path, "/features/")
	slug = strings.Trim(slug, "/")
	if slug == "" {
		http.Redirect(w, r, "/features", http.StatusSeeOther)
		return
	}

	feature, ok := features.BySlug(slug)
	if !ok {
		http.NotFound(w, r)
		return
	}

	data := s.basePageData(
		r,
		feature.Title+" | Realtek Connect+",
		feature.Summary,
	)
	data.Feature = feature
	s.render(w, http.StatusOK, "feature.html", data)
}

func (s *Server) handleContact(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.render(w, http.StatusOK, "contact.html", s.basePageData(
			r,
			"Contact | Realtek Connect+",
			"Contact the Realtek Connect+ team about provisioning, OTA, fleet operations, app SDKs, insights, or private cloud evaluation.",
		))
	case http.MethodPost:
		s.submitContact(w, r)
	default:
		methodNotAllowed(w)
	}
}

func (s *Server) submitContact(w http.ResponseWriter, r *http.Request) {
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
		data := s.basePageData(
			r,
			"Contact | Realtek Connect+",
			"Contact the Realtek Connect+ team about provisioning, OTA, fleet operations, app SDKs, insights, or private cloud evaluation.",
		)
		data.Form = form
		data.Errors = map[string]string{
			"form": "Request could not be processed.",
		}
		s.render(w, http.StatusBadRequest, "contact.html", data)
		return
	}

	errors := validateContact(form)
	if len(errors) > 0 {
		data := s.basePageData(
			r,
			"Contact | Realtek Connect+",
			"Contact the Realtek Connect+ team about provisioning, OTA, fleet operations, app SDKs, insights, or private cloud evaluation.",
		)
		data.Form = form
		data.Errors = errors
		s.render(w, http.StatusBadRequest, "contact.html", data)
		return
	}

	if s.contactLimit != nil && !s.contactLimit.Allow(contactSubmissionKey(r)) {
		data := s.basePageData(
			r,
			"Contact | Realtek Connect+",
			"Contact the Realtek Connect+ team about provisioning, OTA, fleet operations, app SDKs, insights, or private cloud evaluation.",
		)
		data.Form = form
		data.Errors = map[string]string{
			"form": "Too many requests from this address. Please wait a few minutes and try again.",
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

	data := s.basePageData(
		r,
		"Contact | Realtek Connect+",
		"Contact the Realtek Connect+ team about provisioning, OTA, fleet operations, app SDKs, insights, or private cloud evaluation.",
	)
	data.Success = true
	data.SubmittedFor = form.Name
	s.render(w, http.StatusOK, "contact.html", data)
}

func validateContact(form contactForm) map[string]string {
	errors := map[string]string{}
	if form.Name == "" {
		errors["name"] = "Name is required."
	} else if len(form.Name) > contactNameMaxLength {
		errors["name"] = "Name must be 120 characters or fewer."
	}
	if form.Email == "" {
		errors["email"] = "Email is required."
	} else if !emailPattern.MatchString(form.Email) {
		errors["email"] = "Enter a valid email address."
	} else if len(form.Email) > contactEmailMaxLength {
		errors["email"] = "Email must be 254 characters or fewer."
	}
	if form.Interest == "" {
		errors["interest"] = "Select an area of interest."
	} else if len(form.Interest) > contactInterestMaxLength {
		errors["interest"] = "Interest must be 120 characters or fewer."
	}
	if len(form.Company) > contactCompanyMaxLength {
		errors["company"] = "Company must be 160 characters or fewer."
	}
	if len(form.Message) > contactMessageMaxLength {
		errors["message"] = "Message must be 2000 characters or fewer."
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
		"formatTime": formatTime,
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

func securityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		next.ServeHTTP(w, r)
	})
}
