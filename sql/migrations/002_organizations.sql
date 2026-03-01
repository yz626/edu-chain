-- =====================================================
-- EduChain 组织管理模块数据库表
-- MySQL 8.0
-- =====================================================

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- =====================================================
-- 2.1 组织表 (organizations)
-- =====================================================
DROP TABLE IF EXISTS `organizations`;
CREATE TABLE `organizations` (
    `id` VARCHAR(36) NOT NULL COMMENT '组织ID (UUID)',
    `name` VARCHAR(128) NOT NULL COMMENT '组织名称',
    `code` VARCHAR(64) NOT NULL COMMENT '组织代码 (唯一标识)',
    `type` TINYINT NOT NULL DEFAULT 5 COMMENT '组织类型: 1-监管机构, 2-高校, 3-企业, 4-政府机构, 5-其他',
    `short_name` VARCHAR(64) DEFAULT NULL COMMENT '组织简称',
    `english_name` VARCHAR(128) DEFAULT NULL COMMENT '英文名称',
    `logo_url` VARCHAR(512) DEFAULT NULL COMMENT 'Logo URL',
    `banner_url` VARCHAR(512) DEFAULT NULL COMMENT '横幅图片URL',
    `description` TEXT COMMENT '组织描述',
    `address` TEXT COMMENT '地址',
    `province` VARCHAR(32) DEFAULT NULL COMMENT '省份',
    `city` VARCHAR(32) DEFAULT NULL COMMENT '城市',
    `district` VARCHAR(32) DEFAULT NULL COMMENT '区县',
    `postal_code` VARCHAR(10) DEFAULT NULL COMMENT '邮政编码',
    `website` VARCHAR(256) DEFAULT NULL COMMENT '官方网站',
    `email` VARCHAR(128) DEFAULT NULL COMMENT '联系邮箱',
    `phone` VARCHAR(20) DEFAULT NULL COMMENT '联系电话',
    `fax` VARCHAR(20) DEFAULT NULL COMMENT '传真号码',
    `contact_name` VARCHAR(64) DEFAULT NULL COMMENT '联系人姓名',
    `contact_phone` VARCHAR(20) DEFAULT NULL COMMENT '联系人电话',
    `contact_email` VARCHAR(128) DEFAULT NULL COMMENT '联系人邮箱',
    `contact_position` VARCHAR(64) DEFAULT NULL COMMENT '联系人职位',
    `industry` VARCHAR(64) DEFAULT NULL COMMENT '所属行业',
    `scale` VARCHAR(32) DEFAULT NULL COMMENT '组织规模',
    `established_date` DATE DEFAULT NULL COMMENT '成立日期',
    `is_active` TINYINT(1) DEFAULT 1 COMMENT '是否启用: 0-禁用, 1-启用',
    `is_verified` TINYINT(1) DEFAULT 0 COMMENT '是否已认证: 0-未认证, 1-已认证',
    `verified_at` DATETIME(3) DEFAULT NULL COMMENT '认证时间',
    `verified_by` VARCHAR(36) DEFAULT NULL COMMENT '认证人ID',
    `parent_id` VARCHAR(36) DEFAULT NULL COMMENT '上级组织ID',
    `level` INT DEFAULT 1 COMMENT '组织层级 (1=一级)',
    `sort_order` INT DEFAULT 0 COMMENT '排序序号',
    `extra_data` JSON DEFAULT NULL COMMENT '扩展数据 (JSON格式)',
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
    `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
    `deleted_at` DATETIME(3) DEFAULT NULL COMMENT '删除时间',
    `version` INT DEFAULT 1 COMMENT '版本号 (乐观锁)',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_organizations_code` (`code`),
    KEY `idx_organizations_type` (`type`),
    KEY `idx_organizations_active` (`is_active`),
    KEY `idx_organizations_parent` (`parent_id`),
    KEY `idx_organizations_name` (`name`(64)),
    KEY `idx_organizations_province_city` (`province`, `city`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='组织表 (高校/企业/机构)';

-- =====================================================
-- 2.2 组织成员关联表 (organization_users)
-- =====================================================
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
    `is_primary` TINYINT(1) DEFAULT 0 COMMENT '是否为主组织: 0-否, 1-是',
    `is_admin` TINYINT(1) DEFAULT 0 COMMENT '是否管理员: 0-否, 1-是',
    `joined_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '加入时间',
    `left_at` DATETIME(3) DEFAULT NULL COMMENT '离开时间',
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
    `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_organization_users` (`organization_id`, `user_id`),
    KEY `idx_organization_users_org` (`organization_id`),
    KEY `idx_organization_users_user` (`user_id`),
    KEY `idx_organization_users_dept` (`department_id`),
    KEY `idx_organization_users_primary` (`is_primary`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='组织成员关联表';

-- =====================================================
-- 2.3 部门表 (departments)
-- =====================================================
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
    `deputy_ids` JSON DEFAULT NULL COMMENT '副职负责人ID列表',
    `sort_order` INT DEFAULT 0 COMMENT '排序序号',
    `is_active` TINYINT(1) DEFAULT 1 COMMENT '是否启用: 0-禁用, 1-启用',
    `level` INT DEFAULT 1 COMMENT '部门层级',
    `path` VARCHAR(512) DEFAULT NULL COMMENT '部门路径 (如: 学院/计算机系)',
    `extra_data` JSON DEFAULT NULL COMMENT '扩展数据 (JSON格式)',
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
    `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
    `deleted_at` DATETIME(3) DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`),
    KEY `idx_departments_org` (`organization_id`),
    KEY `idx_departments_parent` (`parent_id`),
    KEY `idx_departments_path` (`path`(255))
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='部门表';

SET FOREIGN_KEY_CHECKS = 1;