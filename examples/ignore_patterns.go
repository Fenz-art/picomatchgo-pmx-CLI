//go:build ignore
// +build ignore

package main

import (
	"fmt"

	picomatch "github.com/debayansamal/port-mortem-picomatch-go"
)

func main() {
	// Ignore patterns option example
	opts := &picomatch.Options{
		Ignore: []string{"*.min.js", "vendor/**"},
	}

	isMatch, _ := picomatch.IsMatch("bundle.min.js", "*.js", opts)
	fmt.Printf("IsMatch('bundle.min.js', '*.js', Ignore) = %v\n", isMatch)
}
