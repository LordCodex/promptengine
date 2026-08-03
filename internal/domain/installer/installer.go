package installer

import (
	"fmt"
	"path/filepath"

	"github.com/LordCodex/promptengine/internal/filesystem"
)

// PackageKind identifies what is being installed
type PackageKind string

const (
	KindPlugin         PackageKind = "plugin"
	KindTechnologyPack PackageKind = "technology-pack"
	KindOrgPack        PackageKind = "org-pack"
	KindTemplate       PackageKind = "template"
	KindWorkflowPack   PackageKind = "workflow-pack"
	KindPromptLibrary  PackageKind = "prompt-library"
	KindAIProvider     PackageKind = "ai-provider"
	KindStandard       PackageKind = "standard"
)

// PackageManifest is the on-disk descriptor for any installable package
type PackageManifest struct {
	ID           string      `json:"id"`
	Name         string      `json:"name"`
	Kind         PackageKind `json:"kind"`
	Version      string      `json:"version"`
	Author       string      `json:"author"`
	Description  string      `json:"description"`
	Dependencies []string    `json:"dependencies"`
	MinCoreVer   string      `json:"min_core_ver"`
	Files        []string    `json:"files"` // relative paths to install
}

// InstallRecord stores what was installed and where (for rollback)
type InstallRecord struct {
	Manifest    PackageManifest
	InstalledAt string
	Files       []string
	Enabled     bool
}

// Installer is the common installation contract
type Installer interface {
	Install(manifest PackageManifest) (*InstallRecord, error)
	Uninstall(id string) error
	IsInstalled(id string) bool
	Enable(id string) error
	Disable(id string) error
	Upgrade(manifest PackageManifest) (*InstallRecord, error)
}

// LocalInstaller installs packages from the local filesystem
type LocalInstaller struct {
	fs      filesystem.FileSystem
	destDir string // base installation directory (e.g. ".promptengine/plugins")
	records map[string]*InstallRecord
}

func NewLocalInstaller(fs filesystem.FileSystem, destDir string) *LocalInstaller {
	return &LocalInstaller{
		fs:      fs,
		destDir: destDir,
		records: make(map[string]*InstallRecord),
	}
}

func (i *LocalInstaller) Install(manifest PackageManifest) (*InstallRecord, error) {
	if i.IsInstalled(manifest.ID) {
		return nil, fmt.Errorf("package '%s' is already installed", manifest.ID)
	}

	// Verify declared files exist (source validation)
	for _, f := range manifest.Files {
		if !i.fs.Exists(f) || i.fs.IsDir(f) {
			return nil, fmt.Errorf("install failed: declared file '%s' not found for package '%s'", f, manifest.ID)
		}
	}
	for _, f := range manifest.Files {
		data, err := i.fs.ReadFile(f)
		if err != nil {
			return nil, fmt.Errorf("install failed: read declared file '%s': %w", f, err)
		}
		destination := filepath.Join(i.destDir, manifest.ID, filepath.Base(filepath.Clean(f)))
		if err := i.fs.WriteFile(destination, data, 0644); err != nil {
			return nil, fmt.Errorf("install failed: write package file '%s': %w", f, err)
		}
	}

	// Write marker file so IsInstalled works without a DB
	markerPath := fmt.Sprintf("%s/%s/.installed", i.destDir, manifest.ID)
	if err := i.fs.WriteFile(markerPath, []byte(manifest.Version), 0644); err != nil {
		return nil, fmt.Errorf("failed to write install marker: %w", err)
	}

	record := &InstallRecord{
		Manifest: manifest,
		Files:    manifest.Files,
	}
	i.records[manifest.ID] = record
	return record, nil
}

func (i *LocalInstaller) Uninstall(id string) error {
	if !i.IsInstalled(id) {
		return fmt.Errorf("package '%s' is not installed", id)
	}
	if err := i.fs.RemoveAll(filepath.Join(i.destDir, id)); err != nil {
		return fmt.Errorf("failed to remove package '%s': %w", id, err)
	}
	delete(i.records, id)
	return nil
}

func (i *LocalInstaller) IsInstalled(id string) bool {
	markerPath := fmt.Sprintf("%s/%s/.installed", i.destDir, id)
	if !i.fs.Exists(markerPath) {
		return false
	}
	data, err := i.fs.ReadFile(markerPath)
	if err != nil {
		return false
	}
	return string(data) != "uninstalled"
}

func (i *LocalInstaller) Enable(id string) error {
	if !i.IsInstalled(id) {
		return fmt.Errorf("package '%s' is not installed", id)
	}
	if record, ok := i.records[id]; ok {
		record.Enabled = true
	}
	return i.fs.WriteFile(fmt.Sprintf("%s/%s/.enabled", i.destDir, id), []byte("enabled"), 0644)
}

func (i *LocalInstaller) Disable(id string) error {
	if !i.IsInstalled(id) {
		return fmt.Errorf("package '%s' is not installed", id)
	}
	if record, ok := i.records[id]; ok {
		record.Enabled = false
	}
	return i.fs.WriteFile(fmt.Sprintf("%s/%s/.enabled", i.destDir, id), []byte("disabled"), 0644)
}

func (i *LocalInstaller) Upgrade(manifest PackageManifest) (*InstallRecord, error) {
	if !i.IsInstalled(manifest.ID) {
		return nil, fmt.Errorf("package '%s' is not installed", manifest.ID)
	}
	_ = i.Uninstall(manifest.ID)
	return i.Install(manifest)
}
