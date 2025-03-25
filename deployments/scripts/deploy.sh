#!/bin/bash

# 确保脚本在出错时退出
set -e

# 添加错误处理函数
handle_error() {
    local exit_code=$?
    local line_number=$1
    echo "错误发生在第 ${line_number} 行，退出码: ${exit_code}"
    exit $exit_code
}

# 在脚本开头添加
trap 'handle_error ${LINENO}' ERR

# 添加日志函数
log_info() {
    echo "[INFO] $(date '+%Y-%m-%d %H:%M:%S') - $1"
}

log_error() {
    echo "[ERROR] $(date '+%Y-%m-%d %H:%M:%S') - $1" >&2
}

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

# 设置环境变量（更新默认值）
export MYSQL_ROOT_PASSWORD=${MYSQL_ROOT_PASSWORD}
export DB_PASSWORD=${DB_PASSWORD}
export DB_HOST=${DB_HOST:-localhost}  # 默认使用本地主机
export DB_PORT=${DB_PORT:-3306}
export DB_USER=${DB_USER:-todo_user}
export DB_NAME=${DB_NAME:-todo_db}
export REDIS_HOST=${REDIS_HOST:-localhost}  # 默认使用本地主机
export REDIS_PORT=${REDIS_PORT:-6379}
export REDIS_PASSWORD=${REDIS_PASSWORD}
export JWT_SECRET=${JWT_SECRET}
export SWAGGER_HOST=${SWAGGER_HOST:-api.example.com}
export APP_ENV=prod
export LOG_LEVEL=info
export CONFIG_FILE=/app/configs/config.prod.yaml
export TZ=Asia/Shanghai  # 设置时区为中国时区

# 检查 MySQL 是否健康
check_mysql_health() {
    local max_retries=3
    local retry_count=0
    
    while [ $retry_count -lt $max_retries ]; do
        if mysqladmin ping -h"${DB_HOST}" -u"${DB_USER}" -p"${DB_PASSWORD}" --silent > /dev/null 2>&1; then
            return 0
        fi
        retry_count=$((retry_count + 1))
        [ $retry_count -lt $max_retries ] && sleep 2
    done
    
    log_error "MySQL 连接失败，已重试 ${max_retries} 次"
    return 1
}

# 检查 Redis 是否健康
check_redis_health() {
    local max_attempts=30
    local attempt=1
    local wait_time=2

    if [ -z "${REDIS_PASSWORD}" ]; then
        while [ $attempt -le $max_attempts ]; do
            if redis-cli -h "${REDIS_HOST}" ping > /dev/null 2>&1; then
                echo "Redis 连接成功（无密码）"
                return 0
            fi
            echo "尝试连接 Redis 中... ($attempt/$max_attempts)"
            sleep $wait_time
            attempt=$((attempt + 1))
        done
    else
        while [ $attempt -le $max_attempts ]; do
            if redis-cli -h "${REDIS_HOST}" -a "${REDIS_PASSWORD}" ping > /dev/null 2>&1; then
                echo "Redis 连接成功（带密码）"
                return 0
            fi
            echo "尝试连接 Redis 中... ($attempt/$max_attempts)"
            sleep $wait_time
            attempt=$((attempt + 1))
        done
    fi

    echo "Redis 健康检查失败"
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

# 检查 MySQL 连接和初始化
echo "检查 MySQL 连接..."

# 首先使用 root 用户创建数据库和用户
echo "使用 root 用户初始化数据库..."
initialize_database

# 然后检查普通用户连接
if check_mysql_health; then
    echo "MySQL 连接正常"
    
    # 检查数据库是否需要初始化
    TABLES_EXIST=$(mysql -h"${DB_HOST}" -u"${DB_USER}" -p"${DB_PASSWORD}" -N -e "
        SELECT COUNT(*) FROM information_schema.tables 
        WHERE table_schema = '${DB_NAME}' 
        AND table_name IN ('users', 'todos', 'categories', 'reminders');
    " 2>/dev/null || echo "0")

    if [ "$TABLES_EXIST" = "0" ]; then
        echo "数据库为空，开始初始化..."
        mysql -h"${DB_HOST}" -u"${DB_USER}" -p"${DB_PASSWORD}" "${DB_NAME}" < scripts/init.sql
        if [ $? -eq 0 ]; then
            echo "数据库初始化完成"
        else
            echo "数据库初始化失败"
            exit 1
        fi
    else
        echo "数据库表已存在，跳过初始化"
    fi
else
    echo "无法连接到 MySQL，请检查服务是否运行以及连接参数是否正确"
    exit 1
fi

# 检查 Redis 连接
echo "检查 Redis 连接..."
if ! check_redis_health; then
    echo "无法连接到 Redis，请检查服务是否运行以及连接参数是否正确"
    exit 1
fi

# 添加配置文件验证函数
validate_config() {
    local config_file=$1
    
    if [ ! -f "$config_file" ]; then
        log_error "配置文件不存在: $config_file"
        return 1
    }
    
    # 验证配置文件语法（如果是 YAML）
    if command -v yamllint >/dev/null 2>&1; then
        if ! yamllint -d relaxed "$config_file"; then
            log_error "配置文件语法验证失败"
            return 1
        fi
    fi
    
    return 0
}

# 优化配置文件处理
process_config() {
    local config_file="configs/config.prod.yaml"
    local temp_file="configs/config.prod.yaml.tmp"
    
    log_info "处理配置文件..."
    
    if ! validate_config "$config_file"; then
        return 1
    fi
    
    if ! envsubst < "$config_file" > "$temp_file"; then
        log_error "环境变量替换失败"
        return 1
    }
    
    mv "$temp_file" "$config_file"
    log_info "配置文件处理完成"
}

