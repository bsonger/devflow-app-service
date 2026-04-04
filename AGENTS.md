# Repository Guidelines

## Boundary

- This repository is `devflow-app-service`.
- Public surface is `Project` and `Application` only.
- Do not reintroduce `Manifest`, `Release`, `Intent`, `Configuration`, or `Verify` routes, models, router modules, or runtime/bootstrap logic.
- `Application` may keep `active_manifest` writeback semantics, but only as part of the application boundary.

## Structure

- `cmd/main.go` uses shared bootstrap from `../devflow-service-common`.
- `pkg/api/` contains project and application handlers.
- `pkg/service/` contains project and application CRUD logic.
- `pkg/router/` contains app-only routes and middleware.
- `pkg/config/` initializes config, observability, Mongo, and local store state.
- `docs/` contains the repository-level architecture, API, constraints, observability, and harness docs.

## Required Rules

- Any outbound service or external call must emit `metrics + trace + structured log`.
- Do not add high-cardinality business IDs to metrics labels.
- Default harness is `Planner -> Generator -> Evaluator`.
- When the runtime supports delegation, the harness must spawn those roles as separate sub-agents.
- Non-trivial work should use a run directory under `agents/runs/`.

## Doc And API Hygiene

- Regenerate Swagger after route or handler changes.
- Keep `README.md`, `AGENTS.md`, `agents/protocols/startup.md`, and `docs/*.md` aligned with the actual boundary.
- Do not reintroduce dead service/model/router/bootstrap files.
