# IM Reliability Implementation Plan (superpowers execution)

## Step 1. Red tests

Add failing tests in `demo/tests/public_test.go`:

1. `TestSendMessageDeduplicatesByClientMsgID`
2. `TestStaleAttemptCannotRollbackSuccessfulMessage`
3. `TestDuplicateSuccessCallbackDoesNotOvercountUnread`

## Step 2. Green minimal implementation

1. Repository interface add `GetAttempt`.
2. `SendMessage` checks existing message by `(sender_id, client_msg_id)` before create.
3. `CompleteAttempt` adds no-op guards for:
  - already finished attempt
  - non-active stale attempt

## Step 3. Refactor/cleanup

1. Keep logic minimal in service layer.
2. Keep callback convergence in repository state mutation path.

## Step 4. Verification

1. Run `cd demo && go test ./...`
2. If toolchain unavailable, mark verification gap explicitly in final report.
