package discovery

import (
	"context"
	"path/filepath"
	"strings"

	"github.com/LordCodex/promptengine/internal/eventbus"
	"github.com/LordCodex/promptengine/internal/filesystem"
	"github.com/LordCodex/promptengine/pkg/manifest"
)

type TechStage struct {
	Registry *RuleRegistry
	Manifest *manifest.Engine
	Events   *eventbus.EventBus
}

func (s *TechStage) Name() string { return "technology_detection" }

func (s *TechStage) Execute(ctx context.Context, fs filesystem.FileSystem, pm *ProjectModel) error {
	registry := s.Registry
	if registry == nil {
		registry = NewRuleRegistry()
	}
	if s.Manifest != nil {
		registry.RegisterManifestTechnologies(s.Manifest.ActiveManifest())
	}

	for _, rule := range registry.Detect(fs, pm.RootDir) {
		applyRule(pm, rule)
		s.publish(rule)
	}

	s.detectPackageJSON(fs, pm)
	s.detectComposer(fs, pm)
	s.detectEnvDatabases(fs, pm)
	s.detectConfigFiles(pm)
	pm.SyncLegacyFields()
	return nil
}

func (s *TechStage) publish(rule DetectionRule) {
	if s.Events == nil {
		return
	}
	s.Events.Publish(eventbus.Event{Type: eventbus.TechnologyDetected, Message: "technology detected", Payload: rule.Name})
}

func (s *TechStage) detectPackageJSON(fs filesystem.FileSystem, pm *ProjectModel) {
	if !fs.Exists(filepath.Join(pm.RootDir, "package.json")) {
		return
	}
	pm.Technology.Languages = addUnique(pm.Technology.Languages, "JavaScript")
	pm.Technology.Runtimes = addUnique(pm.Technology.Runtimes, "Node.js")
	pm.Technology.PackageManagers = addUnique(pm.Technology.PackageManagers, detectNodePackageManager(fs, pm.RootDir))
	if packageHas(fs, pm.RootDir, "typescript") {
		pm.Technology.Languages = addUnique(pm.Technology.Languages, "TypeScript")
	}
	frameworks := map[string][]string{
		"Vue":          {"vue"},
		"Nuxt":         {"nuxt"},
		"React":        {"react"},
		"Next.js":      {"next"},
		"Angular":      {"@angular/core"},
		"Express":      {"express"},
		"React Native": {"react-native"},
		"Inertia":      {"@inertiajs/vue3", "@inertiajs/react", "@inertiajs/svelte"},
	}
	for framework, deps := range frameworks {
		if packageHas(fs, pm.RootDir, deps...) {
			pm.Technology.Frameworks = addUnique(pm.Technology.Frameworks, framework)
			s.publish(DetectionRule{Name: framework})
		}
	}
}

func (s *TechStage) detectComposer(fs filesystem.FileSystem, pm *ProjectModel) {
	if !fs.Exists(filepath.Join(pm.RootDir, "composer.json")) {
		return
	}
	pm.Technology.Languages = addUnique(pm.Technology.Languages, "PHP")
	pm.Technology.PackageManagers = addUnique(pm.Technology.PackageManagers, "composer")
	if fs.Exists(filepath.Join(pm.RootDir, "artisan")) {
		pm.Technology.Frameworks = addUnique(pm.Technology.Frameworks, "Laravel")
		s.publish(DetectionRule{Name: "Laravel"})
	}
}

func (s *TechStage) detectEnvDatabases(fs filesystem.FileSystem, pm *ProjectModel) {
	for _, ef := range []string{".env", ".env.example", ".env.local"} {
		data, err := fs.ReadFile(filepath.Join(pm.RootDir, ef))
		if err != nil {
			continue
		}
		content := strings.ToLower(string(data))
		switch {
		case strings.Contains(content, "db_connection=pgsql"), strings.Contains(content, "postgres://"), strings.Contains(content, "postgresql://"):
			pm.Technology.Databases = addUnique(pm.Technology.Databases, "PostgreSQL")
		case strings.Contains(content, "db_connection=mysql"), strings.Contains(content, "mysql://"):
			pm.Technology.Databases = addUnique(pm.Technology.Databases, "MySQL")
		case strings.Contains(content, "mongodb://"), strings.Contains(content, "mongo_url"):
			pm.Technology.Databases = addUnique(pm.Technology.Databases, "MongoDB")
		case strings.Contains(content, "db_connection=sqlite"), strings.Contains(content, "sqlite://"):
			pm.Technology.Databases = addUnique(pm.Technology.Databases, "SQLite")
		}
		if strings.Contains(content, "redis_host=") || strings.Contains(content, "redis://") {
			pm.Technology.Databases = addUnique(pm.Technology.Databases, "Redis")
		}
		return
	}
}

func (s *TechStage) detectConfigFiles(pm *ProjectModel) {
	for _, file := range pm.Repository.ConfigurationFiles {
		switch filepath.Base(file) {
		case "phpunit.xml":
			pm.Technology.Testing = addUnique(pm.Technology.Testing, "PHPUnit")
		case ".gitlab-ci.yml":
			pm.Technology.Infrastructure = addUnique(pm.Technology.Infrastructure, "GitLab CI")
		}
		if strings.HasPrefix(file, ".github/workflows/") {
			pm.Technology.Infrastructure = addUnique(pm.Technology.Infrastructure, "GitHub Actions")
		}
		if strings.HasSuffix(file, ".tf") {
			pm.Technology.Infrastructure = addUnique(pm.Technology.Infrastructure, "Terraform")
		}
	}
}

func detectNodePackageManager(fs filesystem.FileSystem, root string) string {
	switch {
	case fs.Exists(filepath.Join(root, "pnpm-lock.yaml")):
		return "pnpm"
	case fs.Exists(filepath.Join(root, "yarn.lock")):
		return "yarn"
	default:
		return "npm"
	}
}
