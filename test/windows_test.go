package test

import (
	"testing"

	picomatch "github.com/Fenz-art/picomatchgo-pmx-CLI"
)

func TestWindowsPathNormalization(t *testing.T) {
	opts := &picomatch.Options{Windows: true}

	tests := []struct {
		pattern string
		input   string
		want    bool
	}{
		{"foo/*", "foo\\bar", true},
		{"foo/*/*.js", "foo\\bar\\baz.js", true},
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
