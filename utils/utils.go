package utils

import "strings"

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
