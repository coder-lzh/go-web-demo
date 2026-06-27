CREATE DATABASE IF NOT EXISTS `db_sass` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;

USE `db_sass`;

CREATE TABLE IF NOT EXISTS `user` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '用户ID',
    `username` VARCHAR(50) NOT NULL COMMENT '用户名',
    `password` VARCHAR(255) NOT NULL COMMENT '密码（加密存储）',
    `nickname` VARCHAR(50) DEFAULT '' COMMENT '昵称',
    `email` VARCHAR(100) DEFAULT '' COMMENT '邮箱',
    `phone` VARCHAR(20) DEFAULT '' COMMENT '手机号',
    `avatar` VARCHAR(255) DEFAULT '' COMMENT '头像URL',
    `gender` TINYINT UNSIGNED DEFAULT 0 COMMENT '性别：0-未知，1-男，2-女',
    `status` TINYINT UNSIGNED DEFAULT 1 COMMENT '状态：0-禁用，1-启用',
    `is_deleted` TINYINT UNSIGNED DEFAULT 0 COMMENT '是否删除：0-未删除，1-已删除',
    `deleted_at` DATETIME DEFAULT NULL COMMENT '删除时间',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    -- 联合唯一索引：未删除的记录必须唯一
    UNIQUE KEY `uk_username_deleted` (`username`, `deleted_at`),
    UNIQUE KEY `uk_email_deleted` (`email`, `deleted_at`),
    UNIQUE KEY `uk_phone_deleted` (`phone`, `deleted_at`),
    -- 普通索引用于查询
    KEY `idx_status` (`status`),
    KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='用户表';