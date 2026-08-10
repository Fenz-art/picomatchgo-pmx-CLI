package picomatch

import "testing"

func TestDefaultOptions(t *testing.T) {
	opts := DefaultOptions()

	if opts.MaxLength != 65536 {
		t.Errorf("Expected MaxLength 65536, got %d", opts.MaxLength)
	}

	if opts.MaxExtglobDepth != 0 {
		t.Errorf("Expected MaxExtglobDepth 0, got %d", opts.MaxExtglobDepth)
	}

	if !opts.Fastpaths {
		t.Errorf("Expected Fastpaths true, got %v", opts.Fastpaths)
	}

	if opts.Dot {
		t.Errorf("Expected Dot false by default, got %v", opts.Dot)
	}

	if opts.Nocase {
		t.Errorf("Expected Nocase false by default, got %v", opts.Nocase)
	}
}

func TestOptionsCustomConfiguration(t *testing.T) {
	opts := Options{
		Dot:       true,
		Nocase:    true,
		Windows:   true,
		Ignore:    []string{"*.tmp", "*.bak"},
		MaxLength: 1024,
	}

	if !opts.Dot {
		t.Errorf("Expected Dot to be true")
	}

	if !opts.Nocase {
		t.Errorf("Expected Nocase to be true")
	}

	if !opts.Windows {
		t.Errorf("Expected Windows to be true")
	}

	if len(opts.Ignore) != 2 {
		t.Errorf("Expected 2 ignore patterns, got %d", len(opts.Ignore))
	}
}
