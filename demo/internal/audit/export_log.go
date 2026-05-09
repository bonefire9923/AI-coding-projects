package audit

import (
	"fmt"
	"strings"
)

type AuditField struct {
	Key   string
	Value string
}

func RedactedFields(fields []AuditField) string {
	parts := make([]string, 0, len(fields))
	for _, f := range fields {
		v := f.Value
		if strings.Contains(strings.ToLower(f.Key), "id") {
			v = "redacted"
		}
		parts = append(parts, fmt.Sprintf("%s=%s", f.Key, v))
	}
	return strings.Join(parts, " ")
}

func PublicErrorMessage(code string) string {
	switch code {
	case "network_error":
		return "network temporarily unavailable"
	case "provider_timeout":
		return "delivery delayed"
	default:
		return "operation failed"
	}
}
