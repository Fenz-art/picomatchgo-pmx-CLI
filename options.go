package picomatch

// DefaultOptions returns a new Options struct initialized with default settings.
func DefaultOptions() Options {
	return Options{
		MaxLength:       65536,
		MaxExtglobDepth: 0,
		Fastpaths:       true,
	}
}
