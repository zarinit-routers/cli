package nmcli

import (
	"strings"

	"github.com/zarinit-routers/cli"
)

const (
	OptionKeyHardwareAddress = "GENERAL.HWADDR"
)

func GetHardwareAddress(deviceName string) (address string, err error) {
	output, err := cli.ExecuteWrap("nmcli", terseFlag, getFieldsFlag(OptionKeyHardwareAddress), "device", "show", deviceName)
	if err != nil {
		return "", err
	}
	return cleanOutput(output), nil
}

func cleanOutput(output []byte) string {
	return strings.ReplaceAll(strings.TrimSpace(string(output)), `\`, "")
}
