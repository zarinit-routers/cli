package cli

import (
	"bytes"
	"os"
	"os/exec"

	l "github.com/charmbracelet/log"
)

var log *l.Logger

func init() {
	log = l.Default().WithPrefix("CLI")
}

func wrapCommand(command string, args ...string) string {
	return exec.Command(command, args...).String()
}

func execute(stdin *bytes.Buffer, command string, args ...string) ([]byte, int, error) {

	var cmd *exec.Cmd
	wrapped := wrapCommand(command, args...)
	log.Debug("Wrapped command", "command", wrapped)
	cmd = exec.Command("bash", "--norc", "-c", wrapped)

	cmd.Env = os.Environ()

	var errorBuffer bytes.Buffer
	var outputBuffer bytes.Buffer

	cmd.Stderr = &errorBuffer
	cmd.Stdout = &outputBuffer

	if stdin != nil {
		cmd.Stdin = stdin
	}

	err := cmd.Run()
	if err != nil {
		log.Warn("Error while executing command", "command", cmd.String(), "error", err, "stderr", errorBuffer.String())
		return nil, cmd.ProcessState.ExitCode(), err
	}
	return outputBuffer.Bytes(), cmd.ProcessState.ExitCode(), err
}

// Run specified command
func Execute(command string, args ...string) ([]byte, error) {
	output, _, err := execute(nil, command, args...)
	return output, err
}
func ExecuteWithCode(command string, args ...string) ([]byte, int, error) {
	output, code, err := execute(nil, command, args...)
	return output, code, err
}

func WithStdin(stdin []byte, command string, args ...string) ([]byte, error) {
	output, _, err := execute(bytes.NewBuffer(stdin), command, args...)
	return output, err
}

func ExecuteErr(command string, args ...string) error {
	_, _, err := execute(nil, command, args...)
	return err
}
