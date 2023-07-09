package testdata

import "github.com/malt3/abstractfs-core/api"

// Tree returns a test tree.
// It has no attributes set.
func Tree() api.Tree {
	return api.Tree{
		Root: &api.Node{
			Stat: api.Stat{Name: "", Kind: "directory"},
			Children: []*api.Node{
				{Stat: api.Stat{Name: "dev", Kind: "directory"}},
				{
					Stat: api.Stat{Name: "etc", Kind: "directory"},
					Children: []*api.Node{
						{
							Stat: api.Stat{
								Name: "passwd", Kind: "file",
								Payload: "sha256-tiBQcxLF6XVmo8bPr5kUT+/Big2n2UFAHfoPX1j7A2g=",
							},
						},
						{
							Stat: api.Stat{
								Name: "resolv.conf", Kind: "symlink",
								Payload: "../run/systemd/resolve/stub-resolv.conf",
							},
						},
					},
				},
				{
					Stat: api.Stat{Name: "home", Kind: "directory"},
					Children: []*api.Node{
						{
							Stat: api.Stat{Name: "malte", Kind: "directory"},
							Children: []*api.Node{
								{Stat: api.Stat{Name: ".cache", Kind: "directory"}},
								{Stat: api.Stat{Name: ".config", Kind: "directory"}},
								{Stat: api.Stat{Name: ".local", Kind: "directory"}},
								{Stat: api.Stat{Name: ".ssh", Kind: "directory"}},
								{Stat: api.Stat{Name: "Downloads", Kind: "directory"}},
							},
						},
					},
				},
				{Stat: api.Stat{Name: "root", Kind: "directory"}},
				{
					Stat: api.Stat{Name: "usr", Kind: "directory"},
					Children: []*api.Node{
						{
							Stat: api.Stat{Name: "bin", Kind: "directory"},
							Children: []*api.Node{
								{
									Stat: api.Stat{
										Name: "ls", Kind: "file",
										Payload: "sha256-tiBQcxLF6XVmo8bPr5kUT+/Big2n2UFAHfoPX1j7A2g=",
									},
								},
							},
						},
						{Stat: api.Stat{Name: "lib", Kind: "directory"}},
						{Stat: api.Stat{Name: "local", Kind: "directory"}},
						{
							Stat: api.Stat{Name: "sbin", Kind: "directory"},
							Children: []*api.Node{
								{
									Stat: api.Stat{
										Name: "init", Kind: "symlink",
										Payload: "../lib/systemd/systemd",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}
