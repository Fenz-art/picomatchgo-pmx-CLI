package main

import (
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
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

func TestPMXScan(t *testing.T) {
	cmd := exec.Command("go", "run", ".", "scan", "**/parser/*.go")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("pmx scan failed: %v\n%s", err, out)
	}

	content := string(out)
	for _, token := range []string{"PICOMATCH SCAN", "Segments", "globstar"} {
		if !strings.Contains(content, token) {
			t.Fatalf("expected scan output to contain %q, got:\n%s", token, content)
		}
	}
}

func TestPMXParse(t *testing.T) {
	cmd := exec.Command("go", "run", ".", "parse", "foo/{bar,baz}/@(a|b).go")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("pmx parse failed: %v\n%s", err, out)
	}

	content := string(out)
	for _, token := range []string{"PICOMATCH PARSE", "BRACE", "EXTGLOB"} {
		if !strings.Contains(content, token) {
			t.Fatalf("expected parse output to contain %q, got:\n%s", token, content)
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

func TestPMXHelpListsAllCommands(t *testing.T) {
	cmd := exec.Command("go", "run", ".", "help")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("pmx help failed: %v\n%s", err, out)
	}

	content := string(out)
	for _, token := range []string{"pmx match", "pmx scan", "pmx parse", "pmx explain", "pmx validate", "pmx compat", "pmx bench", "pmx fuzz", "pmx ci", "pmx doctor"} {
		if !strings.Contains(content, token) {
			t.Fatalf("expected help output to contain %q, got:\n%s", token, content)
		}
	}
}

func TestPMXDoctor(t *testing.T) {
	cmd := exec.Command("go", "run", ".", "doctor")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("pmx doctor failed: %v\n%s", err, out)
	}

	content := string(out)
	for _, token := range []string{"PMX DOCTOR", "Project", "Diagnostics", "Summary"} {
		if !strings.Contains(content, token) {
			t.Fatalf("expected doctor output to contain %q, got:\n%s", token, content)
		}
	}
}

func TestPMXDoctorConfig(t *testing.T) {
	cmd := exec.Command("go", "run", ".", "doctor", "config")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("pmx doctor config failed: %v\n%s", err, out)
	}

	content := string(out)
	for _, token := range []string{"Configuration", "package.json", "go.mod", "Cargo.toml"} {
		if !strings.Contains(content, token) {
			t.Fatalf("expected config output to contain %q, got:\n%s", token, content)
		}
	}
}

func TestPMXDoctorJSON(t *testing.T) {
	cmd := exec.Command("go", "run", ".", "doctor", "--json")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("pmx doctor --json failed: %v\n%s", err, out)
	}

	content := string(out)
	for _, token := range []string{"\"version\"", "\"project\"", "\"diagnostics\""} {
		if !strings.Contains(content, token) {
			t.Fatalf("expected JSON output to contain %q, got:\n%s", token, content)
		}
	}
}

func TestPMXDoctorCI(t *testing.T) {
	cmd := exec.Command("go", "run", ".", "doctor", "--ci")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("pmx doctor --ci failed: %v\n%s", err, out)
	}

	content := string(out)
	for _, token := range []string{"PMX Doctor CI", "Result:"} {
		if !strings.Contains(content, token) {
			t.Fatalf("expected CI output to contain %q, got:\n%s", token, content)
		}
	}
}

func TestPMXAgentInspectJSON(t *testing.T) {
	cmd := exec.Command("go", "run", ".", "agent", "inspect", "--json")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("pmx agent inspect --json failed: %v\n%s", err, out)
	}

	var report map[string]interface{}
	if err := json.Unmarshal(out, &report); err != nil {
		t.Fatalf("invalid JSON output from pmx agent inspect --json: %v\n%s", err, out)
	}

	for _, key := range []string{"version", "project", "diagnostics"} {
		if _, ok := report[key]; !ok {
			t.Fatalf("missing %q in agent inspect JSON: %s", key, out)
		}
	}
	if got := report["version"]; got != "1" {
		t.Fatalf("version = %v; want 1", got)
	}
	if _, ok := report["project"].(map[string]interface{}); !ok {
		t.Fatalf("project should be an object: %T", report["project"])
	}
	if _, ok := report["diagnostics"].([]interface{}); !ok {
		t.Fatalf("diagnostics should be an array: %T", report["diagnostics"])
	}
}

