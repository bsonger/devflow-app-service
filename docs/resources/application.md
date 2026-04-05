# Application

## Ownership

- owner repo: `devflow-app-service`
- authoritative model file: `pkg/model/application.go`
- authoritative API doc: `docs/api-spec.md`
- generated swagger: `docs/swagger.yaml` (transitional; still reflects legacy handler layer until API migration)

## Purpose

`Application` 是应用元数据资源，并维护应用侧的 `active_manifest` 绑定。
服务暴露信息已经拆分到独立的 `ServiceResource` 子资源。

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
| `project_id` | `uuid.UUID` | required | user | 关联项目 ID |
| `name` | `string` | required | user | 应用名 |
| `repo_address` | `string` | required | user | 代码仓库地址 |
| `active_manifest_id` | `*uuid.UUID` | optional | system/user | 当前绑定的 manifest ID |
| `replica` | `*int32` | optional | user | 副本数 |
| `type` | `ReleaseType` | required | user | 发布策略类型 |
| `status` | `string` | optional | user/system | 应用状态 |

## Nested types

### `ReleaseType`
- `normal`
- `canary`
- `blue-green`

### `Internet`
- `internal`
- `external`

### `Port`
- `name: string`
- `port: int`
- `target_port: int`

## Related child resource: `ServiceResource`

`Application` 不再直接承载服务暴露字段；相关信息迁移为独立子资源：

| Field | Type | Description |
|---|---|---|
| `application_id` | `uuid.UUID` | 所属应用 |
| `name` | `string` | 服务资源名 |
| `internet` | `Internet` | 内外网属性 |
| `ports` | `[]Port` | 端口集合 |
| `status` | `string` | 服务资源状态 |

## Active manifest binding

`active_manifest` 不是 build/release 编排 owner，只表示应用当前绑定的活动版本。

相关字段：
- `active_manifest_id`

专用接口：
- `PATCH /api/v1/applications/{id}/active_manifest`

该接口请求体：
- `manifest_id: string`（required，当前 handler 层仍以字符串承载）

附加规则：
- `manifest_id` 必须引用属于当前 application 的 manifest

## Create / update rules

### Create
- target relational contract:
  - required: `project_id`, `name`, `repo_address`, `type`
  - `project_id` 必须引用存在的 `Project`
- server-managed fields:
  - `id`, `created_at`, `updated_at`

### Update
- mutable fields:
  - `name`, `repo_address`, `active_manifest_id`, `replica`, `type`, `status`
- immutable/system-managed fields:
  - `id`, `created_at`, `deleted_at`
- special update path:
  - `active_manifest` 建议通过专用 patch 接口更新

## Validation notes

- `project_id` 必须引用存在的 `Project`
- `manifest_id` patch 时必须引用属于该应用的 manifest
- 服务暴露相关校验由 `ServiceResource` 独立承担

## Source pointers

- router: `pkg/router/application.go`
- handler: `pkg/api/application.go`
- service: `pkg/service/application.go`
- model: `pkg/model/application.go`
- manifest reference helper: `pkg/model/manifest.go`
