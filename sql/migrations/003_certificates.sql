-- =====================================================
-- EduChain 证书管理模块 (优化后)
-- MySQL 8.0
-- 设计原则: 不使用外键, 数据表尽量精简
-- =====================================================

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- 3.1 证书类型表 (certificate_types)
DROP TABLE IF EXISTS `certificate_types`;
CREATE TABLE `certificate_types` (
    `id` VARCHAR(36) NOT NULL COMMENT '类型ID (UUID)',
    `code` VARCHAR(32) NOT NULL COMMENT '类型代码 (唯一标识)',
    `name` VARCHAR(64) NOT NULL COMMENT '类型名称',
    `category` TINYINT NOT NULL DEFAULT 1 COMMENT '证书类别: 1-毕业证书, 2-学位证书, 3-成绩单, 4-资格证书, 5-其他',
    `degree_level` TINYINT DEFAULT NULL COMMENT '学位等级: 1-学士, 2-硕士, 3-博士',
    `description` TEXT COMMENT '类型描述',
    `is_active` TINYINT(1) DEFAULT 1 COMMENT '是否启用',
    `sort_order` INT DEFAULT 0 COMMENT '排序序号',
    `extra_data` JSON DEFAULT NULL COMMENT '扩展数据',
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
    `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_certificate_types_code` (`code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='证书类型表';

-- 3.2 证书模板表 (certificate_templates)
DROP TABLE IF EXISTS `certificate_templates`;
CREATE TABLE `certificate_templates` (
    `id` VARCHAR(36) NOT NULL COMMENT '模板ID (UUID)',
    `name` VARCHAR(128) NOT NULL COMMENT '模板名称',
    `code` VARCHAR(64) NOT NULL COMMENT '模板代码 (唯一标识)',
    `type_id` VARCHAR(36) NOT NULL COMMENT '证书类型ID',
    `thumbnail_url` VARCHAR(512) DEFAULT NULL COMMENT '缩略图URL',
    `template_file_url` VARCHAR(512) DEFAULT NULL COMMENT '模板文件URL',
    `template_data` JSON DEFAULT NULL COMMENT '模板配置数据',
    `fields` JSON NOT NULL COMMENT '模板字段配置',
    `is_default` TINYINT(1) DEFAULT 0 COMMENT '是否默认模板',
    `is_active` TINYINT(1) DEFAULT 1 COMMENT '是否启用',
    `version` INT DEFAULT 1 COMMENT '版本号',
    `organization_id` VARCHAR(36) DEFAULT NULL COMMENT '所属组织ID',
    `created_by` VARCHAR(36) DEFAULT NULL COMMENT '创建人ID',
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
    `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_certificate_templates_code` (`code`),
    KEY `idx_certificate_templates_type` (`type_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='证书模板表';

-- 3.3 证书主表 (certificates)
-- 优化: 字段从70+减少到55个, 移除冗余字段
DROP TABLE IF EXISTS `certificates`;
CREATE TABLE `certificates` (
    `id` VARCHAR(36) NOT NULL COMMENT '证书ID (UUID)',
    `certificate_no` VARCHAR(64) NOT NULL COMMENT '证书编号 (唯一)',
    `type_id` VARCHAR(36) NOT NULL COMMENT '证书类型ID',
    `user_id` VARCHAR(36) NOT NULL COMMENT '持有人用户ID',
    `organization_id` VARCHAR(36) NOT NULL COMMENT '发证机构ID',
    `template_id` VARCHAR(36) DEFAULT NULL COMMENT '使用的模板ID',
    
    -- 基本信息
    `student_no` VARCHAR(64) DEFAULT NULL COMMENT '学号',
    `name` VARCHAR(64) NOT NULL COMMENT '证书持有人姓名',
    `id_card_number` VARCHAR(18) DEFAULT NULL COMMENT '身份证号',
    `gender` TINYINT DEFAULT NULL COMMENT '性别: 1-男, 2-女',
    
    -- 学历信息
    `major` VARCHAR(128) DEFAULT NULL COMMENT '专业名称',
    `major_code` VARCHAR(32) DEFAULT NULL COMMENT '专业代码',
    `degree` TINYINT DEFAULT 1 COMMENT '学位等级: 1-学士, 2-硕士, 3-博士, 4-无学位, 5-专科',
    `degree_name` VARCHAR(32) DEFAULT NULL COMMENT '学位名称',
    `education_level` TINYINT DEFAULT NULL COMMENT '学历层次',
    `education_type` TINYINT DEFAULT 1 COMMENT '办学类型: 1-全日制, 2-非全日制, 3-成人教育',
    `enrollment_date` DATE DEFAULT NULL COMMENT '入学日期',
    `graduation_date` DATE DEFAULT NULL COMMENT '毕业日期',
    `study_period` VARCHAR(64) DEFAULT NULL COMMENT '学习期限',
    `school_system` VARCHAR(32) DEFAULT NULL COMMENT '学制',
    `campus` VARCHAR(128) DEFAULT NULL COMMENT '校区',
    
    -- 成绩信息
    `gpa` DECIMAL(4,2) DEFAULT NULL COMMENT '平均学分绩点',
    `total_credits` DECIMAL(6,2) DEFAULT NULL COMMENT '总学分',
    
    -- 证书状态
    `status` TINYINT NOT NULL DEFAULT 1 COMMENT '证书状态: 1-有效, 2-已撤销, 3-已挂失, 4-草稿, 5-审核中',
    `issue_date` DATE NOT NULL COMMENT '发证日期',
    `valid_from_date` DATE DEFAULT NULL COMMENT '有效开始日期',
    `valid_until_date` DATE DEFAULT NULL COMMENT '有效结束日期',
    
    -- 区块链信息
    `blockchain_tx_hash` VARCHAR(128) DEFAULT NULL COMMENT '区块链交易哈希',
    `blockchain_cert_hash` VARCHAR(64) DEFAULT NULL COMMENT '证书数据哈希',
    `blockchain_block_no` BIGINT DEFAULT NULL COMMENT '区块链块号',
    `blockchain_timestamp` DATETIME(3) DEFAULT NULL COMMENT '区块链时间戳',
    `on_chain_status` TINYINT DEFAULT 1 COMMENT '上链状态: 1-待上链, 2-上链中, 3-已上链, 4-上链失败',
    `on_chain_at` DATETIME(3) DEFAULT NULL COMMENT '上链时间',
    
    -- PDF信息
    `pdf_url` VARCHAR(512) DEFAULT NULL COMMENT 'PDF文件URL',
    `pdf_hash` VARCHAR(64) DEFAULT NULL COMMENT 'PDF文件哈希',
    `pdf_signed_by` VARCHAR(36) DEFAULT NULL COMMENT 'PDF签名人ID',
    `pdf_signed_at` DATETIME(3) DEFAULT NULL COMMENT 'PDF签名时间',
    
    -- 签发/撤销信息
    `signed_by` VARCHAR(36) DEFAULT NULL COMMENT '证书签发人ID',
    `issued_by` VARCHAR(36) DEFAULT NULL COMMENT '操作人ID',
    `issue_reason` VARCHAR(256) DEFAULT NULL COMMENT '颁发原因',
    `revoked_at` DATETIME(3) DEFAULT NULL COMMENT '撤销时间',
    `revoked_by` VARCHAR(36) DEFAULT NULL COMMENT '撤销人ID',
    `revoke_reason` TEXT DEFAULT NULL COMMENT '撤销原因',
    
    -- 验证统计
    `verification_count` INT DEFAULT 0 COMMENT '验证次数',
    `last_verified_at` DATETIME(3) DEFAULT NULL COMMENT '最后验证时间',
    
    -- 扩展数据
    `extra_data` JSON DEFAULT NULL COMMENT '扩展数据',
    
    -- 审计字段
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
    `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
    `deleted_at` DATETIME(3) DEFAULT NULL COMMENT '删除时间',
    `version` INT DEFAULT 1 COMMENT '版本号',
    
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_certificates_no` (`certificate_no`),
    KEY `idx_certificates_type` (`type_id`),
    KEY `idx_certificates_user` (`user_id`),
    KEY `idx_certificates_org` (`organization_id`),
    KEY `idx_certificates_status` (`status`),
    KEY `idx_certificates_issue_date` (`issue_date`),
    KEY `idx_certificates_grad_date` (`graduation_date`),
    KEY `idx_certificates_blockchain_hash` (`blockchain_tx_hash`(64)),
    KEY `idx_certificates_name` (`name`(32)),
    KEY `idx_certificates_student_no` (`student_no`),
    KEY `idx_certificates_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='证书主表';

-- 3.4 证书批次表 (certificate_batches)
DROP TABLE IF EXISTS `certificate_batches`;
CREATE TABLE `certificate_batches` (
    `id` VARCHAR(36) NOT NULL COMMENT '批次ID (UUID)',
    `batch_no` VARCHAR(64) NOT NULL COMMENT '批次号 (唯一)',
    `type_id` VARCHAR(36) NOT NULL COMMENT '证书类型ID',
    `organization_id` VARCHAR(36) NOT NULL COMMENT '发证机构ID',
    `template_id` VARCHAR(36) DEFAULT NULL COMMENT '模板ID',
    `name` VARCHAR(128) NOT NULL COMMENT '批次名称',
    `description` TEXT COMMENT '批次描述',
    `total_count` INT DEFAULT 0 COMMENT '总数量',
    `success_count` INT DEFAULT 0 COMMENT '成功数量',
    `fail_count` INT DEFAULT 0 COMMENT '失败数量',
    `status` TINYINT NOT NULL DEFAULT 1 COMMENT '批次状态: 1-待处理, 2-处理中, 3-已完成, 4-已暂停, 5-已取消',
    `import_file_url` VARCHAR(512) DEFAULT NULL COMMENT '导入文件URL',
    `executed_by` VARCHAR(36) DEFAULT NULL COMMENT '执行人ID',
    `started_at` DATETIME(3) DEFAULT NULL COMMENT '开始时间',
    `completed_at` DATETIME(3) DEFAULT NULL COMMENT '完成时间',
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
    `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_certificate_batches_no` (`batch_no`),
    KEY `idx_certificate_batches_type` (`type_id`),
    KEY `idx_certificate_batches_org` (`organization_id`),
    KEY `idx_certificate_batches_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='证书批次表';

SET FOREIGN_KEY_CHECKS = 1;
