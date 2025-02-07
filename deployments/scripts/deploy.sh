#!/bin/bash

# 确保脚本在出错时退出
set -e

# 检查必要的环境变量
check_required_env() {
    local missing_vars=()
    
    if [ -z "${MYSQL_ROOT_PASSWORD}" ]; then
        missing_vars+=("MYSQL_ROOT_PASSWORD")
    fi
    if [ -z "${DB_PASSWORD}" ]; then
        missing_vars+=("DB_PASSWORD")
    fi
    if [ -z "${JWT_SECRET}" ]; then
        missing_vars+=("JWT_SECRET")
    fi
    
    if [ ${#missing_vars[@]} -ne 0 ]; then
        echo "Error: Missing required environment variables:"
        printf '%s\n' "${missing_vars[@]}"
        exit 1
    fi
}

# 检查环境变量
check_required_env

# 设置环境变量
export MYSQL_ROOT_PASSWORD=${MYSQL_ROOT_PASSWORD}
export DB_PASSWORD=${DB_PASSWORD}
export REDIS_PASSWORD=${REDIS_PASSWORD}
export DB_HOST=mysql
export REDIS_HOST=redis

# 检查服务是否已经存在并运行
check_service_exists() {
    local service_name=$1
    if docker ps --format '{{.Names}}' | grep -q "^${service_name}$"; then
        return 0  # 服务存在且运行中
    fi
    return 1  # 服务不存在或未运行
}

# 检查服务是否健康
check_mysql_health() {
    if docker exec todo-mysql mysqladmin ping -h mysql -u"root" -p"${MYSQL_ROOT_PASSWORD}" --silent > /dev/null 2>&1; then
        return 0  # MySQL 健康
    fi
    return 1  # MySQL 不健康
}

check_redis_health() {
    if [ -z "${REDIS_PASSWORD}" ]; then
        if docker exec todo-redis redis-cli -h redis ping > /dev/null 2>&1; then
            return 0  # Redis 健康（无密码）
        fi
    else
        if docker exec todo-redis redis-cli -h redis -a "${REDIS_PASSWORD}" ping > /dev/null 2>&1; then
            return 0  # Redis 健康（有密码）
        fi
    fi
    return 1  # Redis 不健康
}

# 等待 MySQL 就绪的函数
wait_for_mysql() {
    echo "Waiting for MySQL to be ready..."
    for i in {1..60}; do
        if check_mysql_health; then
            echo "MySQL is ready!"
            return 0
        fi
        echo "Waiting for MySQL to be ready... ($i/60)"
        sleep 2
    done
    echo "MySQL did not become ready in time"
    return 1
}

# 等待 Redis 就绪的函数
wait_for_redis() {
    echo "Waiting for Redis to be ready..."
    for i in {1..60}; do
        if check_redis_health; then
            echo "Redis is ready!"
            return 0
        fi
        echo "Waiting for Redis to be ready... ($i/60)"
        sleep 2
    done
    echo "Redis did not become ready in time"
    return 1
}

# 创建必要的目录
mkdir -p logs

# 创建 docker 网络（如果不存在）
docker network create todo-network 2>/dev/null || true

# 检查并启动 MySQL
if check_service_exists "todo-mysql" && check_mysql_health; then
    echo "MySQL is already running and healthy, skipping deployment"
else
    echo "Starting MySQL..."
    docker-compose up -d mysql
    wait_for_mysql || exit 1
    
    # 初始化数据库
    echo "Initializing database..."
    docker-compose exec -T mysql mysql -uroot -p"${MYSQL_ROOT_PASSWORD}" -e "
    CREATE DATABASE IF NOT EXISTS todo_db;
    CREATE USER IF NOT EXISTS 'todo_user'@'%' IDENTIFIED BY '${DB_PASSWORD}';
    GRANT ALL PRIVILEGES ON todo_db.* TO 'todo_user'@'%';
    FLUSH PRIVILEGES;
    "

    # 导入初始化 SQL
    echo "Importing database schema..."
    docker-compose exec -T mysql mysql -uroot -p"${MYSQL_ROOT_PASSWORD}" todo_db < scripts/init.sql
fi

# 检查并启动 Redis
if check_service_exists "todo-redis" && check_redis_health; then
    echo "Redis is already running and healthy, skipping deployment"
else
    echo "Starting Redis..."
    docker-compose up -d redis
    wait_for_redis || exit 1
fi

# 停止并重新启动应用
echo "Redeploying application..."
docker-compose rm -sf app || true
docker-compose up -d app

# 等待应用就绪
echo "Waiting for application to be ready..."
for i in {1..30}; do
    if curl -s http://localhost:8081/health > /dev/null; then
        echo "Application is ready!"
        break
    fi
    echo "Waiting for application to be ready... ($i/30)"
    sleep 2
done

# 检查服务状态
echo "Checking service status..."
docker-compose ps

# 检查应用日志
echo "Checking application logs..."
docker-compose logs --tail=50 app

echo "Deployment completed successfully!" 