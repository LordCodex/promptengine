package discovery

import (
	"context"
	"path/filepath"
	"strings"

	"github.com/LordCodex/promptengine/internal/filesystem"
)

type RepositoryScanner struct {
	MaxFiles int
}

func NewRepositoryScanner() *RepositoryScanner {
	return &RepositoryScanner{MaxFiles: 10000}
}

func (s *RepositoryScanner) Scan(ctx context.Context, fs filesystem.FileSystem, root string) (RepositoryInfo, error) {
	if root == "" {
		root = "."
	}
	info := RepositoryInfo{RootPath: root}
	if !fs.Exists(root) && root != "." {
		return info, nil
	}
	err := s.walk(ctx, fs, filepath.Clean(root), filepath.Clean(root), &info)
	info.IsMonorepo = detectMonorepo(info)
	return info, err
}

func (s *RepositoryScanner) walk(ctx context.Context, fs filesystem.FileSystem, root, dir string, info *RepositoryInfo) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	if s.MaxFiles > 0 && len(info.Files) >= s.MaxFiles {
		return nil
	}
	entries, err := fs.ReadDir(dir)
	if err != nil {
		info.PermissionErrors = append(info.PermissionErrors, dir)
		return nil
	}
	for _, entry := range entries {
		name := entry.Name()
		full := filepath.Join(dir, name)
		rel, relErr := filepath.Rel(root, full)
		if relErr != nil {
			rel = full
		}
		rel = filepath.ToSlash(rel)
		if shouldIgnore(name, rel) {
			info.IgnoredFiles = append(info.IgnoredFiles, rel)
			continue
		}
		if entry.IsDir() {
			info.Directories = append(info.Directories, rel)
			if err := s.walk(ctx, fs, root, full, info); err != nil {
				return err
			}
			continue
		}
		info.Files = append(info.Files, rel)
		if isConfigFile(name, rel) {
			info.ConfigurationFiles = append(info.ConfigurationFiles, rel)
		}
		if isDocumentationFile(name, rel) {
			info.DocumentationFiles = append(info.DocumentationFiles, rel)
		}
	}
	return nil
}

func shouldIgnore(name, rel string) bool {
	ignored := map[string]bool{
		".git": true, "node_modules": true, "vendor": true, "build": true,
		"dist": true, ".next": true, ".nuxt": true, ".dart_tool": true,
		"coverage": true, ".idea": true, ".vscode": true, ".gocache": true,
	}
	if ignored[name] {
		return true
	}
	return strings.HasPrefix(rel, ".git/")
}

func isConfigFile(name, rel string) bool {
	configNames := map[string]bool{
		"go.mod": true, "package.json": true, "composer.json": true, "pubspec.yaml": true,
		"requirements.txt": true, "Pipfile": true, "Gemfile": true, "pom.xml": true,
		"Dockerfile": true, "docker-compose.yml": true, "docker-compose.yaml": true,
		".env": true, ".env.example": true, ".gitlab-ci.yml": true, "Cargo.toml": true,
		"kustomization.yaml": true,
	}
	if configNames[name] {
		return true
	}
	return strings.HasPrefix(rel, ".github/workflows/") || strings.HasSuffix(name, ".tf") || strings.HasSuffix(name, ".k8s.yaml")
}

func isDocumentationFile(name, rel string) bool {
	lower := strings.ToLower(name)
	if lower == "readme.md" || lower == "agents.md" || strings.HasSuffix(lower, ".md") {
		return true
	}
	return strings.HasPrefix(rel, "docs/")
}

func detectMonorepo(info RepositoryInfo) bool {
	roots := 0
	for _, file := range info.Files {
		switch filepath.Base(file) {
		case "package.json", "composer.json", "go.mod", "pubspec.yaml", "pom.xml":
			if strings.Contains(file, "/") {
				roots++
			}
		}
	}
	for _, dir := range info.Directories {
		if dir == "packages" || dir == "apps" || dir == "services" {
			return true
		}
	}
	return roots >= 2
}
