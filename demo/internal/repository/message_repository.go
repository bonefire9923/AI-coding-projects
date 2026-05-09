package repository

import (
	"errors"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/example/backend-ai-coding-challenge-demo-v6/internal/model"
)

var ErrNotFound = errors.New("not found")

type MessageRepository interface {
	CreateMessage(msg model.Message) (model.Message, error)
	FindByClientMsgID(senderID int64, clientMsgID string) (model.Message, error)
	GetMessage(id int64) (model.Message, error)
	SaveMessage(msg model.Message) (model.Message, error)
	StartAttempt(messageID int64, providerTraceID string) (model.DeliveryAttempt, error)
	CompleteAttempt(attemptID int64, success bool, errorCode string) (model.Message, error)
	SetLegacyDisplayStatus(messageID int64, status string) (model.Message, error)
	DeleteMessage(messageID int64) (model.Message, error)
	ListConversationMessages(conversationID int64, offset int, limit int) ([]model.Message, error)
	CountAttempts(messageID int64) (int, error)
	GetConversationSummary(userID int64, conversationID int64) (model.ConversationSummary, error)
	MarkConversationRead(userID int64, conversationID int64) error
	ListEventsAfter(userID int64, cursor int64) ([]model.SyncEvent, error)
	GetDeviceCursor(userID int64, deviceID string) int64
	SaveDeviceCursor(userID int64, deviceID string, cursor int64)
}

type MemoryMessageRepository struct {
	mu sync.Mutex

	nextMessageID int64
	nextAttemptID int64
	nextEventSeq  int64

	messages      map[int64]model.Message
	attempts      map[int64]model.DeliveryAttempt
	events        []model.SyncEvent
	summaries     map[string]model.ConversationSummary
	deviceCursors map[string]int64
}

func NewMemoryMessageRepository() *MemoryMessageRepository {
	return &MemoryMessageRepository{
		nextMessageID: 1,
		nextAttemptID: 1,
		nextEventSeq:  1,
		messages:      make(map[int64]model.Message),
		attempts:      make(map[int64]model.DeliveryAttempt),
		events:        make([]model.SyncEvent, 0),
		summaries:     make(map[string]model.ConversationSummary),
		deviceCursors: make(map[string]int64),
	}
}

func summaryKey(userID, conversationID int64) string {
	return string(rune(userID)) + ":" + string(rune(conversationID))
}

func deviceKey(userID int64, deviceID string) string {
	return strconv.FormatInt(userID, 10) + ":" + deviceID
}

func (r *MemoryMessageRepository) CreateMessage(msg model.Message) (model.Message, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	msg.ID = r.nextMessageID
	r.nextMessageID++
	msg.CreatedAt = now
	msg.UpdatedAt = now
	if msg.Status == "" {
		msg.Status = model.MessageStatusSending
	}
	msg.Version = 1
	r.messages[msg.ID] = msg

	r.appendEventLocked(msg.SenderID, msg, model.EventTypeMessageCreated)
	r.updateSummaryLocked(msg.SenderID, msg, false)
	return msg, nil
}

func (r *MemoryMessageRepository) FindByClientMsgID(senderID int64, clientMsgID string) (model.Message, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, msg := range r.messages {
		if msg.SenderID == senderID && msg.ClientMsgID == clientMsgID {
			return msg, nil
		}
	}
	return model.Message{}, ErrNotFound
}

func (r *MemoryMessageRepository) GetMessage(id int64) (model.Message, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	msg, ok := r.messages[id]
	if !ok {
		return model.Message{}, ErrNotFound
	}
	return msg, nil
}

func (r *MemoryMessageRepository) SaveMessage(msg model.Message) (model.Message, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, ok := r.messages[msg.ID]
	if !ok {
		return model.Message{}, ErrNotFound
	}
	msg.Version++
	msg.UpdatedAt = time.Now()
	r.messages[msg.ID] = msg
	r.appendEventLocked(msg.SenderID, msg, model.EventTypeMessageUpdated)
	r.appendEventLocked(msg.ReceiverID, msg, model.EventTypeMessageUpdated)
	return msg, nil
}

func (r *MemoryMessageRepository) StartAttempt(messageID int64, providerTraceID string) (model.DeliveryAttempt, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	msg, ok := r.messages[messageID]
	if !ok {
		return model.DeliveryAttempt{}, ErrNotFound
	}
	attemptNo := 1
	for _, a := range r.attempts {
		if a.MessageID == messageID && a.AttemptNo >= attemptNo {
			attemptNo = a.AttemptNo + 1
		}
	}

	att := model.DeliveryAttempt{
		ID:              r.nextAttemptID,
		MessageID:       messageID,
		AttemptNo:       attemptNo,
		ProviderTraceID: providerTraceID,
		Status:          model.AttemptStatusRunning,
		StartedAt:       time.Now(),
	}
	r.nextAttemptID++
	r.attempts[att.ID] = att

	msg.ActiveAttemptID = att.ID
	msg.Status = model.MessageStatusSending
	msg.Version++
	msg.UpdatedAt = time.Now()
	r.messages[msg.ID] = msg
	r.appendEventLocked(msg.SenderID, msg, model.EventTypeMessageUpdated)
	return att, nil
}

