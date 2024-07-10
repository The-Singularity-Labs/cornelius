package sync

import (
	"bytes"
	"fmt"
	"os/exec"
	gosync "sync"
)

var globalLock gosync.Mutex

// ExecCmd executes a command and returns the combined output and error.
func ExecCmd(cmd string, args ...string) ([]byte, error) {
	var combinedOutput bytes.Buffer
	command := exec.Command(cmd, args...)
	command.Stdout = &combinedOutput
	command.Stderr = &combinedOutput

	err := command.Run()
	if err != nil {
		// Combine stdout and stderr for non-zero exit codes
		return nil, fmt.Errorf("command failed: %w\n%s", err, combinedOutput.String())
	}

	return combinedOutput.Bytes(), nil
}
