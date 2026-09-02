package rulesources

import (
	"context"
	"testing"
	"time"

	"github.com/LordCodex/promptengine/internal/filesystem"
)

type fakeSourceFetcher struct {
	files map[string]map[string][]byte
}

func (f fakeSourceFetcher) Fetch(ctx context.Context, source Source) (map[string][]byte, error) {
	return f.files[source.Repository], nil
}

func TestSynchronizerSyncProfileCachesPinnedSources(t *testing.T) {
	registry := &Registry{
		Version: 1,
		Sources: map[string]Source{
			"universal": {Repository: "LordCodex/engineering-ai-rules", Ref: "abc123", Inherits: nil},
			"vue":       {Repository: "LordCodex/vue-ai-rules", Ref: "def456", Inherits: []string{"universal"}},
		},
	}
	if err := registry.Validate(); err != nil {
		t.Fatal(err)
	}
	profile := &Profile{ID: "vue", Version: 1, Inheritance: []string{"vue"}}
	fsys := filesystem.NewMockFileSystem()
	fetcher := fakeSourceFetcher{files: map[string]map[string][]byte{
		"LordCodex/engineering-ai-rules": {"RULES.md": []byte("universal"), "docs/SECURITY.md": []byte("security")},
		"LordCodex/vue-ai-rules":         {"RULES.md": []byte("vue")},
	}}
	syncer := NewSynchronizer(registry, fsys, fetcher)
	syncer.Now = func() time.Time { return time.Date(2026, 9, 2, 7, 0, 0, 0, time.UTC) }

	report, err := syncer.SyncProfile(context.Background(), profile)
	if err != nil {
		t.Fatal(err)
	}
	if len(report.Sources) != 2 {
		t.Fatalf("expected two resolved sources, got %#v", report.Sources)
	}
	if !fsys.Exists(".promptengine/rules/universal/abc123/RULES.md") {
		t.Fatal("expected universal rules to be cached")
	}
	if !fsys.Exists(".promptengine/rules/vue/def456/RULES.md") {
		t.Fatal("expected Vue rules to be cached")
	}
	snapshot, err := LoadSnapshot(fsys, DefaultCacheRoot, "vue", registry.Sources["vue"])
	if err != nil {
		t.Fatal(err)
	}
	if snapshot.Ref != "def456" || len(snapshot.Files) != 1 {
		t.Fatalf("unexpected snapshot: %#v", snapshot)
	}
}

func TestMatchProfilePrefersMostSpecificStack(t *testing.T) {
	profiles := []Profile{
		{ID: "vue", Version: 1, Match: ProfileMatch{RequiredTechnologies: []string{"vue"}}, Inheritance: []string{"vue"}},
		{ID: "laravel", Version: 1, Match: ProfileMatch{RequiredTechnologies: []string{"php", "laravel"}}, Inheritance: []string{"laravel"}},
		{ID: "laravel-inertia-vue", Version: 1, Match: ProfileMatch{RequiredTechnologies: []string{"php", "laravel", "inertia", "vue"}}, Inheritance: []string{"laravel", "vue"}},
	}
	profile, ok := MatchProfile(profiles, []string{"PHP", "Laravel", "Inertia", "Vue", "JavaScript"})
	if !ok {
		t.Fatal("expected a profile match")
	}
	if profile.ID != "laravel-inertia-vue" {
		t.Fatalf("expected most specific profile, got %q", profile.ID)
	}
}

func TestArchiveRelativePathRejectsTraversal(t *testing.T) {
	for _, candidate := range []string{"../RULES.md", "/absolute/RULES.md", "root/../../RULES.md"} {
		if rel, ok := archiveRelativePath(candidate); ok {
			t.Fatalf("expected %q to be rejected, got %q", candidate, rel)
		}
	}
	if rel, ok := archiveRelativePath("repo-sha/docs/SECURITY.md"); !ok || rel != "docs/SECURITY.md" {
		t.Fatalf("expected safe archive path, got %q, %v", rel, ok)
	}
}
