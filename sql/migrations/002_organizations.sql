-- =====================================================
-- EduChain 组织架构模块 (优化后)
-- MySQL 8.0
-- 设计原则: 不使用外键, 数据表尽量精简
-- =====================================================

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- 2.1 组织表 (organizations)
DROP TABLE IF EXISTS `organizations`;
CREATE TABLE `organizations` (
    `id` VARCHAR(36) NOT NULL COMMENT '组织ID (UUID)',
    `name` VARCHAR(128) NOT NULL COMMENT '组织名称',
    `code` VARCHAR(64) NOT NULL COMMENT '组织代码 (唯一标识)',
    `type` TINYINT NOT NULL DEFAULT 5 COMMENT '组织类型: 1-监管机构, 2-高校, 3-企业, 4-政府机构, 5-其他',
    `short_name` VARCHAR(64) DEFAULT NULL COMMENT '组织简称',
    `logo_url` VARCHAR(512) DEFAULT NULL COMMENT 'Logo URL',
    `description` TEXT COMMENT '组织描述',
    `address` TEXT COMMENT '地址',
    `province` VARCHAR(32) DEFAULT NULL COMMENT '省份',
    `city` VARCHAR(32) DEFAULT NULL COMMENT '城市',
    `district` VARCHAR(32) DEFAULT NULL COMMENT '区县',
    `postal_code` VARCHAR(10) DEFAULT NULL COMMENT '邮政编码',
    `website` VARCHAR(256) DEFAULT NULL COMMENT '官方网站',
    `email` VARCHAR(128) DEFAULT NULL COMMENT '联系邮箱',
    `phone` VARCHAR(20) DEFAULT NULL COMMENT '联系电话',
    `contact_name` VARCHAR(64) DEFAULT NULL COMMENT '联系人姓名',
    `contact_phone` VARCHAR(20) DEFAULT NULL COMMENT '联系人电话',
    `contact_email` VARCHAR(128) DEFAULT NULL COMMENT '联系人邮箱',
    `industry` VARCHAR(64) DEFAULT NULL COMMENT '所属行业',
    `scale` VARCHAR(32) DEFAULT NULL COMMENT '组织规模',
    `established_date` DATE DEFAULT NULL COMMENT '成立日期',
    `is_active` TINYINT(1) DEFAULT 1 COMMENT '是否启用',
    `is_verified` TINYINT(1) DEFAULT 0 COMMENT '是否已认证',
    `verified_at` DATETIME(3) DEFAULT NULL COMMENT '认证时间',
    `parent_id` VARCHAR(36) DEFAULT NULL COMMENT '上级组织ID',
    `level` INT DEFAULT 1 COMMENT '组织层级',
    `sort_order` INT DEFAULT 0 COMMENT '排序序号',
    `extra_data` JSON DEFAULT NULL COMMENT '扩展数据',
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
    `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
    `deleted_at` DATETIME(3) DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_organizations_code` (`code`),
    KEY `idx_organizations_type` (`type`),
    KEY `idx_organizations_parent` (`parent_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='组织表';

-- 2.2 组织成员关联表 (organization_users)
DROP TABLE IF EXISTS `organization_users`;
CREATE TABLE `organization_users` (
    `id` VARCHAR(36) NOT NULL COMMENT '记录ID (UUID)',
    `organization_id` VARCHAR(36) NOT NULL COMMENT '组织ID',
    `user_id` VARCHAR(36) NOT NULL COMMENT '用户ID',
    `department_id` VARCHAR(36) DEFAULT NULL COMMENT '部门ID',
    `department_name` VARCHAR(128) DEFAULT NULL COMMENT '部门名称',
    `position` VARCHAR(64) DEFAULT NULL COMMENT '职位',
    `position_title` VARCHAR(64) DEFAULT NULL COMMENT '职称',
    `employee_id` VARCHAR(64) DEFAULT NULL COMMENT '员工号',
    `hire_date` DATE DEFAULT NULL COMMENT '入职日期',
    `resignation_date` DATE DEFAULT NULL COMMENT '离职日期',
    `employment_type` TINYINT DEFAULT 1 COMMENT '用工类型: 1-正式, 2-合同, 3-实习, 4-外包',
    `work_status` TINYINT DEFAULT 1 COMMENT '工作状态: 1-在职, 2-离职, 3-退休, 4-停薪留职',
    `is_primary` TINYINT(1) DEFAULT 0 COMMENT '是否为主组织',
    `is_admin` TINYINT(1) DEFAULT 0 COMMENT '是否管理员',
    `joined_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '加入时间',
    `left_at` DATETIME(3) DEFAULT NULL COMMENT '离开时间',
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
    `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_organization_users` (`organization_id`, `user_id`),
    KEY `idx_organization_users_org` (`organization_id`),
    KEY `idx_organization_users_user` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='组织成员关联表';

-- 2.3 部门表 (departments)
DROP TABLE IF EXISTS `departments`;
CREATE TABLE `departments` (
    `id` VARCHAR(36) NOT NULL COMMENT '部门ID (UUID)',
    `organization_id` VARCHAR(36) NOT NULL COMMENT '所属组织ID',
    `parent_id` VARCHAR(36) DEFAULT NULL COMMENT '上级部门ID',
    `name` VARCHAR(128) NOT NULL COMMENT '部门名称',
    `code` VARCHAR(64) DEFAULT NULL COMMENT '部门代码',
    `type` TINYINT DEFAULT 1 COMMENT '部门类型: 1-职能部门, 2-教学部门, 3-行政部门, 4-其他',
    `description` TEXT COMMENT '部门描述',
    `leader_id` VARCHAR(36) DEFAULT NULL COMMENT '负责人ID',
    `sort_order` INT DEFAULT 0 COMMENT '排序序号',
    `is_active` TINYINT(1) DEFAULT 1 COMMENT '是否启用',
    `level` INT DEFAULT 1 COMMENT '部门层级',
    `path` VARCHAR(512) DEFAULT NULL COMMENT '部门路径',
    `extra_data` JSON DEFAULT NULL COMMENT '扩展数据',
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
    `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
    `deleted_at` DATETIME(3) DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`),
    KEY `idx_departments_org` (`organization_id`),
    KEY `idx_departments_parent` (`parent_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='部门表';

SET FOREIGN_KEY_CHECKS = 1;
