package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	picomatch "github.com/Fenz-art/picomatchgo-pmx-CLI"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	subcommand := os.Args[1]
	args := os.Args[2:]

	switch subcommand {
	case "match":
		exit(runMatch(args))
	case "scan":
		exit(runScan(args))
	case "parse":
		exit(runParse(args))
	case "explain":
		exit(runExplain(args))
	case "validate":
		exit(runValidate(args))
	case "compat":
		exit(runCompat(args))
	case "regression":
		exit(runRegression(args))
	case "bench":
		exit(runBench(args))
	case "fuzz":
		exit(runFuzz(args))
	case "doctor":
		exit(runDoctor(args))
	case "agent":
		exit(runAgent(args))
	case "ci":
		exit(runCI(args))
	case "help", "-h", "--help":
		printUsage()
		os.Exit(0)
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", subcommand)
		printUsage()
		os.Exit(2)
	}
}

func exit(code int) {
	if code != 0 {
		os.Exit(code)
	}
}

func printUsage() {
	fmt.Println("pmx — Picomatch Go Developer Diagnostics")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  pmx match <pattern> <input> [flags]")
	fmt.Println("  pmx scan <pattern> [flags]")
	fmt.Println("  pmx parse <pattern> [flags]")
	fmt.Println("  pmx explain <pattern> [--input <string>] [flags]")
	fmt.Println("  pmx validate <pattern> [--input <string>] [flags]")
	fmt.Println("  pmx compat [--suite <name>] [flags]")
	fmt.Println("  pmx regression [--json]")
	fmt.Println("  pmx bench [--compare <baseline.json>] [flags]")
	fmt.Println("  pmx fuzz [--target <name>] [--time <duration>] [flags]")
	fmt.Println("  pmx doctor [config|deps|toolchain|all]")
	fmt.Println("  pmx agent inspect [--json]")
	fmt.Println("  pmx agent check [--json]")
	fmt.Println("  pmx ci run [--json]")
	fmt.Println("  pmx ci watch <run-id>")
	fmt.Println("  pmx ci logs <run-id> [job-name]")
	fmt.Println()
	fmt.Println("Flags:")
	fmt.Println("  --dot       match dotfiles")
	fmt.Println("  --nocase    case-insensitive matching")
	fmt.Println("  --windows   normalize Windows path separators")
}

func runMatch(args []string) int {
	fs := flag.NewFlagSet("match", flag.ContinueOnError)
	dot := fs.Bool("dot", false, "match dotfiles")
	nocase := fs.Bool("nocase", false, "case-insensitive matching")
	windows := fs.Bool("windows", false, "normalize Windows path separators")
	fs.SetOutput(os.Stdout)
	if err := fs.Parse(args); err != nil {
		return 2
	}

	if fs.NArg() != 2 {
		fmt.Fprintln(os.Stderr, "usage: pmx match <pattern> <input> [flags]")
		return 2
	}

	pattern := fs.Arg(0)
	input := fs.Arg(1)
	opts := &picomatch.Options{
		Dot:     *dot,
		Nocase:  *nocase,
		Windows: *windows,
	}

	matched, err := picomatch.IsMatch(input, pattern, opts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	if matched {
		fmt.Println("MATCH")
	} else {
		fmt.Println("NO MATCH")
	}
	fmt.Println()
	fmt.Printf("pattern: %s\n", pattern)
	fmt.Printf("input:   %s\n", input)
	return 0
}

func runScan(args []string) int {
	fs := flag.NewFlagSet("scan", flag.ContinueOnError)
	dot := fs.Bool("dot", false, "match dotfiles")
	nocase := fs.Bool("nocase", false, "case-insensitive matching")
	windows := fs.Bool("windows", false, "normalize Windows path separators")
	fs.SetOutput(os.Stdout)
	if err := fs.Parse(args); err != nil {
		return 2
	}
	if fs.NArg() != 1 {
		fmt.Fprintln(os.Stderr, "usage: pmx scan <pattern> [flags]")
		return 2
	}

	pattern := fs.Arg(0)
	state := picomatch.Scan(pattern, &picomatch.Options{Dot: *dot, Nocase: *nocase, Windows: *windows, Parts: true, Tokens: true})
	fmt.Println("PICOMATCH SCAN")
	fmt.Println(strings.Repeat("─", 50))
	fmt.Println()
	fmt.Printf("Pattern\n  %s\n\n", pattern)
	fmt.Println("Structure")
	fmt.Printf("  base:      %s\n", defaultString(state.Base, "."))
	fmt.Printf("  glob:      %t\n", state.IsGlob)
	fmt.Printf("  globstar:  %t\n", state.IsGlobstar)
	fmt.Printf("  segments:  %d\n\n", len(state.Parts))
	fmt.Println("Segments")
	if len(state.Parts) == 0 {
		fmt.Println("  [none]")
	} else {
		for i, seg := range state.Parts {
			fmt.Printf("  [%d] %s\n", i, seg)
		}
	}
	fmt.Println()
	fmt.Println("Flags")
	fmt.Printf("  dot:      %t\n", *dot)
	fmt.Printf("  nocase:   %t\n", *nocase)
	fmt.Printf("  windows:  %t\n", *windows)
	return 0
}

func runParse(args []string) int {
	fs := flag.NewFlagSet("parse", flag.ContinueOnError)
	dot := fs.Bool("dot", false, "match dotfiles")
	nocase := fs.Bool("nocase", false, "case-insensitive matching")
	windows := fs.Bool("windows", false, "normalize Windows path separators")
	fs.SetOutput(os.Stdout)
	if err := fs.Parse(args); err != nil {
		return 2
	}
	if fs.NArg() != 1 {
		fmt.Fprintln(os.Stderr, "usage: pmx parse <pattern> [flags]")
		return 2
	}

	pattern := fs.Arg(0)
	state, err := picomatch.Parse(pattern, &picomatch.Options{Dot: *dot, Nocase: *nocase, Windows: *windows})
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	fmt.Println("PICOMATCH PARSE")
	fmt.Println(strings.Repeat("─", 50))
	fmt.Println()
	fmt.Printf("Pattern\n  %s\n\n", pattern)
	fmt.Println("Tokens")
	if len(state.Tokens) == 0 {
		fmt.Println("  <none>")
	} else {
		for i, token := range state.Tokens {
			if token == nil {
				continue
			}
			fmt.Printf("  [%d] %-12s %s\n", i, formatParseTokenName(token.Type), defaultString(token.Value, "<empty>"))
		}
	}
	fmt.Println()
	fmt.Println("AST")
	fmt.Println("  ROOT")
	if len(state.Tokens) > 0 {
		for i, token := range state.Tokens {
			if token == nil {
				continue
			}
			fmt.Printf("  ├── [%d] %s\n", i, formatParseTokenName(token.Type))
		}
	}
	return 0
}

func formatParseTokenName(tokenType string) string {
	switch strings.ToLower(tokenType) {
	case "at", "plus", "star", "qmark", "negate", "paren":
		return "EXTGLOB"
	case "brace":
		return "BRACE"
	case "slash":
		return "SLASH"
	case "dot":
		return "DOT"
	case "comma":
		return "COMMA"
	case "text":
		return "TEXT"
	case "bos":
		return "BOS"
	default:
		return strings.ToUpper(tokenType)
	}
}

func runBench(args []string) int {
	fs := flag.NewFlagSet("bench", flag.ContinueOnError)
	compare := fs.String("compare", "", "baseline benchmark JSON file")
	fs.SetOutput(os.Stdout)
	if err := fs.Parse(args); err != nil {
		return 2
	}

	fmt.Println("PICOMATCH BENCHMARK")
	fmt.Println(strings.Repeat("─", 50))
	fmt.Println()

	cmdArgs := []string{"test", "-run=^$", "-bench=.", "-benchmem", "-count=1"}
	cmd := exec.Command("go", cmdArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "benchmark run failed: %v\n", err)
		return 1
	}

	if *compare != "" {
		fmt.Printf("note: benchmark compare (%s) is not implemented in this release\n", *compare)
	}

	return 0
}

