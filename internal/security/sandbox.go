package security

import (
	"fmt"
	"path/filepath"
	"strings"
)

// ValidateSafePath checks if the target path is strictly contained within the base directory.
// This prevents directory traversal attacks.
func ValidateSafePath(baseDir, targetPath string) (string, error) {
	absBase, err := filepath.Abs(baseDir)
	if err != nil {
		return "", fmt.Errorf("invalid base path: %w", err)
	}

	absTarget, err := filepath.Abs(targetPath)
	if err != nil {
		return "", fmt.Errorf("invalid target path: %w", err)
	}

	rel, err := filepath.Rel(absBase, absTarget)
	if err != nil {
		return "", fmt.Errorf("failed to compute relative path: %w", err)
	}

	if strings.HasPrefix(rel, "..") || strings.HasPrefix(rel, "/") {
		return "", fmt.Errorf("path traversal detected: path '%s' goes outside base directory '%s'", targetPath, baseDir)
	}

	return absTarget, nil
}
