package javascript

import (
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/Fenz-art/picomatchgo-pmx-CLI/ecosystem/core"
)

// Adapter implements validation for JavaScript and TypeScript projects.
type Adapter struct{}

func NewJavaScriptAdapter() *Adapter { return &Adapter{} }

func (a *Adapter) Name() string { return "javascript" }

func (a *Adapter) Detect(dir string) bool {
	for _, name := range []string{
		"package.json",
		"pnpm-lock.yaml",
		"yarn.lock",
		"package-lock.json",
		"bun.lock",
		"bun.lockb",
		"tsconfig.json",
		"jsconfig.json",
		"eslint.config.js",
		".eslintrc",
		".eslintrc.json",
	} {
		if _, err := os.Stat(filepath.Join(dir, name)); err == nil {
			return true
		}
	}
	return false
}

func (a *Adapter) Inspect(dir string) core.EcosystemInfo {
	info := core.EcosystemInfo{
		Ecosystem: "javascript",
		Language:  "javascript",
		Name:      "javascript",
		Root:      core.NormalizeDir(dir),
		Details:   map[string]string{},
	}
	if _, err := os.Stat(filepath.Join(dir, "tsconfig.json")); err == nil {
		info.Ecosystem = "javascript/typescript"
		info.Language = "typescript"
		info.Name = "javascript/typescript"
		info.Details["typescript"] = "true"
	}
	for _, file := range []string{"package.json", "pnpm-lock.yaml", "yarn.lock", "package-lock.json", "bun.lock", "bun.lockb"} {
		if _, err := os.Stat(filepath.Join(dir, file)); err == nil {
			info.Manifest = filepath.Join(dir, file)
			break
		}
	}
	if info.Manifest == "" {
		info.Manifest = filepath.Join(dir, "package.json")
	}
	if data, err := os.ReadFile(filepath.Join(dir, "package.json")); err == nil {
		var manifest map[string]any
		if err := json.Unmarshal(data, &manifest); err == nil {
			if pm, ok := manifest["packageManager"].(string); ok && pm != "" {
				if idx := strings.Index(pm, "@"); idx > 0 {
					info.PackageManager = pm[:idx]
				} else {
					info.PackageManager = pm
				}
			}
			if name, ok := manifest["name"].(string); ok && name != "" {
				info.Details["package_name"] = name
			}
		}
	}
	for _, file := range []string{"pnpm-lock.yaml", "yarn.lock", "package-lock.json", "bun.lock", "bun.lockb"} {
		if _, err := os.Stat(filepath.Join(dir, file)); err == nil {
			pm := strings.TrimSuffix(filepath.Base(file), filepath.Ext(file))
			if strings.HasSuffix(pm, "-lock") {
				pm = strings.TrimSuffix(pm, "-lock")
			}
			info.PackageManager = pm
			break
		}
	}
	if info.PackageManager == "" {
		if _, err := os.Stat(filepath.Join(dir, "package.json")); err == nil {
			info.PackageManager = "npm"
		}
	}
	if files, err := os.ReadDir(dir); err == nil {
		info.FileCount = len(files)
	}
	return info
}

func (a *Adapter) ValidateEnvironment(dir string) []core.Diagnostic {
	if _, err := os.Stat(filepath.Join(dir, "package.json")); err != nil {
		return nil
	}
	if _, err := execLookPath("node"); err != nil {
		return []core.Diagnostic{{
			ID:         "PMX-ENV-001",
			Severity:   "fail",
			Category:   "environment",
			Title:      "Node.js is not installed",
			File:       "package.json",
			Message:    "A JavaScript project was detected, but Node.js is not available on PATH.",
			Suggestion: "Install Node.js and a supported package manager before running JavaScript validation.",
		}}
	}
	return nil
}

func (a *Adapter) ValidateDependencies(dir string) []core.Diagnostic {
	variations := 0
	for _, name := range []string{"pnpm-lock.yaml", "yarn.lock", "package-lock.json", "bun.lock", "bun.lockb"} {
		if _, err := os.Stat(filepath.Join(dir, name)); err == nil {
			variations++
		}
	}
	if variations > 1 {
		return []core.Diagnostic{{
			ID:         "PMX-PKG-001",
			Severity:   "warn",
			Category:   "dependencies",
			Title:      "Multiple lockfiles detected",
			File:       "package.json",
			Message:    "More than one package manager lockfile exists in the project root.",
			Suggestion: "Keep a single canonical lockfile to avoid dependency drift.",
		}}
	}
	return nil
}

func (a *Adapter) ValidateConfiguration(dir string) []core.Diagnostic {
	var diags []core.Diagnostic
	if _, err := os.Stat(filepath.Join(dir, "tsconfig.json")); err == nil {
		data, err := os.ReadFile(filepath.Join(dir, "tsconfig.json"))
		if err == nil && !regexp.MustCompile(`(?i)"strict"\s*:\s*(true|1)`).Match(data) {
			diags = append(diags, core.Diagnostic{
				ID:         "PMX-TS-001",
				Severity:   "warn",
				Category:   "configuration",
				Title:      "TypeScript strict mode is disabled",
				File:       "tsconfig.json",
				Message:    "TypeScript strict mode is not enabled, reducing static safety.",
				Suggestion: "Enable strict mode in tsconfig.json for stronger checks.",
			})
		}
	}
	for _, name := range []string{".eslintrc", ".eslintrc.json"} {
		if _, err := os.Stat(filepath.Join(dir, name)); err == nil {
			diags = append(diags, core.Diagnostic{
				ID:         "PMX-ESLINT-001",
				Severity:   "warn",
				Category:   "configuration",
				Title:      "Legacy ESLint configuration detected",
				File:       name,
				Message:    "Legacy ESLint configuration may require migration to flat config.",
				Suggestion: "Review the ESLint config and update to the flat-config format if needed.",
			})
			break
		}
	}
	return diags
}

func (a *Adapter) ValidateToolchain(dir string) []core.Diagnostic {
	if _, err := os.Stat(filepath.Join(dir, "package.json")); err != nil {
		return nil
	}
	data, err := os.ReadFile(filepath.Join(dir, "package.json"))
	if err != nil {
		return nil
	}
	var manifest map[string]any
	if err := json.Unmarshal(data, &manifest); err != nil {
		return nil
	}
	engines, ok := manifest["engines"].(map[string]any)
	if !ok {
		return nil
	}
	if nodeConstraint, ok := engines["node"].(string); ok && nodeConstraint != "" {
		if _, err := execLookPath("node"); err == nil {
			if out, err := execCommand("node", "-p", "process.versions.node"); err == nil {
				ok, compareErr := compareNodeVersion(strings.TrimSpace(out), nodeConstraint)
				if compareErr == nil && !ok {
					return []core.Diagnostic{{
						ID:         "PMX-TOOL-001",
						Severity:   "warn",
						Category:   "toolchain",
						Title:      "Node.js engine constraint mismatch",
						File:       "package.json",
						Message:    "The installed Node.js version does not satisfy the configured engine constraint.",
						Suggestion: "Install a Node.js version that satisfies the package engine requirement.",
					}}
				}
			}
		}
	}
	return nil
}

func (a *Adapter) ValidateProject(dir string) []core.Diagnostic { return nil }
