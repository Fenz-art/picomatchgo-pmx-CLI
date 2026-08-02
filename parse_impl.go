package picomatch

import (
	"regexp"
	"strings"
)

func parseFastpaths(input string) (string, error) {
	if input == "" {
		return "", nil
	}

	if !strings.ContainsAny(input, "*?[]{}()!+") {
		return EscapeRegex(input), nil
	}

	return "", nil
}

func parsePatternFragment(pattern string, opts *Options) string {
	if pattern == "" {
		return ""
	}

	if strings.Contains(pattern, "{") {
		if braceStart := strings.IndexByte(pattern, '{'); braceStart >= 0 {
			if braceEnd := findMatchingBrace(pattern, braceStart); braceEnd > braceStart {
				inner := pattern[braceStart+1 : braceEnd]
				if strings.Contains(inner, ",") {
					prefix := pattern[:braceStart]
					suffix := pattern[braceEnd+1:]
					parts := strings.Split(inner, ",")
					alters := make([]string, 0, len(parts))
					for _, part := range parts {
						alters = append(alters, parsePatternFragment(prefix+part+suffix, opts))
					}
					return "(?:" + strings.Join(alters, "|") + ")"
				}
			}
		}
	}

	pattern = strings.ReplaceAll(pattern, "**/", "__GLOBSTAR_PATH__")
	pattern = strings.ReplaceAll(pattern, "**", "__GLOBSTAR_ANY__")

	segments := strings.Split(pattern, "/")
	parts := make([]string, 0, len(segments))
	for _, segment := range segments {
		parts = append(parts, parseSegment(segment, opts))
	}
	return strings.Join(parts, "/")
}

func parseSegment(segment string, opts *Options) string {
	if segment == "" {
		return ""
	}

	var b strings.Builder
	if shouldProtectDot(segment, opts) {
		if strings.HasPrefix(segment, "*") || strings.HasPrefix(segment, "?") {
			b.WriteString("[^.]")
		}
	}

	for i := 0; i < len(segment); {
		if strings.HasPrefix(segment[i:], "__GLOBSTAR_PATH__") {
			b.WriteString("(?:[^/]+/)*")
			i += len("__GLOBSTAR_PATH__")
			continue
		}
		if strings.HasPrefix(segment[i:], "__GLOBSTAR_ANY__") {
			b.WriteString("(?:[^/]+/)*")
			i += len("__GLOBSTAR_ANY__")
			continue
		}

		ch := segment[i]
		if ch == '\\' && i+1 < len(segment) {
			b.WriteString(regexp.QuoteMeta(string(segment[i+1])))
			i += 2
			continue
		}

		switch ch {
		case '*':
			b.WriteString("[^/]*")
			i++
		case '?':
			b.WriteString("[^/]")
			i++
		case '{':
			end := strings.IndexByte(segment[i+1:], '}')
			if end >= 0 {
				inner := segment[i+1 : i+1+end]
				parts := strings.Split(inner, ",")
				if len(parts) > 1 {
					b.WriteString("(?:")
					for idx, part := range parts {
						if idx > 0 {
							b.WriteString("|")
						}
						b.WriteString(parseSegment(part, opts))
					}
					b.WriteString(")")
					i += end + 1
					continue
				}
			}
			b.WriteString(regexp.QuoteMeta(string(ch)))
			i++
		case '[':
			if mapped, next, ok := parsePosixClass(segment, i); ok {
				b.WriteString(mapped)
				i = next
				continue
			}
			b.WriteString("[")
			i++
		default:
			b.WriteString(regexp.QuoteMeta(string(ch)))
			i++
		}
	}

	return b.String()
}

func shouldProtectDot(segment string, opts *Options) bool {
	if opts == nil || opts.Dot || segment == "" || strings.HasPrefix(segment, ".") || strings.HasPrefix(segment, "**") {
		return false
	}
	return strings.ContainsAny(segment, "*?[]{}()+!@")
}

func findMatchingBrace(pattern string, open int) int {
	depth := 0
	for i := open; i < len(pattern); i++ {
		switch pattern[i] {
		case '{':
			depth++
		case '}':
			depth--
			if depth == 0 {
				return i
			}
		}
	}
	return -1
}

func parsePosixClass(pattern string, start int) (string, int, bool) {
	if start >= len(pattern) || pattern[start] != '[' {
		return "", 0, false
	}

	if start+1 >= len(pattern) || pattern[start+1] != '[' {
		return "", 0, false
	}
	if start+2 >= len(pattern) || pattern[start+2] != ':' {
		return "", 0, false
	}

	end := strings.Index(pattern[start+3:], ":]")
	if end < 0 {
		return "", 0, false
	}

	name := pattern[start+3 : start+3+end]
	mapped, ok := posixClassToRegex(name)
	if !ok {
		return "", 0, false
	}

	next := start + 3 + end + 3
	return mapped, next, true
}

func posixClassToRegex(name string) (string, bool) {
	switch name {
	case "alnum":
		return "[A-Za-z0-9]", true
	case "alpha":
		return "[A-Za-z]", true
	case "ascii":
		return "[\\x00-\\x7F]", true
	case "blank":
		return "[ \\t]", true
	case "digit":
		return "[0-9]", true
	case "lower":
		return "[a-z]", true
	case "space":
		return "[ \\t\\r\\n\\v\\f]", true
	case "upper":
		return "[A-Z]", true
	case "word":
		return "[A-Za-z0-9_]", true
	case "xdigit":
		return "[A-Fa-f0-9]", true
	default:
		return "", false
	}
}

// Parse converts a glob pattern into a ParseState with a regex source fragment.
func Parse(input string, options *Options) (ParseState, error) {
	opts := options
	if opts == nil {
		opts = &Options{}
	}

	if input == "" {
		return ParseState{Input: input, Output: ""}, nil
	}

	negated := false
	if strings.HasPrefix(input, "!") {
		negated = true
		input = input[1:]
	}

	if opts.Windows {
		input = strings.ReplaceAll(input, `\`, "/")
	}

	output := parsePatternFragment(input, opts)
	return ParseState{Input: input, Output: output, Negated: negated}, nil
}
