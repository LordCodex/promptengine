package cache

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Cache manages state caching to improve CLI performance
type Cache struct {
	mu       sync.RWMutex
	cacheDir string
}

type cacheEntry struct {
	Value      []byte    `json:"value"`
	Expiration time.Time `json:"expiration"`
	SourceHash string    `json:"source_hash"`
}

func NewCache(rootDir string) *Cache {
	return &Cache{
		cacheDir: filepath.Join(rootDir, ".promptengine", "cache"),
	}
}

// ComputeHash generates SHA256 of files to detect drift
func (c *Cache) ComputeHash(filePath string) (string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", err
	}
	return fmt.Sprintf("%x", sha256.Sum256(data)), nil
}

func (c *Cache) Get(key string, target interface{}, currentHash string) (bool, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	cPath := filepath.Join(c.cacheDir, key+".json")
	data, err := os.ReadFile(cPath)
	if err != nil {
		return false, nil
	}

	var entry cacheEntry
	if err := json.Unmarshal(data, &entry); err != nil {
		return false, err
	}

	if !entry.Expiration.IsZero() && time.Now().After(entry.Expiration) {
		return false, nil
	}

	if currentHash != "" && entry.SourceHash != currentHash {
		return false, nil // cached source model has drifted
	}

	if err := json.Unmarshal(entry.Value, target); err != nil {
		return false, err
	}

	return true, nil
}

func (c *Cache) Set(key string, val interface{}, duration time.Duration, sourceHash string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if err := os.MkdirAll(c.cacheDir, 0755); err != nil {
		return err
	}

	rawVal, err := json.Marshal(val)
	if err != nil {
		return err
	}

	var exp time.Time
	if duration > 0 {
		exp = time.Now().Add(duration)
	}

	entry := cacheEntry{
		Value:      rawVal,
		Expiration: exp,
		SourceHash: sourceHash,
	}

	data, err := json.Marshal(entry)
	if err != nil {
		return err
	}

	cPath := filepath.Join(c.cacheDir, key+".json")
	return os.WriteFile(cPath, data, 0644)
}

func (c *Cache) Clear() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return os.RemoveAll(c.cacheDir)
}
