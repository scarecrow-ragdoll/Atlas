//go:build tools

// FILE: apps/api/tools.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Track Go command dependencies used by API development and code generation.
//   SCOPE: Tool-only module imports for local and CI reproducibility; excludes runtime API dependencies.
//   DEPENDS: github.com/99designs/gqlgen, github.com/pressly/goose/v3/cmd/goose, github.com/sqlc-dev/sqlc/cmd/sqlc.
//   LINKS: M-API / M-WORKSPACE / V-M-API / V-M-WORKSPACE.
//   ROLE: CONFIG
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   imports - Pins gqlgen, goose, and sqlc command packages through Go module resolution.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added sqlc command dependency for API SQL codegen.
// END_CHANGE_SUMMARY

package tools

import (
	_ "github.com/99designs/gqlgen"
	_ "github.com/pressly/goose/v3/cmd/goose"
	_ "github.com/sqlc-dev/sqlc/cmd/sqlc"
)
