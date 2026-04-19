# API Spec

## Purpose

`devflow-app-service` defines the converged public metadata API surface for:

- `Project`
- `Application`
- `Cluster`
- `Environment`

It also owns the narrow `Application.active_image` binding endpoint.
`Environment` remains an app-owned deploy-target record, but it does not expose a writable `namespace` field.

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

### `Cluster`
- `POST /api/v1/clusters`
- `GET /api/v1/clusters`
- `GET /api/v1/clusters/{id}`
- `PUT /api/v1/clusters/{id}`
- `DELETE /api/v1/clusters/{id}`

### `Environment`
- `POST /api/v1/environments`
- `GET /api/v1/environments`
- `GET /api/v1/environments/{id}`
- `PUT /api/v1/environments/{id}`
- `DELETE /api/v1/environments/{id}`

## Request Rules

- all list endpoints use `page` and `page_size`
- `GET /api/v1/clusters` accepts optional `name`
- `GET /api/v1/environments` accepts optional `name` and `cluster_id`
- `PATCH /api/v1/applications/{id}/active_image` requires `image_id`
- `active_image` only represents the currently bound image, not build/release orchestration state
- `Application` owns stable app metadata such as `project_id`, `name`, `repo_address`, `labels`, and `active_image_id`
- `Cluster` owns deploy-target server identity plus sensitive connection material such as `kubeconfig`
- `Environment` owns deploy-target naming and `cluster_id` selection, but it must not accept or persist a writable `namespace`
- app-service does not own environment variables, images, releases, configuration revisions, or verification records

## Response Rules

- create endpoints return `201` with `{ "data": ... }`
- get endpoints return `200` with `{ "data": ... }`
- list endpoints return `{ "data": [...], "pagination": { "page", "page_size", "total" } }`
- `PUT`, `PATCH`, and `DELETE` return `204 No Content`
- project/application/cluster/environment payloads include `created_at` and `updated_at`
- `Environment` payloads include `cluster_id` and never include `namespace`
- `labels` use ordered `[{ "key": "team", "value": "platform" }]`
- legacy map-shaped labels are still tolerated when reading older rows from storage, but new writes use the array contract
- `404` means the resource does not exist
- `400` means invalid request input or invalid reference
- `409` means state or uniqueness conflict
- delete behavior is soft-delete oriented

## Error Rules

- invalid ID or malformed request body -> `400`
- missing required cluster/environment fields -> `400`
- invalid referenced `cluster_id` -> `400`
- resource not found -> `404`
- unique-name conflict -> `409`
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
