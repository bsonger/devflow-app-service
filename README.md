# DevFlow App Service

`devflow-app-service` owns `Project`, `Application`, `Cluster`, and `Environment` metadata.

## Role

- own `Project`
- own `Application`
- own `Cluster`
- own `Environment`
- maintain `Application.active_image` binding
- provide application/project/deploy-target catalog queries for other services and the platform console

## Current contract highlights

- `Project.labels`, `Application.labels`, `Cluster.labels`, and `Environment.labels` use ordered `[{ key, value }]`
- responses include `created_at` and `updated_at`
- `Environment` uses `cluster_id` and does not expose a writable `namespace`
- `Cluster` owns destination server and sensitive connection material such as `kubeconfig`
- `Application.active_image_id` remains the only app-service release-facing binding outside the metadata CRUD surface

## Key Commands

- `go run ./cmd`
- `go build ./cmd/main.go`
- `go test ./...`
- Swagger UI: `/swagger/index.html`
- Staging Swagger UI: `/api/v1/app/swagger/index.html`

## Key Docs

- `docs/README.md`
- `scripts/README.md`
- `docs/architecture.md`
- `docs/constraints.md`
- `docs/observability.md`
- `docs/api-spec.md`
- `docs/resources/README.md`
- `docs/generated/swagger/swagger.yaml`
