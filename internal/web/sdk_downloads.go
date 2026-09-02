package web

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"realtek-connect/internal/analytics"
	"realtek-connect/internal/content"
)

const sdkSessionCookie = "rtk_sdk_session"

// SDK IDs use base64url, which can start with '-' or '_'. Keep accepting the
// existing cookie alphabet without applying the analytics first-character rule.
var sdkSessionIDPattern = regexp.MustCompile(`^[A-Za-z0-9_-][A-Za-z0-9._:-]{0,63}$`)

func (s *Server) handleSDKCatalogAPI(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/api/sdk/catalog" {
		http.NotFound(w, r)
		return
	}
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if s.sdkDownloads == nil {
		w.Header().Set("Cache-Control", "no-store")
		http.Error(w, "SDK catalog is unavailable", http.StatusServiceUnavailable)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()
	catalog, err := s.sdkDownloads.PublicCatalog(ctx)
	if err != nil {
		w.Header().Set("Cache-Control", "no-store")
		http.Error(w, "SDK catalog is unavailable", http.StatusServiceUnavailable)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "public, max-age=60")
	_ = json.NewEncoder(w).Encode(catalog)
}

func (s *Server) attachSDKCatalog(r *http.Request, data *pageData) {
	if data == nil || s.sdkDownloads == nil {
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()
	catalog, err := s.sdkDownloads.Catalog(ctx)
	if err != nil {
		data.SDKDownloadError = "SDK downloads are temporarily unavailable."
		return
	}
	data.SDKCatalog = catalog
	data.SDKDownloadsEnabled = true
}

func (s *Server) handleSDKTerms(w http.ResponseWriter, r *http.Request, locale content.Locale, publicPath string) {
	if r.Method != http.MethodGet {
		methodNotAllowed(w)
		return
	}
	data := s.basePageData(r, locale, publicPath, "SDK Evaluation Terms | Realtek Connect+", "Evaluation terms for public Realtek Connect+ SDK downloads.")
	if s.sdkDownloads == nil {
		s.sdkDownloadMetrics.recordError()
		data.SDKDownloadError = "SDK evaluation terms are temporarily unavailable."
	} else {
		ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
		defer cancel()
		terms, version, err := s.sdkDownloads.Terms(ctx)
		if err != nil {
			data.SDKDownloadError = "SDK evaluation terms are temporarily unavailable."
		} else {
			data.SDKTerms = terms
			data.SDKTermsVersion = version
		}
	}
	s.render(w, http.StatusOK, "sdk_terms.html", data)
}

func (s *Server) handleSDKDownload(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store")
	if r.Method != http.MethodPost {
		methodNotAllowed(w)
		return
	}
	if s.sdkDownloads == nil {
		http.Error(w, "SDK downloads are unavailable", http.StatusServiceUnavailable)
		return
	}
	if err := r.ParseForm(); err != nil {
		s.sdkDownloadMetrics.recordError()
		http.Error(w, "invalid download request", http.StatusBadRequest)
		return
	}
	if r.FormValue("accepted") != "true" {
		s.sdkDownloadMetrics.recordError()
		http.Error(w, "evaluation terms must be accepted", http.StatusBadRequest)
		return
	}
	version := strings.TrimSpace(r.FormValue("version"))
	packageSlug := strings.TrimSpace(r.FormValue("package"))
	termsVersion := strings.TrimSpace(r.FormValue("terms_version"))
	sessionID := sdkSessionID(w, r)
	if s.sdkDownloadLimit != nil && !s.sdkDownloadLimit.Allow(contactSubmissionKey(r)+":"+sessionID) {
		s.sdkDownloadMetrics.recordError()
		http.Error(w, "too many SDK download requests", http.StatusTooManyRequests)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()
	artifact, downloadURL, err := s.sdkDownloads.DownloadURL(ctx, version, packageSlug, termsVersion, s.sdkDownloadURLTTL)
	if err != nil {
		s.sdkDownloadMetrics.recordError()
		http.Error(w, "invalid or unavailable SDK download", http.StatusBadRequest)
		return
	}
	requestID := newOpaqueID()
	s.sdkDownloadMetrics.recordAcceptance(version, artifact.Slug)
	if s.analyticsStore != nil {
		if err := s.analyticsStore.InsertSDKDownloadAcceptance(ctx, analytics.SDKDownloadAcceptance{
			AcceptedAt: time.Now().UTC(), SessionID: sessionID, TermsVersion: termsVersion,
			SDKVersion: version, Package: artifact.Slug, RequestID: requestID,
		}); err != nil {
			s.sdkDownloadMetrics.recordError()
			http.Error(w, "could not record SDK terms acceptance", http.StatusInternalServerError)
			return
		}
	}
	w.Header().Set("X-Request-ID", requestID)
	s.sdkDownloadMetrics.recordRedirect(version, artifact.Slug)
	http.Redirect(w, r, downloadURL, http.StatusSeeOther)
}

func sdkSessionID(w http.ResponseWriter, r *http.Request) string {
	if cookie, err := r.Cookie(sdkSessionCookie); err == nil && sdkSessionIDPattern.MatchString(cookie.Value) {
		return cookie.Value
	}
	value := newOpaqueID()
	http.SetCookie(w, &http.Cookie{Name: sdkSessionCookie, Value: value, Path: "/manual/sdk", MaxAge: 86400, HttpOnly: true, Secure: r.TLS != nil, SameSite: http.SameSiteLaxMode})
	return value
}

func newOpaqueID() string {
	buffer := make([]byte, 18)
	if _, err := rand.Read(buffer); err == nil {
		return base64.RawURLEncoding.EncodeToString(buffer)
	}
	return fmt.Sprintf("sdk-%d", time.Now().UnixNano())
}

func formatBytes(size int64) string {
	const unit = int64(1024)
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	value := float64(size)
	units := []string{"KB", "MB", "GB"}
	for _, suffix := range units {
		value /= float64(unit)
		if value < float64(unit) || suffix == "GB" {
			return fmt.Sprintf("%.1f %s", value, suffix)
		}
	}
	return fmt.Sprintf("%d B", size)
}
