# k3s创建KubeConfig

本文描述了创建一个以token方式连接到K3S集群的方法。

## 创建KubeConfig

本步骤直接以default:default账号来进行创建，如若需要创建其他serviceAccount的KubeConfig，请自行创建。

### 获取default ServiceAccount

```shell
kubectl describe serviceAccounts default
```

output:

```shell
Name:                default
Namespace:           default
Labels:              <none>
Annotations:         <none>
Image pull secrets:  <none>
Mountable secrets:   default-token-hrr7q
Tokens:              default-token-hrr7q
Events:              <none>
```

### 获取token内容

```shell
kubectl describe secrets default-token-hrr7q
```

output:

```shell
Name:         default-token-hrr7q
Namespace:    default
Labels:       <none>
Annotations:  kubernetes.io/service-account.name: default
              kubernetes.io/service-account.uid: c345707c-450e-405e-b862-cb94de95f0f0

Type:  kubernetes.io/service-account-token

Data
====
token:      eyJhbGciOiJSUzI1NiIsImtpZCI6InZyRGVVT2l6S0pxNlM5SDAzLVRmZm1NVC1iZlVHQzdWbm9hbXdFT0E4NGsifQ.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJkZWZhdWx0Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9zZWNyZXQubmFtZSI6ImRlZmF1bHQtdG9rZW4taHJyN3EiLCJrdWJlcm5ldGVzLmlvL3NlcnZpY2VhY2NvdW50L3NlcnZpY2UtYWNjb3VudC5uYW1lIjoiZGVmYXVsdCIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50LnVpZCI6ImMzNDU3MDdjLTQ1MGUtNDA1ZS1iODYyLWNiOTRkZTk1ZjBmMCIsInN1YiI6InN5c3RlbTpzZXJ2aWNlYWNjb3VudDpkZWZhdWx0OmRlZmF1bHQifQ.jjmgNN3CleDtaMirAYteFWUcMFcwq38g6WqYgMD8xnODh-slQRWhFhI0t2o4DG73if9Tpqjwe5nqnn4wZ_-zXlU1fJjE6cHgyVlo1t4rFmp4lwurmBXv3nVhVhnaU2Mfw8ZImI2lJNHNq8Mn7jol8dkpqyewli4sD8E-Mv_f_8kmJToufpsJRc-C67dI6Z0hdBEilznHotIjY1wIYY-xNUW4f1ZCh8dRHXhXYfnsTbH5qY-i_H9nOsqM2gnbzl0euXIsJg6HHS7heM7kI4CpPPaZv8o9T0hS1P6DxTEbW98v_zzazASBaRtTkumbvEnDVYojpRr6f7mdJ1F1kf7xC1
ca.crt:     526 bytes
namespace:  7 bytes
```

保存token节的内容。

### 获取cluster证书

```shell
kubectl config view --flatten --minify
```

output:

```shell
apiVersion: v1
clusters:
- cluster:
    certificate-authority-data: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUJWekNCL3FBREFnRUNBZ0VBTUFvR0NDcUdTTTQ5QkFNQ01DTXhJVEFmQmdOVkJBTU1HR3N6Y3kxelpYSjIKWlhJdFkyRkFNVFU1TmpZNE1EUTJOakFlRncweU1EQTRNRFl3TWpJeE1EWmFGdzB6TURBNE1EUXdNakl4TURaYQpNQ014SVRBZkJnTlZCQU1NR0dzemN5MXpaWEoyWlhJdFkyRkFNVFU1TmpZNE1EUTJOakJaTUJNR0J5cUdTTTQ5CkFnRUdDQ3FHU000OUF3RUhBMElBQkMxK2dkQ1I0bWNjR016Z3h6Z2cyTkxla1grMUhSUVdoN1JhMXM3bDZnWUgKZDIzbTV3aFBSTTZseWpOWnhiN3l6R29IUFQreUNrWlNCZlVuZFZEZFdQYWpJekFoTUE0R0ExVWREd0VCL3dRRQpBd0lDcERBUEJnTlZIUk1CQWY4RUJUQURBUUgvTUFvR0NDcUdTTTQ5QkFNQ0EwZ0FNRVVDSUNLd1VWK2hYZlB0Cm4wWE1LOEV0MFcvNHRDWlZqK3QvZ284bENiNzFSelIzQWlFQWhtNk9TVE5rYkNISDFEU2hjVytFdnFtL1J3YkYKVGc4eUVwUWFHUHo4WmVjPQotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0tCa==
    server: https://127.0.0.1:6443
  name: default
contexts:
- context:
    cluster: default
    user: default
  name: default
current-context: default
kind: Config
preferences: {}
users:
- name: default
  user:
    password: 7fbb30d3a1e5d639d75e6134e9eea519
    username: admin
```

