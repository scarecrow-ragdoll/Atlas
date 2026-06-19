package logger_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"

	"monorepo-template/libs/go/logger"
)

func TestRequestID_GeneratesUUID_WhenNoHeader(t *testing.T) {
	core, _ := observer.New(zap.DebugLevel)
	base := zap.New(core)

	var capturedID string
	handler := logger.RequestID(base)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedID = w.Header().Get("X-Request-ID")
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	assert.NotEmpty(t, capturedID)
	assert.Len(t, capturedID, 36) // UUID v4 format
	assert.Equal(t, capturedID, rec.Header().Get("X-Request-ID"))
}

func TestRequestID_PassesThrough_ExistingHeader(t *testing.T) {
	core, _ := observer.New(zap.DebugLevel)
	base := zap.New(core)

	var capturedID string
	handler := logger.RequestID(base)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedID = w.Header().Get("X-Request-ID")
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Request-ID", "existing-id-123")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	assert.Equal(t, "existing-id-123", capturedID)
	assert.Equal(t, "existing-id-123", rec.Header().Get("X-Request-ID"))
}

func TestRequestID_PutsLoggerWithRequestID_InContext(t *testing.T) {
	core, logs := observer.New(zap.DebugLevel)
	base := zap.New(core)

	handler := logger.RequestID(base)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l := logger.FromContext(r.Context())
		l.Info("test message")
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Request-ID", "test-req-id")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	require.Equal(t, 1, logs.Len())
	entry := logs.All()[0]
	assert.Equal(t, "test message", entry.Message)
	fields := entry.ContextMap()
	assert.Equal(t, "test-req-id", fields["request_id"])
}

func TestLogging_LogsRequestFields(t *testing.T) {
	core, logs := observer.New(zap.DebugLevel)
	base := zap.New(core)

	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
	})

	// Chain: RequestID -> Logging -> inner
	handler := logger.RequestID(base)(logger.Logging()(inner))

	req := httptest.NewRequest(http.MethodPost, "/graphql", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	require.Equal(t, 1, logs.Len())
	entry := logs.All()[0]
	assert.Equal(t, "request", entry.Message)

	fields := entry.ContextMap()
	assert.Equal(t, "POST", fields["method"])
	assert.Equal(t, "/graphql", fields["path"])
	assert.Equal(t, int64(201), fields["status"])
	assert.Contains(t, fields, "duration")
	assert.Contains(t, fields, "request_id")
	assert.Contains(t, fields, "remote_addr")
}

func TestLogging_DefaultStatus200(t *testing.T) {
	core, logs := observer.New(zap.DebugLevel)
	base := zap.New(core)

	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// no explicit WriteHeader — default 200
		_, _ = w.Write([]byte("ok"))
	})

	handler := logger.RequestID(base)(logger.Logging()(inner))

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	require.Equal(t, 1, logs.Len())
	fields := logs.All()[0].ContextMap()
	assert.Equal(t, int64(200), fields["status"])
}
