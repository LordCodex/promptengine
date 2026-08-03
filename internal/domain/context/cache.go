package context

import (
	"crypto/sha256"
	"fmt"
	"strings"
	"sync"

	"github.com/LordCodex/promptengine/internal/filesystem"
)

type Cache struct {
	mu      sync.RWMutex
	entries map[string]cacheEntry
}

type cacheEntry struct {
	Fingerprint string
	Package     *ContextPackage
}

func NewCache() *Cache {
	return &Cache{entries: map[string]cacheEntry{}}
}

func (c *Cache) Get(key, fingerprint string) (*ContextPackage, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	entry, ok := c.entries[key]
	if !ok || entry.Fingerprint != fingerprint {
		return nil, false
	}
	pkg := clonePackage(entry.Package)
	pkg.Summary.CacheHit = true
	return pkg, true
}

func (c *Cache) Set(key, fingerprint string, pkg *ContextPackage) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries[key] = cacheEntry{Fingerprint: fingerprint, Package: clonePackage(pkg)}
}

func fingerprint(fs filesystem.FileSystem, paths []string, salt string) string {
	h := sha256.New()
	h.Write([]byte(salt))
	for _, path := range unique(paths) {
		h.Write([]byte(path))
		if data, err := fs.ReadFile(path); err == nil {
			sum := sha256.Sum256(data)
			h.Write(sum[:])
		}
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}

func cacheKey(req ContextRequest) string {
	return strings.Join([]string{string(req.TaskType), req.WorkflowType, req.UserIntent, req.RequestedOperation, strings.Join(req.AffectedFiles, ","), fmt.Sprintf("%.2f", req.MinRelevanceScore), fmt.Sprintf("%t", req.Explain)}, "|")
}

func clonePackage(pkg *ContextPackage) *ContextPackage {
	if pkg == nil {
		return nil
	}
	cp := *pkg
	cp.Items = append([]ContextItem(nil), pkg.Items...)
	cp.ExcludedItems = append([]ContextItem(nil), pkg.ExcludedItems...)
	cp.Documents = append([]DocumentItem(nil), pkg.Documents...)
	cp.SelectedFiles = append([]string(nil), pkg.SelectedFiles...)
	cp.SelectedDocs = append([]string(nil), pkg.SelectedDocs...)
	cp.RelevantStandards = append([]string(nil), pkg.RelevantStandards...)
	cp.RelatedPlaybooks = append([]string(nil), pkg.RelatedPlaybooks...)
	cp.Reasoning = append([]string(nil), pkg.Reasoning...)
	cp.Summary.DroppedFiles = append([]string(nil), pkg.Summary.DroppedFiles...)
	cp.Explanations = map[string]string{}
	for k, v := range pkg.Explanations {
		cp.Explanations[k] = v
	}
	cp.ProjectMetadata = map[string]string{}
	for k, v := range pkg.ProjectMetadata {
		cp.ProjectMetadata[k] = v
	}
	return &cp
}
