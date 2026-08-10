package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	picomatch "github.com/Fenz-art/picomatchgo-pmx-CLI"
)

type compatCase struct {
	Pattern string          `json:"pattern"`
	Input   string          `json:"input"`
	Options map[string]bool `json:"options,omitempty"`
	Expect  bool            `json:"expect"`
}

type compatSuite struct {
	Name  string       `json:"name"`
	Cases []compatCase `json:"cases"`
}

func runCompat(args []string) int {
	fs := flag.NewFlagSet("compat", flag.ContinueOnError)
	suiteName := fs.String("suite", "basic", "compatibility suite to run")
	dot := fs.Bool("dot", false, "match dotfiles")
	nocase := fs.Bool("nocase", false, "case-insensitive matching")
	windows := fs.Bool("windows", false, "normalize Windows path separators")
	fs.SetOutput(os.Stdout)
	if err := fs.Parse(args); err != nil {
		return 2
	}

	suiteFile, err := findCompatSuiteFile(*suiteName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error locating suite file: %v\n", err)
		return 1
	}

	suite, err := loadCompatSuite(suiteFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading suite: %v\n", err)
		return 1
	}

	fmt.Printf("Compatibility Suite: %s\n", suite.Name)
	fmt.Println(strings.Repeat("─", 50))
	fmt.Println()

	pass := 0
	fail := 0
	for _, c := range suite.Cases {
		opts := &picomatch.Options{
			Dot:     *dot || c.Options["dot"],
			Nocase:  *nocase || c.Options["nocase"],
			Windows: *windows || c.Options["windows"],
		}

		result, err := picomatch.IsMatch(c.Input, c.Pattern, opts)
		expected := c.Expect
		if err != nil {
			fmt.Printf("FAIL %-30s input=%-20s error=%v\n", c.Pattern, c.Input, err)
			fail++
			continue
		}

		if result == expected {
			pass++
		} else {
			fmt.Printf("MISMATCH %-30s input=%-20s got=%v want=%v\n", c.Pattern, c.Input, result, expected)
			fail++
		}
	}

	fmt.Println()
	fmt.Printf("Cases: %d\n", len(suite.Cases))
	fmt.Printf("Pass:  %d\n", pass)
	fmt.Printf("Fail:  %d\n", fail)

	if fail > 0 {
		fmt.Println("Behavior: MISMATCH")
		return 1
	}

	fmt.Println("Behavior: EQUIVALENT")
	return 0
}

func findCompatSuiteFile(name string) (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	current := cwd
	for i := 0; i < 4; i++ {
		candidate := filepath.Join(current, "test", "compatibility", name+".json")
		if _, err := os.Stat(candidate); err == nil {
			return candidate, nil
		}
		parent := filepath.Dir(current)
		if parent == current {
			break
		}
		current = parent
	}

	return "", fmt.Errorf("compatibility suite not found: %s", name)
}

func loadCompatSuite(path string) (*compatSuite, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var suite compatSuite
	if err := json.Unmarshal(data, &suite); err != nil {
		return nil, err
	}

	return &suite, nil
}
