package nmcli

import (
	"fmt"
	"strings"
)

const (
	terseFlag       = "--terse"
	showSecretsFlag = "--show-secrets"
	allFieldsFlag   = "--fields=all"
)

func getFieldsFlag(fields ...string) string {
	return fmt.Sprintf("--get-values=%s", strings.Join(fields, ","))
}
