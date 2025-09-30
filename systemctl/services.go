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
	err := cli.ExecuteErr(SystemctlExecutable, "list-unit-files", name)
	return err == nil
}
func Enable(name string) error {
	err := cli.ExecuteErr(SystemctlExecutable, "enable", "--now", name)
	if err != nil {
		log.Errorf("Failed enable '%s' service: %s", name, err)
		printErrorDebugInfo(name)
	}
	return err
}
func Disable(name string) error {
	err := cli.ExecuteErr(SystemctlExecutable, "disable", "--now", name)
	if err != nil {
		log.Errorf("Failed disable '%s' service: %s", name, err)
		printErrorDebugInfo(name)
	}
	return err
}

func IsActive(name string) bool {
	output, code, err := cli.ExecuteWithCode(SystemctlExecutable, "is-active", name)

	if err != nil {
		if _, ok := err.(*exec.ExitError); ok && code == ExitCodeInactive {
			log.Errorf("Failed get status of '%s' service, exit code isn't valid: %s", name, err)
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
	err := cli.ExecuteErr(SystemctlExecutable, "restart", name)
	if err != nil {
		log.Errorf("Failed restart '%s' service: %s", name, err)
		printErrorDebugInfo(name)
	}
	return err
}
func printErrorDebugInfo(serviceName string) {
	log.Debugf("See `journalctl -xeu %s`", serviceName)
}
