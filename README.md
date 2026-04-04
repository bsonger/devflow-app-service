# Devflow App Service

`devflow-app-service` is the backend owner for `Project` and `Application` metadata.

## Backend Role

- own `Project`
- own `Application`
- maintain `Application.active_manifest` binding
- provide application/project catalog queries for other services and the future platform

## Backend Architecture

This repo uses a **layered metadata-service backend**:

```text
cmd
 -> config
 -> router
 -> api
 -> service
 -> store
 -> model
```

### Package responsibilities

- `cmd/`: service startup
- `pkg/config`: config loading and runtime init
- `pkg/router`: Gin router and middleware wiring
- `pkg/api`: HTTP handlers and status mapping
- `pkg/service`: metadata rules and resource behavior
- `pkg/store`: Mongo access
- `pkg/model`: `Project` / `Application` models

## Non-Goals

- no `Manifest` ownership
- no `Release` ownership
- no `Intent` ownership
- no `Configuration` ownership
- no verify ingress
- no Tekton / Argo execution orchestration

## Key Docs

- `docs/architecture.md`
- `docs/api-spec.md`
- `docs/constraints.md`
- `docs/resources/README.md`

## Local Run

- `go run ./cmd`
- `go build ./cmd/main.go`
- `go test ./...`
