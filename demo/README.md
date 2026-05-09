# Demo：一对一私聊消息后端原型

这是一个早期后端 Demo，用于本次 AI Coding 工程评测。

业务上，它模拟的是一个最基础的一对一私聊消息系统：用户 A 给用户 B 发送文本消息，消息进入两人的聊天详情页，同时影响最近聊天列表中的最后一条消息、最近更新时间和未读角标。

这个 Demo 不是完整 IM 产品，而是一段早期链路。真实项目里，历史兼容、旧客户端迁移、第三方回执、已读回执、同步游标、实验逻辑和投影数据可能同时存在。你的目标不是重写一个新系统，而是判断哪些代码路径会影响用户可见的可靠性，并做必要改造。

## 运行

```bash
go test ./...
go run ./cmd/server
```

## 主要目录

- `internal/service`：业务服务层
- `internal/repository`：内存数据层
- `internal/model`：核心模型
- `internal/receipt`、`internal/delivery`、`internal/readmodel`、`internal/checkpoint`：容易混淆的回执、已读、投递、同步相关模型
- `internal/legacy`、`internal/compat`、`internal/experiments`、`internal/history`、`internal/migration`、`internal/noise`：历史兼容、迁移、实验和运维上下文

不是所有看起来相关的文件都应该进入核心路径。真实项目中，AI 也经常会把历史逻辑、实验逻辑和核心逻辑混在一起。
