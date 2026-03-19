-- =====================================================
-- EduChain 验证服务模块 (优化后)
-- MySQL 8.0
-- 设计原则: 不使用外键, 数据表尽量精简
-- =====================================================

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- 5.1 验证记录表 (verifications)
-- 优化: 合并verification_requests和verification_reports到本表
DROP TABLE IF EXISTS `verifications`;
CREATE TABLE `verifications` (
    `id` VARCHAR(36) NOT NULL COMMENT '验证记录ID (UUID)',
    `verification_no` VARCHAR(64) NOT NULL COMMENT '验证流水号 (唯一)',
    `certificate_id` VARCHAR(36) NOT NULL COMMENT '证书ID',
    `user_id` VARCHAR(36) DEFAULT NULL COMMENT '验证人用户ID',
    `verifier_id` VARCHAR(36) DEFAULT NULL COMMENT '验证操作人ID',
    `verifier_org_id` VARCHAR(36) DEFAULT NULL COMMENT '验证人所属组织ID',
    `verifier_org_name` VARCHAR(128) DEFAULT NULL COMMENT '验证人所属组织名称',
    
    -- 验证类型与用途
    `verification_type` TINYINT NOT NULL DEFAULT 1 COMMENT '验证类型: 1-本人查询, 2-机构验证, 3-管理员查询, 4-批量验证',
    `purpose` VARCHAR(256) DEFAULT NULL COMMENT '验证用途',
    `input_type` TINYINT DEFAULT 1 COMMENT '输入类型: 1-证书编号, 2-身份证号, 3-姓名+身份证, 4-扫码验证',
    `input_data` JSON DEFAULT NULL COMMENT '输入数据',
    
    -- 验证结果
    `result` TINYINT NOT NULL COMMENT '验证结果: 1-真实, 2-可疑, 3-未匹配, 4-已撤销',
    `result_details` JSON DEFAULT NULL COMMENT '验证结果详情',
    `matched_fields` JSON DEFAULT NULL COMMENT '匹配的字段列表',
    `mismatch_fields` JSON DEFAULT NULL COMMENT '不匹配的字段列表',
    
    -- 区块链验证
    `blockchain_verified` TINYINT(1) DEFAULT 0 COMMENT '区块链是否已验证',
    `blockchain_result` JSON DEFAULT NULL COMMENT '区块链验证结果',
    
    -- 报告信息
    `report_url` VARCHAR(512) DEFAULT NULL COMMENT '验证报告URL',
    `report_id` VARCHAR(36) DEFAULT NULL COMMENT '报告ID',
    
    -- 风险评估
    `risk_level` TINYINT DEFAULT 1 COMMENT '风险等级: 1-低, 2-中, 3-高, 4-极高, 5-已确认欺诈',
    `risk_factors` JSON DEFAULT NULL COMMENT '风险因素列表',
    
    -- 审核信息
    `reviewed_by` VARCHAR(36) DEFAULT NULL COMMENT '审核人ID',
    `reviewed_at` DATETIME(3) DEFAULT NULL COMMENT '审核时间',
    `review_result` TINYINT DEFAULT NULL COMMENT '审核结果',
    `review_notes` TEXT COMMENT '审核备注',
    
    -- 请求信息
    `ip_address` VARCHAR(45) DEFAULT NULL COMMENT '请求IP地址',
    `user_agent` TEXT COMMENT 'User-Agent',
    `location` VARCHAR(128) DEFAULT NULL COMMENT '地理位置',
    
    -- 状态
    `status` TINYINT DEFAULT 1 COMMENT '状态: 1-待验证, 2-验证中, 3-已完成, 4-已过期',
    `expires_at` DATETIME(3) DEFAULT NULL COMMENT '过期时间',
    `verified_at` DATETIME(3) DEFAULT NULL COMMENT '验证完成时间',
    
    -- 审计字段
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
    `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
    
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_verifications_no` (`verification_no`),
    KEY `idx_verifications_cert` (`certificate_id`),
    KEY `idx_verifications_user` (`user_id`),
    KEY `idx_verifications_verifier` (`verifier_id`),
    KEY `idx_verifications_type` (`verification_type`),
    KEY `idx_verifications_result` (`result`),
    KEY `idx_verifications_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='验证记录表';

SET FOREIGN_KEY_CHECKS = 1;
