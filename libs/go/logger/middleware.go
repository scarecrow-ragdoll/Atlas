package logger

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

// RequestID middleware:
// 1. Extracts X-Request-ID from request header, or generates UUID v4
// 2. Sets X-Request-ID in response header
// 3. Creates base.With(zap.String("request_id", id)) and puts it in ctx via WithContext
func RequestID(base *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id := extractID(r.Header.Get(HeaderRequestID))
			if id == "" {
				id = generateID()
			}

			w.Header().Set(HeaderRequestID, id)

			l := base.With(zap.String("request_id", id))
			ctx := WithContext(r.Context(), l)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// Logging middleware logs each request with method, path, status, duration, remote_addr.
// Takes logger from ctx (set by RequestID middleware).
// Wraps http.ResponseWriter internally to capture the written status code.
// ORDERING: RequestID must be registered before Logging in the middleware chain.
// If no logger is found in ctx, falls back to zap.NewNop() (silent, no panic).
func Logging() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
			next.ServeHTTP(wrapped, r)

			l := FromContext(r.Context())
			l.Info("request",
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.Int("status", wrapped.statusCode),
				zap.Duration("duration", time.Since(start)),
				zap.String("remote_addr", r.RemoteAddr),
			)
		})
	}
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
