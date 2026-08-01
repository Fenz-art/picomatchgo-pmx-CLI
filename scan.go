package picomatch

import (
	"math"
)

// ScanToken represents a single path segment token produced during glob scanning.
// Each slash-separated component of the glob becomes a token, carrying metadata
// about whether it contains glob syntax, its directory depth, etc.
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

// ScanState is the full result of scanning a glob pattern. It separates the
// non-glob base path from the glob portion and detects features like braces,
// brackets, extglobs, globstars, and negation patterns.
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

	// Populated when opts.Tokens or opts.Parts is true.
	MaxDepth int
	Tokens   []ScanToken
	Slashes  []int
	Parts    []string
}

// isPathSeparator returns true if code is a forward or backward slash.
func isPathSeparator(code byte) bool {
	return code == CHAR_FORWARD_SLASH || code == CHAR_BACKWARD_SLASH
}

// setDepth sets the directory depth on a token. Globstar tokens get
// MaxInt32 (representing Infinity in the JS source), others get 1.
func setDepth(token *ScanToken) {
	if !token.IsPrefix {
		if token.IsGlobstar {
			token.Depth = math.MaxInt32
		} else {
			token.Depth = 1
		}
	}
}

// charAt safely returns the byte at position i, or 0 if out of bounds.
func charAt(s string, i int) byte {
	if i < 0 || i >= len(s) {
		return 0
	}
	return s[i]
}

