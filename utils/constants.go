package utils

import (
	pm "github.com/debayansamal/port-mortem-picomatch-go"
	"regexp"
)

type Options = pm.Options

const (
	DEFAULT_MAX_EXTGLOB_RECURSION = 0
	MAX_LENGTH                    = 1024 * 64

	CHAR_0                          = 48
	CHAR_9                          = 57
	CHAR_UPPERCASE_A                = 65
	CHAR_LOWERCASE_A                = 97
	CHAR_UPPERCASE_Z                = 90
	CHAR_LOWERCASE_Z                = 122
	CHAR_LEFT_PARENTHESES           = 40
	CHAR_RIGHT_PARENTHESES          = 41
	CHAR_ASTERISK                   = 42
	CHAR_AMPERSAND                  = 38
	CHAR_AT                         = 64
	CHAR_BACKWARD_SLASH             = 92
	CHAR_CARRIAGE_RETURN            = 13
	CHAR_CIRCUMFLEX_ACCENT          = 94
	CHAR_COLON                      = 58
	CHAR_COMMA                      = 44
	CHAR_DOT                        = 46
	CHAR_DOUBLE_QUOTE               = 34
	CHAR_EQUAL                      = 61
	CHAR_EXCLAMATION_MARK           = 33
	CHAR_FORM_FEED                  = 12
	CHAR_FORWARD_SLASH              = 47
	CHAR_GRAVE_ACCENT               = 96
	CHAR_HASH                       = 35
	CHAR_HYPHEN_MINUS               = 45
	CHAR_LEFT_ANGLE_BRACKET         = 60
	CHAR_LEFT_CURLY_BRACE           = 123
	CHAR_LEFT_SQUARE_BRACKET        = 91
	CHAR_LINE_FEED                  = 10
	CHAR_NO_BREAK_SPACE             = 160
	CHAR_PERCENT                    = 37
	CHAR_PLUS                       = 43
	CHAR_QUESTION_MARK              = 63
	CHAR_RIGHT_ANGLE_BRACKET        = 62
	CHAR_RIGHT_CURLY_BRACE          = 125
	CHAR_RIGHT_SQUARE_BRACKET       = 93
	CHAR_SEMICOLON                  = 59
	CHAR_SINGLE_QUOTE               = 39
	CHAR_SPACE                      = 32
	CHAR_TAB                        = 9
	CHAR_UNDERSCORE                 = 95
	CHAR_VERTICAL_LINE              = 124
	CHAR_ZERO_WIDTH_NOBREAK_SPACE   = 65279
)

var (
	REGEX_BACKSLASH = regexp.MustCompile(`\\(?![*+?^${}(|)[\]])`)
	REGEX_NON_SPECIAL_CHARS = regexp.MustCompile(`^[^@![\].,$*+?^{}()|\\/]+`)
	REGEX_SPECIAL_CHARS = regexp.MustCompile(`[-*+?.^${}(|)[\]]`)
	REGEX_SPECIAL_CHARS_BACKREF = regexp.MustCompile(`(\\?)((\W)(\3*))`)
	REGEX_SPECIAL_CHARS_GLOBAL = regexp.MustCompile(`([-*+?.^${}(|)[\]])`)
	REGEX_REMOVE_BACKSLASH = regexp.MustCompile(`(?:\[.*?[^\\]\]|\\(?=.))`)

	REPLACEMENTS = map[string]string{
		"***":      "*",
		"**/**":    "**",
		"**/**/**": "**",
	}

	POSIX_REGEX_SOURCE = map[string]string{
		"alnum": "a-zA-Z0-9",
		"alpha": "a-zA-Z",
		"ascii": "\\x00-\\x7F",
		"blank": " \\t",
		"cntrl": "\\x00-\\x1F\\x7F",
		"digit": "0-9",
		"graph": "\\x21-\\x7E",
		"lower": "a-z",
		"print": "\\x20-\\x7E ",
		"punct": "\\-!\"#$%&'()\\*+,./:;<=>?@[\\]^_`{|}~",
		"space": " \\t\\r\\n\\v\\f",
		"upper": "A-Z",
		"word": "A-Za-z0-9_",
		"xdigit": "A-Fa-f0-9",
	}
)

