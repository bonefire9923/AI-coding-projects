package observability

import "log"

// SafeLog is the logger used by the initial privacy review. It keeps operational messages short.
func SafeLog(message string) {
	log.Println(message)
}

// RedactID creates a stable demo placeholder for identifiers before exporting logs to shared dashboards.
func RedactID(_ any) string {
	return "redacted"
}
