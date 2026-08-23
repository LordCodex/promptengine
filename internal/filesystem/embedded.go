package filesystem

import (
	"io/fs"
	"os"
)

// EmbeddedFileSystem adapts an fs.FS to PromptEngine's FileSystem boundary.
// It is intentionally read-only and is used for the bundled standards library.
type EmbeddedFileSystem struct {
	fs fs.FS
}

func NewEmbeddedFileSystem(source fs.FS) *EmbeddedFileSystem {
	return &EmbeddedFileSystem{fs: source}
}

func (e *EmbeddedFileSystem) Exists(path string) bool {
	if e == nil || e.fs == nil {
		return false
	}
	_, err := fs.Stat(e.fs, path)
	return err == nil
}

func (e *EmbeddedFileSystem) IsDir(path string) bool {
	if e == nil || e.fs == nil {
		return false
	}
	info, err := fs.Stat(e.fs, path)
	return err == nil && info.IsDir()
}

func (e *EmbeddedFileSystem) ReadFile(path string) ([]byte, error) {
	if e == nil || e.fs == nil {
		return nil, fs.ErrNotExist
	}
	return fs.ReadFile(e.fs, path)
}

func (e *EmbeddedFileSystem) ReadDir(path string) ([]os.DirEntry, error) {
	if e == nil || e.fs == nil {
		return nil, fs.ErrNotExist
	}
	entries, err := fs.ReadDir(e.fs, path)
	if err != nil {
		return nil, err
	}
	out := make([]os.DirEntry, 0, len(entries))
	for _, entry := range entries {
		out = append(out, entry)
	}
	return out, nil
}

func (e *EmbeddedFileSystem) WriteFile(string, []byte, os.FileMode) error { return fs.ErrPermission }
func (e *EmbeddedFileSystem) MkdirAll(string, os.FileMode) error          { return fs.ErrPermission }
func (e *EmbeddedFileSystem) Remove(string) error                         { return fs.ErrPermission }
func (e *EmbeddedFileSystem) RemoveAll(string) error                      { return fs.ErrPermission }
func (e *EmbeddedFileSystem) IsSafePath(_, target string) bool {
	return target != "" && fs.ValidPath(target)
}

// OverlayFileSystem reads project files from Primary first and falls back to
// Library for bundled PromptEngine resources. All writes and directory walks
// remain project-only so bundled standards never appear as project source code.
type OverlayFileSystem struct {
	Primary FileSystem
	Library FileSystem
}

func NewOverlayFileSystem(primary, library FileSystem) *OverlayFileSystem {
	return &OverlayFileSystem{Primary: primary, Library: library}
}

func (o *OverlayFileSystem) Exists(path string) bool {
	return (o.Primary != nil && o.Primary.Exists(path)) || (o.Library != nil && o.Library.Exists(path))
}

func (o *OverlayFileSystem) IsDir(path string) bool {
	if o.Primary != nil && o.Primary.Exists(path) {
		return o.Primary.IsDir(path)
	}
	return o.Library != nil && o.Library.IsDir(path)
}

func (o *OverlayFileSystem) ReadFile(path string) ([]byte, error) {
	if o.Primary != nil && o.Primary.Exists(path) {
		return o.Primary.ReadFile(path)
	}
	if o.Library != nil {
		return o.Library.ReadFile(path)
	}
	return nil, fs.ErrNotExist
}

// ReadDir intentionally exposes only the project workspace. Discovery must not
// treat the bundled PromptEngine library as files belonging to the user's app.
func (o *OverlayFileSystem) ReadDir(path string) ([]os.DirEntry, error) {
	if o.Primary == nil {
		return nil, fs.ErrNotExist
	}
	return o.Primary.ReadDir(path)
}

func (o *OverlayFileSystem) WriteFile(path string, data []byte, perm os.FileMode) error {
	if o.Primary == nil {
		return fs.ErrPermission
	}
	return o.Primary.WriteFile(path, data, perm)
}

func (o *OverlayFileSystem) MkdirAll(path string, perm os.FileMode) error {
	if o.Primary == nil {
		return fs.ErrPermission
	}
	return o.Primary.MkdirAll(path, perm)
}

func (o *OverlayFileSystem) Remove(path string) error {
	if o.Primary == nil {
		return fs.ErrPermission
	}
	return o.Primary.Remove(path)
}

func (o *OverlayFileSystem) RemoveAll(path string) error {
	if o.Primary == nil {
		return fs.ErrPermission
	}
	return o.Primary.RemoveAll(path)
}

func (o *OverlayFileSystem) IsSafePath(base, target string) bool {
	if o.Primary != nil && o.Primary.Exists(target) {
		return o.Primary.IsSafePath(base, target)
	}
	return o.Library != nil && o.Library.IsSafePath(base, target)
}
