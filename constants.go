package picomatch

import (
	"regexp"
)

// Character constants (ASCII code points)
const (
	CHAR_0 = 48 // '0'
	CHAR_9 = 57 // '9'

	CHAR_UPPERCASE_A = 65  // 'A'
	CHAR_LOWERCASE_A = 97  // 'a'
	CHAR_UPPERCASE_Z = 90  // 'Z'
	CHAR_LOWERCASE_Z = 122 // 'z'

	CHAR_LEFT_PARENTHESES  = 40 // '('
	CHAR_RIGHT_PARENTHESES = 41 // ')'
	CHAR_ASTERISK          = 42 // '*'

	CHAR_AMPERSAND                = 38    // '&'
	CHAR_AT                       = 64    // '@'
	CHAR_BACKWARD_SLASH           = 92    // '\\'
	CHAR_CARRIAGE_RETURN          = 13    // '\r'
	CHAR_CIRCUMFLEX_ACCENT        = 94    // '^'
	CHAR_COLON                    = 58    // ':'
	CHAR_COMMA                    = 44    // ','
	CHAR_DOT                      = 46    // '.'
	CHAR_DOUBLE_QUOTE             = 34    // '"'
	CHAR_EQUAL                    = 61    // '='
	CHAR_EXCLAMATION_MARK         = 33    // '!'
	CHAR_FORM_FEED                = 12    // '\f'
	CHAR_FORWARD_SLASH            = 47    // '/'
	CHAR_GRAVE_ACCENT             = 96    // '`'
	CHAR_HASH                     = 35    // '#'
	CHAR_HYPHEN_MINUS             = 45    // '-'
	CHAR_LEFT_ANGLE_BRACKET       = 60    // '<'
	CHAR_LEFT_CURLY_BRACE         = 123   // '{'
	CHAR_LEFT_SQUARE_BRACKET      = 91    // '['
	CHAR_LINE_FEED                = 10    // '\n'
	CHAR_NO_BREAK_SPACE           = 160   // '\u00A0'
	CHAR_PERCENT                  = 37    // '%'
	CHAR_PLUS                     = 43    // '+'
	CHAR_QUESTION_MARK            = 63    // '?'
	CHAR_RIGHT_ANGLE_BRACKET      = 62    // '>'
	CHAR_RIGHT_CURLY_BRACE        = 125   // '}'
	CHAR_RIGHT_SQUARE_BRACKET     = 93    // ']'
	CHAR_SEMICOLON                = 59    // ';'
	CHAR_SINGLE_QUOTE             = 39    // '\''
	CHAR_SPACE                    = 32    // ' '
	CHAR_TAB                      = 9     // '\t'
	CHAR_UNDERSCORE               = 95    // '_'
	CHAR_VERTICAL_LINE            = 124   // '|'
	CHAR_ZERO_WIDTH_NOBREAK_SPACE = 65279 // '\uFEFF'

	MAX_LENGTH                    = 1024 * 64
	DEFAULT_MAX_EXTGLOB_RECURSION = 0
	Version                       = "1.0.0"
)

// Posix glob regex fragments
const (
	DOT_LITERAL   = `\.`
	PLUS_LITERAL  = `\+`
	QMARK_LITERAL = `\?`
	SLASH_LITERAL = `\/`
	ONE_CHAR      = ``
	QMARK         = `[^/]`
	END_ANCHOR    = `(?:\/|$)`
	START_ANCHOR  = `(?:^|\/)`
	DOTS_SLASH    = `\.{1,2}(?:\/|$)`
	NO_DOT        = ``
	NO_DOTS       = ``
	NO_DOT_SLASH  = ``
	NO_DOTS_SLASH = ``
	QMARK_NO_DOT  = `[^./]`
	STAR          = `[^/]*?`
	SEP           = `/`
)

// GlobChars holds regex tokens for POSIX or Windows path handling.
type GlobChars struct {
	DotLiteral   string
	PlusLiteral  string
	QmarkLiteral string
	SlashLiteral string
	OneChar      string
	Qmark        string
	EndAnchor    string
	DotsSlash    string
	NoDot        string
	NoDots       string
	NoDotSlash   string
	NoDotsSlash  string
	QmarkNoDot   string
	Star         string
	StartAnchor  string
	Sep          string
}

