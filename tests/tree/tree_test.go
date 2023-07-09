package tree_test

import (
	"testing"

	"github.com/malt3/abstractfs-core/api"
	"github.com/malt3/abstractfs-core/tests/internal/testdata"
	coretree "github.com/malt3/abstractfs-core/tree"
	"github.com/stretchr/testify/assert"
)

func TestFlatten(t *testing.T) {
	assert := assert.New(t)
	tree := testdata.Tree()
	wantFlat := testdata.Flat()
	gotFlat := coretree.Flatten(tree)
	assert.Equal(wantFlat, gotFlat)
}

func TestUnflatten(t *testing.T) {
	assert := assert.New(t)
	flat := testdata.Flat()
	wantTree := testdata.Tree()
	gotTree := coretree.Unflatten(flat)
	assert.Equal(wantTree, gotTree)
}

func TestFlattenUnflatten(t *testing.T) {
	assert := assert.New(t)
	tree := testdata.Tree()
	gotTree := coretree.Unflatten(coretree.Flatten(tree))
	assert.Equal(tree, gotTree)
}

func TestUnflattenFlatten(t *testing.T) {
	assert := assert.New(t)
	flat := testdata.Flat()
	gotFlat := coretree.Flatten(coretree.Unflatten(flat))
	assert.Equal(flat, gotFlat)
}

func TestDeepCopyInto(t *testing.T) {
	assert := assert.New(t)
	original := testdata.Tree()
	gotTree := api.Tree{
		Root: &api.Node{},
	}
	coretree.DeepCopyInto(original.Root, gotTree.Root)
	assert.Equal(testdata.Tree(), gotTree)

	// Modify original tree. Deep copied tree should not be affected.
	original.Root.Stat.Name = "foo"
	assert.NotEqual(testdata.Tree(), original)
	assert.Equal(testdata.Tree(), gotTree)

}
