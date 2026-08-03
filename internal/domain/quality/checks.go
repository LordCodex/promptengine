package quality

import (
	"encoding/json"

	"github.com/LordCodex/promptengine/internal/filesystem"
)

// CheckManifestExists verifies if playbook-manifest.json exists in workspace
func CheckManifestExists(fs filesystem.FileSystem) bool {
	return fs.Exists("playbook-manifest.json")
}

// CheckConfigExists verifies if .promptengine directory exists
func CheckConfigExists(fs filesystem.FileSystem) bool {
	return fs.Exists(".promptengine")
}

// CheckDocsDirExists verifies if docs directory exists
func CheckDocsDirExists(fs filesystem.FileSystem) bool {
	return fs.Exists("docs")
}

// CheckManifestNonEmpty verifies if playbook-manifest.json is non-empty and valid
func CheckManifestNonEmpty(fs filesystem.FileSystem) bool {
	if !CheckManifestExists(fs) {
		return false
	}
	data, err := fs.ReadFile("playbook-manifest.json")
	if err != nil {
		return false
	}
	// Verify it contains minimally structured JSON
	var val interface{}
	if err := json.Unmarshal(data, &val); err != nil {
		return false
	}
	return len(data) >= 2
}

// CheckCoreDocsComplete checks if required core docs exist
func CheckCoreDocsComplete(fs filesystem.FileSystem) []string {
	required := []string{
		"docs/Architecture.md",
		"docs/Database.md",
		"docs/API.md",
	}
	var missing []string
	for _, p := range required {
		if !fs.Exists(p) {
			missing = append(missing, p)
		}
	}
	return missing
}
