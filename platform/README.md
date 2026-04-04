# Platform Notes

This repository only owns the `devflow-app-service` boundary.

Runtime shape:

- `cmd/main.go` uses shared bootstrap from `../devflow-service-common`
- `pkg/router/` exposes only project/application routes
- `pkg/api/project.go` and `pkg/api/application.go` are the only HTTP handler surfaces
- `pkg/service/` only contains app metadata logic

Shared infra:

- pagination and response helpers come from `devflow-service-common/httpx`
- middleware and telemetry helpers come from `devflow-service-common/routercore` and `devflow-service-common/observability`

Operational rules:

- outbound service or external calls must emit `metrics + trace + structured log`
- `Planner -> Generator -> Evaluator` is the default harness
- when delegation is supported, sub-agents must be spawned
