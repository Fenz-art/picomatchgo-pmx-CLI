package javascript

import (
	"context"
	"os/exec"
	"strings"
	"time"
)

// execLookPath wraps exec.LookPath for testability.
var execLookPath = exec.LookPath

// getToolVersionExec executes `name --version` and returns the version string.
func getToolVersionExec(name string) string {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var cmd *exec.Cmd
	switch name {
	case "tsc":
		cmd = exec.CommandContext(ctx, name, "--version")
	default:
		cmd = exec.CommandContext(ctx, name, "--version")
	}

	out, err := cmd.CombinedOutput()
	if err != nil {
		// Try -v fallback
		cmd = exec.CommandContext(ctx, name, "-v")
		out, err = cmd.CombinedOutput()
		if err != nil {
			return ""
		}
	}

	// Parse version from output
	text := strings.TrimSpace(string(out))
	// Common patterns:
	//   v22.22.1
	//   9.15.4
	//   Version: 5.8.3
	//   tsc --version => TypeScript 5.8.3
	lines := strings.Split(text, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		// Strip common prefixes
		for _, prefix := range []string{"v", "Version: ", "TypeScript "} {
			if strings.HasPrefix(line, prefix) {
				line = strings.TrimPrefix(line, prefix)
			}
		}
		// Check if it looks like a version number
		if isVersionString(line) {
			return line
		}
	}

	// Return the first line as a best-effort
	if len(lines) > 0 {
		return lines[0]
	}
	return text
}

func isVersionString(s string) bool {
	if s == "" {
		return false
	}
	// Must start with a digit
	if s[0] < '0' || s[0] > '9' {
		return false
	}
	// Must contain at least one dot
	return strings.Contains(s, ".")
}

func init() {
	// Override the getToolVersion method in the adapter
	// This connects the pure-logic adapter.go to the exec-based version detection
	getToolVersionFunc = getToolVersionExec
}

// getToolVersionFunc is the function pointer used by JSAdapter.getToolVersion
var getToolVersionFunc = func(name string) string { return "" }
