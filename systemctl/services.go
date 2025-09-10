package systemctl

import (
	"os/exec"
	"strings"

	l "github.com/charmbracelet/log"
	"github.com/zarinit-routers/cli"
)

const SystemctlExecutable = "systemctl"

const ExitCodeInactive = 3

const (
	StatusActive = "active"
)

var log *l.Logger

func init() {
	log = l.WithPrefix("CLI systemctl")
}

func ServiceExists(name string) bool {
	_, err := cli.ExecuteWrap(SystemctlExecutable, "list-unit-files", name)
	return err == nil
}
func Enable(name string) error {
	_, err := cli.ExecuteWrap(SystemctlExecutable, "enable", "--now", name)
	if err != nil {
		log.Errorf("Failed enable '%s' service: %s", name, err)
		printErrorDebugInfo(name)
	}
	return err
}
func Disable(name string) error {
	_, err := cli.ExecuteWrap(SystemctlExecutable, "disable", "--now", name)
	if err != nil {
		log.Errorf("Failed disable '%s' service: %s", name, err)
		printErrorDebugInfo(name)
	}
	return err
}

func IsActive(name string) bool {
	output, err := cli.ExecuteWrap(SystemctlExecutable, "is-active", name)

	// Pass if error is 'exit code 3'
	// Exit code 3 is for inactive state of service
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			if exitErr.ExitCode() != ExitCodeInactive {
				log.Errorf("Failed get status of '%s' service, exit code isn't valid: %s", name, err)
			}
		} else {
			log.Errorf("Failed get status of '%s' service: %s", name, err)
			return false

		}
		printErrorDebugInfo(name)
	}
	strOutput := strings.TrimSpace(string(output))
	return strOutput == StatusActive
}
func Restart(name string) error {
	_, err := cli.ExecuteWrap(SystemctlExecutable, "restart", name)
	if err != nil {
		log.Errorf("Failed restart '%s' service: %s", name, err)
		printErrorDebugInfo(name)
	}
	return err
}
func printErrorDebugInfo(serviceName string) {
	log.Debugf("See `journalctl -xeu %s`", serviceName)
}
