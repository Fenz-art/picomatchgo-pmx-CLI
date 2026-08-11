package ecosystem

import (
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/Fenz-art/picomatchgo-pmx-CLI/ecosystem/core"
	jsadapter "github.com/Fenz-art/picomatchgo-pmx-CLI/ecosystem/javascript"
	pyadapter "github.com/Fenz-art/picomatchgo-pmx-CLI/ecosystem/python"
	rustadapter "github.com/Fenz-art/picomatchgo-pmx-CLI/ecosystem/rust"
)

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

// AllAdapters returns every registered adapter. Add new adapters here.
func AllAdapters() []core.EcosystemAdapter {
	return []core.EcosystemAdapter{
		&GoAdapter{},
		jsadapter.NewJavaScriptAdapter(),
		rustadapter.NewRustAdapter(),
		pyadapter.NewPythonAdapter(),
	}
}

// DetectEcosystems scans a project directory and returns all applicable ecosystems.
func DetectEcosystems(dir string, adapters []core.EcosystemAdapter) ProjectDetection {
	result := ProjectDetection{Directory: dir}
	for _, adapter := range adapters {
		if adapter.Detect(dir) {
			ref := EcosystemRef{Name: adapter.Name(), Adapter: adapter.Name()}
			result.Ecosystems = append(result.Ecosystems, ref)
		}
	}
	if len(result.Ecosystems) > 0 {
		result.Primary = result.Ecosystems[0].Name
		result.Ecosystems[0].Primary = true
	}
	return result
}

// ValidateAll runs all validation layers for all detected adapters.
func ValidateAll(dir string, adapters []core.EcosystemAdapter) []core.Diagnostic {
	var all []core.Diagnostic
	for _, adapter := range adapters {
		if !adapter.Detect(dir) {
			continue
		}
		all = append(all, adapter.ValidateEnvironment(dir)...)
		all = append(all, adapter.ValidateDependencies(dir)...)
		all = append(all, adapter.ValidateConfiguration(dir)...)
		all = append(all, adapter.ValidateToolchain(dir)...)
		all = append(all, adapter.ValidateProject(dir)...)
	}
	return all
}

// InspectAll runs Inspect() for all detected adapters.
func InspectAll(dir string, adapters []core.EcosystemAdapter) []core.EcosystemInfo {
	var infos []core.EcosystemInfo
	for _, adapter := range adapters {
		if adapter.Detect(dir) {
			infos = append(infos, adapter.Inspect(dir))
		}
	}
	return infos
}

// GoAdapter validates Go projects.
type GoAdapter struct{}

func (g *GoAdapter) Name() string { return "go" }

func (g *GoAdapter) Detect(dir string) bool {
	_, err := os.Stat(filepath.Join(dir, "go.mod"))
	return err == nil
}

func (g *GoAdapter) Inspect(dir string) core.EcosystemInfo {
	info := core.EcosystemInfo{
		Ecosystem:      "go",
		PackageManager: "go",
		Language:       "go",
		Root:           core.NormalizeDir(dir),
		Manifest:       filepath.Join(dir, "go.mod"),
		Detected:       []string{"go.mod"},
		Details:        map[string]string{},
		Name:           "go",
	}
	if fileExists(filepath.Join(dir, "go.sum")) {
		info.Dependencies.Lockfile = "go.sum"
	}
	if ver, _ := detectGoVersion(); ver != "" {
		info.Toolchains = append(info.Toolchains, core.ToolchainInfo{Name: "go", Version: ver, Status: "available"})
	}
	return info
}

func (g *GoAdapter) ValidateEnvironment(dir string) []core.Diagnostic {
	if ver, err := detectGoVersion(); err != nil || ver == "" {
		return []core.Diagnostic{{
			ID:         "PMX-GO-001",
			Severity:   "fail",
			Category:   "environment",
			Title:      "Go toolchain not found",
			File:       "go.mod",
			Message:    "Go is not installed or not in PATH. This Go project cannot be validated without the Go toolchain.",
			Evidence:   []string{"go version check failed"},
			Suggestion: "Install Go from https://go.dev/dl before running project validation.",
		}}
	}
	return nil
}

func (g *GoAdapter) ValidateDependencies(dir string) []core.Diagnostic  { return nil }
func (g *GoAdapter) ValidateConfiguration(dir string) []core.Diagnostic { return nil }
func (g *GoAdapter) ValidateToolchain(dir string) []core.Diagnostic     { return nil }
func (g *GoAdapter) ValidateProject(dir string) []core.Diagnostic {
	if _, err := os.Stat(filepath.Join(dir, "go.mod")); err != nil {
		return nil
	}
	return []core.Diagnostic{{
		ID:       "PMX-PROJECT-001",
		Severity: "pass",
		Category: "project",
		Title:    "Go project validated",
		File:     "go.mod",
		Message:  "Go project evidence has been collected through the adapter contract.",
		Evidence: []string{"go.mod detected"},
	}}
}

func detectGoVersion() (string, error) {
	if _, err := exec.LookPath("go"); err != nil {
		return "", err
	}
	out, err := exec.Command("go", "version").CombinedOutput()
	if err != nil {
		return "", err
	}
	text := strings.TrimSpace(string(out))
	if text == "" {
		return "", os.ErrNotExist
	}
	return text, nil
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func firstMatch(content, pattern string) string {
	matches := regexpMustCompile(pattern).FindStringSubmatch(content)
	if len(matches) < 2 {
		return ""
	}
	return matches[1]
}

func regexpMustCompile(pattern string) *regexp.Regexp {
	res, err := regexp.Compile(pattern)
	if err != nil {
		panic(err)
	}
	return res
}
