# Cluster

## Ownership

- owner repo: `devflow-app-service`
- authoritative model file: `pkg/domain/cluster.go`
- authoritative API doc: `docs/api-spec.md`
- generated swagger: `docs/generated/swagger/swagger.yaml`

## Purpose

`Cluster` is the app-owned deploy target resource for Kubernetes cluster metadata.
It carries the destination server, sensitive connection material, and Argo CD naming data needed by later deploy-target resolution and onboarding work.
The public owner API exposes standalone CRUD/list endpoints at `/api/v1/clusters`.

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
| `name` | `string` | required | user | 集群显示名 |
| `server` | `string` | required | user | Kubernetes API server 地址 |
| `kubeconfig` | `string` | required | user/secret | 集群连接配置；敏感字段 |
| `argocd_cluster_name` | `string` | optional | user | Argo CD 中使用的 cluster 标识 |
| `description` | `string` | optional | user | 集群描述 |
| `labels` | `[]LabelItem` | optional | user | 扩展标签 |

## Create / update rules

### Create
- practical required fields:
  - `name`
  - `server`
  - `kubeconfig`
- server-managed fields:
  - `id`
  - `created_at`
  - `updated_at`

### Update
- mutable fields:
  - `name`, `server`, `kubeconfig`, `argocd_cluster_name`, `description`, `labels`
- immutable/system-managed fields:
  - `id`, `created_at`, `deleted_at`

## API surface

- `POST /api/v1/clusters`
- `GET /api/v1/clusters`
- `GET /api/v1/clusters/{id}`
- `PUT /api/v1/clusters/{id}`
- `DELETE /api/v1/clusters/{id}`

## Validation notes

- `name`、`server`、`kubeconfig` 不能为空
- `server` 必须是可解析的目标 cluster API 地址
- `kubeconfig` 属于敏感连接材料，不应出现在平台读模型或普通日志中
- `argocd_cluster_name` 用于后续 Argo CD cluster 对应关系，允许先为空后补齐
- 重复 `name` 会返回冲突错误而不是静默覆盖

## Source pointers

- router: `pkg/router/cluster.go`
- handler: `pkg/api/cluster.go`
- service: `pkg/app/cluster.go`
- model: `pkg/domain/cluster.go`
