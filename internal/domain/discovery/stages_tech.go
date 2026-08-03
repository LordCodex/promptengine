package discovery

import (
	"context"
	"path/filepath"
	"strings"

	"github.com/LordCodex/promptengine/internal/filesystem"
)

// TechStage detects active languages, frameworks, package managers, and databases
type TechStage struct{}

func (s *TechStage) Name() string { return "tech_stage" }
func (s *TechStage) Execute(ctx context.Context, fs filesystem.FileSystem, pm *ProjectModel) error {
	// 1. Language and Package Manager files scanning
	if fs.Exists(filepath.Join(pm.RootDir, "go.mod")) {
		pm.Languages = append(pm.Languages, "Go")
		pm.PackageManagers = append(pm.PackageManagers, "go mod")
	}

	if fs.Exists(filepath.Join(pm.RootDir, "composer.json")) {
		pm.Languages = append(pm.Languages, "PHP")
		pm.PackageManagers = append(pm.PackageManagers, "composer")
	}

	if fs.Exists(filepath.Join(pm.RootDir, "package.json")) {
		pm.Languages = append(pm.Languages, "JavaScript")
		pm.PackageManagers = append(pm.PackageManagers, "npm")

		// Let's inspect package.json to trace TypeScript or specific frameworks
		data, err := fs.ReadFile(filepath.Join(pm.RootDir, "package.json"))
		if err == nil {
			content := string(data)
			if strings.Contains(content, `"typescript"`) {
				pm.Languages = append(pm.Languages, "TypeScript")
			}
			if strings.Contains(content, `"vue"`) {
				pm.Frameworks = append(pm.Frameworks, "Vue")
			}
			if strings.Contains(content, `"nuxt"`) {
				pm.Frameworks = append(pm.Frameworks, "Nuxt")
			}
			if strings.Contains(content, `"react"`) {
				pm.Frameworks = append(pm.Frameworks, "React")
			}
			if strings.Contains(content, `"next"`) {
				pm.Frameworks = append(pm.Frameworks, "Next.js")
			}
			if strings.Contains(content, `"express"`) {
				pm.Frameworks = append(pm.Frameworks, "Express")
			}
			if strings.Contains(content, `"@nestjs/core"`) {
				pm.Frameworks = append(pm.Frameworks, "NestJS")
			}
		}
	}

	if fs.Exists(filepath.Join(pm.RootDir, "pubspec.yaml")) {
		pm.Languages = append(pm.Languages, "Dart")
		pm.Frameworks = append(pm.Frameworks, "Flutter")
		pm.PackageManagers = append(pm.PackageManagers, "pub")
	}

	if fs.Exists(filepath.Join(pm.RootDir, "requirements.txt")) || fs.Exists(filepath.Join(pm.RootDir, "Pipfile")) {
		pm.Languages = append(pm.Languages, "Python")
		pm.PackageManagers = append(pm.PackageManagers, "pip")
	}

	if fs.Exists(filepath.Join(pm.RootDir, "Cargo.toml")) {
		pm.Languages = append(pm.Languages, "Rust")
		pm.PackageManagers = append(pm.PackageManagers, "cargo")
	}

	// 2. Specific Framework checks
	if fs.Exists(filepath.Join(pm.RootDir, "artisan")) {
		pm.Frameworks = append(pm.Frameworks, "Laravel")
	}

	// 3. Database detection from configs/env files
	envFiles := []string{".env", ".env.example", ".env.local"}
	for _, ef := range envFiles {
		fullPath := filepath.Join(pm.RootDir, ef)
		if fs.Exists(fullPath) {
			data, err := fs.ReadFile(fullPath)
			if err == nil {
				content := string(data)
				if strings.Contains(content, "DB_CONNECTION=pgsql") || strings.Contains(content, "postgres://") {
					pm.Databases = append(pm.Databases, "PostgreSQL")
				}
				if strings.Contains(content, "DB_CONNECTION=mysql") || strings.Contains(content, "mysql://") {
					pm.Databases = append(pm.Databases, "MySQL")
				}
				if strings.Contains(content, "DB_CONNECTION=sqlite") || strings.Contains(content, "sqlite://") {
					pm.Databases = append(pm.Databases, "SQLite")
				}
				if strings.Contains(content, "REDIS_HOST=") || strings.Contains(content, "redis://") {
					pm.Databases = append(pm.Databases, "Redis")
				}
			}
			break
		}
	}

	// 4. CI/CD pipelines
	if fs.Exists(filepath.Join(pm.RootDir, ".github/workflows")) {
		pm.CIs = append(pm.CIs, "GitHub Actions")
	}
	if fs.Exists(filepath.Join(pm.RootDir, ".gitlab-ci.yml")) {
		pm.CIs = append(pm.CIs, "GitLab CI")
	}

	// 5. Testing Frameworks
	if fs.Exists(filepath.Join(pm.RootDir, "phpunit.xml")) {
		pm.TestingFrames = append(pm.TestingFrames, "PHPUnit")
	}

	return nil
}
