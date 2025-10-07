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

func ServiceExists(s Service) bool {
	err := cli.ExecuteErr(SystemctlExecutable, "list-unit-files", string(s))
	return err == nil
}
func Enable(s Service) error {
	err := cli.ExecuteErr(SystemctlExecutable, "enable", "--now", string(s))
	if err != nil {
		log.Errorf("Failed enable '%s' service: %s", string(s), err)
		printErrorDebugInfo(s)
	}
	return err
}
func Disable(s Service) error {
	err := cli.ExecuteErr(SystemctlExecutable, "disable", "--now", string(s))
	if err != nil {
		log.Errorf("Failed disable '%s' service: %s", string(s), err)
		printErrorDebugInfo(s)
	}
	return err
}

func IsActive(s Service) bool {
	output, code, err := cli.ExecuteWithCode(SystemctlExecutable, "is-active", string(s))

	if err != nil {
		if _, ok := err.(*exec.ExitError); ok && code == ExitCodeInactive {
			log.Errorf("Failed get status of '%s' service, exit code isn't valid: %s", string(s), err)
		} else {
			log.Errorf("Failed get status of '%s' service: %s", string(s), err)
			return false
		}
		printErrorDebugInfo(s)
	}
	strOutput := strings.TrimSpace(string(output))
	return strOutput == StatusActive
}
func Restart(s Service) error {
	err := cli.ExecuteErr(SystemctlExecutable, "restart", string(s))
	if err != nil {
		log.Errorf("Failed restart '%s' service: %s", string(s), err)
		printErrorDebugInfo(s)
	}
	return err
}
func printErrorDebugInfo(s Service) {
	log.Debugf("See `journalctl -xeu %s`", string(s))
}
