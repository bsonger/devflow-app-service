# Architecture

## Purpose

`devflow-app-service` is the metadata owner for:

- `Project`
- `Environment`
- `Application`
- `ServiceResource`

It provides project/environment/application relationships, application repository identity, static service metadata, and the narrow `active_image` binding.

## Architecture Style

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
- `Application` 1 -> N `ServiceResource`
- `Application.repo_address` is the unified repository locator
- `ServiceResource` stores `description`, `labels`, and `ports`

## Request Flow

```text
Client
  -> router
  -> project/application/service-resource handler
  -> metadata service logic
  -> persistence store
  -> HTTP response
```

## Internal Package Layout

- `cmd/main.go`
  - process entrypoint only
- `pkg/infra/config`
  - config loading
  - runtime initialization
- `pkg/router`
  - route registration
  - middleware wiring
- `pkg/api`
  - project/application/service-resource handlers
- `pkg/app`
  - metadata behavior
  - `active_image` binding rules
- `pkg/infra/store`
  - repo-owned metadata persistence
- `pkg/domain`
  - `Project`, `Application`, `ServiceResource`

## External Dependencies

- `Gin`
- PostgreSQL persistence
- `devflow-service-common`

## Non-Goals

- `Manifest`
- `Release`
- `Intent`
- `Configuration`
- `ConfigurationRevision`
- environment-variable ownership
- verify ingress
- Tekton / Argo / Kubernetes execution orchestration

## Swagger generation

- `Dockerfile` runs `swag init -g cmd/main.go --parseDependency -o docs/generated/swagger` before building.
- Keep the generated bundle under `docs/generated/swagger`; rerun `scripts/regen-swagger.sh` when handlers change.
- `scripts/build.sh` wraps regeneration plus `go build` for locals, and `scripts/export_service_repo.sh` always copies the generated folder.
