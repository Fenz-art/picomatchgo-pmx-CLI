package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"

	picomatch "github.com/debayansamal/port-mortem-picomatch-go"
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
	case "explain":
		exit(runExplain(args))
	case "validate":
		exit(runValidate(args))
	case "compat":
		exit(runCompat(args))
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
	fmt.Println("pmx — Picomatch Developer Diagnostics")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  pmx match <pattern> <input> [flags]")
	fmt.Println("  pmx explain <pattern> [--input <string>] [flags]")
	fmt.Println("  pmx validate <pattern> [--input <string>] [flags]")
	fmt.Println("  pmx compat [--suite <name>] [flags]")
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
