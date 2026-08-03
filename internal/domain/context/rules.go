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
	if normalizeTask(req.TaskType) == "architecture_review" {
		return []string{"docs/Architecture.md"}
	}
	return nil
}

func (r *RuleSet) ConditionalDocuments(req ContextRequest) map[string]string {
	signals := requestSignals(req)
	docs := map[string]string{}
	task := normalizeTask(req.TaskType)
	if task == normalizeTask(TaskBugFix) || hasAnySignal(signals, "bug", "error", "failure", "troubleshoot", "incident") {
		docs["docs/Troubleshooting.md"] = "The task is a bug investigation or failure path, so troubleshooting history may contain directly relevant fixes."
	}
	if task == normalizeTask(TaskRefactor) || hasAnySignal(signals, "architecture", "service", "module", "boundary", "dependency", "refactor") {
		docs["docs/Architecture.md"] = "The task affects architecture, module boundaries, or dependencies."
	}
	if hasAnySignal(signals, "database", "migration", "model", "query", "persistence", "schema", "table", "repository", "eloquent") {
		docs["docs/Database.md"] = "The task mentions persistence, models, queries, migrations, or schema behavior."
	}
	if hasAnySignal(signals, "api", "endpoint", "route", "controller", "request", "response", "contract", "auth", "validation") {
		docs["docs/API.md"] = "The task affects endpoints, contracts, request validation, authentication, or responses."
	}
	if task == normalizeTask(TaskAddFeature) || hasAnySignal(signals, "business", "rule", "payment", "billing", "subscription", "refund", "policy") {
		docs["docs/BusinessRules.md"] = "The task may change domain behavior or business rules."
	}
	if hasAnySignal(signals, "deploy", "deployment", "environment", "infra", "docker", "kubernetes", "runtime", "config") {
		docs["docs/Deployment.md"] = "The task affects deployment, infrastructure, environment, or runtime configuration."
	}
	if hasAnySignal(signals, "security", "auth", "authorization", "secret", "password", "token", "permission", "payment", "personal", "input", "upload") {
		docs["docs/Security.md"] = "The task touches security-sensitive behavior such as auth, secrets, payments, personal data, user input, or file access."
	}
	if task == "documentation_update" || hasAnySignal(signals, "documentation", "docs", "readme") {
		docs["README.md"] = "The task explicitly affects documentation or onboarding material."
	}
	return docs
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

func requestSignals(req ContextRequest) string {
	return strings.ToLower(strings.Join([]string{
		string(req.TaskType),
		req.WorkflowType,
		req.UserIntent,
		req.RequestedOperation,
		strings.Join(req.AffectedFiles, " "),
	}, " "))
}

func hasAnySignal(signals string, terms ...string) bool {
	for _, term := range terms {
		if strings.Contains(signals, strings.ToLower(term)) {
			return true
		}
	}
	return false
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
