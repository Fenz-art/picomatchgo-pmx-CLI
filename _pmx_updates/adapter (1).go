package javascript

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Fenz-art/picomatchgo-pmx-CLI/ecosystem/core"
)

// JSAdapter provides deep JavaScript/TypeScript project validation.
// This is the first genuinely deep ecosystem adapter for PMX.
type JSAdapter struct{}

func (j *JSAdapter) Name() string { return "javascript" }

// Detect checks for JS/TS ecosystem sentinel files.
func (j *JSAdapter) Detect(dir string) bool {
	sentinels := []string{
		"package.json", "pnpm-lock.yaml", "yarn.lock",
		"package-lock.json", "bun.lock", "bun.lockb",
		"tsconfig.json", "jsconfig.json",
		"eslint.config.js", "eslint.config.mjs", "eslint.config.ts",
		".eslintrc", ".eslintrc.json", ".eslintrc.js", ".eslintrc.cjs",
	}
	for _, s := range sentinels {
		if fileInDir(dir, s) {
			return true
		}
	}
	return false
}

// Inspect returns a structured snapshot of the JS/TS core.
func (j *JSAdapter) Inspect(dir string) core.EcosystemInfo {
	info := core.EcosystemInfo{
		Ecosystem: "javascript",
		Language:  "javascript",
		Detected:  []string{},
	}

	// Detect TypeScript
	hasTS := fileInDir(dir, "tsconfig.json") || fileInDir(dir, "jsconfig.json")
	if hasTS {
		info.Ecosystem = "javascript/typescript"
		info.Language = "typescript"
		info.Detected = append(info.Detected, "typescript")
	}

	// Detect package manager
	pm := j.detectPackageManager(dir)
	info.PackageManager = pm
	info.Detected = append(info.Detected, pm)

	// Detect framework
	fw := j.detectFramework(dir)
	if fw != "" {
		info.Framework = fw
		info.Detected = append(info.Detected, fw)
	}

	// Detect configuration files
	info.ConfigFiles = j.detectConfigFiles(dir)

	// Dependency summary
	if fileInDir(dir, "package.json") {
		info.Dependencies.Manifest = "package.json"
		pkg := j.readPackageJSON(dir)
		if pkg != nil {
			info.Dependencies.Total = j.countDependencies(pkg)
		}
	}
	info.Dependencies.Lockfile = j.detectLockfile(dir)

	// Detect toolchains
	info.Toolchains = j.detectToolchains(dir, info)

	return info
}

// ValidateEnvironment checks that required JS/TS toolchains are available.
func (j *JSAdapter) ValidateEnvironment(dir string) []core.Diagnostic {
	var diags []core.Diagnostic

	// Check Node.js
	nodeStatus := j.checkToolchain("node", "Node.js", ">= 20")
	if nodeStatus != nil {
		diags = append(diags, *nodeStatus)
	}

	// Check package manager
	pm := j.detectPackageManager(dir)
	switch pm {
	case "pnpm":
		if d := j.checkToolchain("pnpm", "pnpm", ">= 9"); d != nil {
			diags = append(diags, *d)
		}
	case "yarn":
		if d := j.checkToolchain("yarn", "Yarn", ">= 4"); d != nil {
			diags = append(diags, *d)
		}
	case "bun":
		if d := j.checkToolchain("bun", "Bun", ""); d != nil {
			diags = append(diags, *d)
		}
	}

	// Check TypeScript if tsconfig exists
	if fileInDir(dir, "tsconfig.json") {
		if d := j.checkToolchain("tsc", "TypeScript compiler", ">= 5"); d != nil {
			diags = append(diags, *d)
		}
	}

	return diags
}