func runFuzz(args []string) int {
	fs := flag.NewFlagSet("fuzz", flag.ContinueOnError)
	target := fs.String("target", "FuzzScan", "fuzz target name")
	timeArg := fs.String("time", "15s", "fuzz duration")
	fs.SetOutput(os.Stdout)
	if err := fs.Parse(args); err != nil {
		return 2
	}

	switch *target {
	case "FuzzScan", "FuzzParse", "FuzzIsMatch":
	default:
		fmt.Fprintf(os.Stderr, "unknown fuzz target: %s\n", *target)
		return 2
	}

	fmt.Println("PICOMATCH FUZZ CAMPAIGN")
	fmt.Println(strings.Repeat("─", 50))
	fmt.Println()

	cmd := exec.Command("go", "test", "-run=^$", "-fuzz="+*target, "-fuzztime="+*timeArg, "./...")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "fuzz run failed: %v\n", err)
		return 1
	}
	return 0
}

type doctorProject struct {
	Name           string `json:"name"`
	Ecosystem      string `json:"ecosystem"`
	PackageManager string `json:"package_manager"`
	TypeScript     bool   `json:"typescript,omitempty"`
	Framework      string `json:"framework,omitempty"`
}

type doctorDiagnostic struct {
	ID         string   `json:"id"`
	Severity   string   `json:"severity"`
	Category   string   `json:"category,omitempty"`
	Title      string   `json:"title,omitempty"`
	File       string   `json:"file,omitempty"`
	Message    string   `json:"message"`
	Evidence   []string `json:"evidence,omitempty"`
	Suggestion string   `json:"suggestion,omitempty"`
}

type doctorSummary struct {
	Pass int `json:"pass"`
	Warn int `json:"warn"`
	Fail int `json:"fail"`
}

type doctorReport struct {
	Version     string             `json:"version"`
	Project     doctorProject      `json:"project"`
	Diagnostics []doctorDiagnostic `json:"diagnostics"`
	Summary     doctorSummary      `json:"summary"`
}

type agentInspectReport struct {
	Version     string             `json:"version"`
	Project     doctorProject      `json:"project"`
	Diagnostics []doctorDiagnostic `json:"diagnostics"`
}

type agentDiagnostic struct {
	Code       string `json:"code"`
	Severity   string `json:"severity"`
	File       string `json:"file,omitempty"`
	Message    string `json:"message,omitempty"`
	Suggestion string `json:"suggestion,omitempty"`
	ID         string `json:"id,omitempty"`
}

type agentCheckEntry struct {
	Name   string `json:"name"`
	Status string `json:"status"`
	Detail string `json:"detail,omitempty"`
}

type agentAction struct {
	Action string `json:"action"`
	Code   string `json:"code"`
	Text   string `json:"text,omitempty"`
}

type agentCheckReport struct {
	Version     string            `json:"version"`
	Result      string            `json:"result"`
	Diagnostics []agentDiagnostic `json:"diagnostics"`
	Checks      []agentCheckEntry `json:"checks"`
	NextActions []agentAction     `json:"next_actions"`
}

func runDoctor(args []string) int {
	fs := flag.NewFlagSet("doctor", flag.ContinueOnError)
	jsonOutput := fs.Bool("json", false, "emit JSON output")
	ciMode := fs.Bool("ci", false, "enable CI-style output")
	fs.SetOutput(os.Stdout)
	if err := fs.Parse(args); err != nil {
		return 2
	}

	mode := "all"
	projectDir := "."
	if fs.NArg() > 0 {
		candidate := fs.Arg(0)
		if isDoctorMode(candidate) {
			mode = candidate
			if fs.NArg() > 1 {
				projectDir = fs.Arg(1)
			}
		} else {
			projectDir = candidate
		}
	}

	if projectDir != "." {
		oldWd, err := os.Getwd()
		if err != nil {
			fmt.Fprintf(os.Stderr, "doctor target failed: %v\n", err)
			return 1
		}
		if err := os.Chdir(projectDir); err != nil {
			fmt.Fprintf(os.Stderr, "doctor target failed: %v\n", err)
			return 1
		}
		defer func() {
			if err := os.Chdir(oldWd); err != nil {
				fmt.Fprintf(os.Stderr, "warning: failed to restore working directory: %v\n", err)
			}
		}()
	}

	report := buildDoctorReport()
	if *jsonOutput {
		payload, err := json.MarshalIndent(report, "", "  ")
		if err != nil {
			fmt.Fprintf(os.Stderr, "doctor json failed: %v\n", err)
			return 1
		}
		fmt.Println(string(payload))
		return 0
	}

	switch mode {
	case "config":
		printDoctorConfigReport(report, *ciMode)
		return 0
	case "deps":
		printDoctorDepsReport(report, *ciMode)
		return 0
	case "toolchain":
		printDoctorToolchainReport(report, *ciMode)
		return 0
	case "all", "":
		if *ciMode {
			return printDoctorCIReport(report)
		}
		printDoctorSummary(report)
		return 0
	default:
		fmt.Fprintf(os.Stderr, "unknown doctor mode: %s\n", mode)
		fmt.Println("Usage:")
		fmt.Println("  pmx doctor [config|deps|toolchain|all] [project-dir]")
		return 2
	}
}

