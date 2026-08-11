package javascript

import (
	"context"
	"os/exec"
	"strings"
	"time"
)

var execLookPathImpl = exec.LookPath

func execLookPath(name string) (string, error) { return execLookPathImpl(name) }

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
