package web

import (
	"fmt"
	"net/http"
	"sort"
	"strings"

	"realtek-connect/internal/leads"
)

func (s *Server) handlePrometheusMetrics(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/metrics/prometheus" {
		http.NotFound(w, r)
		return
	}
	if r.Method != http.MethodGet {
		methodNotAllowed(w)
		return
	}

	w.Header().Set("Content-Type", "text/plain; version=0.0.4; charset=utf-8")

	var b strings.Builder
	b.WriteString("# HELP rtk_cloud_frontend_up Whether the Cloud Frontend app is serving metrics.\n")
	b.WriteString("# TYPE rtk_cloud_frontend_up gauge\n")
	b.WriteString("rtk_cloud_frontend_up 1\n")

	if s.leadStore != nil {
		count, err := s.leadStore.Count(r.Context(), leads.ListFilter{})
		if err != nil {
			b.WriteString("# HELP rtk_cloud_frontend_leads_query_error Whether querying lead metrics failed.\n")
			b.WriteString("# TYPE rtk_cloud_frontend_leads_query_error gauge\n")
			b.WriteString("rtk_cloud_frontend_leads_query_error 1\n")
		} else {
			b.WriteString("# HELP rtk_cloud_frontend_leads_total Total lead records visible to the frontend lead store.\n")
			b.WriteString("# TYPE rtk_cloud_frontend_leads_total gauge\n")
			_, _ = fmt.Fprintf(&b, "rtk_cloud_frontend_leads_total %d\n", count)
		}
	}

	acceptances, redirects, downloadErrors := s.sdkDownloadMetrics.snapshot()
	b.WriteString("# HELP rtk_cloud_frontend_sdk_download_acceptances_total Accepted SDK evaluation terms by version and package.\n")
	b.WriteString("# TYPE rtk_cloud_frontend_sdk_download_acceptances_total counter\n")
	b.WriteString("# HELP rtk_cloud_frontend_sdk_download_redirects_total Issued SDK presigned redirects by version and package.\n")
	b.WriteString("# TYPE rtk_cloud_frontend_sdk_download_redirects_total counter\n")
	keys := make([]sdkMetricKey, 0, len(acceptances)+len(redirects))
	seen := map[sdkMetricKey]bool{}
	for key := range acceptances {
		seen[key] = true
		keys = append(keys, key)
	}
	for key := range redirects {
		if !seen[key] {
			keys = append(keys, key)
		}
	}
	sort.Slice(keys, func(i, j int) bool {
		if keys[i].Version == keys[j].Version {
			return keys[i].Package < keys[j].Package
		}
		return keys[i].Version < keys[j].Version
	})
	for _, key := range keys {
		_, _ = fmt.Fprintf(&b, "rtk_cloud_frontend_sdk_download_acceptances_total{version=%q,package=%q} %d\n", key.Version, key.Package, acceptances[key])
		_, _ = fmt.Fprintf(&b, "rtk_cloud_frontend_sdk_download_redirects_total{version=%q,package=%q} %d\n", key.Version, key.Package, redirects[key])
	}
	b.WriteString("# HELP rtk_cloud_frontend_sdk_download_errors_total Failed SDK catalog, validation, persistence, or signing operations.\n")
	b.WriteString("# TYPE rtk_cloud_frontend_sdk_download_errors_total counter\n")
	_, _ = fmt.Fprintf(&b, "rtk_cloud_frontend_sdk_download_errors_total %d\n", downloadErrors)

	_, _ = w.Write([]byte(b.String()))
}
