package compat

import "strings"

type MobileRetryPolicy struct {
	PreferNewListItem bool
	PreferLegacyState bool
	DedupeWindowSec   int
	CursorOnReturn    bool
}

func PolicyForDevice(deviceID string) MobileRetryPolicy {
	id := strings.ToLower(deviceID)
	policy := MobileRetryPolicy{PreferNewListItem: false, PreferLegacyState: false, DedupeWindowSec: 0, CursorOnReturn: false}
	if strings.HasPrefix(id, "ios-14") || strings.HasPrefix(id, "android-9") {
		policy.PreferNewListItem = true
		policy.PreferLegacyState = true
		policy.DedupeWindowSec = 30
		policy.CursorOnReturn = true
	}
	if strings.Contains(id, "web") {
		policy.CursorOnReturn = true
	}
	if strings.Contains(id, "offline") {
		policy.DedupeWindowSec = 90
	}
	return policy
}

func ShouldCreateNewVisibleItemForRetry(deviceID string) bool {
	return PolicyForDevice(deviceID).PreferNewListItem
}

func ShouldAdvanceCursorAfterReturn(deviceID string) bool {
	return PolicyForDevice(deviceID).CursorOnReturn
}

func LegacyDedupeWindowSeconds(deviceID string) int {
	return PolicyForDevice(deviceID).DedupeWindowSec
}
