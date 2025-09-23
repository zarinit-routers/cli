package nmcli

import (
	l "github.com/charmbracelet/log"
)

var log *l.Logger

func init() {
	log = l.WithPrefix("NMCLI")
	log.SetLevel(l.DebugLevel)

}
