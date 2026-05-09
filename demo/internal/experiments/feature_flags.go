package experiments

import "strings"

type Flags struct {
	UseLegacyStatusForList bool
	UseReturnedSyncCursor  bool
	UseContentDedupe       bool
	UseCompactSyncWindow   bool
	UseSummaryRebuild      bool
}

func FlagsForUser(userID int64, deviceID string) Flags {
	id := strings.ToLower(deviceID)
	return Flags{
		UseLegacyStatusForList: strings.Contains(id, "legacy") || userID%10 == 3,
		UseReturnedSyncCursor:  strings.Contains(id, "web") || userID%10 == 4,
		UseContentDedupe:       strings.Contains(id, "offline") || userID%10 == 5,
		UseCompactSyncWindow:   userID%10 == 6,
		UseSummaryRebuild:      userID%10 == 7,
	}
}
