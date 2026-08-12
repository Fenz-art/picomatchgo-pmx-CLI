package rust

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Fenz-art/picomatchgo-pmx-CLI/ecosystem/core"
)

func TestRustAdapterDetectAndInspect(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "Cargo.toml"), []byte("[package]\nname = \"demo\"\nversion = \"0.1.0\"\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	adapter := NewRustAdapter()
	if !adapter.Detect(dir) {
		t.Fatalf("expected rust adapter to detect Cargo.toml")
	}

	info := adapter.Inspect(dir)
	if info.Ecosystem != "rust" {
		t.Fatalf("expected ecosystem rust, got %q", info.Ecosystem)
	}
	if info.Dependencies.Manifest != "Cargo.toml" {
		t.Fatalf("expected Dependencies.Manifest to be Cargo.toml, got %q", info.Dependencies.Manifest)
	}
	if len(info.Detected) == 0 {
		t.Fatalf("expected Detected to be non-empty")
	}

	var _ core.EcosystemAdapter = adapter
}

func TestRustAdapterValidateProjectContract(t *testing.T) {
	adapter := NewRustAdapter()
	_ = adapter.ValidateProject(t.TempDir())
}
