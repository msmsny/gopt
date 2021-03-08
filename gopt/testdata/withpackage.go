package testdata

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
