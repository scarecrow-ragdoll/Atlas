// FILE: apps/api/internal/atlas/models/workout_graphql.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Define GraphQL transport-specific workout input models that preserve omitted versus explicit null fields.
//   SCOPE: WAVE-03 nullable update inputs for gqlgen only; excludes service-layer validation, mutation behavior, and resolver mapping.
//   DEPENDS: github.com/99designs/gqlgen/graphql.
//   LINKS: M-API / V-M-API / WAVE-03.
//   ROLE: TYPES
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   UpdateWorkoutExerciseGraphQLInput - GraphQL update input preserving position and notes field presence.
//   UpdateWorkoutSetGraphQLInput - GraphQL update input preserving set value and notes field presence.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added GraphQL omittable input models for WAVE-03 update operations.
// END_CHANGE_SUMMARY

package models

import "github.com/99designs/gqlgen/graphql"

type UpdateWorkoutExerciseGraphQLInput struct {
	Position graphql.Omittable[*int32]  `json:"position"`
	Notes    graphql.Omittable[*string] `json:"notes"`
}

type UpdateWorkoutSetGraphQLInput struct {
	SetNumber graphql.Omittable[*int32]   `json:"setNumber"`
	Weight    graphql.Omittable[*float64] `json:"weight"`
	Reps      graphql.Omittable[*int32]   `json:"reps"`
	RPE       graphql.Omittable[*float64] `json:"rpe"`
	RIR       graphql.Omittable[*int32]   `json:"rir"`
	Notes     graphql.Omittable[*string]  `json:"notes"`
}
