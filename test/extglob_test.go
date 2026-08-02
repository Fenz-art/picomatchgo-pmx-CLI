package test

import (
	"github.com/debayansamal/port-mortem-picomatch-go"
	"testing"
)

func TestExtglobPatterns(t *testing.T) {
	tests := []struct {
		pattern string
		input   string
		want    bool
	}{
		{"@(a|b)", "a", true},
		{"@(a|b)", "c", false},
		{"+(a|b)", "aab", true},
		{"*(a|b)", "ab", true},
		{"?(a|b)", "a", true},
		{"!(a|b)", "c", true},
	}

	for _, tc := range tests {
		t.Run(tc.pattern+"_"+tc.input, func(t *testing.T) {
			got, err := picomatch.IsMatch(tc.input, tc.pattern, nil)
			if err != nil {
				t.Fatalf("IsMatch error: %v", err)
			}
			if got != tc.want {
				t.Errorf("IsMatch(%q, %q) = %v, want %v", tc.input, tc.pattern, got, tc.want)
			}
		})
	}
}
