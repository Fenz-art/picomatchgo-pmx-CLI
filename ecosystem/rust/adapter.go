package rust

import (
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/Fenz-art/picomatchgo-pmx-CLI/ecosystem/core"
)

// Adapter implements validation for Rust projects.
type Adapter struct{}

func NewRustAdapter() *Adapter { return &Adapter{} }

func (a *Adapter) Name() string { return "rust" }

func (a *Adapter) Detect(dir string) bool {
	_, err := os.Stat(filepath.Join(dir, "Cargo.toml"))
	return err == nil
}

func (a *Adapter) Inspect(dir string) core.EcosystemInfo {
	manifest := filepath.Join(dir, "Cargo.toml")
	info := core.EcosystemInfo{
		Ecosystem:      "rust",
		PackageManager: "cargo",
		Language:       "rust",
		Name:           "rust",
		Root:           core.NormalizeDir(dir),
		Manifest:       manifest,
		Details:        map[string]string{"manifest": manifest},
	}
	if data, err := os.ReadFile(manifest); err == nil {
		if name := firstMatch(string(data), `(?m)^name\s*=\s*"([^"]+)"`); name != "" {
			info.Details["package_name"] = name
		}
		if edition := firstMatch(string(data), `(?m)^edition\s*=\s*"([^"]+)"`); edition != "" {
			info.Details["edition"] = edition
		}
	}
	if files, err := os.ReadDir(dir); err == nil {
		info.FileCount = len(files)
	}
	return info
}

func (a *Adapter) ValidateEnvironment(dir string) []core.Diagnostic {
	if _, err := exec.LookPath("cargo"); err != nil {
		return []core.Diagnostic{{
			ID:         "PMX-ENV-001",
			Severity:   "fail",
			Category:   "environment",
			Title:      "Rust toolchain is missing",
			File:       "Cargo.toml",
			Message:    "A Rust project was detected, but the cargo toolchain is not installed or not on PATH.",
			Suggestion: "Install Rust and ensure cargo and rustc are available on PATH.",
		}}
	}
	if _, err := exec.LookPath("rustc"); err != nil {
		return []core.Diagnostic{{
			ID:         "PMX-ENV-003",
			Severity:   "fail",
			Category:   "environment",
			Title:      "Rust compiler is missing",
			File:       "Cargo.toml",
			Message:    "The Rust compiler is not available, so the project cannot be built or validated.",
			Suggestion: "Install rustc before running cargo test or clippy.",
		}}
	}
	return nil
}

func (a *Adapter) ValidateDependencies(dir string) []core.Diagnostic {
	if _, err := os.Stat(filepath.Join(dir, "Cargo.lock")); err != nil {
		return []core.Diagnostic{{
			ID:         "PMX-DEP-001",
			Severity:   "warn",
			Category:   "dependencies",
			Title:      "Cargo.lock is missing",
			File:       "Cargo.toml",
			Message:    "This Rust project does not currently have a Cargo.lock file in the repo root.",
			Suggestion: "Run `cargo generate-lockfile` to record the exact dependency graph for reproducible builds.",
		}}
	}
	return nil
}

func (a *Adapter) ValidateConfiguration(dir string) []core.Diagnostic {
	manifestPath := filepath.Join(dir, "Cargo.toml")
	data, err := os.ReadFile(manifestPath)
	if err != nil {
		return nil
	}
	content := string(data)
	if edition := firstMatch(content, `(?m)^edition\s*=\s*"([^"]+)"`); edition != "" && edition != "2021" && edition != "2024" {
		return []core.Diagnostic{{
			ID:         "PMX-CONFIG-002",
			Severity:   "warn",
			Category:   "configuration",
			Title:      "Unsupported Rust edition detected",
			File:       "Cargo.toml",
			Message:    "The manifest uses a Rust edition that is outside the currently supported modern toolchain set.",
			Suggestion: "Use edition = \"2021\" or edition = \"2024\" to stay aligned with the current toolchain.",
		}}
	}
	if !strings.Contains(content, "[package]") {
		return []core.Diagnostic{{
			ID:         "PMX-CONFIG-001",
			Severity:   "warn",
			Category:   "configuration",
			Title:      "Rust package manifest is incomplete",
			File:       "Cargo.toml",
			Message:    "The manifest is missing the standard [package] section.",
			Suggestion: "Add the package metadata and edition declaration expected by Cargo.",
		}}
	}
	return nil
}

func (a *Adapter) ValidateToolchain(dir string) []core.Diagnostic {
	if _, err := exec.LookPath("cargo"); err != nil {
		return nil
	}
	if _, err := exec.LookPath("rustc"); err != nil {
		return nil
	}
	out, err := exec.Command("cargo", "test", "--quiet", "--no-run").CombinedOutput()
	if err != nil {
		return []core.Diagnostic{{
			ID:         "PMX-TOOL-001",
			Severity:   "warn",
			Category:   "toolchain",
			Title:      "Cargo test could not complete",
			File:       "Cargo.toml",
			Message:    "The project manifest or toolchain configuration is preventing cargo test from finishing.",
			Suggestion: "Run `cargo test -- --nocapture` to inspect the specific build or test failure.",
			Evidence:   []string{strings.TrimSpace(string(out))},
		}}
	}
	if _, err := exec.LookPath("cargo-clippy"); err != nil {
		return []core.Diagnostic{{
			ID:         "PMX-TOOL-002",
			Severity:   "warn",
			Category:   "toolchain",
			Title:      "Clippy is not installed",
			File:       "Cargo.toml",
			Message:    "Cargo Clippy is not installed, so lint checks cannot run in this environment.",
			Suggestion: "Install clippy via `rustup component add clippy`.",
		}}
	}
	return nil
}

func (a *Adapter) ValidateProject(dir string) []core.Diagnostic {
	if _, err := os.Stat(filepath.Join(dir, "Cargo.toml")); err != nil {
		return nil
	}
	return []core.Diagnostic{{
		ID:       "PMX-PROJECT-003",
		Severity: "pass",
		Category: "project",
		Title:    "Rust project validated",
		File:     "Cargo.toml",
		Message:  "Rust project evidence has been collected through the adapter contract.",
		Evidence: []string{"Cargo.toml detected"},
	}}
}

func firstMatch(content, pattern string) string {
	matches := regexp.MustCompile(pattern).FindStringSubmatch(content)
	if len(matches) < 2 {
		return ""
	}
	return matches[1]
}
