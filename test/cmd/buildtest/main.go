// +build test

package main

import (
	"fmt"
	"time"
)

// This main command is a test code which creates test instance and prints.
//
// In order to build command, generate options
// $ gopt --name sample --options foo:string,bar:int,baz:bool,qux:duration --package main -o test/cmd/buildtest/options.go
//
// Build test/cmd/buildtest and execute
// $ go build -tags test ./test/cmd/buildtest
// $ ./buildtest
// &main.test{required:"required", foo:"foo", bar:1, baz:true, qux:1000000000}
func main() {
	test := newTest(
		"required",
		WithFoo("foo"),
		WithBar(1),
		WithBaz(true),
		WithQux(1000000000),
	)
	fmt.Printf("%#v\n", test)
}

type test struct {
	required string
	foo      string
	bar      int
	baz      bool
	qux      time.Duration
}

func newTest(required string, options ...SampleOption) *test {
	opts := evaluateOptions(options)

	return &test{
		required: required,
		foo:      opts.foo,
		bar:      opts.bar,
		baz:      opts.baz,
		qux:      opts.qux,
	}
}
