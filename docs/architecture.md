# Architecture

## Purpose

`devflow-app-service` is the metadata owner for `Project` and `Application`.
It provides application identity, project/application relationships, and active-manifest binding metadata.

## Architecture Style

This repo uses a **layered metadata-service backend**:

```text
router -> api -> service -> store
                    \-> model
```

Where:
- `api` binds HTTP requests and maps status codes
- `service` owns metadata rules and cross-resource checks
- `store` persists repo-owned metadata in Mongo

## Request Flow

```text
Client
  -> router
  -> project/application handler
  -> project/application service
  -> Mongo store
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
  - project/application handlers
- `pkg/service`
  - metadata behavior
  - `active_manifest` binding rules
- `pkg/store`
  - Mongo access
- `pkg/model`
  - `Project` and `Application` models

## External Dependencies

- `Gin`
- `MongoDB`
- `devflow-service-common`

## Non-Goals

- `Manifest`
- `Release`
- `Intent`
- `Configuration`
- verify ingress
- Tekton / Argo / Kubernetes execution orchestration