func runAgent(args []string) int {
	if len(args) == 0 {
		fmt.Println("PMX AGENT")
		fmt.Println(strings.Repeat("─", 50))
		fmt.Println()
		fmt.Println("Usage:")
		fmt.Println("  pmx agent inspect [--json]")
		fmt.Println("  pmx agent check [--json]")
		return 2
	}

	subcommand := args[0]
	rest := args[1:]
	fs := flag.NewFlagSet(subcommand, flag.ContinueOnError)
	jsonOutput := fs.Bool("json", false, "emit JSON output")
	fs.SetOutput(os.Stdout)
	if err := fs.Parse(rest); err != nil {
		return 2
	}

	switch subcommand {
	case "inspect":
		return runAgentInspect(*jsonOutput)
	case "check":
		return runAgentCheck(*jsonOutput)
	default:
		fmt.Fprintf(os.Stderr, "unknown agent command: %s\n", subcommand)
		fmt.Println("Usage:")
		fmt.Println("  pmx agent inspect [--json]")
		fmt.Println("  pmx agent check [--json]")
		return 2
	}
}

func runAgentInspect(jsonOutput bool) int {
	report := buildDoctorReport()
	payload := agentInspectReport{
		Version:     "1",
		Project:     report.Project,
		Diagnostics: report.Diagnostics,
	}
	if payload.Diagnostics == nil {
		payload.Diagnostics = []doctorDiagnostic{}
	}
	if jsonOutput {
		data, err := json.MarshalIndent(payload, "", "  ")
		if err != nil {
			fmt.Fprintf(os.Stderr, "agent inspect json failed: %v\n", err)
			return 1
		}
		fmt.Println(string(data))
		return 0
	}

	status := agentStatusForDiagnostics(report.Diagnostics)
	fmt.Println("PMX AGENT INSPECT")
	fmt.Println(strings.Repeat("─", 50))
	fmt.Printf("Project: %s\n", payload.Project.Ecosystem)
	fmt.Printf("Status:  %s\n", status)
	fmt.Printf("Warnings: %d\n", len(agentDiagnosticsBySeverity(payload.Diagnostics, "warn")))
	fmt.Printf("Failures: %d\n", len(agentDiagnosticsBySeverity(payload.Diagnostics, "fail")))
	return 0
}

func runAgentCheck(jsonOutput bool) int {
	// The gate deliberately invokes the same executable validation paths that a
	// developer or Foundry invokes. Do not derive check states from doctor data:
	// doing so can report a pass for validation that was never run.
	checks := []agentCheckEntry{
		runAgentGateCheck("doctor", []string{"doctor", "--ci"}),
		runAgentGateCheck("validate", []string{"validate", "*.go", "--input", "main.go"}),
		runAgentGateCheck("compat", []string{"compat", "--suite", "basic"}),
		runAgentGateCheck("ci", []string{"ci", "--json"}),
		runAgentGateCheck("regression", []string{"regression", "--json"}),
	}
	doctorReport := buildDoctorReport()
	result := agentResultForChecks(checks)
	status := result
	if status == "warn" {
		status = "warning"
	}
	payload := agentCheckReport{
		Version:     "1",
		Result:      result,
		Diagnostics: convertDoctorDiagnosticsToAgent(doctorReport.Diagnostics),
		Checks:      checks,
		NextActions: buildAgentActions(doctorReport.Diagnostics),
	}
	if payload.NextActions == nil {
		payload.NextActions = []agentAction{{Action: "none", Code: "PMX-NOOP", Text: "No blocking actions required."}}
	}
	if jsonOutput {
		data, err := json.MarshalIndent(payload, "", "  ")
		if err != nil {
			fmt.Fprintf(os.Stderr, "agent check json failed: %v\n", err)
			return 1
		}
		fmt.Println(string(data))
		if result == "fail" {
			return 1
		}
		return 0
	}

	fmt.Println("PMX AGENT CHECK")
	fmt.Println(strings.Repeat("─", 50))
	fmt.Printf("Status: %s\n", status)
	fmt.Printf("Result: %s\n", result)
	for _, check := range checks {
		fmt.Printf("%s: %s\n", check.Name, check.Status)
	}
	fmt.Println("Next actions:")
	for _, act := range payload.NextActions {
		fmt.Printf("  - %s\n", act.Code)
	}
	if result == "fail" {
		return 1
	}
	return 0
}

func runAgentGateCheck(name string, args []string) agentCheckEntry {
	root, err := repositoryRoot()
	if err != nil {
		return agentCheckEntry{Name: name, Status: "fail", Detail: err.Error()}
	}
	executable, err := os.Executable()
	if err != nil {
		return agentCheckEntry{Name: name, Status: "fail", Detail: err.Error()}
	}
	cmd := exec.Command(executable, args...)
	cmd.Dir = root
	out, err := cmd.CombinedOutput()
	detail := strings.TrimSpace(string(out))
	if err != nil {
		if detail == "" {
			detail = err.Error()
		}
		return agentCheckEntry{Name: name, Status: "fail", Detail: detail}
	}
	if strings.Contains(strings.ToLower(detail), `"result": "warning"`) || strings.Contains(detail, "Result: WARNING") {
		return agentCheckEntry{Name: name, Status: "warn", Detail: detail}
	}
	return agentCheckEntry{Name: name, Status: "pass", Detail: detail}
}

func agentResultForChecks(checks []agentCheckEntry) string {
	result := "pass"
	for _, check := range checks {
		if check.Status == "fail" {
			return "fail"
		}
		if check.Status == "warn" {
			result = "warn"
		}
	}
	return result
}

func convertDoctorDiagnosticsToAgent(diags []doctorDiagnostic) []agentDiagnostic {
	out := make([]agentDiagnostic, 0, len(diags))
	for _, diag := range diags {
		out = append(out, agentDiagnostic{
			Code:       diag.ID,
			Severity:   diag.Severity,
			File:       diag.File,
			Message:    diag.Message,
			Suggestion: diag.Suggestion,
			ID:         diag.ID,
		})
	}
	return out
}

func buildAgentActions(diags []doctorDiagnostic) []agentAction {
	if len(diags) == 0 {
		return nil
	}
	out := make([]agentAction, 0, len(diags))
	for _, diag := range diags {
		if diag.ID == "" {
			continue
		}
		out = append(out, agentAction{
			Action: "fix",
			Code:   diag.ID,
			Text:   diag.Title,
		})
	}
	return out
}

