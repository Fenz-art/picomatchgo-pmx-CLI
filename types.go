package picomatch

// Shared types used across the package.

// Options configures the matching and scanning behavior of picomatch.
type Options struct {
	// Baseline matching flags
	Dot             bool     // Match dotfiles when glob does not explicitly start with a dot.
	NoGlobstar      bool     // Disable ** matching across directories.
	Posix           bool     // Enable POSIX character classes like [:alnum:].
	Windows         bool     // Enable Windows path separator (\) support.
	StrictSlashes   bool     // Enforce strict slash matching.
	Ignore          []string // Patterns to ignore.
	MaxLength       int      // Max length of glob pattern (default 65536).
	MaxExtglobDepth int      // Max recursion depth for extglobs.

	// Callbacks and Formatters
	ExpandRange func(min, max string) string
	Format      func(input string) string
	OnMatch     func(match MatchResult)
	OnIgnore    func(match MatchResult)
	OnResult    func(result MatchResult)

	// Internal regex & matcher flags
	Contains  bool
	Flags     string
	Nocase    bool
	Nonegate  bool
	Fastpaths bool

	// Scan-related options
	Noext     bool // Disable extglob support during scan.
	Noparen   bool // Disable parentheses as glob characters.
	ScanToEnd bool // Continue scanning to the end of the pattern.
	Parts     bool // Return parts of the pattern in scan results.
	Tokens    bool // Return tokens in scan results.
	Unescape  bool // Remove backslashes from base and glob.

	// Parser/compiler options (present in upstream behavior)
	Prepend         string
	Capture         bool
	Bash            bool
	Noextglob       bool
	KeepQuotes      bool
	Nobracket       bool
	StrictBrackets  bool
	LiteralBrackets bool
	Nobrace         bool
	Regex           bool
	NoparenAlias    bool

	// Match-related options
	MatchBase     bool // Match the basename of the filepath.
	Basename      bool // Alias for MatchBase.
	CaptureGroups bool // Return capture groups.
	Debug         bool // Return errors on invalid regex.
}

// MatchResult contains metadata passed to callback functions.
type MatchResult struct {
	Glob    string
	Input   string
	Match   string
	Output  string
	Regex   string
	IsMatch bool
}

// ScanToken represents a single path segment token produced during glob scanning.
type ScanToken struct {
	Value       string
	Depth       int
	IsGlob      bool
	IsGlobstar  bool
	IsExtglob   bool
	IsBrace     bool
	IsBracket   bool
	IsPrefix    bool
	Negated     bool
	Backslashes bool
}

// ScanState is the result of scanning a glob pattern.
type ScanState struct {
	Prefix         string
	Input          string
	Start          int
	Base           string
	Glob           string
	IsBrace        bool
	IsBracket      bool
	IsGlob         bool
	IsExtglob      bool
	IsGlobstar     bool
	Negated        bool
	NegatedExtglob bool

	MaxDepth int
	Tokens   []ScanToken
	Slashes  []int
	Parts    []string
}

// ParseState is the intermediate state produced while translating a glob pattern
// into regex source.
type ParseState struct {
	Input          string
	Index          int
	Start          int
	Dot            bool
	Consumed       string
	Output         string
	Prefix         string
	Backtrack      bool
	Negated        bool
	NegatedExtglob bool
	Brackets       int
	Braces         int
	Parens         int
	Quotes         int
	Globstar       bool
	Tokens         []*ParseToken
}

// ParseToken represents a single parsed token from a glob pattern.
type ParseToken struct {
	Type        string
	Value       string
	Output      string
	Prev        *ParseToken
	Extglob     bool
	StartIndex  int
	TokensIndex int
	Parens      int
	Inner       string
	Conditions  int
	Suffix      string
	Backslashes bool
	Open        string
	Close       string
	Posix       bool
	Star        bool
	Dots        bool
	Comma       bool
	OutputIndex int
}

// ParseRepeatedExtglobMatch captures the repeated extglob match body and end
// index discovered during parsing.
type ParseRepeatedExtglobMatch struct {
	Type string
	Body string
	End  int
}

// RepeatedExtglobAnalysis describes whether a repeated extglob match is safe
// to emit as regex output.
type RepeatedExtglobAnalysis struct {
	Risky      bool
	SafeOutput string
}