# 优化应用部署函数
deploy_application() {
    log_info "开始部署应用..."
    
    # 停止并删除旧容器（如果存在）
    if docker ps -a | grep -q "todo-api"; then
        log_info "停止并删除旧容器..."
        docker stop todo-api || true
        docker rm todo-api || true
    fi
    
    # 构建新镜像
    log_info "构建应用镜像..."
    if ! docker-compose build --no-cache app; then
        log_error "应用构建失败"
        return 1
    fi
    
    # 启动应用
    log_info "启动应用容器..."
    if ! docker-compose up -d app; then
        log_error "应用启动失败"
        return 1
    fi
    
    log_info "应用部署完成"
    return 0
}

# 优化健康检查函数
check_application_health() {
    local max_attempts=30
    local attempt=1
    local wait_time=10
    
    log_info "等待应用就绪..."
    sleep $wait_time  # 初始等待
    
    while [ $attempt -le $max_attempts ]; do
        if curl -s -f http://localhost:8081/health | grep -q "healthy"; then
            log_info "应用已就绪"
            return 0
        fi
        
        # 检查容器状态
        local container_status=$(docker inspect --format='{{.State.Status}}' todo-api 2>/dev/null)
        if [ "$container_status" != "running" ]; then
            log_error "容器状态异常: $container_status"
            docker logs todo-api
            return 1
        fi
        
        log_info "等待应用就绪中... ($attempt/$max_attempts)"
        sleep 2
        attempt=$((attempt + 1))
    done
    
    log_error "应用未能在指定时间内就绪"
    docker logs todo-api
    return 1
}

# 主部署流程
main() {
    # 检查必要的环境变量
    check_required_env
    
    # 检查系统参数
    check_system_params
    
    # 检查 MySQL 连接
    log_info "检查 MySQL 连接..."
    if ! wait_for_mysql; then
        log_error "MySQL 连接失败"
        exit 1
    fi
    
    # 检查 Redis 连接
    log_info "检查 Redis 连接..."
    if ! wait_for_redis; then
        log_error "Redis 连接失败"
        exit 1
    fi
    
    # 处理配置文件
    if ! process_config; then
        log_error "配置文件处理失败"
        exit 1
    fi
    
    # 部署应用
    if ! deploy_application; then
        log_error "应用部署失败"
        exit 1
    fi
    
    # 检查应用健康状态
    if ! check_application_health; then
        log_error "应用健康检查失败"
        exit 1
    fi
    
    log_info "部署成功完成！"
}

# 执行主函数
main

# 优化数据库初始化部分
initialize_database() {
    log_info "开始初始化数据库..."
    
    # 使用 -v 参数增加详细输出
    if ! mysql -v -h"${DB_HOST}" -u"root" -p"${MYSQL_ROOT_PASSWORD}" -e "
        CREATE DATABASE IF NOT EXISTS ${DB_NAME};
        CREATE USER IF NOT EXISTS '${DB_USER}'@'%' IDENTIFIED BY '${DB_PASSWORD}';
        GRANT ALL PRIVILEGES ON ${DB_NAME}.* TO '${DB_USER}'@'%';
        FLUSH PRIVILEGES;
    "; then
        log_error "数据库初始化失败"
        return 1
    fi
    
    log_info "数据库初始化成功"
    return 0
}

# 添加清理函数
cleanup() {
    log_info "开始清理..."
    
    # 清理临时文件
    rm -f configs/*.tmp
    
    # 清理旧的日志文件（保留最近7天）
    find logs/ -type f -mtime +7 -delete
    
    log_info "清理完成"
}

# 在脚本结束时调用清理
trap cleanup EXIT 