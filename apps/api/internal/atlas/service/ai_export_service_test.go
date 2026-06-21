// FILE: apps/api/internal/atlas/service/ai_export_service_test.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Unit tests for BuildPrompt covering user context, persistent AI context, one-time comment, week flags, empty date range, no-data-in-period, all profile fields, and nil profile.
//   SCOPE: Pure function tests for prompt generation logic. Does not cover Generate/GetByID/List/Delete service methods (those require repo mocks and are covered by integration tests).
//   DEPENDS: apps/api/internal/atlas/service (BuildPrompt, UserProfileExport, SectionToggles).
//   LINKS: M-API / V-M-API / WAVE-07.
//   ROLE: TEST
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added BuildPrompt pure function tests for WAVE-07.
// END_CHANGE_SUMMARY

package service_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"monorepo-template/apps/api/internal/atlas/service"
)

func TestBuildPrompt_IncludesUserContext(t *testing.T) {
	goal := "Build muscle"
	profile := &service.UserProfileExport{Goal: &goal}
	prompt := service.BuildPrompt(profile, "2026-01-01", "2026-01-28", service.SectionToggles{}, nil, nil, "Some data")
	assert.Contains(t, prompt, "Build muscle")
	assert.Contains(t, prompt, "2026-01-01")
	assert.Contains(t, prompt, "2026-01-28")
	assert.Contains(t, prompt, "Analysis Requests")
}

func TestBuildPrompt_WithPersistentContext(t *testing.T) {
	ctx := "Focus on progressive overload and recovery"
	profile := &service.UserProfileExport{PersistentAiContext: &ctx}
	prompt := service.BuildPrompt(profile, "2026-01-01", "2026-01-28", service.SectionToggles{}, nil, nil, "")
	assert.Contains(t, prompt, "progressive overload")
}

func TestBuildPrompt_WithOneTimeComment(t *testing.T) {
	comment := "Trying a new deload protocol this month"
	profile := &service.UserProfileExport{}
	prompt := service.BuildPrompt(profile, "2026-01-01", "2026-01-28", service.SectionToggles{}, &comment, nil, "")
	assert.Contains(t, prompt, "deload protocol")
}

func TestBuildPrompt_WithWeekFlags(t *testing.T) {
	profile := &service.UserProfileExport{}
	flags := []string{"POOR_SLEEP", "HIGH_STRESS"}
	prompt := service.BuildPrompt(profile, "2026-01-01", "2026-01-28", service.SectionToggles{}, nil, flags, "")
	assert.Contains(t, prompt, "POOR_SLEEP")
	assert.Contains(t, prompt, "HIGH_STRESS")
}

func TestBuildPrompt_EmptyDateRange(t *testing.T) {
	profile := &service.UserProfileExport{}
	prompt := service.BuildPrompt(profile, "", "", service.SectionToggles{}, nil, nil, "")
	assert.NotEmpty(t, prompt)
	assert.Contains(t, prompt, "Analysis Requests")
}

func TestBuildPrompt_NoDataInPeriod(t *testing.T) {
	profile := &service.UserProfileExport{}
	prompt := service.BuildPrompt(profile, "2026-01-01", "2026-01-28", service.SectionToggles{}, nil, nil, "")
	assert.NotEmpty(t, prompt)
	assert.Contains(t, prompt, "Analysis Requests")
}

func TestBuildPrompt_WithAllProfileFields(t *testing.T) {
	goal := "Lose weight"
	height := 175.0
	birthDate := "1990-06-15"
	exp := "Advanced"
	split := "Upper/Lower"
	prog := "Linear"
	nutri := "Low Carb"
	aiCtx := "Focus on Zone 2 cardio"

	profile := &service.UserProfileExport{
		Goal:                      &goal,
		Height:                    &height,
		BirthDate:                 &birthDate,
		TrainingExperience:        &exp,
		CurrentTrainingSplit:      &split,
		PreferredProgressionStyle: &prog,
		NutritionStrategy:         &nutri,
		PersistentAiContext:       &aiCtx,
	}
	prompt := service.BuildPrompt(profile, "2026-01-01", "2026-01-28", service.SectionToggles{}, nil, nil, "")
	assert.Contains(t, prompt, "Lose weight")
	assert.Contains(t, prompt, "175.0")
	assert.Contains(t, prompt, "Advanced")
	assert.Contains(t, prompt, "Upper/Lower")
	assert.Contains(t, prompt, "Low Carb")
	assert.Contains(t, prompt, "Zone 2 cardio")
}

func TestBuildPrompt_NilProfile(t *testing.T) {
	prompt := service.BuildPrompt(nil, "2026-01-01", "2026-01-28", service.SectionToggles{}, nil, nil, "")
	assert.NotEmpty(t, prompt)
	assert.Contains(t, prompt, "Analysis Requests")
}