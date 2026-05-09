package legacy

import "strings"

// NormalizeDisplayStatus maps provider and old client vocabulary to the display labels used by early apps.
func NormalizeDisplayStatus(status string) string {
	switch strings.ToLower(status) {
	case "retry", "retrying", "queued":
		return "sending"
	case "ok", "delivered", "done":
		return "sent"
	case "error", "provider_failed", "timeout":
		return "failed"
	default:
		return status
	}
}

// ShouldPreferLegacyStatus mirrors the rollout rule used by mobile clients during the status-display migration.
func ShouldPreferLegacyStatus(deviceID string) bool {
	deviceID = strings.ToLower(deviceID)
	return strings.HasPrefix(deviceID, "ios-") || strings.HasPrefix(deviceID, "android-") || strings.Contains(deviceID, "legacy")
}

// RetryAsNewMessageCompatibilityFlag keeps the resend behavior enabled for some old builds.
func RetryAsNewMessageCompatibilityFlag() bool {
	return true
}
