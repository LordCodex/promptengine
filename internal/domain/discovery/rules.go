package discovery

import (
	"encoding/json"
	"path/filepath"
	"strings"

	"github.com/LordCodex/promptengine/internal/filesystem"
	"github.com/LordCodex/promptengine/pkg/manifest"
)

type DetectionRule struct {
	ID              string
	Name            string
	Kind            string
	Files           []string
	Directories     []string
	FileContains    map[string][]string
	Language        string
	Framework       string
	Runtime         string
	PackageManager  string
	Database        string
	Infrastructure  string
	Testing         string
	RelatedPlaybook []string
}

type RuleRegistry struct {
	rules []DetectionRule
}

func NewRuleRegistry() *RuleRegistry {
	r := &RuleRegistry{}
	r.Register(defaultRules()...)
	return r
}

func (r *RuleRegistry) Register(rules ...DetectionRule) {
	r.rules = append(r.rules, rules...)
}

func (r *RuleRegistry) Rules() []DetectionRule {
	return append([]DetectionRule(nil), r.rules...)
}

func (r *RuleRegistry) RegisterManifestTechnologies(m *manifest.Manifest) {
	if m == nil {
		return
	}
	for _, tech := range m.Technologies {
		r.Register(DetectionRule{
			ID:              "manifest:" + firstNonEmptyLocal(tech.ID, tech.Framework, tech.Language, tech.Stack),
			Name:            firstNonEmptyLocal(tech.Framework, tech.Language, tech.Stack, tech.ID),
			Kind:            "manifest",
			Language:        tech.Language,
			Framework:       tech.Framework,
			RelatedPlaybook: tech.RelatedPlaybooks,
		})
	}
}

func (r *RuleRegistry) Detect(fs filesystem.FileSystem, root string) []DetectionRule {
	var matched []DetectionRule
	for _, rule := range r.rules {
		if ruleMatches(fs, root, rule) {
			matched = append(matched, rule)
		}
	}
	return matched
}

func ruleMatches(fs filesystem.FileSystem, root string, rule DetectionRule) bool {
	if rule.Kind == "manifest" {
		return false
	}
	for _, file := range rule.Files {
		if !fs.Exists(filepath.Join(root, file)) {
			return false
		}
	}
	for _, dir := range rule.Directories {
		if !fs.Exists(filepath.Join(root, dir)) {
			return false
		}
	}
	for file, needles := range rule.FileContains {
		data, err := fs.ReadFile(filepath.Join(root, file))
		if err != nil {
			return false
		}
		content := strings.ToLower(string(data))
		found := false
		for _, needle := range needles {
			if strings.Contains(content, strings.ToLower(needle)) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return len(rule.Files) > 0 || len(rule.Directories) > 0 || len(rule.FileContains) > 0
}

func applyRule(pm *ProjectModel, rule DetectionRule) {
	pm.Technology.Languages = addUnique(pm.Technology.Languages, rule.Language)
	pm.Technology.Frameworks = addUnique(pm.Technology.Frameworks, rule.Framework)
	pm.Technology.Runtimes = addUnique(pm.Technology.Runtimes, rule.Runtime)
	pm.Technology.PackageManagers = addUnique(pm.Technology.PackageManagers, rule.PackageManager)
	pm.Technology.Databases = addUnique(pm.Technology.Databases, rule.Database)
	pm.Technology.Infrastructure = addUnique(pm.Technology.Infrastructure, rule.Infrastructure)
	pm.Technology.Testing = addUnique(pm.Technology.Testing, rule.Testing)
}

func packageHas(fs filesystem.FileSystem, root string, deps ...string) bool {
	data, err := fs.ReadFile(filepath.Join(root, "package.json"))
	if err != nil {
		return false
	}
	var pkg struct {
		Dependencies    map[string]any `json:"dependencies"`
		DevDependencies map[string]any `json:"devDependencies"`
	}
	if json.Unmarshal(data, &pkg) != nil {
		content := strings.ToLower(string(data))
		for _, dep := range deps {
			if strings.Contains(content, strings.ToLower(dep)) {
				return true
			}
		}
		return false
	}
	for _, dep := range deps {
		if _, ok := pkg.Dependencies[dep]; ok {
			return true
		}
		if _, ok := pkg.DevDependencies[dep]; ok {
			return true
		}
	}
	return false
}

func composerHas(fs filesystem.FileSystem, root string, packages ...string) bool {
	data, err := fs.ReadFile(filepath.Join(root, "composer.json"))
	if err != nil {
		return false
	}
	var composer struct {
		Require    map[string]any `json:"require"`
		RequireDev map[string]any `json:"require-dev"`
	}
	if json.Unmarshal(data, &composer) != nil {
		content := strings.ToLower(string(data))
		for _, pkg := range packages {
			if strings.Contains(content, strings.ToLower(pkg)) {
				return true
			}
		}
		return false
	}
	for _, pkg := range packages {
		if _, ok := composer.Require[pkg]; ok {
			return true
		}
		if _, ok := composer.RequireDev[pkg]; ok {
			return true
		}
	}
	return false
}

func defaultRules() []DetectionRule {
	return []DetectionRule{
		{ID: "go", Name: "Go", Files: []string{"go.mod"}, Language: "Go", Runtime: "Go", PackageManager: "go mod"},
		{ID: "laravel", Name: "Laravel", Files: []string{"composer.json", "artisan"}, Language: "PHP", Framework: "Laravel", Runtime: "PHP", PackageManager: "composer"},
		{ID: "django", Name: "Django", FileContains: map[string][]string{"requirements.txt": {"django"}}, Language: "Python", Framework: "Django", Runtime: "Python", PackageManager: "pip"},
		{ID: "rails", Name: "Rails", Files: []string{"Gemfile"}, FileContains: map[string][]string{"Gemfile": {"rails"}}, Language: "Ruby", Framework: "Rails", Runtime: "Ruby", PackageManager: "bundler"},
		{ID: "spring", Name: "Spring", Files: []string{"pom.xml"}, FileContains: map[string][]string{"pom.xml": {"spring-boot"}}, Language: "Java", Framework: "Spring", Runtime: "JVM", PackageManager: "maven"},
		{ID: "flutter", Name: "Flutter", Files: []string{"pubspec.yaml"}, FileContains: map[string][]string{"pubspec.yaml": {"flutter:"}}, Language: "Dart", Framework: "Flutter", Runtime: "Dart", PackageManager: "pub"},
		{ID: "docker", Name: "Docker", Files: []string{"Dockerfile"}, Infrastructure: "Docker"},
		{ID: "compose", Name: "Docker Compose", Files: []string{"docker-compose.yml"}, Infrastructure: "Docker Compose"},
		{ID: "github-actions", Name: "GitHub Actions", Directories: []string{".github/workflows"}, Infrastructure: "GitHub Actions"},
		{ID: "kubernetes", Name: "Kubernetes", Directories: []string{"k8s"}, Infrastructure: "Kubernetes"},
	}
}

func firstNonEmptyLocal(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}
