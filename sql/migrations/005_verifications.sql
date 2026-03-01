-- =====================================================
-- EduChain 验证服务模块数据库表
-- MySQL 8.0
-- =====================================================

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- =====================================================
-- 5.1 验证记录表 (verifications)
-- =====================================================
DROP TABLE IF EXISTS `verifications`;
CREATE TABLE `verifications` (
    `id` VARCHAR(36) NOT NULL COMMENT '验证记录ID (UUID)',
    `verification_no` VARCHAR(64) NOT NULL COMMENT '验证流水号 (唯一)',
    `certificate_id` VARCHAR(36) NOT NULL COMMENT '证书ID',
    `user_id` VARCHAR(36) DEFAULT NULL COMMENT '验证人用户ID',
    `verifier_id` VARCHAR(36) DEFAULT NULL COMMENT '验证操作人ID',
    `verifier_name` VARCHAR(64) DEFAULT NULL COMMENT '验证人姓名',
    `verifier_org_id` VARCHAR(36) DEFAULT NULL COMMENT '验证人所属组织ID',
    `verifier_org_name` VARCHAR(128) DEFAULT NULL COMMENT '验证人所属组织名称',
    `verification_type` TINYINT NOT NULL DEFAULT 1 COMMENT '验证类型: 1-本人查询, 2-机构验证, 3-管理员查询, 4-批量验证',
    `purpose` VARCHAR(256) DEFAULT NULL COMMENT '验证用途',
    `request_data` JSON DEFAULT NULL COMMENT '请求数据 (JSON)',
    
    -- 验证输入
    `input_type` TINYINT DEFAULT 1 COMMENT '输入类型: 1-证书编号, 2-身份证号, 3-姓名+身份证, 4-扫码验证',
    `input_data` JSON DEFAULT NULL COMMENT '输入数据 (JSON)',
    
    -- 验证结果
    `result` TINYINT NOT NULL COMMENT '验证结果: 1-真实, 2-可疑, 3-未匹配, 4-已撤销',
    `result_details` JSON DEFAULT NULL COMMENT '验证结果详情 (JSON)',
    `matched_fields` JSON DEFAULT NULL COMMENT '匹配的字段列表',
    `mismatch_fields` JSON DEFAULT NULL COMMENT '不匹配的字段列表',
    
    -- 区块链验证
    `blockchain_verified` TINYINT(1) DEFAULT 0 COMMENT '区块链是否已验证: 0-否, 1-是',
    `blockchain_result` JSON DEFAULT NULL COMMENT '区块链验证结果 (JSON)',
    
    -- 报告信息
    `report_url` VARCHAR(512) DEFAULT NULL COMMENT '验证报告URL',
    `report_id` VARCHAR(36) DEFAULT NULL COMMENT '报告ID',
    `report_generated_at` DATETIME(3) DEFAULT NULL COMMENT '报告生成时间',
    
    -- 风险评估
    `risk_level` TINYINT DEFAULT 1 COMMENT '风险等级: 1-低, 2-中, 3-高, 4-极高, 5-已确认欺诈',
    `risk_factors` JSON DEFAULT NULL COMMENT '风险因素列表',
    
    -- 审核信息
    `reviewed_by` VARCHAR(36) DEFAULT NULL COMMENT '审核人ID',
    `reviewed_at` DATETIME(3) DEFAULT NULL COMMENT '审核时间',
    `review_result` TINYINT DEFAULT NULL COMMENT '审核结果',
    `review_notes` TEXT DEFAULT NULL COMMENT '审核备注',
    
    -- 请求信息
    `ip_address` VARCHAR(45) DEFAULT NULL COMMENT '请求IP地址',
    `user_agent` TEXT DEFAULT NULL COMMENT 'User-Agent',
    `location` VARCHAR(128) DEFAULT NULL COMMENT '地理位置',
    `device_info` JSON DEFAULT NULL COMMENT '设备信息 (JSON)',
    
    -- 状态
    `status` TINYINT DEFAULT 1 COMMENT '状态: 1-待验证, 2-验证中, 3-已完成, 4-已过期',
    `expires_at` DATETIME(3) DEFAULT NULL COMMENT '过期时间',
    
    -- 时间
    `verified_at` DATETIME(3) DEFAULT NULL COMMENT '验证完成时间',
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
    `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_verifications_no` (`verification_no`),
    KEY `idx_verifications_cert` (`certificate_id`),
    KEY `idx_verifications_user` (`user_id`),
    KEY `idx_verifications_verifier` (`verifier_id`),
    KEY `idx_verifications_type` (`verification_type`),
    KEY `idx_verifications_result` (`result`),
    KEY `idx_verifications_risk` (`risk_level`),
    KEY `idx_verifications_ip` (`ip_address`(45)),
    KEY `idx_verifications_created_at` (`created_at`),
    KEY `idx_verifications_cert_result` (`certificate_id`, `result`),
    KEY `idx_verifications_cert_created` (`certificate_id`, `created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='验证记录表';

-- =====================================================
-- 5.2 验证请求日志表 (verification_requests)
-- =====================================================
DROP TABLE IF EXISTS `verification_requests`;
CREATE TABLE `verification_requests` (
    `id` VARCHAR(36) NOT NULL COMMENT '请求ID (UUID)',
    `verification_id` VARCHAR(36) NOT NULL COMMENT '验证记录ID',
    `input_no` VARCHAR(64) DEFAULT NULL COMMENT '输入编号',
    `input_type` TINYINT NOT NULL COMMENT '输入类型',
    `input_data` JSON DEFAULT NULL COMMENT '输入数据 (JSON)',
    `processed` TINYINT(1) DEFAULT 0 COMMENT '是否已处理: 0-否, 1-是',
    `processed_at` DATETIME(3) DEFAULT NULL COMMENT '处理时间',
    `processing_time_ms` INT DEFAULT NULL COMMENT '处理耗时 (毫秒)',
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
    PRIMARY KEY (`id`),
    KEY `idx_verification_requests_verify` (`verification_id`),
    KEY `idx_verification_requests_no` (`input_no`),
    KEY `idx_verification_requests_processed` (`processed`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='验证请求日志表';

-- =====================================================
-- 5.3 验证报告表 (verification_reports)
-- =====================================================
DROP TABLE IF EXISTS `verification_reports`;
CREATE TABLE `verification_reports` (
    `id` VARCHAR(36) NOT NULL COMMENT '报告ID (UUID)',
    `report_no` VARCHAR(64) NOT NULL COMMENT '报告编号 (唯一)',
    `verification_id` VARCHAR(36) NOT NULL COMMENT '验证记录ID',
    `certificate_id` VARCHAR(36) NOT NULL COMMENT '证书ID',
    `report_type` TINYINT DEFAULT 1 COMMENT '报告类型: 1-验证报告, 2-详细报告, 3-原始数据',
    `report_data` JSON NOT NULL COMMENT '报告数据 (JSON)',
    `pdf_url` VARCHAR(512) DEFAULT NULL COMMENT 'PDF报告URL',
    `pdf_generated_at` DATETIME(3) DEFAULT NULL COMMENT 'PDF生成时间',
    `is_authorized` TINYINT(1) DEFAULT 0 COMMENT '是否已授权: 0-否, 1-是',
    `authorized_by` VARCHAR(36) DEFAULT NULL COMMENT '授权人ID',
    `authorized_at` DATETIME(3) DEFAULT NULL COMMENT '授权时间',
    `access_count` INT DEFAULT 0 COMMENT '访问次数',
    `last_accessed_at` DATETIME(3) DEFAULT NULL COMMENT '最后访问时间',
    `expires_at` DATETIME(3) DEFAULT NULL COMMENT '过期时间',
    `status` TINYINT DEFAULT 1 COMMENT '状态: 1-有效, 2-已过期, 3-已撤销',
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
    `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_verification_reports_no` (`report_no`),
    KEY `idx_verification_reports_verify` (`verification_id`),
    KEY `idx_verification_reports_cert` (`certificate_id`),
    KEY `idx_verification_reports_expires` (`expires_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='验证报告表';

SET FOREIGN_KEY_CHECKS = 1;