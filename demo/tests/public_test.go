package tests

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/example/backend-ai-coding-challenge-demo-v6/internal/model"
	"github.com/example/backend-ai-coding-challenge-demo-v6/internal/repository"
	"github.com/example/backend-ai-coding-challenge-demo-v6/internal/service"
)

func newService() *service.MessageService {
	repo := repository.NewMemoryMessageRepository()
	return service.NewMessageService(repo)
}

func newFileService(t *testing.T, dataPath string) *service.MessageService {
	t.Helper()
	repo, err := repository.NewFileMessageRepository(dataPath)
	if err != nil {
		t.Fatalf("new file repository failed: %v", err)
	}
	return service.NewMessageService(repo)
}

func send(t *testing.T, svc *service.MessageService, clientID string) (model.Message, model.DeliveryAttempt) {
	t.Helper()
	msg, att, err := svc.SendMessage(service.SendMessageRequest{
		RequestID:      "req-" + clientID,
		SenderID:       1,
		ReceiverID:     2,
		DeviceID:       "device-a",
		ConversationID: 100,
		ClientMsgID:    clientID,
		Content:        "hello",
	})
	if err != nil {
		t.Fatalf("send message failed: %v", err)
	}
	return msg, att
}

func TestSendMessageCreatesAttempt(t *testing.T) {
	svc := newService()
	msg, attempt := send(t, svc, "local-001")
	if msg.ID <= 0 || attempt.ID <= 0 {
		t.Fatalf("expected ids, msg=%d attempt=%d", msg.ID, attempt.ID)
	}
	if msg.Status != model.MessageStatusSending {
		t.Fatalf("expected sending status, got %s", msg.Status)
	}
}

func TestCompleteAttemptSuccess(t *testing.T) {
	svc := newService()
	_, attempt := send(t, svc, "local-002")
	updated, err := svc.CompleteAttempt(service.CompleteAttemptRequest{RequestID: "callback-1", AttemptID: attempt.ID, Success: true})
	if err != nil {
		t.Fatalf("complete attempt failed: %v", err)
	}
	if updated.Status != model.MessageStatusSent {
		t.Fatalf("expected sent status, got %s", updated.Status)
	}
}

func TestRetryMessageCreatesNewAttempt(t *testing.T) {
	svc := newService()
	msg, attempt := send(t, svc, "local-retry")
	_, err := svc.CompleteAttempt(service.CompleteAttemptRequest{AttemptID: attempt.ID, Success: false, ErrorCode: "network_error"})
	if err != nil {
		t.Fatalf("complete attempt failed: %v", err)
	}
	retried, retryAttempt, err := svc.RetryMessage(msg.ID)
	if err != nil {
		t.Fatalf("retry failed: %v", err)
	}
	if retryAttempt.ID <= attempt.ID {
		t.Fatalf("expected new attempt")
	}
	if retried.Status != model.MessageStatusSending {
		t.Fatalf("expected sending after retry, got %s", retried.Status)
	}
}

func TestListConversationMessagesAndSummary(t *testing.T) {
	svc := newService()
	for i := 0; i < 3; i++ {
		msg, attempt := send(t, svc, fmt.Sprintf("local-list-%d", i))
		_, err := svc.CompleteAttempt(service.CompleteAttemptRequest{AttemptID: attempt.ID, Success: true})
		if err != nil {
			t.Fatalf("complete attempt failed for msg %d: %v", msg.ID, err)
		}
	}
	items, err := svc.ListConversationMessages(100, 0, 20)
	if err != nil {
		t.Fatalf("list failed: %v", err)
	}
	if len(items) != 3 {
		t.Fatalf("expected 3 messages, got %d", len(items))
	}
	summary, err := svc.GetConversationSummary(2, 100)
	if err != nil {
		t.Fatalf("summary failed: %v", err)
	}
	if summary.UnreadCount == 0 {
		t.Fatalf("expected receiver unread count")
	}
}

func TestSyncReturnsEvents(t *testing.T) {
	svc := newService()
	_, attempt := send(t, svc, "local-sync")
	_, err := svc.CompleteAttempt(service.CompleteAttemptRequest{AttemptID: attempt.ID, Success: true})
	if err != nil {
		t.Fatalf("complete attempt failed: %v", err)
	}
	events, cursor, err := svc.Sync(service.SyncRequest{UserID: 2, DeviceID: "device-b", Cursor: 0})
	if err != nil {
		t.Fatalf("sync failed: %v", err)
	}
	if len(events) == 0 || cursor == 0 {
		t.Fatalf("expected sync events and cursor")
	}
}

