package testdata

import (
	"time"

	"go.uber.org/zap"
)

func evaluateOptions(options []SampleOption) *sampleOptions {
	opts := &sampleOptions{}
	for _, option := range options {
		option.apply(opts)
	}

	return opts
}

type SampleOption interface{ apply(*sampleOptions) }

type sampleOptions struct {
	foo string
	bar time.Duration
	baz *zap.Logger
}

type fooOption string

func (o fooOption) apply(opts *sampleOptions) { opts.foo = string(o) }

func WithFoo(value string) SampleOption { return fooOption(value) }

type barOption time.Duration

func (o barOption) apply(opts *sampleOptions) { opts.bar = time.Duration(o) }

func WithBar(value time.Duration) SampleOption { return barOption(value) }

type loggerOption struct {
	logger *zap.Logger
}

func (o loggerOption) apply(opts *sampleOptions) { opts.baz = o.logger }

func WithBaz(logger *zap.Logger) SampleOption { return loggerOption{logger: logger} }
