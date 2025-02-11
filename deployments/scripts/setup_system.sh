#!/bin/bash

# 检查系统参数状态
check_system_params() {
    echo "检查系统参数..."
    local issues=()
    
    # 检查内存分配策略
    if [ "$(sysctl -n vm.overcommit_memory)" != "1" ]; then
        issues+=("vm.overcommit_memory 未设置为推荐值 1")
    fi
    
    # 检查网络连接数
    if [ "$(sysctl -n net.core.somaxconn)" -lt "512" ]; then
        issues+=("net.core.somaxconn 小于推荐值 512")
    fi
    
    # 检查 THP 状态
    if [ -f /sys/kernel/mm/transparent_hugepage/enabled ]; then
        if ! grep -q "\[never\]" /sys/kernel/mm/transparent_hugepage/enabled; then
            issues+=("透明大页面(THP)未禁用")
        fi
    fi

    # 如果发现问题，返回非零值
    if [ ${#issues[@]} -gt 0 ]; then
        echo "❌ 发现以下问题："
        printf '%s\n' "${issues[@]}"
        return 1
    fi
    echo "✅ 系统参数检查通过"
    return 0
}

# 检查是否具有 root 权限
if [ "$EUID" -ne 0 ]; then
    echo "请使用 root 权限运行此脚本"
    exit 1
fi

# 检查当前系统参数状态
check_system_params
if [ $? -eq 0 ] && [ "$1" != "--force" ]; then
    echo "系统参数已经是最优配置，无需修改"
    exit 0
fi

echo "正在设置系统参数..."

# 1. 内存相关配置
# 设置 Redis 内存分配策略
echo "vm.overcommit_memory = 1" > /etc/sysctl.d/redis-memory.conf
sysctl -p /etc/sysctl.d/redis-memory.conf

# 禁用透明大页面(THP)
echo never > /sys/kernel/mm/transparent_hugepage/enabled
echo never > /sys/kernel/mm/transparent_hugepage/defrag

# 创建持久化的 THP 设置服务
cat > /etc/systemd/system/disable-thp.service << EOF
[Unit]
Description=禁用透明大页面(THP)
After=network.target

[Service]
Type=oneshot
ExecStart=/bin/sh -c 'echo never > /sys/kernel/mm/transparent_hugepage/enabled && echo never > /sys/kernel/mm/transparent_hugepage/defrag'
RemainAfterExit=yes

[Install]
WantedBy=multi-user.target
EOF

# 2. 网络相关配置
echo "net.core.somaxconn = 1024" > /etc/sysctl.d/network-tuning.conf
sysctl -p /etc/sysctl.d/network-tuning.conf

# 3. 启用服务
systemctl daemon-reload
systemctl enable disable-thp
systemctl start disable-thp

echo "系统参数设置完成"

# 验证设置
echo "验证系统参数设置："
check_system_params 