package rulesources

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"sort"
	"strings"
	"time"

	"github.com/LordCodex/promptengine/internal/filesystem"
	"gopkg.in/yaml.v3"
)

const (
	DefaultCacheRoot       = ".promptengine/rules"
	defaultArchiveMaxBytes = 32 << 20
	defaultRuleMaxBytes    = 2 << 20
)

type SourceFetcher interface {
	Fetch(ctx context.Context, source Source) (map[string][]byte, error)
}

type GitHubArchiveFetcher struct {
	Client          *http.Client
	Token           string
	ArchiveMaxBytes int64
	RuleMaxBytes    int64
}

func NewGitHubArchiveFetcher() *GitHubArchiveFetcher {
	token := strings.TrimSpace(os.Getenv("GITHUB_TOKEN"))
	if token == "" {
		token = strings.TrimSpace(os.Getenv("GH_TOKEN"))
	}
	return &GitHubArchiveFetcher{
		Client:          &http.Client{Timeout: 45 * time.Second},
		Token:           token,
		ArchiveMaxBytes: defaultArchiveMaxBytes,
		RuleMaxBytes:    defaultRuleMaxBytes,
	}
}

func (f *GitHubArchiveFetcher) Fetch(ctx context.Context, source Source) (map[string][]byte, error) {
	if err := validateRepository(source.Repository); err != nil {
		return nil, err
	}
	if err := validateCacheComponent(source.Ref); err != nil {
		return nil, fmt.Errorf("invalid source ref: %w", err)
	}

	client := f.Client
	if client == nil {
		client = &http.Client{Timeout: 45 * time.Second}
	}
	archiveLimit := f.ArchiveMaxBytes
	if archiveLimit <= 0 {
		archiveLimit = defaultArchiveMaxBytes
	}
	ruleLimit := f.RuleMaxBytes
	if ruleLimit <= 0 {
		ruleLimit = defaultRuleMaxBytes
	}

	url := fmt.Sprintf("https://api.github.com/repos/%s/tarball/%s", source.Repository, source.Ref)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("build GitHub archive request: %w", err)
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	req.Header.Set("User-Agent", "PromptEngine-rule-sync")
	if strings.TrimSpace(f.Token) != "" {
		req.Header.Set("Authorization", "Bearer "+strings.TrimSpace(f.Token))
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("download %s@%s: %w", source.Repository, source.Ref, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		if resp.StatusCode == http.StatusNotFound && strings.TrimSpace(f.Token) == "" {
			return nil, fmt.Errorf("download %s@%s: GitHub returned 404; the repository may be private, set GITHUB_TOKEN or GH_TOKEN with read access", source.Repository, source.Ref)
		}
		return nil, fmt.Errorf("download %s@%s: GitHub returned %s", source.Repository, source.Ref, resp.Status)
	}

	limited := io.LimitReader(resp.Body, archiveLimit+1)
	gz, err := gzip.NewReader(limited)
	if err != nil {
		return nil, fmt.Errorf("open GitHub archive for %s: %w", source.Repository, err)
	}
	defer gz.Close()

	tr := tar.NewReader(gz)
	files := map[string][]byte{}
	var total int64
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("read GitHub archive for %s: %w", source.Repository, err)
		}
		if header.Typeflag != tar.TypeReg && header.Typeflag != tar.TypeRegA {
			continue
		}
		rel, ok := archiveRelativePath(header.Name)
		if !ok || !isRuleTextFile(rel) {
			continue
		}
		if header.Size < 0 || header.Size > ruleLimit {
			return nil, fmt.Errorf("rule file %q in %s exceeds the %d-byte limit", rel, source.Repository, ruleLimit)
		}
		data, err := io.ReadAll(io.LimitReader(tr, ruleLimit+1))
		if err != nil {
			return nil, fmt.Errorf("read rule file %q from %s: %w", rel, source.Repository, err)
		}
		if int64(len(data)) > ruleLimit {
			return nil, fmt.Errorf("rule file %q in %s exceeds the %d-byte limit", rel, source.Repository, ruleLimit)
		}
		total += int64(len(data))
		if total > archiveLimit {
			return nil, fmt.Errorf("rule content from %s exceeds the %d-byte aggregate limit", source.Repository, archiveLimit)
		}
		files[rel] = data
	}
	if len(files) == 0 {
		return nil, fmt.Errorf("no rule text files found in %s@%s", source.Repository, source.Ref)
	}
	return files, nil
}

