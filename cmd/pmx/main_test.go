package main

import (
	"os/exec"
	"strings"
	"testing"
)

func TestPMXMatch(t *testing.T) {
	cmd := exec.Command("go", "run", ".", "match", "*.js", "foo.js")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("pmx match failed: %v\n%s", err, out)
	}

	if !strings.Contains(string(out), "MATCH") {
		t.Fatalf("expected MATCH output, got:\n%s", out)
	}
}

func TestPMXMatchNoMatch(t *testing.T) {
	cmd := exec.Command("go", "run", ".", "match", "*.js", "foo.txt")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("pmx match failed: %v\n%s", err, out)
	}

	if !strings.Contains(string(out), "NO MATCH") {
		t.Fatalf("expected NO MATCH output, got:\n%s", out)
	}
}

func TestPMXExplain(t *testing.T) {
	cmd := exec.Command("go", "run", ".", "explain", "**/*.go", "--input", "src/parser/scan.go")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("pmx explain failed: %v\n%s", err, out)
	}

	content := string(out)
	for _, token := range []string{"Scanner", "Parser", "Compiler", "Matcher", "globstar", "segments", "regex"} {
		if !strings.Contains(content, token) {
			t.Fatalf("expected explain output to contain %q, got:\n%s", token, content)
		}
	}
}

func TestPMXValidateValidPattern(t *testing.T) {
	cmd := exec.Command("go", "run", ".", "validate", "*.go", "--input", "foo.go")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("pmx validate failed: %v\n%s", err, out)
	}

	if !strings.Contains(string(out), "Result           VALID") {
		t.Fatalf("expected valid result, got:\n%s", out)
	}
}

func TestPMXValidateInvalidPattern(t *testing.T) {
	cmd := exec.Command("go", "run", "./cmd/pmx", "validate", "[", "--input", "foo.go")
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatalf("expected failure for invalid pattern, got success\n%s", out)
	}
}

func TestPMXCompatSuite(t *testing.T) {
	cmd := exec.Command("go", "run", ".", "compat", "--suite", "basic")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("pmx compat failed: %v\n%s", err, out)
	}

	if !strings.Contains(string(out), "Behavior: EQUIVALENT") {
		t.Fatalf("expected compatibility equivalent, got:\n%s", out)
	}
}
