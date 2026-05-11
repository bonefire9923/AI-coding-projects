# AI_CONVERSATION

## Session Date

- 2026-05-11

## User Prompt (key intent)

- 记住文档里的要求，利用 gstack 进行设计架构，然后交给 superpowers 进行开发，并以 `README.md` 文档要求作为测试案例进行 TDD 开发。
- 额外确认：`AI_CONVERSATION.md` 需要边做边记录。

## Assistant Context Gathering

1. Read top-level `README.md` and extracted business reliability goals:
   - dedupe, no rollback, retry isolation, multi-device consistency, summary/unread correctness, observability.
2. Read demo code:
   - `demo/internal/service/message_service.go`
   - `demo/internal/repository/message_repository.go`
   - `demo/tests/public_test.go`
   - `demo/internal/model/message.go`
3. Read compatibility docs:
   - `demo/docs/provider_callback_contract.md`
   - `demo/docs/legacy_migration_notes.md`

## Architecture/Design Decisions (gstack-style output)

Design doc written to:

- `docs/superpowers/specs/2026-05-11-im-reliability-design.md`

Chosen architecture:

1. Idempotent send by `(sender_id, client_msg_id)` when `client_msg_id` is present.
2. Callback convergence with active-attempt guard:
   - ignore callback if attempt already finished;
   - ignore callback if attempt is not current `active_attempt_id`.
3. Retry continues to create new attempt; stale old attempts become non-authoritative.
4. Unread correctness achieved through callback idempotency.

Alternative considered but rejected:

- Full event-sourcing/projection rewrite.
- Rejected because scope too large for this task and not required by README acceptance goals.

## superpowers-style Implementation Plan

Plan doc written to:

- `docs/superpowers/plans/2026-05-11-im-reliability-implementation-plan.md`

Execution pattern:

1. Add tests first (RED).
2. Apply minimal code change (GREEN).
3. Keep changes focused, avoid unrelated refactor.

## TDD Test Additions

Added in `demo/tests/public_test.go`:

1. `TestSendMessageDeduplicatesByClientMsgID`
2. `TestStaleAttemptCannotRollbackSuccessfulMessage`
3. `TestDuplicateSuccessCallbackDoesNotOvercountUnread`

## Code Changes Applied

1. `demo/internal/repository/message_repository.go`
   - repository interface adds `GetAttempt(id int64)`.
   - in-memory implementation adds `GetAttempt`.
   - `CompleteAttempt` now no-ops for already-finished attempts.
   - `CompleteAttempt` now no-ops for stale (non-active) attempts.
2. `demo/internal/service/message_service.go`
   - `SendMessage` now deduplicates by existing `(sender_id, client_msg_id)` when `client_msg_id` non-empty.
   - when dedup hit and active attempt exists, returns existing active attempt.

## Verification Attempt and Limitation

Attempted:

- `cd demo && go test ./...`

Environment result:

- `go: command not found`

Conclusion:

- Tests were authored in TDD order, but runtime verification is blocked in current environment due to missing Go toolchain.

## Human Judgment Notes

1. Prioritized user-visible reliability faults over broad refactors.
2. Kept compatibility behavior for empty `client_msg_id` (legacy clients).
3. Chose attempt-authoritative convergence because provider callbacks can be delayed/retried per docs.

## Follow-up: README supplement + local Go install

User requested:

- Continue by supplementing `demo/README.md` run instructions and install Go directly in current environment.

Actions:

1. Updated `demo/README.md`:
   - added Go version prerequisite (`1.22.x`);
   - added local install method to `.tools/go`;
   - added explicit `cd demo` run commands;
   - added focused reliability test run command.
2. Installed Go locally:
   - downloaded `go1.22.12.linux-amd64.tar.gz`;
   - extracted to `<repo>/.tools/go`;
   - verified with `go version go1.22.12 linux/amd64`.
3. Executed tests with writable cache paths in this sandbox:
   - `GOCACHE=/tmp/gocache`
   - `GOMODCACHE=/tmp/gomodcache`
   - full: `go test ./...` passed.
   - focused: three newly added reliability tests all passed.

## Follow-up: 3-hour gap-closing pass

User decision:

- Agreed to execute the 3-hour minimal gap list.

Implemented:

1. Dedupe correctness hardening:
   - repository `FindByClientMsgID` changed from `(sender_id, client_msg_id)` to `(sender_id, conversation_id, client_msg_id)`.
   - avoids cross-conversation accidental dedupe.
2. Multi-device and summary consistency tests added:
   - `TestDedupDoesNotCrossConversationBoundary`
   - `TestSyncAcrossDevicesConvergesWithoutDuplicates`
   - `TestMarkReadThenSyncDoesNotReinflateUnread`
   - `TestSummaryMatchesConversationLastMessageAfterRetrySuccess`
