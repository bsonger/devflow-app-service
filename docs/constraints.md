# Constraints

## Ownership

- `Project` is the top-level workspace resource.
- `Application` belongs to one `Project`.
- `Application.active_image` may only reference the current active image binding for that application.

## Hard constraints

- do not introduce `Manifest`, `Release`, `Intent`, `Configuration`, or `Verify` as public resources in this repo
- do not move execution-state writeback behavior into app-service
- do not copy other service domain models into app-service just for convenience

## Data rules

- deletion must follow the repo's existing soft-delete semantics
- list and detail responses must honor soft-delete filtering rules
- `active_image` updates must remain idempotent

## Dependency rules

- any outbound dependency call must emit metrics, traces, and structured logs together
- resource identifiers may appear in logs and trace attributes, but not as metric labels

## Non-goals

- merging execution-plane objects into `Application`
- building direct database update logic inside HTTP handlers
- treating historical Swagger entries as the current ownership boundary
