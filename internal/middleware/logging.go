package middleware

import (
	"net"
	"net/http"
	"strings"
	"time"

	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"

	"fiozap/internal/logger"
)

type responseWriter struct {
	http.ResponseWriter
	status int
	size   int
}

func (rw *responseWriter) WriteHeader(status int) {
	rw.status = status
	rw.ResponseWriter.WriteHeader(status)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	size, err := rw.ResponseWriter.Write(b)
	rw.size += size
	return size, err
}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		rw := &responseWriter{
			ResponseWriter: w,
			status:         http.StatusOK,
		}

		next.ServeHTTP(rw, r)

		duration := time.Since(start)
		log := logger.Get()

		var event *zerolog.Event
		switch {
		case rw.status >= 500:
			event = log.Error()
		case rw.status >= 400:
			event = log.Warn()
		default:
			event = log.Info()
		}

		event.
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Int("status", rw.status).
			Int("bytes", rw.size).
			Dur("latency", duration).
			Str("ip", getClientIP(r)).
			Str("reqId", chiMiddleware.GetReqID(r.Context())).
			Msg("http")
	})
}

func getClientIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		if idx := strings.Index(xff, ","); idx != -1 {
			return strings.TrimSpace(xff[:idx])
		}
		return strings.TrimSpace(xff)
	}

	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}
