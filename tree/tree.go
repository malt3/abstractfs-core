// This package provides operations on the tree data structure defined in the api package.
package tree

import (
	"io"
	"path"
	"sort"
	"strings"

	"github.com/malt3/abstractfs-core/api"
	"github.com/malt3/abstractfs-core/traverse"
)

// FromSource returns a tree representation of the source.
func FromSource(source api.Source) (tree api.Tree, err error) {
	tree.Root = &api.Node{
		Stat: api.Stat{
			Kind: "directory",
		},
	}
	for {
		node, err := source.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return tree, err
		}
		var dir string
		dir, node.Stat.Name = normalizeFlatName(node.Stat.Name, node.Stat.Kind)
		Insert(tree, dir, node.Stat)
	}
	return tree, nil
}

// Flatten returns a flat representation of the tree.
func Flatten(tree api.Tree) (flat api.Flat) {
	visit := func(dir string, node *api.Node) {
		stat := node.Stat
		stat.Name = path.Join("/", dir, stat.Name)
		flat.Files = append(flat.Files, stat)
	}
	traverse.BFS(tree.Root, visit)
	return
}

// Unflatten returns a tree representation of the flat.
func Unflatten(flat api.Flat) (tree api.Tree) {
	tree.Root = &api.Node{
		Stat: api.Stat{
			Kind: "directory",
		},
	}
	for _, file := range flat.Files {
		var dir string
		dir, file.Name = normalizeFlatName(file.Name, file.Kind)
		Insert(tree, dir, file)
	}
	return
}

// Insert inserts an individual node into the tree.
// If any parent of the node does not exist, it will be created with default values.
// If the node already exists, it will be overwritten.
func Insert(tree api.Tree, dir string, stat api.Stat) error {
	if (dir == "/" || dir == "") && (stat.Name == "/" || stat.Name == "") {
		tree.Root.Stat = stat
		return nil
	}
	if len(dir) > 0 && dir[0] == '/' {
		dir = dir[1:]
	}
	path := strings.Split(strings.TrimSuffix(dir, "/"), "/")
	if len(path) == 1 && path[0] == "" {
		path = nil
	}
	parent := tree.Root
	for _, base := range path {
		child := findChild(parent, base)
		if child == nil {
			child = &api.Node{Stat: api.Stat{Name: base, Kind: "directory"}}
			parent.Children = append(parent.Children, child)
			sortChildren(parent)
		}
		parent = child
	}
	child := findChild(parent, stat.Name)
	if child == nil {
		parent.Children = append(parent.Children, &api.Node{Stat: stat})
		sortChildren(parent)
		return nil
	}
	child.Stat = stat
	return nil
}

// Get returns the node at the given path.
func Get(tree api.Tree, path string) *api.Node {
	if path == "/" || path == "." || path == "" {
		return tree.Root
	}
	if len(path) > 0 && path[0] == '/' {
		path = path[1:]
	}
	parts := strings.Split(strings.TrimSuffix(path, "/"), "/")
	if len(parts) == 1 && parts[0] == "" {
		parts = nil
	}
	node := tree.Root
	for _, part := range parts {
		node = findChild(node, part)
		if node == nil {
			return nil
		}
	}
	return node
}

// DeepCopyInto copies the tree rooted at src into dst.
func DeepCopyInto(src, dst *api.Node) {
	dst.Stat = src.Stat
	dst.Children = nil
	for _, child := range src.Children {
		newChild := &api.Node{}
		DeepCopyInto(child, newChild)
		dst.Children = append(dst.Children, newChild)
	}
}

func findChild(parent *api.Node, base string) *api.Node {
	for _, child := range parent.Children {
		if child.Stat.Name == base {
			return child
		}
	}
	return nil
}

func sortChildren(node *api.Node) {
	sort.Slice(node.Children, func(i, j int) bool {
		return node.Children[i].Stat.Name < node.Children[j].Stat.Name
	})
}

// normalizeFlatName normalizes a flat name.
// It returns the dir and the name.
func normalizeFlatName(name, kind string) (string, string) {
	dir := path.Dir(name)
	if strings.HasPrefix(dir, "/") {
		dir = dir[1:]
	}
	switch name {
	case "/", "":
		name = ""
	default:
		name = path.Base(name)
	}
	if kind == "directory" {
		dir = strings.TrimSuffix(dir, name)
	}
	return dir, name
}
