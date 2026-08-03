package filesystem

import (
	"io"
	"os"
	"path/filepath"
	"strings"
)

// FileSystem defines mockable file system boundary operations
type FileSystem interface {
	Exists(path string) bool
	IsDir(path string) bool
	ReadFile(path string) ([]byte, error)
	WriteFile(path string, data []byte, perm os.FileMode) error
	MkdirAll(path string, perm os.FileMode) error
	ReadDir(path string) ([]os.DirEntry, error)
	Remove(path string) error
	IsSafePath(base, target string) bool
}

// OSFileSystem implements FileSystem using standard os calls
type OSFileSystem struct{}

func (fs *OSFileSystem) Exists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func (fs *OSFileSystem) IsDir(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

func (fs *OSFileSystem) ReadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

func (fs *OSFileSystem) WriteFile(path string, data []byte, perm os.FileMode) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	return os.WriteFile(path, data, perm)
}

func (fs *OSFileSystem) MkdirAll(path string, perm os.FileMode) error {
	return os.MkdirAll(path, perm)
}

func (fs *OSFileSystem) ReadDir(path string) ([]os.DirEntry, error) {
	return os.ReadDir(path)
}

func (fs *OSFileSystem) Remove(path string) error {
	return os.Remove(path)
}

// IsSafePath checks target file bounds to prevent traversal attacks
func (fs *OSFileSystem) IsSafePath(base, target string) bool {
	cleanBase := filepath.Clean(base)
	cleanTarget := filepath.Clean(target)
	if !filepath.IsAbs(cleanBase) {
		if abs, err := filepath.Abs(cleanBase); err == nil {
			cleanBase = abs
		}
	}
	if !filepath.IsAbs(cleanTarget) {
		if abs, err := filepath.Abs(cleanTarget); err == nil {
			cleanTarget = abs
		}
	}
	rel, err := filepath.Rel(cleanBase, cleanTarget)
	if err != nil || strings.HasPrefix(rel, "..") {
		return false
	}
	return true
}

// MockFileSystem implements FileSystem in memory for unit testing
type MockFileSystem struct {
	Files map[string][]byte
	Dirs  map[string]bool
}

func NewMockFileSystem() *MockFileSystem {
	return &MockFileSystem{
		Files: make(map[string][]byte),
		Dirs:  make(map[string]bool),
	}
}

func (fs *MockFileSystem) Exists(path string) bool {
	if _, ok := fs.Files[path]; ok {
		return true
	}
	if fs.Dirs[path] {
		return true
	}
	prefix := path + "/"
	for k := range fs.Files {
		if strings.HasPrefix(k, prefix) {
			return true
		}
	}
	return false
}

func (fs *MockFileSystem) IsDir(path string) bool {
	return fs.Dirs[path]
}

func (fs *MockFileSystem) ReadFile(path string) ([]byte, error) {
	data, ok := fs.Files[path]
	if !ok {
		return nil, io.EOF
	}
	return data, nil
}

func (fs *MockFileSystem) WriteFile(path string, data []byte, perm os.FileMode) error {
	fs.Files[path] = data
	return nil
}

func (fs *MockFileSystem) MkdirAll(path string, perm os.FileMode) error {
	fs.Dirs[path] = true
	return nil
}

func (fs *MockFileSystem) ReadDir(path string) ([]os.DirEntry, error) {
	return nil, nil
}

func (fs *MockFileSystem) Remove(path string) error {
	delete(fs.Files, path)
	delete(fs.Dirs, path)
	return nil
}

func (fs *MockFileSystem) IsSafePath(base, target string) bool {
	return !strings.Contains(target, "..")
}
