# DevFlow App Service

`devflow-app-service` is the backend owner for `Project` and `Application` metadata.

## Backend Role

- own `Project`
- own `Application`
- maintain `Application.active_manifest` binding
- provide application/project catalog queries for other services and the future platform

## Local Run

- `go run ./cmd`
- `go build ./cmd/main.go`
- `go test ./...`

## Key Docs

- `docs/architecture.md`
- `docs/api-spec.md`
- `docs/constraints.md`
- `docs/resources/README.md`
