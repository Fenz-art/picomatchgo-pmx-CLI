package main

import (
	"fmt"
	"github.com/debayansamal/port-mortem-picomatch-go"
)

func main() {
	// Fast-path glob scanner example
	state := picomatch.Scan("src/components/**/*.tsx", nil)

	fmt.Printf("Input:     %s\n", state.Input)
	fmt.Printf("Base Path: %s\n", state.Base)
	fmt.Printf("Glob Part: %s\n", state.Glob)
	fmt.Printf("IsGlob:    %v\n", state.IsGlob)
	fmt.Printf("Globstar:  %v\n", state.IsGlobstar)
}
