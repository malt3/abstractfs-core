package testdata

import "github.com/malt3/abstractfs-core/api"

// Flat returns a test flat.
// It has no attributes set.
func Flat() api.Flat {
	return api.Flat{
		Files: []api.Stat{
			{Name: "/", Kind: "directory"},
			{Name: "/dev", Kind: "directory"},
			{Name: "/etc", Kind: "directory"},
			{Name: "/home", Kind: "directory"},
			{Name: "/root", Kind: "directory"},
			{Name: "/usr", Kind: "directory"},
			{Name: "/etc/passwd", Kind: "file", Payload: "sha256-tiBQcxLF6XVmo8bPr5kUT+/Big2n2UFAHfoPX1j7A2g="},
			{Name: "/etc/resolv.conf", Kind: "symlink", Payload: "../run/systemd/resolve/stub-resolv.conf"},
			{Name: "/home/malte", Kind: "directory"},
			{Name: "/usr/bin", Kind: "directory"},
			{Name: "/usr/lib", Kind: "directory"},
			{Name: "/usr/local", Kind: "directory"},
			{Name: "/usr/sbin", Kind: "directory"},
			{Name: "/home/malte/.cache", Kind: "directory"},
			{Name: "/home/malte/.config", Kind: "directory"},
			{Name: "/home/malte/.local", Kind: "directory"},
			{Name: "/home/malte/.ssh", Kind: "directory"},
			{Name: "/home/malte/Downloads", Kind: "directory"},
			{Name: "/usr/bin/ls", Kind: "file", Payload: "sha256-tiBQcxLF6XVmo8bPr5kUT+/Big2n2UFAHfoPX1j7A2g="},
			{Name: "/usr/sbin/init", Kind: "symlink", Payload: "../lib/systemd/systemd"},
		},
	}
}
