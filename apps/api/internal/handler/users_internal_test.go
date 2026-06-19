// FILE: apps/api/internal/handler/users_internal_test.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Cover package-internal users handler response encoding failure behavior.
//   SCOPE: Internal response writer edge cases; excludes public route behavior covered by users_test.go.
//   DEPENDS: apps/api/internal/handler.
//   LINKS: M-API / V-M-API.
//   ROLE: TEST
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   TestWriteJSON_WritesEncodeFailureBodyAfterInitialStatus - Covers encoder failure fallback behavior.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added internal handler coverage for response encoding failures.
// END_CHANGE_SUMMARY

package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

type unencodablePayload struct {
	Broken func()
}

func TestWriteJSON_WritesEncodeFailureBodyAfterInitialStatus(t *testing.T) {
	rec := httptest.NewRecorder()

	writeJSON(rec, http.StatusOK, unencodablePayload{})

	require.Equal(t, http.StatusOK, rec.Code)
	require.Contains(t, rec.Body.String(), "encode response")
}
