# Project

## Ownership

- owner repo: `devflow-app-service`
- authoritative model file: `pkg/domain/project.go`
- authoritative API doc: `docs/api-spec.md`
- swagger source: `docs/generated/swagger/swagger.yaml`

## Purpose

`Project` 是顶层空间资源，用于组织 `Application` 元数据。

## Common base fields

| Field | Type | Required | Writable | Description |
|---|---|---|---|---|
| `id` | `uuid.UUID` | server-generated | no | 主键 |
| `created_at` | `time.Time` | server-generated | no | 创建时间 |
| `updated_at` | `time.Time` | server-generated | no | 更新时间 |
| `deleted_at` | `*time.Time` | optional | system-managed | 软删除时间 |

## Field table

| Field | Type | Required | Writable | Description |
|---|---|---|---|---|
| `name` | `string` | expected on create | user | 项目名 |
| `description` | `string` | optional | user | 项目描述 |
| `labels` | `map[string]string` | optional | user | 扩展标签 |

## Create / update rules

### Create
- practical required fields:
  - `name`
- server-managed fields:
  - `id`
  - `created_at`
  - `updated_at`

### Update
- mutable fields:
  - `name`, `description`, `labels`
- immutable/system-managed fields:
  - `id`, `created_at`, `deleted_at`

## Validation notes

- 当前 handler 没有单独的字段级 required 校验
- `id` 必须是合法 UUID 才能用于读取/更新/删除
- 软删除记录默认不会出现在列表里，除非显式包含 deleted 数据

## Source pointers

- router: `pkg/router/project.go`
- handler: `pkg/api/project.go`
- service: `pkg/app/project.go`
- model: `pkg/domain/project.go`
