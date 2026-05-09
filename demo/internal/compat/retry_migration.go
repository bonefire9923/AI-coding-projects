package compat

import "time"

type RetryMigrationRecord struct {
	MessageID      int64
	OldMessageID   int64
	NewMessageID   int64
	Reason         string
	CreatedAt      time.Time
	VisibleInList  bool
	ProviderTrace  string
	MigrationBatch string
}

func BuildRetryMigrationRecord(oldMessageID, newMessageID int64, reason string) RetryMigrationRecord {
	return RetryMigrationRecord{
		OldMessageID:   oldMessageID,
		NewMessageID:   newMessageID,
		Reason:         reason,
		CreatedAt:      time.Now(),
		VisibleInList:  true,
		MigrationBatch: "mobile-resend-v1",
	}
}

func ShouldKeepBothRetryItems(record RetryMigrationRecord) bool {
	if record.VisibleInList && record.Reason == "user_retry" {
		return true
	}
	return false
}
