package cas

import (
	"io"
	"os"

	"github.com/malt3/abstractfs-core/api"
)

// EmptyCAS is a CAS implementation that always returns os.ErrNotExist.
type EmptyCAS struct{}

func (c *EmptyCAS) Open(_ string) (io.ReadCloser, error) {
	return nil, os.ErrNotExist
}

var _ api.CASReader = (*EmptyCAS)(nil)
