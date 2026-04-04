# Application

## Ownership

- owner repo: `devflow-app-service`
- authoritative model file: `pkg/model/application.go`
- authoritative API doc: `docs/api-spec.md`
- swagger source: `docs/swagger.yaml`

## Purpose

`Application` 是应用元数据资源，并维护应用侧的 `active_manifest` 绑定。

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
| `name` | `string` | expected on create | user | 应用名 |
| `project_id` | `*ObjectID` | optional | user | 关联项目 ID |
| `project_name` | `string` | optional/system-synced | user/system | 项目名；有 `project_id` 时由服务同步校正 |
| `repo_url` | `string` | expected in practice | user | 代码仓库地址 |
| `active_manifest_name` | `string` | optional | system/user | 当前绑定的 manifest 名 |
| `active_manifest_id` | `*ObjectID` | optional | system/user | 当前绑定的 manifest ID |
| `replica` | `*int32` | optional | user | 副本数 |
| `type` | `ReleaseType` | expected in practice | user | 发布策略类型 |
| `config_maps` | `[]*ConfigMap` | optional | user | ConfigMap 挂载配置 |
| `service` | `Service` | expected in practice | user | 服务暴露信息 |
| `internet` | `Internet` | expected in practice | user | 内外网类型 |
| `envs` | `map[string][]EnvVar` | optional | user | 环境变量集合 |
| `status` | `string` | optional | user/system | 应用状态 |

## Nested types

### `ReleaseType`
- `normal`
- `canary`
- `blue-green`

### `Internet`
- `internal`
- `external`

### `Service`
- `ports: []Port`

### `Port`
- `name: string`
- `port: int`
- `target_port: int`

### `ConfigMap`
- `name: string`
- `mount_path: string`
- `files_path: map[string]string`

### `EnvVar`
- `name: string`
- `value: string`

## Active manifest binding

`active_manifest` 不是 build/release 编排 owner，只表示应用当前绑定的活动版本。

相关字段：
- `active_manifest_id`
- `active_manifest_name`

专用接口：
- `PATCH /api/v1/applications/{id}/active_manifest`

该接口请求体：
- `manifest_id: string`（required）

附加规则：
- `manifest_id` 必须是合法 ObjectID
- 目标 manifest 必须属于当前 application

## Create / update rules

### Create
- current API behavior:
  - handler 绑定整个 `model.Application`
  - 当前未做大范围字段级 `binding:"required"` 校验
- practical required fields:
  - `name`
  - `repo_url`
  - `type`
  - `service`
  - `internet`
- project reference rule:
  - 若提供 `project_id`，服务会加载 `Project` 并同步 `project_name`
  - 若 `project_name` 与 `project_id` 不匹配，会返回错误
- server-managed fields:
  - `id`, `created_at`, `updated_at`

### Update
- mutable fields:
  - 大部分业务字段可更新
- immutable/system-managed fields:
  - `id`, `created_at`, `deleted_at`
- special update path:
  - `active_manifest` 建议通过专用 patch 接口更新

## Validation notes

- `project_id` 若提供，必须引用存在的 `Project`
- `manifest_id` patch 时必须引用属于该应用的 manifest
- 列表接口对 `project_id` 的过滤要求合法 ObjectID

## Source pointers

- router: `pkg/router/application.go`
- handler: `pkg/api/application.go`
- service: `pkg/service/application.go`
- model: `pkg/model/application.go`
- manifest reference helper: `pkg/model/manifest.go`
