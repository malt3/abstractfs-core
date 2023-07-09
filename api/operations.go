package api

import (
	"io"
	"io/fs"
)

type Source interface {
	// Next returns the next node in the source.
	// If there are no more nodes, it returns io.EOF.
	Next() (SourceNode, error)
}

// Sink is a sink for nodes.
// It is used to write nodes to a destination.
type Sink interface {
	Consume(in fs.FS) error
}

type SourceNode struct {
	Stat Stat
	Open func() (io.ReadCloser, error)
}

type CAS interface {
	CASReader
	CASWriter
}

type CASReader interface {
	// Open returns a reader for the given SRI.
	// If the SRI does not exist, it returns fs.ErrNotExist.
	Open(sri string) (io.ReadCloser, error)
}

type CASWriter interface {
	Write(sri string, r io.Reader) error
}

// CloseWaitFunc is a function that closes a resource and waits for it to be closed.
type CloseWaitFunc func() error
