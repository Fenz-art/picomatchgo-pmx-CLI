package test

import (
	"testing"
	"github.com/debayansamal/port-mortem-picomatch-go"
)

func TestPosixCharacterClasses(t *testing.T) {
	opts := &picomatch.Options{Posix: true}

	tests := []struct {
		pattern string
		input   string
		want    bool
	}{
		{"[[:alnum:]]", "a", true},
		{"[[:alnum:]]", "1", true},
		{"[[:digit:]]", "5", true},
		{"[[:digit:]]", "x", false},
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
