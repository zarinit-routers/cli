package nmcli

import (
	"fmt"
	"net"
	"strings"

	"github.com/zarinit-routers/cli"
)

type Connection struct {
	*keyValOutput

	Name   string
	UUID   string
	Type   ConnectionType
	Device string
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
	return GetConnection(connectionName)
}

const (
	OptionKeyAutoconnect   = "connection.autoconnect"
	OptionKeyIP4Method     = "ipv4.method"
	OptionKeyIP4Addresses  = "ipv4.addresses"
	OptionKeyGeneralState  = "GENERAL.STATE"
	OptionKeyDNSAddresses  = "ipv4.dns"
	OptionKeyDHCPRange     = "ipv4.dhcp-range"
	OptionKeyDHCPLeaseTime = "ipv4.dhcp-lease-time"
	OptionKeyIP4Gateway    = "ipv4.gateway"
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

// TODO: move to net.IP
func (c *Connection) SetDNSAddresses(addresses []string) error {
	return c.setOption(OptionKeyDNSAddresses, strings.Join(addresses, ","))
}

// Deprecated: THis method must not be used
func (c *Connection) SetDHCPRange(from, to net.IP) error {
	return c.setOption(OptionKeyDHCPRange, strings.Join(
		[]string{from.String(), to.String()}, ","))
}

// Deprecated: THis method must not be used
func (c *Connection) SetDHCPLeaseTime(secs int) error {
	return c.setOption(OptionKeyDHCPLeaseTime, fmt.Sprintf("%d", secs))
}

func (c *Connection) GetGateway() net.IP {
	gateway := c.getOption(OptionKeyIP4Gateway)
	return net.ParseIP(gateway)
}
func (c *Connection) SetGateway(gateway net.IP) error {
	return c.setOption(OptionKeyIP4Gateway, gateway.String())
}
func (c *Connection) GetAutoconnect() bool {
	opt := c.getOption(OptionKeyAutoconnect)
	return opt == TrueValue
}

type ConnectionState = string

const (
	ConnectionStateActivated = "activated"
)

func (c *Connection) IsActive() bool {
	state := c.getOption(OptionKeyGeneralState)
	return ConnectionState(state) == ConnectionStateActivated
}

func (c *Connection) setOption(optionName, optionValue string) error {
	log.Debug("Setting option", "option", optionName, "newValue", optionValue, "currentValue", c.options[optionName])
	err := cli.ExecuteErr("nmcli", "connection", "modify", c.Name, optionName, optionValue)
	if err != nil {
		return fmt.Errorf("failed set option %q to %q: %s", optionName, optionValue, err)
	}

	c.options[optionName] = optionValue
	return nil
}

func GetConnection(name string) (*Connection, error) {
	output, err := cli.Execute("nmcli", allFieldsFlag, terseFlag, showSecretsFlag, "connection", "show", fmt.Sprintf("%q", name))
	if err != nil {
		return nil, fmt.Errorf("failed execute NetworkManger: %s", err)
	}
	return parseShowConnectionOutput(output), nil
}

func parseShowConnectionOutput(output []byte) *Connection {
	kv := newKeyValOutput(output)
	if err := kv.ensureOptionsParsed(); err != nil {
		log.Errorf("Failed to parse connection output: %s", err)
	}

	return &Connection{
		keyValOutput: kv,

		Name:   kv.getOption("connection.id"),
		UUID:   kv.getOption("connection.uuid"),
		Type:   ConnectionType(kv.getOption("connection.type")),
		Device: kv.getOption("connection.interface-name"),
	}
}
