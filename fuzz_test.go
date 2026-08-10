package picomatch

import (
	"testing"
)

// FuzzScan tests the fast-path scanner against random byte inputs.
func FuzzScan(f *testing.F) {
	seeds := []string{
		"*.js",
		"foo/**/bar/*.ts",
		"!(foo|bar)",
		"a/{b,c}/d",
		"path/[a-z]/file.*",
		"./!foo/bar/*.js",
	}

	for _, seed := range seeds {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, pattern string) {
		// Ensure Scan never panics or crashes on arbitrary input
		_ = Scan(pattern, nil)
		_ = Scan(pattern, &Options{Parts: true, Tokens: true, Unescape: true})
	})
}

// FuzzParse tests the AST parser against random byte inputs.
func FuzzParse(f *testing.F) {
	seeds := []string{
		"*.js",
		"foo/**/bar/*.ts",
		"!(foo|bar)",
		"a/{b,c}/d",
		"path/[a-z]/file.*",
	}

	for _, seed := range seeds {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, pattern string) {
		// Ensure Parse never panics or crashes on arbitrary input
		_, _ = Parse(pattern, nil)
	})
}

// FuzzIsMatch tests the full matcher pipeline against random pattern and input strings.
func FuzzIsMatch(f *testing.F) {
	f.Add("*.js", "foo.js")
	f.Add("foo/**/*.js", "foo/bar/baz.js")
	f.Add("a/{b,c}/d", "a/b/d")
	f.Add("!(foo|bar)", "baz")

	f.Fuzz(func(t *testing.T, pattern string, input string) {
		// Ensure IsMatch never panics on arbitrary inputs
		_, _ = IsMatch(input, pattern, nil)
	})
}
