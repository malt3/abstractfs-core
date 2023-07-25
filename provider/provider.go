package provider

import (
	"errors"

	"github.com/malt3/abstractfs-core/api"
)

type Provider interface {
	Name() string
	SourceBuilder() SourceBuilder
	SinkBuilder() SinkBuilder
	CAS() (api.CAS, api.CloseWaitFunc, error)
	CASReader() (api.CASReader, api.CloseWaitFunc, error)
	CASWriter() (api.CASWriter, api.CloseWaitFunc, error)
}

type SourceBuilder interface {
	WithSourceRef(string) SourceBuilder
	Build() (api.Source, api.CloseWaitFunc, error)
}

type SinkBuilder interface {
	WithSinkRef(string) SinkBuilder
	Set(string, any) SinkBuilder
	Build() (api.Sink, api.CloseWaitFunc, error)
}

type SourceOptions interface {
	SourceRef() string
}

type SinkOptions interface {
	SinkRef() string
}

type UnsupportedSourceBuilder struct{}

func (b UnsupportedSourceBuilder) WithSourceRef(_ string) SourceBuilder {
	return b
}

func (b UnsupportedSourceBuilder) Set(_ string, _ any) SourceBuilder {
	return b
}

func (b UnsupportedSourceBuilder) Build() (api.Source, api.CloseWaitFunc, error) {
	return nil, nil, ErrUnsupported
}

type UnsupportedSinkBuilder struct{}

func (b UnsupportedSinkBuilder) WithSinkRef(_ string) SinkBuilder {
	return b
}

func (b UnsupportedSinkBuilder) Set(_ string, _ any) SinkBuilder {
	return b
}

func (b UnsupportedSinkBuilder) Build() (api.Sink, api.CloseWaitFunc, error) {
	return nil, nil, ErrUnsupported
}

var ErrUnsupported = errors.New("unsupported")

const (
	OptionCASAlgorithm = "cas-algorithm"
)
