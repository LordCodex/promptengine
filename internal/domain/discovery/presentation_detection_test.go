package discovery

import (
	"testing"

	"github.com/LordCodex/promptengine/internal/filesystem"
)

func TestDiscoveryPipeline_DetectsLaravelLivewireFromComposer(t *testing.T) {
	fs := filesystem.NewMockFileSystem()
	fs.WriteFile("composer.json", []byte(`{"require":{"php":"^8.3","laravel/framework":"^13.0","livewire/livewire":"^4.0"}}`), 0644)
	fs.WriteFile("artisan", []byte("cli"), 0644)

	pm := runDiscovery(t, fs)
	assertContains(t, pm.Frameworks, "Laravel")
	assertContains(t, pm.Frameworks, "Livewire")
}

func TestDiscoveryPipeline_DetectsServerSideInertiaFromComposer(t *testing.T) {
	fs := filesystem.NewMockFileSystem()
	fs.WriteFile("composer.json", []byte(`{"require":{"php":"^8.3","laravel/framework":"^13.0","inertiajs/inertia-laravel":"^3.0"}}`), 0644)
	fs.WriteFile("artisan", []byte("cli"), 0644)

	pm := runDiscovery(t, fs)
	assertContains(t, pm.Frameworks, "Laravel")
	assertContains(t, pm.Frameworks, "Inertia")
}
