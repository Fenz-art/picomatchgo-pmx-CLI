package javascript

import (
	"os/exec"
)

var execLookPathImpl = exec.LookPath

func execLookPath(name string) (string, error) { return execLookPathImpl(name) }
