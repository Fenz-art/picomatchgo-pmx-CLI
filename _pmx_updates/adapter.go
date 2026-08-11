package core

// EcosystemAdapter is the key abstraction for cross-language project validation.
// Each adapter provides ecosystem-specific detection, inspection, and validation.
// Adding a new language (Rust, Python, C/C++) means implementing this interface,
// not redesigning PMX.
type EcosystemAdapter interface {
	// Name returns the adapter identifier (e.g. "javascript", "rust", "python").
	Name() string

	// Detect returns true if this adapter's ecosystem is present in the project directory.
	// Detect should be fast and only check for sentinel files (package.json, Cargo.toml, etc.).
	Detect(dir string) bool

	// Inspect returns a structured snapshot of the detected ecosystem.
	// This is the "pmx inspect" layer — what ecosystem, framework, package manager, etc.
	Inspect(dir string) EcosystemInfo

	// ValidateEnvironment checks runtime toolchain availability and version constraints.
	// This validates that Node/pnpm/TypeScript/etc are installed and meet version requirements.
	ValidateEnvironment(dir string) []Diagnostic

	// ValidateDependencies checks dependency manifests, lockfiles, and compatibility.
	// This goes beyond "does package.json exist" to understand the dependency graph.
	ValidateDependencies(dir string) []Diagnostic

	// ValidateConfiguration checks configuration files and cross-file consistency.
	// This validates tsconfig, eslint config, framework config, and their relationships.
	ValidateConfiguration(dir string) []Diagnostic

	// ValidateToolchain checks toolchain version compatibility with the project.
	ValidateToolchain(dir string) []Diagnostic
}

// EcosystemInfo is the structured result of an adapter's Inspect() call.
// This is what "pmx inspect" or "pmx agent inspect --json" returns per-ecosystem.
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
	Manifest string   `json:"manifest,omitempty"` // "package.json", "Cargo.toml", etc.
	Lockfile string   `json:"lockfile,omitempty"` // "pnpm-lock.yaml", "Cargo.lock", etc.
	Total    int      `json:"total,omitempty"`
	Outdated int      `json:"outdated,omitempty"`
	Missing  bool     `json:"missing,omitempty"`
	Issues   []string `json:"issues,omitempty"`
}

// ConfigFile describes a detected configuration file.
type ConfigFile struct {
	Path     string `json:"path"`
	Exists   bool   `json:"exists"`
	Model    string `json:"model,omitempty"` // "legacy" or "flat" for ESLint, etc.
	Language string `json:"language,omitempty"`
}

// Diagnostic is the unified diagnostic model across all adapters.
// This is what flows through "pmx doctor --json" and "pmx agent check --json".
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
	Primary    string         `json:"primary"` // The dominant ecosystem
}

// EcosystemRef references a detected ecosystem and its adapter.
type EcosystemRef struct {
	Name    string `json:"name"`
	Adapter string `json:"adapter"`
	Primary bool   `json:"primary,omitempty"`
}
