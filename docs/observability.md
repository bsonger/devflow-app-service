# 观测规范

## 总则

本仓库要求 `metrics`、`log`、`trace`、`profile` 统一落地。

## 必须记录的场景

- 每个入站 HTTP 请求
- 每个出站服务调用
- 每个外部系统调用

## HTTP 指标

- 请求计数
- 请求耗时
- 错误计数

## Trace 约定

- 入站请求必须有 server span
- 出站调用必须有 client span
- span attribute 中可包含资源 ID、状态和目标服务名

## 日志约定

- 必须是结构化日志
- 至少包含 `service`、`trace_id`、`span_id`、`request_id`
- 涉及资源时可追加 `project_id`、`application_id`、`active_manifest_id`

## Profile

- 保留 `pprof`
- 仅在需要诊断性能问题时启用

## 禁止项

- 不要把高基数字段放进 metrics label
- 不要把 `/metrics`、`/healthz`、`/readyz`、`/debug/pprof/*` 当作业务流量
