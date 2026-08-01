package options

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
}
