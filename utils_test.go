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

func TestBasename(t *testing.T) {
	tests := []struct {
		path    string
		windows bool
		want    string
	}{
		{"a/b/c", false, "c"},
		{"a/b/c/", false, "c"},
		{"/a/b/c.js", false, "c.js"},
		{"a\\b\\c", true, "c"},
		{"a\\b\\c\\", true, "c"},
		{"foo", false, "foo"},
	}

	for _, tc := range tests {
		t.Run(tc.path, func(t *testing.T) {
			got := Basename(tc.path, tc.windows)
			if got != tc.want {
				t.Errorf("Basename(%q, %v) = %q, want %q", tc.path, tc.windows, got, tc.want)
			}
		})
	}
}

func TestWrapOutput(t *testing.T) {
	tests := []struct {
		name  string
		input string
		state *ScanState
		opts  *Options
		want  string
	}{
		{
			name:  "basic no contains",
			input: "foo",
			state: &ScanState{},
			opts:  &Options{},
			want:  "^(?:foo)$",
		},
		{
			name:  "with contains",
			input: "foo",
			state: &ScanState{},
			opts:  &Options{Contains: true},
			want:  "(?:foo)",
		},
		{
			name:  "negated",
			input: "foo",
			state: &ScanState{Negated: true},
			opts:  &Options{},
			want:  "(?:^(?!^(?:foo)$).*$)",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := WrapOutput(tc.input, tc.state, tc.opts)
			if got != tc.want {
				t.Errorf("WrapOutput(%q) = %q, want %q", tc.input, got, tc.want)
			}
		})
	}
}

func TestRemovePrefix(t *testing.T) {
	state := &ScanState{}
	result := RemovePrefix("./foo/bar", state)
	if result != "foo/bar" {
		t.Errorf("RemovePrefix(\"./foo/bar\") = %q, want %q", result, "foo/bar")
	}
	if state.Prefix != "./" {
		t.Errorf("state.Prefix = %q, want %q", state.Prefix, "./")
	}

	state2 := &ScanState{}
	result2 := RemovePrefix("foo/bar", state2)
	if result2 != "foo/bar" {
		t.Errorf("RemovePrefix(\"foo/bar\") = %q, want %q", result2, "foo/bar")
	}
	if state2.Prefix != "" {
		t.Errorf("state.Prefix = %q, want %q", state2.Prefix, "")
	}
}

func TestEscapeLast(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		char    byte
		lastIdx int
		want    string
	}{
		{"escape last paren", "foo(bar)", ')', -1, `foo(bar\)`},
		{"no match", "foobar", ')', -1, "foobar"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := EscapeLast(tc.input, tc.char, tc.lastIdx)
			if got != tc.want {
				t.Errorf("EscapeLast(%q, %q, %d) = %q, want %q", tc.input, tc.char, tc.lastIdx, got, tc.want)
			}
		})
	}
}
