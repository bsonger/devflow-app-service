# Environment

## Purpose

`Environment` defines a deploy target identity used by configuration, runtime, and release flows.
It is part of the current app-owned metadata model even though the repo does not yet expose standalone public `Environment` CRUD/list endpoints in `docs/api-spec.md`.

## Field table

| Field | Type | Required | Writable | Description |
|---|---|---|---|---|
| `name` | `string` | required | user | 环境名 |
| `cluster` | `string` | optional | user | 集群名 |
| `namespace` | `string` | optional | user | 命名空间 |
| `description` | `string` | optional | user | 环境描述 |
| `labels` | `map[string]string` | optional | user | 扩展标签 |
