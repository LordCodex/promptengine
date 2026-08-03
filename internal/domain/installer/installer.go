package installer

import (
	"fmt"

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
}

// Installer is the common installation contract
type Installer interface {
	Install(manifest PackageManifest) (*InstallRecord, error)
	Uninstall(id string) error
	IsInstalled(id string) bool
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
		if !i.fs.Exists(f) {
			return nil, fmt.Errorf("install failed: declared file '%s' not found for package '%s'", f, manifest.ID)
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
	markerPath := fmt.Sprintf("%s/%s/.installed", i.destDir, id)
	// Overwrite marker to signal removal (full deletion handled by fs layer)
	_ = i.fs.WriteFile(markerPath, []byte("uninstalled"), 0644)
	delete(i.records, id)
	return nil
}

func (i *LocalInstaller) IsInstalled(id string) bool {
	markerPath := fmt.Sprintf("%s/%s/.installed", i.destDir, id)
	return i.fs.Exists(markerPath)
}
