// FILE: apps/api/internal/handler/atlas_media_test.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Verify Atlas media scaffold HTTP handlers return 501 Not Implemented for all operations.
//   SCOPE: POST upload, GET download, DELETE delete; scaffold-only until WAVE-02+.
//   DEPENDS: internal/handler, httptest.
//   LINKS: M-API / V-M-API.
//   ROLE: TEST
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
//
// START_MODULE_MAP
//   atlas media tests - Prove upload, download, and delete endpoints return 501.
// END_MODULE_MAP

package handler_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"monorepo-template/apps/api/internal/handler"
)

func TestAtlasMediaUpload_ReturnsNotImplemented(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/v1/media/upload", nil)
	rec := httptest.NewRecorder()

	handler.AtlasMediaUpload()(rec, req)

	assert.Equal(t, http.StatusNotImplemented, rec.Code)
	assert.Contains(t, rec.Body.String(), "not implemented")
}

func TestAtlasMediaDownload_ReturnsNotImplemented(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/media/123", nil)
	rec := httptest.NewRecorder()

	handler.AtlasMediaDownload()(rec, req)

	assert.Equal(t, http.StatusNotImplemented, rec.Code)
	assert.Contains(t, rec.Body.String(), "not implemented")
}

func TestAtlasMediaDelete_ReturnsNotImplemented(t *testing.T) {
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/media/123", nil)
	rec := httptest.NewRecorder()

	handler.AtlasMediaDelete()(rec, req)

	assert.Equal(t, http.StatusNotImplemented, rec.Code)
	assert.Contains(t, rec.Body.String(), "not implemented")
}