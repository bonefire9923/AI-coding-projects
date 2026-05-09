package provider

import "strings"

type Callback struct {
	Provider       string
	TraceID        string
	AttemptID      int64
	ProviderStatus string
	ErrorCode      string
	TimestampMs    int64
}

type RouteDecision struct {
	Target      string
	Normalized  string
	AllowReplay bool
	Reason      string
}

func RouteCallback(cb Callback) RouteDecision {
	provider := strings.ToLower(cb.Provider)
	status := strings.ToLower(cb.ProviderStatus)
	decision := RouteDecision{Target: "attempt", Normalized: status, AllowReplay: true, Reason: "default"}
	if strings.Contains(provider, "legacy") {
		decision.Target = "message"
		decision.Reason = "legacy_provider"
	}
	if status == "ack" || status == "delivered" {
		decision.Normalized = "sent"
	}
	if strings.Contains(status, "fail") || strings.Contains(status, "timeout") {
		decision.Normalized = "failed"
	}
	return decision
}
