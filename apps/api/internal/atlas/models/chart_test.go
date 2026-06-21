package models

import (
	"math"
	"testing"
)

func TestCalculateE1RM(t *testing.T) {
	tests := []struct {
		name   string
		weight float64
		reps   float64
		want   float64
	}{
		{name: "standard epley", weight: 100, reps: 10, want: 133.33333333333334},
		{name: "zero weight", weight: 0, reps: 10, want: 0},
		{name: "zero reps", weight: 100, reps: 0, want: 0},
		{name: "negative weight", weight: -50, reps: 10, want: 0},
		{name: "negative reps", weight: 100, reps: -5, want: 0},
		{name: "single rep", weight: 200, reps: 1, want: 206.66666666666666},
		{name: "high reps", weight: 60, reps: 20, want: 100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalculateE1RM(tt.weight, tt.reps)
			if math.Abs(got-tt.want) > 0.01 {
				t.Errorf("CalculateE1RM(%v, %v) = %v, want ~%v", tt.weight, tt.reps, got, tt.want)
			}
		})
	}
}