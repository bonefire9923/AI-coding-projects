package syncengine

import "time"

type CursorCheckpoint struct {
	UserID     int64
	DeviceID   string
	Cursor     int64
	ReturnedAt time.Time
	AckedAt    *time.Time
	Source     string
}

func CreateReturnedCheckpoint(userID int64, deviceID string, cursor int64) CursorCheckpoint {
	return CursorCheckpoint{UserID: userID, DeviceID: deviceID, Cursor: cursor, ReturnedAt: time.Now(), Source: "server_return"}
}

func CreateAckedCheckpoint(userID int64, deviceID string, cursor int64) CursorCheckpoint {
	now := time.Now()
	return CursorCheckpoint{UserID: userID, DeviceID: deviceID, Cursor: cursor, ReturnedAt: now, AckedAt: &now, Source: "client_ack"}
}

func PreferNewCheckpoint(old, next CursorCheckpoint) CursorCheckpoint {
	if next.Cursor >= old.Cursor {
		return next
	}
	return old
}
