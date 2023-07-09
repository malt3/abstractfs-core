package traverse_test

import (
	"path"
	"testing"

	"github.com/malt3/abstractfs-core/api"
	"github.com/malt3/abstractfs-core/tests/internal/testdata"
	"github.com/malt3/abstractfs-core/traverse"
	"github.com/stretchr/testify/assert"
)

func TestBFS(t *testing.T) {
	assert := assert.New(t)
	var visited []string
	visit := func(dir string, node *api.Node) {
		visited = append(visited, path.Join("/", dir, node.Stat.Name))
	}
	tree := testdata.Tree()
	traverse.BFS(tree.Root, visit)
	expected := []string{
		"/",
		"/dev",
		"/etc",
		"/home",
		"/root",
		"/usr",
		"/etc/passwd",
		"/etc/resolv.conf",
		"/home/malte",
		"/usr/bin",
		"/usr/lib",
		"/usr/local",
		"/usr/sbin",
		"/home/malte/.cache",
		"/home/malte/.config",
		"/home/malte/.local",
		"/home/malte/.ssh",
		"/home/malte/Downloads",
		"/usr/bin/ls",
		"/usr/sbin/init",
	}
	assert.Equal(expected, visited)
}
