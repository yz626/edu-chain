-- =====================================================
-- EduChain 系统管理模块 (优化后)
-- MySQL 8.0
-- 设计原则: 不使用外键, 数据表尽量精简
-- =====================================================

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- 7.1 系统配置表 (system_configs)
DROP TABLE IF EXISTS `system_configs`;
CREATE TABLE `system_configs` (
    `id` VARCHAR(36) NOT NULL COMMENT '配置ID (UUID)',
    `key` VARCHAR(128) NOT NULL COMMENT '配置键 (唯一)',
    `value` TEXT NOT NULL COMMENT '配置值',
    `type` VARCHAR(32) NOT NULL DEFAULT 'string' COMMENT '配置类型: string, integer, boolean, json',
    `description` TEXT COMMENT '配置描述',
    `category` VARCHAR(64) DEFAULT 'general' COMMENT '配置分类',
    `is_encrypted` TINYINT(1) DEFAULT 0 COMMENT '是否加密存储',
    `is_editable` TINYINT(1) DEFAULT 1 COMMENT '是否可编辑',
    `sort_order` INT DEFAULT 0 COMMENT '排序序号',
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
    `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_system_configs_key` (`key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='系统配置表';

-- 7.2 字典表 (dictionaries)
DROP TABLE IF EXISTS `dictionaries`;
CREATE TABLE `dictionaries` (
    `id` VARCHAR(36) NOT NULL COMMENT '字典ID (UUID)',
    `type_code` VARCHAR(64) NOT NULL COMMENT '字典类型代码',
    `code` VARCHAR(64) NOT NULL COMMENT '字典项代码',
    `name` VARCHAR(128) NOT NULL COMMENT '字典项名称',
    `value` TEXT DEFAULT NULL COMMENT '字典项值',
    `parent_code` VARCHAR(64) DEFAULT NULL COMMENT '父级代码',
    `description` TEXT COMMENT '描述',
    `sort_order` INT DEFAULT 0 COMMENT '排序序号',
    `is_active` TINYINT(1) DEFAULT 1 COMMENT '是否启用',
    `extra_data` JSON DEFAULT NULL COMMENT '扩展数据',
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
    `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_dictionaries_type_code` (`type_code`, `code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='字典表';

-- 7.3 文件记录表 (file_records)
DROP TABLE IF EXISTS `file_records`;
CREATE TABLE `file_records` (
    `id` VARCHAR(36) NOT NULL COMMENT '文件ID (UUID)',
    `file_name` VARCHAR(256) NOT NULL COMMENT '文件名',
    `file_path` VARCHAR(512) NOT NULL COMMENT '文件路径',
    `file_url` VARCHAR(512) NOT NULL COMMENT '文件访问URL',
    `file_type` VARCHAR(32) NOT NULL COMMENT '文件类型: pdf, image, excel, word, zip',
    `mime_type` VARCHAR(64) DEFAULT NULL COMMENT 'MIME类型',
    `file_size` BIGINT NOT NULL COMMENT '文件大小 (字节)',
    `file_hash` VARCHAR(64) DEFAULT NULL COMMENT '文件哈希值 (SHA256)',
    `storage_type` TINYINT DEFAULT 1 COMMENT '存储类型: 1-本地, 2-OSS, 3-S3, 4-MinIO',
    `storage_bucket` VARCHAR(128) DEFAULT NULL COMMENT '存储桶名称',
    `storage_key` VARCHAR(256) DEFAULT NULL COMMENT '存储键',
    `uploaded_by` VARCHAR(36) DEFAULT NULL COMMENT '上传人ID',
    `organization_id` VARCHAR(36) DEFAULT NULL COMMENT '组织ID',
    `related_type` VARCHAR(64) DEFAULT NULL COMMENT '关联业务类型',
    `related_id` VARCHAR(36) DEFAULT NULL COMMENT '关联业务ID',
    `description` VARCHAR(512) DEFAULT NULL COMMENT '文件描述',
    `is_active` TINYINT(1) DEFAULT 1 COMMENT '是否有效',
    `is_temp` TINYINT(1) DEFAULT 0 COMMENT '是否临时文件',
    `expires_at` DATETIME(3) DEFAULT NULL COMMENT '过期时间',
    `access_count` INT DEFAULT 0 COMMENT '访问次数',
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
    `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY `idx_file_records_type` (`file_type`),
    KEY `idx_file_records_uploaded_by` (`uploaded_by`),
    KEY `idx_file_records_related` (`related_type`(32), `related_id`),
    KEY `idx_file_records_hash` (`file_hash`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='文件记录表';

-- 7.4 任务队列表 (job_queues)
DROP TABLE IF EXISTS `job_queues`;
CREATE TABLE `job_queues` (
    `id` VARCHAR(36) NOT NULL COMMENT '任务ID (UUID)',
    `job_type` VARCHAR(64) NOT NULL COMMENT '任务类型',
    `job_name` VARCHAR(128) NOT NULL COMMENT '任务名称',
    `payload` JSON NOT NULL COMMENT '任务数据 (JSON)',
    `priority` TINYINT DEFAULT 5 COMMENT '优先级: 1-最高, 10-最低',
    `status` TINYINT NOT NULL DEFAULT 1 COMMENT '任务状态: 1-等待, 2-执行中, 3-成功, 4-失败, 5-已取消, 6-已重试',
    `retry_count` TINYINT DEFAULT 0 COMMENT '已重试次数',
    `max_retry_count` TINYINT DEFAULT 3 COMMENT '最大重试次数',
    `next_retry_at` DATETIME(3) DEFAULT NULL COMMENT '下次重试时间',
    `started_at` DATETIME(3) DEFAULT NULL COMMENT '开始执行时间',
    `completed_at` DATETIME(3) DEFAULT NULL COMMENT '完成时间',
    `progress_percent` TINYINT DEFAULT 0 COMMENT '进度百分比 (0-100)',
    `result_data` JSON DEFAULT NULL COMMENT '执行结果',
    `error_message` TEXT COMMENT '错误信息',
    `scheduled_at` DATETIME(3) DEFAULT NULL COMMENT '计划执行时间',
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
    `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY `idx_job_queues_type` (`job_type`),
    KEY `idx_job_queues_status` (`status`),
    KEY `idx_job_queues_priority` (`priority`, `scheduled_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='任务队列表';

SET FOREIGN_KEY_CHECKS = 1;
