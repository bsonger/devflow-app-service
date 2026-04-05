# Architecture

## Purpose

`devflow-app-service` is the metadata owner for:

- `Project`
- `Application`
- `ServiceResource`

It provides project/application relationships, application repository identity, service-exposure metadata, and the narrow `active_manifest` binding.

## Architecture Style

This repo uses a **layered metadata-service backend**:

```text
router -> api -> service -> store
                    \-> model
```

Where:
- `api` binds HTTP requests and maps status codes
- `service` owns metadata rules and cross-resource checks
- `store` persists repo-owned metadata

The converged target resource model is:

- `Project` 1 -> N `Application`
- `Application` 1 -> N `ServiceResource`

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
- `pkg/config`
  - config loading
  - runtime initialization
- `pkg/router`
  - route registration
  - middleware wiring
- `pkg/api`
  - project/application/service-resource handlers
- `pkg/service`
  - metadata behavior
  - `active_manifest` binding rules
- `pkg/store`
  - repo-owned metadata persistence
- `pkg/model`
  - `Project`, `Application`, `ServiceResource`

## External Dependencies

- `Gin`
- PostgreSQL target persistence
- `devflow-service-common`

## Non-Goals

- `Manifest`
- `Release`
- `Intent`
- `Configuration`
- environment-variable ownership
- verify ingress
- Tekton / Argo / Kubernetes execution orchestration
