package core

// EcosystemAdapter is the key abstraction for cross-language project validation.
// Each adapter provides ecosystem-specific detection, inspection, and validation.
// Adding a new language means implementing this interface, not redesigning PMX.
type EcosystemAdapter interface {
	// Name returns the adapter identifier (e.g. "javascript", "rust", "python", "go").
	Name() string

	// Detect returns true if this adapter's ecosystem is present in the project directory.
	Detect(dir string) bool

	// Inspect returns a structured snapshot of the detected ecosystem.
	Inspect(dir string) EcosystemInfo

	// ValidateEnvironment checks runtime toolchain availability and version constraints.
	ValidateEnvironment(dir string) []Diagnostic

	// ValidateDependencies checks dependency manifests, lockfiles, and compatibility.
	ValidateDependencies(dir string) []Diagnostic

	// ValidateConfiguration checks configuration files and cross-file consistency.
	ValidateConfiguration(dir string) []Diagnostic

	// ValidateToolchain checks toolchain version compatibility with the project.
	ValidateToolchain(dir string) []Diagnostic

	// ValidateProject performs a project-level validation that is rooted at the
	// directory being inspected rather than a single file or config layer.
	ValidateProject(dir string) []Diagnostic
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
}

// ToolchainInfo describes a detected toolchain and its version.
type ToolchainInfo struct {
	Name    string `json:"name"`
	Version string `json:"version,omitempty"`
	Path    string `json:"path,omitempty"`
	Status  string `json:"status"` // "available", "missing", "incompatible"
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

// Diagnostic is the unified diagnostic model across all adapters.
type Diagnostic struct {
	ID         string   `json:"id"`
	Severity   string   `json:"severity"` // "pass", "warn", "fail"
	Category   string   `json:"category,omitempty"`
	Title      string   `json:"title,omitempty"`
	File       string   `json:"file,omitempty"`
	Message    string   `json:"message"`
	Evidence   []string `json:"evidence,omitempty"`
	Suggestion string   `json:"suggestion,omitempty"`
}

// ProjectDetection is the result of scanning a directory for all applicable ecosystems.
type ProjectDetection struct {
	Directory  string         `json:"directory"`
	Ecosystems []EcosystemRef `json:"ecosystems"`
	Primary    string         `json:"primary"`
}

// EcosystemRef references a detected ecosystem and its adapter.
type EcosystemRef struct {
	Name    string `json:"name"`
	Adapter string `json:"adapter"`
	Primary bool   `json:"primary,omitempty"`
}
