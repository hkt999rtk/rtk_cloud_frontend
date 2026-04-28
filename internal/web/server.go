package web

import (
	"context"
	"html/template"
	"net/http"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"realtek-connect/internal/docs"
	"realtek-connect/internal/features"
	"realtek-connect/internal/leads"
)

type LeadStore interface {
	Insert(context.Context, leads.Lead) error
}

type Config struct {
	TemplatesDir string
	StaticDir    string
	LeadStore    LeadStore
}

type Server struct {
	templatesDir string
	staticDir    string
	leadStore    LeadStore
}

type pageData struct {
	Title        string
	CurrentPath  string
	Docs         []docs.Section
	Doc          docs.Section
	Features     []features.Feature
	Feature      features.Feature
	Form         contactForm
	Errors       map[string]string
	Success      bool
	SubmittedFor string
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
	}, nil
}

func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(s.staticDir))))
	mux.HandleFunc("/", s.handleHome)
	mux.HandleFunc("/docs", s.handleDocs)
	mux.HandleFunc("/docs/", s.handleDocDetail)
	mux.HandleFunc("/features", s.handleFeatures)
	mux.HandleFunc("/features/", s.handleFeatureDetail)
	mux.HandleFunc("/contact", s.handleContact)
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
		Docs:        docs.All(),
		Features:    features.All(),
	})
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
	s.render(w, http.StatusOK, "docs.html", pageData{
		Title:       "Developer Docs | Realtek Connect+",
		CurrentPath: r.URL.Path,
		Docs:        docs.All(),
		Features:    features.All(),
	})
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

	s.render(w, http.StatusOK, "doc.html", pageData{
		Title:       doc.Title + " | Realtek Connect+ Docs",
		CurrentPath: r.URL.Path,
		Docs:        docs.All(),
		Doc:         doc,
		Features:    features.All(),
	})
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
		Docs:        docs.All(),
		Features:    features.All(),
	})
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
		Docs:        docs.All(),
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
			Docs:        docs.All(),
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
			Docs:        docs.All(),
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
		Docs:         docs.All(),
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
	tmpl, err := template.ParseFiles(files...)
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
