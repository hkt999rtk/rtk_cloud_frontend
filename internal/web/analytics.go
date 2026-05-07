package web

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"realtek-connect/internal/analytics"
	"realtek-connect/internal/docs"
	"realtek-connect/internal/features"
)

const analyticsRequestBodyLimit = 4 << 10

var (
	analyticsAllowedPages = func() map[string]struct{} {
		pages := map[string]struct{}{
			"home":     {},
			"features": {},
			"docs":     {},
			"contact":  {},
			"privacy":  {},
		}
		for _, feature := range features.All() {
			pages[feature.Slug] = struct{}{}
		}
		for _, section := range docs.All() {
			pages[section.Slug] = struct{}{}
		}
		return pages
	}()
	analyticsAllowedCTAs = map[string]struct{}{
		"contact_us":               {},
		"see_plans_limits":         {},
		"watch_brand_film":         {},
		"talk_to_sales":            {},
		"talk_to_platform_team":    {},
		"see_app_platform_context": {},
		"evaluate_feature":         {},
		"submit_request":           {},
		"review_features":          {},
		"apply_filters":            {},
		"clear_filters":            {},
		"export_csv":               {},
		"previous_page":            {},
		"next_page":                {},
		"open_docs":                {},
		"view_feature":             {},
		"contact_submit":           {},
		"doc_cta_primary":          {},
		"docs_cta_primary":         {},
		"feature_cta_all":          {},
		"feature_cta_primary":      {},
		"home_cta_band":            {},
		"home_cta_discuss":         {},
		"home_cta_primary":         {},
		"home_cta_secondary":       {},
		"home_feature_discuss":     {},
		"nav_cta_contact":          {},
	}
	analyticsAllowedVariants = map[string]struct{}{
		"control":       {},
		"default":       {},
		"experiment_a":  {},
		"experiment_b":  {},
		"variant_a":     {},
		"variant_b":     {},
		"home_hero_a":   {},
		"home_hero_b":   {},
		"contact_cta_a": {},
		"contact_cta_b": {},
	}
	analyticsSessionIDPattern = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._:-]{0,63}$`)
)

type analyticsEventPayload struct {
	Event     string `json:"event"`
	Page      string `json:"page"`
	CTA       string `json:"cta,omitempty"`
	Percent   *int   `json:"percent,omitempty"`
	Duration  *int   `json:"duration,omitempty"`
	Variant   string `json:"variant,omitempty"`
	SessionID string `json:"session_id"`
}

func (s *Server) handleAnalyticsEvent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store")

	if r.URL.Path != "/api/event" {
		http.NotFound(w, r)
		return
	}
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if s.analyticsStore == nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	mediaType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil || mediaType != "application/json" {
		http.Error(w, "content type must be application/json", http.StatusUnsupportedMediaType)
		return
	}

	event, err := s.parseAnalyticsEvent(r)
	if err != nil {
		status := http.StatusBadRequest
		if errors.Is(err, errAnalyticsBodyTooLarge) {
			status = http.StatusRequestEntityTooLarge
		}
		http.Error(w, err.Error(), status)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()
	if err := s.analyticsStore.InsertEvent(ctx, event); err != nil {
		http.Error(w, "could not save analytics event", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusAccepted)
	_, _ = w.Write([]byte(`{"status":"ok"}`))
}

var errAnalyticsBodyTooLarge = errors.New("analytics event body too large")

func (s *Server) parseAnalyticsEvent(r *http.Request) (analytics.Event, error) {
	defer r.Body.Close()

	limited := io.LimitedReader{R: r.Body, N: analyticsRequestBodyLimit + 1}
	body, err := io.ReadAll(&limited)
	if err != nil {
		return analytics.Event{}, fmt.Errorf("invalid analytics event payload")
	}
	if len(body) > analyticsRequestBodyLimit {
		return analytics.Event{}, errAnalyticsBodyTooLarge
	}

	decoder := json.NewDecoder(bytes.NewReader(body))
	var payload analyticsEventPayload
	if err := decoder.Decode(&payload); err != nil {
		return analytics.Event{}, fmt.Errorf("invalid analytics event payload")
	}

	if err := ensureNoTrailingJSON(decoder); err != nil {
		return analytics.Event{}, fmt.Errorf("invalid analytics event payload")
	}

	return s.validateAnalyticsEvent(r, payload)
}

func (s *Server) validateAnalyticsEvent(r *http.Request, payload analyticsEventPayload) (analytics.Event, error) {
	now := time.Now().UTC()

	eventType := strings.TrimSpace(payload.Event)
	if _, ok := allowedAnalyticsEventTypes[eventType]; !ok {
		return analytics.Event{}, fmt.Errorf("invalid analytics event type")
	}

	page := strings.TrimSpace(payload.Page)
	if _, ok := analyticsAllowedPages[page]; !ok {
		return analytics.Event{}, fmt.Errorf("invalid analytics page")
	}

	sessionID := strings.TrimSpace(payload.SessionID)
	if !analyticsSessionIDPattern.MatchString(sessionID) {
		return analytics.Event{}, fmt.Errorf("invalid analytics session id")
	}

	referrerOrigin := sanitizeReferrerOrigin(r.Header.Get("Referer"))

	event := analytics.Event{
		TS:             now.Unix(),
		Type:           eventType,
		Page:           page,
		ReferrerOrigin: referrerOrigin,
		SessionID:      sessionID,
		CreatedAt:      now,
	}

	variant := strings.TrimSpace(payload.Variant)
	if variant != "" {
		if _, ok := analyticsAllowedVariants[variant]; !ok {
			return analytics.Event{}, fmt.Errorf("invalid analytics variant")
		}
		event.Variant = variant
	}

	if event.Variant != "" && strings.ContainsAny(event.Variant, " \t\r\n") {
		return analytics.Event{}, fmt.Errorf("invalid analytics variant")
	}

	switch eventType {
	case "page_view":
		if payload.CTA != "" || payload.Percent != nil || payload.Duration != nil {
			return analytics.Event{}, fmt.Errorf("invalid analytics event fields")
		}
	case "click_cta":
		cta := strings.TrimSpace(payload.CTA)
		if _, ok := analyticsAllowedCTAs[cta]; !ok {
			return analytics.Event{}, fmt.Errorf("invalid analytics cta")
		}
		if payload.Percent != nil || payload.Duration != nil {
			return analytics.Event{}, fmt.Errorf("invalid analytics event fields")
		}
		event.CTA = cta
	case "scroll":
		if payload.CTA != "" || payload.Duration != nil {
			return analytics.Event{}, fmt.Errorf("invalid analytics event fields")
		}
		percent := payload.Percent
		if percent == nil {
			return analytics.Event{}, fmt.Errorf("invalid analytics percent")
		}
		if _, ok := allowedAnalyticsPercentages[*percent]; !ok {
			return analytics.Event{}, fmt.Errorf("invalid analytics percent")
		}
		event.Percent = percent
	case "engaged":
		if payload.CTA != "" || payload.Percent != nil {
			return analytics.Event{}, fmt.Errorf("invalid analytics event fields")
		}
		duration := payload.Duration
		if duration == nil {
			return analytics.Event{}, fmt.Errorf("invalid analytics duration")
		}
		if _, ok := allowedAnalyticsDurations[*duration]; !ok {
			return analytics.Event{}, fmt.Errorf("invalid analytics duration")
		}
		event.Duration = duration
	}

	return event, nil
}

var (
	allowedAnalyticsEventTypes = map[string]struct{}{
		"page_view": {},
		"click_cta": {},
		"scroll":    {},
		"engaged":   {},
	}
	allowedAnalyticsPercentages = map[int]struct{}{
		25:  {},
		50:  {},
		75:  {},
		100: {},
	}
	allowedAnalyticsDurations = map[int]struct{}{
		10: {},
		30: {},
		60: {},
	}
)

func ensureNoTrailingJSON(decoder *json.Decoder) error {
	var extra json.RawMessage
	if err := decoder.Decode(&extra); err != io.EOF {
		return err
	}
	return nil
}

func sanitizeReferrerOrigin(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return ""
	}

	parsed, err := url.Parse(raw)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return ""
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return ""
	}
	return parsed.Scheme + "://" + parsed.Host
}
