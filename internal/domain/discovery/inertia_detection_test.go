package discovery

import (
	"context"
	"testing"

	"github.com/LordCodex/promptengine/internal/filesystem"
)

func TestDiscoveryPipeline_LaravelInertiaVue(t *testing.T) {
	fs := filesystem.NewMockFileSystem()
	fs.WriteFile("composer.json", []byte(`{"require":{"php":"^8.2","laravel/framework":"^13.0","inertiajs/inertia-laravel":"^3.0"}}`), 0644)
	fs.WriteFile("package.json", []byte(`{"dependencies":{"vue":"^3.5","@inertiajs/vue3":"^3.0"},"devDependencies":{"typescript":"^5.9"}}`), 0644)
	fs.WriteFile("artisan", []byte("cli entry"), 0644)

	pipeline := NewDefaultPipeline(nil, nil)
	pm, err := pipeline.Execute(context.Background(), fs, ".")
	if err != nil {
		t.Fatalf("expected discovery to run, got %v", err)
	}

	assertContains(t, pm.Languages, "PHP")
	assertContains(t, pm.Languages, "JavaScript")
	assertContains(t, pm.Languages, "TypeScript")
	assertContains(t, pm.Frameworks, "Laravel")
	assertContains(t, pm.Frameworks, "Inertia")
	assertContains(t, pm.Frameworks, "Vue")
}
