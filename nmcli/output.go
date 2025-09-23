package nmcli

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/log"
)

type keyValOutput struct {
	output  []byte
	options map[string]string
}

func newKeyValOutput(output []byte) *keyValOutput {
	return &keyValOutput{
		output: output,
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

func (c *keyValOutput) ensureOptionsParsed() error {
	if c.options != nil {
		return nil
	}

	if c.output == nil {
		return fmt.Errorf("no output to parse options from")
	}

	c.options = parseKeyValOutput(c.output)
	return nil
}

func (c *keyValOutput) getOption(optionName string) string {
	if err := c.ensureOptionsParsed(); err != nil {
		log.Error("Failed to parse options", "error", err)
		return ""
	}
	return c.options[optionName]
}