// PosixChars is the token set for POSIX paths.
var PosixChars = GlobChars{
	DotLiteral:   DOT_LITERAL,
	PlusLiteral:  PLUS_LITERAL,
	QmarkLiteral: QMARK_LITERAL,
	SlashLiteral: SLASH_LITERAL,
	OneChar:      ONE_CHAR,
	Qmark:        QMARK,
	EndAnchor:    END_ANCHOR,
	DotsSlash:    DOTS_SLASH,
	NoDot:        NO_DOT,
	NoDots:       NO_DOTS,
	NoDotSlash:   NO_DOT_SLASH,
	NoDotsSlash:  NO_DOTS_SLASH,
	QmarkNoDot:   QMARK_NO_DOT,
	Star:         STAR,
	StartAnchor:  START_ANCHOR,
	Sep:          SEP,
}

// WindowsChars is the token set for Windows paths.
var WindowsChars = GlobChars{
	DotLiteral:   DOT_LITERAL,
	PlusLiteral:  PLUS_LITERAL,
	QmarkLiteral: QMARK_LITERAL,
	SlashLiteral: `[\\/]`,
	OneChar:      "",
	Qmark:        `[^\\/]`,
	EndAnchor:    `(?:[\\/]|$)`,
	DotsSlash:    `\.{1,2}(?:[\\/]|$)`,
	NoDot:        "",
	NoDots:       "",
	NoDotSlash:   "",
	NoDotsSlash:  "",
	QmarkNoDot:   `[^.\\/]`,
	Star:         `[^\\/]*?`,
	StartAnchor:  `(?:^|[\\/])`,
	Sep:          "\\",
}

// GetGlobChars returns the platform-specific glob character set for POSIX or
// Windows-style path separators.
func GetGlobChars(win32 bool) GlobChars {
	if win32 {
		return WindowsChars
	}
	return PosixChars
}

// PosixRegexSource maps POSIX character class names to regex character ranges.
var PosixRegexSource = map[string]string{
	"alnum":  `a-zA-Z0-9`,
	"alpha":  `a-zA-Z`,
	"ascii":  `\x00-\x7F`,
	"blank":  ` \t`,
	"cntrl":  `\x00-\x1F\x7F`,
	"digit":  `0-9`,
	"graph":  `\x21-\x7E`,
	"lower":  `a-z`,
	"print":  `\x20-\x7E `,
	"punct":  `\-!"#$%&'()\*+,./:;<=>?@[\]^_` + "`{|}~",
	"space":  ` \t\r\n\v\f`,
	"upper":  `A-Z`,
	"word":   `A-Za-z0-9_`,
	"xdigit": `A-Fa-f0-9`,
}

// Compiled regexes for utils and scanner
var (
	REGEX_BACKSLASH             = regexp.MustCompile(`\\`)
	REGEX_NON_SPECIAL_CHARS     = regexp.MustCompile(`^[^@!\\[\].,$*+?^{}(|\\/]+`)
	REGEX_SPECIAL_CHARS         = regexp.MustCompile(`[-*+?.^${}(|)\\[\\]]`)
	REGEX_SPECIAL_CHARS_BACKREF = regexp.MustCompile(`(\\?)(\W)`)
	REGEX_SPECIAL_CHARS_GLOBAL  = regexp.MustCompile(`([-*+?.^${}(|)\\[\\]])`)
	REGEX_REMOVE_BACKSLASH      = regexp.MustCompile(`\\.`)
)

// Replacements simplifies glob patterns.
var Replacements = map[string]string{
	"***":      "*",
	"**/**":    "**",
	"**/**/**": "**",
}

// ExtglobType describes an extglob operator and its regex delimiters.
type ExtglobType struct {
	Type  string
	Open  string
	Close string
}

// ExtglobChars returns the extglob character map for the supplied glob
// character set.
func ExtglobChars(chars GlobChars) map[byte]ExtglobType {
	return map[byte]ExtglobType{
		'!': {Type: "negate", Open: "(?:", Close: ")" + chars.Star + ")"},
		'?': {Type: "qmark", Open: "(?:", Close: ")?"},
		'+': {Type: "plus", Open: "(?:", Close: ")+"},
		'*': {Type: "star", Open: "(?:", Close: ")*"},
		'@': {Type: "at", Open: "(?:", Close: ")"},
	}
}
