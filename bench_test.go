package picomatch

import (
	"testing"
)

// BenchmarkSimple measures basic wildcards.
func BenchmarkSimple(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = IsMatch("foo.js", "*.js", nil)
	}
}

// BenchmarkGlobstar measures recursive deep directory matching.
func BenchmarkGlobstar(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = IsMatch("foo/bar/baz/qux/test.js", "foo/**/*.js", nil)
	}
}

// BenchmarkBraces measures brace expansions.
func BenchmarkBraces(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = IsMatch("foo/bar/test.js", "foo/{bar,baz}/*.js", nil)
	}
}

// BenchmarkExtglob measures extglob parsing and matching.
func BenchmarkExtglob(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = IsMatch("foo/bar.js", "foo/@(bar|baz).js", nil)
	}
}

// BenchmarkPosixClass measures POSIX character class matching.
func BenchmarkPosixClass(b *testing.B) {
	opts := &Options{Posix: true}
	for i := 0; i < b.N; i++ {
		_, _ = IsMatch("foo/a123.js", "foo/[[:alnum:]]*.js", opts)
	}
}
