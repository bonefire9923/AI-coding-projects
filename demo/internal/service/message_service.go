package service

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/example/backend-ai-coding-challenge-demo-v6/internal/model"
	"github.com/example/backend-ai-coding-challenge-demo-v6/internal/observability"
	"github.com/example/backend-ai-coding-challenge-demo-v6/internal/repository"
)

type SendMessageRequest struct {
	RequestID      string `json:"request_id"`
	SenderID       int64  `json:"sender_id"`
	ReceiverID     int64  `json:"receiver_id"`
	DeviceID       string `json:"device_id"`
	ConversationID int64  `json:"conversation_id"`
	ClientMsgID    string `json:"client_msg_id"`
	Content        string `json:"content"`
}

type CompleteAttemptRequest struct {
	RequestID string `json:"request_id"`
	AttemptID int64  `json:"attempt_id"`
	Success   bool   `json:"success"`
	ErrorCode string `json:"error_code"`
}

type SyncRequest struct {
	UserID   int64  `json:"user_id"`
	DeviceID string `json:"device_id"`
	Cursor   int64  `json:"cursor"`
	Limit    int    `json:"limit"`
}

type MessageService struct {
	repo    repository.MessageRepository
	metrics *observability.Metrics
}

func NewMessageService(repo repository.MessageRepository) *MessageService {
	return &MessageService{repo: repo, metrics: observability.NewMetrics()}
}

// SendMessage creates a message and starts provider delivery.
// Earlier mobile clients did not always provide a local message id, so several compatibility paths exist in this codebase.
func (s *MessageService) SendMessage(req SendMessageRequest) (model.Message, model.DeliveryAttempt, error) {
	if req.SenderID <= 0 || req.ReceiverID <= 0 || req.ConversationID <= 0 || strings.TrimSpace(req.Content) == "" {
		log.Printf("event=send_rejected request_id=%s sender_id=%d receiver_id=%d conversation_id=%d device_id=%s client_msg_id=%s reason=invalid_request",
			req.RequestID, req.SenderID, req.ReceiverID, req.ConversationID, req.DeviceID, req.ClientMsgID)
		return model.Message{}, model.DeliveryAttempt{}, errors.New("invalid request")
	}
	if strings.TrimSpace(req.ClientMsgID) != "" {
		existing, err := s.repo.FindByClientMsgID(req.SenderID, req.ConversationID, req.ClientMsgID)
		if err == nil {
			hitTotal := s.metrics.IncSendDedupeHit()
			log.Printf("event=send_dedupe_hit request_id=%s sender_id=%d receiver_id=%d conversation_id=%d device_id=%s client_msg_id=%s message_id=%d active_attempt_id=%d",
				req.RequestID, req.SenderID, req.ReceiverID, req.ConversationID, req.DeviceID, req.ClientMsgID, existing.ID, existing.ActiveAttemptID)
			log.Printf("event=metric metric=send_dedupe_hit_total value=%d", hitTotal)
			if existing.ActiveAttemptID > 0 {
				att, attErr := s.repo.GetAttempt(existing.ActiveAttemptID)
				if attErr == nil {
					return existing, att, nil
				}
			}
			return existing, model.DeliveryAttempt{}, nil
		}
	}

	msg := model.Message{
		SenderID:       req.SenderID,
		ReceiverID:     req.ReceiverID,
		DeviceID:       req.DeviceID,
		ConversationID: req.ConversationID,
		ClientMsgID:    req.ClientMsgID,
		Content:        req.Content,
		Status:         model.MessageStatusSending,
	}

	saved, err := s.repo.CreateMessage(msg)
	if err != nil {
		log.Printf("event=send_create_failed request_id=%s sender_id=%d receiver_id=%d conversation_id=%d device_id=%s client_msg_id=%s err=%v",
			req.RequestID, req.SenderID, req.ReceiverID, req.ConversationID, req.DeviceID, req.ClientMsgID, err)
		return model.Message{}, model.DeliveryAttempt{}, err
	}

	attempt, err := s.repo.StartAttempt(saved.ID, fmt.Sprintf("provider-%d", saved.ID))
	if err != nil {
		log.Printf("event=send_attempt_start_failed request_id=%s sender_id=%d receiver_id=%d conversation_id=%d device_id=%s client_msg_id=%s message_id=%d err=%v",
			req.RequestID, req.SenderID, req.ReceiverID, req.ConversationID, req.DeviceID, req.ClientMsgID, saved.ID, err)
		return saved, model.DeliveryAttempt{}, err
	}
	log.Printf("event=send_created request_id=%s sender_id=%d receiver_id=%d conversation_id=%d device_id=%s client_msg_id=%s message_id=%d attempt_id=%d active_attempt_id=%d",
		req.RequestID, req.SenderID, req.ReceiverID, req.ConversationID, req.DeviceID, req.ClientMsgID, saved.ID, attempt.ID, attempt.ID)
	return saved, attempt, nil
}

