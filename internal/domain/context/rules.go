package context

import (
	"path/filepath"
	"strings"

	"github.com/LordCodex/promptengine/internal/domain/discovery"
)

type RuleSet struct{}

func NewRuleSet() *RuleSet {
	return &RuleSet{}
}

// RequiredDocuments returns only project-level instructions that are universally
// required for safe work. Task-specific project docs are handled as conditional
// candidates so unrelated documentation is not injected into every prompt.
func (r *RuleSet) RequiredDocuments(_ ContextRequest) []string {
	return []string{"AGENTS.md"}
}

// CandidateDocuments returns task-relevant documentation that may be included
// when the request itself provides a concrete reason to do so.
func (r *RuleSet) CandidateDocuments(req ContextRequest) []string {
	switch normalizeTask(req.TaskType) {
	case normalizeTask(TaskAddFeature), "feature", "new_feature", "feature_development":
		return []string{"docs/Architecture.md", "docs/Database.md", "docs/API.md", "docs/BusinessRules.md"}
	case normalizeTask(TaskBugFix), "bug_fix":
		return []string{"docs/Troubleshooting.md", "docs/Architecture.md"}
	case normalizeTask(TaskRefactor), "refactoring":
		return []string{"docs/Architecture.md", "docs/API.md"}
	case normalizeTask(TaskReview), "architecture_review":
		return []string{"docs/Architecture.md", "README.md"}
	case "documentation_update", "documentation_review":
		return []string{"docs/Architecture.md", "docs/API.md", "docs/Database.md"}
	default:
		return []string{"README.md"}
	}
}

// DocumentRelevance returns an additive score for a conditional documentation
// candidate. A document only becomes selectable when the concrete task or
// affected files indicate that the document materially helps the work.
func (r *RuleSet) DocumentRelevance(path string, req ContextRequest) float64 {
	intent := strings.ToLower(strings.Join([]string{
		req.UserIntent,
		req.RequestedOperation,
		strings.Join(req.AffectedFiles, " "),
	}, " "))
	lowerPath := strings.ToLower(filepath.ToSlash(path))
	task := normalizeTask(req.TaskType)

	switch {
	case strings.Contains(lowerPath, "businessrules"):
		if task == normalizeTask(TaskAddFeature) || task == "feature" || task == "new_feature" || task == "feature_development" {
			return 20
		}
		if containsAny(intent, "business", "rule", "requirement", "policy", "commission", "payment", "wallet", "order", "billing") {
			return 25
		}
	case strings.Contains(lowerPath, "troubleshooting"):
		if task == normalizeTask(TaskBugFix) || task == "bug_fix" {
			return 30
		}
	case strings.Contains(lowerPath, "architecture"):
		if task == normalizeTask(TaskRefactor) || task == "refactoring" || task == normalizeTask(TaskReview) || task == "architecture_review" {
			return 25
		}
		if containsAny(intent, "architecture", "architectural", "module", "boundary", "structure", "design", "layer", "dependency") {
			return 30
		}
	case strings.Contains(lowerPath, "database"):
		if containsAny(intent, "database", "schema", "migration", "table", "model", "query", "sql", "persist", "storage", "repository") || affectedPathContains(req.AffectedFiles, "database", "migration", "model", "repository") {
			return 30
		}
	case strings.Contains(lowerPath, "/api") || strings.HasSuffix(lowerPath, "api.md"):
		if containsAny(intent, "api", "endpoint", "route", "controller", "request", "response", "webhook", "rest", "graphql") || affectedPathContains(req.AffectedFiles, "route", "controller", "api") {
			return 30
		}
	case strings.HasSuffix(lowerPath, "readme.md"):
		if task == normalizeTask(TaskNewProject) || task == normalizeTask(TaskExistingProject) || task == normalizeTask(TaskReview) {
			return 20
		}
	}

	return 0
}

func (r *RuleSet) TechnologyPaths(pm *discovery.ProjectModel) []string {
	if pm == nil {
		return nil
	}
	var paths []string
	for _, framework := range pm.Frameworks {
		switch framework {
		case "Laravel":
			paths = append(paths, "routes", "app/Models", "app/Services", "app/Http/Controllers", "database")
		case "React":
			paths = append(paths, "src", "components", "hooks")
		case "Vue", "Nuxt":
			paths = append(paths, "src", "components", "composables", "pages")
		case "Next.js":
			paths = append(paths, "app", "pages", "components", "src")
		case "Flutter":
			paths = append(paths, "lib", "test")
		case "Go":
			paths = append(paths, "cmd", "internal", "pkg")
		}
	}
	for _, lang := range pm.Languages {
		if lang == "Go" {
			paths = append(paths, "cmd", "internal", "pkg")
		}
	}
	return unique(paths)
}

func (r *RuleSet) AffectedRelationships(req ContextRequest, pm *discovery.ProjectModel) []string {
	var out []string
	for _, file := range req.AffectedFiles {
		out = append(out, file)
		dir := filepath.Dir(file)
		base := filepath.Base(file)
		if dir != "." {
			out = append(out, dir)
		}
		if strings.HasSuffix(base, ".go") {
			out = append(out, strings.TrimSuffix(file, ".go")+"_test.go")
		}
		if strings.HasSuffix(base, ".php") {
			out = append(out, strings.Replace(file, "app/", "tests/", 1))
		}
		if strings.HasSuffix(base, ".tsx") || strings.HasSuffix(base, ".ts") || strings.HasSuffix(base, ".jsx") || strings.HasSuffix(base, ".js") {
			out = append(out, strings.TrimSuffix(file, filepath.Ext(file))+".test"+filepath.Ext(file))
		}
	}
	return unique(out)
}

func normalizeTask(task TaskType) string {
	return strings.ReplaceAll(strings.ToLower(string(task)), "-", "_")
}

func unique(items []string) []string {
	seen := map[string]bool{}
	var out []string
	for _, item := range items {
		if item == "" || seen[item] {
			continue
		}
		seen[item] = true
		out = append(out, item)
	}
	return out
}

func containsAny(value string, needles ...string) bool {
	for _, needle := range needles {
		if strings.Contains(value, needle) {
			return true
		}
	}
	return false
}

func affectedPathContains(paths []string, needles ...string) bool {
	for _, path := range paths {
		lower := strings.ToLower(filepath.ToSlash(path))
		if containsAny(lower, needles...) {
			return true
		}
	}
	return false
}