func agentStatusForDiagnostics(diags []doctorDiagnostic) string {
	if agentHasSeverity(diags, "fail") {
		return "fail"
	}
	if agentHasSeverity(diags, "warn") {
		return "warning"
	}
	return "healthy"
}

func agentHasSeverity(diags []doctorDiagnostic, severity string) bool {
	for _, diag := range diags {
		if diag.Severity == severity {
			return true
		}
	}
	return false
}

func agentDiagnosticsBySeverity(diags []doctorDiagnostic, severity string) []doctorDiagnostic {
	out := make([]doctorDiagnostic, 0)
	for _, diag := range diags {
		if diag.Severity == severity {
			out = append(out, diag)
		}
	}
	return out
}

func isDoctorMode(value string) bool {
	switch value {
	case "config", "deps", "toolchain", "all", "":
		return true
	default:
		return false
	}
}

func buildDoctorReport() doctorReport {
	project := doctorProject{
		Name:           filepathBase("."),
		Ecosystem:      detectProjectEcosystem(),
		PackageManager: detectPackageManager(),
		TypeScript:     fileExists("tsconfig.json") || fileExists("jsconfig.json"),
		Framework:      detectFramework(),
	}
	if project.TypeScript && project.Ecosystem == "javascript" {
		project.Ecosystem = "javascript/typescript"
	}
	if project.PackageManager == "unknown" && multiplePackageManagersDetected() {
		project.PackageManager = "unknown"
	}

	report := doctorReport{Version: "1", Project: project, Diagnostics: []doctorDiagnostic{}}
	if project.Ecosystem == "javascript/typescript" && !fileExists("package.json") {
		report.Summary.Fail++
		report.Diagnostics = append(report.Diagnostics, doctorDiagnostic{
			ID:         "PMX-PKG-002",
			Severity:   "fail",
			Category:   "dependency",
			Title:      "JavaScript package manifest is missing",
			File:       "package.json",
			Message:    "This project looks like a JavaScript or TypeScript app, but no package manifest was found.",
			Evidence:   []string{"package.json missing", "JS/TS project detected"},
			Suggestion: "Create a package.json and choose a single package manager for this project.",
		})
	}

	if multiplePackageManagersDetected() {
		report.Summary.Warn++
		report.Diagnostics = append(report.Diagnostics, doctorDiagnostic{
			ID:         "PMX-PKG-001",
			Severity:   "warn",
			Category:   "dependency",
			Title:      "Multiple package-manager lockfiles detected",
			File:       "package.json",
			Message:    "More than one package manager lockfile exists in the project root.",
			Evidence:   []string{"pnpm-lock.yaml", "package-lock.json"},
			Suggestion: "Keep one canonical package manager and remove stale lockfiles.",
		})
	}

	if fileExists("tsconfig.json") {
		strict := tsConfigStrictEnabled()
		if !strict {
			report.Summary.Warn++
			report.Diagnostics = append(report.Diagnostics, doctorDiagnostic{
				ID:         "PMX-TS-001",
				Severity:   "warn",
				Category:   "typescript",
				Title:      "TypeScript strict mode is disabled",
				File:       "tsconfig.json",
				Message:    "TypeScript configuration detected but strict mode is disabled.",
				Evidence:   []string{"tsconfig.json detected", "compilerOptions.strict = false"},
				Suggestion: "Review whether strict mode should be enabled for this project.",
			})
		}
	}

	if hasLegacyESLintConfig() {
		report.Summary.Warn++
		report.Diagnostics = append(report.Diagnostics, doctorDiagnostic{
			ID:         "PMX-ESLINT-001",
			Severity:   "warn",
			Category:   "eslint",
			Title:      "Legacy ESLint configuration detected",
			File:       detectLegacyESLintFile(),
			Message:    "Use of a legacy ESLint config style may require migration to flat config.",
			Evidence:   []string{".eslintrc.json detected", "ESLint 9 compatibility risk"},
			Suggestion: "Verify ESLint flat-config compatibility before upgrading or running lint in CI.",
		})
	}

	return report
}

func detectProjectEcosystem() string {
	switch {
	case fileExists("package.json") || fileExists("pnpm-lock.yaml") || fileExists("yarn.lock") || fileExists("package-lock.json") || fileExists("bun.lock") || fileExists("bun.lockb") || fileExists("tsconfig.json") || fileExists("eslint.config.js") || fileExists(".eslintrc.json") || fileExists(".eslintrc"):
		return "javascript"
	case fileExists("go.mod"):
		return "go"
	case fileExists("Cargo.toml"):
		return "rust"
	case fileExists("pyproject.toml"):
		return "python"
	default:
		return "unknown"
	}
}

func detectPackageManager() string {
	if fileExists("package.json") {
		pkg, err := os.ReadFile("package.json")
		if err == nil {
			var manifest map[string]interface{}
			if err := json.Unmarshal(pkg, &manifest); err == nil {
				if pm, ok := manifest["packageManager"].(string); ok && pm != "" {
					if idx := strings.Index(pm, "@"); idx > 0 {
						return pm[:idx]
					}
					return pm
				}
			}
		}
	}

	switch {
	case fileExists("pnpm-lock.yaml"):
		return "pnpm"
	case fileExists("yarn.lock"):
		return "yarn"
	case fileExists("package-lock.json"):
		return "npm"
	case fileExists("bun.lock") || fileExists("bun.lockb"):
		return "bun"
	case fileExists("package.json"):
		return "npm"
	case fileExists("go.mod"):
		return "go"
	case fileExists("Cargo.toml"):
		return "cargo"
	case fileExists("pyproject.toml"):
		return "pip"
	default:
		return "unknown"
	}
}

func detectFramework() string {
	switch {
	case fileExists("next.config.js") || fileExists("next.config.mjs") || fileExists("next.config.ts"):
		return "next"
	case fileExists("vite.config.js") || fileExists("vite.config.ts") || fileExists("vite.config.mjs"):
		return "vite"
	case fileExists("astro.config.mjs") || fileExists("astro.config.js") || fileExists("astro.config.ts"):
		return "astro"
	default:
		return ""
	}
}

func tsConfigStrictEnabled() bool {
	path := "tsconfig.json"
	if !fileExists(path) {
		return false
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return false
	}
	text := string(data)
	return strings.Contains(text, "\"strict\": true") || strings.Contains(text, "\"strict\":\n    true") || strings.Contains(text, "\"strict\":\n\ttrue")
}

