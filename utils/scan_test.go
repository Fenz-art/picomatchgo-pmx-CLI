package utils

import "testing"

func TestScanSimple(t *testing.T) {
	state := Scan("foo/bar/*.js", nil)
	if !state.IsGlob {
		t.Fatal("expected IsGlob true")
	}
	if state.Base != "foo/bar" {
		t.Fatalf("expected base foo/bar, got %q", state.Base)
	}
	if state.Glob != "*.js" {
		t.Fatalf("expected glob *.js, got %q", state.Glob)
	}
}

func TestScanNegated(t *testing.T) {
	state := Scan("!./foo/*.js", nil)
	if !state.IsGlob {
		t.Fatal("expected IsGlob true")
	}
	if !state.Negated {
		t.Fatal("expected Negated true")
	}
	if state.Prefix != "!./" {
		t.Fatalf("expected prefix !./, got %q", state.Prefix)
	}
}
