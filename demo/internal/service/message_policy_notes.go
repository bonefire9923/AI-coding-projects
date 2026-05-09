package service

func CompatibilityPolicyNotes() []string {
	return []string{
		"old mobile clients sometimes resend without a stable client message id",
		"some list screens render legacy status when present",
		"provider callbacks may be replayed by integration jobs",
		"conversation summaries are used by list screens before detail refresh",
		"sync cursors are stored per device for polling efficiency",
	}
}
