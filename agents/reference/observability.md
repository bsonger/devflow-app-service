# Observability Reference

## Purpose

约束 `app-service` 的统一观测方式。

## Scope

适用于：

- 所有入站 HTTP 链路
- 所有对 Mongo 以及其他外部系统的出站调用
- 所有异步 worker / controller / executor 路径

不适用于：

- 非服务化的一次性本地脚本
- 与当前仓库无关的外部观测平台实现细节

## Must

- `trace` 统一使用 OpenTelemetry
- `metrics` 统一使用 OpenTelemetry
- `log` 必须是结构化日志，并注入 `trace_id` / `span_id`
- Go 服务必须保留 `pprof`，集群环境建议接入 Pyroscope
- 服务统一向 OTel Collector 输出 OTLP，不允许各服务直连不同观测厂商 exporter
- 所有 HTTP / gRPC / async worker 调用必须传递 trace context
- 发起下游调用时创建 client span，接收请求时创建 server span
- 跨异步边界时至少保留 `trace_id` 和业务主键，必要时用 span link 衔接
- 必须观测的关键链路包括：
  - `GET /api/v1/projects`
  - `POST /api/v1/projects`
  - `GET /api/v1/projects/:id`
  - `PUT /api/v1/projects/:id`
  - `DELETE /api/v1/projects/:id`
  - `GET /api/v1/projects/:id/applications`
  - `POST /api/v1/applications`
  - `GET /api/v1/applications/:id`
  - `PUT /api/v1/applications/:id`
  - `DELETE /api/v1/applications/:id`
  - `PATCH /api/v1/applications/:id/active_manifest`
  - 所有对 Mongo 和其他外部系统的出站调用
- span 名称必须稳定，业务 ID 放入 span attribute 而不是 span 名称
- 错误必须记录为 span status 和 error event
- 所有服务至少提供：
  - 入站请求数、耗时、错误数
  - 出站请求数、耗时、错误数
  - Mongo 调用耗时和错误数
- 日志至少包含：
  - `ts`
  - `level`
  - `msg`
  - `service`
  - `trace_id`
  - `span_id`
  - `request_id`
- 涉及控制面资源时，日志追加 `project_id`、`application_id`、`active_manifest_id`

## Must Not

- 不让每个服务直接接不同 exporter
- 不把 `project_id`、`application_id` 这类高基数字段放进 metrics label
- 不把日志当成状态真相来源
- 不把 token、secret、kubeconfig、完整 webhook body 直接写日志
- 不把每个 debug 字段都做成正式指标
- 不把 `/metrics`、`/healthz`、`/readyz`、`/debug/pprof/*` 计入业务指标

## Outputs

- 可关联的 traces、logs、metrics、profiles
- 能从 `project_id` / `application_id` 追踪到数据库写入链路
- 明确的 OTel Collector 输出路径

## Pass/Fail

- `Pass`：关键链路可追踪，日志可关联，指标低基数，profiling 可按需启用
- `Fail`：trace 断裂、日志不可关联、metrics 高基数失控，或观测出口分裂
