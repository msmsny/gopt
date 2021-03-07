package main

import (
	"fmt"
	"os"

	"github.com/msmsny/gopt/gopt"
)

func main() {
	if err := gopt.NewGoptCommand().Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}
}
