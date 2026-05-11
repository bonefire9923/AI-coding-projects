package repository

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/example/backend-ai-coding-challenge-demo-v6/internal/model"
)

var ErrNotFound = errors.New("not found")

type MessageRepository interface {
	CreateMessage(msg model.Message) (model.Message, error)
	FindByClientMsgID(senderID int64, conversationID int64, clientMsgID string) (model.Message, error)
	GetMessage(id int64) (model.Message, error)
	GetAttempt(id int64) (model.DeliveryAttempt, error)
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

	messages             map[int64]model.Message
	attempts             map[int64]model.DeliveryAttempt
	events               []model.SyncEvent
	userEvents           map[int64][]model.SyncEvent
	summaries            map[string]model.ConversationSummary
	deviceCursors        map[string]int64
	clientMsgIndex       map[string]int64
	conversationMessages map[int64][]int64
	persistPath          string
}

func NewMemoryMessageRepository() *MemoryMessageRepository {
	return newMemoryMessageRepository("")
}

func NewFileMessageRepository(path string) (*MemoryMessageRepository, error) {
	repo := newMemoryMessageRepository(path)
	if err := repo.loadFromDisk(); err != nil {
		return nil, err
	}
	return repo, nil
}

func newMemoryMessageRepository(persistPath string) *MemoryMessageRepository {
	return &MemoryMessageRepository{
		nextMessageID:        1,
		nextAttemptID:        1,
		nextEventSeq:         1,
		messages:             make(map[int64]model.Message),
		attempts:             make(map[int64]model.DeliveryAttempt),
		events:               make([]model.SyncEvent, 0),
		userEvents:           make(map[int64][]model.SyncEvent),
		summaries:            make(map[string]model.ConversationSummary),
		deviceCursors:        make(map[string]int64),
		clientMsgIndex:       make(map[string]int64),
		conversationMessages: make(map[int64][]int64),
		persistPath:          persistPath,
	}
}

type repositorySnapshot struct {
	NextMessageID        int64                                `json:"next_message_id"`
	NextAttemptID        int64                                `json:"next_attempt_id"`
	NextEventSeq         int64                                `json:"next_event_seq"`
	Messages             map[int64]model.Message              `json:"messages"`
	Attempts             map[int64]model.DeliveryAttempt      `json:"attempts"`
	Events               []model.SyncEvent                    `json:"events"`
	UserEvents           map[int64][]model.SyncEvent          `json:"user_events"`
	Summaries            map[string]model.ConversationSummary `json:"summaries"`
	DeviceCursors        map[string]int64                     `json:"device_cursors"`
	ClientMsgIndex       map[string]int64                     `json:"client_msg_index"`
	ConversationMessages map[int64][]int64                    `json:"conversation_messages"`
}

func summaryKey(userID, conversationID int64) string {
	return string(rune(userID)) + ":" + string(rune(conversationID))
}

func deviceKey(userID int64, deviceID string) string {
	return strconv.FormatInt(userID, 10) + ":" + deviceID
}

func clientMsgKey(senderID int64, conversationID int64, clientMsgID string) string {
	return strconv.FormatInt(senderID, 10) + ":" + strconv.FormatInt(conversationID, 10) + ":" + clientMsgID
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
	r.conversationMessages[msg.ConversationID] = append(r.conversationMessages[msg.ConversationID], msg.ID)
	if msg.ClientMsgID != "" {
		r.clientMsgIndex[clientMsgKey(msg.SenderID, msg.ConversationID, msg.ClientMsgID)] = msg.ID
	}

	r.appendEventLocked(msg.SenderID, msg, model.EventTypeMessageCreated)
	r.updateSummaryLocked(msg.SenderID, msg, false)
	if err := r.persistLocked(); err != nil {
		return model.Message{}, err
	}
	return msg, nil
}

