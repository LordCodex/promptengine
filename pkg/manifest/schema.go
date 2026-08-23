package manifest

import "time"

const SupportedSchemaVersion = "2.0.0"

type ProjectMetadata struct {
	Name             string    `json:"name" yaml:"name"`
	Version          string    `json:"version" yaml:"version"`
	SchemaVersion    string    `json:"schema_version" yaml:"schema_version"`
	GeneratedAt      time.Time `json:"generated_at" yaml:"generated_at"`
	GeneratedAtRaw   string    `json:"-" yaml:"-"`
	PromptEnginePath string    `json:"promptengine_path,omitempty" yaml:"promptengine_path,omitempty"`
}

type TechnologyDefinition struct {
	ID                 string   `json:"id" yaml:"id"`
	Language           string   `json:"language" yaml:"language"`
	Framework          string   `json:"framework,omitempty" yaml:"framework,omitempty"`
	Stack              string   `json:"stack,omitempty" yaml:"stack,omitempty"`
	VersionConstraints []string `json:"version_constraints,omitempty" yaml:"version_constraints,omitempty"`
	RelatedStandards   []string `json:"related_standards,omitempty" yaml:"related_standards,omitempty"`
	RelatedPlaybooks   []string `json:"related_playbooks,omitempty" yaml:"related_playbooks,omitempty"`
}

type PlaybookCategory string

const (
	CategoryCore          PlaybookCategory = "core"
	CategoryStacks        PlaybookCategory = "stacks"
	CategorySecurity      PlaybookCategory = "security"
	CategoryPerformance   PlaybookCategory = "performance"
	CategoryDesign        PlaybookCategory = "design"
	CategoryWorkflows     PlaybookCategory = "workflows"
	CategoryProject       PlaybookCategory = "project"
	CategoryBridge        PlaybookCategory = "bridge"
	CategoryChecklist     PlaybookCategory = "checklist"
	CategoryDecisionGuide PlaybookCategory = "decision-guide"
	CategoryGuide         PlaybookCategory = "guide"
	CategoryPrompt        PlaybookCategory = "prompt"
	CategoryAI            PlaybookCategory = "ai"
	CategoryCLI           PlaybookCategory = "cli"
)

type PlaybookDefinition struct {
	ID          string           `json:"id" yaml:"id"`
	Name        string           `json:"name" yaml:"name"`
	Category    PlaybookCategory `json:"category" yaml:"category"`
	Location    string           `json:"location" yaml:"location"`
	Description string           `json:"description,omitempty" yaml:"description,omitempty"`
	Priority    int              `json:"priority" yaml:"priority"`
}

type WorkflowDefinition struct {
	ID                string   `json:"id" yaml:"id"`
	Steps             []string `json:"steps" yaml:"steps"`
	RequiredContext   []string `json:"required_context,omitempty" yaml:"required_context,omitempty"`
	RequiredPlaybooks []string `json:"required_playbooks,omitempty" yaml:"required_playbooks,omitempty"`
	OptionalPlaybooks []string `json:"optional_playbooks,omitempty" yaml:"optional_playbooks,omitempty"`
	Validators        []string `json:"validators,omitempty" yaml:"validators,omitempty"`
	Prompts           []string `json:"prompts,omitempty" yaml:"prompts,omitempty"`
}

type PromptMapping struct {
	TaskType          string   `json:"task_type" yaml:"task_type"`
	PromptTemplate    string   `json:"prompt_template" yaml:"prompt_template"`
	RequiredContext   []string `json:"required_context,omitempty" yaml:"required_context,omitempty"`
	AIFormattingRules []string `json:"ai_formatting_rules,omitempty" yaml:"ai_formatting_rules,omitempty"`
}

type TemplateDefinition struct {
	Name               string   `json:"name" yaml:"name"`
	Location           string   `json:"location" yaml:"location"`
	Version            string   `json:"version" yaml:"version"`
	Variables          []string `json:"variables,omitempty" yaml:"variables,omitempty"`
	SupportedWorkflows []string `json:"supported_workflows,omitempty" yaml:"supported_workflows,omitempty"`
	Type               string   `json:"type,omitempty" yaml:"type,omitempty"`
}

