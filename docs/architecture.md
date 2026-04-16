# Architecture

## Purpose

`devflow-app-service` is the metadata owner for:

- `Project`
- `Environment`
- `Application`

It provides project/environment/application relationships, application repository identity, shared environment vocabulary, and the narrow `active_image` binding.
Its current public API surface remains intentionally narrower than the full metadata model and exposes only `Project`, `Application`, and the `active_image` binding.

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
- `Environment` defines a stable deploy-target identity reused by runtime and release flows
- `Application.repo_address` is the unified repository locator

## Request flow

```text
Client
  -> router
  -> project/application handler
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
  - project/application handlers
- `pkg/app`
  - metadata behavior
  - `active_image` binding rules
- `pkg/infra/store`
  - repo-owned metadata persistence
- `pkg/domain`
  - `Project`, `Environment`, `Application`

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
