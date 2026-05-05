package web

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"realtek-connect/internal/analytics"
)

var (
	allowedAnalyticsEvents = map[string]struct{}{
		"page_view": {},
		"click_cta": {},
		"scroll":    {},
		"engaged":   {},
	}

	allowedAnalyticsPages = regexp.MustCompile(`^[a-z0-9][a-z0-9\/\-]{0,127}$`)
	allowedScrollPercent  = map[int]struct{}{
		25:  {},
		50:  {},
		75:  {},
		100: {},
	}
	allowedEngagedDuration = map[int]struct{}{
		10: {},
		30: {},
		60: {},
	}
	allowedAnalyticsCta = map[string]struct{}{
		"home_cta_primary":     {},
		"home_cta_secondary":   {},
		"home_cta_discuss":     {},
		"home_cta_band":        {},
		"home_feature_discuss": {},
		"feature_cta_primary":  {},
		"feature_cta_all":      {},
		"docs_cta_primary":     {},
		"docs_cta_secondary":   {},
		"doc_cta_primary":      {},
		"nav_cta_contact":      {},
		"contact_submit":       {},
	}
	allowedAnalyticsVariant = map[string]struct{}{
		"control": {},
	}
	sessionIDPattern = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._-]{0,63}$`)
)

type analyticsEventPayload struct {
	Event    string `json:"event"`
	Page     string `json:"page"`
	Cta      string `json:"cta"`
	Percent  int    `json:"percent"`
	Duration int    `json:"duration"`
	Variant  string `json:"variant"`
	Session  string `json:"session_id"`
}

func (s *Server) handleAnalyticsEvent(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/api/event" {
		http.NotFound(w, r)
		return
	}
	if r.Method != http.MethodPost {
		methodNotAllowed(w)
		return
	}
	if s.analyticsStore == nil {
		http.NotFound(w, r)
		return
	}
	if !isSameOrigin(r) {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	mediaType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil || mediaType != "application/json" {
		http.Error(w, "invalid content type", http.StatusBadRequest)
		return
	}

	if r.ContentLength > analyticsMaxPayloadBytes {
		http.Error(w, "payload too large", http.StatusBadRequest)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, analyticsMaxPayloadBytes)
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if s.analyticsLimiter != nil {
		if !s.analyticsLimiter.Allow(contactSubmissionKey(r)) {
			http.Error(w, "too many requests", http.StatusTooManyRequests)
			return
		}
	}

	var payload analyticsEventPayload
	if err := decoder.Decode(&payload); err != nil {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}
	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}

	event, err := normalizeAndValidateAnalyticsEvent(payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	event.ReferrerOrigin = sanitizedReferrerOrigin(r)
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	if err := s.analyticsStore.InsertEvent(ctx, event); err != nil {
		http.Error(w, "could not store analytics event", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusAccepted)
	_, _ = w.Write([]byte(`{"status":"ok"}`))
}

func (s *Server) handleAdminAnalytics(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/admin/analytics" {
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
	if s.analyticsStore == nil {
		http.NotFound(w, r)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	conversionRate, err := s.analyticsStore.ConversionRate(ctx)
	if err != nil {
		http.Error(w, "could not load analytics", http.StatusInternalServerError)
		return
	}
	refs, err := s.analyticsStore.TopReferrerOrigins(ctx, adminAnalyticsDefaultLimit)
	if err != nil {
		http.Error(w, "could not load analytics", http.StatusInternalServerError)
		return
	}
	scrolls, err := s.analyticsStore.ScrollDistribution(ctx)
	if err != nil {
		http.Error(w, "could not load analytics", http.StatusInternalServerError)
		return
	}
	ctaClicks, err := s.analyticsStore.CTAClicksByPage(ctx, adminAnalyticsDefaultLimit)
	if err != nil {
		http.Error(w, "could not load analytics", http.StatusInternalServerError)
		return
	}

	data := s.adminPageData(r, "Analytics | Realtek Connect+", "Protected analytics aggregate dashboard.")
	data.Analytics.Enabled = true
	data.Analytics.ConversionRate = conversionRate
	data.Analytics.TopReferrerOrigins = refs
	data.Analytics.ScrollDistribution = scrolls
	data.Analytics.CTAClicksByPage = ctaClicks

	s.render(w, http.StatusOK, "admin_analytics.html", data)
}

func normalizeAndValidateAnalyticsEvent(payload analyticsEventPayload) (analytics.Event, error) {
	if !isAllowedAnalyticsEvent(payload.Event) {
		return analytics.Event{}, fmt.Errorf("invalid event")
	}
	page := normalizeAnalyticsPage(payload.Page)
	if !isAllowedAnalyticsPage(page) {
		return analytics.Event{}, fmt.Errorf("invalid page")
	}
	if payload.Session == "" || !sessionIDPattern.MatchString(payload.Session) {
		return analytics.Event{}, fmt.Errorf("invalid session_id")
	}

	event := analytics.Event{
		Event:     payload.Event,
		Page:      page,
		SessionID: payload.Session,
	}

	if payload.Variant != "" {
		if !isAllowedAnalyticsVariant(payload.Variant) {
			return analytics.Event{}, fmt.Errorf("invalid variant")
		}
		event.Variant = payload.Variant
	}

	switch payload.Event {
	case "page_view":
		if payload.Percent != 0 || payload.Duration != 0 || payload.Cta != "" || payload.Variant != "" {
			return analytics.Event{}, fmt.Errorf("invalid page_view payload")
		}
	case "click_cta":
		if payload.Cta == "" || !isAllowedCTA(payload.Cta) {
			return analytics.Event{}, fmt.Errorf("invalid cta")
		}
		event.CTA = payload.Cta
	case "scroll":
		if payload.Percent == 0 || !isAllowedAnalyticsPercent(payload.Percent) {
			return analytics.Event{}, fmt.Errorf("invalid percent")
		}
		event.Percent = payload.Percent
	case "engaged":
		if payload.Duration == 0 || !isAllowedAnalyticsDuration(payload.Duration) {
			return analytics.Event{}, fmt.Errorf("invalid duration")
		}
		event.Duration = payload.Duration
	}

	return event, nil
}

func normalizeAnalyticsPage(value string) string {
	return strings.ToLower(strings.TrimSpace(strings.Trim(value, "/")))
}

func isAllowedAnalyticsEvent(value string) bool {
	_, ok := allowedAnalyticsEvents[value]
	return ok
}

func isAllowedAnalyticsPage(value string) bool {
	return allowedAnalyticsPages.MatchString(value)
}

func isAllowedCTA(value string) bool {
	_, ok := allowedAnalyticsCta[value]
	return ok
}

func isAllowedAnalyticsPercent(value int) bool {
	_, ok := allowedScrollPercent[value]
	return ok
}

func isAllowedAnalyticsDuration(value int) bool {
	_, ok := allowedEngagedDuration[value]
	return ok
}

func isAllowedAnalyticsVariant(value string) bool {
	_, ok := allowedAnalyticsVariant[value]
	return ok
}

func sanitizedReferrerOrigin(r *http.Request) string {
	referer := strings.TrimSpace(r.Header.Get("Referer"))
	if referer == "" {
		return ""
	}
	parsed, err := url.Parse(referer)
	if err != nil {
		return ""
	}
	if parsed.Scheme == "" || parsed.Host == "" {
		return ""
	}
	return parsed.Scheme + "://" + parsed.Host
}

func isSameOrigin(r *http.Request) bool {
	origin := strings.TrimSpace(r.Header.Get("Origin"))
	if origin == "" {
		return true
	}
	if origin == "null" {
		return false
	}
	originURL, err := url.Parse(origin)
	if err != nil || originURL.Scheme == "" || originURL.Host == "" {
		return false
	}

	host := strings.TrimSpace(r.Header.Get("X-Forwarded-Host"))
	if host == "" {
		host = r.Host
	}
	if host == "" {
		return false
	}

	expected := requestScheme(r) + "://" + host
	return originURL.String() == expected
}
