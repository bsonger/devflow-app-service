# App Service

职责：

- 提供 `Project` 的 CRUD 与应用归属入口
- 提供 `Application` 的 CRUD
- 维护 `active_manifest` 绑定
- 作为应用元数据的查询入口

当前复用的现有实现：

- `pkg/api/project.go`
- `pkg/service/project.go`
- `pkg/router/project.go`
- `pkg/api/application.go`
- `pkg/service/application.go`
- `pkg/router/application.go`

建议端口：

- `APP_SERVICE_PORT`
- `APP_SERVICE_METRICS_PORT`
- `APP_SERVICE_PPROF_PORT`

运行时：

- 上报的 OTel `service.name` 为 `app-service`

接口：

- `GET /api/v1/projects`
- `POST /api/v1/projects`
- `GET /api/v1/projects/:id`
- `PUT /api/v1/projects/:id`
- `DELETE /api/v1/projects/:id`
- `GET /api/v1/projects/:id/applications`
- `GET /api/v1/applications`
- `POST /api/v1/applications`
