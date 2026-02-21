package validation

import (
	"testing"

	"github.com/go-playground/validator/v10"
)

func TestIsValidMonthYear(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want bool
	}{
		{"valid 07-2025", "07-2025", true},
		{"valid 01-2024", "01-2024", true},
		{"valid 12-2030", "12-2030", true},
		{"empty", "", false},
		{"invalid month 00", "00-2025", false},
		{"invalid month 13", "13-2025", false},
		{"wrong format no dash", "072025", false},
		{"wrong format short", "7-2025", false},
		{"wrong format year short", "07-25", false},
		{"invalid chars", "ab-cdef", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidMonthYear(tt.s); got != tt.want {
				t.Errorf("IsValidMonthYear(%q) = %v, want %v", tt.s, got, tt.want)
			}
		})
	}
}

func TestMonthYear_Validator(t *testing.T) {
	v := newTestValidator(t)
	tests := []struct {
		name  string
		value string
		valid bool
	}{
		{"valid", "07-2025", true},
		{"invalid month", "13-2025", false},
		{"invalid format", "7-2025", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.Var(tt.value, "month_year")
			if tt.valid && err != nil {
				t.Errorf("expected valid, got %v", err)
			}
			if !tt.valid && err == nil {
				t.Errorf("expected validation error for %q", tt.value)
			}
		})
	}
}

func newTestValidator(t *testing.T) *validator.Validate {
	t.Helper()
	v := validator.New()
	if err := RegisterMonthYear(v); err != nil {
		t.Fatal(err)
	}
	return v
}
