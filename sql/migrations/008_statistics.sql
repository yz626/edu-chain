-- =====================================================
-- EduChain 统计与分析模块数据库表
-- MySQL 8.0
-- =====================================================

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- =====================================================
-- 8.1 每日统计表 (daily_statistics)
-- =====================================================
DROP TABLE IF EXISTS `daily_statistics`;
CREATE TABLE `daily_statistics` (
    `id` VARCHAR(36) NOT NULL COMMENT '记录ID (UUID)',
    `stat_date` DATE NOT NULL COMMENT '统计日期 (唯一)',
    `module` VARCHAR(64) NOT NULL COMMENT '统计模块: users, certificates, verifications, blockchain, api',
    
    -- 用户统计
    `new_users` INT DEFAULT 0 COMMENT '新增用户数',
    `active_users` INT DEFAULT 0 COMMENT '活跃用户数',
    `login_count` INT DEFAULT 0 COMMENT '登录次数',
    
    -- 证书统计
    `certificates_issued` INT DEFAULT 0 COMMENT '证书颁发数',
    `certificates_revoked` INT DEFAULT 0 COMMENT '证书撤销数',
    `certificates_total` INT DEFAULT 0 COMMENT '证书总数',
    
    -- 验证统计
    `verifications_count` INT DEFAULT 0 COMMENT '验证请求数',
    `verifications_success` INT DEFAULT 0 COMMENT '验证成功数',
    `verifications_failed` INT DEFAULT 0 COMMENT '验证失败数',
    `verifications_by_type` JSON DEFAULT NULL COMMENT '按类型统计 (JSON)',
    
    -- 区块链统计
    `on_chain_count` INT DEFAULT 0 COMMENT '上链请求数',
    `on_chain_success` INT DEFAULT 0 COMMENT '上链成功数',
    `on_chain_failed` INT DEFAULT 0 COMMENT '上链失败数',
    
    -- API统计
    `api_requests` INT DEFAULT 0 COMMENT 'API请求数',
    `api_errors` INT DEFAULT 0 COMMENT 'API错误数',
    `api_avg_duration` DECIMAL(10,2) DEFAULT NULL COMMENT '平均响应时间 (毫秒)',
    
    -- 扩展数据
    `extra_data` JSON DEFAULT NULL COMMENT '扩展数据 (JSON)',
    
    -- 审计字段
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
    `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_daily_statistics_date_module` (`stat_date`, `module`),
    KEY `idx_daily_statistics_date` (`stat_date`),
    KEY `idx_daily_statistics_module_date` (`module`, `stat_date`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='每日统计表';

SET FOREIGN_KEY_CHECKS = 1;