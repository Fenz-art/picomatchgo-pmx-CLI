package picomatch

import (
	"testing"
)

func TestRootMakeReAndIsMatch(t *testing.T) {
	regex, err := MakeRe("*.js", nil)
	if err != nil {
		t.Fatalf("MakeRe failed: %v", err)
	}
	if !regex.MatchString("foo.js") {
		t.Fatal("expected foo.js to match")
	}
	if regex.MatchString(".foo.js") {
		t.Fatal("expected .foo.js not to match")
	}
}

func TestRootIsMatchArray(t *testing.T) {
	ok, err := IsMatch("a.a", []string{"b.*", "*.a"}, nil)
	if err != nil {
		t.Fatalf("IsMatch failed: %v", err)
	}
	if !ok {
		t.Fatal("expected a.a to match patterns")
	}
}

func TestRootDotfileHandling(t *testing.T) {
	ok, err := IsMatch("foo/.bar", "foo/*", nil)
	if err != nil {
		t.Fatalf("IsMatch failed: %v", err)
	}
	if ok {
		t.Fatal("expected foo/.bar to not match foo/*")
	}
}