func TestPMXAgentCheckJSON(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "go", "run", ".", "agent", "check", "--json")
	out, err := cmd.CombinedOutput()
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			t.Fatalf("pmx agent check --json exceeded its 45s test budget; the agent gate must not recursively run the full CI suite\n%s", out)
		}
		t.Fatalf("pmx agent check --json failed: %v\n%s", err, out)
	}

	content := string(out)
	for _, token := range []string{"\"version\"", "\"result\"", "\"diagnostics\"", "\"checks\"", "\"next_actions\""} {
		if !strings.Contains(content, token) {
			t.Fatalf("expected agent check JSON to contain %q, got:\n%s", token, content)
		}
	}
}

func TestPMXAgentCheckStrictContract(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "go", "run", ".", "agent", "check", "--json")
	out, err := cmd.CombinedOutput()
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			t.Fatalf("pmx agent check --json exceeded its 45s test budget; the agent gate must not recursively run the full CI suite\n%s", out)
		}
		t.Fatalf("pmx agent check --json failed: %v\n%s", err, out)
	}

	var report map[string]interface{}
	if err := json.Unmarshal(out, &report); err != nil {
		t.Fatalf("invalid JSON output from pmx agent check --json: %v\n%s", err, out)
	}

	wantKeys := []string{"version", "result", "diagnostics", "checks", "next_actions"}
	for _, key := range wantKeys {
		if _, ok := report[key]; !ok {
			t.Fatalf("missing %q in agent check JSON: %s", key, out)
		}
	}
	for _, key := range []string{"status", "project", "validation", "phase"} {
		if _, ok := report[key]; ok {
			t.Fatalf("legacy field %q should not be present in strict ADLC contract: %s", key, out)
		}
	}

	if got := report["version"]; got != "1" {
		t.Fatalf("version = %v; want 1", got)
	}
	if got := report["result"]; got != "pass" && got != "warn" && got != "fail" {
		t.Fatalf("result = %v; want pass|warn|fail", got)
	}
	if _, ok := report["checks"].([]interface{}); !ok {
		t.Fatalf("checks should be an array: %T", report["checks"])
	}
	if _, ok := report["next_actions"].([]interface{}); !ok {
		t.Fatalf("next_actions should be an array: %T", report["next_actions"])
	}
	checks := report["checks"].([]interface{})
	wantChecks := map[string]bool{"doctor": false, "validate": false, "compat": false, "regression": false}
	for _, item := range checks {
		check, ok := item.(map[string]interface{})
		if !ok {
			t.Fatalf("check should be an object: %T", item)
		}
		name, _ := check["name"].(string)
		if name == "ci" {
			t.Fatalf("agent check must not invoke pmx ci; it recursively re-enters the agent test path: %s", out)
		}
		if _, known := wantChecks[name]; known {
			wantChecks[name] = true
		}
	}
	for name, seen := range wantChecks {
		if !seen {
			t.Fatalf("missing executed agent check %q: %s", name, out)
		}
	}
}

func TestPMXDoctorFixturePath(t *testing.T) {
	cmd := exec.Command("go", "run", "./cmd/pmx", "doctor", "--json", "fixtures/js-ts-fail")
	cmd.Dir = filepath.Join("..", "..")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("pmx doctor --json fixture path failed: %v\n%s", err, out)
	}

	var report map[string]interface{}
	if err := json.Unmarshal(out, &report); err != nil {
		t.Fatalf("invalid JSON output from fixture path: %v\n%s", err, out)
	}
	project := report["project"].(map[string]interface{})
	if got := project["ecosystem"]; got != "javascript/typescript" {
		t.Fatalf("ecosystem = %v; want javascript/typescript", got)
	}
	if got := project["package_manager"]; got != "pnpm" {
		t.Fatalf("package_manager = %v; want pnpm", got)
	}
	if got := report["summary"].(map[string]interface{})["warn"]; got != float64(3) {
		t.Fatalf("warn summary = %v; want 3", got)
	}
	if got := report["summary"].(map[string]interface{})["fail"]; got != float64(0) {
		t.Fatalf("fail summary = %v; want 0", got)
	}
}

