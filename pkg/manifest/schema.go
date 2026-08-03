package manifest

// WorkflowDef defines workflow pipelines in the manifest schema
type WorkflowDef struct {
	ID             string   `json:"id"`
	Name           string   `json:"name"`
	Description    string   `json:"description"`
	Preconditions  []string `json:"preconditions"`
	Postconditions []string `json:"postconditions"`
}

// StandardDef defines codebase guidelines and conventions
type StandardDef struct {
	ID         string `json:"id"`
	Title      string `json:"title"`
	Category   string `json:"category"`
	Technology string `json:"technology"`
	Priority   int    `json:"priority"` // weight score
}

// TechDef defines technology profile templates
type TechDef struct {
	ID               string   `json:"id"`
	RelatedStandards []string `json:"related_standards"`
	DocsScaffolds    []string `json:"docs_scaffolds"`
}

// PromptDef defines variable-injected AI prompts metadata
type PromptDef struct {
	ID          string   `json:"id"`
	WorkflowID  string   `json:"workflow_id"`
	Variables   []string `json:"variables"`
	TemplatePath string   `json:"template_path"`
}

// HealthRuleDef defines metrics variables used by the Health scoring engine
type HealthRuleDef struct {
	ID       string `json:"id"`
	Category string `json:"category"`
	Weight   int    `json:"weight"`
}

// VersionCompatibilityDef outlines CLI/schema target alignments
type VersionCompatibilityDef struct {
	MinCLIVersion     string `json:"min_cli_version"`
	MaxCLIVersion     string `json:"max_cli_version"`
	ManifestSchemaVer string `json:"manifest_schema_ver"`
}

// DeclarativeManifest is the versioned schema representable on disk
type DeclarativeManifest struct {
	SchemaVersion string                  `json:"schema_version"`
	Workflows     map[string]WorkflowDef  `json:"workflows"`
	Standards     map[string]StandardDef  `json:"standards"`
	Technologies  map[string]TechDef      `json:"technologies"`
	Prompts       map[string]PromptDef    `json:"prompts"`
	HealthRules   map[string]HealthRuleDef `json:"health_rules"`
	Compatibility VersionCompatibilityDef `json:"compatibility"`
}
