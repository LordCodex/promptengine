package filesystem

import (
	"io"
	"io/fs"
	"os"
	pathpkg "path"
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
	RemoveAll(path string) error
	IsSafePath(base, target string) bool
}

// OSFileSystem implements FileSystem using standard os calls bounded to BaseDir.
type OSFileSystem struct {
	BaseDir string
}

func (fs *OSFileSystem) Exists(path string) bool {
	safe, err := fs.safePath(path)
	if err != nil {
		return false
	}
	_, err = os.Stat(safe)
	return !os.IsNotExist(err)
}

func (fs *OSFileSystem) IsDir(path string) bool {
	safe, err := fs.safePath(path)
	if err != nil {
		return false
	}
	info, err := os.Stat(safe)
	if err != nil {
		return false
	}
	return info.IsDir()
}

func (fs *OSFileSystem) ReadFile(path string) ([]byte, error) {
	safe, err := fs.safePath(path)
	if err != nil {
		return nil, err
	}
	return os.ReadFile(safe)
}

func (fs *OSFileSystem) WriteFile(path string, data []byte, perm os.FileMode) error {
	safe, err := fs.safePath(path)
	if err != nil {
		return err
	}
	dir := filepath.Dir(safe)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	return os.WriteFile(safe, data, perm)
}

func (fs *OSFileSystem) MkdirAll(path string, perm os.FileMode) error {
	safe, err := fs.safePath(path)
	if err != nil {
		return err
	}
	return os.MkdirAll(safe, perm)
}

func (fs *OSFileSystem) ReadDir(path string) ([]os.DirEntry, error) {
	safe, err := fs.safePath(path)
	if err != nil {
		return nil, err
	}
	return os.ReadDir(safe)
}

func (fs *OSFileSystem) Remove(path string) error {
	safe, err := fs.safePath(path)
	if err != nil {
		return err
	}
	return os.Remove(safe)
}

func (fs *OSFileSystem) RemoveAll(path string) error {
	safe, err := fs.safePath(path)
	if err != nil {
		return err
	}
	return os.RemoveAll(safe)
}

func (fs *OSFileSystem) AppendFile(path string, data []byte, perm os.FileMode) error {
	safe, err := fs.safePath(path)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(safe), 0755); err != nil {
		return err
	}
	f, err := os.OpenFile(safe, os.O_CREATE|os.O_WRONLY|os.O_APPEND, perm)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(data)
	return err
}

func (fs *OSFileSystem) safePath(target string) (string, error) {
	base := fs.BaseDir
	if base == "" {
		base = "."
	}
	cleanBase, err := filepath.Abs(filepath.Clean(base))
	if err != nil {
		return "", err
	}
	resolved := filepath.Clean(target)
	if !filepath.IsAbs(resolved) {
		resolved = filepath.Join(cleanBase, resolved)
	}
	if !fs.IsSafePath(cleanBase, resolved) {
		return "", os.ErrPermission
	}
	return resolved, nil
}

// IsSafePath checks target file bounds to prevent traversal attacks
func (fs *OSFileSystem) IsSafePath(base, target string) bool {
	cleanBase, err := filepath.Abs(filepath.Clean(base))
	if err != nil {
		return false
	}
	cleanTarget, err := filepath.Abs(filepath.Clean(target))
	if err != nil {
		return false
	}
	if resolvedBase, err := filepath.EvalSymlinks(cleanBase); err == nil {
		cleanBase = resolvedBase
	}
	if resolvedTarget, err := resolveExistingPath(cleanTarget); err == nil {
		cleanTarget = resolvedTarget
	}
	rel, err := filepath.Rel(cleanBase, cleanTarget)
	if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return false
	}
	return true
}

func resolveExistingPath(path string) (string, error) {
	if resolved, err := filepath.EvalSymlinks(path); err == nil {
		return resolved, nil
	}
	parent := filepath.Dir(path)
	resolvedParent, err := filepath.EvalSymlinks(parent)
	if err != nil {
		return path, err
	}
	return filepath.Join(resolvedParent, filepath.Base(path)), nil
}

// MockFileSystem implements FileSystem in memory for unit testing. Paths are
// stored in slash form so fixtures behave identically on Windows, macOS, and Linux.
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

func normalizeMockPath(value string) string {
	value = strings.ReplaceAll(value, "\\", "/")
	clean := pathpkg.Clean(value)
	if clean == "" {
		return "."
	}
	return clean
}

