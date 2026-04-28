package web

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	DefaultReadTimeout     = 5 * time.Second
	DefaultWriteTimeout    = 15 * time.Second
	DefaultIdleTimeout     = 60 * time.Second
	DefaultShutdownTimeout = 10 * time.Second
)

type HTTPServerConfig struct {
	Addr         string
	Handler      http.Handler
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

func NewHTTPServer(cfg HTTPServerConfig) *http.Server {
	readTimeout := cfg.ReadTimeout
	if readTimeout <= 0 {
		readTimeout = DefaultReadTimeout
	}

	writeTimeout := cfg.WriteTimeout
	if writeTimeout <= 0 {
		writeTimeout = DefaultWriteTimeout
	}

	idleTimeout := cfg.IdleTimeout
	if idleTimeout <= 0 {
		idleTimeout = DefaultIdleTimeout
	}

	return &http.Server{
		Addr:         cfg.Addr,
		Handler:      cfg.Handler,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		IdleTimeout:  idleTimeout,
	}
}

func LoggingMiddleware(logger *log.Logger) func(http.Handler) http.Handler {
	if logger == nil {
		logger = log.New(io.Discard, "", 0)
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			started := time.Now()
			recorder := &statusRecorder{
				ResponseWriter: w,
				status:         http.StatusOK,
			}

			next.ServeHTTP(recorder, r)

			logger.Printf("%s %s %d %s", r.Method, sanitizedRequestURI(r.URL), recorder.status, time.Since(started).Round(time.Millisecond))
		})
	}
}

func sanitizedRequestURI(requestURL *url.URL) string {
	if requestURL == nil {
		return ""
	}
	if requestURL.RawQuery == "" {
		return requestURL.RequestURI()
	}

	values, err := url.ParseQuery(requestURL.RawQuery)
	if err != nil {
		return requestURL.EscapedPath()
	}

	redacted := false
	for key := range values {
		if !isSensitiveQueryParam(key) {
			continue
		}
		values.Set(key, "REDACTED")
		redacted = true
	}
	if !redacted {
		return requestURL.RequestURI()
	}

	sanitized := *requestURL
	sanitized.RawQuery = values.Encode()
	return sanitized.RequestURI()
}

func isSensitiveQueryParam(key string) bool {
	switch strings.ToLower(key) {
	case "token", "access_token", "api_key", "apikey":
		return true
	default:
		return false
	}
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}
