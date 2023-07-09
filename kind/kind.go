package kind

import (
	"os"

	"github.com/malt3/abstractfs-core/api"
)

func FromMode(mode os.FileMode) string {
	if mode.IsDir() {
		return api.KindDirectory
	}

	if mode&os.ModeSymlink != 0 {
		return api.KindSymlink
	}

	return api.KindRegular
}
