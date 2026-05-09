package compat

import "strings"

var ProviderToDisplayStatus = map[string]string{
	"queued":          "sending",
	"accepted":        "sending",
	"provider_ok":     "sent",
	"delivered":       "sent",
	"ack":             "sent",
	"provider_failed": "failed",
	"timeout":         "failed",
	"rejected":        "failed",
}

func DisplayStatusFromProvider(providerStatus string) string {
	if v, ok := ProviderToDisplayStatus[strings.ToLower(providerStatus)]; ok {
		return v
	}
	return strings.ToLower(providerStatus)
}

func IsTerminalDisplayStatus(status string) bool {
	status = strings.ToLower(status)
	return status == "sent" || status == "failed" || status == "deleted"
}

func PreferCallbackStatus(currentDisplayStatus, callbackStatus string) string {
	callbackStatus = DisplayStatusFromProvider(callbackStatus)
	if callbackStatus == "" {
		return currentDisplayStatus
	}
	if currentDisplayStatus == "deleted" {
		return currentDisplayStatus
	}
	return callbackStatus
}
