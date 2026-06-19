// FILE: apps/api/internal/handler/atlas_media.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Provide Atlas media REST scaffold handlers for WAVE-01. All operations return 501 Not Implemented.
//   SCOPE: POST /api/v1/media/upload, GET /api/v1/media/{id}, DELETE /api/v1/media/{id}; scaffold only, passes through Atlas guarded chain (PIN enforced when enabled).
//   DEPENDS: none (scaffold only).
//   LINKS: M-API / V-M-API.
//   ROLE: RUNTIME
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added Atlas media REST scaffold for WAVE-01.
// END_CHANGE_SUMMARY

package handler

import (
	"net/http"
)

func AtlasMediaUpload() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "not implemented", http.StatusNotImplemented)
	}
}

func AtlasMediaDownload() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "not implemented", http.StatusNotImplemented)
	}
}

func AtlasMediaDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "not implemented", http.StatusNotImplemented)
	}
}