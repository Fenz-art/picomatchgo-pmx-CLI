package ecosystem

import (
	"os"
	"path/filepath"

	"github.com/Fenz-art/picomatchgo-pmx-CLI/ecosystem/core"
	"github.com/Fenz-art/picomatchgo-pmx-CLI/ecosystem/javascript"
)

// AllAdapters returns every registered adapter. Add new adapters here.
func AllAdapters() []core.EcosystemAdapter {
	return []core.EcosystemAdapter{
		&javascript.JSAdapter{},
		&GoAdapter{},
	}
}

// DetectEcosystems scans a project directory and returns all applicable ecosystems.
func DetectEcosystems(dir string, adapters []core.EcosystemAdapter) core.ProjectDetection {
	result := core.ProjectDetection{Directory: dir}

	for _, adapter := range adapters {
		if adapter.Detect(dir) {
			ref := core.EcosystemRef{Name: adapter.Name(), Adapter: adapter.Name()}
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

// GoAdapter is the inline runtime adapter integration for the Go module model.
type GoAdapter struct{}

func (g *GoAdapter) Name() string { return "go" }

func (g *GoAdapter) Detect(dir string) bool {
	return fileInDir(dir, "go.mod")
}

func (g *GoAdapter) Inspect(dir string) core.EcosystemInfo {
	info := core.EcosystemInfo{
		Ecosystem:      "go",
		PackageManager: "go",
		Language:       "go",
		Detected:       []string{"go.mod"},
	}
	if fileInDir(dir, "go.sum") {
		info.Dependencies.Lockfile = "go.sum"
	}
	return info
}

func (g *GoAdapter) ValidateEnvironment(dir string) []core.Diagnostic   { return nil }
func (g *GoAdapter) ValidateDependencies(dir string) []core.Diagnostic  { return nil }
func (g *GoAdapter) ValidateConfiguration(dir string) []core.Diagnostic { return nil }
func (g *GoAdapter) ValidateToolchain(dir string) []core.Diagnostic     { return nil }
func (g *GoAdapter) ValidateProject(dir string) []core.Diagnostic       { return nil }

func fileInDir(dir, name string) bool {
	_, err := os.Stat(filepath.Join(dir, name))
	return err == nil
}
