package context

import "strings"

// inferredTaskTypes keeps the user's primary workflow intact while deriving
// additional standards domains from concrete request language. This lets a
// generic feature task activate payment/API/security/etc. playbooks without an
// AI classification call.
func inferredTaskTypes(req ContextRequest) []string {
	primary := normalizeTask(req.TaskType)
	intent := strings.ToLower(strings.Join([]string{
		string(req.TaskType),
		req.WorkflowType,
		req.UserIntent,
		req.RequestedOperation,
		strings.Join(req.AffectedFiles, " "),
	}, " "))

	seen := map[string]bool{}
	var out []string
	add := func(task string) {
		task = normalizeTask(TaskType(task))
		if task == "" || seen[task] {
			return
		}
		seen[task] = true
		out = append(out, task)
	}
	add(primary)

	infer := []struct {
		task     string
		keywords []string
	}{
		{"authentication", []string{"authentication", "login", "sign in", "signin", "password reset", "two factor", "2fa", "session"}},
		{"authorization", []string{"authorization", "permission", "permissions", "role", "roles", "access control", "policy"}},
		{"api_development", []string{"api", "endpoint", "webhook", "route", "controller", "request", "response", "graphql", "rest"}},
		{"database_design", []string{"database", "schema", "migration", "table", "foreign key", "relationship", "model"}},
		{"database_optimization", []string{"slow query", "query performance", "database performance", "index", "indexes", "explain plan"}},
		{"payment_integration", []string{"payment", "payments", "refund", "withdrawal", "payout", "stripe", "paystack", "flutterwave", "monnify", "wallet", "billing"}},
		{"file_upload", []string{"file upload", "upload file", "uploads", "attachment", "image upload"}},
		{"queue_processing", []string{"queue", "queued", "worker", "job queue"}},
		{"background_jobs", []string{"background job", "background task", "scheduled job", "scheduler", "cron"}},
		{"notifications", []string{"notification", "notifications", "push notification"}},
		{"email", []string{"email", "mailer", "mailing", "smtp"}},
		{"search", []string{"search", "full text", "full-text", "filter results"}},
		{"dashboard", []string{"dashboard", "analytics dashboard", "admin dashboard"}},
		{"admin_panel", []string{"admin panel", "admin portal", "back office"}},
		{"mobile_screen", []string{"mobile screen", "flutter screen", "app screen"}},
		{"frontend_component", []string{"component", "vue component", "react component", "frontend component"}},
		{"ui_design", []string{"ui", "ux", "interface design", "responsive", "accessibility"}},
		{"testing", []string{"test", "tests", "testing", "coverage", "regression"}},
		{"security_review", []string{"security review", "vulnerability", "threat model", "xss", "csrf", "injection", "secret exposure"}},
		{"performance_optimization", []string{"performance", "optimize", "optimization", "latency", "memory usage", "slow"}},
		{"deployment", []string{"deploy", "deployment", "production release", "hosting"}},
		{"ci_cd", []string{"ci/cd", "ci cd", "github actions", "pipeline", "continuous integration"}},
		{"refactor", []string{"refactor", "restructure", "cleanup architecture"}},
		{"bug_fix", []string{"bug", "fix error", "broken", "regression", "exception", "failing"}},
		{"code_review", []string{"code review", "review changes", "review code", "pull request review"}},
	}

	for _, candidate := range infer {
		if containsIntentPhrase(intent, candidate.keywords...) {
			add(candidate.task)
		}
	}
	return out
}

func containsIntentPhrase(intent string, phrases ...string) bool {
	for _, phrase := range phrases {
		if strings.Contains(intent, phrase) {
			return true
		}
	}
	return false
}
