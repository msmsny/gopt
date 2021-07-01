# gopt

[![Go Report Card](https://goreportcard.com/badge/github.com/msmsny/gopt)](https://goreportcard.com/report/github.com/msmsny/gopt)
[![Test](https://github.com/msmsny/gopt/actions/workflows/test.yml/badge.svg)](https://github.com/msmsny/gopt/actions/workflows/test.yml)
[![Coverage Status](https://coveralls.io/repos/github/msmsny/gopt/badge.svg?branch=master)](https://coveralls.io/github/msmsny/gopt?branch=master)

Functional options code generator.

## Install

```bash
$ go get github.com/msmsny/gopt
```

## Usage

```bash
$ gopt --help
gopt generates functional options pattern code

Usage:
  gopt [flags]

Flags:
      --name string       functional options name to specify variadic functions arguments (required)
      --options strings   option names and values, e.g.: foo:string,bar:int,baz:bool
      --package string    output package name
  -o, --output string     output file name
      --evaluate          output evaluateOptions (default true)
      --format-imports    format import statement by goimports
      --header            generated code header with signature "Code generated..."
                          this option is enabled only if the output option is not empty (default true)
  -h, --help              help for gopt
```

Output stdout

```bash
$ gopt --name sample --options foo:string,bar:int,baz:bool
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
	bar int
	baz bool
}

type fooOption string

func (o fooOption) apply(opts *sampleOptions) { opts.foo = string(o) }

func WithFoo(value string) SampleOption { return fooOption(value) }

type barOption int

func (o barOption) apply(opts *sampleOptions) { opts.bar = int(o) }

func WithBar(value int) SampleOption { return barOption(value) }

type bazOption bool

func (o bazOption) apply(opts *sampleOptions) { opts.baz = bool(o) }

func WithBaz(value bool) SampleOption { return bazOption(value) }
```

Use generated code

```go
// implement new function
type Sample struct {
	foo string
	bar int
	baz bool
}

func NewSample(requiredArg string, opts ...SampleOption) *Sample {
	opt := evaluateOptions(opts)

    // ...

	return &Sample{
		foo: opt.foo,
		bar: opt.bar,
		baz: opt.baz,
	}
}

// call function
func callNewSample() {
	sample := NewSample(
		"required arg",
		WithFoo("foo"),
		WithBar(1),
		WithBaz(true),
	)
	// ...
}
```

### Command options

Output file

```bash
$ gopt --name sample --options foo:string,bar:int,baz:bool -o sample.go
```

Output file includes the header with generated code [signature](https://golang.org/cmd/go/#hdr-Generate_Go_files_by_processing_source).

With package name

```bash
$ gopt --name sample --options foo:string,bar:int,baz:bool --package gopt
package gopt

func evaluateOptions(options []SampleOption) *sampleOptions {
	opts := &sampleOptions{}
	for _, option := range options {
		option.apply(opts)
	}

	return opts
}

type SampleOption interface{ apply(*sampleOptions) }

// ...
```

Without evaluateOptions

```bash
$ gopt --name sample --options foo:string,bar:int,baz:bool --package gopt --evaluate=false
package gopt

type SampleOption interface{ apply(*sampleOptions) }

// ...
```

With package type

```bash
$ gopt --name sample --options 'foo:string,bar:int,baz:*go.uber.org/zap.Logger' --package gopt
package gopt

import (
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
	bar int
	baz *zap.Logger
}

type fooOption string

func (o fooOption) apply(opts *sampleOptions) { opts.foo = string(o) }

func WithFoo(value string) SampleOption { return fooOption(value) }

type barOption int

func (o barOption) apply(opts *sampleOptions) { opts.bar = int(o) }

func WithBar(value int) SampleOption { return barOption(value) }

type loggerOption struct {
	logger *zap.Logger
}

func (o loggerOption) apply(opts *sampleOptions) { opts.baz = o.logger }

func WithBaz(logger *zap.Logger) SampleOption { return loggerOption{logger: logger} }
```

With local package type

```bash
$ gopt --name sample --options 'foo:string,bar:int,baz:*Logger' --package gopt --format-imports
package gopt

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
	bar int
	baz *Logger
}

type fooOption string

func (o fooOption) apply(opts *sampleOptions) { opts.foo = string(o) }

func WithFoo(value string) SampleOption { return fooOption(value) }

type barOption int

func (o barOption) apply(opts *sampleOptions) { opts.bar = int(o) }

func WithBar(value int) SampleOption { return barOption(value) }

type loggerOption struct {
	logger *Logger
}

func (o loggerOption) apply(opts *sampleOptions) { opts.baz = o.logger }

func WithBaz(logger *Logger) SampleOption { return loggerOption{logger: logger} }
```

Format import statement (requires [goimports](https://pkg.go.dev/golang.org/x/tools/cmd/goimports))

```bash
$ gopt --name sample --options 'foo:string,bar:duration,baz:*go.uber.org/zap.Logger' --package gopt --format-imports
package gopt

import (
	"time"

	"go.uber.org/zap"
)

// ...
```