type CommandMapping struct {
	Command  string `json:"command" yaml:"command"`
	Workflow string `json:"workflow" yaml:"workflow"`
}

type TaskRelationship struct {
	TaskType         string   `json:"task_type" yaml:"task_type"`
	RelatedTasks     []string `json:"related_tasks,omitempty" yaml:"related_tasks,omitempty"`
	RequiredWorkflow string   `json:"required_workflow,omitempty" yaml:"required_workflow,omitempty"`
}

type Manifest struct {
	Metadata          ProjectMetadata                `json:"metadata" yaml:"metadata"`
	Technologies      []TechnologyDefinition         `json:"technologies,omitempty" yaml:"technologies,omitempty"`
	Playbooks         []PlaybookDefinition           `json:"playbooks,omitempty" yaml:"playbooks,omitempty"`
	Workflows         []WorkflowDefinition           `json:"workflows,omitempty" yaml:"workflows,omitempty"`
	Prompts           []PromptMapping                `json:"prompts,omitempty" yaml:"prompts,omitempty"`
	Templates         []TemplateDefinition           `json:"templates,omitempty" yaml:"templates,omitempty"`
	CommandMappings   []CommandMapping               `json:"command_mappings,omitempty" yaml:"command_mappings,omitempty"`
	TaskRelationships []TaskRelationship             `json:"task_relationships,omitempty" yaml:"task_relationships,omitempty"`
	PluginData        map[string]map[string]any      `json:"plugin_data,omitempty" yaml:"plugin_data,omitempty"`
	Extensions        map[string][]ExtensionResource `json:"extensions,omitempty" yaml:"extensions,omitempty"`
}

type ExtensionResource struct {
	Kind string         `json:"kind" yaml:"kind"`
	ID   string         `json:"id" yaml:"id"`
	Data map[string]any `json:"data,omitempty" yaml:"data,omitempty"`
}

// Legacy aliases kept so existing packages continue to compile during the phase rollout.
type WorkflowDef = WorkflowDefinition
type StandardDef struct {
	ID         string `json:"id" yaml:"id"`
	Title      string `json:"title" yaml:"title"`
	Category   string `json:"category" yaml:"category"`
	Technology string `json:"technology" yaml:"technology"`
	Priority   int    `json:"priority" yaml:"priority"`
}
type TechDef struct {
	ID               string   `json:"id" yaml:"id"`
	RelatedStandards []string `json:"related_standards" yaml:"related_standards"`
	DocsScaffolds    []string `json:"docs_scaffolds" yaml:"docs_scaffolds"`
}
type PromptDef struct {
	ID           string   `json:"id" yaml:"id"`
	WorkflowID   string   `json:"workflow_id" yaml:"workflow_id"`
	Variables    []string `json:"variables" yaml:"variables"`
	TemplatePath string   `json:"template_path" yaml:"template_path"`
}
type HealthRuleDef struct {
	ID       string `json:"id" yaml:"id"`
	Category string `json:"category" yaml:"category"`
	Weight   int    `json:"weight" yaml:"weight"`
}
type VersionCompatibilityDef struct {
	MinCLIVersion     string `json:"min_cli_version" yaml:"min_cli_version"`
	MaxCLIVersion     string `json:"max_cli_version" yaml:"max_cli_version"`
	ManifestSchemaVer string `json:"manifest_schema_ver" yaml:"manifest_schema_ver"`
}
type DeclarativeManifest struct {
	SchemaVersion string                   `json:"schema_version" yaml:"schema_version"`
	Workflows     map[string]WorkflowDef   `json:"workflows" yaml:"workflows"`
	Standards     map[string]StandardDef   `json:"standards" yaml:"standards"`
	Technologies  map[string]TechDef       `json:"technologies" yaml:"technologies"`
	Prompts       map[string]PromptDef     `json:"prompts" yaml:"prompts"`
	HealthRules   map[string]HealthRuleDef `json:"health_rules" yaml:"health_rules"`
	Compatibility VersionCompatibilityDef  `json:"compatibility" yaml:"compatibility"`
}
