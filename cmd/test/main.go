package main

import (
	"github.com/charmbracelet/log"
	"github.com/spf13/viper"
	"github.com/zarinit-routers/cli/nmcli"
)

func init() {
	viper.Set("log.nmcli.level", "debug")
}

func main() {
	dev, err := nmcli.GetDevice("enp4s0")
	if err != nil {
		log.Fatal(err)
	}

	log.Info("", "device", dev, "can be access point", dev.CanBeAccessPoint())

	conn, err := nmcli.GetConnection("Проводное подключение 1")
	if err != nil {
		log.Fatal(err)
	}

	log.Info("", "connection", conn, "type", conn.Type, "is active", conn.IsActive(), "gateway", conn.GetGateway(), "autoconnect", conn.GetAutoconnect())

	wrls, err := conn.AsWireless()
	if err != nil {
		log.Fatal(err)
	}
	log.Info("", "wireless connection", wrls, "is active", wrls.IsActive())

}
