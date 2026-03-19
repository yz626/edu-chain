-- =====================================================
-- EduChain 审计日志模块 (优化后)
-- MySQL 8.0
-- 设计原则: 不使用外键, 数据表尽量精简
-- =====================================================

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- 6.1 审计日志表 (audit_logs)
-- 优化: 合并login_logs, operation_logs, api_access_logs到本表
DROP TABLE IF EXISTS `audit_logs`;
CREATE TABLE `audit_logs` (
    `id` VARCHAR(36) NOT NULL COMMENT '日志ID (UUID)',
    `user_id` VARCHAR(36) DEFAULT NULL COMMENT '操作用户ID',
    `username` VARCHAR(64) DEFAULT NULL COMMENT '用户名',
    `organization_id` VARCHAR(36) DEFAULT NULL COMMENT '所属组织ID',
    `module` VARCHAR(64) NOT NULL COMMENT '模块名称: login, certificate, verification, system, api',
    `action` VARCHAR(64) NOT NULL COMMENT '操作类型',
    `resource_type` VARCHAR(64) DEFAULT NULL COMMENT '资源类型',
    `resource_id` VARCHAR(36) DEFAULT NULL COMMENT '资源ID',
    `resource_name` VARCHAR(256) DEFAULT NULL COMMENT '资源名称',
    `description` TEXT COMMENT '操作描述',
    `request_data` JSON DEFAULT NULL COMMENT '请求数据',
    `response_data` JSON DEFAULT NULL COMMENT '响应数据',
    `ip_address` VARCHAR(45) DEFAULT NULL COMMENT 'IP地址',
    `user_agent` TEXT COMMENT 'User-Agent',
    `location` VARCHAR(128) DEFAULT NULL COMMENT '地理位置',
    `device_info` JSON DEFAULT NULL COMMENT '设备信息',
    `login_success` TINYINT(1) DEFAULT 1 COMMENT '是否成功',
    `error_message` TEXT COMMENT '错误信息',
    `trace_id` VARCHAR(36) DEFAULT NULL COMMENT '链路追踪ID',
    `duration_ms` INT DEFAULT NULL COMMENT '耗时 (毫秒)',
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
    
    PRIMARY KEY (`id`),
    KEY `idx_audit_logs_user` (`user_id`),
    KEY `idx_audit_logs_module` (`module`),
    KEY `idx_audit_logs_action` (`action`),
    KEY `idx_audit_logs_resource` (`resource_type`(32), `resource_id`),
    KEY `idx_audit_logs_created_at` (`created_at`),
    KEY `idx_audit_logs_trace` (`trace_id`),
    KEY `idx_audit_logs_user_module_at` (`user_id`, `module`, `created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='审计日志表';

SET FOREIGN_KEY_CHECKS = 1;
