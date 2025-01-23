#!/bin/bash

# 创建命名空间
kubectl create namespace todo-app

# 创建配置和密钥
kubectl apply -f configmap.yaml # 应用配置
kubectl apply -f secret.yaml    # 敏感信息

# 部署数据库
kubectl apply -f mysql.yaml # MySQL数据库
kubectl apply -f redis.yaml # Redis缓存

# 等待数据库就绪
echo "Waiting for MySQL to be ready..."
kubectl wait --for=condition=ready pod -l app=mysql -n todo-app --timeout=300s
echo "Waiting for Redis to be ready..."
kubectl wait --for=condition=ready pod -l app=redis -n todo-app --timeout=300s

# 部署应用
kubectl apply -f deployment.yaml # API应用
kubectl apply -f service.yaml    # 服务
kubectl apply -f ingress.yaml    # 入口配置

# 等待应用就绪
echo "Waiting for Todo API to be ready..."
kubectl wait --for=condition=ready pod -l app=todo-api -n todo-app --timeout=300s

echo "Deployment completed successfully!"
