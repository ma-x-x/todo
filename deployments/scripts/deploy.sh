#!/bin/bash

# 确保脚本在出错时退出
set -e

# 检查并设置系统参数
setup_system_params() {
    echo "正在设置系统参数..."
    
    # 设置 vm.overcommit_memory
    if [ "$(sysctl -n vm.overcommit_memory)" != "1" ]; then
        echo "设置 vm.overcommit_memory = 1"
        sudo sysctl -w vm.overcommit_memory=1
        echo "vm.overcommit_memory = 1" | sudo tee -a /etc/sysctl.conf
    fi
    
    # 设置 somaxconn
    if [ "$(sysctl -n net.core.somaxconn)" -lt "512" ]; then
        echo "设置 net.core.somaxconn = 512"
        sudo sysctl -w net.core.somaxconn=512
        echo "net.core.somaxconn = 512" | sudo tee -a /etc/sysctl.conf
    fi
    
    # 禁用 THP (Transparent Huge Pages)
    if [ -f /sys/kernel/mm/transparent_hugepage/enabled ]; then
        echo never | sudo tee /sys/kernel/mm/transparent_hugepage/enabled
        echo never | sudo tee /sys/kernel/mm/transparent_hugepage/defrag
    fi
}

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

# 设置系统参数
setup_system_params

# 设置环境变量
export MYSQL_ROOT_PASSWORD=${MYSQL_ROOT_PASSWORD}
export DB_PASSWORD=${DB_PASSWORD}
export DB_HOST=${DB_HOST:-mysql}
export DB_PORT=3306
export DB_USER=${DB_USER:-todo_user}
export DB_NAME=todo_db
export REDIS_HOST=redis
export REDIS_PORT=6379
export REDIS_PASSWORD=${REDIS_PASSWORD}
export JWT_SECRET=${JWT_SECRET}
export APP_ENV=prod
export LOG_LEVEL=debug
export CONFIG_FILE=/app/configs/config.prod.yaml

# 检查服务是否已经存在并运行
check_service_exists() {
    local service_name=$1
    if docker ps --format '{{.Names}}' | grep -q "^${service_name}$"; then
        return 0  # 服务正在运行
    fi
    return 1  # 服务未运行
}

# 检查服务是否健康
check_mysql_health() {
    if docker exec todo-mysql mysqladmin ping -h mysql -u"root" -p"${MYSQL_ROOT_PASSWORD}" --silent > /dev/null 2>&1; then
        return 0  # MySQL 正常
    fi
    return 1  # MySQL 异常
}

check_redis_health() {
    if [ -z "${REDIS_PASSWORD}" ]; then
        if docker exec todo-redis redis-cli -h redis ping > /dev/null 2>&1; then
            return 0  # Redis 正常（无密码）
        fi
    else
        if docker exec todo-redis redis-cli -h redis -a "${REDIS_PASSWORD}" ping > /dev/null 2>&1; then
            return 0  # Redis 正常（有密码）
        fi
    fi
    return 1  # Redis 异常
}

# 等待 MySQL 就绪的函数
wait_for_mysql() {
    echo "等待 MySQL 就绪..."
    for i in {1..60}; do
        if check_mysql_health; then
            echo "MySQL 已就绪！"
            return 0
        fi
        echo "等待 MySQL 就绪中... ($i/60)"
        sleep 2
    done
    echo "MySQL 未能在指定时间内就绪"
    return 1
}

# 等待 Redis 就绪的函数
wait_for_redis() {
    echo "等待 Redis 就绪..."
    for i in {1..60}; do
        if check_redis_health; then
            echo "Redis 已就绪！"
            return 0
        fi
        echo "等待 Redis 就绪中... ($i/60)"
        sleep 2
    done
    echo "Redis 未能在指定时间内就绪"
    return 1
}

# 创建必要的目录
mkdir -p logs

# 创建 docker 网络（如果不存在）
docker network create todo-network 2>/dev/null || true

# 检查并启动 MySQL
if check_service_exists "todo-mysql" && check_mysql_health; then
    echo "MySQL 已在运行且状态正常，跳过部署"
else
    echo "正在启动 MySQL..."
    docker-compose up -d mysql
    wait_for_mysql || exit 1
    
    # 初始化数据库
    echo "正在初始化数据库..."
    docker-compose exec -T mysql mysql -uroot -p"${MYSQL_ROOT_PASSWORD}" -e "
    CREATE DATABASE IF NOT EXISTS todo_db;
    CREATE USER IF NOT EXISTS 'todo_user'@'%' IDENTIFIED BY '${DB_PASSWORD}';
    GRANT ALL PRIVILEGES ON todo_db.* TO 'todo_user'@'%';
    FLUSH PRIVILEGES;
    "

    # 导入初始化 SQL
    echo "正在导入数据库架构..."
    docker-compose exec -T mysql mysql -uroot -p"${MYSQL_ROOT_PASSWORD}" todo_db < scripts/init.sql
fi

# 检查并启动 Redis
if check_service_exists "todo-redis" && check_redis_health; then
    echo "Redis 已在运行且状态正常，跳过部署"
else
    echo "正在启动 Redis..."
    docker-compose up -d redis
    wait_for_redis || exit 1
fi

# 停止并重新启动应用
echo "正在重新部署应用..."
docker-compose rm -sf app || true
docker-compose up -d app

# 等待应用就绪
echo "等待应用就绪..."
for i in {1..30}; do
    if curl -s http://localhost:8081/health > /dev/null; then
        echo "应用已就绪！"
        break
    fi
    echo "等待应用就绪中... ($i/30)"
    sleep 2
done

# 检查服务状态
echo "检查服务状态..."
docker-compose ps

# 检查应用日志
echo "检查应用日志..."
docker-compose logs --tail=50 app

echo "部署成功完成！"

# 检查系统参数
check_system_params() {
    echo "检查系统参数..."
    
    # 检查 vm.overcommit_memory
    if [ "$(sysctl -n vm.overcommit_memory)" != "1" ]; then
        echo "警告：vm.overcommit_memory 未设置为 1"
    fi
    
    # 检查 somaxconn
    if [ "$(sysctl -n net.core.somaxconn)" -lt "512" ]; then
        echo "警告：net.core.somaxconn 小于 512"
    fi
    
    # 检查 THP
    if [ -f /sys/kernel/mm/transparent_hugepage/enabled ]; then
        if ! grep -q "\[never\]" /sys/kernel/mm/transparent_hugepage/enabled; then
            echo "警告：透明大页面(THP)未设置为 never"
        fi
    fi
}

# 在部署前检查系统参数
check_system_params 