func TestPMXDoctorDetectsJSProject(t *testing.T) {
	tempDir := t.TempDir()
	if err := os.WriteFile(tempDir+"/package.json", []byte(`{"name":"demo","version":"1.0.0","scripts":{"build":"next build","lint":"eslint ."}}`), 0o600); err != nil {
		t.Fatalf("write package.json: %v", err)
	}
	if err := os.WriteFile(tempDir+"/pnpm-lock.yaml", []byte("lockfileVersion: '9.0'\n"), 0o600); err != nil {
		t.Fatalf("write pnpm-lock.yaml: %v", err)
	}
	if err := os.WriteFile(tempDir+"/tsconfig.json", []byte(`{"compilerOptions":{"strict":false}}`), 0o600); err != nil {
		t.Fatalf("write tsconfig.json: %v", err)
	}
	if err := os.WriteFile(tempDir+"/.eslintrc.json", []byte(`{"root":true}`), 0o600); err != nil {
		t.Fatalf("write .eslintrc.json: %v", err)
	}

	binPath := testBinaryPath(tempDir, "pmx-test")
	buildCmd := exec.Command("go", "build", "-o", binPath, ".")
	buildCmd.Dir = "."
	if out, err := buildCmd.CombinedOutput(); err != nil {
		t.Fatalf("build pmx binary failed: %v\n%s", err, out)
	}

	cmd := exec.Command(binPath, "doctor", "--json")
	cmd.Dir = tempDir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("pmx doctor --json in JS project failed: %v\n%s", err, out)
	}

	content := string(out)
	for _, token := range []string{"\"ecosystem\": \"javascript/typescript\"", "\"package_manager\": \"pnpm\"", "\"severity\": \"warn\"", "\"id\": \"PMX-TS-001\""} {
		if !strings.Contains(content, token) {
			t.Fatalf("expected JSON report to contain %q, got:\n%s", token, content)
		}
	}
}

func buildDoctorBinary(t *testing.T) string {
	t.Helper()
	binPath := testBinaryPath(t.TempDir(), "pmx-doctor-test")
	cmd := exec.Command("go", "build", "-o", binPath, ".")
	cmd.Dir = filepath.Join("..", "..", "cmd", "pmx")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("build pmx binary failed: %v\n%s", err, out)
	}
	return binPath
}

func testBinaryPath(dir, name string) string {
	if runtime.GOOS == "windows" {
		name += ".exe"
	}
	return filepath.Join(dir, name)
}

func TestPMXDoctorBrokenFixture(t *testing.T) {
	binPath := buildDoctorBinary(t)
	fixtureDir := filepath.Join("..", "..", "fixtures", "js-ts-broken")
	cmd := exec.Command(binPath, "doctor")
	cmd.Dir = fixtureDir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("pmx doctor in broken fixture failed: %v\n%s", err, out)
	}

	content := string(out)
	for _, token := range []string{"PMX DOCTOR", "PMX-PKG-001", "PMX-TS-001", "PMX-ESLINT-001", "Result: WARNING"} {
		if !strings.Contains(content, token) {
			t.Fatalf("expected doctor output to contain %q, got:\n%s", token, content)
		}
	}
}

