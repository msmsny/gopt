// +build test

package main

import (
	"fmt"
	"net/http"
	"time"
)

// This main command is a test code which creates test instance and prints.
//
// In order to build command, generate options
// $ gopt --name sample --options 'foo:string,bar:int,baz:bool,qux:duration,quux:*net/http.Client' --package main -o test/cmd/buildtest/options.go
//
// Build test/cmd/buildtest and execute
// $ go build -tags test ./test/cmd/buildtest
// $ ./buildtest
// &main.test{required:"required", foo:"foo", bar:1, baz:true, qux:1000000000, quux:(*http.Client)(nil)}, quux: &http.Client{Transport:http.RoundTripper(nil), CheckRedirect:(func(*http.Request, []*http.Request) error)(nil), Jar:http.CookieJar(nil), Timeout:0}
func main() {
	test := newTest(
		"required",
		WithFoo("foo"),
		WithBar(1),
		WithBaz(true),
		WithQux(1000000000),
		WithQuux(&http.Client{}),
	)
	// dump quux directly for pointer address
	quux := test.quux
	test.quux = nil
	fmt.Printf("%#v, quux: %#v\n", test, quux)
}

type test struct {
	required string
	foo      string
	bar      int
	baz      bool
	qux      time.Duration
	quux     *http.Client
}

func newTest(required string, options ...SampleOption) *test {
	opts := evaluateOptions(options)

	return &test{
		required: required,
		foo:      opts.foo,
		bar:      opts.bar,
		baz:      opts.baz,
		qux:      opts.qux,
		quux:     opts.quux,
	}
}