// ValidateDependencies performs deep dependency validation.
func (j *JSAdapter) ValidateDependencies(dir string) []core.Diagnostic {
	var diags []core.Diagnostic

	// Check: package.json must exist for JS/TS projects
	if !fileInDir(dir, "package.json") {
		diags = append(diags, core.Diagnostic{
			ID:         "PMX-PKG-002",
			Severity:   "fail",
			Category:   "dependency",
			Title:      "JavaScript package manifest is missing",
			File:       "package.json",
			Message:    "This project looks like a JavaScript or TypeScript app, but no package manifest was found.",
			Evidence:   []string{"package.json missing", "JS/TS project detected"},
			Suggestion: "Create a package.json and choose a single package manager for this project.",
		})
		return diags
	}

	// Check: multiple package managers
	if j.multiplePackageManagers(dir) {
		lockfiles := j.presentLockfiles(dir)
		diags = append(diags, core.Diagnostic{
			ID:         "PMX-PKG-001",
			Severity:   "warn",
			Category:   "dependency",
			Title:      "Multiple package-manager lockfiles detected",
			File:       "package.json",
			Message:    "More than one package manager lockfile exists in the project root.",
			Evidence:   lockfiles,
			Suggestion: "Keep one canonical package manager and remove stale lockfiles.",
		})
	}

	// Check: lockfile consistency
	pm := j.detectPackageManager(dir)
	lockfile := j.detectLockfile(dir)
	if lockfile == "" && pm != "npm" {
		// pnpm, yarn, bun all create lockfiles — missing one is a warning
		diags = append(diags, core.Diagnostic{
			ID:         "PMX-DEP-001",
			Severity:   "warn",
			Category:   "dependency",
			Title:      "Package manager lockfile is missing",
			File:       "package.json",
			Message:    "The detected package manager (" + pm + ") typically uses a lockfile, but none was found.",
			Evidence:   []string{"package manager: " + pm, "no lockfile detected"},
			Suggestion: "Run the appropriate install command to generate the lockfile before committing.",
		})
	}

	// Deep package.json analysis
	pkg := j.readPackageJSON(dir)
	if pkg == nil {
		return diags
	}

	// Check: dependencies vs devDependencies overlap
	diags = append(diags, j.checkDependencyOverlap(pkg)...)

	// Check: ESLint version vs config model compatibility
	diags = append(diags, j.checkESLintCompatibility(dir, pkg)...)

	return diags
}

// ValidateConfiguration validates JS/TS configuration files and cross-file consistency.
func (j *JSAdapter) ValidateConfiguration(dir string) []core.Diagnostic {
	var diags []core.Diagnostic

	// TypeScript configuration
	if fileInDir(dir, "tsconfig.json") {
		diags = append(diags, j.validateTSConfig(dir)...)
	}

	// ESLint configuration
	diags = append(diags, j.validateESLintConfig(dir)...)

	// Framework + TypeScript consistency
	if j.detectFramework(dir) != "" && fileInDir(dir, "tsconfig.json") {
		diags = append(diags, j.validateFrameworkTSConsistency(dir)...)
	}

	return diags
}

// ValidateToolchain checks toolchain version compatibility with the project.
func (j *JSAdapter) ValidateToolchain(dir string) []core.Diagnostic {
	var diags []core.Diagnostic

	// Check engines field in package.json
	pkg := j.readPackageJSON(dir)
	if pkg == nil {
		return diags
	}

	engines, ok := pkg["engines"].(map[string]interface{})
	if !ok {
		return diags
	}

	for tool, constraint := range engines {
		constraintStr, _ := constraint.(string)
		if constraintStr != "" {
			// We document the constraint even if we can't verify it without exec
			diags = append(diags, core.Diagnostic{
				ID:       "PMX-TOOL-001",
				Severity: "pass",
				Category: "toolchain",
				Title:    "Engine constraint declared",
				File:     "package.json",
				Message:  tool + " constraint: " + constraintStr,
				Evidence: []string{"engines." + tool + " = " + constraintStr},
			})
		}
	}

	return diags
}

// --- Internal helpers ---

func (j *JSAdapter) detectPackageManager(dir string) string {
	// Check packageManager field first (corepack)
	if fileInDir(dir, "package.json") {
		pkg := j.readPackageJSON(dir)
		if pkg != nil {
			if pm, ok := pkg["packageManager"].(string); ok && pm != "" {
				if idx := strings.Index(pm, "@"); idx > 0 {
					return pm[:idx]
				}
				return pm
			}
		}
	}

	switch {
	case fileInDir(dir, "pnpm-lock.yaml"):
		return "pnpm"
	case fileInDir(dir, "yarn.lock"):
		return "yarn"
	case fileInDir(dir, "package-lock.json"):
		return "npm"
	case fileInDir(dir, "bun.lock") || fileInDir(dir, "bun.lockb"):
		return "bun"
	case fileInDir(dir, "package.json"):
		return "npm"
	default:
		return "unknown"
	}
}

