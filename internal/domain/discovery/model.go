package discovery

// ProjectClassification defines project archetypes
type ProjectClassification string

const (
	ClassGreenfield         ProjectClassification = "greenfield"
	ClassExisting           ProjectClassification = "existing"
	ClassPromptEngine       ProjectClassification = "promptengine_project"
	ClassNonPromptEngine    ProjectClassification = "non_promptengine_project"
	ClassMonorepo           ProjectClassification = "monorepo"
	ClassMultiService       ProjectClassification = "multi_service"
	ClassLibrary            ProjectClassification = "library"
	ClassPackage            ProjectClassification = "package"
	ClassCLIApplication     ProjectClassification = "cli_application"
	ClassBackendAPI         ProjectClassification = "backend_api"
	ClassFrontendSPA        ProjectClassification = "frontend_spa"
	ClassSSRApplication     ProjectClassification = "ssr_application"
	ClassMobileApplication  ProjectClassification = "mobile_application"
	ClassDesktopApplication ProjectClassification = "desktop_application"
	ClassHybrid             ProjectClassification = "hybrid"
)

// DocSpec defines localized target specifications
type DocSpec struct {
	Name       string `json:"name"`
	Exists     bool   `json:"exists"`
	Path       string `json:"path"`
	Completeness float64 `json:"completeness"` // percentage 0.0 - 100.0
}

// ArchitectureInference defines architectural style scores
type ArchitectureInference struct {
	Style      string  `json:"style"`      // e.g. "MVC", "Clean Architecture", "Hexagonal"
	Confidence float64 `json:"confidence"` // score 0.0 - 1.0
	Reason     string  `json:"reason"`
}

// PromptEngineStatus outlines PromptEngine installation properties
type PromptEngineStatus struct {
	Installed       bool   `json:"installed"`
	AgentsMDPresent bool   `json:"agents_md_present"`
	HasConfig       bool   `json:"has_config"`
	Version         string `json:"version"`
	ConfigVersion   string `json:"config_version"`
}

// ProjectModel represents everything discovered about a repository
type ProjectModel struct {
	RootDir         string                  `json:"root_dir"`
	Classifications []ProjectClassification `json:"classifications"`
	Languages       []string                `json:"languages"`
	Frameworks      []string                `json:"frameworks"`
	Databases       []string                `json:"databases"`
	PackageManagers []string                `json:"package_managers"`
	TestingFrames   []string                `json:"testing_frameworks"`
	CIs             []string                `json:"ci_cd"`
	HasGit          bool                    `json:"has_git"`
	HasDocker       bool                    `json:"has_docker"`
	Architectures   []ArchitectureInference `json:"architectures"`
	Docs            map[string]DocSpec      `json:"documentation"`
	PromptEngine    PromptEngineStatus      `json:"promptengine"`
}

func NewProjectModel(rootDir string) *ProjectModel {
	return &ProjectModel{
		RootDir:         rootDir,
		Classifications: make([]ProjectClassification, 0),
		Languages:       make([]string, 0),
		Frameworks:      make([]string, 0),
		Databases:       make([]string, 0),
		PackageManagers: make([]string, 0),
		TestingFrames:   make([]string, 0),
		CIs:             make([]string, 0),
		Architectures:   make([]ArchitectureInference, 0),
		Docs:            make(map[string]DocSpec),
	}
}

func (pm *ProjectModel) HasClassification(c ProjectClassification) bool {
	for _, item := range pm.Classifications {
		if item == c {
			return true
		}
	}
	return false
}
