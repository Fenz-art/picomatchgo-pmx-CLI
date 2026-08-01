package utils

import (
	"errors"
	"regexp"
	"strings"
)

type MatchResult struct {
	Glob    string
	State   ParseState
	Regex   *regexp.Regexp
	Posix   bool
	Input   string
	Output  string
	Match   []string
	IsMatch bool
}

func Test(input string, regex *regexp.Regexp, options *Options, glob string, posix bool) (bool, []string, string, error) {
	if input == "" {
		return false, nil, "", nil
	}

	opts := options
	if opts == nil {
		opts = &Options{}
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

func MatchBase(input, glob string, options *Options, posix bool) bool {
	r, err := regexp.Compile(glob)
	if err == nil {
		return r.MatchString(Basename(input, posix))
	}

	regex, err := MakeRe(glob, options)
	if err != nil {
		return false
	}
	return regex.MatchString(Basename(input, posix))
}

func IsMatch(str string, patterns interface{}, options *Options) (bool, error) {
	switch p := patterns.(type) {
	case string:
		regex, err := MakeRe(p, options)
		if err != nil {
			return false, err
		}
		return regex.MatchString(str), nil
	case []string:
		for _, pattern := range p {
			ok, err := IsMatch(str, pattern, options)
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
