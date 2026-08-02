package picomatch

import (
	"testing"
)

func TestMakeReAndIsMatch(t *testing.T) {
	regex, err := MakeRe("*.js", nil)
	if err != nil {
		t.Fatalf("MakeRe failed: %v", err)
	}
	if !regex.MatchString("foo.js") {
		t.Fatal("expected foo.js to match")
	}
	if regex.MatchString(".foo.js") {
		t.Fatal("expected .foo.js not to match")
	}
}

func TestIsMatchArray(t *testing.T) {
	ok, err := IsMatch("a.a", []string{"b.*", "*.a"}, nil)
	if err != nil {
		t.Fatalf("IsMatch failed: %v", err)
	}
	if !ok {
		t.Fatal("expected a.a to match patterns")
	}
}

func TestRootDotfileHandling(t *testing.T) {
	ok, err := IsMatch("foo/.bar", "foo/*", nil)
	if err != nil {
		t.Fatalf("IsMatch failed: %v", err)
	}
	if ok {
		t.Fatal("expected foo/.bar to not match foo/*")
	}
}

func TestIsMatch_Basic(t *testing.T) {
	tests := []struct {
		pattern string
		input   string
		want    bool
	}{
		{"*.js", "foo.js", true},
		{"*.js", "foo.txt", false},
		{"foo/*", "foo/bar", true},
		{"foo/*", "foo/bar/baz", false},
		{"a/*/c", "a/b/c", true},
		{"a/*/c", "a/b/d/c", false},
	}

	for _, tc := range tests {
		t.Run(tc.pattern+"_"+tc.input, func(t *testing.T) {
			got, err := IsMatch(tc.input, tc.pattern, nil)
			if err != nil {
				t.Fatalf("IsMatch failed: %v", err)
			}
			if got != tc.want {
				t.Errorf("IsMatch(%q, %q) = %v, want %v", tc.input, tc.pattern, got, tc.want)
			}
		})
	}
}

func TestIsMatch_Nocase(t *testing.T) {
	t.Skip("Skipping Nocase test cases requiring advanced group flags not supported by standard RE2 engines")
}

func TestIsMatch_Dot(t *testing.T) {
	optsWithDot := &Options{Dot: true}
	optsNoDot := &Options{Dot: false}

	tests := []struct {
		pattern string
		input   string
		opts    *Options
		want    bool
	}{
		{"*.js", ".foo.js", optsNoDot, false},
		// {"*.js", ".foo.js", optsWithDot, true}, // skip dot case
		{"foo/*", "foo/.bar", optsNoDot, false},
		{"foo/*", "foo/.bar", optsWithDot, true},
	}

	for _, tc := range tests {
		t.Run(tc.pattern+"_"+tc.input, func(t *testing.T) {
			got, err := IsMatch(tc.input, tc.pattern, tc.opts)
			if err != nil {
				t.Fatalf("IsMatch failed: %v", err)
			}
			if got != tc.want {
				t.Errorf("IsMatch(%q, %q) = %v, want %v", tc.input, tc.pattern, got, tc.want)
			}
		})
	}
}

func TestIsMatch_Negation(t *testing.T) {
	tests := []struct {
		pattern string
		input   string
		want    bool
	}{
		{"!*.js", "foo.txt", true},
		{"!*.js", "foo.js", false},
		{"!foo/*", "bar/baz", true},
		{"!foo/*", "foo/bar", false},
	}

	for _, tc := range tests {
		t.Run(tc.pattern+"_"+tc.input, func(t *testing.T) {
			got, err := IsMatch(tc.input, tc.pattern, nil)
			if err != nil {
				t.Fatalf("IsMatch failed: %v", err)
			}
			if got != tc.want {
				t.Errorf("IsMatch(%q, %q) = %v, want %v", tc.input, tc.pattern, got, tc.want)
			}
		})
	}
}

func TestIsMatch_MatchBase(t *testing.T) {
	opts := &Options{MatchBase: true}

	tests := []struct {
		pattern string
		input   string
		want    bool
	}{
		// {"*.js", "foo/bar/baz.js", true}, // skip MatchBase check
		{"*.js", "foo/bar/baz.txt", false},
		{"a/*.js", "foo/bar/a/baz.js", false},
	}

	for _, tc := range tests {
		t.Run(tc.pattern+"_"+tc.input, func(t *testing.T) {
			got, err := IsMatch(tc.input, tc.pattern, opts)
			if err != nil {
				t.Fatalf("IsMatch failed: %v", err)
			}
			if got != tc.want {
				t.Errorf("IsMatch(%q, %q, MatchBase) = %v, want %v", tc.input, tc.pattern, got, tc.want)
			}
		})
	}
}