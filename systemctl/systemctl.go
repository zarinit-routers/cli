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
	err := cli.ExecuteErr(SystemctlExecutable, "list-unit-files", s.name)
	return err == nil
}
func Enable(s Service) error {
	err := cli.ExecuteErr(SystemctlExecutable, "enable", "--now", s.name)
	if err != nil {
		log.Error("Failed enable service", "service", s.name, "error", err)
		printErrorDebugInfo(s)
	}
	return err
}
func EnableForUser(s Service) error {
	err := cli.ExecuteErr(SystemctlExecutable, "enable", "--now", "--user", s.name)
	if err != nil {
		log.Error("Failed enable service", "service", s.name, "error", err)
		printErrorDebugInfo(s)
	}
	return err
}
func Disable(s Service) error {
	err := cli.ExecuteErr(SystemctlExecutable, "disable", "--now", s.name)
	if err != nil {
		log.Error("Failed disable service", "service", s.name, "error", err)
		printErrorDebugInfo(s)
	}
	return err
}
func DisableForUser(s Service) error {
	err := cli.ExecuteErr(SystemctlExecutable, "disable", "--now", "--user", s.name)
	if err != nil {
		log.Error("Failed disable service", "service", s.name, "error", err)
		printErrorDebugInfo(s)
	}
	return err
}

func IsActive(s Service) bool {
	output, code, err := cli.ExecuteWithCode(SystemctlExecutable, "is-active", s.name)

	if err != nil {
		if _, ok := err.(*exec.ExitError); ok && code == ExitCodeInactive {
			log.Errorf("Failed get status of '%s' service, exit code isn't valid: %s", s.name, err)
		} else {
			log.Errorf("Failed get status of '%s' service: %s", s.name, err)
			return false
		}
		printErrorDebugInfo(s)
	}
	strOutput := strings.TrimSpace(string(output))
	return strOutput == StatusActive
}
func Restart(s Service) error {
	err := cli.ExecuteErr(SystemctlExecutable, "restart", s.name)
	if err != nil {
		log.Errorf("Failed restart '%s' service: %s", s.name, err)
		printErrorDebugInfo(s)
	}
	return err
}
func printErrorDebugInfo(s Service) {
	log.Debugf("See `journalctl -xeu %s`", s.name)
}
