package web

import (
	"context"
	"encoding/csv"
	"html/template"
	"net/http"
	"net/url"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

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
}

type pageData struct {
	Title        string
	CurrentPath  string
	Features     []features.Feature
	Feature      features.Feature
	Form         contactForm
	Errors       map[string]string
	Success      bool
	SubmittedFor string
	Leads        []leads.LeadRecord
	AdminEnabled bool
	AdminCSVHref string
}

type contactForm struct {
	Name     string
	Company  string
	Email    string
	Interest string
	Message  string
}

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
	}, nil
}

func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(s.staticDir))))
	mux.HandleFunc("/", s.handleHome)
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
	s.render(w, http.StatusOK, "home.html", pageData{
		Title:       "Realtek Connect+ | IoT Cloud Platform",
		CurrentPath: r.URL.Path,
		Features:    features.All(),
	})
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
	s.render(w, http.StatusOK, "features.html", pageData{
		Title:       "Features | Realtek Connect+",
		CurrentPath: r.URL.Path,
		Features:    features.All(),
	})
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

	s.render(w, http.StatusOK, "admin_leads.html", pageData{
		Title:        "Leads | Realtek Connect+",
		CurrentPath:  r.URL.Path,
		Features:     features.All(),
		Leads:        records,
		AdminEnabled: s.adminToken != "",
		AdminCSVHref: s.adminCSVHref(r),
	})
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

	s.render(w, http.StatusOK, "feature.html", pageData{
		Title:       feature.Title + " | Realtek Connect+",
		CurrentPath: r.URL.Path,
		Feature:     feature,
		Features:    features.All(),
	})
}

func (s *Server) handleContact(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.render(w, http.StatusOK, "contact.html", pageData{
			Title:       "Contact | Realtek Connect+",
			CurrentPath: r.URL.Path,
			Features:    features.All(),
		})
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
	}

	errors := validateContact(form)
	if len(errors) > 0 {
		s.render(w, http.StatusBadRequest, "contact.html", pageData{
			Title:       "Contact | Realtek Connect+",
			CurrentPath: r.URL.Path,
			Features:    features.All(),
			Form:        form,
			Errors:      errors,
		})
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

	s.render(w, http.StatusOK, "contact.html", pageData{
		Title:        "Contact | Realtek Connect+",
		CurrentPath:  r.URL.Path,
		Features:     features.All(),
		Success:      true,
		SubmittedFor: form.Name,
	})
}

func validateContact(form contactForm) map[string]string {
	errors := map[string]string{}
	if form.Name == "" {
		errors["name"] = "Name is required."
	}
	if form.Email == "" {
		errors["email"] = "Email is required."
	} else if !emailPattern.MatchString(form.Email) {
		errors["email"] = "Enter a valid email address."
	}
	if form.Interest == "" {
		errors["interest"] = "Select an area of interest."
	}
	return errors
}

var emailPattern = regexp.MustCompile(`^[^@\s]+@[^@\s]+\.[^@\s]+$`)

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