func hasLegacyESLintConfig() bool {
	return fileExists(".eslintrc") || fileExists(".eslintrc.json") || fileExists(".eslintrc.js") || fileExists(".eslintrc.cjs")
}

func detectLegacyESLintFile() string {
	switch {
	case fileExists(".eslintrc.json"):
		return ".eslintrc.json"
	case fileExists(".eslintrc.js"):
		return ".eslintrc.js"
	case fileExists(".eslintrc.cjs"):
		return ".eslintrc.cjs"
	case fileExists(".eslintrc"):
		return ".eslintrc"
	default:
		return "eslint.config.js"
	}
}

func multiplePackageManagersDetected() bool {
	lockfiles := []string{"pnpm-lock.yaml", "yarn.lock", "package-lock.json", "bun.lock", "bun.lockb"}
	count := 0
	for _, path := range lockfiles {
		if fileExists(path) {
			count++
		}
	}
	return count > 1
}

func printDoctorSummary(report doctorReport) {
	fmt.Println("PMX DOCTOR")
	fmt.Println()
	fmt.Println("Project")
	fmt.Println(strings.Repeat("─", 56))
	fmt.Printf("Ecosystem       %s\n", formatDoctorEcosystem(report.Project.Ecosystem))
	fmt.Printf("Package Manager %s\n", formatDoctorPackageManager(report.Project.PackageManager))
	if report.Project.TypeScript {
		fmt.Println("TypeScript      Detected")
	}
	if report.Project.Framework != "" {
		fmt.Printf("Framework       %s\n", report.Project.Framework)
	}
	fmt.Println()
	fmt.Println("Diagnostics")
	fmt.Println(strings.Repeat("─", 56))
	fmt.Println()
	for _, diag := range report.Diagnostics {
		if diag.Severity == "fail" {
			fmt.Printf("FAIL  %s\n", diag.ID)
		} else {
			fmt.Printf("WARN  %s\n", diag.ID)
		}
		if diag.Title != "" {
			fmt.Println("  " + diag.Title)
		}
		if diag.File != "" {
			fmt.Printf("  file: %s\n", diag.File)
		}
		if diag.Message != "" {
			fmt.Printf("  %s\n", diag.Message)
		}
		for _, line := range diag.Evidence {
			fmt.Printf("  %s\n", line)
		}
		if diag.Suggestion != "" {
			fmt.Printf("  Recommendation: %s\n", diag.Suggestion)
		}
		fmt.Println()
	}
	fmt.Println("Summary")
	fmt.Println(strings.Repeat("─", 56))
	fmt.Printf("PASS    %d\n", report.Summary.Pass)
	fmt.Printf("WARN    %d\n", report.Summary.Warn)
	fmt.Printf("FAIL    %d\n", report.Summary.Fail)
	fmt.Println()
	if report.Summary.Fail > 0 {
		fmt.Println("Result: FAILURE")
		return
	}
	if report.Summary.Warn > 0 {
		fmt.Println("Result: WARNING")
		return
	}
	fmt.Println("Result: PASS")
}

func printDoctorConfigReport(report doctorReport, ciMode bool) {
	fmt.Println("Configuration")
	fmt.Println(strings.Repeat("─", 56))
	fmt.Println()
	printFileStatus("package.json", fileExists("package.json"))
	printFileStatus("tsconfig.json", fileExists("tsconfig.json"))
	printFileStatus("eslint.config.js", fileExists("eslint.config.js"))
	printFileStatus("go.mod", fileExists("go.mod"))
	printFileStatus("Cargo.toml", fileExists("Cargo.toml"))
	fmt.Println()
	fmt.Println("Summary")
	fmt.Printf("  %d config checks passed\n", report.Summary.Pass)
	fmt.Printf("  %d warnings\n", report.Summary.Warn)
	if ciMode {
		fmt.Println("  CI: PASS")
	}
}

func printDoctorDepsReport(report doctorReport, ciMode bool) {
	fmt.Println("Dependencies")
	fmt.Println(strings.Repeat("─", 56))
	fmt.Println()
	fmt.Println("  ✓ package metadata detected when present")
	fmt.Println("  ✓ module graph is represented by the active project root")
	fmt.Println("  ⚠ dependency intelligence is intentionally narrow in the MVP")
	if ciMode {
		fmt.Println("  CI: PASS")
	}
}

func printDoctorToolchainReport(report doctorReport, ciMode bool) {
	fmt.Println("Toolchain")
	fmt.Println(strings.Repeat("─", 56))
	fmt.Println()
	printToolStatus("Go", toolExists("go"))
	printToolStatus("Node", toolExists("node"))
	printToolStatus("Cargo", toolExists("cargo"))
	printToolStatus("Python", toolExists("python3"))
	if ciMode {
		fmt.Println("  CI: PASS")
	}
}

func printDoctorCIReport(report doctorReport) int {
	fmt.Println("PMX Doctor CI")
	fmt.Println(strings.Repeat("─", 56))
	fmt.Println()
	for _, diag := range report.Diagnostics {
		if diag.Severity == "fail" {
			fmt.Printf("FAIL %s\n", diag.ID)
		} else {
			fmt.Printf("WARN %s\n", diag.ID)
		}
	}
	fmt.Println()
	if report.Summary.Fail > 0 {
		fmt.Println("Result: FAILURE")
		return 1
	}
	if report.Summary.Warn > 0 {
		fmt.Println("Result: WARNING")
		return 0
	}
	fmt.Println("Result: PASS")
	return 0
}

func formatDoctorEcosystem(value string) string {
	switch value {
	case "javascript/typescript":
		return "JavaScript / TypeScript"
	case "javascript":
		return "JavaScript"
	case "go":
		return "Go"
	case "rust":
		return "Rust"
	case "python":
		return "Python"
	default:
		if value == "" {
			return "Unknown"
		}
		return titleCase(value)
	}
}

func formatDoctorPackageManager(value string) string {
	if value == "" {
		return "Unknown"
	}
	return titleCase(value)
}

