package context

import (
	"context"
	"sort"

	"github.com/LordCodex/promptengine/internal/domain/discovery"
	"github.com/LordCodex/promptengine/internal/filesystem"
)

// Engine compiles context packages for tasks
type Engine struct {
	fs filesystem.FileSystem
}

func NewEngine(fs filesystem.FileSystem) *Engine {
	return &Engine{fs: fs}
}

func (e *Engine) GenerateContext(ctx context.Context, task TaskType, pm *discovery.ProjectModel, budget BudgetType) (*ContextPackage, error) {
	pkg := NewContextPackage(string(task), budget)
	limit := GetBudgetLimit(budget)

	// 1. Gather candidate files
	candidates := e.gatherCandidates(pm)

	// 2. Score candidate files
	scored := e.scoreCandidates(candidates, task, pm)

	// 3. Optimize and slice to fit budget limit
	e.optimizeContext(pkg, scored, limit)

	// 4. Set unified system prompt context headers
	pkg.SystemPrompt = "You are an AI coding assistant. Follow codebase specifications strictly. Use the attached playbooks context:\n"
	for _, doc := range pkg.Documents {
		pkg.SystemPrompt += "\n--- File: " + doc.Path + " ---\n" + doc.Content + "\n"
	}

	return pkg, nil
}

func (e *Engine) gatherCandidates(pm *discovery.ProjectModel) []DocumentItem {
	var list []DocumentItem

	// Add discovered documents
	for label, doc := range pm.Docs {
		if doc.Exists {
			content, err := e.fs.ReadFile(doc.Path)
			if err == nil {
				category := "generic"
				switch label {
				case "BusinessRules":
					category = "business_rules"
				case "Architecture":
					category = "architecture"
				case "Database", "API":
					category = "workflow" // design specs
				case "Roadmap", "Troubleshooting":
					category = "roadmap"
				}

				list = append(list, DocumentItem{
					Path:     doc.Path,
					Category: category,
					Content:  string(content),
					Size:     len(content),
				})
			}
		}
	}

	// Add AGENTS.md if present
	if pm.PromptEngine.AgentsMDPresent {
		content, err := e.fs.ReadFile("AGENTS.md")
		if err == nil {
			list = append(list, DocumentItem{
				Path:     "AGENTS.md",
				Category: "agents",
				Content:  string(content),
				Size:     len(content),
			})
		}
	}

	return list
}

func (e *Engine) scoreCandidates(items []DocumentItem, task TaskType, pm *discovery.ProjectModel) []DocumentItem {
	for i, item := range items {
		score := 50.0 // base score

		// Priority Category weights
		switch item.Category {
		case "agents":
			score += 50.0 // Always include project constitution
			items[i].Explanation = "Mandatory Project Constitution (AGENTS.md)."
		case "business_rules":
			score += 40.0 // Business Rules outrank others
			items[i].Explanation = "Business Rules take highest domain precedence."
		case "architecture":
			score += 30.0 // Architecture outranks roadmap
			items[i].Explanation = "Architecture details carry heavy weight."
		case "workflow":
			score += 25.0
			items[i].Explanation = "Relevant task workflows specification."
		case "stack":
			score += 20.0
			items[i].Explanation = "Technology-specific playbook."
		case "roadmap":
			score += 10.0
			items[i].Explanation = "Roadmap/Progress context guidelines."
		default:
			items[i].Explanation = "Standard project guideline."
		}

		// Task correlation adjustments
		if task == TaskBugFix && item.Path == "docs/Troubleshooting.md" {
			score += 30.0
			items[i].Explanation += " Elevated: Troubleshooting file is critical for resolving bug-fix issues."
		}
		if (task == TaskDatabaseChanges || task == TaskAddFeature) && item.Path == "docs/Database.md" {
			score += 30.0
			items[i].Explanation += " Elevated: Database spec is crucial for data schema migrations."
		}

		items[i].Score = score
	}

	// Sort by score descending
	sort.Slice(items, func(i, j int) bool {
		return items[i].Score > items[j].Score
	})

	return items
}

func (e *Engine) optimizeContext(pkg *ContextPackage, candidates []DocumentItem, limit int) {
	currentSize := 0
	pkg.Summary.InitialCount = len(candidates)
	pkg.Summary.BudgetLimit = limit
	pkg.Summary.Deduplicated = true

	// Deduplication mapping (ensure no duplicate paths)
	seen := make(map[string]bool)

	for _, item := range candidates {
		pkg.Summary.InitialSize += item.Size

		if seen[item.Path] {
			continue
		}
		seen[item.Path] = true

		if currentSize+item.Size <= limit {
			pkg.Documents = append(pkg.Documents, item)
			currentSize += item.Size
			pkg.Explanations[item.Path] = "INCLUDED: " + item.Explanation
		} else {
			pkg.Summary.DroppedFiles = append(pkg.Summary.DroppedFiles, item.Path)
			pkg.Explanations[item.Path] = "EXCLUDED: Exceeded token budget constraints. Priority score: " + string(rune(item.Score))
		}
	}

	pkg.Summary.FinalCount = len(pkg.Documents)
	pkg.Summary.FinalSize = currentSize
}
