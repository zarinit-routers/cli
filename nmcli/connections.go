package nmcli

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/zarinit-routers/cli"
)

type Connection struct {
	Name    string
	UUID    string
	Type    ConnectionType
	Device  string
	options map[string]string
}
type WirelessConnection struct {
	*Connection
}

type ConnectionType string

const (
	ConnectionTypeWIFI     ConnectionType = "wifi"
	ConnectionTypeWireless ConnectionType = "802-11-wireless" // Access point connection
	ConnectionTypeEthernet ConnectionType = "ethernet"
)

func GetConnections() ([]Connection, error) {
	output, err := cli.Execute("nmcli", terseFlag, "connection")
	if err != nil {
		return nil, err
	}

	connections := parseConnections(output)
	return connections, nil
}

func parseConnections(cliOutput []byte) []Connection {
	lines := strings.Split(
		string(cliOutput), "\n",
	)
	connections := []Connection{}
	for _, line := range lines {
		conn, err := parseConn(line)
		if err != nil {
			log.Warnf("Bad connection: %s", err)
			continue
		}
		connections = append(connections, conn)
	}
	return connections
}

var ErrTooLittleCols = fmt.Errorf("too little cols specified")

func parseConn(line string) (Connection, error) {
	words := strings.Split(line, ":")
	if len(words) < 4 {
		log.Debugf("Bad connection '%s'", line)
		return Connection{}, ErrTooLittleCols
	}
	return Connection{
		Name:   words[0],
		UUID:   words[1],
		Type:   ConnectionType(words[2]),
		Device: words[3],
	}, nil

}

func createConnection(
	t ConnectionType,
	deviceName string,
	connectionName string, additionalCliParams []string) (*Connection, error) {

	params := []string{"connection", "add", "type", string(t), "ifname", deviceName, "con-name", connectionName}
	params = append(params, additionalCliParams...)
	err := cli.ExecuteErr("nmcli", params...)
	if err != nil {
		return nil, err
	}
	return &Connection{
		Name:   connectionName,
		Type:   t,
		Device: deviceName,
	}, nil
}

const (
	OptionKeyIP4Method    = "ipv4.method"
	OptionKeyIP4Addresses = "ipv4.addresses"
	OptionKeyGeneralState = "GENERAL.STATE"
)

type IP4Method = string

const (
	ConnectionIP4MethodShared IP4Method = "shared"
)

func (c *Connection) SetIP4Method(method IP4Method) error {
	return c.setOption(OptionKeyIP4Method, string(method))
}
func (c *Connection) SetIP4Address(address string) error {
	return c.setOption(OptionKeyIP4Addresses, address)
}

func (c *Connection) Up() error {
	return cli.ExecuteErr("nmcli", "connection", "up", c.Name)
}
func (c *Connection) Down() error {
	return cli.ExecuteErr("nmcli", "connection", "down", c.Name)
}

type ConnectionState = string

const (
	ConnectionStateActivated = "activated"
)

func (c *Connection) IsActive() bool {
	state := c.getOption(OptionKeyGeneralState)
	return ConnectionState(state) == ConnectionStateActivated
}

func (c *Connection) getOption(optionName string) string {
	c.ensureOptionsParsed()
	return c.options[optionName]
}

func (c *Connection) ensureOptionsParsed() error {
	if c.options != nil {
		return nil
	}

	output, err := cli.ExecuteWrap("nmcli", terseFlag, showSecretsFlag, "connection", "show", c.Name)
	if err != nil {
		return err
	}
	c.options = parseKeyValOutput(output)
	return nil
}

func (c *Connection) setOption(optionName, optionValue string) error {
	log.Debugf("Setting option %s to '%s', current value is '%s'", optionName, optionValue, c.options[optionName])
	return cli.ExecuteErr("nmcli", "connection", "modify", c.Name, optionName, optionValue)
}

func GetConnection(name string) (*Connection, error) {
	output, err := cli.ExecuteWrap("nmcli", terseFlag, showSecretsFlag, "connection", "show", name)
	if err != nil {
		return nil, err
	}
	return parseShowConnectionOutput(output), nil
}

func parseShowConnectionOutput(output []byte) *Connection {
	dict := parseKeyValOutput(output)
	return &Connection{
		Name:    dict["connection.id"],
		UUID:    dict["connection.uuid"],
		Type:    ConnectionType(dict["connection.type"]),
		Device:  dict["connection.interface-name"],
		options: dict,
	}
}

func parseKeyValOutput(output []byte) map[string]string {
	dict := map[string]string{}
	lines := strings.Split(string(output), "\n")
	for _, l := range lines {
		words := strings.Split(l, ":")
		if len(words) < 2 {
			dict[words[0]] = ""
			continue
		}

		dict[words[0]] = words[1]
	}
	return dict
}