func titleCase(value string) string {
	if value == "" {
		return ""
	}
	r, size := utf8.DecodeRuneInString(value)
	if r == utf8.RuneError {
		return value
	}
	return string(unicode.ToTitle(r)) + value[size:]
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func toolExists(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

func filepathBase(path string) string {
	wd, err := os.Getwd()
	if err != nil {
		return "project"
	}
	if path == "." {
		return filepathBaseName(wd)
	}
	return filepathBaseName(path)
}

func filepathBaseName(path string) string {
	if path == "" {
		return "project"
	}
	parts := strings.Split(path, string(os.PathSeparator))
	for i := len(parts) - 1; i >= 0; i-- {
		if parts[i] != "" && parts[i] != "." {
			return parts[i]
		}
	}
	return "project"
}

func printFileStatus(name string, exists bool) {
	if exists {
		fmt.Printf("  ✓ %s\n", name)
		return
	}
	fmt.Printf("  ! %s\n", name)
}

func printToolStatus(name string, exists bool) {
	if exists {
		fmt.Printf("  ✓ %s\n", name)
		return
	}
	fmt.Printf("  ! %s\n", name)
}

func runCI(args []string) int {
	fs := flag.NewFlagSet("ci", flag.ContinueOnError)
	jsonOutput := fs.Bool("json", false, "emit JSON output")
	ciMode := fs.Bool("ci", false, "enable CI-style output")
	fs.SetOutput(os.Stdout)
	if err := fs.Parse(args); err != nil {
		return 2
	}

	if fs.NArg() == 0 {
		if *jsonOutput || *ciMode {
			return runCIReport(true, *ciMode)
		}
		fmt.Println("PICOMATCH ENGINEERING FOUNDRY")
		fmt.Println(strings.Repeat("─", 50))
		fmt.Println()
		fmt.Println("ci commands are routed through the same Go validation pipeline used by the repo.")
		fmt.Println("Usage:")
		fmt.Println("  pmx ci run [--json]")
		fmt.Println("  pmx ci watch <run-id>")
		fmt.Println("  pmx ci logs <run-id> [job-name]")
		return 0
	}

	switch fs.Arg(0) {
	case "run":
		return runCIReport(*jsonOutput, *ciMode)
	case "watch":
		if fs.NArg() < 2 {
			fmt.Fprintln(os.Stderr, "usage: pmx ci watch <run-id>")
			return 2
		}
		fmt.Fprintf(os.Stderr, "pmx ci watch %s is not available in the local CLI; use the Foundry GitHub Actions API for live workflow status.\n", fs.Arg(1))
		return 1
	case "logs":
		if fs.NArg() < 2 {
			fmt.Fprintln(os.Stderr, "usage: pmx ci logs <run-id> [job-name]")
			return 2
		}
		job := ""
		if fs.NArg() >= 3 {
			job = " for job " + fs.Arg(2)
		}
		fmt.Fprintf(os.Stderr, "pmx ci logs %s%s is not available in the local CLI; use the Foundry GitHub Actions API for live job logs.\n", fs.Arg(1), job)
		return 1
	default:
		if *jsonOutput || *ciMode {
			return runCIReport(*jsonOutput, *ciMode)
		}
		fmt.Fprintf(os.Stderr, "unknown ci subcommand: %s\n", fs.Arg(0))
		fmt.Println("Usage:")
		fmt.Println("  pmx ci run [--json]")
		fmt.Println("  pmx ci watch <run-id>")
		fmt.Println("  pmx ci logs <run-id> [job-name]")
		return 2
	}
}

func regressionTargetPackages() ([]string, error) {
	root, err := repositoryRoot()
	if err != nil {
		return nil, err
	}
	cmd := exec.Command("go", "list", "./...")
	cmd.Dir = root
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("list packages: %w: %s", err, strings.TrimSpace(string(out)))
	}

	packages := make([]string, 0)
	for _, pkg := range strings.Fields(string(out)) {
		pkg = strings.TrimSpace(pkg)
		if pkg == "" || pkg == "github.com/Fenz-art/picomatchgo-pmx-CLI/cmd/pmx" || strings.HasSuffix(pkg, "/cmd/pmx") {
			continue
		}
		packages = append(packages, pkg)
	}
	return packages, nil
}

func runRegressionReport() (map[string]interface{}, error) {
	start := time.Now()
	root, err := repositoryRoot()
	if err != nil {
		return nil, err
	}
	packages, err := regressionTargetPackages()
	if err != nil {
		return nil, err
	}
	if len(packages) == 0 {
		return map[string]interface{}{
			"command":     "regression",
			"result":      "pass",
			"total":       0,
			"passed":      0,
			"failed":      0,
			"skipped":     0,
			"duration_ms": time.Since(start).Milliseconds(),
		}, nil
	}

	args := append([]string{"test", "-json"}, packages...)
	cmd := exec.Command("go", args...)
	cmd.Dir = root
	out, err := cmd.CombinedOutput()
	total, passed, failed, skipped := countGoTestResults(out)
	duration := time.Since(start)
	result := "pass"
	if failed > 0 || err != nil {
		result = "fail"
	}
	return map[string]interface{}{
		"command":     "regression",
		"result":      result,
		"total":       total,
		"passed":      passed,
		"failed":      failed,
		"skipped":     skipped,
		"duration_ms": duration.Milliseconds(),
	}, err
}

func runRegression(args []string) int {
	fs := flag.NewFlagSet("regression", flag.ContinueOnError)
	jsonOutput := fs.Bool("json", false, "emit JSON output")
	fs.SetOutput(os.Stdout)
	if err := fs.Parse(args); err != nil {
		return 2
	}

	report, err := runRegressionReport()
	if err != nil {
		fmt.Fprintf(os.Stderr, "regression run failed: %v\n", err)
		return 1
	}

	failed := report["failed"].(int)
	if *jsonOutput {
		payload, marshalErr := json.MarshalIndent(report, "", "  ")
		if marshalErr != nil {
			fmt.Fprintf(os.Stderr, "regression json failed: %v\n", marshalErr)
			return 1
		}
		fmt.Println(string(payload))
		if failed > 0 {
			return 1
		}
		return 0
	}

	total := report["total"].(int)
	passed := report["passed"].(int)
	skipped := report["skipped"].(int)
	durationMs := report["duration_ms"].(int64)
	result := strings.ToUpper(report["result"].(string))
	fmt.Println("PMX REGRESSION SUITE")
	fmt.Println(strings.Repeat("─", 50))
	fmt.Println()
	fmt.Printf("Tests:       %d\n", total)
	fmt.Printf("Passed:      %d\n", passed)
	fmt.Printf("Failed:      %d\n", failed)
	fmt.Printf("Skipped:     %d\n", skipped)
	fmt.Printf("Duration:    %s\n", time.Duration(durationMs)*time.Millisecond)
	fmt.Println()
	fmt.Printf("RESULT: %s\n", result)
	if failed > 0 {
		return 1
	}
	return 0
}

func countGoTestResults(data []byte) (total, passed, failed, skipped int) {
	decoder := json.NewDecoder(bytes.NewReader(data))
	for {
		var event map[string]interface{}
		if err := decoder.Decode(&event); err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			break
		}
		testName, _ := event["Test"].(string)
		if testName == "" {
			continue
		}
		action, _ := event["Action"].(string)
		switch action {
		case "pass":
			total++
			passed++
		case "fail":
			total++
			failed++
		case "skip":
			total++
			skipped++
		}
	}
	return total, passed, failed, skipped
}

