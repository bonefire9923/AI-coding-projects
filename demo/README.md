# Demo：一对一私聊消息后端原型

这是一个早期后端 Demo，用于本次 AI Coding 工程评测。

业务上，它模拟的是一个最基础的一对一私聊消息系统：用户 A 给用户 B 发送文本消息，消息进入两人的聊天详情页，同时影响最近聊天列表中的最后一条消息、最近更新时间和未读角标。

这个 Demo 不是完整 IM 产品，而是一段早期链路。真实项目里，历史兼容、旧客户端迁移、第三方回执、已读回执、同步游标、实验逻辑和投影数据可能同时存在。你的目标不是重写一个新系统，而是判断哪些代码路径会影响用户可见的可靠性，并做必要改造。

## 运行环境

- Go `1.22.x`（建议 `1.22.12` 或以上 patch）
- Linux/macOS

如果本机没有 `go` 命令，可在仓库根目录执行以下“本地安装到 `.tools/go`”方式（不污染系统环境）：

```bash
cd ..
mkdir -p .tools
curl -fsSL https://go.dev/dl/go1.22.12.linux-amd64.tar.gz -o /tmp/go1.22.12.linux-amd64.tar.gz
tar -C .tools -xzf /tmp/go1.22.12.linux-amd64.tar.gz
export PATH="$(pwd)/.tools/go/bin:$PATH"
go version
```

如需长期生效，请把 `export PATH="<repo>/.tools/go/bin:$PATH"` 写入你的 shell profile。

## 运行

```bash
cd demo
go test ./...
go run ./cmd/server
```

服务默认监听 `:8080`。

启用文件持久化（重启保留消息、attempt、去重索引）：

```bash
cd demo
DEMO_DATA_PATH=./data/repo.json go run ./cmd/server
```

## 快速自检

```bash
cd demo
go test ./tests -run 'TestSendMessageDeduplicatesByClientMsgID|TestStaleAttemptCannotRollbackSuccessfulMessage|TestDuplicateSuccessCallbackDoesNotOvercountUnread' -v
```

一致性补充检查：

```bash
cd demo
go test ./tests -run 'TestSyncAcrossDevicesConvergesWithoutDuplicates|TestMarkReadThenSyncDoesNotReinflateUnread|TestSummaryMatchesConversationLastMessageAfterRetrySuccess|TestDedupDoesNotCrossConversationBoundary' -v
```

长期一致性压力补充：

```bash
cd demo
go test ./tests -run 'TestSyncWithLimitPaginatesWithoutLoss|TestMultipleRetriesWithLateCallbacksConvergesToLatestSuccess' -v
```

持久化恢复补充：

```bash
cd demo
go test ./tests -run 'TestFileRepositoryPersistsDataAndDedupeIndex' -v
```

## 排障建议

当出现“消息重复 / 成功后回退 / 未读异常”时，优先检索以下日志字段：

- `event`
- `request_id`
- `sender_id`
- `receiver_id`
- `conversation_id`
- `device_id`
- `client_msg_id`
- `message_id`
- `attempt_id`
- `active_attempt_id`
- `success`
- `error_code`

关键事件：

- `send_dedupe_hit`：重复发送被幂等命中
- `attempt_complete_ignored`：回调被忽略（`attempt_finished` 或 `stale_attempt`）
- `attempt_completed`：当前 authoritative attempt 生效
- `event=metric`：关键计数器输出（如去重命中累计、忽略回调累计）

## 主要目录

- `internal/service`：业务服务层
- `internal/repository`：内存数据层
- `internal/model`：核心模型
- `internal/receipt`、`internal/delivery`、`internal/readmodel`、`internal/checkpoint`：容易混淆的回执、已读、投递、同步相关模型
- `internal/legacy`、`internal/compat`、`internal/experiments`、`internal/history`、`internal/migration`、`internal/noise`：历史兼容、迁移、实验和运维上下文

不是所有看起来相关的文件都应该进入核心路径。真实项目中，AI 也经常会把历史逻辑、实验逻辑和核心逻辑混在一起。
