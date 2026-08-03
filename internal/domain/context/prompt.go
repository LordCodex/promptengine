package context

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Provider string

const (
	ProviderChatGPT  Provider = "chatgpt"
	ProviderClaude   Provider = "claude"
	ProviderGemini   Provider = "gemini"
	ProviderCursor   Provider = "cursor"
	ProviderWindsurf Provider = "windsurf"
)

// Formatter formats context packages for target developers/IDE platforms
type Formatter struct {
	provider Provider
}

func NewFormatter(p Provider) *Formatter {
	return &Formatter{provider: p}
}

func (f *Formatter) Format(pkg *ContextPackage) (string, error) {
	switch f.provider {
	case ProviderCursor:
		// Generate JSON payload representing .cursorrules keys
		rules := map[string]interface{}{
			"rules":         pkg.SystemPrompt,
			"task_type":     pkg.TaskType,
			"documentation": f.extractPaths(pkg),
		}
		data, err := json.MarshalIndent(rules, "", "  ")
		if err != nil {
			return "", err
		}
		return string(data), nil

	case ProviderWindsurf:
		// Generate Windsurf .windsurfrules instructions envelope
		var sb strings.Builder
		sb.WriteString("# Windsurf Instructions\n\n")
		sb.WriteString(pkg.SystemPrompt)
		sb.WriteString("\n\n## Monitored Files\n")
		for _, doc := range pkg.Documents {
			sb.WriteString(fmt.Sprintf("- %s\n", doc.Path))
		}
		return sb.String(), nil

	case ProviderClaude:
		// Generate XML tag blocks representing Claude instructions formats
		var sb strings.Builder
		sb.WriteString("<system_instructions>\n")
		sb.WriteString(pkg.SystemPrompt)
		sb.WriteString("</system_instructions>\n")
		return sb.String(), nil

	default:
		// Standard output formatting for ChatGPT / Gemini
		return pkg.SystemPrompt, nil
	}
}

func (f *Formatter) extractPaths(pkg *ContextPackage) []string {
	var paths []string
	for _, doc := range pkg.Documents {
		paths = append(paths, doc.Path)
	}
	return paths
}
