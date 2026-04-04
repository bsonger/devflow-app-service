# Devflow App Service

`devflow-app-service` 只负责 `Project` 和 `Application`。

边界：

- 仅保留 `Project`、`Application` 的 HTTP、service、model、store 和配置加载
- `Application` 允许维护 `active_manifest` 关系
- 不提供 `Manifest`、`Job`、`Intent`、`Configuration`、`Verify` 对外资源
- 启动、路由中间件、分页、响应和观测基础设施来自 `../devflow-service-common`

仓库文档：

- [架构](docs/architecture.md)
- [接口规范](docs/api-spec.md)
- [约束](docs/constraints.md)
- [观测规范](docs/observability.md)
- [Harness](docs/harness.md)

运行约定：

- 任何调用其他服务或外部系统的代码都必须同时产出 `metrics + trace + structured log`
- 默认 harness 为 `Planner -> Generator -> Evaluator`
- 运行时支持 delegation 时，必须真实启动对应 sub-agent，不允许只在单 agent 内口头模拟

常用命令：

- `go run ./cmd`
- `go build ./cmd/main.go`
- `go test ./...`
- `swag init -g cmd/main.go --parseDependency`
