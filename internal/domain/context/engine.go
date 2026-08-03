package context

import (
	"context"
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/LordCodex/promptengine/internal/domain/discovery"
	"github.com/LordCodex/promptengine/internal/eventbus"
	"github.com/LordCodex/promptengine/internal/filesystem"
	"github.com/LordCodex/promptengine/internal/security"
	"github.com/LordCodex/promptengine/pkg/manifest"
)

type Engine struct {
	fs       filesystem.FileSystem
	manifest *manifest.QueryEngine
	events   *eventbus.EventBus
	rules    *RuleSet
	cache    *Cache
}

type EngineOption func(*Engine)

func WithManifestQuery(q *manifest.QueryEngine) EngineOption {
	return func(e *Engine) { e.manifest = q }
}

func WithEventBus(events *eventbus.EventBus) EngineOption {
	return func(e *Engine) { e.events = events }
}

func WithCache(cache *Cache) EngineOption {
	return func(e *Engine) { e.cache = cache }
}

func NewEngine(fs filesystem.FileSystem, opts ...EngineOption) *Engine {
	e := &Engine{fs: fs, rules: NewRuleSet(), cache: NewCache()}
	for _, opt := range opts {
		opt(e)
	}
	return e
}

func (e *Engine) Build(ctx context.Context, req ContextRequest) (*ContextPackage, error) {
	e.publish(eventbus.ContextBuildStarted, "context build started", req)
	if err := ctx.Err(); err != nil {
		e.publish(eventbus.ContextBuildFailed, "context build failed", err)
		return nil, err
	}
	if req.Budget == "" {
		req.Budget = BudgetSmall
	}
	limit := req.MaxBytes
	if limit <= 0 {
		limit = GetBudgetLimit(req.Budget)
	}
	threshold := req.MinRelevanceScore
	if threshold <= 0 {
		threshold = 45
	}

	candidates := e.gatherCandidates(req)
	paths := make([]string, 0, len(candidates))
	for _, item := range candidates {
		if item.Path != "" {
			paths = append(paths, item.Path)
		}
	}
	sort.Strings(paths)
	fp := fingerprint(e.fs, paths, e.cacheSalt(req))
	key := cacheKey(req)
	if cached, ok := e.cache.Get(key, fp); ok {
		return cached, nil
	}

	ranked := e.rankCandidates(candidates, req)
	pkg := NewContextPackage(string(req.TaskType), req.Budget)
	pkg.WorkflowType = req.WorkflowType
	if req.Project != nil {
		pkg.ProjectMetadata["root"] = req.Project.RootDir
		pkg.ProjectMetadata["detected_type"] = req.Project.Project.DetectedType
	}
	e.optimize(pkg, ranked, limit, threshold, req.Explain)
	pkg.SystemPrompt = e.renderSystemPrompt(pkg)
	e.cache.Set(key, fp, pkg)
	e.publish(eventbus.ContextBuilt, "context built", pkg)
	return pkg, nil
}

func (e *Engine) GenerateContext(ctx context.Context, task TaskType, pm *discovery.ProjectModel, budget BudgetType) (*ContextPackage, error) {
	return e.Build(ctx, ContextRequest{TaskType: task, WorkflowType: string(task), Project: pm, Budget: budget})
}

func (e *Engine) gatherCandidates(req ContextRequest) []ContextItem {
	var items []ContextItem
	addPath := func(path string, typ ContextSourceType, score float64, level InclusionLevel, reason string) {
		if path == "" {
			return
		}
		if security.IsSensitivePath(path) {
			return
		}
		data, err := e.fs.ReadFile(path)
		if err != nil {
			return
		}
		content, redacted := security.RedactSecrets(string(data))
		if redacted {
			reason += " Sensitive values redacted."
		}
		items = append(items, ContextItem{Path: path, Type: typ, RelevanceScore: score, Reason: reason, Size: len(content), Content: content, InclusionLevel: level})
	}

	for _, path := range e.rules.RequiredDocuments(req) {
		addPath(path, ContextDocumentation, 90, InclusionRequired, "The requested task cannot be completed correctly without this document.")
	}
	for path, reason := range e.rules.ConditionalDocuments(req) {
		addPath(path, ContextDocumentation, 62, InclusionConditional, reason)
	}
	if req.Project != nil {
		if req.Project.PromptEngine.AgentsMDPresent {
			addPath("AGENTS.md", ContextDocumentation, 88, InclusionRequired, "Project agent instructions define repository-specific constraints for this task.")
		}
		for _, rel := range e.rules.TechnologyPaths(req.Project) {
			items = append(items, e.itemsUnderPath(req.Project, rel, 20, InclusionConditional, "Technology path candidate; included only if task signals match this file.")...)
		}
		for _, file := range req.Project.Repository.Files {
			if item := e.fileCandidate(file, 18, InclusionConditional, "Repository file candidate; included only if task-specific relevance is proven."); item.Path != "" {
				items = append(items, item)
			}
		}
	}
	for _, path := range e.rules.AffectedRelationships(req, req.Project) {
		addPath(path, ContextFile, 92, InclusionRequired, "The file is explicitly affected or is a directly related source/test counterpart.")
	}
	if e.manifest != nil {
		for _, playbook := range e.manifest.PlaybooksByTask(string(req.TaskType)) {
			score, reason := playbookRelevance(req, playbook)
			addPath(playbook.Location, ContextStandard, score, InclusionConditional, reason)
			items = append(items, ContextItem{Path: playbook.ID, Type: ContextManifestEntry, RelevanceScore: score - 6, Reason: reason, InclusionLevel: InclusionConditional})
		}
		for _, prompt := range e.manifest.PromptsByWorkflow(firstNonEmpty(req.WorkflowType, string(req.TaskType))) {
			items = append(items, ContextItem{Path: prompt.PromptTemplate, Type: ContextWorkflow, RelevanceScore: 50, Reason: "Manifest prompt mapping matches the requested workflow.", InclusionLevel: InclusionConditional})
		}
	}
	return items
}

