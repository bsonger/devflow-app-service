# Environment

## Ownership

- owner repo: `devflow-app-service`
- authoritative model file: `pkg/domain/environment.go`
- authoritative API doc: `docs/api-spec.md`
- generated swagger: `docs/generated/swagger/swagger.yaml`

## Purpose

`Environment` defines deploy semantics plus the selected target cluster.
It is part of the app-owned metadata model and now exposes standalone public CRUD/list endpoints at `/api/v1/environments`.
`Environment` does not store or accept a user-managed namespace field.

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
| `name` | `string` | required | user | 环境名 |
| `cluster_id` | `uuid.UUID` | required | user | 目标集群 ID |
| `description` | `string` | optional | user | 环境描述 |
| `labels` | `[]LabelItem` | optional | user | 扩展标签 |

## Cross-resource references

| Field | Refers to | Owning repo | Notes |
|---|---|---|---|
| `cluster_id` | `Cluster.id` | `devflow-app-service` | 绑定环境到具体 deploy target cluster |

## Create / update rules

### Create
- practical required fields:
  - `name`
  - `cluster_id`
- server-managed fields:
  - `id`
  - `created_at`
  - `updated_at`

### Update
- mutable fields:
  - `name`, `cluster_id`, `description`, `labels`
- immutable/system-managed fields:
  - `id`, `created_at`, `deleted_at`
- forbidden fields:
  - `namespace`

## API surface

- `POST /api/v1/environments`
- `GET /api/v1/environments`
- `GET /api/v1/environments/{id}`
- `PUT /api/v1/environments/{id}`
- `DELETE /api/v1/environments/{id}`

## Validation notes

- `name` 和 `cluster_id` 不能为空
- `cluster_id` 必须引用存在的 `Cluster`
- `namespace` 不是 `Environment` 输入字段；deploy-target namespace 由 normalized `project.name + environment.name` 推导
- 当 `environment.name == production` 时，推导结果为 normalized `project.name`
- 非生产环境推导结果为 normalized `project.name-environment.name`
- 重复 `name` 会返回冲突错误而不是静默覆盖

## Source pointers

- router: `pkg/router/environment.go`
- handler: `pkg/api/environment.go`
- service: `pkg/app/environment.go`
- model: `pkg/domain/environment.go`
