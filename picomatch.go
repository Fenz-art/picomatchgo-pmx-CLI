package picomatch

import (
	"regexp"

	"github.com/debayansamal/port-mortem-picomatch-go/options"
	"github.com/debayansamal/port-mortem-picomatch-go/utils"
)

// IsMatch returns true if the input string matches the given glob pattern(s).
// patterns can be a string or a []string of patterns.
func IsMatch(str string, patterns interface{}, opts *options.Options) (bool, error) {
	return utils.IsMatch(str, patterns, opts)
}

// MakeRe compiles a glob pattern into a Go regular expression.
func MakeRe(input string, opts *options.Options) (*regexp.Regexp, error) {
	return utils.MakeRe(input, opts)
}

// CompileRe compiles the ParseState output into a Go regular expression.
func CompileRe(state ParseState, opts *options.Options) (*regexp.Regexp, error) {
	return utils.CompileRe(state, opts)
}

// ToRegex compiles a raw regex source string into a Go regular expression,
// applying Nocase or Flags options as specified.
func ToRegex(source string, opts *options.Options) (*regexp.Regexp, error) {
	return utils.ToRegex(source, opts)
}

// Test matches the input string against a compiled regexp, returning matching info
// such as whether it matched, the matched string slices, and formatting outputs.
func Test(input string, regex *regexp.Regexp, opts *options.Options, glob string, posix bool) (bool, []string, string, error) {
	return utils.Test(input, regex, opts, glob, posix)
}

// MatchBase matches the basename of the input filepath against a glob pattern.
func MatchBase(input, glob string, opts *options.Options, posix bool) bool {
	return utils.MatchBase(input, glob, opts, posix)
}
