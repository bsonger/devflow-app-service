# App Service Platform Notes

## Purpose

This file is the repo-local runtime note for `devflow-app-service`.
For public API shape, ownership, and resource details, prefer:
- `../README.md`
- `../docs/`
- `../docs/resources/`

## Runtime entrypoints

- process entry: `cmd/main.go`
- shared bootstrap: `../devflow-service-common/bootstrap`
- router root: `pkg/router/router.go`

## Main local code paths

- project routes: `pkg/router/project.go`
- application routes: `pkg/router/application.go`
- project handlers: `pkg/api/project.go`
- application handlers: `pkg/api/application.go`
- project logic: `pkg/service/project.go`
- application logic: `pkg/service/application.go`

## Platform dependencies

- shared response / pagination: `devflow-service-common/httpx`
- shared middleware: `devflow-service-common/routercore`
- shared observability: `devflow-service-common/observability`

## Service identity

- OTel `service.name`: `app-service`
- typical ports:
  - `APP_SERVICE_PORT`
  - `APP_SERVICE_METRICS_PORT`
  - `APP_SERVICE_PPROF_PORT`
