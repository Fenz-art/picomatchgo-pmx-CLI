package picomatch

// DefaultOptions returns a new Options struct initialized with safe defaults for
// glob scanning and matching.
func DefaultOptions() Options {
	return Options{
		MaxLength:       65536,
		MaxExtglobDepth: 0,
		Fastpaths:       true,
	}
}
