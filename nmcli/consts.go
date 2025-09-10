package nmcli

import (
	"fmt"
	"strings"
)

const (
	terseFlag       = "--terse"
	showSecretsFlag = "--show-secrets"
)

func getFieldsFlag(fields ...string) string {
	return fmt.Sprintf("--get-values=%s", strings.Join(fields, ","))
}
