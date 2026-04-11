# API Spec

## Purpose

`devflow-app-service` defines the converged public metadata API surface for:

- `Project`
- `Application`

It also owns the narrow `Application.active_image` binding endpoint.

## Swagger

- local UI: `/swagger/index.html`
- app-scoped UI alias: `/api/v1/app/swagger/index.html`
- generated source: `docs/generated/swagger/swagger.yaml`

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
- `PATCH /api/v1/applications/{id}/active_image`

## Request Rules

- all list endpoints use `page` and `page_size`
- `PATCH /api/v1/applications/{id}/active_image` requires `image_id`
- `active_image` only represents the currently bound image, not build/release orchestration state
- `Application` owns stable app metadata such as `project_id`, `name`, `repo_address`, `labels`, and `active_image_id`
- app-service does not own environment variables or environment-specific deployment state
- app-service does not own images, releases, configuration revisions, or verification records

## Response Rules

- create endpoints return `201` with `{ "data": ... }`
- get endpoints return `200` with `{ "data": ... }`
- list endpoints return `{ "data": [...], "pagination": { "page", "page_size", "total" } }`
- `PUT`, `PATCH`, and `DELETE` return `204 No Content`
- project/application payloads include `created_at` and `updated_at`
- `labels` use ordered `[{ "key": "team", "value": "platform" }]`
- legacy map-shaped labels are still tolerated when reading older rows from storage, but new writes use the array contract
- `404` means the resource does not exist
- `400` means invalid request input or invalid reference
- `409` is reserved for state or boundary conflicts
- delete behavior is soft-delete oriented

## Error Rules

- invalid ID or malformed request body -> `400`
- resource not found -> `404`
- state/reference conflict -> `409`
- internal/storage error -> `500`

Error responses use:

```json
{
  "error": {
    "code": "invalid_argument",
    "message": "invalid project_id"
  }
}
```

## Boundary Note

For repo scope and non-goals, see `docs/architecture.md`.

## Swagger Note

Generated Swagger artifacts must stay aligned with the current PostgreSQL-backed API contract. Regenerate them after route, request, or response changes.
