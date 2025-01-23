# Kubernetes 部署指南

本文档详细介绍如何在 Kubernetes 集群中部署待办事项管理系统。

## 目录结构

```bash
deployments/kubernetes/
├── README.md           # 部署说明文档
├── configmap.yaml      # 应用配置
├── deployment.yaml     # 应用部署配置
├── ingress.yaml        # 入口配置
├── mysql.yaml          # MySQL数据库配置
├── redis.yaml          # Redis缓存配置
├── redis-config.yaml   # Redis详细配置
├── secret.yaml         # 敏感信息配置
├── service.yaml        # 服务配置
├── prometheus.yaml     # 监控配置
├── grafana.yaml        # 监控面板配置
└── deploy.sh          # 部署脚本
```

## 前置要求

在开始部署之前，请确保您的环境满足以下条件：

1. 已安装并配置 Kubernetes 集群（版本 >= 1.18）
2. 已安装 kubectl 命令行工具
3. 已安装 Ingress Controller（推荐 nginx-ingress）
4. 已配置默认 StorageClass（用于持久化存储）

可以使用以下命令检查环境：

```bash
# 检查 kubectl 是否正确配置
kubectl cluster-info

# 检查 Kubernetes 版本
kubectl version --short

# 检查 StorageClass
kubectl get storageclass

# 检查 Ingress Controller
kubectl get pods -n ingress-nginx
```

## 快速部署

### 1. 配置准备

1. 克隆项目代码：
```bash
git clone <项目地址>
cd deployments/kubernetes
```

2. 修改配置文件：
   - 检查并修改 `configmap.yaml` 中的应用配置
   - 更新 `secret.yaml` 中的敏感信息（注意使用 base64 编码）
   - 根据需要调整 `ingress.yaml` 中的域名配置

### 2. 执行部署

可以使用部署脚本一键部署：

```bash
chmod +x deploy.sh
./deploy.sh
```

或者手动执行部署步骤：

```bash
# 创建命名空间
kubectl create namespace todo-app

# 部署配置和密钥
kubectl apply -f configmap.yaml
kubectl apply -f secret.yaml

# 部署数据库和缓存
kubectl apply -f mysql.yaml
kubectl apply -f redis-config.yaml
kubectl apply -f redis.yaml

# 部署应用
kubectl apply -f deployment.yaml
kubectl apply -f service.yaml
kubectl apply -f ingress.yaml

# 部署监控组件（可选）
kubectl apply -f prometheus.yaml
kubectl apply -f grafana.yaml
```

### 3. 验证部署

1. 检查 Pod 状态：
```bash
kubectl get pods -n todo-app
```

2. 检查服务状态：
```bash
kubectl get svc -n todo-app
```

3. 检查 Ingress 配置：
```bash
kubectl get ingress -n todo-app
```

## 监控配置

系统集成了 Prometheus 和 Grafana 用于监控：

1. 访问 Grafana：
```bash
kubectl port-forward svc/grafana 3000:3000 -n todo-app
```
然后访问 http://localhost:3000

2. 访问 Prometheus：
```bash
kubectl port-forward svc/prometheus 9090:9090 -n todo-app
```
然后访问 http://localhost:9090

## 常见问题

### 1. Pod 无法启动
- 检查 Pod 日志：
```bash
kubectl logs <pod-name> -n todo-app
```
- 检查 Pod 描述：
```bash
kubectl describe pod <pod-name> -n todo-app
```

### 2. 数据库连接失败
- 确认 MySQL Pod 是否正常运行
- 检查密码配置是否正确
- 验证网络连接：
```bash
kubectl exec -it <pod-name> -n todo-app -- ping mysql
```

### 3. 存储问题
- 检查 PVC 状态：
```bash
kubectl get pvc -n todo-app
```
- 确认 StorageClass 配置正确

## 生产环境注意事项

1. 安全配置
   - 更新所有默认密码
   - 配置 TLS 证书
   - 限制网络访问策略
   - 配置资源配额

2. 高可用配置
   - 增加应用副本数
   - 配置 Pod 反亲和性
   - 使用多可用区部署
   - 配置自动扩缩容

3. 备份策略
   - 配置数据库定期备份
   - 使用持久化存储
   - 配置日志采集

## 扩展配置

### 配置自动扩缩容

```bash
kubectl autoscale deployment todo-api \
  --cpu-percent=80 \
  --min=3 \
  --max=10 \
  -n todo-app
```

### 配置资源配额

```bash
kubectl apply -f - <<EOF
apiVersion: v1
kind: ResourceQuota
metadata:
  name: todo-quota
  namespace: todo-app
spec:
  hard:
    requests.cpu: "4"
    requests.memory: 8Gi
    limits.cpu: "8"
    limits.memory: 16Gi
EOF
```

## 清理部署

如需清理所有资源：

```bash
kubectl delete namespace todo-app
```

## 参考文档

- [Kubernetes 官方文档](https://kubernetes.io/docs/)
- [Prometheus 文档](https://prometheus.io/docs/)
- [Grafana 文档](https://grafana.com/docs/)
