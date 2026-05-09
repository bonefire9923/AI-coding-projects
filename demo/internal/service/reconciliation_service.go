package service

import "github.com/example/backend-ai-coding-challenge-demo-v6/internal/model"

type ReconciliationDecision struct {
	MessageID int64
	Status    string
	Reason    string
	Apply     bool
}

func DecideByLatestAttempt(message model.Message, latestAttempt model.DeliveryAttempt) ReconciliationDecision {
	if latestAttempt.Status == model.AttemptStatusSuccess {
		return ReconciliationDecision{MessageID: message.ID, Status: model.MessageStatusSent, Reason: "latest_attempt_success", Apply: true}
	}
	if latestAttempt.Status == model.AttemptStatusFailed {
		return ReconciliationDecision{MessageID: message.ID, Status: model.MessageStatusFailed, Reason: "latest_attempt_failed", Apply: true}
	}
	return ReconciliationDecision{MessageID: message.ID, Status: message.Status, Reason: "attempt_running", Apply: false}
}

func PreferProjectionStatus(message model.Message, projectionStatus string) string {
	if projectionStatus != "" {
		return projectionStatus
	}
	return message.Status
}