func TestSendMessageDeduplicatesByClientMsgID(t *testing.T) {
	svc := newService()
	msg1, att1 := send(t, svc, "local-dedupe")
	msg2, att2 := send(t, svc, "local-dedupe")
	if msg1.ID != msg2.ID {
		t.Fatalf("expected same message id, got %d and %d", msg1.ID, msg2.ID)
	}
	if att1.ID != att2.ID {
		t.Fatalf("expected same active attempt id, got %d and %d", att1.ID, att2.ID)
	}
	items, err := svc.ListConversationMessages(100, 0, 20)
	if err != nil {
		t.Fatalf("list failed: %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("expected one visible message after duplicate sends, got %d", len(items))
	}
}

func TestStaleAttemptCannotRollbackSuccessfulMessage(t *testing.T) {
	svc := newService()
	msg, firstAttempt := send(t, svc, "local-stale-callback")
	_, err := svc.CompleteAttempt(service.CompleteAttemptRequest{AttemptID: firstAttempt.ID, Success: false, ErrorCode: "network_error"})
	if err != nil {
		t.Fatalf("complete first attempt failed: %v", err)
	}
	msg, secondAttempt, err := svc.RetryMessage(msg.ID)
	if err != nil {
		t.Fatalf("retry failed: %v", err)
	}
	msg, err = svc.CompleteAttempt(service.CompleteAttemptRequest{AttemptID: secondAttempt.ID, Success: true})
	if err != nil {
		t.Fatalf("complete second attempt failed: %v", err)
	}
	if msg.Status != model.MessageStatusSent {
		t.Fatalf("expected sent after successful retry, got %s", msg.Status)
	}
	msg, err = svc.CompleteAttempt(service.CompleteAttemptRequest{AttemptID: firstAttempt.ID, Success: false, ErrorCode: "late_error"})
	if err != nil {
		t.Fatalf("late callback should be ignored, got err: %v", err)
	}
	if msg.Status != model.MessageStatusSent {
		t.Fatalf("expected stale callback not to rollback status, got %s", msg.Status)
	}
}

func TestDuplicateSuccessCallbackDoesNotOvercountUnread(t *testing.T) {
	svc := newService()
	_, attempt := send(t, svc, "local-duplicate-callback")
	_, err := svc.CompleteAttempt(service.CompleteAttemptRequest{AttemptID: attempt.ID, Success: true})
	if err != nil {
		t.Fatalf("first callback failed: %v", err)
	}
	_, err = svc.CompleteAttempt(service.CompleteAttemptRequest{AttemptID: attempt.ID, Success: true})
	if err != nil {
		t.Fatalf("duplicate callback should be ignored, got err: %v", err)
	}
	summary, err := svc.GetConversationSummary(2, 100)
	if err != nil {
		t.Fatalf("summary failed: %v", err)
	}
	if summary.UnreadCount != 1 {
		t.Fatalf("expected unread count 1 after duplicate callback, got %d", summary.UnreadCount)
	}
}

func TestDedupDoesNotCrossConversationBoundary(t *testing.T) {
	svc := newService()
	req := service.SendMessageRequest{
		RequestID:      "req-cross-1",
		SenderID:       1,
		ReceiverID:     2,
		DeviceID:       "device-a",
		ConversationID: 100,
		ClientMsgID:    "dup-cross",
		Content:        "hello",
	}
	msg1, _, err := svc.SendMessage(req)
	if err != nil {
		t.Fatalf("send 1 failed: %v", err)
	}
	req.RequestID = "req-cross-2"
	req.ConversationID = 101
	msg2, _, err := svc.SendMessage(req)
	if err != nil {
		t.Fatalf("send 2 failed: %v", err)
	}
	if msg1.ID == msg2.ID {
		t.Fatalf("expected different messages across conversations")
	}
}

func TestSyncAcrossDevicesConvergesWithoutDuplicates(t *testing.T) {
	svc := newService()
	_, attempt := send(t, svc, "local-sync-multi-device")
	_, err := svc.CompleteAttempt(service.CompleteAttemptRequest{RequestID: "cb-sync-md", AttemptID: attempt.ID, Success: true})
	if err != nil {
		t.Fatalf("complete attempt failed: %v", err)
	}
	eventsA1, cursorA1, err := svc.Sync(service.SyncRequest{UserID: 2, DeviceID: "device-b", Cursor: 0})
	if err != nil {
		t.Fatalf("first sync device-b failed: %v", err)
	}
	if len(eventsA1) == 0 || cursorA1 == 0 {
		t.Fatalf("expected events on first sync")
	}
	eventsA2, cursorA2, err := svc.Sync(service.SyncRequest{UserID: 2, DeviceID: "device-b", Cursor: cursorA1})
	if err != nil {
		t.Fatalf("second sync device-b failed: %v", err)
	}
	if len(eventsA2) != 0 || cursorA2 != cursorA1 {
		t.Fatalf("expected no duplicates on second sync")
	}
	eventsWeb, cursorWeb, err := svc.Sync(service.SyncRequest{UserID: 2, DeviceID: "device-web", Cursor: 0})
	if err != nil {
		t.Fatalf("sync device-web failed: %v", err)
	}
	if len(eventsWeb) != len(eventsA1) || cursorWeb != cursorA1 {
		t.Fatalf("expected another device to converge to same cursor and event count")
	}
}

func TestMarkReadThenSyncDoesNotReinflateUnread(t *testing.T) {
	svc := newService()
	_, attempt := send(t, svc, "local-read-then-sync")
	_, err := svc.CompleteAttempt(service.CompleteAttemptRequest{RequestID: "cb-read-sync", AttemptID: attempt.ID, Success: true})
	if err != nil {
		t.Fatalf("complete attempt failed: %v", err)
	}
	if err := svc.MarkConversationRead(2, 100); err != nil {
		t.Fatalf("mark read failed: %v", err)
	}
	_, _, err = svc.Sync(service.SyncRequest{UserID: 2, DeviceID: "device-b", Cursor: 0})
	if err != nil {
		t.Fatalf("sync failed: %v", err)
	}
	summary, err := svc.GetConversationSummary(2, 100)
	if err != nil {
		t.Fatalf("summary failed: %v", err)
	}
	if summary.UnreadCount != 0 {
		t.Fatalf("expected unread to stay zero after sync, got %d", summary.UnreadCount)
	}
}

func TestSummaryMatchesConversationLastMessageAfterRetrySuccess(t *testing.T) {
	svc := newService()
	original, firstAttempt := send(t, svc, "local-summary-retry")
	_, err := svc.CompleteAttempt(service.CompleteAttemptRequest{RequestID: "cb-retry-fail", AttemptID: firstAttempt.ID, Success: false, ErrorCode: "timeout"})
	if err != nil {
		t.Fatalf("complete first attempt failed: %v", err)
	}
	msg, retryAttempt, err := svc.RetryMessage(original.ID)
	if err != nil {
		t.Fatalf("retry failed: %v", err)
	}
	_, err = svc.CompleteAttempt(service.CompleteAttemptRequest{RequestID: "cb-retry-success", AttemptID: retryAttempt.ID, Success: true})
	if err != nil {
		t.Fatalf("complete retry attempt failed: %v", err)
	}
	items, err := svc.ListConversationMessages(100, 0, 1)
	if err != nil {
		t.Fatalf("list failed: %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("expected one message in conversation")
	}
	summarySender, err := svc.GetConversationSummary(1, 100)
	if err != nil {
		t.Fatalf("sender summary failed: %v", err)
	}
	if summarySender.LastMessageID != items[0].ID || summarySender.LastMessagePreview != items[0].Content {
		t.Fatalf("sender summary mismatch with conversation last message")
	}
	if msg.ID != items[0].ID {
		t.Fatalf("retry should keep same message id")
	}
}

func TestSyncWithLimitPaginatesWithoutLoss(t *testing.T) {
	svc := newService()
	for i := 0; i < 3; i++ {
		msg, attempt := send(t, svc, fmt.Sprintf("local-sync-limit-%d", i))
		_, err := svc.CompleteAttempt(service.CompleteAttemptRequest{
			RequestID: fmt.Sprintf("cb-sync-limit-%d", msg.ID),
			AttemptID: attempt.ID,
			Success:   true,
		})
		if err != nil {
			t.Fatalf("complete attempt failed: %v", err)
		}
	}
	page1, cursor1, err := svc.Sync(service.SyncRequest{UserID: 2, DeviceID: "device-limit", Cursor: 0, Limit: 2})
	if err != nil {
		t.Fatalf("sync page1 failed: %v", err)
	}
	if len(page1) != 2 {
		t.Fatalf("expected 2 events on page1, got %d", len(page1))
	}
	page2, cursor2, err := svc.Sync(service.SyncRequest{UserID: 2, DeviceID: "device-limit", Cursor: cursor1, Limit: 10})
	if err != nil {
		t.Fatalf("sync page2 failed: %v", err)
	}
	if len(page2) == 0 {
		t.Fatalf("expected remaining events on page2")
	}
	if cursor2 <= cursor1 {
		t.Fatalf("expected cursor progress across pages")
	}
}

func TestMultipleRetriesWithLateCallbacksConvergesToLatestSuccess(t *testing.T) {
	svc := newService()
	msg, attempt1 := send(t, svc, "local-retry-late-many")
	_, err := svc.CompleteAttempt(service.CompleteAttemptRequest{RequestID: "cb-retry-late-1", AttemptID: attempt1.ID, Success: false, ErrorCode: "net_1"})
	if err != nil {
		t.Fatalf("complete attempt1 failed: %v", err)
	}
	_, attempt2, err := svc.RetryMessage(msg.ID)
	if err != nil {
		t.Fatalf("retry 1 failed: %v", err)
	}
	_, err = svc.CompleteAttempt(service.CompleteAttemptRequest{RequestID: "cb-retry-late-2", AttemptID: attempt2.ID, Success: false, ErrorCode: "net_2"})
	if err != nil {
		t.Fatalf("complete attempt2 failed: %v", err)
	}
	_, attempt3, err := svc.RetryMessage(msg.ID)
	if err != nil {
		t.Fatalf("retry 2 failed: %v", err)
	}
	msg, err = svc.CompleteAttempt(service.CompleteAttemptRequest{RequestID: "cb-retry-late-3", AttemptID: attempt3.ID, Success: true})
	if err != nil {
		t.Fatalf("complete attempt3 failed: %v", err)
	}
	if msg.Status != model.MessageStatusSent {
		t.Fatalf("expected sent after latest success, got %s", msg.Status)
	}
	msg, err = svc.CompleteAttempt(service.CompleteAttemptRequest{RequestID: "cb-retry-late-old", AttemptID: attempt1.ID, Success: false, ErrorCode: "late_old"})
	if err != nil {
		t.Fatalf("late callback should be ignored: %v", err)
	}
	if msg.Status != model.MessageStatusSent {
		t.Fatalf("expected late old callback not to rollback latest success")
	}
}

func TestFileRepositoryPersistsDataAndDedupeIndex(t *testing.T) {
	tmp := t.TempDir()
	dataPath := filepath.Join(tmp, "repo.json")

	svc1 := newFileService(t, dataPath)
	msg1, attempt1 := send(t, svc1, "local-persist-001")
	_, err := svc1.CompleteAttempt(service.CompleteAttemptRequest{RequestID: "cb-persist-1", AttemptID: attempt1.ID, Success: true})
	if err != nil {
		t.Fatalf("complete attempt failed: %v", err)
	}

	svc2 := newFileService(t, dataPath)
	msg2, attempt2, err := svc2.SendMessage(service.SendMessageRequest{
		RequestID:      "req-local-persist-001-replay",
		SenderID:       1,
		ReceiverID:     2,
		DeviceID:       "device-a",
		ConversationID: 100,
		ClientMsgID:    "local-persist-001",
		Content:        "hello",
	})
	if err != nil {
		t.Fatalf("send on recovered repo failed: %v", err)
	}
	if msg2.ID != msg1.ID {
		t.Fatalf("expected dedupe after restart, ids: %d vs %d", msg1.ID, msg2.ID)
	}
	if attempt2.ID != attempt1.ID {
		t.Fatalf("expected same active attempt after restart, attempts: %d vs %d", attempt1.ID, attempt2.ID)
	}
	summary, err := svc2.GetConversationSummary(2, 100)
	if err != nil {
		t.Fatalf("summary on recovered repo failed: %v", err)
	}
	if summary.UnreadCount != 1 {
		t.Fatalf("expected unread to persist after restart, got %d", summary.UnreadCount)
	}
}
