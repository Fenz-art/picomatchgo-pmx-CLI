package utils

import "strings"

type ScanToken struct {
	Value       string
	Depth       int
	IsGlob      bool
	IsPrefix    bool
	Backslashes bool
	IsBrace     bool
	IsBracket   bool
	IsExtglob   bool
	IsGlobstar  bool
	Negated     bool
}

type ScanState struct {
	Prefix        string
	Input         string
	Start         int
	Base          string
	Glob          string
	IsBrace       bool
	IsBracket     bool
	IsGlob        bool
	IsExtglob     bool
	IsGlobstar    bool
	Negated       bool
	NegatedExtglob bool
	Slashes       []int
	Parts         []string
	Tokens        []ScanToken
	MaxDepth      int
}

func isPathSeparator(code int) bool {
	return code == CHAR_FORWARD_SLASH || code == CHAR_BACKWARD_SLASH
}

func depth(token *ScanToken) {
	if token.IsPrefix != true {
		if token.IsGlobstar {
			token.Depth = 0
		} else {
			token.Depth = 1
		}
	}
}

func Scan(input string, options *Options) ScanState {
	opts := options
	if opts == nil {
		opts = &Options{}
	}

	length := len(input) - 1
	scanToEnd := opts.Parts || opts.ScanToEnd
	slashes := make([]int, 0, 4)
	tokens := make([]ScanToken, 0, 8)
	parts := make([]string, 0, 8)

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
	var prev int
	var code int
	token := ScanToken{Value: "", Depth: 0, IsGlob: false}

	eos := func() bool { return index >= length }
	peek := func() int {
		if index+1 < len(str) {
			return int(str[index+1])
		}
		return -1
	}
	advance := func() int {
		prev = code
		index++
		if index >= len(str) {
			return -1
		}
		return int(str[index])
	}

	for index < length {
		code = advance()
		if code == -1 {
			break
		}

		if code == CHAR_BACKWARD_SLASH {
			backslashes = true
			token.Backslashes = true
			code = advance()

			if code == CHAR_LEFT_CURLY_BRACE {
				braceEscaped = true
			}
			continue
		}

		if braceEscaped || code == CHAR_LEFT_CURLY_BRACE {
			braces++

			for !eos() {
				code = advance()
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

				if !braceEscaped && code == CHAR_DOT && peek() == CHAR_DOT {
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

		if !opts.Noext {
			isExtglobChar := code == CHAR_PLUS || code == CHAR_AT || code == CHAR_ASTERISK || code == CHAR_QUESTION_MARK || code == CHAR_EXCLAMATION_MARK

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

		if code == CHAR_QUESTION_MARK {
			isGlob = true
			token.IsGlob = true
			finished = true

			if scanToEnd {
				continue
			}
			break
		}

		if code == CHAR_LEFT_SQUARE_BRACKET {
			for !eos() {
				next := advance()
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

		if !opts.Nonegate && code == CHAR_EXCLAMATION_MARK && index == start {
			negated = true
			token.Negated = true
			start++
			continue
		}

		if !opts.Noparen && code == CHAR_LEFT_PARENTHESES {
			isGlob = true
			token.IsGlob = true

			if scanToEnd {
				for !eos() {
					code = advance()
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

	if opts.Noext {
		isExtglob = false
		isGlob = false
	}

	base := str
	prefix := ""
	glob := ""

	if start > 0 {
		prefix = str[:start]
		str = str[start:]
		lastIndex -= start
	}

	if base != "" && isGlob && lastIndex > 0 {
		base = str[:lastIndex]
		glob = str[lastIndex:]
	} else if isGlob {
		base = ""
		glob = str
	} else {
		base = str
	}

	if base != "" && base != "/" && base != str {
		if isPathSeparator(int(base[len(base)-1])) {
			base = base[:len(base)-1]
		}
	}

	if opts.Unescape {
		if glob != "" {
			glob = RemoveBackslashes(glob)
		}
		if base != "" && backslashes {
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

	if opts.Tokens {
		state.MaxDepth = 0
		if !isPathSeparator(code) {
			tokens = append(tokens, token)
		}
		state.Tokens = tokens
	}

	if opts.Parts || opts.Tokens {
		var prevIndex int
		for idx, slashPos := range slashes {
			n := 0
			if idx > 0 {
				n = prevIndex + 1
			} else {
				n = start
			}
			value := input[n:slashPos]
			if opts.Tokens {
				if idx == 0 && start != 0 {
					tokens[idx].IsPrefix = true
					tokens[idx].Value = prefix
				} else {
					tokens[idx].Value = value
				}
				depth(&tokens[idx])
				state.MaxDepth += tokens[idx].Depth
			}
			if idx != 0 || value != "" {
				parts = append(parts, value)
			}
			prevIndex = slashPos
		}

		if prevIndex+1 < len(input) {
			value := input[prevIndex+1:]
			parts = append(parts, value)

			if opts.Tokens {
				tokenIndex := len(tokens) - 1
				tokens[tokenIndex].Value = value
				depth(&tokens[tokenIndex])
				state.MaxDepth += tokens[tokenIndex].Depth
			}
		}

		state.Slashes = slashes
		state.Parts = parts
	}

	return state
}
