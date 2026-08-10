//go:build ignore

package main

import (
	"fmt"

	picomatch "github.com/Fenz-art/picomatchgo-pmx-CLI"
)

func main() {
	// Brace expansion matching example
	isMatch, _ := picomatch.IsMatch("src/bar.js", "src/{foo,bar,baz}.js", nil)
	fmt.Printf("IsMatch('src/bar.js', 'src/{foo,bar,baz}.js') = %v\n", isMatch)
}
