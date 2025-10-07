// Documentation for wireless nmcli:
//
// - https://www.networkmanager.dev/docs/api/latest/settings-802-11-wireless.html
//
// - https://www.networkmanager.dev/docs/api/latest/settings-802-11-wireless-security.html
package nmcli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/zarinit-routers/cli"
)

const (
	OptionKeyWirelessSSID                  = "802-11-wireless.ssid"
	OptionKeyWirelessHidden                = "802-11-wireless.hidden"
	OptionKeyWirelessChanel                = "802-11-wireless.channel"
	OptionKeyWirelessMode                  = "802-11-wireless.mode"
	OptionKeyWirelessBand                  = "802-11-wireless.band"
	OptionKeyWirelessSecurityPassword      = "802-11-wireless-security.psk"
	OptionKeyWirelessSecurityKeyManagement = "802-11-wireless-security.key-mgmt"
	OptionKeyWirelessSecurityProto         = "802-11-wireless-security.proto"
	OptionKeyWirelessSecurityGroup         = "802-11-wireless-security.group"
	OptionKeyWirelessSecurityPairwise      = "802-11-wireless-security.pairwise"
	OptionKeyWirelessSeenBSSIDs            = "802-11-wireless.seen-bssids"
)

func (c *Connection) AsWireless() (*WirelessConnection, error) {
	if c.Type != ConnectionTypeWireless {
		return nil, fmt.Errorf("connection %q is not a wireless connection but a %s connection", c.Name, c.Type)
	}
	return &WirelessConnection{c}, nil
}
func CreateWirelessConnection(deviceName string, connectionName string, password string) (*WirelessConnection, error) {
	if len(password) < 8 {
		return nil, fmt.Errorf("invalid password: must be at least 8 characters long")
	}

	dev, err := GetDevice(deviceName)
	if err != nil {
		return nil, fmt.Errorf("can't get device %q: %s", deviceName, err)
	}
	if !dev.CanBeAccessPoint() {
		return nil, fmt.Errorf("device %q can't be access point", deviceName)
	}

	conn, err := createConnection(
		ConnectionTypeWIFI, deviceName, connectionName,
		[]string{
			"autoconnect", TrueValue,
			"ssid", connectionName,
			OptionKeyWirelessSecurityPassword, password,
			OptionKeyWirelessSecurityKeyManagement, KeyManagementWPA2_3Personal,
			OptionKeyWirelessSecurityProto, ProtoAllowWPA2RSN,
			OptionKeyWirelessSecurityGroup, EncryptionAlgCcmp,
			OptionKeyWirelessSecurityPairwise, EncryptionAlgCcmp,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed create base connection: %s", err)
	}
	wireless := WirelessConnection{conn}
	if _, err := cli.Execute("nmcli", "device", "wifi", "hotspot",
		"ifname", deviceName,
		"conn-name", conn.Name,
		"ssid", wireless.GetSSID(),
		"password", wireless.GetPassword(),
	); err != nil {
		return nil, fmt.Errorf("can't create hotspot: %s", err)
	}

	return &wireless, nil
}

const (
	EncryptionAlgTkip   = "tkip"
	EncryptionAlgCcmp   = "ccmp"
	EncryptionAlgWep40  = "wep40"
	EncryptionAlgWep104 = "wep104"
)

type Proto = string

const (
	ProtoAllowWPA2RSN Proto = "rsn"
	ProtoAllowWPA     Proto = "wpa"
)

type KeyManagement = string

const (
	KeyManagementNone             KeyManagement = "none"
	KeyManagementWPA2_3Personal   KeyManagement = "wpa-psk"             // WPA2 + WPA3 personal
	KeyManagementWPA3Personal     KeyManagement = "sae"                 // WPA3 personal only
	KeyManagementWPA2_3Enterprise KeyManagement = "wpa-eap"             // WPA2 + WPA3 enterprise
	KeyManagementWPA3Enterprise   KeyManagement = "wpa-eap-suite-b-192" // WPA3 enterprise only
)

type WirelessMode string

const (
	WirelessModeAccessPoint    WirelessMode = "ap"
	WirelessModeInfrastructure WirelessMode = "infrastructure"
	WirelessModeMesh           WirelessMode = "mesh"
	WirelessModeAdhoc          WirelessMode = "adhoc"
)

func (c *WirelessConnection) SetMode(mode WirelessMode) error {
	return c.setOption(OptionKeyWirelessMode, string(mode))
}

type WirelessBand = string

const (
	WirelessBand2GHz WirelessBand = "bg"
	WirelessBand5GHz WirelessBand = "a"
)

func (c *WirelessConnection) SetBand(band WirelessBand) error {
	return c.setOption(OptionKeyWirelessBand, string(band))
}
func (c *WirelessConnection) GetBand() WirelessBand {
	return WirelessBand(c.getOption(OptionKeyWirelessBand))
}

func (c *WirelessConnection) GetSSID() string {
	return c.getOption(OptionKeyWirelessSSID)
}
func (c *WirelessConnection) SetSSID(ssid string) error {
	err := c.setOption(OptionKeyWirelessSSID, ssid)
	if err == nil {
		return c.ensureOptionsParsed()
	}
	return err
}
func (c *WirelessConnection) GetChanel() int {
	opt := c.getOption(OptionKeyWirelessChanel)
	value, err := strconv.Atoi(opt)
	if err != nil {
		return 0
	}
	return value
}
func (c *WirelessConnection) SetChannel(chanel int) error {
	return c.setOption(OptionKeyWirelessChanel, strconv.Itoa(chanel))
}
func (c *WirelessConnection) GetPassword() string {
	return c.getOption(OptionKeyWirelessSecurityPassword)
}
func (c *WirelessConnection) SetPassword(password string) error {
	return c.setOption(OptionKeyWirelessSecurityPassword, password)
}

const (
	WirelessHiddenValue    = TrueValue
	WirelessNotHiddenValue = "no"
)

func (c *WirelessConnection) IsHidden() bool {
	return c.getOption(OptionKeyWirelessHidden) == WirelessHiddenValue
}
func (c *WirelessConnection) SetHidden(hide bool) error {
	var value string
	if hide {
		value = WirelessHiddenValue
	} else {
		value = WirelessNotHiddenValue
	}
	err := c.setOption(OptionKeyWirelessHidden, value)
	if err == nil {
		return c.ensureOptionsParsed()
	}
	return err
}

func (c *WirelessConnection) GetBSSID() string {
	seen := c.getOption(OptionKeyWirelessSeenBSSIDs)
	firstBssid := strings.Split(seen, ",")[0]
	return firstBssid
}

type DeviceDataKey string

const (
	DeviceDataKeySSID           DeviceDataKey = "SSID"
	DeviceDataKeyBSSID          DeviceDataKey = "BSSID"
	DeviceDataKeySignalStrength DeviceDataKey = "SIGNAL"
	DeviceDataKeyRate           DeviceDataKey = "RATE"
)

func (c *WirelessConnection) GetSignalStrength() uint {
	bssid := c.GetBSSID()

	val, err := c.getDeviceData(DeviceDataKeySignalStrength)

	if err != nil {
		log.Errorf("Failed get wifi signal strength for BSSID '%s': %s", bssid, err)
		return 0
	}

	strength, err := strconv.Atoi(string(val))
	if err != nil {
		log.Errorf("Failed parse wifi signal strength from '%s': %s", string(val), err)
	}
	return uint(strength)
}

func (c *WirelessConnection) getDeviceData(key DeviceDataKey) (string, error) {
	bssid := c.GetBSSID()

	val, err := cli.Execute(
		"nmcli", terseFlag, getFieldsFlag(string(key)),
		"device", "wifi", "list",
		"bssid", bssid,
	)
	if err != nil {
		log.Error("Failed get device data for for connection", "connectionName", c.Name, "connectionBSSID", bssid, "err", err)
		return "", err
	}
	return string(val), nil
}

func (c *WirelessConnection) GetNetworkRate() string {
	bssid := c.GetBSSID()

	val, err := c.getDeviceData(DeviceDataKeyRate)

	if err != nil {
		log.Errorf("Failed get wifi signal strength for BSSID '%s': %s", bssid, err)
	}

	return string(val)
}
