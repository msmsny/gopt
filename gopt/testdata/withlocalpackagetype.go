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
	baz *LocalPackageType
}

type fooOption string

func (o fooOption) apply(opts *sampleOptions) { opts.foo = string(o) }

func WithFoo(value string) SampleOption { return fooOption(value) }

type barOption int

func (o barOption) apply(opts *sampleOptions) { opts.bar = int(o) }

func WithBar(value int) SampleOption { return barOption(value) }

type localPackageTypeOption struct {
	localPackageType *LocalPackageType
}

func (o localPackageTypeOption) apply(opts *sampleOptions) { opts.baz = o.localPackageType }

func WithBaz(localPackageType *LocalPackageType) SampleOption {
	return localPackageTypeOption{localPackageType: localPackageType}
}
