package javascript

import (
	"context"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

var execLookPathImpl = exec.LookPath

func execLookPath(name string) (string, error) { return execLookPathImpl(name) }

func execCommand(name string, args ...string) (string, error) {
	out, err := exec.Command(name, args...).CombinedOutput()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func compareNodeVersion(actual, constraint string) (bool, error) {
	actual = strings.TrimSpace(actual)
	if actual == "" {
		return false, fmt.Errorf("invalid node version: %q", actual)
	}
	actual = strings.TrimPrefix(actual, "v")
	actual = strings.TrimPrefix(actual, "V")
	if idx := strings.IndexAny(actual, " .-"); idx >= 0 {
		actual = actual[:idx]
	}
	major, err := strconv.Atoi(actual)
	if err != nil {
		return false, err
	}
	constraint = strings.TrimSpace(constraint)
	if constraint == "" {
		return true, nil
	}
	if strings.HasPrefix(constraint, ">=") {
		min, err := strconv.Atoi(strings.TrimPrefix(constraint, ">="))
		if err != nil {
			return false, err
		}
		return major >= min, nil
	}
	if strings.HasPrefix(constraint, ">") {
		min, err := strconv.Atoi(strings.TrimPrefix(constraint, ">"))
		if err != nil {
			return false, err
		}
		return major > min, nil
	}
	if strings.HasPrefix(constraint, "<=") {
		max, err := strconv.Atoi(strings.TrimPrefix(constraint, "<="))
		if err != nil {
			return false, err
		}
		return major <= max, nil
	}
	if strings.HasPrefix(constraint, "<") {
		max, err := strconv.Atoi(strings.TrimPrefix(constraint, "<"))
		if err != nil {
			return false, err
		}
		return major < max, nil
	}
	if constraint == "*" || constraint == "x" || constraint == "X" {
		return true, nil
	}
	return major == 0 || constraint == "*", nil
}

// getToolVersionExec executes `name --version` and returns the version string.
func getToolVersionExec(name string) string {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, name, "--version")
	out, err := cmd.CombinedOutput()
	if err != nil {
		cmd = exec.CommandContext(ctx, name, "-v")
		out, err = cmd.CombinedOutput()
		if err != nil {
			return ""
		}
	}

	text := strings.TrimSpace(string(out))
	lines := strings.Split(text, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		for _, prefix := range []string{"v", "Version: ", "TypeScript "} {
			line = strings.TrimPrefix(line, prefix)
		}
		if len(line) > 0 && line[0] >= '0' && line[0] <= '9' && strings.Contains(line, ".") {
			return line
		}
	}
	if len(lines) > 0 {
		return lines[0]
	}
	return text
}

var getToolVersionFunc = func(name string) string { return "" }

func init() {
	getToolVersionFunc = getToolVersionExec
}