3. Observability upgrade:
   - `message_service` logs now include stable correlation fields:
     `request_id sender_id receiver_id conversation_id device_id client_msg_id message_id attempt_id active_attempt_id success error_code`.
   - explicit ignored callback reason logs:
     - `reason=attempt_finished`
     - `reason=stale_attempt`
4. README troubleshooting section added:
   - reliability-focused test commands
   - key log fields and key event names.

Rejected options (with reasons):

1. Full projection/event-sourcing rewrite.
   - Rejected as out-of-scope for 3-hour pass and unnecessary for README minimum acceptance.
2. Introducing persistent DB/index migration in this round.
   - Rejected due to large migration surface and high regression risk versus current in-memory demo baseline.
3. Rebuilding sync protocol with new event schema.
   - Rejected because existing cursor/event model already sufficient for current acceptance tests; focused on guardrails and regression tests first.

Verification:

1. Full suite passed:
   - `cd demo && GOCACHE=/tmp/gocache GOMODCACHE=/tmp/gomodcache ../.tools/go/bin/go test ./...`
2. Gap-focused suite passed:
   - `cd demo && GOCACHE=/tmp/gocache GOMODCACHE=/tmp/gomodcache ../.tools/go/bin/go test ./tests -run 'TestSyncAcrossDevicesConvergesWithoutDuplicates|TestMarkReadThenSyncDoesNotReinflateUnread|TestSummaryMatchesConversationLastMessageAfterRetrySuccess|TestDedupDoesNotCrossConversationBoundary' -v`

## Follow-up: shortboard optimization pass (phase 1)

User decision:

- Approved full optimization of remaining shortboards in incremental mode.

Interaction timeline (near-original):

1. User asked whether all shortboards can be optimized.
2. Assistant answered yes, split into phase-1 (1-2 days) and phase-2 (3-7 days), and asked whether to start phase-1 now.
3. User replied "同意".
4. Assistant implemented phase-1 in this session.

Phase-1 implementation details:

1. Performance indexing in memory repository:
   - Added `clientMsgIndex` map keyed by `sender_id + conversation_id + client_msg_id`.
   - Added `conversationMessages` map for direct conversation message traversal.
   - Added `userEvents` map for direct per-user event scan in `ListEventsAfter`.
   - Result: removes linear full-map scans on common paths.
2. Additional consistency/high-pressure tests:
   - `TestSyncWithLimitPaginatesWithoutLoss`
   - `TestMultipleRetriesWithLateCallbacksConvergesToLatestSuccess`
3. Observability metrics counters:
   - Added lightweight counters in `internal/observability/metrics.go`.
   - Emitted `event=metric` logs for:
     - `send_dedupe_hit_total`
     - `attempt_ignored_finished_total`
     - `attempt_ignored_stale_total`
4. README verification guidance extended:
   - Added long-horizon consistency test command set.
   - Documented `event=metric` in troubleshooting.

Additional rejected options in this pass:

1. Introducing background reconciliation worker.
   - Rejected for now because current acceptance criteria can be covered by deterministic synchronous state guards + tests.
2. Converting to persistent metrics backend (Prometheus/OpenTelemetry).
   - Rejected due to scope and infra dependencies; used log-based counters as minimal viable observability.

Verification (phase-1):

1. `cd demo && ../.tools/go/bin/gofmt -w internal/repository/message_repository.go internal/service/message_service.go internal/observability/metrics.go tests/public_test.go`
2. `cd demo && GOCACHE=/tmp/gocache GOMODCACHE=/tmp/gomodcache ../.tools/go/bin/go test ./...` passed.

## Follow-up: persistent unique constraints/indexes pass

User decision:

- Agreed to start with "persistent unique constraints/indexes" immediately.

Implemented:

1. Added file-backed repository constructor:
   - `repository.NewFileMessageRepository(path)`
   - repository state is persisted atomically to JSON snapshot (`.tmp` + `rename`).
2. Persisted and recovered state includes:
   - messages / attempts / summaries / device cursors
   - `clientMsgIndex` unique key map (`sender_id + conversation_id + client_msg_id`)
   - `conversationMessages` and `userEvents` indexes
3. Wired runtime switch:
   - `DEMO_DATA_PATH` enables file-backed repository in `cmd/server/main.go`.
4. Added recovery regression test:
   - `TestFileRepositoryPersistsDataAndDedupeIndex`
   - verifies restart keeps unread summary and dedupe index behavior.

Rejected option in this pass:

1. Direct migration to SQL database with DDL constraints.
   - Rejected in this step to keep incremental risk low and avoid introducing infra/tooling surface before finishing file-backed persistence baseline.

Verification:

1. `cd demo && ../.tools/go/bin/gofmt -w internal/repository/message_repository.go cmd/server/main.go tests/public_test.go`
2. `cd demo && GOCACHE=/tmp/gocache GOMODCACHE=/tmp/gomodcache ../.tools/go/bin/go test ./...` passed.
