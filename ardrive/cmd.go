package ardrive

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"sync"
)

var globalLock sync.Mutex

// ExecCmd executes a command and returns the combined output and error.
func ExecCmd(ctx context.Context, cmd string, args ...string) ([]byte, error) {
	var combinedOutput bytes.Buffer
	command := exec.CommandContext(ctx, cmd, args...)
	command.Stdout = &combinedOutput
	command.Stderr = &combinedOutput

	err := command.Run()
	if err != nil {
		// Combine stdout and stderr for non-zero exit codes
		return nil, fmt.Errorf("command failed: %w\n%s", err, combinedOutput.String())
	}

	return combinedOutput.Bytes(), nil
}
