package picomatch

import (
	"errors"
	"regexp"
	"strings"
)

func shouldSkipDotfile(input string, pattern string, opts *Options) bool {
	if opts != nil && opts.Dot {
		return false
	}
	if input == "" || pattern == "" {
		return false
	}

	inputSegs := strings.Split(input, "/")
	patternSegs := strings.Split(pattern, "/")
	for i, seg := range inputSegs {
		if seg == "" || !strings.HasPrefix(seg, ".") {
			continue
		}
		if i >= len(patternSegs) {
			return true
		}
		pseg := patternSegs[i]
		if pseg == "" || pseg == "**" {
			continue
		}
		if strings.HasPrefix(pseg, ".") {
			continue
		}
		return true
	}
	return false
}

func Test(input string, regex *regexp.Regexp, options *Options, glob string, posix bool) (bool, []string, string, error) {
	if input == "" {
		return false, nil, "", nil
	}

	opts := options
	if opts == nil {
		opts = &Options{}
	}

	if shouldSkipDotfile(input, glob, opts) {
		return false, nil, input, nil
	}

	var format func(string) string
	if opts.Format != nil {
		format = opts.Format
	} else if posix {
		format = ToPosixSlashes
	}

	match := input == glob
	output := input
	if match && format != nil {
		output = format(input)
	}

	if !match {
		if format != nil {
			output = format(input)
		}
		match = output == glob
	}

	var found []string
	if !match || opts.Capture {
		if opts.MatchBase || opts.Basename {
			match = MatchBase(input, regex, options, posix)
		} else {
			found = regex.FindStringSubmatch(output)
			match = found != nil
		}
	}

	return match, found, output, nil
}

func MatchBase(input string, globOrRegex interface{}, options *Options, posix bool) bool {
	switch v := globOrRegex.(type) {
	case *regexp.Regexp:
		return v.MatchString(Basename(input, posix))
	case string:
		r, err := regexp.Compile(v)
		if err == nil {
			return r.MatchString(Basename(input, posix))
		}
		regex, err := MakeRe(v, options)
		if err != nil {
			return false
		}
		return regex.MatchString(Basename(input, posix))
	default:
		return false
	}
}

// IsMatch reports whether the provided input matches one of the supplied glob patterns.
func IsMatch(str string, patterns interface{}, options *Options) (bool, error) {
	opts := options
	if opts == nil {
		opts = &Options{}
	}

	switch p := patterns.(type) {
	case string:
		if shouldSkipDotfile(str, p, opts) {
			return false, nil
		}
		regex, err := MakeRe(p, opts)
		if err != nil {
			return false, err
		}
		return regex.MatchString(str), nil
	case []string:
		for _, pattern := range p {
			ok, err := IsMatch(str, pattern, opts)
			if err != nil {
				return false, err
			}
			if ok {
				return true, nil
			}
		}
		return false, nil
	default:
		return false, errors.New("patterns must be string or []string")
	}
}

// CompileRe builds a regexp.Regexp from parsed glob state.
func CompileRe(state ParseState, options *Options) (*regexp.Regexp, error) {
	opts := options
	if opts == nil {
		opts = &Options{}
	}

	prepend := "^"
	append := "$"
	if opts.Contains {
		prepend = ""
		append = ""
	}

	source := prepend + "(?:" + state.Output + ")" + append
	if state.Negated {
		source = `^(?!` + source + `).*$`
	}

	return ToRegex(source, options)
}

// ToRegex compiles a regex source string into a regexp.Regexp using the supplied options.
func ToRegex(source string, options *Options) (*regexp.Regexp, error) {
	opts := options
	if opts == nil {
		opts = &Options{}
	}

	if opts.Flags != "" {
		if strings.Contains(opts.Flags, "i") {
			source = `(?i:` + source + `)`
		}
	} else if opts.Nocase {
		source = `(?i:` + source + `)`
	}

	return regexp.Compile(source)
}

// MakeRe parses and compiles a glob pattern into a regexp.Regexp.
func MakeRe(input string, options *Options) (*regexp.Regexp, error) {
	if input == "" {
		return nil, errors.New("expected a non-empty string")
	}

	opts := options
	if opts == nil {
		opts = &Options{}
	}

	var output string
	if opts.Fastpaths {
		s, err := parseFastpaths(input, options)
		if err == nil && s != "" {
			output = s
		}
	}

	var state ParseState
	var err error
	if output == "" {
		state, err = Parse(input, options)
		if err != nil {
			return nil, err
		}
	} else {
		state = ParseState{Output: output, Negated: false}
	}

	return CompileRe(state, options)
}