func (j *JSAdapter) detectFramework(dir string) string {
	switch {
	case fileInDir(dir, "next.config.js") || fileInDir(dir, "next.config.mjs") || fileInDir(dir, "next.config.ts"):
		return "next"
	case fileInDir(dir, "vite.config.js") || fileInDir(dir, "vite.config.ts") || fileInDir(dir, "vite.config.mjs"):
		return "vite"
	case fileInDir(dir, "astro.config.mjs") || fileInDir(dir, "astro.config.js") || fileInDir(dir, "astro.config.ts"):
		return "astro"
	case fileInDir(dir, "nuxt.config.ts") || fileInDir(dir, "nuxt.config.js"):
		return "nuxt"
	case fileInDir(dir, "remix.config.js"):
		return "remix"
	}

	// Check package.json dependencies for framework hints
	pkg := j.readPackageJSON(dir)
	if pkg == nil {
		return ""
	}
	allDeps := j.allDependencyNames(pkg)
	for _, name := range allDeps {
		switch name {
		case "next":
			return "next"
		case "vite":
			return "vite"
		case "astro":
			return "astro"
		case "nuxt":
			return "nuxt"
		case "@remix-run/react":
			return "remix"
		case "react":
			// Only report react if no other framework found
		}
	}
	return ""
}

func (j *JSAdapter) detectConfigFiles(dir string) []core.ConfigFile {
	var files []core.ConfigFile

	configs := []struct {
		path     string
		model    string
		language string
	}{
		{"package.json", "npm", "json"},
		{"tsconfig.json", "tsc", "json"},
		{"jsconfig.json", "tsc", "json"},
		{"eslint.config.js", "flat", "javascript"},
		{"eslint.config.mjs", "flat", "javascript"},
		{"eslint.config.ts", "flat", "typescript"},
		{".eslintrc.json", "legacy", "json"},
		{".eslintrc.js", "legacy", "javascript"},
		{".eslintrc.cjs", "legacy", "javascript"},
		{".eslintrc", "legacy", ""},
		{"next.config.js", "next", "javascript"},
		{"next.config.mjs", "next", "javascript"},
		{"next.config.ts", "next", "typescript"},
		{"vite.config.js", "vite", "javascript"},
		{"vite.config.ts", "vite", "typescript"},
	}

	for _, c := range configs {
		files = append(files, core.ConfigFile{
			Path:     c.path,
			Exists:   fileInDir(dir, c.path),
			Model:    c.model,
			Language: c.language,
		})
	}

	return files
}

func (j *JSAdapter) detectLockfile(dir string) string {
	switch {
	case fileInDir(dir, "pnpm-lock.yaml"):
		return "pnpm-lock.yaml"
	case fileInDir(dir, "yarn.lock"):
		return "yarn.lock"
	case fileInDir(dir, "package-lock.json"):
		return "package-lock.json"
	case fileInDir(dir, "bun.lock"):
		return "bun.lock"
	case fileInDir(dir, "bun.lockb"):
		return "bun.lockb"
	default:
		return ""
	}
}

func (j *JSAdapter) detectToolchains(dir string, info core.EcosystemInfo) []core.ToolchainInfo {
	var tools []core.ToolchainInfo

	// We declare what we expect; actual version detection happens in ValidateEnvironment
	tools = append(tools, core.ToolchainInfo{Name: "node", Status: "expected"})
	if info.PackageManager == "pnpm" {
		tools = append(tools, core.ToolchainInfo{Name: "pnpm", Status: "expected"})
	}
	if info.Ecosystem == "javascript/typescript" {
		tools = append(tools, core.ToolchainInfo{Name: "typescript", Status: "expected"})
	}

	return tools
}

func (j *JSAdapter) multiplePackageManagers(dir string) bool {
	lockfiles := []string{"pnpm-lock.yaml", "yarn.lock", "package-lock.json", "bun.lock", "bun.lockb"}
	count := 0
	for _, lf := range lockfiles {
		if fileInDir(dir, lf) {
			count++
		}
	}
	return count > 1
}

func (j *JSAdapter) presentLockfiles(dir string) []string {
	var present []string
	lockfiles := []string{"pnpm-lock.yaml", "yarn.lock", "package-lock.json", "bun.lock", "bun.lockb"}
	for _, lf := range lockfiles {
		if fileInDir(dir, lf) {
			present = append(present, lf)
		}
	}
	return present
}

func (j *JSAdapter) readPackageJSON(dir string) map[string]interface{} {
	data, err := os.ReadFile(filepath.Join(dir, "package.json"))
	if err != nil {
		return nil
	}
	var pkg map[string]interface{}
	if err := json.Unmarshal(data, &pkg); err != nil {
		return nil
	}
	return pkg
}

func (j *JSAdapter) countDependencies(pkg map[string]interface{}) int {
	count := 0
	for _, key := range []string{"dependencies", "devDependencies", "peerDependencies", "optionalDependencies"} {
		if deps, ok := pkg[key].(map[string]interface{}); ok {
			count += len(deps)
		}
	}
	return count
}

