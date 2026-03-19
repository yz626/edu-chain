-- =====================================================
-- EduChain 用户与权限模块 (优化后)
-- MySQL 8.0
-- 设计原则: 不使用外键, 数据表尽量精简
-- =====================================================

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- 1.1 用户基础表 (users)
-- 优化: 合并user_profiles字段到users表, 精简冗余字段
DROP TABLE IF EXISTS `users`;
CREATE TABLE `users` (
    `id` VARCHAR(36) NOT NULL COMMENT '用户ID (UUID)',
    `username` VARCHAR(64) NOT NULL COMMENT '用户名',
    `email` VARCHAR(128) NOT NULL COMMENT '邮箱地址',
    `phone` VARCHAR(20) DEFAULT NULL COMMENT '手机号',
    `password_hash` VARCHAR(255) NOT NULL COMMENT '加密后的密码哈希',
    `salt` VARCHAR(64) DEFAULT NULL COMMENT '密码盐值',
    `status` TINYINT NOT NULL DEFAULT 1 COMMENT '用户状态: 1-正常, 2-禁用, 3-待审核, 4-锁定',
    `user_type` TINYINT NOT NULL DEFAULT 1 COMMENT '用户类型: 1-普通用户, 2-管理员, 3-系统用户',
    `source` TINYINT DEFAULT 1 COMMENT '注册来源: 1-注册, 2-导入, 3-第三方, 4-API',
    `real_name` VARCHAR(64) DEFAULT NULL COMMENT '真实姓名',
    `nickname` VARCHAR(64) DEFAULT NULL COMMENT '昵称',
    `avatar_url` VARCHAR(512) DEFAULT NULL COMMENT '头像URL',
    `gender` TINYINT DEFAULT NULL COMMENT '性别: 0-未知, 1-男, 2-女',
    `birth_date` DATE DEFAULT NULL COMMENT '出生日期',
    `id_card_number` VARCHAR(18) DEFAULT NULL COMMENT '身份证号',
    `bio` TEXT COMMENT '个人简介',
    `address` TEXT COMMENT '地址',
    `last_login_at` DATETIME(3) DEFAULT NULL COMMENT '最后登录时间',
    `last_login_ip` VARCHAR(45) DEFAULT NULL COMMENT '最后登录IP',
    `login_count` INT DEFAULT 0 COMMENT '登录次数',
    `failed_attempts` TINYINT DEFAULT 0 COMMENT '失败尝试次数',
    `locked_until` DATETIME(3) DEFAULT NULL COMMENT '锁定截止时间',
    `email_verified` TINYINT(1) DEFAULT 0 COMMENT '邮箱是否已验证',
    `phone_verified` TINYINT(1) DEFAULT 0 COMMENT '手机是否已验证',
    `two_factor_enabled` TINYINT(1) DEFAULT 0 COMMENT '是否启用双因素认证',
    `two_factor_secret` VARCHAR(128) DEFAULT NULL COMMENT '双因素认证密钥',
    `extra_data` JSON DEFAULT NULL COMMENT '扩展数据',
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
    `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
    `deleted_at` DATETIME(3) DEFAULT NULL COMMENT '删除时间',
    `version` INT DEFAULT 1 COMMENT '版本号 (乐观锁)',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_users_username` (`username`),
    UNIQUE KEY `uk_users_email` (`email`),
    UNIQUE KEY `uk_users_phone` (`phone`),
    KEY `idx_users_status` (`status`),
    KEY `idx_users_user_type` (`user_type`),
    KEY `idx_users_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户基础表';

-- 1.2 角色表 (roles)
DROP TABLE IF EXISTS `roles`;
CREATE TABLE `roles` (
    `id` VARCHAR(36) NOT NULL COMMENT '角色ID (UUID)',
    `name` VARCHAR(64) NOT NULL COMMENT '角色名称',
    `code` VARCHAR(64) NOT NULL COMMENT '角色代码 (唯一标识)',
    `type` TINYINT NOT NULL DEFAULT 1 COMMENT '角色类型: 1-超级管理员, 2-审计员, 3-操作员, 4-校管理员, 5-院系管理员, 6-企业HR, 7-学生',
    `description` TEXT COMMENT '角色描述',
    `level` INT DEFAULT 100 COMMENT '角色级别 (数值越小级别越高)',
    `is_system` TINYINT(1) DEFAULT 0 COMMENT '是否系统角色',
    `is_active` TINYINT(1) DEFAULT 1 COMMENT '是否启用',
    `permissions` JSON DEFAULT NULL COMMENT '权限列表 (JSON数组)',
    `priority` INT DEFAULT 0 COMMENT '排序优先级',
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
    `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
    `deleted_at` DATETIME(3) DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_roles_code` (`code`),
    KEY `idx_roles_type` (`type`),
    KEY `idx_roles_level` (`level`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='角色表';

-- 1.3 权限表 (permissions)
DROP TABLE IF EXISTS `permissions`;
CREATE TABLE `permissions` (
    `id` VARCHAR(36) NOT NULL COMMENT '权限ID (UUID)',
    `code` VARCHAR(64) NOT NULL COMMENT '权限代码',
    `name` VARCHAR(64) NOT NULL COMMENT '权限名称',
    `resource` VARCHAR(64) NOT NULL COMMENT '资源类型',
    `action` VARCHAR(32) NOT NULL COMMENT '操作类型',
    `description` TEXT COMMENT '权限描述',
    `path_pattern` VARCHAR(256) DEFAULT NULL COMMENT 'API路径匹配模式',
    `method` VARCHAR(16) DEFAULT NULL COMMENT 'HTTP方法',
    `is_api` TINYINT(1) DEFAULT 1 COMMENT '是否为API权限',
    `is_menu` TINYINT(1) DEFAULT 0 COMMENT '是否菜单权限',
    `menu_icon` VARCHAR(64) DEFAULT NULL COMMENT '菜单图标',
    `menu_order` INT DEFAULT 0 COMMENT '菜单排序',
    `parent_id` VARCHAR(36) DEFAULT NULL COMMENT '父级权限ID',
    `level` INT DEFAULT 100 COMMENT '权限级别',
    `is_active` TINYINT(1) DEFAULT 1 COMMENT '是否启用',
    `extra_data` JSON DEFAULT NULL COMMENT '扩展数据',
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
    `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
    `deleted_at` DATETIME(3) DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_permissions_code` (`code`),
    KEY `idx_permissions_resource` (`resource`),
    KEY `idx_permissions_action` (`action`),
    KEY `idx_permissions_parent` (`parent_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='权限表';

-- 1.4 用户角色关联表 (user_roles)
DROP TABLE IF EXISTS `user_roles`;
CREATE TABLE `user_roles` (
    `id` VARCHAR(36) NOT NULL COMMENT '记录ID (UUID)',
    `user_id` VARCHAR(36) NOT NULL COMMENT '用户ID',
    `role_id` VARCHAR(36) NOT NULL COMMENT '角色ID',
    `granted_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '授权时间',
    `granted_by` VARCHAR(36) DEFAULT NULL COMMENT '授权人ID',
    `expires_at` DATETIME(3) DEFAULT NULL COMMENT '过期时间',
    `is_active` TINYINT(1) DEFAULT 1 COMMENT '是否有效',
    `remark` VARCHAR(256) DEFAULT NULL COMMENT '备注',
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
    `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_user_roles_user_role` (`user_id`, `role_id`),
    KEY `idx_user_roles_user` (`user_id`),
    KEY `idx_user_roles_role` (`role_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户角色关联表';

-- 1.5 角色权限关联表 (role_permissions)
DROP TABLE IF EXISTS `role_permissions`;
CREATE TABLE `role_permissions` (
    `id` VARCHAR(36) NOT NULL COMMENT '记录ID (UUID)',
    `role_id` VARCHAR(36) NOT NULL COMMENT '角色ID',
    `permission_id` VARCHAR(36) NOT NULL COMMENT '权限ID',
    `is_granted` TINYINT(1) DEFAULT 1 COMMENT '是否授予',
    `conditions` JSON DEFAULT NULL COMMENT '授权条件',
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_role_permissions` (`role_id`, `permission_id`),
    KEY `idx_role_permissions_role` (`role_id`),
    KEY `idx_role_permissions_permission` (`permission_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='角色权限关联表';

-- 1.6 令牌表 (refresh_tokens)
-- 优化: 精简字段,移除冗余设备信息
DROP TABLE IF EXISTS `refresh_tokens`;
CREATE TABLE `refresh_tokens` (
    `id` VARCHAR(36) NOT NULL COMMENT '令牌ID (UUID)',
    `user_id` VARCHAR(36) NOT NULL COMMENT '用户ID',
    `token_hash` VARCHAR(128) NOT NULL COMMENT '令牌哈希值',
    `expires_at` DATETIME(3) NOT NULL COMMENT '过期时间',
    `is_revoked` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否已撤销',
    `revoked_at` DATETIME(3) DEFAULT NULL COMMENT '撤销时间',
    `device_info` JSON DEFAULT NULL COMMENT '设备信息',
    `ip_address` VARCHAR(45) DEFAULT NULL COMMENT 'IP地址',
    `location` VARCHAR(128) DEFAULT NULL COMMENT '地理位置',
    `last_used_at` DATETIME(3) DEFAULT NULL COMMENT '最后使用时间',
    `use_count` INT DEFAULT 0 COMMENT '使用次数',
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_refresh_tokens_token_hash` (`token_hash`),
    KEY `idx_refresh_tokens_user` (`user_id`),
    KEY `idx_refresh_tokens_expires` (`expires_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='令牌表';

SET FOREIGN_KEY_CHECKS = 1;