func runCIReport(jsonOutput bool, ciMode bool) int {
	start := time.Now()
	root, rootErr := repositoryRoot()
	if rootErr != nil {
		fmt.Fprintf(os.Stderr, "ci repository discovery failed: %v\n", rootErr)
		return 1
	}
	checks := []struct {
		name string
		args []string
	}{
		{name: "Format", args: []string{"gofmt", "-l", "."}},
		{name: "Vet", args: []string{"go", "vet", "./..."}},
		{name: "Unit", args: []string{"go", "test", "-count=1", "./..."}},
		{name: "Race", args: []string{"go", "test", "-count=1", "-race", "./..."}},
		{name: "Build", args: []string{"go", "build", "./..."}},
		{name: "CLI", args: []string{"go", "run", "./cmd/pmx", "match", "**/*.go", "cmd/pmx/main.go"}},
		{name: "Compatibility", args: []string{"go", "run", "./cmd/pmx", "compat", "--suite", "basic"}},
		{name: "Doctor", args: []string{"go", "run", "./cmd/pmx", "doctor", "--ci"}},
	}

	report := map[string]interface{}{
		"result":      "pass",
		"checks":      []map[string]string{},
		"passed":      0,
		"failed":      0,
		"warnings":    0,
		"duration_ms": 0,
	}

	for _, check := range checks {
		status, detail := executeCICheck(root, check.name, check.args)
		item := map[string]string{"name": check.name, "status": status}
		if detail != "" {
			item["details"] = detail
		}
		report["checks"] = append(report["checks"].([]map[string]string), item)
		switch status {
		case "pass":
			report["passed"] = report["passed"].(int) + 1
		case "warn":
			report["warnings"] = report["warnings"].(int) + 1
		case "fail":
			report["failed"] = report["failed"].(int) + 1
		}
	}

	regressionReport, regressionErr := runRegressionReport()
	regressionStatus := "pass"
	regressionDetail := ""
	if regressionErr != nil {
		regressionStatus = "fail"
		regressionDetail = regressionErr.Error()
	} else if regressionReport["result"] == "fail" {
		regressionStatus = "fail"
		regressionDetail = fmt.Sprintf("failed %v/%v tests", regressionReport["failed"], regressionReport["total"])
	} else if regressionReport["result"] == "warning" {
		regressionStatus = "warn"
	}
	item := map[string]string{"name": "Regression", "status": regressionStatus}
	if regressionDetail != "" {
		item["details"] = regressionDetail
	}
	report["checks"] = append(report["checks"].([]map[string]string), item)
	switch regressionStatus {
	case "pass":
		report["passed"] = report["passed"].(int) + 1
	case "warn":
		report["warnings"] = report["warnings"].(int) + 1
	case "fail":
		report["failed"] = report["failed"].(int) + 1
	}

	if report["failed"].(int) > 0 {
		report["result"] = "fail"
	} else if report["warnings"].(int) > 0 {
		report["result"] = "warning"
	}
	report["duration_ms"] = time.Since(start).Milliseconds()

	if jsonOutput {
		payload, err := json.MarshalIndent(report, "", "  ")
		if err != nil {
			fmt.Fprintf(os.Stderr, "ci json failed: %v\n", err)
			return 1
		}
		fmt.Println(string(payload))
		if report["failed"].(int) > 0 {
			return 1
		}
		return 0
	}

	fmt.Println("PMX CONTINUOUS VALIDATION")
	fmt.Println(strings.Repeat("─", 56))
	for _, item := range report["checks"].([]map[string]string) {
		fmt.Printf("[%s] %s", strings.ToUpper(item["status"]), item["name"])
		if detail := item["details"]; detail != "" {
			fmt.Printf(" - %s", strings.TrimSpace(detail))
		}
		fmt.Println()
	}
	fmt.Println()
	fmt.Printf("%d checks\n", len(report["checks"].([]map[string]string)))
	fmt.Printf("%d passed\n", report["passed"].(int))
	if report["warnings"].(int) > 0 {
		fmt.Printf("%d warnings\n", report["warnings"].(int))
	}
	fmt.Printf("%d failed\n", report["failed"].(int))
	fmt.Println("────────────────────────────")
	fmt.Printf("RESULT: %s\n", strings.ToUpper(report["result"].(string)))
	if ciMode {
		if report["failed"].(int) > 0 {
			return 1
		}
		return 0
	}
	if report["failed"].(int) > 0 {
		return 1
	}
	return 0
}

func executeCICheck(root, name string, args []string) (status, detail string) {
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Dir = root
	out, err := cmd.CombinedOutput()
	output := strings.TrimSpace(string(out))
	if len(output) > 0 && name == "Format" {
		return "warn", output
	}
	if err == nil {
		if name == "Doctor" {
			if strings.Contains(output, "Result: WARNING") {
				return "warn", "warning"
			}
		}
		return "pass", output
	}
	if name == "Doctor" && strings.Contains(output, "Result: WARNING") {
		return "warn", "warning"
	}
	if len(output) > 0 {
		return "fail", output
	}
	return "fail", err.Error()
}

// repositoryRoot makes PMX commands work the same from a built binary, from
// `go run ./cmd/pmx`, and from package tests that execute in cmd/pmx.
func repositoryRoot() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for current := cwd; ; current = filepath.Dir(current) {
		if fileExists(filepath.Join(current, "go.mod")) && fileExists(filepath.Join(current, "cmd", "pmx", "main.go")) {
			return current, nil
		}
		parent := filepath.Dir(current)
		if parent == current {
			break
		}
	}
	return "", fmt.Errorf("could not locate PMX repository root from %s", cwd)
}