func (j *JSAdapter) allDependencyNames(pkg map[string]interface{}) []string {
	var names []string
	seen := map[string]bool{}
	for _, key := range []string{"dependencies", "devDependencies"} {
		if deps, ok := pkg[key].(map[string]interface{}); ok {
			for name := range deps {
				if !seen[name] {
					names = append(names, name)
					seen[name] = true
				}
			}
		}
	}
	return names
}

func (j *JSAdapter) checkDependencyOverlap(pkg map[string]interface{}) []core.Diagnostic {
	var diags []core.Diagnostic

	deps, _ := pkg["dependencies"].(map[string]interface{})
	devDeps, _ := pkg["devDependencies"].(map[string]interface{})

	if deps == nil || devDeps == nil {
		return diags
	}

	var overlap []string
	for name := range deps {
		if _, inDev := devDeps[name]; inDev {
			overlap = append(overlap, name)
		}
	}

	if len(overlap) > 0 {
		diags = append(diags, core.Diagnostic{
			ID:         "PMX-DEP-002",
			Severity:   "warn",
			Category:   "dependency",
			Title:      "Dependency appears in both dependencies and devDependencies",
			File:       "package.json",
			Message:    "Some packages are listed in both dependencies and devDependencies, which can cause unexpected behavior.",
			Evidence:   overlap,
			Suggestion: "Move runtime packages to dependencies only and development-only packages to devDependencies.",
		})
	}

	return diags
}

func (j *JSAdapter) checkESLintCompatibility(dir string, pkg map[string]interface{}) []core.Diagnostic {
	var diags []core.Diagnostic

	hasLegacy := fileInDir(dir, ".eslintrc") || fileInDir(dir, ".eslintrc.json") || fileInDir(dir, ".eslintrc.js") || fileInDir(dir, ".eslintrc.cjs")
	hasFlat := fileInDir(dir, "eslint.config.js") || fileInDir(dir, "eslint.config.mjs") || fileInDir(dir, "eslint.config.ts")

	if !hasLegacy && !hasFlat {
		return diags // No ESLint config detected
	}

	// Check ESLint version from dependencies
	var eslintVersion string
	for _, key := range []string{"dependencies", "devDependencies"} {
		if deps, ok := pkg[key].(map[string]interface{}); ok {
			if ver, ok := deps["eslint"].(string); ok {
				eslintVersion = ver
				break
			}
		}
	}

	// Detect ESLint 9+ with legacy config
	if hasLegacy && eslintVersion != "" {
		major := j.extractMajorVersion(eslintVersion)
		if major >= 9 {
			legacyFile := j.detectLegacyESLintFile(dir)
			diags = append(diags, core.Diagnostic{
				ID:         "PMX-DEP-004",
				Severity:   "warn",
				Category:   "dependency",
				Title:      "ESLint 9+ with legacy configuration model",
				File:       legacyFile,
				Message:    "ESLint 9 uses flat config by default. Legacy configuration model may cause compatibility issues.",
				Evidence:   []string{"eslint: " + eslintVersion, "config: " + legacyFile, "ESLint 9 flat-config mode is default"},
				Suggestion: "Migrate to eslint.config.js (flat config) or verify legacy compatibility.",
			})
		}
	}

	// Both legacy and flat configs present
	if hasLegacy && hasFlat {
		diags = append(diags, core.Diagnostic{
			ID:         "PMX-DEP-005",
			Severity:   "warn",
			Category:   "configuration",
			Title:      "Both legacy and flat ESLint configs detected",
			File:       "eslint.config.js",
			Message:    "Having both legacy (.eslintrc*) and flat (eslint.config.*) ESLint configurations can cause confusion.",
			Evidence:   []string{j.detectLegacyESLintFile(dir), "eslint.config.*"},
			Suggestion: "Remove the legacy configuration and use flat config only.",
		})
	}

	return diags
}

