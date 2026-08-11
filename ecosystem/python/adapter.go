package python

import (
	"os"
	"path/filepath"

	"github.com/Fenz-art/picomatchgo-pmx-CLI/ecosystem/core"
)

// Adapter implements a minimal validation surface for Python projects.
type Adapter struct{}

func NewPythonAdapter() *Adapter { return &Adapter{} }

func (a *Adapter) Name() string { return "python" }

func (a *Adapter) Detect(dir string) bool {
	for _, name := range []string{"pyproject.toml", "requirements.txt", "setup.py", "Pipfile"} {
		if _, err := os.Stat(filepath.Join(dir, name)); err == nil {
			return true
		}
	}
	return false
}

func (a *Adapter) Inspect(dir string) core.EcosystemInfo {
	info := core.EcosystemInfo{
		Ecosystem:      "python",
		PackageManager: "pip",
		Language:       "python",
		Name:           "python",
		Root:           core.NormalizeDir(dir),
		Details:        map[string]string{},
	}
	for _, file := range []string{"pyproject.toml", "requirements.txt", "setup.py", "Pipfile"} {
		if _, err := os.Stat(filepath.Join(dir, file)); err == nil {
			info.Manifest = filepath.Join(dir, file)
			break
		}
	}
	if files, err := os.ReadDir(dir); err == nil {
		info.FileCount = len(files)
	}
	return info
}

func (a *Adapter) ValidateEnvironment(dir string) []core.Diagnostic   { return nil }
func (a *Adapter) ValidateDependencies(dir string) []core.Diagnostic  { return nil }
func (a *Adapter) ValidateConfiguration(dir string) []core.Diagnostic { return nil }
func (a *Adapter) ValidateToolchain(dir string) []core.Diagnostic     { return nil }
func (a *Adapter) ValidateProject(dir string) []core.Diagnostic       { return nil }
