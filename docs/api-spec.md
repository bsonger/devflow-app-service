# API Spec

## Purpose

`devflow-app-service` defines the converged public metadata API surface for:

- `Project`
- `Application`
- `ServiceResource`

It also owns the narrow `Application.active_manifest` binding endpoint.

## Endpoint Groups

### `Project`
- `POST /api/v1/projects`
- `GET /api/v1/projects`
- `GET /api/v1/projects/{id}`
- `PUT /api/v1/projects/{id}`
- `DELETE /api/v1/projects/{id}`
- `GET /api/v1/projects/{id}/applications`

### `Application`
- `POST /api/v1/applications`
- `GET /api/v1/applications`
- `GET /api/v1/applications/{id}`
- `PUT /api/v1/applications/{id}`
- `DELETE /api/v1/applications/{id}`
- `PATCH /api/v1/applications/{id}/active_manifest`

### `ServiceResource`
- `GET /api/v1/applications/{id}/services`
- `POST /api/v1/applications/{id}/services`
- `GET /api/v1/services/{id}`
- `PUT /api/v1/services/{id}`
- `DELETE /api/v1/services/{id}`

## Request Rules

- all list endpoints support pagination parameters consistent with Swagger
- `PATCH /api/v1/applications/{id}/active_manifest` requires `manifest_id`
- `active_manifest` only represents the currently bound manifest, not build/release orchestration state
- `Application` owns stable app metadata such as `project_id`, `name`, `repo_address`, and release strategy type
- `ServiceResource` is a child resource under `Application`; one application may own multiple service resources
- app-service does not own environment variables

## Response Rules

- list endpoints follow the common pagination shape used in this repo
- `404` means the resource does not exist
- `400` means invalid request input or invalid reference
- `409` is reserved for state or boundary conflicts
- delete behavior is soft-delete oriented

## Error Rules

- invalid ID or malformed request body -> `400`
- resource not found -> `404`
- state/reference conflict -> `409`
- internal/storage error -> `500`

## Boundary Note

For repo scope and non-goals, see `docs/architecture.md`.
