# DevFlow App Service

`devflow-app-service` owns `Project` and `Application` metadata.

## Current contract highlights

- `Project.labels` and `Application.labels` now use ordered `[{ key, value }]`
- responses include `created_at` and `updated_at`
- `Application.active_image_id` remains the only app-service release-facing binding
- environment-specific console views are **not** owned here

## Backend Role

- own `Project`
- own `Application`
- maintain `Application.active_image` binding
- provide application/project catalog queries for other services and the platform console

## Local Run

- `go run ./cmd`
- `go build ./cmd/main.go`
- `go test ./...`

## Key Docs

- `docs/architecture.md`
- `docs/api-spec.md`
- `docs/constraints.md`
- `docs/resources/README.md`
