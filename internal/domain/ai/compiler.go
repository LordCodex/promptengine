package ai

import (
	"fmt"
	"strings"

	ctxengine "github.com/LordCodex/promptengine/internal/domain/context"
	"github.com/LordCodex/promptengine/pkg/manifest"
)

type CompileInput struct {
	Provider             string
	Model                string
	UserRequest          string
	SystemInstructions   string
	WorkflowRequirements []string
	Standards            []string
	ContextPackage       *ctxengine.ContextPackage
	Manifest             *manifest.Manifest
	Temperature          float32
	MaxTokens            int
	Metadata             map[string]string
}

type PromptCompiler struct{}

func NewPromptCompiler() *PromptCompiler { return &PromptCompiler{} }

func (c *PromptCompiler) Compile(input CompileInput) (Request, error) {
	if strings.TrimSpace(input.UserRequest) == "" {
		return Request{}, AIError{Category: ErrorInvalidResponse, Message: "user request is required", RecommendedAction: "Provide a task request before compiling a prompt."}
	}
	var contextBuilder strings.Builder
	if input.ContextPackage != nil {
		contextBuilder.WriteString("Selected context:\n")
		for _, item := range input.ContextPackage.Items {
			contextBuilder.WriteString(fmt.Sprintf("\n--- %s: %s ---\n", item.Type, item.Path))
			if item.Summary != "" {
				contextBuilder.WriteString(item.Summary)
			} else {
				contextBuilder.WriteString(item.Content)
			}
			contextBuilder.WriteString("\n")
		}
	}
	if len(input.WorkflowRequirements) > 0 {
		contextBuilder.WriteString("\nWorkflow requirements:\n- ")
		contextBuilder.WriteString(strings.Join(input.WorkflowRequirements, "\n- "))
	}
	if len(input.Standards) > 0 {
		contextBuilder.WriteString("\nStandards:\n- ")
		contextBuilder.WriteString(strings.Join(input.Standards, "\n- "))
	}
	if input.Manifest != nil {
		contextBuilder.WriteString("\nManifest project: ")
		contextBuilder.WriteString(input.Manifest.Metadata.Name)
	}
	return Request{
		Provider:           input.Provider,
		Model:              input.Model,
		Prompt:             input.UserRequest,
		Context:            contextBuilder.String(),
		SystemInstructions: input.SystemInstructions,
		Temperature:        input.Temperature,
		MaxTokens:          input.MaxTokens,
		Metadata:           input.Metadata,
	}, nil
}