type Snapshot struct {
	SourceID   string   `yaml:"source_id"`
	Repository string   `yaml:"repository"`
	Ref        string   `yaml:"ref"`
	SyncedAt   string   `yaml:"synced_at"`
	Files      []string `yaml:"files"`
}

type SourceSyncResult struct {
	SourceID   string `json:"source_id" yaml:"source_id"`
	Repository string `json:"repository" yaml:"repository"`
	Ref        string `json:"ref" yaml:"ref"`
	Files      int    `json:"files" yaml:"files"`
	Bytes      int    `json:"bytes" yaml:"bytes"`
	CachePath  string `json:"cache_path" yaml:"cache_path"`
}

type SyncReport struct {
	Profile string             `json:"profile" yaml:"profile"`
	Sources []SourceSyncResult `json:"sources" yaml:"sources"`
}

type Synchronizer struct {
	Registry  *Registry
	FS        filesystem.FileSystem
	Fetcher   SourceFetcher
	CacheRoot string
	Now       func() time.Time
}

func NewSynchronizer(registry *Registry, fsys filesystem.FileSystem, fetcher SourceFetcher) *Synchronizer {
	if fetcher == nil {
		fetcher = NewGitHubArchiveFetcher()
	}
	return &Synchronizer{
		Registry:  registry,
		FS:        fsys,
		Fetcher:   fetcher,
		CacheRoot: DefaultCacheRoot,
		Now:       time.Now,
	}
}

func (s *Synchronizer) SyncProfile(ctx context.Context, profile *Profile) (*SyncReport, error) {
	if s == nil || s.Registry == nil {
		return nil, fmt.Errorf("rule synchronizer has no registry")
	}
	if s.FS == nil {
		return nil, fmt.Errorf("rule synchronizer has no filesystem")
	}
	if s.Fetcher == nil {
		return nil, fmt.Errorf("rule synchronizer has no fetcher")
	}
	ids, err := s.Registry.ResolveProfile(profile)
	if err != nil {
		return nil, err
	}
	report := &SyncReport{Profile: profile.ID}
	for _, id := range ids {
		result, err := s.SyncSource(ctx, id)
		if err != nil {
			return nil, fmt.Errorf("sync rule source %q: %w", id, err)
		}
		report.Sources = append(report.Sources, result)
	}
	return report, nil
}

func (s *Synchronizer) SyncSource(ctx context.Context, sourceID string) (SourceSyncResult, error) {
	if s == nil || s.Registry == nil {
		return SourceSyncResult{}, fmt.Errorf("rule synchronizer has no registry")
	}
	source, ok := s.Registry.Sources[sourceID]
	if !ok {
		return SourceSyncResult{}, fmt.Errorf("unknown rule source %q", sourceID)
	}
	if err := validateCacheComponent(sourceID); err != nil {
		return SourceSyncResult{}, fmt.Errorf("invalid source id %q: %w", sourceID, err)
	}
	files, err := s.Fetcher.Fetch(ctx, source)
	if err != nil {
		return SourceSyncResult{}, err
	}

	root := s.cacheRoot()
	base := path.Join(root, sourceID, source.Ref)
	paths := make([]string, 0, len(files))
	bytesWritten := 0
	for rel := range files {
		paths = append(paths, rel)
	}
	sort.Strings(paths)
	for _, rel := range paths {
		if !safeRelativeRulePath(rel) {
			return SourceSyncResult{}, fmt.Errorf("unsafe rule path %q returned for %s", rel, source.Repository)
		}
		data := files[rel]
		if err := s.FS.WriteFile(path.Join(base, rel), data, 0644); err != nil {
			return SourceSyncResult{}, fmt.Errorf("cache rule file %q: %w", rel, err)
		}
		bytesWritten += len(data)
	}

	now := time.Now().UTC()
	if s.Now != nil {
		now = s.Now().UTC()
	}
	snapshot := Snapshot{
		SourceID:   sourceID,
		Repository: source.Repository,
		Ref:        source.Ref,
		SyncedAt:   now.Format(time.RFC3339),
		Files:      paths,
	}
	metadata, err := yaml.Marshal(snapshot)
	if err != nil {
		return SourceSyncResult{}, fmt.Errorf("encode source snapshot: %w", err)
	}
	if err := s.FS.WriteFile(path.Join(base, ".source.yaml"), metadata, 0644); err != nil {
		return SourceSyncResult{}, fmt.Errorf("write source snapshot: %w", err)
	}

	return SourceSyncResult{
		SourceID:   sourceID,
		Repository: source.Repository,
		Ref:        source.Ref,
		Files:      len(paths),
		Bytes:      bytesWritten,
		CachePath:  base,
	}, nil
}

