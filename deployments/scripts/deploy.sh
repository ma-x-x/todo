#!/bin/bash

# 确保脚本在出错时退出
set -e

# 系统参数检查函数
check_system_params() {
    echo "检查系统参数..."
    local warnings=()
    
    # 检查系统参数并收集警告信息
    if [ "$(sysctl -n vm.overcommit_memory)" != "1" ]; then
        warnings+=("vm.overcommit_memory 未设置为推荐值 1，可能影响 Redis 性能")
    fi
    
    if [ "$(sysctl -n net.core.somaxconn)" -lt "1024" ]; then
        warnings+=("net.core.somaxconn 小于推荐值 1024，可能影响高并发处理")
    fi
    
    if [ -f /sys/kernel/mm/transparent_hugepage/enabled ]; then
        if ! grep -q "\[never\]" /sys/kernel/mm/transparent_hugepage/enabled; then
            warnings+=("透明大页面(THP)未禁用，可能导致 Redis 性能问题")
        fi
    fi

    # 如果有警告，统一显示
    if [ ${#warnings[@]} -gt 0 ]; then
        echo "⚠️ 性能优化建议："
        printf '%s\n' "${warnings[@]}"
    fi
}

# 检查环境变量
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

# 部署前的系统检查
check_system_params

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
export SWAGGER_HOST=${SWAGGER_HOST:-api.example.com}
export APP_ENV=prod
export LOG_LEVEL=info
export CONFIG_FILE=/app/configs/config.prod.yaml
export TZ=Asia/Shanghai  # 设置时区为中国时区

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
    local max_attempts=30
    local attempt=1
    local wait_time=2

    if [ -z "${REDIS_PASSWORD}" ]; then
        while [ $attempt -le $max_attempts ]; do
            if docker exec todo-redis redis-cli -h redis ping > /dev/null 2>&1; then
                echo "Redis 连接成功（无密码）"
                return 0
            fi
            echo "尝试连接 Redis 中... ($attempt/$max_attempts)"
            sleep $wait_time
            attempt=$((attempt + 1))
        done
    else
        while [ $attempt -le $max_attempts ]; do
            if docker exec todo-redis redis-cli -h redis -a "${REDIS_PASSWORD}" ping > /dev/null 2>&1; then
                echo "Redis 连接成功（带密码）"
                return 0
            fi
            echo "尝试连接 Redis 中... ($attempt/$max_attempts)"
            sleep $wait_time
            attempt=$((attempt + 1))
        done
    fi

    echo "Redis 健康检查失败，查看日志："
    docker logs todo-redis
    return 1
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

    # 等待数据库完全就绪
    sleep 5

    # 检查数据库是否需要初始化
    TABLES_EXIST=$(docker-compose exec -T mysql mysql -uroot -p"${MYSQL_ROOT_PASSWORD}" -N -e "
        SELECT COUNT(*) FROM information_schema.tables 
        WHERE table_schema = 'todo_db' 
        AND table_name IN ('users', 'todos', 'categories', 'reminders');
    " todo_db)

    if [ "$TABLES_EXIST" = "0" ]; then
        echo "数据库为空，开始初始化..."
        docker-compose exec -T mysql mysql -uroot -p"${MYSQL_ROOT_PASSWORD}" todo_db < scripts/init.sql
        echo "数据库初始化完成"
    else
        echo "数据库表已存在，跳过初始化"
    fi
fi

# Redis 启动前临时设置系统参数
echo "临时设置系统参数..."
if [ -x "$(command -v sudo)" ]; then
    sudo sysctl -w net.core.somaxconn=1024 || echo "警告: 无法设置 somaxconn"
fi

# 检查并启动 Redis
if check_service_exists "todo-redis" && check_redis_health; then
    echo "Redis 已在运行且状态正常，跳过部署"
else
    echo "正在启动 Redis..."
    echo "Redis 配置信息："
    echo "- 密码已设置: $([ -n "${REDIS_PASSWORD}" ] && echo "是" || echo "否")"
    echo "- 数据持久化: 已启用"
    echo "- 最大内存: 2GB"
    
    docker-compose up -d redis
    
    echo "等待 Redis 启动..."
    sleep 5
    
    echo "Redis 容器日志："
    docker-compose logs redis
    
    wait_for_redis || exit 1
fi

# 停止并重新启动应用
echo "正在重新部署应用..."
# 验证配置文件
echo "验证配置文件..."
if [ -f "configs/config.prod.yaml" ]; then
    # 使用 envsubst 处理配置文件中的环境变量
    envsubst < configs/config.prod.yaml > configs/config.prod.yaml.tmp
    mv configs/config.prod.yaml.tmp configs/config.prod.yaml
    echo "配置文件内容:"
    cat configs/config.prod.yaml
else
    echo "错误: 找不到配置文件 configs/config.prod.yaml"
    exit 1
fi

docker-compose build --no-cache app
docker-compose up -d --force-recreate app

# 等待并检查应用日志
echo "检查应用启动日志..."
sleep 5
docker-compose logs app

# 等待应用就绪
echo "等待应用就绪..."
for i in {1..180}; do  # 最多等待3分钟
    # 检查健康检查接口，确保返回包含 "healthy" 的响应
    if curl -s -f http://localhost:8081/health | grep -q "healthy"; then
        echo "应用已就绪！"
        break
    fi
    # 如果达到最大重试次数，则输出日志并退出
    if [ $i -eq 180 ]; then
        echo "应用未能在指定时间内就绪，检查日志..."
        docker-compose logs --tail=100 app
        exit 1
    fi
    echo "等待应用就绪中... ($i/180)"
    sleep 2  # 每次检查间隔2秒
done

# 检查容器状态
echo "检查容器状态..."
docker-compose ps

# 如果应用不健康，输出详细日志
if [ "$(docker inspect --format='{{.State.Health.Status}}' todo-api)" != "healthy" ]; then
    echo "应用健康检查失败，输出详细日志："
    docker-compose logs --tail=200 app
    exit 1
fi

# 检查服务状态
echo "检查服务状态..."
docker-compose ps

# 检查应用日志
echo "检查应用日志..."
docker-compose logs --tail=50 app

echo "部署成功完成！"

# 在部署开始时添加
if [ -x "$(command -v sudo)" ] && [ -f "./setup_system.sh" ]; then
    echo "检查系统参数..."
    if ! sudo ./setup_system.sh; then
        echo "⚠️ 系统参数可能未达到最优状态，可能会影响性能"
        echo "建议在部署后运行 sudo ./setup_system.sh 进行优化"
    fi
fi 