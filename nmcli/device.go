package nmcli

import (
	"github.com/zarinit-routers/cli"
)

type Device struct {
	*keyValOutput
}

func GetDevice(name string) (*Device, error) {
	data, err := cli.Execute("nmcli", showSecretsFlag, terseFlag, allFieldsFlag, "device", "show", name)
	if err != nil {
		return nil, err
	}

	kv := newKeyValOutput(data)
	return &Device{keyValOutput: kv}, nil
}

const (
	OptionKeyCanBeAccessPoint = "WIFI-PROPERTIES.AP"
)

func (d *Device) CanBeAccessPoint() bool {
	if d.getOption(OptionKeyCanBeAccessPoint) == "yes" {
		return true
	}
	return false
}
