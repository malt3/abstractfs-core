package api

// Tree is a abstractfs tree.
type Tree struct {
	// Root is the root node of the tree.
	Root *Node
}

// Flat is a flat representation of a tree.
// It is used to serialize a tree.
type Flat struct {
	Files []Stat
}
