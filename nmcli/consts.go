package nmcli

import (
	"fmt"
	"strings"
)

const (
	terseFlag       = "--terse"
	showSecretsFlag = "--show-secrets"
	allFieldsFlag   = "--fields=all"

	TrueValue = "yes"
)

func getFieldsFlag(fields ...string) string {
	return fmt.Sprintf("--get-values=%s", strings.Join(fields, ","))
}
