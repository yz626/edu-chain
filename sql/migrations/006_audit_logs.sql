-- =====================================================
-- EduChain 审计与日志模块数据库表
-- MySQL 8.0
-- =====================================================

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- =====================================================
-- 6.1 审计日志表 (audit_logs)
-- =====================================================
DROP TABLE IF EXISTS `audit_logs`;
CREATE TABLE `audit_logs` (
    `id` VARCHAR(36) NOT NULL COMMENT '日志ID (UUID)',
    `user_id` VARCHAR(36) DEFAULT NULL COMMENT '操作用户ID',
    `username` VARCHAR(64) DEFAULT NULL COMMENT '用户名',
    `real_name` VARCHAR(64) DEFAULT NULL COMMENT '真实姓名',
    `organization_id` VARCHAR(36) DEFAULT NULL COMMENT '所属组织ID',
    `module` VARCHAR(64) NOT NULL COMMENT '模块名称',
    `action` VARCHAR(64) NOT NULL COMMENT '操作类型',
    `resource_type` VARCHAR(64) DEFAULT NULL COMMENT '资源类型',
    `resource_id` VARCHAR(36) DEFAULT NULL COMMENT '资源ID',
    `resource_name` VARCHAR(256) DEFAULT NULL COMMENT '资源名称',
    `description` TEXT COMMENT '操作描述',
    `request_data` JSON DEFAULT NULL COMMENT '请求数据 (JSON)',
    `response_data` JSON DEFAULT NULL COMMENT '响应数据 (JSON)',
    `ip_address` VARCHAR(45) DEFAULT NULL COMMENT 'IP地址',
    `user_agent` TEXT DEFAULT NULL COMMENT 'User-Agent',
    `location` VARCHAR(128) DEFAULT NULL COMMENT '地理位置',
    `device_info` JSON DEFAULT NULL COMMENT '设备信息 (JSON)',
    `status` TINYINT DEFAULT 1 COMMENT '状态: 1-成功, 2-失败, 3-部分成功',
    `error_message` TEXT DEFAULT NULL COMMENT '错误信息',
    `trace_id` VARCHAR(36) DEFAULT NULL COMMENT '链路追踪ID',
    `span_id` VARCHAR(36) DEFAULT NULL COMMENT 'Span ID',
    `parent_span_id` VARCHAR(36) DEFAULT NULL COMMENT '父Span ID',
    `duration_ms` INT DEFAULT NULL COMMENT '耗时 (毫秒)',
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
    PRIMARY KEY (`id`),
    KEY `idx_audit_logs_user` (`user_id`),
    KEY `idx_audit_logs_module` (`module`),
    KEY `idx_audit_logs_action` (`action`),
    KEY `idx_audit_logs_resource` (`resource_type`(32), `resource_id`),
    KEY `idx_audit_logs_ip` (`ip_address`(45)),
    KEY `idx_audit_logs_created_at` (`created_at`),
    KEY `idx_audit_logs_trace` (`trace_id`),
    KEY `idx_audit_logs_user_module_at` (`user_id`, `module`, `created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='审计日志表';

-- =====================================================
-- 6.2 操作日志表 (operation_logs)
-- =====================================================
DROP TABLE IF EXISTS `operation_logs`;
CREATE TABLE `operation_logs` (
    `id` VARCHAR(36) NOT NULL COMMENT '日志ID (UUID)',
    `operator_id` VARCHAR(36) DEFAULT NULL COMMENT '操作人ID',
    `operator_name` VARCHAR(64) DEFAULT NULL COMMENT '操作人姓名',
    `operator_org_id` VARCHAR(36) DEFAULT NULL COMMENT '操作人组织ID',
    `module` VARCHAR(64) DEFAULT NULL COMMENT '模块',
    `operation` VARCHAR(128) NOT NULL COMMENT '操作名称',
    `target_type` VARCHAR(64) DEFAULT NULL COMMENT '目标类型',
    `target_id` VARCHAR(36) DEFAULT NULL COMMENT '目标ID',
    `target_name` VARCHAR(256) DEFAULT NULL COMMENT '目标名称',
    `operation_type` TINYINT NOT NULL DEFAULT 1 COMMENT '操作类型: 1-新增, 2-修改, 3-删除, 4-查询, 5-其他',
    `before_data` JSON DEFAULT NULL COMMENT '修改前数据 (JSON)',
    `after_data` JSON DEFAULT NULL COMMENT '修改后数据 (JSON)',
    `changed_fields` JSON DEFAULT NULL COMMENT '变更字段列表',
    `ip_address` VARCHAR(45) DEFAULT NULL COMMENT 'IP地址',
    `location` VARCHAR(128) DEFAULT NULL COMMENT '地理位置',
    `remark` VARCHAR(512) DEFAULT NULL COMMENT '备注',
    `status` TINYINT DEFAULT 1 COMMENT '状态: 1-成功, 2-失败',
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
    PRIMARY KEY (`id`),
    KEY `idx_operation_logs_operator` (`operator_id`),
    KEY `idx_operation_logs_module` (`module`),
    KEY `idx_operation_logs_target` (`target_type`(32), `target_id`),
    KEY `idx_operation_logs_created_at` (`created_at`),
    KEY `idx_operation_logs_type_at` (`operation_type`, `created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='操作日志表';

-- =====================================================
-- 6.3 API访问日志表 (api_access_logs)
-- =====================================================
DROP TABLE IF EXISTS `api_access_logs`;
CREATE TABLE `api_access_logs` (
    `id` VARCHAR(36) NOT NULL COMMENT '日志ID (UUID)',
    `request_id` VARCHAR(36) NOT NULL COMMENT '请求ID (唯一)',
    `trace_id` VARCHAR(36) DEFAULT NULL COMMENT '链路追踪ID',
    `span_id` VARCHAR(36) DEFAULT NULL COMMENT 'Span ID',
    `parent_span_id` VARCHAR(36) DEFAULT NULL COMMENT '父Span ID',
    `user_id` VARCHAR(36) DEFAULT NULL COMMENT '用户ID',
    `path` VARCHAR(256) NOT NULL COMMENT 'API路径',
    `method` VARCHAR(16) NOT NULL COMMENT 'HTTP方法: GET, POST, PUT, DELETE等',
    `query_params` JSON DEFAULT NULL COMMENT '查询参数 (JSON)',
    `header_data` JSON DEFAULT NULL COMMENT '请求头 (JSON)',
    `body_data` LONGTEXT DEFAULT NULL COMMENT '请求体',
    `response_status` INT DEFAULT NULL COMMENT '响应状态码',
    `response_body` LONGTEXT DEFAULT NULL COMMENT '响应体',
    `response_size` BIGINT DEFAULT NULL COMMENT '响应大小 (字节)',
    `duration_ms` INT NOT NULL COMMENT '请求耗时 (毫秒)',
    `ip_address` VARCHAR(45) DEFAULT NULL COMMENT 'IP地址',
    `user_agent` TEXT DEFAULT NULL COMMENT 'User-Agent',
    `location` VARCHAR(128) DEFAULT NULL COMMENT '地理位置',
    `device_info` JSON DEFAULT NULL COMMENT '设备信息 (JSON)',
    `error_message` TEXT DEFAULT NULL COMMENT '错误信息',
    `stack_trace` TEXT DEFAULT NULL COMMENT '堆栈信息',
    `is_sampled` TINYINT(1) DEFAULT 1 COMMENT '是否采样: 0-否, 1-是',
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_api_access_logs_request` (`request_id`),
    KEY `idx_api_access_logs_user` (`user_id`),
    KEY `idx_api_access_logs_path` (`path`(64)),
    KEY `idx_api_access_logs_method` (`method`),
    KEY `idx_api_access_logs_status` (`response_status`),
    KEY `idx_api_access_logs_duration` (`duration_ms`),
    KEY `idx_api_access_logs_created_at` (`created_at`),
    KEY `idx_api_access_logs_trace` (`trace_id`),
    KEY `idx_api_access_logs_path_time` (`path`(64), `created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='API访问日志表';

SET FOREIGN_KEY_CHECKS = 1;