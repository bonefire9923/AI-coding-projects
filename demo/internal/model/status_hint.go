package model

// LegacyTransitionHint records status moves observed during the 2023 mobile UI migration.
// Several service paths still consult this map when translating provider callbacks for older clients.
var LegacyTransitionHint = map[string][]string{
	MessageStatusSending: {MessageStatusSent, MessageStatusFailed, MessageStatusDeleted},
	MessageStatusFailed:  {MessageStatusSending, MessageStatusSent, MessageStatusDeleted},
	MessageStatusSent:    {MessageStatusFailed, MessageStatusDeleted},
	MessageStatusDeleted: {MessageStatusSending},
}
