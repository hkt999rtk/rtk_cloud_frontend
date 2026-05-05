package web

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"realtek-connect/internal/analytics"
)

type adminAnalyticsView struct {
	Enabled            bool
	HasData            bool
	ConversionRate     string
	Metrics            []adminAnalyticsMetric
	TopReferrers       []adminAnalyticsRow
	ScrollDistribution []adminAnalyticsRow
	CtaByPage          []adminAnalyticsRow
}

type adminAnalyticsMetric struct {
	Label string
	Value string
	Hint  string
}

type adminAnalyticsRow struct {
	Primary   string
	Secondary string
	Value     string
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

	data := s.adminPageData(
		r,
		"Analytics | Realtek Connect+",
		"Protected aggregate analytics view for first-party website metrics.",
	)
	data.AdminEnabled = s.adminToken != ""
	data.AdminLeadsHref = s.adminLeadsHref(strings.TrimSpace(r.URL.Query().Get("token")), adminLeadFilters{})
	data.AdminAnalyticsHref = s.adminAnalyticsHref(strings.TrimSpace(r.URL.Query().Get("token")))

	if s.analyticsStore == nil {
		data.AdminAnalytics = adminAnalyticsView{
			Enabled: false,
			Metrics: defaultAnalyticsMetrics(),
		}
		s.render(w, http.StatusOK, "admin_analytics.html", data)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	summary, err := s.analyticsStore.Summary(ctx)
	if err != nil {
		http.Error(w, "could not load analytics summary", http.StatusInternalServerError)
		return
	}
	topReferrers, err := s.analyticsStore.TopReferrerOrigins(ctx, 5)
	if err != nil {
		http.Error(w, "could not load analytics summary", http.StatusInternalServerError)
		return
	}
	scrollDistribution, err := s.analyticsStore.ScrollDistribution(ctx)
	if err != nil {
		http.Error(w, "could not load analytics summary", http.StatusInternalServerError)
		return
	}
	ctaByPage, err := s.analyticsStore.CTAByPage(ctx, 8)
	if err != nil {
		http.Error(w, "could not load analytics summary", http.StatusInternalServerError)
		return
	}

	conversionRate := "0%"
	if summary.PageViews > 0 {
		conversionRate = fmt.Sprintf("%.1f%%", float64(summary.ClickCTAs)/float64(summary.PageViews)*100)
	}

	data.AdminAnalytics = adminAnalyticsView{
		Enabled:            true,
		HasData:            summary.PageViews > 0 || summary.ClickCTAs > 0 || summary.Scrolls > 0 || summary.Engaged > 0,
		ConversionRate:     conversionRate,
		Metrics:            analyticsMetrics(summary, conversionRate),
		TopReferrers:       renderReferrerRows(topReferrers),
		ScrollDistribution: renderScrollRows(scrollDistribution),
		CtaByPage:          renderCTARows(ctaByPage),
	}
	s.render(w, http.StatusOK, "admin_analytics.html", data)
}

func (s *Server) adminAnalyticsHref(token string) string {
	token = strings.TrimSpace(token)
	if token == "" {
		return "/admin/analytics"
	}
	values := url.Values{}
	values.Set("token", token)
	return "/admin/analytics?" + values.Encode()
}

func defaultAnalyticsMetrics() []adminAnalyticsMetric {
	return []adminAnalyticsMetric{
		{Label: "Page views", Value: "0", Hint: "Total tracked page_view events."},
		{Label: "CTA clicks", Value: "0", Hint: "Total tracked click_cta events."},
		{Label: "Conversion rate", Value: "0%", Hint: "click_cta events divided by page_view events."},
		{Label: "Scroll events", Value: "0", Hint: "Total tracked scroll events."},
	}
}

func analyticsMetrics(summary analytics.Summary, conversionRate string) []adminAnalyticsMetric {
	return []adminAnalyticsMetric{
		{Label: "Page views", Value: strconv.Itoa(summary.PageViews), Hint: "Total tracked page_view events."},
		{Label: "CTA clicks", Value: strconv.Itoa(summary.ClickCTAs), Hint: "Total tracked click_cta events."},
		{Label: "Conversion rate", Value: conversionRate, Hint: "click_cta events divided by page_view events."},
		{Label: "Scroll events", Value: strconv.Itoa(summary.Scrolls), Hint: "Total tracked scroll events."},
		{Label: "Engaged events", Value: strconv.Itoa(summary.Engaged), Hint: "Total tracked engaged events."},
	}
}

func renderReferrerRows(rows []analytics.ReferrerOriginCount) []adminAnalyticsRow {
	result := make([]adminAnalyticsRow, 0, len(rows))
	for _, row := range rows {
		label := row.ReferrerOrigin
		if label == "" {
			label = "(direct)"
		}
		result = append(result, adminAnalyticsRow{
			Primary: label,
			Value:   strconv.Itoa(row.Count),
		})
	}
	return result
}

func renderScrollRows(rows []analytics.ScrollDistribution) []adminAnalyticsRow {
	result := make([]adminAnalyticsRow, 0, len(rows))
	for _, row := range rows {
		result = append(result, adminAnalyticsRow{
			Primary: strconv.Itoa(row.Percent) + "%",
			Value:   strconv.Itoa(row.Count),
		})
	}
	return result
}

func renderCTARows(rows []analytics.CTAByPage) []adminAnalyticsRow {
	result := make([]adminAnalyticsRow, 0, len(rows))
	for _, row := range rows {
		result = append(result, adminAnalyticsRow{
			Primary:   row.Page,
			Secondary: row.CTA,
			Value:     strconv.Itoa(row.Count),
		})
	}
	return result
}