// Scan quickly scans a glob pattern and returns a ScanState containing useful
// metadata: the non-glob base path, the glob pattern, flags for braces,
// brackets, extglobs, globstars, negation, and optionally tokens and parts.
//
// This is a faithful port of lib/scan.js from micromatch/picomatch.
func Scan(input string, opts *Options) ScanState {
	if opts == nil {
		opts = &Options{}
	}

	length := len(input) - 1
	scanToEnd := opts.Parts || opts.ScanToEnd

	slashes := []int{}
	tokens := []ScanToken{}
	parts := []string{}

	str := input
	index := -1
	start := 0
	lastIndex := 0

	isBrace := false
	isBracket := false
	isGlob := false
	isExtglob := false
	isGlobstar := false
	braceEscaped := false
	backslashes := false
	negated := false
	negatedExtglob := false
	finished := false
	braces := 0

	var prev byte
	var code byte

	token := ScanToken{Value: "", Depth: 0, IsGlob: false}

	eos := func() bool {
		return index >= length
	}

	peek := func() byte {
		return charAt(str, index+1)
	}

	advance := func() byte {
		prev = code
		index++
		if index < len(str) {
			return str[index]
		}
		return 0
	}

	// Main scanning loop
	for index < length {
		code = advance()
		var next byte
		_ = next

		// Handle backslash escapes
		if code == CHAR_BACKWARD_SLASH {
			backslashes = true
			token.Backslashes = true
			code = advance()

			if code == CHAR_LEFT_CURLY_BRACE {
				braceEscaped = true
			}
			continue
		}

		// Handle braces: {a,b} and {1..5}
		if braceEscaped || code == CHAR_LEFT_CURLY_BRACE {
			braces++

			for !eos() {
				code = advance()
				if code == 0 {
					break
				}

				if code == CHAR_BACKWARD_SLASH {
					backslashes = true
					token.Backslashes = true
					advance()
					continue
				}

				if code == CHAR_LEFT_CURLY_BRACE {
					braces++
					continue
				}

				if !braceEscaped && code == CHAR_DOT {
					code = advance()
					if code == CHAR_DOT {
						isBrace = true
						token.IsBrace = true
						isGlob = true
						token.IsGlob = true
						finished = true

						if scanToEnd {
							continue
						}
						break
					}
				}

				if !braceEscaped && code == CHAR_COMMA {
					isBrace = true
					token.IsBrace = true
					isGlob = true
					token.IsGlob = true
					finished = true

					if scanToEnd {
						continue
					}
					break
				}

				if code == CHAR_RIGHT_CURLY_BRACE {
					braces--
					if braces == 0 {
						braceEscaped = false
						isBrace = true
						token.IsBrace = true
						finished = true
						break
					}
				}
			}

			if scanToEnd {
				continue
			}
			break
		}

		// Handle forward slashes (path separators)
		if code == CHAR_FORWARD_SLASH {
			slashes = append(slashes, index)
			tokens = append(tokens, token)
			token = ScanToken{Value: "", Depth: 0, IsGlob: false}

			if finished {
				continue
			}
			if prev == CHAR_DOT && index == (start+1) {
				start += 2
				continue
			}

			lastIndex = index + 1
			continue
		}

		// Handle extglob characters: +(, @(, *(, ?(, !(
		if !opts.Noext {
			isExtglobChar := code == CHAR_PLUS ||
				code == CHAR_AT ||
				code == CHAR_ASTERISK ||
				code == CHAR_QUESTION_MARK ||
				code == CHAR_EXCLAMATION_MARK

			if isExtglobChar && peek() == CHAR_LEFT_PARENTHESES {
				isGlob = true
				token.IsGlob = true
				isExtglob = true
				token.IsExtglob = true
				finished = true

				if code == CHAR_EXCLAMATION_MARK && index == start {
					negatedExtglob = true
				}

				if scanToEnd {
					for !eos() {
						code = advance()
						if code == 0 {
							break
						}

						if code == CHAR_BACKWARD_SLASH {
							backslashes = true
							token.Backslashes = true
							code = advance()
							continue
						}

						if code == CHAR_RIGHT_PARENTHESES {
							isGlob = true
							token.IsGlob = true
							finished = true
							break
						}
					}
					continue
				}
				break
			}
		}

		// Handle asterisks
		if code == CHAR_ASTERISK {
			if prev == CHAR_ASTERISK {
				isGlobstar = true
				token.IsGlobstar = true
			}
			isGlob = true
			token.IsGlob = true
			finished = true

			if scanToEnd {
				continue
			}
			break
		}

		// Handle question marks
		if code == CHAR_QUESTION_MARK {
			isGlob = true
			token.IsGlob = true
			finished = true

			if scanToEnd {
				continue
			}
			break
		}

		// Handle square brackets: [a-z]
		if code == CHAR_LEFT_SQUARE_BRACKET {
			for !eos() {
				next = advance()
				if next == 0 {
					break
				}

				if next == CHAR_BACKWARD_SLASH {
					backslashes = true
					token.Backslashes = true
					advance()
					continue
				}

				if next == CHAR_RIGHT_SQUARE_BRACKET {
					isBracket = true
					token.IsBracket = true
					isGlob = true
					token.IsGlob = true
					finished = true
					break
				}
			}

			if scanToEnd {
				continue
			}
			break
		}

		// Handle negation: leading !
		if !opts.Nonegate && code == CHAR_EXCLAMATION_MARK && index == start {
			negated = true
			token.Negated = true
			start++
			continue
		}

		// Handle parentheses as glob characters
		if !opts.Noparen && code == CHAR_LEFT_PARENTHESES {
			isGlob = true
			token.IsGlob = true

			if scanToEnd {
				for !eos() {
					code = advance()
					if code == 0 {
						break
					}

					if code == CHAR_LEFT_PARENTHESES {
						backslashes = true
						token.Backslashes = true
						code = advance()
						continue
					}

					if code == CHAR_RIGHT_PARENTHESES {
						finished = true
						break
					}
				}
				continue
			}
			break
		}

		if isGlob {
			finished = true

			if scanToEnd {
				continue
			}
			break
		}
	}

	// When noext is set, clear extglob and glob flags
	if opts.Noext {
		isExtglob = false
		isGlob = false
	}

	// Compute base path and glob pattern
	base := str
	prefix := ""
	glob := ""

	if start > 0 {
		prefix = str[:start]
		str = str[start:]
		lastIndex -= start
	}

	if isGlob && lastIndex > 0 {
		base = str[:lastIndex]
		glob = str[lastIndex:]
	} else if isGlob {
		base = ""
		glob = str
	} else {
		base = str
	}

	// Trim trailing path separator from base, unless it's the root or full path
	if base != "" && base != "/" && base != str {
		if isPathSeparator(base[len(base)-1]) {
			base = base[:len(base)-1]
		}
	}

	// Optionally strip backslashes from base and glob
	if opts.Unescape {
		if glob != "" {
			glob = RemoveBackslashes(glob)
		}
		if backslashes {
			base = RemoveBackslashes(base)
		}
	}

	state := ScanState{
		Prefix:         prefix,
		Input:          input,
		Start:          start,
		Base:           base,
		Glob:           glob,
		IsBrace:        isBrace,
		IsBracket:      isBracket,
		IsGlob:         isGlob,
		IsExtglob:      isExtglob,
		IsGlobstar:     isGlobstar,
		Negated:        negated,
		NegatedExtglob: negatedExtglob,
	}

	// Populate tokens if requested
	if opts.Tokens {
		state.MaxDepth = 0
		if !isPathSeparator(code) {
			tokens = append(tokens, token)
		}
		state.Tokens = tokens
	}

	// Populate parts and token values if requested
	if opts.Parts || opts.Tokens {
		var prevIndex int
		hasPrevIndex := false

		for idx := 0; idx < len(slashes); idx++ {
			n := start
			if hasPrevIndex {
				n = prevIndex + 1
			}
			i := slashes[idx]
			value := input[n:i]

			if opts.Tokens {
				if idx == 0 && start != 0 {
					tokens[idx].IsPrefix = true
					tokens[idx].Value = prefix
				} else {
					tokens[idx].Value = value
				}
				setDepth(&tokens[idx])
				state.MaxDepth += tokens[idx].Depth
			}

			if idx != 0 || value != "" {
				parts = append(parts, value)
			}

			prevIndex = i
			hasPrevIndex = true
		}

		if hasPrevIndex && prevIndex+1 < len(input) {
			value := input[prevIndex+1:]
			parts = append(parts, value)

			if opts.Tokens {
				lastToken := &tokens[len(tokens)-1]
				lastToken.Value = value
				setDepth(lastToken)
				state.MaxDepth += lastToken.Depth
			}
		}

		state.Slashes = slashes
		state.Parts = parts
	}

	return state
}
