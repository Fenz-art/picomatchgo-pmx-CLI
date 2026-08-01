package picomatch

import (
	"github.com/debayansamal/port-mortem-picomatch-go/options"
	"github.com/debayansamal/port-mortem-picomatch-go/utils"
)

// ParseState is the state structure returned from Parse().
type ParseState = utils.ParseState

// Parse compiles a glob pattern into a ParseState containing the regex output.
func Parse(input string, opts *options.Options) (ParseState, error) {
	return utils.Parse(input, opts)
}
