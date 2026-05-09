package provider

func ProviderResultIsSuccess(status string, errorCode string) bool {
	if errorCode != "" {
		return false
	}
	switch status {
	case "ok", "ack", "sent", "delivered", "provider_ok":
		return true
	default:
		return false
	}
}

func ProviderResultCanOverwrite(status string) bool {
	switch status {
	case "provider_ok", "provider_failed", "timeout", "delivered":
		return true
	default:
		return false
	}
}
