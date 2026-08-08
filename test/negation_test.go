package test

import (
	"testing"

	picomatch "github.com/Fenz-art/picomatchgo-pmx-CLI"
)

func TestNegationPatterns(t *testing.T) {
	tests := []struct {
		pattern string
		input   string
		want    bool
	}{
		{"!*.js", "foo.txt", true},
		{"!*.js", "foo.js", false},
		{"!foo/*", "bar/baz", true},
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
