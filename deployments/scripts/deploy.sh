#!/bin/bash
#
# 部署脚本 - 用于部署 Todo 应用
#
# 功能：
# - 检查系统环境和必要参数
# - 初始化和验证数据库连接
# - 部署和监控应用状态
#
# 使用方法：
# ./deploy.sh
#
# 环境变量：
# - MYSQL_ROOT_PASSWORD: MySQL root 密码
# - DB_PASSWORD: 应用数据库用户密码
# - JWT_SECRET: JWT 密钥
# - 其他可选环境变量见脚本内说明
#

set -e  # 确保脚本在出错时退出

###################
# 常量定义
###################
readonly MAX_RETRY=3
readonly MAX_WAIT_TIME=60
readonly APP_PORT=8081
readonly CONFIG_FILE="configs/config.prod.yaml"

###################
# 日志和错误处理
###################
# 错误处理函数
handle_error() {
    local exit_code=$?
    local line_number=$1
    log_error "错误发生在第 ${line_number} 行，退出码: ${exit_code}"
    exit $exit_code
}

# 日志函数
log_info() {
    echo "[INFO] $(date '+%Y-%m-%d %H:%M:%S') - $1"
}

log_error() {
    echo "[ERROR] $(date '+%Y-%m-%d %H:%M:%S') - $1" >&2
}

log_warning() {
    echo "[WARN] $(date '+%Y-%m-%d %H:%M:%S') - $1" >&2
}

# 设置错误处理
trap 'handle_error ${LINENO}' ERR

