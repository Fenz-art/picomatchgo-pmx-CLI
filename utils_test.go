package picomatch

import "testing"

func TestHasRegexChars(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"abc", false},
		{"a*b", true},
		{"a?b", true},
		{"a+b", true},
		{"a.b", true},
		{"a[b", true},
		{"foo/bar", false},
	}

	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			got := HasRegexChars(tc.input)
			if got != tc.want {
				t.Errorf("HasRegexChars(%q) = %v, want %v", tc.input, got, tc.want)
			}
		})
	}
}

func TestIsRegexChar(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"*", true},
		{"?", true},
		{"+", true},
		{"a", false},
		{"ab", false},
		{"", false},
	}

	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			got := IsRegexChar(tc.input)
			if got != tc.want {
				t.Errorf("IsRegexChar(%q) = %v, want %v", tc.input, got, tc.want)
			}
		})
	}
}

func TestEscapeRegex(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"abc", "abc"},
		{"a*b", `a\*b`},
		{"a.b", `a\.b`},
		{"(foo)", `\(foo\)`},
	}

	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			got := EscapeRegex(tc.input)
			if got != tc.want {
				t.Errorf("EscapeRegex(%q) = %q, want %q", tc.input, got, tc.want)
			}
		})
	}
}

func TestToPosixSlashes(t *testing.T) {
	input := `foo\bar\baz`
	want := `foo/bar/baz`
	got := ToPosixSlashes(input)
	if got != want {
		t.Errorf("ToPosixSlashes(%q) = %q, want %q", input, got, want)
	}
}
