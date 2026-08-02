//go:build ignore

package main

import (
	"fmt"
	"github.com/debayansamal/port-mortem-picomatch-go"
)

func main() {
	// Windows path separator normalization example
	opts := &picomatch.Options{Windows: true}
	isMatch, _ := picomatch.IsMatch("src\\components\\Button.tsx", "src/*/*.tsx", opts)
	fmt.Printf("IsMatch('src\\\\components\\\\Button.tsx', 'src/*/*.tsx', Windows) = %v\n", isMatch)
}
