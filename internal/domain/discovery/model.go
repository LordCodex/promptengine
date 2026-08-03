package discovery

import "time"

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

type ProjectInfo struct {
	Name         string            `json:"name"`
	RootPath     string            `json:"root_path"`
	DetectedType string            `json:"detected_type"`
	Metadata     map[string]string `json:"metadata,omitempty"`
	GeneratedAt  time.Time         `json:"generated_at"`
}

type TechnologyInfo struct {
	Languages       []string `json:"languages"`
	Frameworks      []string `json:"frameworks"`
	Runtimes        []string `json:"runtimes"`
	PackageManagers []string `json:"package_managers"`
	Databases       []string `json:"databases"`
	Infrastructure  []string `json:"infrastructure"`
	Testing         []string `json:"testing"`
}

type RepositoryInfo struct {
	RootPath           string   `json:"root_path"`
	Directories        []string `json:"directories"`
	Files              []string `json:"files"`
	ConfigurationFiles []string `json:"configuration_files"`
	DocumentationFiles []string `json:"documentation_files"`
	IgnoredFiles       []string `json:"ignored_files"`
	PermissionErrors   []string `json:"permission_errors"`
	IsMonorepo         bool     `json:"is_monorepo"`
}

type ArchitectureInfo struct {
	Backend        bool     `json:"backend"`
	Frontend       bool     `json:"frontend"`
	Mobile         bool     `json:"mobile"`
	Services       []string `json:"services"`
	Infrastructure bool     `json:"infrastructure"`
}

type DocSpec struct {
	Name         string  `json:"name"`
	Exists       bool    `json:"exists"`
	Path         string  `json:"path"`
	Completeness float64 `json:"completeness"`
}

type ArchitectureInference struct {
	Style      string  `json:"style"`
	Confidence float64 `json:"confidence"`
	Reason     string  `json:"reason"`
}

type PromptEngineStatus struct {
	Installed       bool   `json:"installed"`
	AgentsMDPresent bool   `json:"agents_md_present"`
	HasConfig       bool   `json:"has_config"`
	Version         string `json:"version"`
	ConfigVersion   string `json:"config_version"`
}

type ProjectModel struct {
	Project      ProjectInfo      `json:"project"`
	Technology   TechnologyInfo   `json:"technology"`
	Repository   RepositoryInfo   `json:"repository"`
	Architecture ArchitectureInfo `json:"architecture"`

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
	if rootDir == "" {
		rootDir = "."
	}
	return &ProjectModel{
		Project: ProjectInfo{
			RootPath:    rootDir,
			Metadata:    map[string]string{},
			GeneratedAt: time.Now().UTC(),
		},
		Repository: RepositoryInfo{RootPath: rootDir},
		RootDir:    rootDir,
		Docs:       map[string]DocSpec{},
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

func (pm *ProjectModel) SyncLegacyFields() {
	pm.Technology.Languages = uniqueStrings(append(pm.Technology.Languages, pm.Languages...))
	pm.Technology.Frameworks = uniqueStrings(append(pm.Technology.Frameworks, pm.Frameworks...))
	pm.Technology.Databases = uniqueStrings(append(pm.Technology.Databases, pm.Databases...))
	pm.Technology.PackageManagers = uniqueStrings(append(pm.Technology.PackageManagers, pm.PackageManagers...))
	pm.Technology.Testing = uniqueStrings(append(pm.Technology.Testing, pm.TestingFrames...))
	pm.Technology.Infrastructure = uniqueStrings(append(pm.Technology.Infrastructure, pm.CIs...))
	if pm.HasDocker {
		pm.Technology.Infrastructure = addUnique(pm.Technology.Infrastructure, "Docker")
	}

	pm.Languages = uniqueStrings(pm.Technology.Languages)
	pm.Frameworks = uniqueStrings(pm.Technology.Frameworks)
	pm.Databases = uniqueStrings(pm.Technology.Databases)
	pm.PackageManagers = uniqueStrings(pm.Technology.PackageManagers)
	pm.TestingFrames = uniqueStrings(pm.Technology.Testing)
	pm.CIs = uniqueStrings(pm.Technology.Infrastructure)
}

func addUnique(items []string, item string) []string {
	if item == "" {
		return items
	}
	for _, existing := range items {
		if existing == item {
			return items
		}
	}
	return append(items, item)
}

func uniqueStrings(items []string) []string {
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
