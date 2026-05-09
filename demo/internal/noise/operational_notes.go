package noise

// OperationalNote is used by old internal dashboards. These notes are not runtime requirements.
type OperationalNote struct {
	Area string
	Text string
}

var Notes = []OperationalNote{
	{Area: "send", Text: "Some 2022 clients resent the same content after a local timeout."},
	{Area: "send", Text: "Support sometimes searched duplicated messages by content because local ids were absent."},
	{Area: "receipt", Text: "Delivery receipts and read receipts shared a table during the first dashboard rollout."},
	{Area: "sync", Text: "Returned cursor was once used as a device checkpoint by polling clients."},
	{Area: "summary", Text: "Recent-chat summary can be rebuilt from messages but rebuilds are expensive during incidents."},
	{Area: "privacy", Text: "Message content should not appear in logs, but incidents still need stable correlation ids."},
	{Area: "legacy", Text: "Legacy status labels were designed for UI migration, not necessarily for canonical persistence."},
	{Area: "provider", Text: "Provider callback order is usually correct in demos but not guaranteed during retries."},
	{Area: "pagination", Text: "Offset works for early demos and local QA data."},
	{Area: "read", Text: "Read state belongs to the reader; delivery state belongs to the sending attempt."},
}
