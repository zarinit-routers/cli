package cli

import (
	"bytes"
	"os/exec"

	"github.com/charmbracelet/log"
)

var logger *log.Logger

func init() {
	logger = log.Default().WithPrefix("CLI")
}

func wrapCommand(command string, args ...string) string {
	return exec.Command(command, args...).String()
}

func execute(wrapInBash bool, stdin *bytes.Buffer, command string, args ...string) ([]byte, error) {

	var cmd *exec.Cmd
	if wrapInBash {
		wrapped := wrapCommand(command, args...)
		logger.Debugf("Wrapped command to execute `%s`", wrapped)
		cmd = exec.Command("bash", "-c", wrapped)
	} else {
		cmd = exec.Command(command, args...)
	}
	logger.Debugf("Command to execute `%s`", cmd.String())

	var errorBuffer bytes.Buffer
	var outputBuffer bytes.Buffer

	cmd.Stderr = &errorBuffer
	cmd.Stdout = &outputBuffer

	if stdin != nil {
		cmd.Stdin = stdin
	}

	err := cmd.Run()
	if err != nil {
		logger.Warn("Error while executing command", "command", cmd.String(), "error", err, "stderr", errorBuffer.String())
		return nil, err
	}
	return outputBuffer.Bytes(), err

}

// Runs specified command in a bash shell
//
// Deprecated: Execute wrap should be removed later.
// Use simple Execute function instead.
// But firstly check if it works the same.
func ExecuteWrap(command string, args ...string) ([]byte, error) {
	return execute(true, nil, command, args...)
}

// Run specified command
func Execute(command string, args ...string) ([]byte, error) {
	return execute(false, nil, command, args...)
}

func WithStdin(stdin []byte, command string, args ...string) ([]byte, error) {
	return execute(false, bytes.NewBuffer(stdin), command, args...)
}

func ExecuteErr(command string, args ...string) error {
	_, err := ExecuteWrap(command, args...)
	return err
}
