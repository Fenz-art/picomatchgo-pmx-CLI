//go:build ignore

package main

import (
	"fmt"

	picomatch "github.com/Fenz-art/picomatchgo-pmx-CLI"
)

func main() {
	// Custom options example: case-insensitive & dotfile matching
	opts := &picomatch.Options{
		Nocase: true,
		Dot:    true,
	}

	isMatch, _ := picomatch.IsMatch(".FOO.JS", "*.js", opts)
	fmt.Printf("IsMatch('.FOO.JS', '*.js', Nocase+Dot) = %v\n", isMatch)
}
