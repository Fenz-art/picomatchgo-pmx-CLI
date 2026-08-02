package test

import (
	"testing"
	"github.com/debayansamal/port-mortem-picomatch-go"
)

func TestGlobstarPatterns(t *testing.T) {
	tests := []struct {
		pattern string
		input   string
		want    bool
	}{
		{"**/*.js", "a/b/c/d.js", true},
		{"foo/**/bar", "foo/a/b/c/bar", true},
		{"foo/**/*.js", "foo/bar.js", true},
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
