# App Service

职责：

- 提供 `Project` 的 CRUD 与应用归属入口
- 提供 `Application` 的 CRUD
- 维护 `active_manifest` 绑定
- 不负责 release / verify 运行态资源

当前实现：

- `cmd/main.go` 通过 `devflow-service-common/bootstrap` 启动
- `pkg/api/project.go`
- `pkg/api/application.go`
- `pkg/service/project.go`
- `pkg/service/application.go`
- `pkg/router/project.go`
- `pkg/router/application.go`

建议端口：

- `APP_SERVICE_PORT`
- `APP_SERVICE_METRICS_PORT`
- `APP_SERVICE_PPROF_PORT`

运行时：

- 上报的 OTel `service.name` 为 `app-service`
- 任何 outbound service / external call 都必须带 `metrics + trace + structured log`
- 默认 harness 为 `Planner -> Generator -> Evaluator`，并且支持 delegation 时必须真实启动 sub-agents

接口：

- `GET /api/v1/projects`
- `POST /api/v1/projects`
- `GET /api/v1/projects/:id`
- `PUT /api/v1/projects/:id`
- `DELETE /api/v1/projects/:id`
- `GET /api/v1/projects/:id/applications`
- `GET /api/v1/applications`
- `POST /api/v1/applications`
- `GET /api/v1/applications/:id`
- `PUT /api/v1/applications/:id`
- `DELETE /api/v1/applications/:id`
- `PATCH /api/v1/applications/:id/active_manifest`