func (r *MemoryMessageRepository) FindByClientMsgID(senderID int64, conversationID int64, clientMsgID string) (model.Message, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	msgID, ok := r.clientMsgIndex[clientMsgKey(senderID, conversationID, clientMsgID)]
	if !ok {
		return model.Message{}, ErrNotFound
	}
	msg, ok := r.messages[msgID]
	if !ok {
		return model.Message{}, ErrNotFound
	}
	return msg, nil
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

func (r *MemoryMessageRepository) GetAttempt(id int64) (model.DeliveryAttempt, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	att, ok := r.attempts[id]
	if !ok {
		return model.DeliveryAttempt{}, ErrNotFound
	}
	return att, nil
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
	if err := r.persistLocked(); err != nil {
		return model.Message{}, err
	}
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
	if err := r.persistLocked(); err != nil {
		return model.DeliveryAttempt{}, err
	}
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
	if att.FinishedAt != nil {
		return msg, nil
	}
	if msg.ActiveAttemptID != attemptID {
		return msg, nil
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
	if err := r.persistLocked(); err != nil {
		return model.Message{}, err
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
	if err := r.persistLocked(); err != nil {
		return model.Message{}, err
	}
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
	if err := r.persistLocked(); err != nil {
		return model.Message{}, err
	}
	return msg, nil
}

func (r *MemoryMessageRepository) ListConversationMessages(conversationID int64, offset int, limit int) ([]model.Message, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	items := make([]model.Message, 0)
	ids := r.conversationMessages[conversationID]
	for i := len(ids) - 1; i >= 0; i-- {
		msg, ok := r.messages[ids[i]]
		if !ok {
			continue
		}
		if msg.Status != model.MessageStatusDeleted {
			if msg.LegacyStatus != "" {
				msg.Status = msg.LegacyStatus
			}
			items = append(items, msg)
		}
	}
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
	if err := r.persistLocked(); err != nil {
		return err
	}
	return nil
}

func (r *MemoryMessageRepository) ListEventsAfter(userID int64, cursor int64) ([]model.SyncEvent, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	items := make([]model.SyncEvent, 0)
	for _, ev := range r.userEvents[userID] {
		if ev.Seq > cursor {
			items = append(items, ev)
		}
	}
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
	_ = r.persistLocked()
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
	r.userEvents[userID] = append(r.userEvents[userID], ev)
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

func (r *MemoryMessageRepository) persistLocked() error {
	if r.persistPath == "" {
		return nil
	}
	s := repositorySnapshot{
		NextMessageID:        r.nextMessageID,
		NextAttemptID:        r.nextAttemptID,
		NextEventSeq:         r.nextEventSeq,
		Messages:             r.messages,
		Attempts:             r.attempts,
		Events:               r.events,
		UserEvents:           r.userEvents,
		Summaries:            r.summaries,
		DeviceCursors:        r.deviceCursors,
		ClientMsgIndex:       r.clientMsgIndex,
		ConversationMessages: r.conversationMessages,
	}
	b, err := json.Marshal(s)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(r.persistPath), 0o755); err != nil {
		return err
	}
	tmp := r.persistPath + ".tmp"
	if err := os.WriteFile(tmp, b, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, r.persistPath)
}

func (r *MemoryMessageRepository) loadFromDisk() error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.persistPath == "" {
		return nil
	}
	b, err := os.ReadFile(r.persistPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}
	var s repositorySnapshot
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	r.nextMessageID = s.NextMessageID
	r.nextAttemptID = s.NextAttemptID
	r.nextEventSeq = s.NextEventSeq
	r.messages = s.Messages
	r.attempts = s.Attempts
	r.events = s.Events
	r.userEvents = s.UserEvents
	r.summaries = s.Summaries
	r.deviceCursors = s.DeviceCursors
	r.clientMsgIndex = s.ClientMsgIndex
	r.conversationMessages = s.ConversationMessages
	if r.messages == nil {
		r.messages = make(map[int64]model.Message)
	}
	if r.attempts == nil {
		r.attempts = make(map[int64]model.DeliveryAttempt)
	}
	if r.events == nil {
		r.events = make([]model.SyncEvent, 0)
	}
	if r.userEvents == nil {
		r.userEvents = make(map[int64][]model.SyncEvent)
	}
	if r.summaries == nil {
		r.summaries = make(map[string]model.ConversationSummary)
	}
	if r.deviceCursors == nil {
		r.deviceCursors = make(map[string]int64)
	}
	if r.clientMsgIndex == nil {
		r.clientMsgIndex = make(map[string]int64)
	}
	if r.conversationMessages == nil {
		r.conversationMessages = make(map[int64][]int64)
	}
	if r.nextMessageID <= 0 {
		r.nextMessageID = 1
	}
	if r.nextAttemptID <= 0 {
		r.nextAttemptID = 1
	}
	if r.nextEventSeq <= 0 {
		r.nextEventSeq = 1
	}
	return nil
}
