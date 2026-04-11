# Observability

## Purpose

`devflow-app-service` emits the shared backend telemetry baseline plus project/application metadata context.

## Logs

Required structured fields:
- `resource`
- `resource_id`
- `project_id`
- `application_id`
- `result`
- `error_code`

## Metrics

- use shared `devflow_http_*` ingress metrics
- do not add ID-like labels to custom metrics
- exclude `/metrics`, `/healthz`, `/readyz`, and `/debug/pprof/*` from business HTTP telemetry

## Tracing

- every business HTTP request should create a server span
- downstream calls, if added later, must emit client spans with propagated trace context
- resource-scoped attributes should prefer project/application identifiers over free-form text

## Health and readiness

- expose `/healthz`, `/readyz`, and `/metrics`
- keep Swagger endpoints and diagnostics outside business metrics aggregation

## Failure modes

Watch for:
- project/application CRUD storage failures
- invalid `active_image` binding requests
- stale metadata assumptions leaking into downstream services

## Dashboards and runbooks

Use the shared backend dashboard/runbook set until repo-specific views exist.
