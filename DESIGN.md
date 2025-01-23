# Todo API 设计文档

## 1. 业务功能设计

### 1.1 用户管理模块

1. 用户注册
   - 用户名唯一性校验
   - 密码加密存储 (bcrypt)
   - 邮箱格式验证
   - 返回注册成功消息
   - 发送欢迎邮件 (异步)

2. 用户登录
   - 用户名密码验证
   - JWT token 生成
   - 刷新 token 机制
   - 登录历史记录
   - 异常登录检测

3. 用户认证
   - JWT token 验证
   - token 过期处理
   - 用户身份识别

4. 用户信息管理
   - 获取用户信息
   - 修改用户信息
   - 修改密码
   - 头像上传

### 1.2 待办事项模块

1. 创建待办事项
   - 标题必填
   - 描述可选
   - 关联当前用户
   - 默认未完成状态

2. 查询待办事项
   - 获取单个待办事项详情
   - 获取用户所有待办事项列表
   - 支持分页查询
   - 支持状态筛选

3. 更新待办事项
   - 修改标题和描述
   - 更新完成状态
   - 仅允许创建者修改

4. 删除待办事项
   - 软删除实现
   - 仅允许创建者删除

5. 待办事项分类
   - 创建分类
   - 修改分类
   - 删除分类
   - 按分类查询

6. 待办事项优先级
   - 高优先级
   - 中优先级
   - 低优先级
   - 按优先级排序

7. 待办事项标签
   - 添加标签
   - 移除标签
   - 按标签筛选

8. 待办事项提醒
   - 设置提醒时间
   - 提醒方式选择
   - 重复提醒设置

## 2. API 设计

### 2.1 认证接口

#### 注册
```http
POST /api/v1/auth/register
Content-Type: application/json

Request:
{
    "username": "string",    // 必填，长度 3-32
    "password": "string",    // 必填，长度 6-32
    "email": "string"        // 必填，有效邮箱
}

Response:
{
    "message": "Registration successful"
}
```

#### 登录
```http
POST /api/v1/auth/login
Content-Type: application/json

Request:
{
    "username": "string",    // 必填
    "password": "string"     // 必填
}

Response:
{
    "token": "string"        // JWT token
}
```

#### 修改密码
```http
PUT /api/v1/auth/password
Authorization: Bearer {token}
Content-Type: application/json

Request:
{
    "old_password": "string",  // 必填
    "new_password": "string"   // 必填，长度 6-32
}

Response:
{
    "message": "Password updated successfully"
}
```

#### 获取用户信息
```http
GET /api/v1/users/me
Authorization: Bearer {token}

Response:
{
    "id": integer,
    "username": "string",
    "email": "string",
    "created_at": "datetime"
}
```

### 2.2 待办事项接口

#### 创建待办事项
```http
POST /api/v1/todos
Authorization: Bearer {token}
Content-Type: application/json

Request:
{
    "title": "string",       // 必填，长度 1-128
    "description": "string"  // 可选，长度 0-1024
}

Response:
{
    "id": integer,
    "title": "string",
    "description": "string",
    "completed": boolean,
    "created_at": "datetime",
    "updated_at": "datetime"
}
```

#### 获取待办事项列表
```http
GET /api/v1/todos
Authorization: Bearer {token}

Query Parameters:
- page: integer        // 页码，默认 1
- page_size: integer   // 每页数量，默认 10
- completed: boolean   // 完成状态筛选，可选

Response:
{
    "total": integer,
    "items": [
        {
            "id": integer,
            "title": "string",
            "description": "string",
            "completed": boolean,
            "created_at": "datetime",
            "updated_at": "datetime"
        }
    ]
}
```

#### 获取待办事项详情
```http
GET /api/v1/todos/{id}
Authorization: Bearer {token}

Response:
{
    "id": integer,
    "title": "string",
    "description": "string",
    "completed": boolean,
    "created_at": "datetime",
    "updated_at": "datetime"
}
```

#### 更新待办事项
```http
PUT /api/v1/todos/{id}
Authorization: Bearer {token}
Content-Type: application/json

Request:
{
    "title": "string",       // 可选
    "description": "string", // 可选
    "completed": boolean     // 可选
}

Response:
{
    "id": integer,
    "title": "string",
    "description": "string",
    "completed": boolean,
    "created_at": "datetime",
    "updated_at": "datetime"
}
```

#### 删除待办事项
```http
DELETE /api/v1/todos/{id}
Authorization: Bearer {token}

Response:
{
    "message": "Todo deleted successfully"
}
```

#### 创建分类
```http
POST /api/v1/categories
Authorization: Bearer {token}
Content-Type: application/json

Request:
{
    "name": "string",        // 必填，长度 1-32
    "color": "string"        // 可选，颜色代码
}

Response:
{
    "id": integer,
    "name": "string",
    "color": "string"
}
```

#### 设置提醒
```http
POST /api/v1/todos/{id}/reminder
Authorization: Bearer {token}
Content-Type: application/json

Request:
{
    "remind_at": "datetime",    // 必填
    "repeat_type": "string",    // 可选：none, daily, weekly, monthly
    "notify_type": "string"     // 可选：email, push
}

Response:
{
    "id": integer,
    "remind_at": "datetime",
    "repeat_type": "string",
    "notify_type": "string"
}
```

## 3. 数据库设计

### 3.1 用户表 (users)