保存certificate-authority-data的内容。

### 创建新KubeConfig

按照以下yaml模板填写好以下三项：

1. cluster节的证书，集群名称，k3s apiserver运行地址及端口
2. contexts节的集群名称，用户名，上下文名称（注意要使用username@clusterName格式）
3. users节的用户名，token

```yaml
apiVersion: v1
kind: Config
clusters:
  - cluster:
  	  # 来自获取cluster证书步骤
      certificate-authority-data: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUJWekNCL3FBREFnRUNBZ0VBTUFvR0NDcUdTTTQ5QkFNQ01DTXhJVEFmQmdOVkJBTU1HR3N6Y3kxelpYSjIKWlhJdFkyRkFNVFl3TWpZMk56QTROVEFlRncweU1ERXdNVFF3T1RFNE1EVmFGdzB6TURFd01USXdPVEU0TURWYQpNQ014SVRBZkJnTlZCQU1NR0dzemN5MXpaWEoyWlhJdFkyRkFNVFl3TWpZMk56QTROVEJaTUJNR0J5cUdTTTQ5CkFnRUdDQ3FHU000OUF3RUhBMElBQk56eUMzNktQWVFyZnhxVFdWU0daRS9UalJqUzVNckx1TGRjdVpRazV1ZnAKQmNDVTJhNm9OT1JQNDZLcTBQdXFrdjlCaVdlVEJULzF2KzJ3bDJ0NCsxT2pJekFoTUE0R0ExVWREd0VCL3dRRQpBd0lDcERBUEJnTlZIUk1CQWY4RUJUQURBUUgvTUFvR0NDcUdTTTQ5QkFNQ0EwZ0FNRVVDSVFESkUwNys3dlJVCjJkcWtqKzdHdVkweTFOUHJFUDFWZDFvSStzZTg0bnh0RGdJZ0JsRlpwZm96UjNYem5hWVRCemVMekUxa2RQUTQKWmhjNnRFWG5GaVRjMllJPQotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0tCa==
      # k3s apiserver 运行地址与端口
      server: https://10.10.10.35:6443
    name: default-cluster
contexts:
  - context:
      cluster: default-cluster
      user: default
    name: default@default-cluster
current-context: default@default-cluster
users:
  - name: default
    user:
      # 来自获取token内容步骤
      token: eyJhbGciOiJSUzI1NiIsImtpZCI6InZyRGVVT2l6S0pxNlM5SDAzLVRmZm1NVC1iZlVHQzdWbm9hbXdFT0E4NGsifQ.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJkZWZhdWx0Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9zZWNyZXQubmFtZSI6ImRlZmF1bHQtdG9rZW4taHJyN3EiLCJrdWJlcm5ldGVzLmlvL3NlcnZpY2VhY2NvdW50L3NlcnZpY2UtYWNjb3VudC5uYW1lIjoiZGVmYXVsdCIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50LnVpZCI6ImMzNDU3MDdjLTQ1MGUtNDA1ZS1iODYyLWNiOTRkZTk1ZjBmMCIsInN1YiI6InN5c3RlbTpzZXJ2aWNlYWNjb3VudDpkZWZhdWx0OmRlZmF1bHQifQ.jjmgNN3CleDtaMirAYteFWUcMFcwq38g6WqYgMD8xnODh-slQRWhFhI0t2o4DG73if9Tpqjwe5nqnn4wZ_-zXlU1fJjE6cHgyVlo1t4rFmp4lwurmBXv3nVhVhnaU2Mfw8ZImI2lJNHNq8Mn7jol8dkpqyewli4sD8E-Mv_f_8kmJToufpsJRc-C67dI6Z0hdBEilznHotIjY1wIYY-xNUW4f1ZCh8dRHXhXYfnsTbH5qY-i_H9nOsqM2gnbzl0euXIsJg6HHS7heM7kI4CpPPaZv8o9T0hS1P6DxTEbW98v_zzazASBaRtTkumbvEnDVYojpRr6f7mdJ1F1kf7xC1
```

## 绑定角色

为system:serviceaccount:default:default在集群中绑定角色

```yaml
# NOTE: The service account `default:default` already exists in k8s cluster.
# You can create a new account following like this:
#---
#apiVersion: v1
#kind: ServiceAccount
#metadata:
#  name: <new-account-name>
#  namespace: <namespace>

---
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: ClusterRoleBinding
metadata:
  name: my-cluster-role-binding
subjects:
  - kind: ServiceAccount
    # Reference to upper's `metadata.name`
    name: default
    # Reference to upper's `metadata.namespace`
    namespace: default
roleRef:
  kind: ClusterRole
  name: cluster-admin	# 指定使用哪个角色
  apiGroup: rbac.authorization.k8s.io
```

完成后可通过新的KubeConfig接入集群。