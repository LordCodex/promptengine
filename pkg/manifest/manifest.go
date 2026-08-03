package manifest

import (
	"encoding/json"
	"os"
)

type Playbook struct {
	ID   string `json:"id"`
	Path string `json:"path"`
}

type TaskMapping struct {
	RequiredPlaybookIDs []string `json:"required_playbook_ids"`
	OptionalPlaybookIDs []string `json:"optional_playbook_ids"`
}

// PlaybookManifest matches the structure of playbook-manifest.json
type PlaybookManifest struct {
	RepositoryNavigation map[string]string      `json:"repository_navigation"`
	CorePlaybooks        []Playbook             `json:"core_playbooks"`
	TechnologyStacks     map[string][]Playbook  `json:"technology_stacks"`
	DomainPlaybooks      []Playbook             `json:"domain_playbooks"`
	ProjectPlaybooks     []Playbook             `json:"project_playbooks"`
	BridgePlaybooks      []Playbook             `json:"bridge_playbooks"`
	Checklists           []Playbook             `json:"checklists"`
	Workflows            []Playbook             `json:"workflows"`
	DecisionGuides       []Playbook             `json:"decision_guides"`
	AIBootstrap          []Playbook             `json:"ai_bootstrap"`
	Guides               []Playbook             `json:"guides"`
	PromptsLibrary       []Playbook             `json:"prompts_library"`
	CliFoundation        []Playbook             `json:"cli_foundation"`
	CliCommandSpecs      []Playbook             `json:"cli_command_specifications"`
	TaskMappings         map[string]TaskMapping `json:"task_mappings"`
}

// Load reads and parses manifest file from disk path
func Load(path string) (*PlaybookManifest, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var m PlaybookManifest
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, err
	}

	return &m, nil
}
