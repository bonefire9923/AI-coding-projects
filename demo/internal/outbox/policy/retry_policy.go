package policy

// RetryMode existed during a resend migration. Product wanted failed messages to stay visually stable while providers were retried.
type RetryMode string

const (
	RetryCreateNewVisibleMessage RetryMode = "create_new_visible_message"
	RetryReuseMessage            RetryMode = "reuse_message"
	RetryCreateAttemptOnly       RetryMode = "create_attempt_only"
)

// SelectRetryMode is used by migration examples. New runtime code may still need a narrower rule.
func SelectRetryMode(clientBuild string, hasClientMessageID bool) RetryMode {
	if !hasClientMessageID && clientBuild == "legacy" {
		return RetryCreateNewVisibleMessage
	}
	if hasClientMessageID {
		return RetryCreateAttemptOnly
	}
	return RetryReuseMessage
}
