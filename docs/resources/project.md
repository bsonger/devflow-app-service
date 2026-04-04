# Project

## Ownership

- owner repo: `devflow-app-service`
- authoritative model file: `pkg/model/project.go`
- authoritative API doc: `docs/api-spec.md`
- swagger source: `docs/swagger.yaml`

## Purpose

`Project` 是顶层空间资源，用于组织 `Application` 元数据。

## Common base fields

| Field | Type | Required | Writable | Description |
|---|---|---|---|---|
| `id` | `ObjectID` | server-generated | no | 主键 |
| `created_at` | `time.Time` | server-generated | no | 创建时间 |
| `updated_at` | `time.Time` | server-generated | no | 更新时间 |
| `deleted_at` | `*time.Time` | optional | system-managed | 软删除时间 |

## Field table

| Field | Type | Required | Writable | Description |
|---|---|---|---|---|
| `name` | `string` | expected on create | user | 项目名 |
| `key` | `string` | expected on create | user | 项目标识 |
| `description` | `string` | optional | user | 项目描述 |
| `namespace` | `string` | optional | user | 命名空间；为空时默认等于 `name` |
| `owner` | `string` | optional | user | 负责人 |
| `labels` | `map[string]string` | optional | user | 扩展标签 |
| `status` | `ProjectStatus` | system-defaulted | user/system | 状态；为空时默认 `active` |

## Lifecycle / status fields

- status field: `status`
- valid values:
  - `active`
  - `archived`
- defaults:
  - create 时为空会被 `ApplyDefaults()` 设置为 `active`
- delete behavior:
  - 删除走软删除，并会把 `status` 设为 `archived`

## Create / update rules

### Create
- current API behavior:
  - handler 会绑定整个 `model.Project`
  - 当前未做额外字段级 `binding:"required"` 校验
- practical required fields:
  - `name`
  - `key`
- server-managed fields:
  - `id`
  - `created_at`
  - `updated_at`
  - `status` 默认值
  - `namespace` 默认值

### Update
- mutable fields:
  - `name`, `key`, `description`, `namespace`, `owner`, `labels`, `status`
- immutable/system-managed fields:
  - `id`, `created_at`, `deleted_at`
- special behavior:
  - 项目名变化时，会同步更新关联 `Application.project_name`

## Validation notes

- 当前 handler 没有单独的字段级 required 校验
- `id` 必须是合法 ObjectID 才能用于读取/更新/删除
- 软删除记录默认不会出现在列表里，除非显式包含 deleted 数据

## Source pointers

- router: `pkg/router/project.go`
- handler: `pkg/api/project.go`
- service: `pkg/service/project.go`
- model: `pkg/model/project.go`
