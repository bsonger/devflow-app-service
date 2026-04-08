# Observability

## Shared Baseline

This repo follows the shared telemetry contract implemented in `devflow-service-common`.

- structured logs with shared runtime fields
- `devflow_http_*` ingress metrics
- standard server/client spans with service-defined business attributes
- optional diagnostics only for `pprof` and Pyroscope

## Repo-Local Focus

`devflow-app-service` should add resource context for:

- `project`
- `application`
- `active_image`

Recommended structured fields:

- `resource`
- `resource_id`
- `project_id`
- `application_id`
- `result`
- `error_code`

## Metrics Notes

- do not add ID-like labels to custom metrics
- `/metrics`, `/healthz`, `/readyz`, and `/debug/pprof/*` are excluded from business HTTP telemetry

## Profile

- `pprof` is disabled by default
- Pyroscope is disabled by default
- both are enabled only through explicit runtime configuration
