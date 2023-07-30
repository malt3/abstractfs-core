package tree

import (
	"io"
	"io/fs"
	"time"

	"github.com/malt3/abstractfs-core/api"
	"github.com/malt3/abstractfs-core/cas/recorder"
	"github.com/malt3/abstractfs-core/sri"
)

// TreeFS implements io/fs.FS and io/fs.ReadDirFS for a tree.
type TreeFS struct {
	api.Tree
	api.CASReader
}

func (t *TreeFS) Open(name string) (fs.File, error) {
	node := Get(t.Tree, name)
	if node == nil {
		return nil, &fs.PathError{Op: "open", Path: name, Err: fs.ErrNotExist}
	}
	if node.Stat.Kind == api.KindDirectory {
		return &readDirFile{node: node}, nil
	}
	if node.Stat.Kind != api.KindRegular {
		return nil, &fs.PathError{Op: "open", Path: name, Err: fs.ErrInvalid}
	}
	inner, err := t.CASReader.Open(node.Stat.Payload)
	if err != nil {
		return nil, &fs.PathError{Op: "open", Path: name, Err: err}
	}
	return &file{node: node, readCloser: inner}, nil
}

func (t *TreeFS) ReadDir(name string) ([]fs.DirEntry, error) {
	node := Get(t.Tree, name)
	if node == nil {
		return nil, &fs.PathError{Op: "readdir", Path: name, Err: fs.ErrNotExist}
	}
	if node.Stat.Kind != api.KindDirectory {
		return nil, &fs.PathError{Op: "readdir", Path: name, Err: fs.ErrInvalid}
	}
	var entries []fs.DirEntry
	for _, child := range node.Children {
		entries = append(entries, dirEntry{child.Stat})
	}
	return entries, nil
}

// Record records all file contents of the tree to a io.Writer.
// The format is compatible with the recorder protocol.
func (t *TreeFS) Record(w io.Writer) error {
	return fs.WalkDir(t, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.Type().IsRegular() {
			return nil
		}

		fInfo, err := d.Info()
		if err != nil {
			return err
		}
		stat, ok := fInfo.Sys().(api.Stat)
		if !ok {
			return &fs.PathError{Op: "record", Path: path, Err: fs.ErrInvalid}
		}
		sri, error := sri.FromString(stat.Payload)
		if error != nil {
			return err
		}

		file, err := t.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		return recorder.Encode(w, sri, stat.Size, file)
	})
}

// file implements fs.File for a node.
type file struct {
	node       *api.Node
	readCloser io.ReadCloser
}

func (f *file) Stat() (fs.FileInfo, error) {
	return fileInfo{f.node.Stat}, nil
}

func (f *file) Read(p []byte) (int, error) {
	return f.readCloser.Read(p)
}

func (f *file) Close() error {
	return f.readCloser.Close()
}

// readDirFile implements fs.File for a directory node.
type readDirFile struct {
	node *api.Node
	pos  int
}

func (f *readDirFile) Stat() (fs.FileInfo, error) {
	return fileInfo{f.node.Stat}, nil
}

func (f *readDirFile) Read([]byte) (int, error) {
	return 0, fs.ErrInvalid
}

func (f *readDirFile) Close() error {
	return nil
}

func (f *readDirFile) ReadDir(n int) ([]fs.DirEntry, error) {
	if f.pos >= len(f.node.Children) {
		return nil, io.EOF
	}
	if n <= 0 {
		n = len(f.node.Children)
	}
	var entries []fs.DirEntry
	for i := f.pos; i < len(f.node.Children) && i < f.pos+n; i++ {
		entries = append(entries, dirEntry{f.node.Children[i].Stat})
	}
	f.pos += n
	return entries, nil
}

type dirEntry struct {
	stat api.Stat
}

func (d dirEntry) Name() string {
	return d.stat.Name
}

func (d dirEntry) IsDir() bool {
	return d.stat.Kind == api.KindDirectory
}

func (d dirEntry) Type() fs.FileMode {
	// TODO: think about support for
	// ModeSetuid
	// ModeSetgid
	// ModeCharDevice
	// ModeSticky
	// ModeIrregular
	return typeMode(d.stat.Kind)
}

func (d dirEntry) Info() (fs.FileInfo, error) {
	return fileInfo{d.stat}, nil
}

type fileInfo struct {
	stat api.Stat
}

func (f fileInfo) Name() string {
	return f.stat.Name
}

func (f fileInfo) Size() int64 {
	return f.stat.Size
}

func (f fileInfo) Mode() fs.FileMode {
	// TODO: think about support for
	// ModeSetuid
	// ModeSetgid
	// ModeCharDevice
	// ModeSticky
	// ModeIrregular
	return typeMode(f.stat.Kind)
}

func (f fileInfo) ModTime() time.Time {
	return f.stat.Attributes.Mtime
}

func (f fileInfo) IsDir() bool {
	return f.stat.Kind == api.KindDirectory
}

func (f fileInfo) Sys() any {
	return f.stat
}

func typeMode(kind string) fs.FileMode {
	switch kind {
	case api.KindDirectory:
		return fs.ModeDir
	case api.KindSymlink:
		return fs.ModeSymlink
	case api.KindRegular:
		return 0
	default:
		panic("not implemented")
	}
}

var _ fs.FS = (*TreeFS)(nil)
var _ fs.ReadDirFS = (*TreeFS)(nil)
