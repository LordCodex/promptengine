package container

import (
	"log/slog"
	"os"

	"github.com/LordCodex/promptengine/internal/cache"
	"github.com/LordCodex/promptengine/internal/config"
	"github.com/LordCodex/promptengine/internal/eventbus"
	"github.com/LordCodex/promptengine/internal/filesystem"
	"github.com/LordCodex/promptengine/internal/output"
	"github.com/LordCodex/promptengine/internal/telemetry"
)

// Container manages dependency injection across CLI services
type Container struct {
	FS        filesystem.FileSystem
	Logger    *slog.Logger
	Config    *config.AppConfig
	Cache     *cache.Cache
	EventBus  *eventbus.EventBus
	Telemetry *telemetry.Telemetry
	Renderer  output.Renderer
}

// NewContainer initializes and registers core platform dependencies
func NewContainer(verbose bool) (*Container, error) {
	fs := &filesystem.OSFileSystem{}

	// Setup logging
	level := slog.LevelInfo
	if verbose {
		level = slog.LevelDebug
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: level}))

	// Load configuration
	loader := config.NewConfigLoader()
	cfg, err := loader.Load(".promptengine.json", "~/.promptengine/config.json")
	if err != nil {
		return nil, err
	}

	// Caching and Telemetry
	c := cache.NewCache(".")
	tel := telemetry.NewTelemetry()
	eb := eventbus.NewEventBus()

	renderer := output.NewConfiguredRenderer(output.FormatHuman, false, verbose)

	return &Container{
		FS:        fs,
		Logger:    logger,
		Config:    cfg,
		Cache:     c,
		EventBus:  eb,
		Telemetry: tel,
		Renderer:  renderer,
	}, nil
}