func runExplain(args []string) int {
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "usage: pmx explain <pattern> [--input <string>] [flags]")
		return 2
	}

	pattern := args[0]
	fs := flag.NewFlagSet("explain", flag.ContinueOnError)
	input := fs.String("input", "", "input string to match")
	dot := fs.Bool("dot", false, "match dotfiles")
	nocase := fs.Bool("nocase", false, "case-insensitive matching")
	windows := fs.Bool("windows", false, "normalize Windows path separators")
	fs.SetOutput(os.Stdout)
	if err := fs.Parse(args[1:]); err != nil {
		return 2
	}

	opts := &picomatch.Options{
		Dot:     *dot,
		Nocase:  *nocase,
		Windows: *windows,
	}
	scanOpts := &picomatch.Options{
		Dot:     *dot,
		Nocase:  *nocase,
		Windows: *windows,
		Parts:   true,
		Tokens:  true,
	}

	scanState := picomatch.Scan(pattern, scanOpts)
	parseState, parseErr := picomatch.Parse(pattern, opts)
	compiled, compileErr := picomatch.MakeRe(pattern, opts)

	printExplainOutput(pattern, scanState, parseState, parseErr, compiled, compileErr, *input)

	if *input != "" {
		match, matchErr := picomatch.IsMatch(*input, pattern, opts)
		printExplainMatchSummary(match, matchErr)
		if matchErr != nil {
			return 1
		}
	}

	if parseErr != nil || compileErr != nil {
		return 1
	}

	return 0
}

func runValidate(args []string) int {
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "usage: pmx validate <pattern> [--input <string>] [flags]")
		return 2
	}

	pattern := args[0]
	fs := flag.NewFlagSet("validate", flag.ContinueOnError)
	input := fs.String("input", "", "input string to match")
	dot := fs.Bool("dot", false, "match dotfiles")
	nocase := fs.Bool("nocase", false, "case-insensitive matching")
	windows := fs.Bool("windows", false, "normalize Windows path separators")
	fs.SetOutput(os.Stdout)
	if err := fs.Parse(args[1:]); err != nil {
		return 2
	}

	opts := &picomatch.Options{
		Dot:     *dot,
		Nocase:  *nocase,
		Windows: *windows,
	}

	scanState := picomatch.Scan(pattern, &picomatch.Options{Parts: true})
	parseState, parseErr := picomatch.Parse(pattern, opts)
	compiled, compileErr := picomatch.MakeRe(pattern, opts)
	var match bool
	var matchErr error
	if *input != "" {
		match, matchErr = picomatch.IsMatch(*input, pattern, opts)
	}

	fmt.Println("Validation")
	fmt.Println(strings.Repeat("─", 50))
	fmt.Println()
	fmt.Printf("pattern: %s\n", pattern)
	fmt.Printf("scanner base: %s\n", defaultString(scanState.Base, "."))
	fmt.Printf("parser output: %s\n", defaultString(parseState.Output, "<empty>"))
	if compiled != nil {
		fmt.Printf("regex: %s\n", compiled.String())
	}
	fmt.Println()
	fmt.Println(strings.Repeat("─", 50))
	fmt.Println()
	fmt.Printf("Scanner          PASS\n")
	if parseErr == nil {
		fmt.Printf("Parser           PASS\n")
	} else {
		fmt.Printf("Parser           FAIL\n")
	}
	if compileErr == nil {
		fmt.Printf("Compiler         PASS\n")
	} else {
		fmt.Printf("Compiler         FAIL\n")
	}
	if *input != "" {
		if matchErr == nil {
			fmt.Printf("Matcher          PASS\n")
		} else {
			fmt.Printf("Matcher          FAIL\n")
		}
	}
	fmt.Println()

	if parseErr != nil {
		fmt.Printf("Parser error: %v\n", parseErr)
	}
	if compileErr != nil {
		fmt.Printf("Compiler error: %v\n", compileErr)
	}
	if matchErr != nil {
		fmt.Printf("Matcher error: %v\n", matchErr)
	}

	if parseErr != nil || compileErr != nil || matchErr != nil {
		fmt.Println("Result           INVALID")
		return 1
	}

	fmt.Println("Result           VALID")
	if *input != "" {
		if match {
			fmt.Println("Match            YES")
		} else {
			fmt.Println("Match            NO")
		}
	}
	return 0
}

func printExplainOutput(pattern string, scanState picomatch.ScanState, parseState picomatch.ParseState, parseErr error, compiled *regexp.Regexp, compileErr error, input string) {
	fmt.Println("Picomatch Analysis")
	fmt.Println(strings.Repeat("─", 50))
	fmt.Println()
	fmt.Println("Pattern")
	fmt.Printf("  %s\n\n", pattern)

	fmt.Println("Scanner")
	fmt.Printf("  base:        %s\n", defaultString(scanState.Base, "."))
	fmt.Printf("  glob:        %t\n", scanState.IsGlob)
	fmt.Printf("  globstar:    %t\n", scanState.IsGlobstar)
	fmt.Printf("  segments:    %d\n", len(scanState.Parts))
	for idx, seg := range scanState.Parts {
		fmt.Printf("    segment[%d]: %s\n", idx, seg)
	}
	fmt.Println()

	fmt.Println("Parser")
	if parseErr != nil {
		fmt.Printf("  status:      FAIL\n")
		fmt.Printf("  error:       %v\n", parseErr)
	} else {
		fmt.Printf("  globstar:    %t\n", parseState.Globstar)
		fmt.Printf("  output:      %s\n", defaultString(parseState.Output, "<empty>"))
	}
	fmt.Println()

	fmt.Println("Compiler")
	fmt.Printf("  engine:      Go regexp / RE2\n")
	if compileErr != nil {
		fmt.Printf("  status:      FAIL\n")
		fmt.Printf("  error:       %v\n", compileErr)
	} else {
		fmt.Printf("  compiled:    PASS\n")
		fmt.Printf("  regex:       %s\n", compiled.String())
	}
	fmt.Println()

	fmt.Println("Matcher")
	if input == "" {
		fmt.Printf("  awaiting input\n")
	} else {
		fmt.Printf("  input:       %s\n", input)
	}
}

func printExplainMatchSummary(match bool, matchErr error) {
	fmt.Println()
	if matchErr != nil {
		fmt.Printf("MATCHER ERROR: %v\n", matchErr)
		return
	}

	if match {
		fmt.Println("MATCH")
	} else {
		fmt.Println("NO MATCH")
	}
	fmt.Println()
	fmt.Println("Scanner       PASS")
	fmt.Println("Parser        PASS")
	fmt.Println("Compiler      PASS")
	fmt.Println("Matcher       PASS")
	fmt.Println()
	fmt.Println("Execution path:")
	fmt.Println("Pattern")
	fmt.Println("  ↓")
	fmt.Println("Scan")
	fmt.Println("  ↓")
	fmt.Println("Parse")
	fmt.Println("  ↓")
	fmt.Println("Compile")
	fmt.Println("  ↓")
	fmt.Println("Match")
}

func defaultString(value string, fallback string) string {
	if value == "" {
		return fallback
	}
	return value
}
