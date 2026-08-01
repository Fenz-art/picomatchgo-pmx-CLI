package picomatch

import (
	"regexp"
	"strings"
)

// HasRegexChars returns true if the string contains any regex-special character.
func HasRegexChars(str string) bool {
	return REGEX_SPECIAL_CHARS.MatchString(str)
}

// IsRegexChar returns true if the single-character string is a regex-special char.
func IsRegexChar(str string) bool {
	return len(str) == 1 && HasRegexChars(str)
}

// EscapeRegex escapes all regex-special characters in str.
func EscapeRegex(str string) string {
	return REGEX_SPECIAL_CHARS_GLOBAL.ReplaceAllString(str, `\$1`)
}

// ToPosixSlashes replaces backslashes with forward slashes for POSIX normalization.
func ToPosixSlashes(str string) string {
	return REGEX_BACKSLASH.ReplaceAllString(str, "/")
}

// RemoveBackslashes strips standalone backslash escapes from a string.
// Backslashes inside character classes like [a\b] are preserved as-is.
func RemoveBackslashes(str string) string {
	return REGEX_REMOVE_BACKSLASH.ReplaceAllStringFunc(str, func(match string) string {
		if match == `\` {
			return ""
		}
		return match
	})
}

// EscapeLast inserts a backslash before the last unescaped occurrence of char
// in the input string, searching backwards from lastIdx. If lastIdx is negative
// it searches the entire string.
func EscapeLast(input string, char byte, lastIdx int) string {
	if lastIdx < 0 {
		lastIdx = len(input) - 1
	}
	if lastIdx < 0 || lastIdx >= len(input) {
		return input
	}
	idx := strings.LastIndexByte(input[:lastIdx+1], char)
	if idx == -1 {
		return input
	}
	if idx > 0 && input[idx-1] == '\\' {
		return EscapeLast(input, char, idx-1)
	}
	return input[:idx] + `\` + input[idx:]
}

// RemovePrefix strips a leading "./" from the input string and records
// the prefix in the given ScanState.
func RemovePrefix(input string, state *ScanState) string {
	output := input
	if strings.HasPrefix(output, "./") {
		output = output[2:]
		state.Prefix = "./"
	}
	return output
}

// WrapOutput wraps the regex source string in anchors (^ and $) and handles
// negation wrapping for inverted match patterns.
func WrapOutput(input string, state *ScanState, opts *Options) string {
	prepend := "^"
	appendStr := "$"
	if opts != nil && opts.Contains {
		prepend = ""
		appendStr = ""
	}

	output := prepend + "(?:" + input + ")" + appendStr
	if state != nil && state.Negated {
		output = "(?:^(?!" + output + ").*$)"
	}
	return output
}

// Basename returns the last path segment of a filepath. If the last segment
// is empty (trailing slash), it returns the second-to-last segment instead.
func Basename(path string, windows bool) string {
	var re *regexp.Regexp
	if windows {
		re = regexp.MustCompile(`[\\/]`)
	} else {
		re = regexp.MustCompile(`/`)
	}
	segs := re.Split(path, -1)
	last := segs[len(segs)-1]
	if last == "" && len(segs) >= 2 {
		return segs[len(segs)-2]
	}
	return last
}