func (j *JSAdapter) validateTSConfig(dir string) []core.Diagnostic {
	var diags []core.Diagnostic

	data, err := os.ReadFile(filepath.Join(dir, "tsconfig.json"))
	if err != nil {
		return diags
	}

	text := string(data)
	strictEnabled := strings.Contains(text, `"strict": true`) ||
		strings.Contains(text, `"strict":\n    true`) ||
		strings.Contains(text, `"strict":\n\ttrue`)

	if !strictEnabled {
		diags = append(diags, core.Diagnostic{
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

	// Check for target/module consistency with framework
	framework := j.detectFramework(dir)
	if framework == "next" {
		// Next.js requires ES modules
		if strings.Contains(text, `"target": "ES5"`) || strings.Contains(text, `"target": "ES3"`) {
			diags = append(diags, core.Diagnostic{
				ID:         "PMX-CONFIG-002",
				Severity:   "warn",
				Category:   "configuration",
				Title:      "TypeScript target may be incompatible with framework",
				File:       "tsconfig.json",
				Message:    "Next.js typically requires ES2017+ target. A lower target may cause runtime issues.",
				Evidence:   []string{"target: ES5 or ES3", "framework: Next.js"},
				Suggestion: "Set compilerOptions.target to ES2017 or higher in tsconfig.json.",
			})
		}
	}

	return diags
}

func (j *JSAdapter) validateESLintConfig(dir string) []core.Diagnostic {
	var diags []core.Diagnostic

	hasLegacy := fileInDir(dir, ".eslintrc") || fileInDir(dir, ".eslintrc.json") || fileInDir(dir, ".eslintrc.js") || fileInDir(dir, ".eslintrc.cjs")

	if hasLegacy {
		legacyFile := j.detectLegacyESLintFile(dir)
		diags = append(diags, core.Diagnostic{
			ID:         "PMX-ESLINT-001",
			Severity:   "warn",
			Category:   "eslint",
			Title:      "Legacy ESLint configuration detected",
			File:       legacyFile,
			Message:    "Use of a legacy ESLint config style may require migration to flat config.",
			Evidence:   []string{legacyFile + " detected", "ESLint 9 compatibility risk"},
			Suggestion: "Verify ESLint flat-config compatibility before upgrading or running lint in CI.",
		})
	}

	return diags
}

func (j *JSAdapter) validateFrameworkTSConsistency(dir string) []core.Diagnostic {
	// Framework + TypeScript cross-validation
	// This is where we detect configuration mismatches between framework
	// requirements and TypeScript settings.
	return nil // Will be expanded as more framework-specific checks are added
}

func (j *JSAdapter) detectLegacyESLintFile(dir string) string {
	switch {
	case fileInDir(dir, ".eslintrc.json"):
		return ".eslintrc.json"
	case fileInDir(dir, ".eslintrc.js"):
		return ".eslintrc.js"
	case fileInDir(dir, ".eslintrc.cjs"):
		return ".eslintrc.cjs"
	case fileInDir(dir, ".eslintrc"):
		return ".eslintrc"
	default:
		return "eslint.config.js"
	}
}

func (j *JSAdapter) checkToolchain(name, displayName, minVersion string) *core.Diagnostic {
	// Check if toolchain is available by looking for it in PATH
	path, err := findExecutable(name)
	if err != nil {
		return &core.Diagnostic{
			ID:         "PMX-ENV-003",
			Severity:   "warn",
			Category:   "environment",
			Title:      displayName + " is not installed",
			Message:    displayName + " is required but was not found in PATH. The project configuration is still validated.",
			Evidence:   []string{name + " not found in PATH", "required: " + displayName + " " + minVersion},
			Suggestion: "Install " + displayName + " before running project validation.",
		}
	}

	// Try to get version
	ver := j.getToolVersion(name)
	if ver != "" {
		return &core.Diagnostic{
			ID:       "PMX-ENV-001",
			Severity: "pass",
			Category: "environment",
			Title:    displayName + " is available",
			Message:  displayName + " " + ver + " found at " + path,
			Evidence: []string{name + " version: " + ver, "path: " + path},
		}
	}

	return &core.Diagnostic{
		ID:       "PMX-ENV-002",
		Severity: "pass",
		Category: "environment",
		Title:    displayName + " is available",
		Message:  displayName + " found at " + path,
		Evidence: []string{"path: " + path},
	}
}

func (j *JSAdapter) extractMajorVersion(version string) int {
	// Strip leading ^, ~, >=, v
	v := strings.TrimLeft(version, "^~>=v")
	parts := strings.SplitN(v, ".", 2)
	if len(parts) > 0 {
		major, err := strconv.Atoi(parts[0])
		if err == nil {
			return major
		}
	}
	return 0
}

func (j *JSAdapter) getToolVersion(name string) string {
	return getToolVersionFunc(name)
}

func fileInDir(dir, name string) bool {
	_, err := os.Stat(filepath.Join(dir, name))
	return err == nil
}

func findExecutable(name string) (string, error) {
	path, err := execLookPath(name)
	if err != nil {
		return "", err
	}
	return path, nil
}
