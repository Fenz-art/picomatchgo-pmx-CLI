package javascript

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

func execLookPath(name string) (string, error) { return exec.LookPath(name) }

func execCommand(name string, args ...string) (string, error) {
	out, err := exec.Command(name, args...).CombinedOutput()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func compareNodeVersion(actual, constraint string) (bool, error) {
	actualParts := strings.Split(actual, ".")
	if len(actualParts) == 0 {
		return false, fmt.Errorf("invalid node version: %q", actual)
	}
	major, err := strconv.Atoi(actualParts[0])
	if err != nil {
		return false, err
	}
	constraint = strings.TrimSpace(constraint)
	if constraint == "" {
		return true, nil
	}
	if strings.HasPrefix(constraint, ">=") {
		min, err := strconv.Atoi(strings.TrimPrefix(constraint, ">="))
		if err != nil {
			return false, err
		}
		return major >= min, nil
	}
	if strings.HasPrefix(constraint, ">") {
		min, err := strconv.Atoi(strings.TrimPrefix(constraint, ">"))
		if err != nil {
			return false, err
		}
		return major > min, nil
	}
	if strings.HasPrefix(constraint, "<=") {
		max, err := strconv.Atoi(strings.TrimPrefix(constraint, "<="))
		if err != nil {
			return false, err
		}
		return major <= max, nil
	}
	if strings.HasPrefix(constraint, "<") {
		max, err := strconv.Atoi(strings.TrimPrefix(constraint, "<"))
		if err != nil {
			return false, err
		}
		return major < max, nil
	}
	return major == 0 || constraint == "*", nil
}
