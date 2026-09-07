package web

import (
	"context"
	"encoding/json"
	"net/http"
	"realtek-connect/internal/analytics"
	"time"
)

func (s *Server) handlePRO2Examples(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store")
	if s.sdkDownloads == nil {
		http.Error(w, "Examples are unavailable", 503)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()
	switch r.URL.Path {
	case "/api/pro2-examples/catalog":
		if r.Method != "GET" {
			methodNotAllowed(w)
			return
		}
		c, e := s.sdkDownloads.Examples(ctx, r.URL.Query().Get("version"))
		if e != nil {
			http.Error(w, "Examples are unavailable", 503)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(c)
	case "/api/pro2-examples/download":
		if r.Method != "POST" {
			methodNotAllowed(w)
			return
		}
		r.Body = http.MaxBytesReader(w, r.Body, 4096)
		if e := r.ParseForm(); e != nil || r.FormValue("accepted") != "true" {
			http.Error(w, "Accept evaluation terms first", 400)
			return
		}
		session := sdkSessionIDForPath(w, r, "/api/pro2-examples")
		if s.sdkDownloadLimit != nil && !s.sdkDownloadLimit.Allow(contactSubmissionKey(r)+":"+session) {
			http.Error(w, "Too many requests", 429)
			return
		}
		version, id, terms := r.FormValue("version"), r.FormValue("artifact"), r.FormValue("terms_version")
		a, u, e := s.sdkDownloads.ExampleDownloadURL(ctx, version, id, terms, 10*time.Minute)
		if e != nil {
			http.Error(w, "Invalid or unavailable example download", 400)
			return
		}
		requestID := newOpaqueID()
		if s.analyticsStore != nil {
			e = s.analyticsStore.InsertSDKDownloadAcceptance(ctx, analytics.SDKDownloadAcceptance{AcceptedAt: time.Now().UTC(), SessionID: session, TermsVersion: terms, SDKVersion: version, Package: "pro2-" + id, RequestID: requestID})
			if e != nil {
				http.Error(w, "Could not record terms acceptance", 500)
				return
			}
		}
		s.sdkDownloadMetrics.recordAcceptance(version, "pro2-"+id)
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{"url": u, "artifact": a, "expires_at": time.Now().UTC().Add(10 * time.Minute)})
	default:
		http.NotFound(w, r)
	}
}
