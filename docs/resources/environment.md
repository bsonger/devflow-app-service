# Environment

## Purpose

`Environment` defines a deploy target identity used by configuration, runtime, and release flows.

## Field table

| Field | Type | Required | Writable | Description |
|---|---|---|---|---|
| `name` | `string` | required | user | 环境名 |
| `cluster` | `string` | optional | user | 集群名 |
| `namespace` | `string` | optional | user | 命名空间 |
| `description` | `string` | optional | user | 环境描述 |
| `labels` | `map[string]string` | optional | user | 扩展标签 |
