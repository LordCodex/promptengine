package installer

import (
	"testing"

	"github.com/LordCodex/promptengine/internal/filesystem"
)

func TestLocalInstaller_InstallAndUninstall(t *testing.T) {
	fs := filesystem.NewMockFileSystem()
	// Create source file declared by the package
	_ = fs.WriteFile("source/my-plugin/plugin.json", []byte(`{}`), 0644)

	inst := NewLocalInstaller(fs, ".promptengine/plugins")

	manifest := PackageManifest{
		ID:      "my-plugin",
		Name:    "My Plugin",
		Kind:    KindPlugin,
		Version: "1.0.0",
		Files:   []string{"source/my-plugin/plugin.json"},
	}

	record, err := inst.Install(manifest)
	if err != nil {
		t.Fatalf("expected successful install, got: %v", err)
	}
	if record.Manifest.ID != "my-plugin" {
		t.Errorf("expected install record ID 'my-plugin', got '%s'", record.Manifest.ID)
	}
	if !inst.IsInstalled("my-plugin") {
		t.Error("expected IsInstalled to return true after install")
	}

	// Duplicate install must fail
	if _, err = inst.Install(manifest); err == nil {
		t.Error("expected error on duplicate install, got nil")
	}

	// Uninstall
	if err = inst.Uninstall("my-plugin"); err != nil {
		t.Errorf("expected clean uninstall, got: %v", err)
	}
}

func TestLocalInstaller_EnableDisableUpgrade(t *testing.T) {
	fs := filesystem.NewMockFileSystem()
	fs.WriteFile("plugin.json", []byte("{}"), 0644)
	installer := NewLocalInstaller(fs, ".promptengine/plugins")
	manifest := PackageManifest{ID: "company", Version: "1.0.0", Files: []string{"plugin.json"}}
	if _, err := installer.Install(manifest); err != nil {
		t.Fatalf("install failed: %v", err)
	}
	if err := installer.Enable("company"); err != nil {
		t.Fatalf("enable failed: %v", err)
	}
	if err := installer.Disable("company"); err != nil {
		t.Fatalf("disable failed: %v", err)
	}
	manifest.Version = "2.0.0"
	if _, err := installer.Upgrade(manifest); err != nil {
		t.Fatalf("upgrade failed: %v", err)
	}
}

func TestLocalInstaller_MissingSourceFile(t *testing.T) {
	fs := filesystem.NewMockFileSystem()
	inst := NewLocalInstaller(fs, ".promptengine/plugins")

	manifest := PackageManifest{
		ID:    "bad-plugin",
		Files: []string{"non/existent/file.json"},
	}
	if _, err := inst.Install(manifest); err == nil {
		t.Error("expected error when declared source file is missing, got nil")
	}
}

func TestLocalInstaller_UninstallNotInstalled(t *testing.T) {
	fs := filesystem.NewMockFileSystem()
	inst := NewLocalInstaller(fs, ".promptengine/plugins")
	if err := inst.Uninstall("ghost-plugin"); err == nil {
		t.Error("expected error when uninstalling a package that is not installed")
	}
}
