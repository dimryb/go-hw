package internalhttp

import (
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func loggingMiddleware(logger Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
			next.ServeHTTP(rw, r)

			clientIP := getClientIP(r)
			latency := time.Since(start)

			logEntry := fmt.Sprintf(
				"%s [%s] %s %s %s %d %d \"%s\"",
				clientIP,
				start.Format("02/Jan/2006:15:04:05 -0700"),
				r.Method,
				r.URL.Path,
				r.Proto,
				rw.statusCode,
				int(latency.Milliseconds()),
				r.UserAgent(),
			)

			logger.Info(logEntry)
		})
	}
}

func getClientIP(r *http.Request) string {
	ip := r.Header.Get("X-Forwarded-For")
	if ip != "" {
		ips := strings.Split(ip, ",")
		return strings.TrimSpace(ips[0])
	}
	ip, _, _ = net.SplitHostPort(r.RemoteAddr)
	return ip
}
