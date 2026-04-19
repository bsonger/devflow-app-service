# devflow-app-service deploy

Apply together with a staging `app-service-config` secret and the shared `aliyun-docker-config` pull secret.

## RBAC prerequisites

The app-service needs permission to write Argo CD cluster-registration secrets in the `argocd` namespace.

Apply in this order:

```bash
kubectl apply -f deploy/serviceaccount.yaml
kubectl apply -f deploy/role-argocd.yaml
kubectl apply -f deploy/rolebinding-argocd.yaml
kubectl apply -f deploy/deployment.yaml
kubectl apply -f deploy/service.yaml
```

Required resources:
- `ServiceAccount` `devflow-app-service` in `devflow-staging`
- `Role` `devflow-app-service-argocd` in `argocd` (secrets: create, get, list, watch, patch, update, delete)
- `RoleBinding` `devflow-app-service-argocd` in `argocd`
