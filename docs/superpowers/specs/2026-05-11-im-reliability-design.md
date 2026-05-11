# IM Reliability Design (gstack-style)

## Scope

Target only the one-to-one message send chain in `demo/`:

1. `SendMessage`
2. `CompleteAttempt` callback convergence
3. `RetryMessage`
4. Conversation summary unread correctness

Out of scope: group chat, provider abstraction rewrite, storage engine rewrite.

## Reliability Goals from README

1. Prevent duplicate visible messages on repeated client send.
2. Prevent status rollback (`sent` should not be overwritten by stale failure callback).
3. Keep retry state isolated from stale older attempts.
4. Keep unread count and list summary aligned with message truth.
5. Preserve debuggability with deterministic state transitions.

## Architecture Decisions

### A. Idempotent send by `sender_id + client_msg_id`

- If the same sender sends with the same non-empty `client_msg_id`, return the existing message instead of creating a new message.
- If existing message has an active attempt, return that attempt as well.

Trade-off:
- Strong dedupe for modern clients.
- Empty `client_msg_id` keeps legacy behavior unchanged.

### B. Callback convergence guard by active attempt

- `CompleteAttempt(attempt_id)` only mutates message when:
1. attempt exists,
2. attempt is not already finished,
3. attempt equals message `active_attempt_id`.

- Otherwise callback is treated as stale/duplicate and ignored (return current message, no mutation).

Trade-off:
- Prevents late callback rollback.
- Gives up on updating history from stale callbacks, but matches product reliability priority.

### C. Retry semantics

- Retry keeps creating a new attempt and sets it as active.
- Older attempts become non-authoritative automatically through rule B.

### D. Unread correctness by idempotent callback handling

- Duplicate success callback for same attempt should not increment unread twice.
- Solved indirectly by rule B because finished attempts are no-op.

## Data/State Invariants

1. At most one authoritative attempt for message status: `message.active_attempt_id`.
2. Attempt completion is one-shot.
3. Message status progression across attempts is monotonic under stale callbacks.
4. Summary unread increments once per authoritative successful delivery event.

## TDD Cases (README-based)

1. Duplicate send with same client id returns same message and attempt, conversation shows one message.
2. Old failed callback after successful retry does not rollback `sent`.
3. Duplicate callback for same successful attempt does not overcount unread.

## Risks

1. In-memory repository still linear-scan for `FindByClientMsgID`; acceptable for demo.
2. Legacy empty `client_msg_id` still allows duplicates by design.
3. No external provider timestamp ordering; current ordering is attempt-authoritative only.