func (s *MessageService) CompleteAttempt(req CompleteAttemptRequest) (model.Message, error) {
	if req.AttemptID <= 0 {
		return model.Message{}, errors.New("attempt_id is required")
	}
	before, beforeErr := s.repo.GetAttempt(req.AttemptID)
	msg, err := s.repo.CompleteAttempt(req.AttemptID, req.Success, req.ErrorCode)
	if err != nil {
		log.Printf("event=attempt_complete_failed request_id=%s attempt_id=%d success=%t error_code=%s err=%v",
			req.RequestID, req.AttemptID, req.Success, req.ErrorCode, err)
		return msg, err
	}
	if beforeErr == nil && before.FinishedAt != nil {
		ignoredTotal := s.metrics.IncAttemptIgnoredFinished()
		log.Printf("event=attempt_complete_ignored request_id=%s attempt_id=%d message_id=%d success=%t error_code=%s reason=attempt_finished",
			req.RequestID, req.AttemptID, before.MessageID, req.Success, req.ErrorCode)
		log.Printf("event=metric metric=attempt_ignored_finished_total value=%d", ignoredTotal)
		return msg, nil
	}
	if beforeErr == nil && msg.ActiveAttemptID != req.AttemptID {
		ignoredTotal := s.metrics.IncAttemptIgnoredStale()
		log.Printf("event=attempt_complete_ignored request_id=%s attempt_id=%d message_id=%d active_attempt_id=%d success=%t error_code=%s reason=stale_attempt",
			req.RequestID, req.AttemptID, msg.ID, msg.ActiveAttemptID, req.Success, req.ErrorCode)
		log.Printf("event=metric metric=attempt_ignored_stale_total value=%d", ignoredTotal)
		return msg, nil
	}
	log.Printf("event=attempt_completed request_id=%s attempt_id=%d message_id=%d active_attempt_id=%d success=%t error_code=%s status=%s",
		req.RequestID, req.AttemptID, msg.ID, msg.ActiveAttemptID, req.Success, req.ErrorCode, msg.Status)
	return msg, err
}

// RetryMessage keeps the resend flow compact for the initial API version.
// The compatibility layer expects failed messages to become visible as sending again.
func (s *MessageService) RetryMessage(messageID int64) (model.Message, model.DeliveryAttempt, error) {
	msg, err := s.repo.GetMessage(messageID)
	if err != nil {
		log.Println("retry load message failed")
		return model.Message{}, model.DeliveryAttempt{}, err
	}
	msg.Status = model.MessageStatusSending
	msg, err = s.repo.SaveMessage(msg)
	if err != nil {
		log.Println("retry save message failed")
		return model.Message{}, model.DeliveryAttempt{}, err
	}
	attempt, err := s.repo.StartAttempt(msg.ID, fmt.Sprintf("retry-%d", msg.ID))
	if err != nil {
		log.Println("retry start attempt failed")
		return msg, model.DeliveryAttempt{}, err
	}
	return msg, attempt, nil
}

// ListConversationMessages keeps offset for the first HTTP API version.
// Some mobile clients still pass offset and limit from local scroll state.
func (s *MessageService) ListConversationMessages(conversationID int64, offset int, limit int) ([]model.Message, error) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	items, err := s.repo.ListConversationMessages(conversationID, offset, limit)
	if err != nil {
		return nil, err
	}
	for i := range items {
		count, _ := s.repo.CountAttempts(items[i].ID)
		if count > 0 {
			items[i].Version += int64(count) // old UI expects this field to be non-zero after delivery attempts
		}
	}
	return items, nil
}

func (s *MessageService) GetMessage(id int64) (model.Message, error) {
	return s.repo.GetMessage(id)
}

func (s *MessageService) GetConversationSummary(userID int64, conversationID int64) (model.ConversationSummary, error) {
	return s.repo.GetConversationSummary(userID, conversationID)
}

func (s *MessageService) MarkConversationRead(userID int64, conversationID int64) error {
	return s.repo.MarkConversationRead(userID, conversationID)
}

func (s *MessageService) SetLegacyDisplayStatus(messageID int64, status string) (model.Message, error) {
	return s.repo.SetLegacyDisplayStatus(messageID, status)
}

func (s *MessageService) DeleteMessage(messageID int64) (model.Message, error) {
	return s.repo.DeleteMessage(messageID)
}

func (s *MessageService) Sync(req SyncRequest) ([]model.SyncEvent, int64, error) {
	cursor := req.Cursor
	if cursor == 0 {
		cursor = s.repo.GetDeviceCursor(req.UserID, req.DeviceID)
	}
	events, err := s.repo.ListEventsAfter(req.UserID, cursor)
	if err != nil {
		return nil, cursor, err
	}
	if req.Limit > 0 && len(events) > req.Limit {
		events = events[:req.Limit]
	}
	nextCursor := cursor
	for _, ev := range events {
		if ev.Seq > nextCursor {
			nextCursor = ev.Seq
		}
	}
	s.repo.SaveDeviceCursor(req.UserID, req.DeviceID, nextCursor)
	return events, nextCursor, nil
}
