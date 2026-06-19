// FILE: apps/api/internal/atlas/models/date_test.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Verify strict Atlas GraphQL Date scalar parsing and marshaling behavior.
//   SCOPE: Tests YYYY-MM-DD success, timestamp rejection, GraphQL quoted-date marshaling, and zero-value null marshaling; excludes downstream schema bindings.
//   DEPENDS: apps/api/internal/atlas/models/date.go, github.com/stretchr/testify.
//   LINKS: M-API / V-M-API / WAVE-03.
//   ROLE: TEST
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   TestDate_UnmarshalStrictYYYYMMDD - Confirms strict YYYY-MM-DD input parses and formats unchanged.
//   TestDate_RejectsTimestamp - Confirms timestamp inputs are rejected.
//   TestDate_MarshalGQLWritesQuotedDate - Confirms Date writes an exact quoted GraphQL date string.
//   TestDate_MarshalGQLZeroValueWritesNull - Confirms zero Date writes GraphQL null.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.1 - Updated Date scalar tests for gqlgen writer-based marshaling.
// END_CHANGE_SUMMARY

package models

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDate_UnmarshalStrictYYYYMMDD(t *testing.T) {
	var d Date

	require.NoError(t, d.UnmarshalGQL("2026-06-19"))
	assert.Equal(t, "2026-06-19", d.String())
}

func TestDate_RejectsTimestamp(t *testing.T) {
	var d Date

	require.Error(t, d.UnmarshalGQL("2026-06-19T10:00:00Z"))
}

func TestDate_MarshalGQLWritesQuotedDate(t *testing.T) {
	d := MustDate("2026-06-19")

	var buf bytes.Buffer
	d.MarshalGQL(&buf)

	assert.Equal(t, `"2026-06-19"`, buf.String())
}

func TestDate_MarshalGQLZeroValueWritesNull(t *testing.T) {
	var d Date

	var buf bytes.Buffer
	d.MarshalGQL(&buf)

	assert.Equal(t, "null", buf.String())
}
