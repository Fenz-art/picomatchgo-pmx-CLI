package javascript

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Fenz-art/picomatchgo-pmx-CLI/ecosystem/core"
)

func TestJSAdapterDetectAndInspect(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "package.json"), []byte(`{"name":"demo"}`), 0o600); err != nil {
		t.Fatal(err)
	}

	adapter := &JSAdapter{}
	if !adapter.Detect(dir) {
		t.Fatalf("expected adapter to detect JavaScript project")
	}

	info := adapter.Inspect(dir)
	if info.Ecosystem != "javascript" {
		t.Fatalf("expected ecosystem javascript, got %q", info.Ecosystem)
	}
	if info.PackageManager != "npm" {
		t.Fatalf("expected default package manager npm, got %q", info.PackageManager)
	}
	if info.Dependencies.Manifest != "package.json" {
		t.Fatalf("expected package.json manifest, got %q", info.Dependencies.Manifest)
	}

	var _ core.EcosystemAdapter = adapter
}

func TestJSAdapterContract(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "package.json"), []byte(`{"name":"demo"}`), 0o600); err != nil {
		t.Fatal(err)
	}

	adapter := &JSAdapter{}
	if diags := adapter.ValidateEnvironment(dir); len(diags) == 0 {
		t.Fatalf("expected environment diagnostics for a JavaScript project")
	}
	if diags := adapter.ValidateDependencies(dir); len(diags) != 0 {
		t.Fatalf("expected no dependency diagnostics for a valid package.json, got %v", diags)
	}
	if diags := adapter.ValidateProject(dir); len(diags) == 0 {
		t.Fatalf("expected project diagnostics for a detected JavaScript project")
	}
}
