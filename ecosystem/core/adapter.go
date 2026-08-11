package core

import "path/filepath"

// Diagnostic is the unified diagnostic model across all adapters.
type Diagnostic struct {
	ID         string   `json:"id"`
	Severity   string   `json:"severity"`
	Category   string   `json:"category,omitempty"`
	Title      string   `json:"title,omitempty"`
	File       string   `json:"file,omitempty"`
	Message    string   `json:"message"`
	Evidence   []string `json:"evidence,omitempty"`
	Suggestion string   `json:"suggestion,omitempty"`
}

// EvidenceRecord captures a single atomic validation fact observed by an adapter.
type EvidenceRecord struct {
	Source    string   `json:"source,omitempty"`
	Evidence  []string `json:"evidence,omitempty"`
	Timestamp string   `json:"timestamp,omitempty"`
}

// CheckResult models the result of an adapter-level health check.
type CheckResult struct {
	Name   string `json:"name"`
	Status string `json:"status"`
	Detail string `json:"detail,omitempty"`
}

// ToolchainInfo describes a detected toolchain and its version.
type ToolchainInfo struct {
	Name    string `json:"name"`
	Version string `json:"version,omitempty"`
	Path    string `json:"path,omitempty"`
	Status  string `json:"status"`
}

// DependencySummary summarizes the dependency state.
type DependencySummary struct {
	Manifest string   `json:"manifest,omitempty"`
	Lockfile string   `json:"lockfile,omitempty"`
	Total    int      `json:"total,omitempty"`
	Outdated int      `json:"outdated,omitempty"`
	Missing  bool     `json:"missing,omitempty"`
	Issues   []string `json:"issues,omitempty"`
}

// ConfigFile describes a detected configuration file.
type ConfigFile struct {
	Path     string `json:"path"`
	Exists   bool   `json:"exists"`
	Model    string `json:"model,omitempty"`
	Language string `json:"language,omitempty"`
}

// EcosystemInfo is the structured result of an adapter's Inspect() call.
type EcosystemInfo struct {
	Ecosystem      string            `json:"ecosystem"`
	PackageManager string            `json:"package_manager"`
	Language       string            `json:"language"`
	Framework      string            `json:"framework,omitempty"`
	Toolchains     []ToolchainInfo   `json:"toolchains,omitempty"`
	Dependencies   DependencySummary `json:"dependencies,omitempty"`
	ConfigFiles    []ConfigFile      `json:"config_files,omitempty"`
	Detected       []string          `json:"detected,omitempty"`
	Name           string            `json:"name,omitempty"`
	Root           string            `json:"root,omitempty"`
	Manifest       string            `json:"manifest,omitempty"`
	FileCount      int               `json:"file_count,omitempty"`
	Details        map[string]string `json:"details,omitempty"`
}

// EcosystemAdapter is the key abstraction for cross-language validation.
type EcosystemAdapter interface {
	Name() string
	Detect(dir string) bool
	Inspect(dir string) EcosystemInfo
	ValidateEnvironment(dir string) []Diagnostic
	ValidateDependencies(dir string) []Diagnostic
	ValidateConfiguration(dir string) []Diagnostic
	ValidateToolchain(dir string) []Diagnostic
	ValidateProject(dir string) []Diagnostic
}

func NormalizeDir(dir string) string {
	if dir == "" {
		return "."
	}
	return filepath.Clean(dir)
}
