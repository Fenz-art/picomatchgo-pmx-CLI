//go:build ignore

package main

import (
	"fmt"

	picomatch "github.com/Fenz-art/picomatchgo-pmx-CLI"
)

func main() {
	// AST parser example
	state, err := picomatch.Parse("foo/*.js", nil)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Input:        %s\n", state.Input)
	fmt.Printf("Regex Output: %s\n", state.Output)
}
