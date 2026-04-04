# 约束

## 资源归属

- `Project` 是顶层空间资源
- `Application` 归属于 `Project`
- `Application.active_manifest` 只能引用应用当前活动版本

## 边界约束

- 不允许在 app-service 中引入 `Manifest`、`Release`、`Intent`、`Configuration`、`Verify` 对外资源
- 不允许把执行面状态写回职责放进 app-service
- 不允许把其他服务的领域模型复制回 app-service

## 状态与删除

- 删除应遵守仓库既有软删除语义
- 对外返回的列表与详情必须遵守软删除过滤规则
- `active_manifest` 的更新必须保持幂等

## 跨服务调用

- 只要出现出站调用，就必须同时产出 `metrics + trace + structured log`
- 出站调用的资源主键只能进日志和 trace attribute，不能进 metrics label

## 禁止事项

- 不要为了便利把执行面对象并入 `Application`
- 不要在 handler 中直接拼接数据库更新逻辑
- 不要把 Swagger 里的历史接口当作当前边界
