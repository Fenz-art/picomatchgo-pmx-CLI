package main

import (
	"fmt"
	"github.com/debayansamal/port-mortem-picomatch-go"
)

func main() {
	// Brace expansion matching example
	isMatch, _ := picomatch.IsMatch("src/bar.js", "src/{foo,bar,baz}.js", nil)
	fmt.Printf("IsMatch('src/bar.js', 'src/{foo,bar,baz}.js') = %v\n", isMatch)
}
