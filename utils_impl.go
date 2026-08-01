package picomatch

import (
	"regexp"
	"strings"
)

// HasRegexChars returns true if the string contains any regex-special character.
func HasRegexChars(str string) bool {
	for _, r := range str {
		switch r {
		case '*', '?', '+', '.', '[', ']', '(', ')', '{', '}', '^', '$', '|', '\\':
			return true
		}
	}
	return false
}

// IsRegexChar returns true if the single-character string is a regex-special char.
func IsRegexChar(str string) bool {
	return len(str) == 1 && HasRegexChars(str)
}

// EscapeRegex escapes all regex-special characters in str.
func EscapeRegex(str string) string {
	var b strings.Builder
	for _, r := range str {
		switch r {
		case '*', '?', '+', '.', '[', ']', '(', ')', '{', '}', '^', '$', '|', '\\':
			b.WriteRune('\\')
		}
		b.WriteRune(r)
	}
	return b.String()
}

// ToPosixSlashes replaces backslashes with forward slashes for POSIX normalization.
func ToPosixSlashes(str string) string {
	var b strings.Builder
	for i := 0; i < len(str); i++ {
		if str[i] != '\\' {
			b.WriteByte(str[i])
			continue
		}
		if i+1 < len(str) {
			b.WriteByte('/')
		} else {
			b.WriteByte('/')
		}
	}
	return b.String()
}

// RemoveBackslashes strips standalone backslash escapes from a string.
// Backslashes inside character classes like [a\b] are preserved as-is.
func RemoveBackslashes(str string) string {
	var b strings.Builder
	for i := 0; i < len(str); i++ {
		if str[i] != '\\' {
			b.WriteByte(str[i])
			continue
		}
		if i+1 < len(str) {
			b.WriteByte(str[i+1])
			i++
		}
	}
	return b.String()
}

// EscapeLast inserts a backslash before the last unescaped occurrence of char
// in the input string, searching backwards from lastIdx.
func EscapeLast(input string, char string, lastIdx int) string {
	if lastIdx < 0 {
		lastIdx = len(input) - 1
	}
	if lastIdx < 0 || lastIdx >= len(input) {
		return input
	}
	idx := strings.LastIndex(input[:lastIdx+1], char)
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
func RemovePrefix(input string, state interface{}) string {
	output := input
	switch s := state.(type) {
	case *ScanState:
		if strings.HasPrefix(output, "./") {
			output = output[2:]
			s.Prefix = "./"
		}
	case *ParseState:
		if strings.HasPrefix(output, "./") {
			output = output[2:]
			s.Prefix = "./"
		}
	}
	return output
}

// WrapOutput wraps the regex source string in anchors (^ and $) and handles
// negation wrapping for inverted match patterns. Accepts either *ParseState or *ScanState.
func WrapOutput(input string, state interface{}, opts *Options) string {
	prepend := "^"
	appendStr := "$"
	if opts != nil && opts.Contains {
		prepend = ""
		appendStr = ""
	}

	output := prepend + "(?:" + input + ")" + appendStr
	var negated bool
	switch s := state.(type) {
	case *ScanState:
		negated = s.Negated
	case *ParseState:
		negated = s.Negated
	}
	if negated {
		return "(?:^(?!" + output + ").*$)"
	}
	return output
}

// Basename returns the last path segment of a filepath.
func Basename(path string, windows bool) string {
	var re *regexp.Regexp
	if windows {
		re = regexp.MustCompile(`[\\\\/]`)
	} else {
		re = regexp.MustCompile(`/`)
	}
	segs := re.Split(path, -1)
	if len(segs) == 0 {
		return ""
	}
	last := segs[len(segs)-1]
	if last == "" && len(segs) >= 2 {
		return segs[len(segs)-2]
	}
	return last
}