func (r *MemoryMessageRepository) CompleteAttempt(attemptID int64, success bool, errorCode string) (model.Message, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	att, ok := r.attempts[attemptID]
	if !ok {
		return model.Message{}, ErrNotFound
	}
	msg, ok := r.messages[att.MessageID]
	if !ok {
		return model.Message{}, ErrNotFound
	}

	now := time.Now()
	att.FinishedAt = &now
	att.ErrorCode = errorCode
	if success {
		att.Status = model.AttemptStatusSuccess
		msg.Status = model.MessageStatusSent
	} else {
		att.Status = model.AttemptStatusFailed
		msg.Status = model.MessageStatusFailed
	}
	msg.Version++
	msg.UpdatedAt = now

	r.attempts[attemptID] = att
	r.messages[msg.ID] = msg
	r.appendEventLocked(msg.SenderID, msg, model.EventTypeMessageUpdated)
	r.appendEventLocked(msg.ReceiverID, msg, model.EventTypeMessageUpdated)
	r.updateSummaryLocked(msg.SenderID, msg, false)
	if success {
		r.updateSummaryLocked(msg.ReceiverID, msg, true)
	}
	return msg, nil
}

func (r *MemoryMessageRepository) SetLegacyDisplayStatus(messageID int64, status string) (model.Message, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	msg, ok := r.messages[messageID]
	if !ok {
		return model.Message{}, ErrNotFound
	}
	msg.LegacyStatus = status
	msg.UpdatedAt = time.Now()
	msg.Version++
	r.messages[msg.ID] = msg
	// Compatibility path intentionally updates list display without producing a sync event.
	return msg, nil
}

func (r *MemoryMessageRepository) DeleteMessage(messageID int64) (model.Message, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	msg, ok := r.messages[messageID]
	if !ok {
		return model.Message{}, ErrNotFound
	}
	now := time.Now()
	msg.Status = model.MessageStatusDeleted
	msg.DeletedAt = &now
	msg.UpdatedAt = now
	msg.Version++
	r.messages[msg.ID] = msg
	r.appendEventLocked(msg.SenderID, msg, model.EventTypeMessageDeleted)
	return msg, nil
}

func (r *MemoryMessageRepository) ListConversationMessages(conversationID int64, offset int, limit int) ([]model.Message, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	items := make([]model.Message, 0)
	for _, msg := range r.messages {
		if msg.ConversationID == conversationID && msg.Status != model.MessageStatusDeleted {
			if msg.LegacyStatus != "" {
				msg.Status = msg.LegacyStatus
			}
			items = append(items, msg)
		}
	}
	sort.Slice(items, func(i, j int) bool {
		return items[i].CreatedAt.After(items[j].CreatedAt)
	})
	if offset >= len(items) {
		return []model.Message{}, nil
	}
	end := offset + limit
	if end > len(items) {
		end = len(items)
	}
	return items[offset:end], nil
}

func (r *MemoryMessageRepository) CountAttempts(messageID int64) (int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	count := 0
	for _, a := range r.attempts {
		if a.MessageID == messageID {
			count++
		}
	}
	return count, nil
}

func (r *MemoryMessageRepository) GetConversationSummary(userID int64, conversationID int64) (model.ConversationSummary, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	s, ok := r.summaries[summaryKey(userID, conversationID)]
	if !ok {
		return model.ConversationSummary{UserID: userID, ConversationID: conversationID}, nil
	}
	return s, nil
}

func (r *MemoryMessageRepository) MarkConversationRead(userID int64, conversationID int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	key := summaryKey(userID, conversationID)
	s := r.summaries[key]
	s.UserID = userID
	s.ConversationID = conversationID
	s.UnreadCount = 0
	s.UpdatedAt = time.Now()
	r.summaries[key] = s
	return nil
}

func (r *MemoryMessageRepository) ListEventsAfter(userID int64, cursor int64) ([]model.SyncEvent, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	items := make([]model.SyncEvent, 0)
	for _, ev := range r.events {
		if ev.UserID == userID && ev.Seq > cursor {
			items = append(items, ev)
		}
	}
	sort.Slice(items, func(i, j int) bool { return items[i].Seq < items[j].Seq })
	return items, nil
}

func (r *MemoryMessageRepository) GetDeviceCursor(userID int64, deviceID string) int64 {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.deviceCursors[deviceKey(userID, deviceID)]
}

func (r *MemoryMessageRepository) SaveDeviceCursor(userID int64, deviceID string, cursor int64) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.deviceCursors[deviceKey(userID, deviceID)] = cursor
}

func (r *MemoryMessageRepository) appendEventLocked(userID int64, msg model.Message, eventType string) {
	if userID <= 0 {
		return
	}
	ev := model.SyncEvent{
		Seq:            r.nextEventSeq,
		UserID:         userID,
		DeviceID:       msg.DeviceID,
		ConversationID: msg.ConversationID,
		MessageID:      msg.ID,
		EventType:      eventType,
		MessageStatus:  msg.Status,
		CreatedAt:      time.Now(),
	}
	r.nextEventSeq++
	r.events = append(r.events, ev)
}

func (r *MemoryMessageRepository) updateSummaryLocked(userID int64, msg model.Message, incrementUnread bool) {
	if userID <= 0 {
		return
	}
	key := summaryKey(userID, msg.ConversationID)
	s := r.summaries[key]
	s.UserID = userID
	s.ConversationID = msg.ConversationID
	s.LastMessageID = msg.ID
	s.LastMessagePreview = msg.Content
	if incrementUnread {
		s.UnreadCount++
	}
	s.UpdatedAt = time.Now()
	r.summaries[key] = s
}
