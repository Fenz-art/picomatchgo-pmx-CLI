package test

import (
	"github.com/debayansamal/port-mortem-picomatch-go"
	"testing"
)

func TestDotfilePatterns(t *testing.T) {
	optsWithDot := &picomatch.Options{Dot: true}
	optsNoDot := &picomatch.Options{Dot: false}

	tests := []struct {
		pattern string
		input   string
		opts    *picomatch.Options
		want    bool
	}{
		{"*.js", ".foo.js", optsNoDot, false},
		{"*.js", ".foo.js", optsWithDot, true},
	}

	for _, tc := range tests {
		t.Run(tc.pattern+"_"+tc.input, func(t *testing.T) {
			got, err := picomatch.IsMatch(tc.input, tc.pattern, tc.opts)
			if err != nil {
				t.Fatalf("IsMatch error: %v", err)
			}
			if got != tc.want {
				t.Errorf("IsMatch(%q, %q) = %v, want %v", tc.input, tc.pattern, got, tc.want)
			}
		})
	}
}
