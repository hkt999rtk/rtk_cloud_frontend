package web

import (
	"fmt"
	"net/http"
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

	_, _ = w.Write([]byte(b.String()))
}
