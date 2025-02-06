-- 创建数据库
CREATE DATABASE IF NOT EXISTS todo_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE todo_db;

-- 创建用户表
CREATE TABLE IF NOT EXISTS users (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(32) NOT NULL UNIQUE,
    password VARCHAR(128) NOT NULL,
    email VARCHAR(128) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

-- 创建分类表
CREATE TABLE IF NOT EXISTS categories (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(32) NOT NULL,
    color VARCHAR(7),
    user_id BIGINT UNSIGNED NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    CONSTRAINT fk_categories_user FOREIGN KEY (user_id) REFERENCES users(id)
);

-- 创建待办事项表
CREATE TABLE IF NOT EXISTS todos (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    title VARCHAR(128) NOT NULL,
    description TEXT,
    completed BOOLEAN DEFAULT FALSE,
    priority VARCHAR(10) DEFAULT 'medium',
    user_id BIGINT UNSIGNED NOT NULL,
    category_id BIGINT UNSIGNED,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    CONSTRAINT fk_todos_user FOREIGN KEY (user_id) REFERENCES users(id),
    CONSTRAINT fk_todos_category FOREIGN KEY (category_id) REFERENCES categories(id)
);

-- 创建提醒表
CREATE TABLE IF NOT EXISTS reminders (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    todo_id BIGINT UNSIGNED NOT NULL,
    remind_at TIMESTAMP NOT NULL,
    remind_type VARCHAR(10) NOT NULL COMMENT 'once/daily/weekly',
    notify_type VARCHAR(10) NOT NULL COMMENT 'email/push',
    status BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    CONSTRAINT fk_reminders_todo FOREIGN KEY (todo_id) REFERENCES todos(id),
    CONSTRAINT chk_remind_type CHECK (remind_type IN ('once', 'daily', 'weekly')),
    CONSTRAINT chk_notify_type CHECK (notify_type IN ('email', 'push'))
);

-- 添加索引
CREATE INDEX idx_reminders_todo_id ON reminders(todo_id);
CREATE INDEX idx_reminders_remind_at ON reminders(remind_at);
CREATE INDEX idx_reminders_todo_remind ON reminders(todo_id, deleted_at);
CREATE INDEX idx_reminders_remind_status ON reminders(remind_at, status, deleted_at); 