func LoadSnapshot(fsys filesystem.FileSystem, cacheRoot, sourceID string, source Source) (*Snapshot, error) {
	if fsys == nil {
		return nil, fmt.Errorf("rule snapshot filesystem is nil")
	}
	if cacheRoot == "" {
		cacheRoot = DefaultCacheRoot
	}
	data, err := fsys.ReadFile(path.Join(cacheRoot, sourceID, source.Ref, ".source.yaml"))
	if err != nil {
		return nil, err
	}
	var snapshot Snapshot
	if err := yaml.Unmarshal(data, &snapshot); err != nil {
		return nil, fmt.Errorf("decode rule snapshot for %q: %w", sourceID, err)
	}
	if snapshot.SourceID != sourceID || snapshot.Repository != source.Repository || snapshot.Ref != source.Ref {
		return nil, fmt.Errorf("rule snapshot for %q does not match pinned source", sourceID)
	}
	return &snapshot, nil
}

func (s *Synchronizer) cacheRoot() string {
	if strings.TrimSpace(s.CacheRoot) == "" {
		return DefaultCacheRoot
	}
	return strings.Trim(strings.ReplaceAll(s.CacheRoot, "\\", "/"), "/")
}

func archiveRelativePath(name string) (string, bool) {
	clean := path.Clean(strings.ReplaceAll(name, "\\", "/"))
	if clean == "." || strings.HasPrefix(clean, "/") || clean == ".." || strings.HasPrefix(clean, "../") {
		return "", false
	}
	parts := strings.Split(clean, "/")
	if len(parts) < 2 {
		return "", false
	}
	rel := path.Clean(strings.Join(parts[1:], "/"))
	if !safeRelativeRulePath(rel) {
		return "", false
	}
	return rel, true
}

func safeRelativeRulePath(value string) bool {
	clean := path.Clean(strings.ReplaceAll(value, "\\", "/"))
	return clean != "." && clean != ".." && !strings.HasPrefix(clean, "../") && !strings.HasPrefix(clean, "/")
}

func isRuleTextFile(value string) bool {
	lower := strings.ToLower(value)
	base := path.Base(lower)
	if strings.HasSuffix(lower, ".md") || strings.HasSuffix(lower, ".mdc") || strings.HasSuffix(lower, ".txt") {
		return true
	}
	switch base {
	case ".windsurfrules", "agents", "rules":
		return true
	default:
		return false
	}
}

func validateRepository(repository string) error {
	parts := strings.Split(repository, "/")
	if len(parts) != 2 {
		return fmt.Errorf("repository %q must be owner/name", repository)
	}
	for _, part := range parts {
		if err := validateCacheComponent(part); err != nil {
			return fmt.Errorf("invalid repository %q: %w", repository, err)
		}
	}
	return nil
}

func validateCacheComponent(value string) error {
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("value is empty")
	}
	for _, r := range value {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '-' || r == '_' || r == '.' {
			continue
		}
		return fmt.Errorf("value %q contains unsupported character %q", value, r)
	}
	if value == "." || value == ".." {
		return fmt.Errorf("value %q is not allowed", value)
	}
	return nil
}