type GlobChars struct {
	DOT_LITERAL   string
	PLUS_LITERAL  string
	QMARK_LITERAL string
	SLASH_LITERAL string
	ONE_CHAR      string
	QMARK         string
	END_ANCHOR    string
	DOTS_SLASH    string
	NO_DOT        string
	NO_DOTS       string
	NO_DOT_SLASH  string
	NO_DOTS_SLASH string
	QMARK_NO_DOT  string
	STAR          string
	START_ANCHOR  string
	SEP           string
}

type ExtglobChar struct {
	Type  string
	Open  string
	Close string
}

func ExtglobChars(chars GlobChars) map[string]ExtglobChar {
	return map[string]ExtglobChar{
		"!": {Type: "negate", Open: `(?:(?!(?:`, Close: `))` + chars.STAR + `)`},
		"?": {Type: "qmark", Open: `(?:`, Close: `)?`},
		"+": {Type: "plus", Open: `(?:`, Close: `)+`},
		"*": {Type: "star", Open: `(?:`, Close: `)*`},
		"@": {Type: "at", Open: `(?:`, Close: `)`},
	}
}

func GlobChars(win32 bool) GlobChars {
	if win32 {
		return windowsGlobChars()
	}
	return posixGlobChars()
}

func posixGlobChars() GlobChars {
	return GlobChars{
		DOT_LITERAL:   `\.`,
		PLUS_LITERAL:  `\+`,
		QMARK_LITERAL: `\?`,
		SLASH_LITERAL: `\/`,
		ONE_CHAR:      `(?=.)`,
		QMARK:         `[^/]`,
		END_ANCHOR:    `(?:\/|$)`,
		DOTS_SLASH:    `\.{1,2}(?:\/|$)`,
		NO_DOT:        `(?!\.)`,
		NO_DOTS:       `(?!(?:^|\/)(?:\.{1,2})(?:\/|$))`,
		NO_DOT_SLASH:  `(?!\.{0,1}(?:\/|$))`,
		NO_DOTS_SLASH: `(?!\.{1,2}(?:\/|$))`,
		QMARK_NO_DOT:  `[^./]`,
		STAR:          `[^/]*?`,
		START_ANCHOR:  `(?:^|\/)`,
		SEP:           `/`,
	}
}

func windowsGlobChars() GlobChars {
	const WIN_SLASH = `\\/`
	const WIN_NO_SLASH = `[^\\/]`

	return GlobChars{
		DOT_LITERAL:   `\.`,
		PLUS_LITERAL:  `\+`,
		QMARK_LITERAL: `\?`,
		SLASH_LITERAL: `[` + WIN_SLASH + `]`,
		ONE_CHAR:      `(?=.)`,
		QMARK:         WIN_NO_SLASH,
		END_ANCHOR:    `(?:[` + WIN_SLASH + `]|$)`,
		DOTS_SLASH:    `\.{1,2}(?:[` + WIN_SLASH + `]|$)`,
		NO_DOT:        `(?!\.)`,
		NO_DOTS:       `(?!(?:^|[` + WIN_SLASH + `])\.{1,2}(?:[` + WIN_SLASH + `]|$))`,
		NO_DOT_SLASH:  `(?!\.{0,1}(?:[` + WIN_SLASH + `]|$))`,
		NO_DOTS_SLASH: `(?!\.{1,2}(?:[` + WIN_SLASH + `]|$))`,
		QMARK_NO_DOT:  `[^.` + WIN_SLASH + `]`,
		STAR:          WIN_NO_SLASH + `*?`,
		START_ANCHOR:  `(?:^|[` + WIN_SLASH + `])`,
		SEP:           `\\`,
	}
}
