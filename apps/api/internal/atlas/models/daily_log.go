// FILE: apps/api/internal/atlas/models/daily_log.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Define minimal DailyLog record model for WAVE-04 cardio FK dependency when WAVE-03 is not deployed.
//   SCOPE: DailyLogRecord used by cardio entry repository for DailyLog auto-creation. Minimal subset — no versioning, no workout exercises/sets.
//   DEPENDS: apps/api/internal/atlas/models/date.go.
//   LINKS: M-API / V-M-API / WAVE-03 / WAVE-04.
//   ROLE: TYPES
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   DailyLogRecord - Minimal DB record for daily_logs table (ID, UserID, Date, Notes, CreatedAt, UpdatedAt).
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added minimal DailyLogRecord for WAVE-04 cardio FK support.
// END_CHANGE_SUMMARY

package models

type DailyLogRecord struct {
	ID        string
	UserID    string
	Date      Date
	Notes     *string
	CreatedAt string
	UpdatedAt string
}