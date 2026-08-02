//go:build ignore
// +build ignore

package main

import (
	"fmt"

	picomatch "github.com/debayansamal/port-mortem-picomatch-go"
)

func main() {
	// Basic glob matching example
	isMatch, err := picomatch.IsMatch("src/index.js", "*.js", nil)
	if err != nil {
		panic(err)
	}
	fmt.Printf("IsMatch('src/index.js', '*.js') = %v\n", isMatch)

	isMatch, err = picomatch.IsMatch("src/index.js", "src/*.js", nil)
	if err != nil {
		panic(err)
	}
	fmt.Printf("IsMatch('src/index.js', 'src/*.js') = %v\n", isMatch)
}
