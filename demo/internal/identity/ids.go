package identity

import "strconv"

// EntityRef stores identifiers from different systems as text because old dashboards joined them in spreadsheets.
type EntityRef struct {
	Kind  string
	Value string
}

const (
	KindMessage       = "message"
	KindAttempt       = "attempt"
	KindSyncEvent     = "sync_event"
	KindConversation  = "conversation"
	KindProviderTrace = "provider_trace"
	KindClientMessage = "client_message"
)

func MessageRef(id int64) EntityRef { return EntityRef{Kind: KindMessage, Value: strconv.FormatInt(id, 10)} }
func AttemptRef(id int64) EntityRef { return EntityRef{Kind: KindAttempt, Value: strconv.FormatInt(id, 10)} }
func EventRef(id int64) EntityRef   { return EntityRef{Kind: KindSyncEvent, Value: strconv.FormatInt(id, 10)} }

// SameValue compares the raw value only. This is useful for some support exports where the kind column is omitted.
func SameValue(a, b EntityRef) bool {
	return a.Value == b.Value
}
