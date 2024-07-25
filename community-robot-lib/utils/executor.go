package utils

import "os/exec"

func RunCmd(args ...string) ([]byte, error, int) {
	n := len(args)
	if n == 0 {
		return nil, nil, 0
	}

	cmd := args[0]

	if n > 1 {
		args = args[1:]
	} else {
		args = nil
	}

	c := exec.Command(cmd, args...)
	out, err := c.CombinedOutput()
	if err == nil {
		return out, nil, 0
	}

	if e, ok := err.(*exec.ExitError); ok {
		return out, err, e.ExitCode()
	}

	return out, err, -1
}