| 字段名     | 类型         | 约束               | 说明     |
| ---------- | ------------ | ------------------ | -------- |
| id         | uint         | PK, AUTO_INCREMENT | 用户ID   |
| username   | varchar(32)  | UNIQUE, NOT NULL   | 用户名   |
| password   | varchar(128) | NOT NULL           | 密码哈希 |
| email      | varchar(128) | NOT NULL           | 邮箱地址 |
| created_at | datetime     | NOT NULL           | 创建时间 |
| updated_at | datetime     | NOT NULL           | 更新时间 |
| deleted_at | datetime     | NULL               | 删除时间 |

索引：
- PRIMARY KEY (`id`)
- UNIQUE KEY `idx_username` (`username`)
- KEY `idx_email` (`email`)

### 3.2 待办事项表 (todos)

| 字段名      | 类型          | 约束                    | 说明       |
| ----------- | ------------- | ----------------------- | ---------- |
| id          | uint          | PK, AUTO_INCREMENT      | 待办事项ID |
| user_id     | uint          | NOT NULL, FK            | 用户ID     |
| title       | varchar(128)  | NOT NULL                | 标题       |
| description | varchar(1024) | NULL                    | 描述       |
| completed   | boolean       | NOT NULL, DEFAULT false | 完成状态   |
| created_at  | datetime      | NOT NULL                | 创建时间   |
| updated_at  | datetime      | NOT NULL                | 更新时间   |
| deleted_at  | datetime      | NULL                    | 删除时间   |

索引：
- PRIMARY KEY (`id`)
- KEY `idx_user_id` (`user_id`)
- KEY `idx_created_at` (`created_at`)

外键约束：
- FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)

### 3.3 分类表 (categories)

| 字段名     | 类型        | 约束               | 说明     |
| ---------- | ----------- | ------------------ | -------- |
| id         | uint        | PK, AUTO_INCREMENT | 分类ID   |
| user_id    | uint        | NOT NULL, FK       | 用户ID   |
| name       | varchar(32) | NOT NULL           | 分类名称 |
| color      | varchar(7)  | NULL               | 颜色代码 |
| created_at | datetime    | NOT NULL           | 创建时间 |
| updated_at | datetime    | NOT NULL           | 更新时间 |
| deleted_at | datetime    | NULL               | 删除时间 |

索引：
- PRIMARY KEY (`id`)
- KEY `idx_user_id` (`user_id`)

### 3.4 提醒表 (reminders)

| 字段名      | 类型        | 约束               | 说明       |
| ----------- | ----------- | ------------------ | ---------- |
| id          | uint        | PK, AUTO_INCREMENT | 提醒ID     |
| todo_id     | uint        | NOT NULL, FK       | 待办事项ID |
| remind_at   | datetime    | NOT NULL           | 提醒时间   |
| repeat_type | varchar(16) | NOT NULL           | 重复类型   |
| notify_type | varchar(16) | NOT NULL           | 通知类型   |
| created_at  | datetime    | NOT NULL           | 创建时间   |
| updated_at  | datetime    | NOT NULL           | 更新时间   |

索引：
- PRIMARY KEY (`id`)
- KEY `idx_todo_id` (`todo_id`)
- KEY `idx_remind_at` (`remind_at`)

## 4. 安全设计

### 4.1 密码安全
- 使用 bcrypt 算法加密存储密码
- 密码长度和复杂度要求
- 登录失败次数限制
- 密码重置功能
- 密码修改限制
- 密码历史记录

### 4.2 认证安全
- JWT token 过期时间设置
- token 签名验证
- 敏感操作二次验证
- 多设备登录控制
- 登录设备管理
- 异常登录通知

### 4.3 数据安全
- 输入数据验证和清洗
- SQL 注入防护
- XSS 防护
- 数据备份策略
- 敏感信息加密
- 操作日志记录

## 5. 性能设计

### 5.1 数据库优化
- 合理的索引设计
- 分页查询实现
- 软删除机制
- 定期数据清理
- 热点数据缓存
- 读写分离

### 5.2 API 优化
- 请求频率限制
- 响应数据缓存
- 数据压缩
- 批量操作接口
- 异步处理机制
- 结果缓存策略

### 5.3 并发处理
- 连接池管理
- goroutine 控制
- 资源限制
- 分布式锁
- 任务队列
- 限流降级

### 5.4 监控告警
- 性能指标监控
- 错误日志告警
- 资源使用告警
- 业务指标监控

## 6. 技术架构设计

### 6.1 整体架构
- 采用分层架构:表现层、业务层、数据访问层
- 使用依赖注入解耦各层
- 统一错误处理机制
- 中间件链式处理

### 6.2 代码组织
- 按功能模块划分包
- 接口优先设计
- 依赖倒置原则
- 单一职责原则

### 6.3 并发处理
- goroutine 池化管理
- 基于 channel 的任务队列
- 分布式锁避免并发冲突
- 优雅关闭处理

### 6.4 缓存设计
- 多级缓存架构
- 缓存更新策略
- 缓存穿透防护
- 缓存雪崩防护

## 7. 部署架构

### 7.1 开发环境
- 本地开发环境
- 测试环境
- 预发布环境
- 生产环境

### 7.2 CI/CD 流程
- 代码检查
- 单元测试
- 集成测试
- 自动部署

### 7.3 监控告警
- 业务监控
- 性能监控
- 资源监控
- 安全监控

## 8. 扩展性设计

### 8.1 水平扩展
- 无状态服务设计
- 会话共享方案
- 数据分片策略
- 负载均衡

### 8.2 垂直扩展
- 微服务拆分
- 领域驱动设计
- 服务治理
- 服务编排
