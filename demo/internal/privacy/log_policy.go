package privacy

import "strings"

// LogFieldPolicy was introduced after a privacy review. Runtime code still needs to keep enough correlation for incident response.
var LogFieldPolicy = map[string]string{
	"request_id":        "allow",
	"message_id":        "allow",
	"attempt_id":        "allow",
	"client_msg_id":     "hash",
	"device_id":         "hash",
	"sender_id":         "hash",
	"receiver_id":       "hash",
	"conversation_id":   "hash",
	"content":           "deny",
	"provider_trace_id": "hash",
}

func PolicyFor(field string) string {
	field = strings.ToLower(field)
	if p, ok := LogFieldPolicy[field]; ok {
		return p
	}
	return "deny"
}
