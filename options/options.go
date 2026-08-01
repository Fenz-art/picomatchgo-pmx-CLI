package options

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

	// Match-related options
	MatchBase bool // Match the basename of the filepath.
	Basename  bool // Alias for MatchBase.
	Capture   bool // Return capture groups.
	Debug     bool // Return errors on invalid regex.
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

// DefaultOptions returns a new Options struct initialized with default settings.
func DefaultOptions() Options {
	return Options{
		MaxLength:       65536,
		MaxExtglobDepth: 0,
		Fastpaths:       true,
	}
}
