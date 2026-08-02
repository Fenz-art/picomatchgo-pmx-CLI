//go:build ignore

package main

import (
	"fmt"
	"github.com/debayansamal/port-mortem-picomatch-go"
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