func TestPMXDoctorBrokenFixtureJSONContract(t *testing.T) {
	binPath := buildDoctorBinary(t)
	fixtureDir := filepath.Join("..", "..", "fixtures", "js-ts-broken")
	cmd := exec.Command(binPath, "doctor", "--json")
	cmd.Dir = fixtureDir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("pmx doctor --json in broken fixture failed: %v\n%s", err, out)
	}

	var report map[string]interface{}
	if err := json.Unmarshal(out, &report); err != nil {
		t.Fatalf("invalid JSON output: %v\n%s", err, out)
	}
	project := report["project"].(map[string]interface{})
	if got := project["ecosystem"]; got != "javascript/typescript" {
		t.Fatalf("ecosystem = %v; want javascript/typescript", got)
	}
	if got := project["package_manager"]; got != "pnpm" {
		t.Fatalf("package_manager = %v; want pnpm", got)
	}
	diagnostics := report["diagnostics"].([]interface{})
	if len(diagnostics) != 3 {
		t.Fatalf("expected 3 diagnostics, got %d: %s", len(diagnostics), out)
	}
	for _, item := range diagnostics {
		d := item.(map[string]interface{})
		if got := d["severity"]; got != "warn" {
			t.Fatalf("diagnostic severity = %v; want warn: %v", got, d)
		}
	}
	if got := report["summary"].(map[string]interface{})["warn"]; got != float64(3) {
		t.Fatalf("warn summary = %v; want 3", got)
	}
}

func TestPMXDoctorCIExitCodes(t *testing.T) {
	binPath := buildDoctorBinary(t)
	warnDir := filepath.Join("..", "..", "fixtures", "js-ts-fail")
	warnCmd := exec.Command(binPath, "doctor", "--ci")
	warnCmd.Dir = warnDir
	warnOut, warnErr := warnCmd.CombinedOutput()
	if warnErr != nil {
		if exitErr, ok := warnErr.(*exec.ExitError); !ok || exitErr.ExitCode() != 0 {
			t.Fatalf("pmx doctor --ci in warning fixture returned unexpected error: %v\n%s", warnErr, warnOut)
		}
	}
	if !strings.Contains(string(warnOut), "Result: WARNING") {
		t.Fatalf("expected warning CI result, got:\n%s", warnOut)
	}

	failDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(failDir, "tsconfig.json"), []byte(`{"compilerOptions":{"strict":false}}`), 0o600); err != nil {
		t.Fatalf("write fail fixture tsconfig: %v", err)
	}
	failCmd := exec.Command(binPath, "doctor", "--ci")
	failCmd.Dir = failDir
	failOut, failErr := failCmd.CombinedOutput()
	if failErr == nil {
		t.Fatalf("expected failed CI fixture to exit non-zero, got:\n%s", failOut)
	}
	if !strings.Contains(string(failOut), "Result: FAILURE") {
		t.Fatalf("expected failure CI result, got:\n%s", failOut)
	}
}

func TestPMXRegression(t *testing.T) {
	cmd := exec.Command("go", "run", ".", "regression")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("pmx regression failed: %v\n%s", err, out)
	}

	content := string(out)
	for _, token := range []string{"PMX REGRESSION SUITE", "RESULT: PASS"} {
		if !strings.Contains(content, token) {
			t.Fatalf("expected regression output to contain %q, got:\n%s", token, content)
		}
	}
}

func TestPMXCIJSON(t *testing.T) {
	cmd := exec.Command("go", "run", ".", "ci", "--json")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("pmx ci --json failed: %v\n%s", err, out)
	}

	var report map[string]interface{}
	if err := json.Unmarshal(out, &report); err != nil {
		t.Fatalf("invalid JSON output from pmx ci --json: %v\n%s", err, out)
	}
	if got := report["result"]; got != "pass" && got != "warning" {
		t.Fatalf("unexpected ci result: %v; output:\n%s", got, out)
	}
	if _, ok := report["checks"]; !ok {
		t.Fatalf("expected checks in JSON report, got:\n%s", out)
	}
}

func TestPMXCIUsage(t *testing.T) {
	cmd := exec.Command("go", "run", ".", "ci")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("pmx ci failed: %v\n%s", err, out)
	}

	if !strings.Contains(string(out), "Usage:") || !strings.Contains(string(out), "pmx ci run") {
		t.Fatalf("expected ci usage output, got:\n%s", out)
	}
}
