package python

import (
	"os"
	"path/filepath"

	"github.com/Fenz-art/picomatchgo-pmx-CLI/ecosystem/core"
)

// PythonAdapter is a minimal placeholder adapter used to keep the adapter
// registry complete in the current source tree.
type PythonAdapter struct{}

func (p *PythonAdapter) Name() string { return "python" }

func (p *PythonAdapter) Detect(dir string) bool {
	return fileInDir(dir, "pyproject.toml") || fileInDir(dir, "requirements.txt")
}

func (p *PythonAdapter) Inspect(dir string) core.EcosystemInfo {
	return core.EcosystemInfo{
		Ecosystem:      "python",
		PackageManager: "pip",
		Language:       "python",
		Detected:       []string{"pyproject.toml", "requirements.txt"},
	}
}

func (p *PythonAdapter) ValidateEnvironment(dir string) []core.Diagnostic   { return nil }
func (p *PythonAdapter) ValidateDependencies(dir string) []core.Diagnostic  { return nil }
func (p *PythonAdapter) ValidateConfiguration(dir string) []core.Diagnostic { return nil }
func (p *PythonAdapter) ValidateToolchain(dir string) []core.Diagnostic     { return nil }
func (p *PythonAdapter) ValidateProject(dir string) []core.Diagnostic       { return nil }

func fileInDir(dir, name string) bool {
	_, err := os.Stat(filepath.Join(dir, name))
	return err == nil
}
