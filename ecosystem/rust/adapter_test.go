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
	if info.Name != "rust" {
		t.Fatalf("expected name rust, got %q", info.Name)
	}
	if info.Root == "" {
		t.Fatal("expected non-empty root")
	}
	if info.FileCount == 0 {
		t.Fatal("expected file count > 0")
	}
	if info.Details["manifest"] == "" {
		t.Fatal("expected manifest detail")
	}

	var _ core.EcosystemAdapter = adapter
}

func TestRustAdapterValidateProjectContract(t *testing.T) {
	adapter := NewRustAdapter()
	_ = adapter.ValidateProject(t.TempDir())
}
