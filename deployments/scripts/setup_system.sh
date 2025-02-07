#!/bin/bash

# 设置系统参数
setup_system() {
    # 检查是否有 root 权限
    if [ "$EUID" -ne 0 ]; then
        echo "Please run as root"
        exit 1
    }

    # 设置系统参数
    cat > /etc/sysctl.d/99-redis.conf << EOF
# Redis 系统优化配置
vm.overcommit_memory = 1
net.core.somaxconn = 512
EOF

    # 应用系统参数
    sysctl -p /etc/sysctl.d/99-redis.conf

    # 禁用 THP
    cat > /etc/systemd/system/disable-thp.service << EOF
[Unit]
Description=Disable Transparent Huge Pages (THP)

[Service]
Type=oneshot
ExecStart=/bin/sh -c 'echo never > /sys/kernel/mm/transparent_hugepage/enabled && echo never > /sys/kernel/mm/transparent_hugepage/defrag'
RemainAfterExit=yes

[Install]
WantedBy=multi-user.target
EOF

    # 启用并启动服务
    systemctl daemon-reload
    systemctl enable disable-thp
    systemctl start disable-thp
}

# 运行设置
setup_system 