func (e *Engine) itemsUnderPath(pm *discovery.ProjectModel, prefix string, score float64, level InclusionLevel, reason string) []ContextItem {
	var out []ContextItem
	cleanPrefix := strings.Trim(filepath.ToSlash(prefix), "/")
	for _, file := range pm.Repository.Files {
		if file == cleanPrefix || strings.HasPrefix(file, cleanPrefix+"/") {
			if security.IsSensitivePath(file) {
				continue
			}
			data, err := e.fs.ReadFile(file)
			if err != nil {
				continue
			}
			content, redacted := security.RedactSecrets(string(data))
			itemReason := reason
			if redacted {
				itemReason += " Sensitive values redacted."
			}
			out = append(out, ContextItem{Path: file, Type: ContextFile, RelevanceScore: score, Reason: itemReason, Size: len(content), Content: content, InclusionLevel: level})
		}
	}
	return out
}

func (e *Engine) fileCandidate(path string, score float64, level InclusionLevel, reason string) ContextItem {
	if path == "" || security.IsSensitivePath(path) {
		return ContextItem{}
	}
	data, err := e.fs.ReadFile(path)
	if err != nil {
		return ContextItem{}
	}
	content, redacted := security.RedactSecrets(string(data))
	if redacted {
		reason += " Sensitive values redacted."
	}
	return ContextItem{Path: path, Type: ContextFile, RelevanceScore: score, Reason: reason, Size: len(content), Content: content, InclusionLevel: level}
}

func (e *Engine) rankCandidates(items []ContextItem, req ContextRequest) []ContextItem {
	intent := strings.ToLower(req.UserIntent + " " + req.RequestedOperation + " " + strings.Join(req.AffectedFiles, " "))
	keywords := intentKeywords(intent)
	for i := range items {
		item := &items[i]
		lowerPath := strings.ToLower(item.Path)
		if item.InclusionLevel == "" {
			item.InclusionLevel = InclusionConditional
		}
		if isBroadManifestCandidate(*item) {
			continue
		}
		for _, keyword := range keywords {
			if strings.Contains(lowerPath, keyword) {
				item.RelevanceScore += 10
				item.Reason += " Task keyword matches path."
			}
			if strings.Contains(strings.ToLower(item.Content), keyword) {
				item.RelevanceScore += 4
				item.Reason += " Task keyword appears in file content."
			}
			if pathHasSegment(lowerPath, keyword) {
				item.RelevanceScore += 20
				item.Reason += " Task keyword directly matches a path segment."
			}
		}
		for _, affected := range req.AffectedFiles {
			if item.Path == affected {
				item.RelevanceScore += 12
				item.Reason += " Exact affected file requested."
				break
			}
			if relatedPath(item.Path, affected) {
				item.RelevanceScore += 14
				item.Reason += " File shares a module, name, or test/source relationship with an affected file."
				break
			}
		}
		if strings.Contains(intent, filepath.Base(lowerPath)) {
			item.RelevanceScore += 20
			item.Reason += " Path matches task intent."
		}
		switch normalizeTask(req.TaskType) {
		case normalizeTask(TaskBugFix), "bug_fix":
			if strings.Contains(lowerPath, "troubleshoot") || item.Type != ContextFile || (strings.Contains(lowerPath, "test") && relatedToAnyAffected(item.Path, req.AffectedFiles)) {
				item.RelevanceScore += 20
				item.Reason += " Bug fixes prioritize tests and troubleshooting."
			}
		case normalizeTask(TaskAddFeature), "feature", "new_feature":
			if strings.Contains(lowerPath, "business") || strings.Contains(lowerPath, "service") {
				item.RelevanceScore += 12
				item.Reason += " Feature work prioritizes the affected domain service and business behavior."
			}
		case normalizeTask(TaskRefactor):
			if item.Type == ContextFile || strings.Contains(lowerPath, "architecture") {
				item.RelevanceScore += 14
				item.Reason += " Refactors prioritize implementation and architecture boundaries."
			}
		}
	}
	sort.SliceStable(items, func(i, j int) bool {
		if items[i].RelevanceScore == items[j].RelevanceScore {
			return items[i].Size < items[j].Size
		}
		return items[i].RelevanceScore > items[j].RelevanceScore
	})
	return items
}

