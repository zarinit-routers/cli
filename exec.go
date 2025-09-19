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

func execute(stdin *bytes.Buffer, command string, args ...string) ([]byte, error) {

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
	return execute(nil, command, args...)
}

// Run specified command
func Execute(command string, args ...string) ([]byte, error) {
	return execute(nil, command, args...)
}

func WithStdin(stdin []byte, command string, args ...string) ([]byte, error) {
	return execute(bytes.NewBuffer(stdin), command, args...)
}

func ExecuteErr(command string, args ...string) error {
	_, err := execute(nil, command, args...)
	return err
}