func (fs *MockFileSystem) Exists(path string) bool {
	path = normalizeMockPath(path)
	if _, ok := fs.Files[path]; ok {
		return true
	}
	if fs.Dirs[path] {
		return true
	}
	prefix := path + "/"
	for k := range fs.Files {
		if strings.HasPrefix(normalizeMockPath(k), prefix) {
			return true
		}
	}
	return false
}

func (fs *MockFileSystem) IsDir(path string) bool {
	return fs.Dirs[normalizeMockPath(path)]
}

func (fs *MockFileSystem) ReadFile(path string) ([]byte, error) {
	data, ok := fs.Files[normalizeMockPath(path)]
	if !ok {
		return nil, io.EOF
	}
	return data, nil
}

func (fs *MockFileSystem) WriteFile(path string, data []byte, perm os.FileMode) error {
	path = normalizeMockPath(path)
	fs.Files[path] = data
	dir := pathpkg.Dir(path)
	for dir != "." && dir != "/" && dir != "" {
		fs.Dirs[dir] = true
		parent := pathpkg.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	fs.Dirs["."] = true
	return nil
}

func (fs *MockFileSystem) MkdirAll(path string, perm os.FileMode) error {
	path = normalizeMockPath(path)
	fs.Dirs[path] = true
	dir := path
	for dir != "." && dir != "/" && dir != "" {
		fs.Dirs[dir] = true
		parent := pathpkg.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	fs.Dirs["."] = true
	return nil
}

func (fs *MockFileSystem) ReadDir(path string) ([]os.DirEntry, error) {
	clean := normalizeMockPath(path)
	if !fs.Exists(clean) && clean != "." {
		return nil, io.EOF
	}
	seen := map[string]mockDirEntry{}
	add := func(candidate string, isDir bool) {
		candidate = normalizeMockPath(candidate)
		parent := pathpkg.Dir(candidate)
		if (parent == "." && clean == ".") || parent == clean {
			name := pathpkg.Base(candidate)
			seen[name] = mockDirEntry{name: name, dir: isDir}
			return
		}
		if clean == "." && !strings.Contains(candidate, "/") {
			seen[candidate] = mockDirEntry{name: candidate, dir: isDir}
		}
	}
	for file := range fs.Files {
		add(file, false)
	}
	for dir := range fs.Dirs {
		if normalizeMockPath(dir) == clean {
			continue
		}
		add(dir, true)
	}
	out := make([]os.DirEntry, 0, len(seen))
	for _, entry := range seen {
		out = append(out, entry)
	}
	return out, nil
}

func (fs *MockFileSystem) Remove(path string) error {
	path = normalizeMockPath(path)
	delete(fs.Files, path)
	delete(fs.Dirs, path)
	return nil
}

func (fs *MockFileSystem) RemoveAll(path string) error {
	clean := normalizeMockPath(path)
	for file := range fs.Files {
		normalized := normalizeMockPath(file)
		if normalized == clean || strings.HasPrefix(normalized, clean+"/") {
			delete(fs.Files, file)
		}
	}
	for dir := range fs.Dirs {
		normalized := normalizeMockPath(dir)
		if normalized == clean || strings.HasPrefix(normalized, clean+"/") {
			delete(fs.Dirs, dir)
		}
	}
	return nil
}

func (fs *MockFileSystem) IsSafePath(base, target string) bool {
	cleanBase := normalizeMockPath(base)
	cleanTarget := normalizeMockPath(target)
	if strings.HasPrefix(cleanTarget, "/") || (len(cleanTarget) >= 2 && cleanTarget[1] == ':') {
		return cleanTarget == cleanBase || strings.HasPrefix(cleanTarget, strings.TrimSuffix(cleanBase, "/")+"/")
	}
	rel, err := pathpkg.Rel(cleanBase, cleanTarget)
	if err != nil {
		return false
	}
	return rel == "." || (rel != ".." && !strings.HasPrefix(rel, "../"))
}

type mockDirEntry struct {
	name string
	dir  bool
}

func (e mockDirEntry) Name() string { return e.name }
func (e mockDirEntry) IsDir() bool  { return e.dir }
func (e mockDirEntry) Type() fs.FileMode {
	if e.dir {
		return fs.ModeDir
	}
	return 0
}
func (e mockDirEntry) Info() (fs.FileInfo, error) { return nil, nil }
