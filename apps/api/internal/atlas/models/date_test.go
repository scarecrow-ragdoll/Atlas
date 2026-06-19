// FILE: apps/api/internal/atlas/models/date_test.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Verify strict Atlas GraphQL Date scalar parsing and marshaling behavior.
//   SCOPE: Tests YYYY-MM-DD success, timestamp rejection, and GraphQL marshaler construction; excludes downstream schema bindings.
//   DEPENDS: apps/api/internal/atlas/models/date.go, github.com/99designs/gqlgen/graphql, github.com/stretchr/testify.
//   LINKS: M-API / V-M-API / WAVE-03.
//   ROLE: TEST
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   TestDate_UnmarshalStrictYYYYMMDD - Confirms strict YYYY-MM-DD input parses and formats unchanged.
//   TestDate_RejectsTimestamp - Confirms timestamp inputs are rejected.
//   TestDate_MarshalGQL - Confirms Date exposes a non-nil GraphQL marshaler.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added Date scalar model tests for WAVE-03.
// END_CHANGE_SUMMARY

package models

import (
	"testing"

	"github.com/99designs/gqlgen/graphql"
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

func TestDate_MarshalGQL(t *testing.T) {
	d := MustDate("2026-06-19")

	var out graphql.Marshaler = d.MarshalGQL()
	assert.NotNil(t, out)
}
