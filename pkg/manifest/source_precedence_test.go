package manifest

import (
	"testing"

	"github.com/LordCodex/promptengine/internal/filesystem"
)

func TestProjectSourceOverridesOrganizationDefinitionWithSameID(t *testing.T) {
	fsys := filesystem.NewMockFileSystem()
	fsys.WriteFile("org-rule.md", []byte("organization"), 0644)
	fsys.WriteFile("project-rule.md", []byte("project"), 0644)

	organization := &Manifest{
		Metadata: ProjectMetadata{Name: "organization", Version: "1", SchemaVersion: SupportedSchemaVersion},
		Playbooks: []PlaybookDefinition{{
			ID:       "shared-security-rule",
			Name:     "Organization rule",
			Category: CategorySecurity,
			Location: "org-rule.md",
		}},
	}
	project := &Manifest{
		Metadata: ProjectMetadata{Name: "project", Version: "1", SchemaVersion: SupportedSchemaVersion},
		Playbooks: []PlaybookDefinition{{
			ID:       "shared-security-rule",
			Name:     "Project specialization",
			Category: CategorySecurity,
			Location: "project-rule.md",
		}},
	}

	engine := NewEngineWithFS(fsys)
	if err := engine.Register("authoritative", SourceOrganization, organization); err != nil {
		t.Fatal(err)
	}
	if err := engine.Register("project", SourceProject, project); err != nil {
		t.Fatal(err)
	}

	active := engine.ActiveManifest()
	if active.Metadata.Name != "project" {
		t.Fatalf("expected project metadata to have highest source precedence, got %q", active.Metadata.Name)
	}
	if len(active.Playbooks) != 1 {
		t.Fatalf("expected same-id playbook to be specialized instead of duplicated, got %#v", active.Playbooks)
	}
	if active.Playbooks[0].Location != "project-rule.md" || active.Playbooks[0].Name != "Project specialization" {
		t.Fatalf("expected project playbook to override organization definition, got %#v", active.Playbooks[0])
	}
}
