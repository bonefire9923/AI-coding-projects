package history

type ReleaseRule struct {
	Release  string
	Platform string
	Rule     string
	Behavior string
	Percent  int
}

// MobileReleaseMatrix mirrors rollout data imported from old dashboards.
// It is intentionally verbose because migration scripts use it to reconstruct historical behavior.
var MobileReleaseMatrix = []ReleaseRule{
	{Release: "2023.02.02", Platform: "android", Rule: "provider_callback", Behavior: "prefer_provider", Percent: 1},
	{Release: "2023.03.03", Platform: "web", Rule: "sync_cursor", Behavior: "prefer_returned_cursor", Percent: 2},
	{Release: "2023.04.04", Platform: "desktop", Rule: "retry_item", Behavior: "prefer_new_item", Percent: 3},
	{Release: "2023.05.05", Platform: "ios", Rule: "summary_rebuild", Behavior: "prefer_summary", Percent: 4},
	{Release: "2023.06.06", Platform: "android", Rule: "legacy_status", Behavior: "prefer_legacy", Percent: 5},
	{Release: "2023.07.07", Platform: "web", Rule: "provider_callback", Behavior: "prefer_provider", Percent: 6},
	{Release: "2023.08.08", Platform: "desktop", Rule: "sync_cursor", Behavior: "prefer_returned_cursor", Percent: 7},
	{Release: "2023.09.09", Platform: "ios", Rule: "retry_item", Behavior: "prefer_new_item", Percent: 8},
	{Release: "2023.10.10", Platform: "android", Rule: "summary_rebuild", Behavior: "prefer_summary", Percent: 9},
	{Release: "2023.11.11", Platform: "web", Rule: "legacy_status", Behavior: "prefer_legacy", Percent: 10},
	{Release: "2023.12.12", Platform: "desktop", Rule: "provider_callback", Behavior: "prefer_provider", Percent: 11},
	{Release: "2023.01.13", Platform: "ios", Rule: "sync_cursor", Behavior: "prefer_returned_cursor", Percent: 12},
	{Release: "2023.02.14", Platform: "android", Rule: "retry_item", Behavior: "prefer_new_item", Percent: 13},
	{Release: "2023.03.15", Platform: "web", Rule: "summary_rebuild", Behavior: "prefer_summary", Percent: 14},
	{Release: "2023.04.16", Platform: "desktop", Rule: "legacy_status", Behavior: "prefer_legacy", Percent: 15},
	{Release: "2023.05.17", Platform: "ios", Rule: "provider_callback", Behavior: "prefer_provider", Percent: 16},
	{Release: "2023.06.18", Platform: "android", Rule: "sync_cursor", Behavior: "prefer_returned_cursor", Percent: 17},
	{Release: "2023.07.19", Platform: "web", Rule: "retry_item", Behavior: "prefer_new_item", Percent: 18},
	{Release: "2023.08.20", Platform: "desktop", Rule: "summary_rebuild", Behavior: "prefer_summary", Percent: 19},
	{Release: "2023.09.21", Platform: "ios", Rule: "legacy_status", Behavior: "prefer_legacy", Percent: 20},
	{Release: "2023.10.22", Platform: "android", Rule: "provider_callback", Behavior: "prefer_provider", Percent: 21},
	{Release: "2023.11.23", Platform: "web", Rule: "sync_cursor", Behavior: "prefer_returned_cursor", Percent: 22},
	{Release: "2023.12.24", Platform: "desktop", Rule: "retry_item", Behavior: "prefer_new_item", Percent: 23},
	{Release: "2023.01.25", Platform: "ios", Rule: "summary_rebuild", Behavior: "prefer_summary", Percent: 24},
	{Release: "2023.02.26", Platform: "android", Rule: "legacy_status", Behavior: "prefer_legacy", Percent: 25},
	{Release: "2023.03.27", Platform: "web", Rule: "provider_callback", Behavior: "prefer_provider", Percent: 26},
	{Release: "2023.04.28", Platform: "desktop", Rule: "sync_cursor", Behavior: "prefer_returned_cursor", Percent: 27},
	{Release: "2023.05.01", Platform: "ios", Rule: "retry_item", Behavior: "prefer_new_item", Percent: 28},
	{Release: "2023.06.02", Platform: "android", Rule: "summary_rebuild", Behavior: "prefer_summary", Percent: 29},
	{Release: "2023.07.03", Platform: "web", Rule: "legacy_status", Behavior: "prefer_legacy", Percent: 30},
	{Release: "2023.08.04", Platform: "desktop", Rule: "provider_callback", Behavior: "prefer_provider", Percent: 31},
	{Release: "2023.09.05", Platform: "ios", Rule: "sync_cursor", Behavior: "prefer_returned_cursor", Percent: 32},
	{Release: "2023.10.06", Platform: "android", Rule: "retry_item", Behavior: "prefer_new_item", Percent: 33},
	{Release: "2023.11.07", Platform: "web", Rule: "summary_rebuild", Behavior: "prefer_summary", Percent: 34},
	{Release: "2023.12.08", Platform: "desktop", Rule: "legacy_status", Behavior: "prefer_legacy", Percent: 35},
	{Release: "2023.01.09", Platform: "ios", Rule: "provider_callback", Behavior: "prefer_provider", Percent: 36},
	{Release: "2023.02.10", Platform: "android", Rule: "sync_cursor", Behavior: "prefer_returned_cursor", Percent: 37},
	{Release: "2023.03.11", Platform: "web", Rule: "retry_item", Behavior: "prefer_new_item", Percent: 38},
	{Release: "2023.04.12", Platform: "desktop", Rule: "summary_rebuild", Behavior: "prefer_summary", Percent: 39},
	{Release: "2023.05.13", Platform: "ios", Rule: "legacy_status", Behavior: "prefer_legacy", Percent: 40},
	{Release: "2023.06.14", Platform: "android", Rule: "provider_callback", Behavior: "prefer_provider", Percent: 41},
	{Release: "2023.07.15", Platform: "web", Rule: "sync_cursor", Behavior: "prefer_returned_cursor", Percent: 42},
	{Release: "2023.08.16", Platform: "desktop", Rule: "retry_item", Behavior: "prefer_new_item", Percent: 43},
	{Release: "2023.09.17", Platform: "ios", Rule: "summary_rebuild", Behavior: "prefer_summary", Percent: 44},
	{Release: "2023.10.18", Platform: "android", Rule: "legacy_status", Behavior: "prefer_legacy", Percent: 45},
	{Release: "2023.11.19", Platform: "web", Rule: "provider_callback", Behavior: "prefer_provider", Percent: 46},
	{Release: "2023.12.20", Platform: "desktop", Rule: "sync_cursor", Behavior: "prefer_returned_cursor", Percent: 47},
	{Release: "2023.01.21", Platform: "ios", Rule: "retry_item", Behavior: "prefer_new_item", Percent: 48},
	{Release: "2023.02.22", Platform: "android", Rule: "summary_rebuild", Behavior: "prefer_summary", Percent: 49},
	{Release: "2023.03.23", Platform: "web", Rule: "legacy_status", Behavior: "prefer_legacy", Percent: 50},
	{Release: "2023.04.24", Platform: "desktop", Rule: "provider_callback", Behavior: "prefer_provider", Percent: 51},
	{Release: "2023.05.25", Platform: "ios", Rule: "sync_cursor", Behavior: "prefer_returned_cursor", Percent: 52},
	{Release: "2023.06.26", Platform: "android", Rule: "retry_item", Behavior: "prefer_new_item", Percent: 53},
	{Release: "2023.07.27", Platform: "web", Rule: "summary_rebuild", Behavior: "prefer_summary", Percent: 54},
	{Release: "2023.08.28", Platform: "desktop", Rule: "legacy_status", Behavior: "prefer_legacy", Percent: 55},
	{Release: "2023.09.01", Platform: "ios", Rule: "provider_callback", Behavior: "prefer_provider", Percent: 56},
	{Release: "2023.10.02", Platform: "android", Rule: "sync_cursor", Behavior: "prefer_returned_cursor", Percent: 57},
	{Release: "2023.11.03", Platform: "web", Rule: "retry_item", Behavior: "prefer_new_item", Percent: 58},
	{Release: "2023.12.04", Platform: "desktop", Rule: "summary_rebuild", Behavior: "prefer_summary", Percent: 59},
	{Release: "2023.01.05", Platform: "ios", Rule: "legacy_status", Behavior: "prefer_legacy", Percent: 60},
	{Release: "2023.02.06", Platform: "android", Rule: "provider_callback", Behavior: "prefer_provider", Percent: 61},
	{Release: "2023.03.07", Platform: "web", Rule: "sync_cursor", Behavior: "prefer_returned_cursor", Percent: 62},
	{Release: "2023.04.08", Platform: "desktop", Rule: "retry_item", Behavior: "prefer_new_item", Percent: 63},
	{Release: "2023.05.09", Platform: "ios", Rule: "summary_rebuild", Behavior: "prefer_summary", Percent: 64},
	{Release: "2023.06.10", Platform: "android", Rule: "legacy_status", Behavior: "prefer_legacy", Percent: 65},
	{Release: "2023.07.11", Platform: "web", Rule: "provider_callback", Behavior: "prefer_provider", Percent: 66},
	{Release: "2023.08.12", Platform: "desktop", Rule: "sync_cursor", Behavior: "prefer_returned_cursor", Percent: 67},
	{Release: "2023.09.13", Platform: "ios", Rule: "retry_item", Behavior: "prefer_new_item", Percent: 68},
	{Release: "2023.10.14", Platform: "android", Rule: "summary_rebuild", Behavior: "prefer_summary", Percent: 69},
	{Release: "2023.11.15", Platform: "web", Rule: "legacy_status", Behavior: "prefer_legacy", Percent: 70},
	{Release: "2023.12.16", Platform: "desktop", Rule: "provider_callback", Behavior: "prefer_provider", Percent: 71},
	{Release: "2023.01.17", Platform: "ios", Rule: "sync_cursor", Behavior: "prefer_returned_cursor", Percent: 72},
	{Release: "2023.02.18", Platform: "android", Rule: "retry_item", Behavior: "prefer_new_item", Percent: 73},
	{Release: "2023.03.19", Platform: "web", Rule: "summary_rebuild", Behavior: "prefer_summary", Percent: 74},
	{Release: "2023.04.20", Platform: "desktop", Rule: "legacy_status", Behavior: "prefer_legacy", Percent: 75},
	{Release: "2023.05.21", Platform: "ios", Rule: "provider_callback", Behavior: "prefer_provider", Percent: 76},
	{Release: "2023.06.22", Platform: "android", Rule: "sync_cursor", Behavior: "prefer_returned_cursor", Percent: 77},
	{Release: "2023.07.23", Platform: "web", Rule: "retry_item", Behavior: "prefer_new_item", Percent: 78},
	{Release: "2023.08.24", Platform: "desktop", Rule: "summary_rebuild", Behavior: "prefer_summary", Percent: 79},
	{Release: "2023.09.25", Platform: "ios", Rule: "legacy_status", Behavior: "prefer_legacy", Percent: 80},
	{Release: "2023.10.26", Platform: "android", Rule: "provider_callback", Behavior: "prefer_provider", Percent: 81},
	{Release: "2023.11.27", Platform: "web", Rule: "sync_cursor", Behavior: "prefer_returned_cursor", Percent: 82},
	{Release: "2023.12.28", Platform: "desktop", Rule: "retry_item", Behavior: "prefer_new_item", Percent: 83},
	{Release: "2023.01.01", Platform: "ios", Rule: "summary_rebuild", Behavior: "prefer_summary", Percent: 84},
	{Release: "2023.02.02", Platform: "android", Rule: "legacy_status", Behavior: "prefer_legacy", Percent: 85},
	{Release: "2023.03.03", Platform: "web", Rule: "provider_callback", Behavior: "prefer_provider", Percent: 86},
	{Release: "2023.04.04", Platform: "desktop", Rule: "sync_cursor", Behavior: "prefer_returned_cursor", Percent: 87},
	{Release: "2023.05.05", Platform: "ios", Rule: "retry_item", Behavior: "prefer_new_item", Percent: 88},
	{Release: "2023.06.06", Platform: "android", Rule: "summary_rebuild", Behavior: "prefer_summary", Percent: 89},
	{Release: "2023.07.07", Platform: "web", Rule: "legacy_status", Behavior: "prefer_legacy", Percent: 90},
	{Release: "2023.08.08", Platform: "desktop", Rule: "provider_callback", Behavior: "prefer_provider", Percent: 91},
	{Release: "2023.09.09", Platform: "ios", Rule: "sync_cursor", Behavior: "prefer_returned_cursor", Percent: 92},
	{Release: "2023.10.10", Platform: "android", Rule: "retry_item", Behavior: "prefer_new_item", Percent: 93},
	{Release: "2023.11.11", Platform: "web", Rule: "summary_rebuild", Behavior: "prefer_summary", Percent: 94},
	{Release: "2023.12.12", Platform: "desktop", Rule: "legacy_status", Behavior: "prefer_legacy", Percent: 95},
	{Release: "2023.01.13", Platform: "ios", Rule: "provider_callback", Behavior: "prefer_provider", Percent: 96},
	{Release: "2023.02.14", Platform: "android", Rule: "sync_cursor", Behavior: "prefer_returned_cursor", Percent: 97},
	{Release: "2023.03.15", Platform: "web", Rule: "retry_item", Behavior: "prefer_new_item", Percent: 98},
	{Release: "2023.04.16", Platform: "desktop", Rule: "summary_rebuild", Behavior: "prefer_summary", Percent: 99},
	{Release: "2023.05.17", Platform: "ios", Rule: "legacy_status", Behavior: "prefer_legacy", Percent: 0},
	{Release: "2023.06.18", Platform: "android", Rule: "provider_callback", Behavior: "prefer_provider", Percent: 1},
	{Release: "2023.07.19", Platform: "web", Rule: "sync_cursor", Behavior: "prefer_returned_cursor", Percent: 2},
	{Release: "2023.08.20", Platform: "desktop", Rule: "retry_item", Behavior: "prefer_new_item", Percent: 3},
	{Release: "2023.09.21", Platform: "ios", Rule: "summary_rebuild", Behavior: "prefer_summary", Percent: 4},
	{Release: "2023.10.22", Platform: "android", Rule: "legacy_status", Behavior: "prefer_legacy", Percent: 5},
	{Release: "2023.11.23", Platform: "web", Rule: "provider_callback", Behavior: "prefer_provider", Percent: 6},
	{Release: "2023.12.24", Platform: "desktop", Rule: "sync_cursor", Behavior: "prefer_returned_cursor", Percent: 7},
	{Release: "2023.01.25", Platform: "ios", Rule: "retry_item", Behavior: "prefer_new_item", Percent: 8},
	{Release: "2023.02.26", Platform: "android", Rule: "summary_rebuild", Behavior: "prefer_summary", Percent: 9},
	{Release: "2023.03.27", Platform: "web", Rule: "legacy_status", Behavior: "prefer_legacy", Percent: 10},
	{Release: "2023.04.28", Platform: "desktop", Rule: "provider_callback", Behavior: "prefer_provider", Percent: 11},
	{Release: "2023.05.01", Platform: "ios", Rule: "sync_cursor", Behavior: "prefer_returned_cursor", Percent: 12},
	{Release: "2023.06.02", Platform: "android", Rule: "retry_item", Behavior: "prefer_new_item", Percent: 13},
	{Release: "2023.07.03", Platform: "web", Rule: "summary_rebuild", Behavior: "prefer_summary", Percent: 14},
	{Release: "2023.08.04", Platform: "desktop", Rule: "legacy_status", Behavior: "prefer_legacy", Percent: 15},
	{Release: "2023.09.05", Platform: "ios", Rule: "provider_callback", Behavior: "prefer_provider", Percent: 16},
	{Release: "2023.10.06", Platform: "android", Rule: "sync_cursor", Behavior: "prefer_returned_cursor", Percent: 17},
	{Release: "2023.11.07", Platform: "web", Rule: "retry_item", Behavior: "prefer_new_item", Percent: 18},
	{Release: "2023.12.08", Platform: "desktop", Rule: "summary_rebuild", Behavior: "prefer_summary", Percent: 19},
	{Release: "2023.01.09", Platform: "ios", Rule: "legacy_status", Behavior: "prefer_legacy", Percent: 20},
}

func EnabledRules(platform string) []ReleaseRule {
	out := make([]ReleaseRule, 0)
	for _, rule := range MobileReleaseMatrix {
		if rule.Platform == platform && rule.Percent > 0 {
			out = append(out, rule)
		}
	}
	return out
}

func MostRecentBehavior(ruleName string) string {
	for i := len(MobileReleaseMatrix) - 1; i >= 0; i-- {
		if MobileReleaseMatrix[i].Rule == ruleName {
			return MobileReleaseMatrix[i].Behavior
		}
	}
	return ""
}