func intentKeywords(intent string) []string {
	aliases := map[string][]string{
		"dogfood":    {"app", "context", "quality", "docs", "prompt"},
		"production": {"app", "quality", "docs", "release"},
		"hardening":  {"quality", "security", "errors", "validation"},
		"prompt":     {"prompt", "ai", "context"},
		"context":    {"context", "discovery"},
		"workflow":   {"workflow"},
		"docs":       {"docs", "documentation"},
	}
	seen := map[string]bool{}
	var out []string
	for _, token := range strings.FieldsFunc(strings.ToLower(intent), func(r rune) bool {
		return (r < 'a' || r > 'z') && (r < '0' || r > '9')
	}) {
		if len(token) < 4 {
			continue
		}
		if !seen[token] {
			seen[token] = true
			out = append(out, token)
		}
		for _, alias := range aliases[token] {
			if !seen[alias] {
				seen[alias] = true
				out = append(out, alias)
			}
		}
	}
	return out
}

func (e *Engine) optimize(pkg *ContextPackage, candidates []ContextItem, limit int, threshold float64, keepExcluded bool) {
	seen := map[string]bool{}
	size := 0
	pkg.Summary.InitialCount = len(candidates)
	pkg.Summary.BudgetLimit = limit
	pkg.Summary.Deduplicated = true
	for _, item := range candidates {
		if item.Path != "" && seen[item.Path] {
			continue
		}
		seen[item.Path] = true
		pkg.Summary.InitialSize += item.Size
		if item.InclusionLevel != InclusionRequired && item.RelevanceScore < threshold {
			item.InclusionLevel = InclusionExcluded
			item.ExclusionReason = fmt.Sprintf("relevance score %.1f is below threshold %.1f", item.RelevanceScore, threshold)
			if keepExcluded {
				pkg.ExcludedItems = append(pkg.ExcludedItems, item)
				pkg.Explanations[item.Path] = "EXCLUDED: " + item.ExclusionReason
			}
			continue
		}
		if item.Size > limit/2 && limit != GetBudgetLimit(BudgetUnlimited) {
			item.Summary = summarize(item.Content, 800)
			item.Content = item.Summary
			item.Size = len(item.Content)
			item.Truncated = true
			item.Reason += " Large file summarized to preserve token budget."
		}
		if size+item.Size > limit {
			pkg.Summary.DroppedFiles = append(pkg.Summary.DroppedFiles, item.Path)
			item.InclusionLevel = InclusionExcluded
			item.ExclusionReason = fmt.Sprintf("over budget with score %.1f", item.RelevanceScore)
			if keepExcluded {
				pkg.ExcludedItems = append(pkg.ExcludedItems, item)
				pkg.Explanations[item.Path] = "EXCLUDED: " + item.ExclusionReason
			}
			continue
		}
		pkg.Items = append(pkg.Items, item)
		size += item.Size
		pkg.Explanations[item.Path] = "INCLUDED: " + item.Reason
		pkg.Reasoning = append(pkg.Reasoning, item.Reason)
		e.publish(eventbus.ContextItemSelected, "context item selected", item)
		switch item.Type {
		case ContextDocumentation:
			pkg.SelectedDocs = append(pkg.SelectedDocs, item.Path)
		case ContextStandard:
			pkg.RelevantStandards = append(pkg.RelevantStandards, item.Path)
		case ContextManifestEntry:
			pkg.RelatedPlaybooks = append(pkg.RelatedPlaybooks, item.Path)
		default:
			pkg.SelectedFiles = append(pkg.SelectedFiles, item.Path)
		}
		pkg.Documents = append(pkg.Documents, DocumentItem{Path: item.Path, Category: string(item.Type), Content: item.Content, Size: item.Size, Score: item.RelevanceScore, Explanation: item.Reason})
	}
	pkg.Summary.FinalCount = len(pkg.Items)
	pkg.Summary.FinalSize = size
	pkg.EstimatedTokens = estimateTokens(size)
}

