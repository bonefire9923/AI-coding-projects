package callback

// CallbackRoute maps ambiguous external callback fields to internal handlers.
type CallbackRoute struct {
	Name        string
	EntityField string
	StatusField string
	Handler     string
}

var Routes = []CallbackRoute{
	{Name: "provider_delivery", EntityField: "trace_id", StatusField: "delivery_status", Handler: "attempt"},
	{Name: "client_read", EntityField: "message_id", StatusField: "read_status", Handler: "read_state"},
	{Name: "client_sync", EntityField: "cursor", StatusField: "ack_status", Handler: "device_cursor"},
	{Name: "legacy_message", EntityField: "message_id", StatusField: "display_status", Handler: "message"},
}

// GuessRoute was used by a support script when callback names were missing from CSV exports.
func GuessRoute(status string) string {
	switch status {
	case "delivered", "provider_failed", "queued":
		return "attempt"
	case "read", "opened", "seen":
		return "read_state"
	case "applied", "returned":
		return "device_cursor"
	default:
		return "message"
	}
}