###################
# 环境检查函数
###################
check_required_env() {
    local required_vars=("DB_HOST" "DB_USER" "MYSQL_ROOT_PASSWORD" "DB_PASSWORD" "JWT_SECRET")
    local missing_vars=()
    
    for var in "${required_vars[@]}"; do
        if [ -z "${!var}" ]; then
            missing_vars+=("$var")
        fi
    done
    
    if [ ${#missing_vars[@]} -ne 0 ]; then
        log_error "缺少必要的环境变量:"
        printf '%s\n' "${missing_vars[@]}"
        return 1
    fi
    
    # 输出环境变量值（敏感信息部分遮蔽）
    log_info "环境变量检查通过，当前配置:"
    echo "DB_HOST=${DB_HOST}"
    echo "DB_USER=${DB_USER}"
    # 对密码类信息只显示前两位和后两位
    echo "MYSQL_ROOT_PASSWORD=${MYSQL_ROOT_PASSWORD:0:2}****${MYSQL_ROOT_PASSWORD: -2}"
    echo "DB_PASSWORD=${DB_PASSWORD:0:2}****${DB_PASSWORD: -2}"
    echo "JWT_SECRET=${JWT_SECRET:0:2}****${JWT_SECRET: -2}"
    
    return 0
}

###################
# 数据库相关函数
###################
check_mysql_health() {
    local retry_count=0
    
    while [ $retry_count -lt $MAX_RETRY ]; do
        if mysqladmin ping -h"${DB_HOST}" -u"${DB_USER}" -p"${DB_PASSWORD}" --silent > /dev/null 2>&1; then
            return 0
        fi
        retry_count=$((retry_count + 1))
        [ $retry_count -lt $MAX_RETRY ] && sleep 2
    done
    
    log_error "MySQL 连接失败，已重试 ${MAX_RETRY} 次"
    return 1
}

initialize_database() {
    log_info "初始化数据库..."
    
    local sql_commands="
        CREATE DATABASE IF NOT EXISTS ${DB_NAME};
        CREATE USER IF NOT EXISTS '${DB_USER}'@'%' IDENTIFIED BY '${DB_PASSWORD}';
        GRANT ALL PRIVILEGES ON ${DB_NAME}.* TO '${DB_USER}'@'%';
        FLUSH PRIVILEGES;
    "
    
    if ! mysql -v -h"${DB_HOST}" -u"root" -p"${MYSQL_ROOT_PASSWORD}" -e "$sql_commands"; then
        log_error "数据库初始化失败"
        return 1
    fi
    
    log_info "数据库初始化成功"
    return 0
}

initialize_tables() {
    log_info "检查数据库表..."
    
    local tables_exist=$(mysql -h"${DB_HOST}" -u"${DB_USER}" -p"${DB_PASSWORD}" -N -e "
        SELECT COUNT(*) FROM information_schema.tables 
        WHERE table_schema = '${DB_NAME}' 
        AND table_name IN ('users', 'todos', 'categories', 'reminders');
    " 2>/dev/null || echo "0")

    if [ "$tables_exist" = "0" ]; then
        log_info "初始化数据库表..."
        if ! mysql -h"${DB_HOST}" -u"${DB_USER}" -p"${DB_PASSWORD}" "${DB_NAME}" < scripts/init.sql; then
            log_error "数据库表初始化失败"
            return 1
        fi
        log_info "数据库表初始化完成"
    else
        log_info "数据库表已存在，跳过初始化"
    fi
}

###################
# Redis 相关函数
###################
check_redis_health() {
    local redis_cmd="redis-cli -h ${REDIS_HOST}"
    [ -n "${REDIS_PASSWORD}" ] && redis_cmd+=" -a ${REDIS_PASSWORD}"
    
    if $redis_cmd ping > /dev/null 2>&1; then
        log_info "Redis 连接成功"
        return 0
    fi
    
    log_error "Redis 连接失败"
    return 1
}

###################
# 配置文件处理
###################
validate_config() {
    local config_file=$1
    
    if [ ! -f "$config_file" ]; then
        log_error "配置文件不存在: $config_file"
        return 1
    fi
    
    if command -v yamllint >/dev/null 2>&1; then
        if ! yamllint -d relaxed "$config_file"; then
            log_error "配置文件语法验证失败"
            return 1
        fi
    fi
    
    return 0
}

process_config() {
    log_info "处理配置文件..."
    
    if ! validate_config "$CONFIG_FILE"; then
        return 1
    fi
    
    local temp_file="${CONFIG_FILE}.tmp"
    if ! envsubst < "$CONFIG_FILE" > "$temp_file"; then
        log_error "环境变量替换失败"
        return 1
    fi
    
    mv "$temp_file" "$CONFIG_FILE"
    log_info "配置文件处理完成"
}

###################
# 应用部署函数
###################
deploy_application() {
    log_info "开始部署应用..."
    
    # 清理旧容器
    if docker ps -a | grep -q "todo-api"; then
        log_info "清理旧容器..."
        docker stop todo-api || true
        docker rm todo-api || true
    fi
    
    # 构建和启动
    log_info "构建应用镜像..."
    if ! docker-compose build --no-cache app; then
        log_error "应用构建失败"
        return 1
    fi
    
    log_info "启动应用容器..."
    if ! docker-compose up -d app; then
        log_error "应用启动失败"
        return 1
    fi
    
    log_info "应用部署完成"
    return 0
}

check_application_health() {
    local max_attempts=30
    local attempt=1
    
    log_info "等待应用就绪..."
    sleep 10  # 初始等待
    
    while [ $attempt -le $max_attempts ]; do
        # 检查健康检查接口
        if curl -s -f "http://localhost:${APP_PORT}/health" | grep -q "healthy"; then
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

###################
# 清理函数
###################
cleanup() {
    log_info "执行清理操作..."
    rm -f configs/*.tmp
    find logs/ -type f -mtime +7 -delete
    log_info "清理完成"
}

###################
# 主函数
###################
main() {
    # 1. 环境检查
    check_required_env || exit 1
    
    # 2. 设置环境变量
    export DB_HOST=${DB_HOST:-localhost}
    export DB_PORT=${DB_PORT:-3306}
    export DB_USER=${DB_USER:-todo_user}
    export DB_NAME=${DB_NAME:-todo_db}
    export REDIS_HOST=${REDIS_HOST:-localhost}
    export REDIS_PORT=${REDIS_PORT:-6379}
    
    # 3. 数据库初始化和检查
    initialize_database || exit 1
    initialize_tables || exit 1
    
    # 4. Redis 检查
    check_redis_health || exit 1
    
    # 5. 配置文件处理
    process_config || exit 1
    
    # 6. 应用部署
    deploy_application || exit 1
    
    # 7. 健康检查
    check_application_health || exit 1
    
    log_info "部署成功完成！"
}

# 注册清理函数
trap cleanup EXIT

# 执行主函数
main 