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

func (r *RuleSet) RequiredDocuments(req ContextRequest) []string {
	switch normalizeTask(req.TaskType) {
	case normalizeTask(TaskAddFeature), "feature", "new_feature":
		return []string{"docs/Architecture.md", "docs/Database.md", "docs/API.md", "docs/BusinessRules.md"}
	case normalizeTask(TaskBugFix), "bug_fix":
		return []string{"docs/Troubleshooting.md", "docs/Architecture.md"}
	case normalizeTask(TaskRefactor):
		return []string{"docs/Architecture.md", "docs/API.md"}
	case normalizeTask(TaskReview), "architecture_review":
		return []string{"docs/Architecture.md", "README.md"}
	case "documentation_update":
		return []string{"docs/Architecture.md", "docs/API.md", "docs/Database.md"}
	default:
		return []string{"README.md", "AGENTS.md"}
	}
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
