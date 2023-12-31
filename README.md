# abstractfs-core

[![Go Report Card](https://goreportcard.com/badge/github.com/malt3/abstractfs-core)](https://goreportcard.com/report/github.com/malt3/abstractfs-core)
[![GoDoc](https://img.shields.io/static/v1?label=godoc&message=reference&color=blue)](https://pkg.go.dev/github.com/malt3/abstractfs-core)

This is the core library of [abstractfs](https://github.com/malt3/abstractfs).
It is split out as a separate module / repository to allow it to have no dependencies (except for the go stdlib).
This allows users of abstractfs to reuse the core and add their own sources, sinks and CAS implementations.

## Tests

Test can be found in [`/tests`](tests). They can only be run from the tests directory, since they are their own go module:

```shell-session
git clone https://github.com/malt3/abstractfs-core
cd abstractfs-core/tests
go test ./...
```