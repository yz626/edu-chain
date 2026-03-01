-- =====================================================
-- EduChain 用户与认证模块数据库表
-- MySQL 8.0
-- =====================================================

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- =====================================================
-- 1.1 用户基础表 (users)
-- =====================================================
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
    `last_login_at` DATETIME(3) DEFAULT NULL COMMENT '最后登录时间',
    `last_login_ip` VARCHAR(45) DEFAULT NULL COMMENT '最后登录IP',
    `login_count` INT DEFAULT 0 COMMENT '登录次数',
    `failed_attempts` TINYINT DEFAULT 0 COMMENT '失败尝试次数',
    `locked_until` DATETIME(3) DEFAULT NULL COMMENT '锁定截止时间',
    `email_verified` TINYINT(1) DEFAULT 0 COMMENT '邮箱是否已验证: 0-未验证, 1-已验证',
    `phone_verified` TINYINT(1) DEFAULT 0 COMMENT '手机是否已验证: 0-未验证, 1-已验证',
    `two_factor_enabled` TINYINT(1) DEFAULT 0 COMMENT '是否启用双因素认证: 0-禁用, 1-启用',
    `two_factor_secret` VARCHAR(128) DEFAULT NULL COMMENT '双因素认证密钥',
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
    KEY `idx_users_created_at` (`created_at`),
    KEY `idx_users_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户基础表';

-- =====================================================
-- 1.2 用户扩展信息表 (user_profiles)
-- =====================================================
DROP TABLE IF EXISTS `user_profiles`;
CREATE TABLE `user_profiles` (
    `id` VARCHAR(36) NOT NULL COMMENT '记录ID (UUID)',
    `user_id` VARCHAR(36) NOT NULL COMMENT '用户ID',
    `real_name` VARCHAR(64) DEFAULT NULL COMMENT '真实姓名',
    `nickname` VARCHAR(64) DEFAULT NULL COMMENT '昵称',
    `avatar_url` VARCHAR(512) DEFAULT NULL COMMENT '头像URL',
    `gender` TINYINT DEFAULT NULL COMMENT '性别: 0-未知, 1-男, 2-女',
    `birth_date` DATE DEFAULT NULL COMMENT '出生日期',
    `id_card_number` VARCHAR(18) DEFAULT NULL COMMENT '身份证号',
    `id_card_front` VARCHAR(512) DEFAULT NULL COMMENT '身份证正面图片URL',
    `id_card_back` VARCHAR(512) DEFAULT NULL COMMENT '身份证背面图片URL',
    `bio` TEXT COMMENT '个人简介',
    `address` TEXT COMMENT '地址',
    `hometown` VARCHAR(128) DEFAULT NULL COMMENT '籍贯',
    `blood_type` TINYINT DEFAULT NULL COMMENT '血型: 1-A, 2-B, 3-AB, 4-O, 5-未知',
    `emergency_contact` VARCHAR(128) DEFAULT NULL COMMENT '紧急联系人',
    `emergency_phone` VARCHAR(20) DEFAULT NULL COMMENT '紧急联系电话',
    `extra_data` JSON DEFAULT NULL COMMENT '扩展数据 (JSON格式)',
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
    `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_user_profiles_user_id` (`user_id`),
    KEY `idx_user_profiles_real_name` (`real_name`),
    KEY `idx_user_profiles_id_card` (`id_card_number`(18))
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户扩展信息表';

-- =====================================================
-- 1.3 角色表 (roles)
-- =====================================================
DROP TABLE IF EXISTS `roles`;
CREATE TABLE `roles` (
    `id` VARCHAR(36) NOT NULL COMMENT '角色ID (UUID)',
    `name` VARCHAR(64) NOT NULL COMMENT '角色名称',
    `code` VARCHAR(64) NOT NULL COMMENT '角色代码 (唯一标识)',
    `type` TINYINT NOT NULL DEFAULT 1 COMMENT '角色类型: 1-超级管理员, 2-审计员, 3-操作员, 4-校管理员, 5-院系管理员, 6-企业HR, 7-学生',
    `description` TEXT COMMENT '角色描述',
    `level` INT DEFAULT 100 COMMENT '角色级别 (数值越小级别越高)',
    `is_system` TINYINT(1) DEFAULT 0 COMMENT '是否系统角色: 0-否, 1-是',
    `is_active` TINYINT(1) DEFAULT 1 COMMENT '是否启用: 0-禁用, 1-启用',
    `permissions` JSON DEFAULT NULL COMMENT '权限列表 (JSON数组)',
    `priority` INT DEFAULT 0 COMMENT '排序优先级',
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
    `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
    `deleted_at` DATETIME(3) DEFAULT NULL COMMENT '删除时间',
    `version` INT DEFAULT 1 COMMENT '版本号 (乐观锁)',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_roles_code` (`code`),
    KEY `idx_roles_type` (`type`),
    KEY `idx_roles_level` (`level`),
    KEY `idx_roles_active` (`is_active`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='角色表';

-- =====================================================
-- 1.4 权限表 (permissions)
-- =====================================================
DROP TABLE IF EXISTS `permissions`;
CREATE TABLE `permissions` (
    `id` VARCHAR(36) NOT NULL COMMENT '权限ID (UUID)',
    `code` VARCHAR(64) NOT NULL COMMENT '权限代码',
    `name` VARCHAR(64) NOT NULL COMMENT '权限名称',
    `resource` VARCHAR(64) NOT NULL COMMENT '资源类型: user, certificate, verification, role, organization等',
    `action` VARCHAR(32) NOT NULL COMMENT '操作类型: create, read, update, delete, list, export等',
    `description` TEXT COMMENT '权限描述',
    `path_pattern` VARCHAR(256) DEFAULT NULL COMMENT 'API路径匹配模式',
    `method` VARCHAR(16) DEFAULT NULL COMMENT 'HTTP方法: GET, POST, PUT, DELETE等',
    `is_api` TINYINT(1) DEFAULT 1 COMMENT '是否为API权限: 0-否, 1-是',
    `is_menu` TINYINT(1) DEFAULT 0 COMMENT '是否菜单权限: 0-否, 1-是',
    `menu_icon` VARCHAR(64) DEFAULT NULL COMMENT '菜单图标',
    `menu_order` INT DEFAULT 0 COMMENT '菜单排序',
    `parent_id` VARCHAR(36) DEFAULT NULL COMMENT '父级权限ID',
    `level` INT DEFAULT 100 COMMENT '权限级别',
    `is_active` TINYINT(1) DEFAULT 1 COMMENT '是否启用: 0-禁用, 1-启用',
    `extra_data` JSON DEFAULT NULL COMMENT '扩展数据 (JSON格式)',
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
    `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
    `deleted_at` DATETIME(3) DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_permissions_code` (`code`),
    KEY `idx_permissions_resource` (`resource`),
    KEY `idx_permissions_action` (`action`),
    KEY `idx_permissions_parent` (`parent_id`),
    UNIQUE KEY `uk_permissions_resource_action` (`resource`, `action`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='权限表';

-- =====================================================
-- 1.5 用户角色关联表 (user_roles)
-- =====================================================
DROP TABLE IF EXISTS `user_roles`;
CREATE TABLE `user_roles` (
    `id` VARCHAR(36) NOT NULL COMMENT '记录ID (UUID)',
    `user_id` VARCHAR(36) NOT NULL COMMENT '用户ID',
    `role_id` VARCHAR(36) NOT NULL COMMENT '角色ID',
    `granted_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '授权时间',
    `granted_by` VARCHAR(36) DEFAULT NULL COMMENT '授权人ID',
    `expires_at` DATETIME(3) DEFAULT NULL COMMENT '过期时间 (可选)',
    `is_active` TINYINT(1) DEFAULT 1 COMMENT '是否有效: 0-无效, 1-有效',
    `remark` VARCHAR(256) DEFAULT NULL COMMENT '备注',
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
    `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_user_roles_user_role` (`user_id`, `role_id`),
    KEY `idx_user_roles_user` (`user_id`),
    KEY `idx_user_roles_role` (`role_id`),
    KEY `idx_user_roles_expires` (`expires_at`),
    KEY `idx_user_roles_active` (`is_active`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户角色关联表';

-- =====================================================
-- 1.6 角色权限关联表 (role_permissions)
-- =====================================================
DROP TABLE IF EXISTS `role_permissions`;
CREATE TABLE `role_permissions` (
    `id` VARCHAR(36) NOT NULL COMMENT '记录ID (UUID)',
    `role_id` VARCHAR(36) NOT NULL COMMENT '角色ID',
    `permission_id` VARCHAR(36) NOT NULL COMMENT '权限ID',
    `is_granted` TINYINT(1) DEFAULT 1 COMMENT '是否授予: 0-未授予, 1-已授予',
    `conditions` JSON DEFAULT NULL COMMENT '授权条件 (JSON格式)',
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_role_permissions` (`role_id`, `permission_id`),
    KEY `idx_role_permissions_role` (`role_id`),
    KEY `idx_role_permissions_permission` (`permission_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='角色权限关联表';

-- =====================================================
-- 1.7 登录日志表 (login_logs)
-- =====================================================
DROP TABLE IF EXISTS `login_logs`;
CREATE TABLE `login_logs` (
    `id` VARCHAR(36) NOT NULL COMMENT '日志ID (UUID)',
    `user_id` VARCHAR(36) DEFAULT NULL COMMENT '用户ID (登录失败可能为空)',
    `session_id` VARCHAR(36) DEFAULT NULL COMMENT '会话ID',
    `ip_address` VARCHAR(45) NOT NULL COMMENT 'IP地址',
    `port` INT DEFAULT NULL COMMENT '端口号',
    `user_agent` TEXT COMMENT 'User-Agent字符串',
    `device_type` VARCHAR(32) DEFAULT NULL COMMENT '设备类型: pc, mobile, tablet',
    `browser` VARCHAR(64) DEFAULT NULL COMMENT '浏览器名称',
    `os` VARCHAR(64) DEFAULT NULL COMMENT '操作系统',
    `platform` VARCHAR(64) DEFAULT NULL COMMENT '平台信息',
    `location` VARCHAR(128) DEFAULT NULL COMMENT '登录地点',
    `success` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否登录成功: 0-失败, 1-成功',
    `failure_reason` VARCHAR(128) DEFAULT NULL COMMENT '失败原因',
    `login_type` TINYINT DEFAULT 1 COMMENT '登录类型: 1-密码, 2-验证码, 3-第三方, 4-Token',
    `captcha_required` TINYINT(1) DEFAULT 0 COMMENT '是否需要验证码: 0-否, 1-是',
    `mfa_required` TINYINT(1) DEFAULT 0 COMMENT '是否需要MFA: 0-否, 1-是',
    `mfa_verified` TINYINT(1) DEFAULT 0 COMMENT 'MFA是否已验证: 0-未验证, 1-已验证',
    `request_id` VARCHAR(36) DEFAULT NULL COMMENT '请求ID',
    `trace_id` VARCHAR(36) DEFAULT NULL COMMENT '链路追踪ID',
    `login_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '登录时间',
    `logout_at` DATETIME(3) DEFAULT NULL COMMENT '登出时间',
    `session_duration` INT DEFAULT NULL COMMENT '会话持续时间 (秒)',
    PRIMARY KEY (`id`),
    KEY `idx_login_logs_user` (`user_id`),
    KEY `idx_login_logs_ip` (`ip_address`(45)),
    KEY `idx_login_logs_success` (`success`),
    KEY `idx_login_logs_login_at` (`login_at`),
    KEY `idx_login_logs_user_at` (`user_id`, `login_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='登录日志表';

-- =====================================================
-- 1.8 刷新令牌表 (refresh_tokens)
-- =====================================================
DROP TABLE IF EXISTS `refresh_tokens`;
CREATE TABLE `refresh_tokens` (
    `id` VARCHAR(36) NOT NULL COMMENT '令牌ID (UUID)',
    `user_id` VARCHAR(36) NOT NULL COMMENT '用户ID',
    `token_hash` VARCHAR(128) NOT NULL COMMENT '令牌哈希值',
    `token_signature` VARCHAR(256) DEFAULT NULL COMMENT '令牌签名',
    `expires_at` DATETIME(3) NOT NULL COMMENT '过期时间',
    `is_revoked` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否已撤销: 0-否, 1-是',
    `revoked_at` DATETIME(3) DEFAULT NULL COMMENT '撤销时间',
    `revoked_by` VARCHAR(36) DEFAULT NULL COMMENT '撤销人ID',
    `revoke_reason` VARCHAR(256) DEFAULT NULL COMMENT '撤销原因',
    `device_type` VARCHAR(32) DEFAULT NULL COMMENT '设备类型',
    `device_name` VARCHAR(128) DEFAULT NULL COMMENT '设备名称',
    `ip_address` VARCHAR(45) DEFAULT NULL COMMENT 'IP地址',
    `user_agent` TEXT COMMENT 'User-Agent',
    `location` VARCHAR(128) DEFAULT NULL COMMENT '地理位置',
    `last_used_at` DATETIME(3) DEFAULT NULL COMMENT '最后使用时间',
    `use_count` INT DEFAULT 0 COMMENT '使用次数',
    `max_use_count` INT DEFAULT 1 COMMENT '最大使用次数',
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_refresh_tokens_token_hash` (`token_hash`),
    KEY `idx_refresh_tokens_user` (`user_id`),
    KEY `idx_refresh_tokens_expires` (`expires_at`),
    KEY `idx_refresh_tokens_revoked` (`is_revoked`, `expires_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='刷新令牌表 (JWT Token管理)';

-- =====================================================
-- 1.9 验证码表 (captchas)
-- =====================================================
DROP TABLE IF EXISTS `captchas`;
CREATE TABLE `captchas` (
    `id` VARCHAR(36) NOT NULL COMMENT '验证码ID (UUID)',
    `captcha_key` VARCHAR(64) NOT NULL COMMENT '验证码键',
    `captcha_code` VARCHAR(16) NOT NULL COMMENT '验证码内容',
    `captcha_image` TEXT COMMENT '验证码图片 (Base64编码)',
    `expire_at` DATETIME(3) NOT NULL COMMENT '过期时间',
    `is_used` TINYINT(1) DEFAULT 0 COMMENT '是否已使用: 0-否, 1-是',
    `used_at` DATETIME(3) DEFAULT NULL COMMENT '使用时间',
    `used_count` TINYINT DEFAULT 0 COMMENT '已使用次数',
    `max_use_count` TINYINT DEFAULT 1 COMMENT '最大使用次数',
    `ip_address` VARCHAR(45) DEFAULT NULL COMMENT '请求IP',
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_captchas_key` (`captcha_key`),
    KEY `idx_captchas_expire` (`expire_at`),
    KEY `idx_captchas_used` (`is_used`, `expire_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='验证码表';

SET FOREIGN_KEY_CHECKS = 1;