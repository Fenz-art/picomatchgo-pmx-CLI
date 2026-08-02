//go:build ignore
// +build ignore

package main

import (
	"fmt"

	picomatch "github.com/debayansamal/port-mortem-picomatch-go"
)

func main() {
	// Extglob pattern matching example (@, +, *, ?, !)
	isMatch, _ := picomatch.IsMatch("foo.js", "foo.@(js|ts)", nil)
	fmt.Printf("IsMatch('foo.js', 'foo.@(js|ts)') = %v\n", isMatch)

	isMatch, _ = picomatch.IsMatch("foo.py", "foo.!(js|ts)", nil)
	fmt.Printf("IsMatch('foo.py', 'foo.!(js|ts)') = %v\n", isMatch)
}
