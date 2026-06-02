package web

import (
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	cloudlogger "github.com/hkt999rtk/rtk_cloud_logger"
	"go.uber.org/zap"
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

func LoggingMiddleware(logger *zap.Logger) func(http.Handler) http.Handler {
	if logger == nil {
		logger = zap.NewNop()
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			started := time.Now()
			recorder := &statusRecorder{
				ResponseWriter: w,
				status:         http.StatusOK,
			}

			next.ServeHTTP(recorder, r)

			fields := []zap.Field{
				zap.String("method", r.Method),
				zap.String("path", sanitizedRequestURI(r.URL)),
				zap.Int("status", recorder.status),
				zap.Float64("duration_ms", float64(time.Since(started).Microseconds())/1000.0),
				zap.String("remote_addr", remoteAddr(r.RemoteAddr)),
			}
			if requestID := strings.TrimSpace(r.Header.Get("X-Request-Id")); requestID != "" {
				fields = append(fields, zap.String("request_id", requestID))
			}
			if traceID := strings.TrimSpace(r.Header.Get("X-Trace-Id")); traceID != "" {
				fields = append(fields, zap.String("trace_id", traceID))
			}

			logger.Info("http request", fields...)
		})
	}
}

func sanitizedRequestURI(requestURL *url.URL) string {
	if requestURL == nil {
		return ""
	}
	return cloudlogger.SanitizePath(requestURL.RequestURI())
}

func remoteAddr(addr string) string {
	host, _, err := net.SplitHostPort(addr)
	if err == nil {
		return host
	}
	return addr
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}