func (e *Engine) renderSystemPrompt(pkg *ContextPackage) string {
	var b strings.Builder
	b.WriteString("Selected Context:\n")
	for _, item := range pkg.Items {
		b.WriteString(item.Path)
		if item.InclusionLevel != "" {
			b.WriteString(" [")
			b.WriteString(string(item.InclusionLevel))
			b.WriteString("]")
		}
		b.WriteString("\n")
		b.WriteString("Reason: ")
		b.WriteString(item.Reason)
		b.WriteString("\n")
	}
	b.WriteString("\n")
	return b.String()
}

func (e *Engine) cacheSalt(req ContextRequest) string {
	parts := []string{string(req.TaskType), req.WorkflowType, req.UserIntent, req.RequestedOperation}
	if req.Project != nil {
		parts = append(parts, req.Project.RootDir, strings.Join(req.Project.Languages, ","), strings.Join(req.Project.Frameworks, ","))
	}
	return strings.Join(parts, "|")
}

func relatedPath(path, affected string) bool {
	pathBase := strings.TrimSuffix(strings.ToLower(filepath.Base(path)), strings.ToLower(filepath.Ext(path)))
	affectedBase := strings.TrimSuffix(strings.ToLower(filepath.Base(affected)), strings.ToLower(filepath.Ext(affected)))
	if pathBase == "" || affectedBase == "" {
		return false
	}
	if pathBase == affectedBase || strings.Contains(pathBase, affectedBase) || strings.Contains(affectedBase, pathBase) {
		return true
	}
	return filepath.Dir(path) == filepath.Dir(affected)
}

func relatedToAnyAffected(path string, affected []string) bool {
	for _, item := range affected {
		if relatedPath(path, item) {
			return true
		}
	}
	return false
}

func pathHasSegment(path, keyword string) bool {
	for _, segment := range strings.FieldsFunc(path, func(r rune) bool {
		return r == '/' || r == '\\' || r == '-' || r == '_' || r == '.'
	}) {
		if segment == keyword {
			return true
		}
	}
	return false
}

func playbookRelevance(req ContextRequest, playbook manifest.PlaybookDefinition) (float64, string) {
	signals := requestSignals(req)
	haystack := strings.ToLower(strings.Join([]string{playbook.ID, playbook.Name, string(playbook.Category), playbook.Location, playbook.Description}, " "))
	signalTerms := []string{
		"security", "auth", "permission", "secret", "token",
		"database", "migration", "model", "query", "schema", "persistence",
		"api", "endpoint", "route", "request", "response", "contract",
		"frontend", "ui", "component", "vue", "react", "css", "screen",
		"deployment", "deploy", "environment", "infra", "docker", "runtime",
		"documentation", "docs", "readme",
		"architecture", "boundary", "dependency", "service", "module",
		"test", "testing", "bug", "failure", "troubleshoot",
		"feature",
	}
	for _, term := range signalTerms {
		if strings.Contains(signals, term) && strings.Contains(haystack, term) {
			return 64, "Manifest standard matches a concrete task signal: " + term + "."
		}
	}
	if normalizeTask(req.TaskType) == normalizeTask(TaskBugFix) && (strings.Contains(haystack, "test") || strings.Contains(haystack, "troubleshoot")) {
		return 64, "Bug-fix workflow requires directly relevant testing or troubleshooting guidance."
	}
	if normalizeTask(req.TaskType) == normalizeTask(TaskRefactor) && strings.Contains(haystack, "architecture") {
		return 64, "Refactor task affects architecture boundaries, so this standard is directly relevant."
	}
	if playbook.Category == manifest.CategoryWorkflows && strings.Contains(haystack, normalizeTask(req.TaskType)) {
		return 60, "Workflow playbook directly names the requested task."
	}
	return 30, "Manifest maps this broad playbook, but no concrete task signal requires it."
}

func isBroadManifestCandidate(item ContextItem) bool {
	return (item.Type == ContextStandard || item.Type == ContextManifestEntry) && strings.HasPrefix(item.Reason, "Manifest maps this broad playbook")
}

func estimateTokens(bytes int) int {
	if bytes <= 0 {
		return 0
	}
	return (bytes + 3) / 4
}

func (e *Engine) publish(t eventbus.EventType, msg string, payload any) {
	if e.events == nil {
		return
	}
	e.events.Publish(eventbus.Event{Type: t, Message: msg, Payload: payload})
}

func summarize(content string, max int) string {
	if len(content) <= max {
		return content
	}
	return content[:max] + "\n\n[truncated summary: large file shortened for context budget]"
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}
