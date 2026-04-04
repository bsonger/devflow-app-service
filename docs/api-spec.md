# 接口规范

`devflow-app-service` 的对外 HTTP 接口只围绕 `Project` 和 `Application`。

## Project

- `POST /api/v1/projects`
- `GET /api/v1/projects`
- `GET /api/v1/projects/{id}`
- `PUT /api/v1/projects/{id}`
- `DELETE /api/v1/projects/{id}`
- `GET /api/v1/projects/{id}/applications`

## Application

- `POST /api/v1/applications`
- `GET /api/v1/applications`
- `GET /api/v1/applications/{id}`
- `PUT /api/v1/applications/{id}`
- `DELETE /api/v1/applications/{id}`
- `PATCH /api/v1/applications/{id}/active_manifest`

## 语义约定

- 所有列表接口默认支持分页参数，返回值与 Swagger 保持一致
- `active_manifest` 只表示当前应用绑定的活动版本，不承载构建或发布编排语义
- 404 表示资源不存在，400 表示请求参数非法，409 表示状态或边界冲突
- 删除语义以逻辑删除为主，查询接口需要遵守仓库中的软删除约定

## 不提供

- `Manifest`
- `Release`
- `Intent`
- `Configuration`
- `Verify`
