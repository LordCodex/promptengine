package promptengine

import "embed"

// StandardsFS contains the versioned PromptEngine engineering knowledge library.
//
// The library is embedded in the CLI so installed binaries can apply PromptEngine
// standards to any project without requiring a separate clone of this repository.
// Project source files and project-specific documentation are never sourced from
// this filesystem; they remain in the developer's workspace.
//
//go:embed playbook-manifest.json core stacks security performance design project bridges checklists workflows decision-guides ai architecture environment legacy guides prompts cli
var StandardsFS embed.FS
