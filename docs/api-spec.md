# API Spec

## Purpose

`devflow-app-service` only exposes public HTTP APIs for `Project` and `Application`.

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

## Request Rules

- all list endpoints support pagination parameters consistent with Swagger
- `PATCH /api/v1/applications/{id}/active_manifest` requires `manifest_id`
- `active_manifest` only represents the currently bound manifest, not build/release orchestration state

## Response Rules

- list endpoints follow the common pagination shape used in this repo
- `404` means the resource does not exist
- `400` means invalid request input or invalid reference
- `409` is reserved for state or boundary conflicts
- delete behavior is soft-delete oriented

## Error Rules

- invalid ObjectID or malformed request body -> `400`
- resource not found -> `404`
- state/reference conflict -> `409`
- internal/storage error -> `500`

## Non-Goals

This repo does not expose public CRUD for:
- `Manifest`
- `Release`
- `Intent`
- `Configuration`
- `Verify`
