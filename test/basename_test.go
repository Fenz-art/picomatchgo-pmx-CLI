package test

import (
	"github.com/debayansamal/port-mortem-picomatch-go"
	"testing"
)

func TestBasenameMatching(t *testing.T) {
	opts := &picomatch.Options{MatchBase: true}

	tests := []struct {
		pattern string
		input   string
		want    bool
	}{
		{"*.js", "foo/bar/baz.js", true},
		{"*.js", "foo/bar/baz.txt", false},
	}

	for _, tc := range tests {
		t.Run(tc.pattern+"_"+tc.input, func(t *testing.T) {
			got, err := picomatch.IsMatch(tc.input, tc.pattern, opts)
			if err != nil {
				t.Fatalf("IsMatch error: %v", err)
			}
			if got != tc.want {
				t.Errorf("IsMatch(%q, %q) = %v, want %v", tc.input, tc.pattern, got, tc.want)
			}
		})
	}
}
