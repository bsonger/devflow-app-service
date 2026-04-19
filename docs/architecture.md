# Architecture

## Purpose

`devflow-app-service` is the metadata owner for:

- `Project`
- `Application`
- `Cluster`
- `Environment`

It provides project/application relationships, shared deploy-target vocabulary, cluster destination metadata, and the narrow `active_image` binding.
Its public API surface now exposes CRUD/list endpoints for `Project`, `Application`, `Cluster`, and `Environment`, plus the dedicated `active_image` binding.

## Architecture style

This repo uses a **layered metadata-service backend**:

```text
router -> api -> app -> infra/store
                \-> domain
```

Where:
- `api` binds HTTP requests and maps status codes
- `app` owns metadata rules and cross-resource checks
- `infra/store` persists repo-owned metadata

The converged target resource model is:

- `Project` 1 -> N `Application`
- `Cluster` defines the deploy-target Kubernetes API server and connection material
- `Environment` defines deploy semantics and selects one `Cluster`
- `Application.repo_address` is the unified repository locator

## Request flow

```text
Client
  -> router
  -> project/application/cluster/environment handler
  -> metadata service logic
  -> persistence store
  -> HTTP response
```

## Internal package layout

- `cmd/main.go`
  - process entrypoint only
- `pkg/infra/config`
  - config loading
  - runtime initialization
- `pkg/router`
  - route registration
  - middleware wiring
- `pkg/api`
  - project/application/cluster/environment handlers
- `pkg/app`
  - metadata behavior
  - cross-resource reference validation
  - `active_image` binding rules
- `pkg/infra/store`
  - repo-owned metadata persistence
- `pkg/domain`
  - `Project`, `Application`, `Cluster`, `Environment`

## External dependencies

- `Gin`
- PostgreSQL persistence
- `devflow-service-common`

## Swagger generation

- `Dockerfile` runs `swag init -g cmd/main.go --parseDependency -o docs/generated/swagger` before building.
- Keep the generated bundle under `docs/generated/swagger`; rerun `scripts/regen-swagger.sh` when handlers change.
- `scripts/build.sh` wraps regeneration plus `go build` for locals, and `scripts/export_service_repo.sh` always copies the generated folder.

## Non-goals

- `Image`
- `Release`
- `Intent`
- `Configuration`
- `ConfigurationRevision`
- environment-variable ownership
- verify ingress
- Tekton / Argo / Kubernetes execution orchestration
