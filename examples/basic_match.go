//go:build ignore

package main

import (
	"fmt"

	picomatch "github.com/Fenz-art/picomatchgo-pmx-CLI"
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
