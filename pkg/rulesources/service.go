package rulesources

import (
	"context"
	"fmt"
	"io/fs"

	"github.com/LordCodex/promptengine/internal/filesystem"
	"github.com/LordCodex/promptengine/pkg/manifest"
)

const DefaultRegistryPath = "sources/rules-sources.yaml"

type Service struct {
	Registry     *Registry
	Profiles     []Profile
	Resolver     *Resolver
	Synchronizer *Synchronizer
}

func NewService(library fs.FS, projectFS filesystem.FileSystem) (*Service, error) {
	registry, err := LoadRegistryFS(library, DefaultRegistryPath)
	if err != nil {
		return nil, err
	}
	profiles, err := LoadProfilesFS(library, DefaultProfilesDir)
	if err != nil {
		return nil, err
	}
	resolver := NewResolver(registry, profiles, projectFS)
	return &Service{
		Registry:     registry,
		Profiles:     profiles,
		Resolver:     resolver,
		Synchronizer: NewSynchronizer(registry, projectFS, nil),
	}, nil
}

func (s *Service) Match(technologies []string) (*Profile, bool) {
	if s == nil {
		return nil, false
	}
	return MatchProfile(s.Profiles, technologies)
}

func (s *Service) Resolve(technologies []string, intent string) (*Resolution, error) {
	if s == nil || s.Resolver == nil {
		return nil, fmt.Errorf("rule source service is not initialized")
	}
	return s.Resolver.Resolve(technologies, intent)
}

func (s *Service) BuildManifest(technologies []string) (*manifest.Manifest, *Resolution, error) {
	if s == nil || s.Resolver == nil {
		return nil, nil, fmt.Errorf("rule source service is not initialized")
	}
	return s.Resolver.BuildManifest(technologies)
}

func (s *Service) Sync(ctx context.Context, technologies []string) (*SyncReport, error) {
	if s == nil || s.Synchronizer == nil {
		return nil, fmt.Errorf("rule source service is not initialized")
	}
	profile, ok := s.Match(technologies)
	if !ok {
		return nil, fmt.Errorf("no authoritative rule profile matches detected technologies")
	}
	return s.Synchronizer.SyncProfile(ctx, profile)
}
