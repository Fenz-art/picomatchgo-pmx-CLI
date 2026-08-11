package javascript

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/Fenz-art/picomatchgo-pmx-CLI/ecosystem/core"
)

// JSAdapter provides a minimal JavaScript/TypeScript adapter contract shape.
type JSAdapter struct{}

func (j *JSAdapter) Name() string { return "javascript" }

func (j *JSAdapter) Detect(dir string) bool {
	if fileInDir(dir, "package.json") {
		return true
	}
	return false
}

func (j *JSAdapter) Inspect(dir string) core.EcosystemInfo {
	info := core.EcosystemInfo{
		Ecosystem:      "javascript/typescript",
		PackageManager: "pnpm",
		Language:       "typescript",
		Detected:       []string{"package.json"},
	}
	if fileInDir(dir, "package.json") {
		info.Dependencies.Manifest = "package.json"
	}
	if fileInDir(dir, "pnpm-lock.yaml") {
		info.Dependencies.Lockfile = "pnpm-lock.yaml"
	}
	return info
}

func (j *JSAdapter) ValidateEnvironment(dir string) []core.Diagnostic   { return nil }
func (j *JSAdapter) ValidateDependencies(dir string) []core.Diagnostic  { return nil }
func (j *JSAdapter) ValidateConfiguration(dir string) []core.Diagnostic { return nil }
func (j *JSAdapter) ValidateToolchain(dir string) []core.Diagnostic     { return nil }
func (j *JSAdapter) ValidateProject(dir string) []core.Diagnostic       { return nil }

func fileInDir(dir, name string) bool {
	_, err := os.Stat(filepath.Join(dir, name))
	return err == nil
}

// execLookPath wraps exec.LookPath for testability.
var execLookPath = exec.LookPath

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
			if strings.HasPrefix(line, prefix) {
				line = strings.TrimPrefix(line, prefix)
			}
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
