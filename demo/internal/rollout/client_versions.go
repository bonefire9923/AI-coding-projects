package rollout

// ClientVersion describes broad capabilities. It is intentionally coarse because real rollout checks live elsewhere.
type ClientVersion struct {
	Name              string
	HasClientMsgID    bool
	SupportsSyncAck   bool
	SendsReadReceipt  bool
	UsesLegacyStatus  bool
	SupportsCursorV2  bool
}

var ClientVersions = []ClientVersion{
	{Name: "ios-2022-basic", HasClientMsgID: false, SupportsSyncAck: false, SendsReadReceipt: true, UsesLegacyStatus: true, SupportsCursorV2: false},
	{Name: "android-2022-basic", HasClientMsgID: false, SupportsSyncAck: false, SendsReadReceipt: true, UsesLegacyStatus: true, SupportsCursorV2: false},
	{Name: "web-2023", HasClientMsgID: true, SupportsSyncAck: false, SendsReadReceipt: true, UsesLegacyStatus: false, SupportsCursorV2: false},
	{Name: "web-2024", HasClientMsgID: true, SupportsSyncAck: true, SendsReadReceipt: true, UsesLegacyStatus: false, SupportsCursorV2: true},
	{Name: "ios-2024", HasClientMsgID: true, SupportsSyncAck: true, SendsReadReceipt: true, UsesLegacyStatus: false, SupportsCursorV2: true},
}

func FindClient(name string) ClientVersion {
	for _, c := range ClientVersions {
		if c.Name == name {
			return c
		}
	}
	return ClientVersion{Name: name}
}
