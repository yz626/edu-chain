-- =====================================================
-- EduChain 区块链模块数据库表
-- MySQL 8.0
-- =====================================================

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- =====================================================
-- 4.1 区块链网络表 (blockchain_networks)
-- =====================================================
DROP TABLE IF EXISTS `blockchain_networks`;
CREATE TABLE `blockchain_networks` (
    `id` VARCHAR(36) NOT NULL COMMENT '网络ID (UUID)',
    `name` VARCHAR(64) NOT NULL COMMENT '网络名称',
    `code` VARCHAR(32) NOT NULL COMMENT '网络代码 (唯一标识)',
    `type` TINYINT NOT NULL DEFAULT 1 COMMENT '区块链类型: 1-Fabric, 2-Ethereum',
    `chain_id` INT DEFAULT NULL COMMENT '链ID',
    `endpoint_url` VARCHAR(256) NOT NULL COMMENT '节点RPC endpoint URL',
    `explorer_url` VARCHAR(256) DEFAULT NULL COMMENT '区块链浏览器URL',
    `status` TINYINT DEFAULT 1 COMMENT '状态: 1-正常, 2-维护中, 3-停用',
    `is_default` TINYINT(1) DEFAULT 0 COMMENT '是否默认网络: 0-否, 1-是',
    `extra_data` JSON DEFAULT NULL COMMENT '扩展数据 (JSON格式)',
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
    `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_blockchain_networks_code` (`code`),
    KEY `idx_blockchain_networks_type` (`type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='区块链网络表';

-- =====================================================
-- 4.2 区块链交易表 (blockchain_transactions)
-- =====================================================
DROP TABLE IF EXISTS `blockchain_transactions`;
CREATE TABLE `blockchain_transactions` (
    `id` VARCHAR(36) NOT NULL COMMENT '交易ID (UUID)',
    `tx_hash` VARCHAR(128) NOT NULL COMMENT '交易哈希 (唯一)',
    `network_id` VARCHAR(36) NOT NULL COMMENT '区块链网络ID',
    `certificate_id` VARCHAR(36) DEFAULT NULL COMMENT '关联证书ID',
    `tx_type` TINYINT NOT NULL DEFAULT 1 COMMENT '交易类型: 1-存证, 2-撤销, 3-查询, 4-转让, 5-其他',
    `from_address` VARCHAR(128) DEFAULT NULL COMMENT '发起方地址',
    `to_address` VARCHAR(128) DEFAULT NULL COMMENT '接收方地址',
    `data` TEXT DEFAULT NULL COMMENT '交易数据',
    `value` DECIMAL(38,0) DEFAULT NULL COMMENT '交易金额',
    `gas_used` BIGINT DEFAULT NULL COMMENT 'Gas消耗',
    `gas_price` DECIMAL(38,0) DEFAULT NULL COMMENT 'Gas价格',
    `tx_fee` DECIMAL(38,8) DEFAULT NULL COMMENT '交易手续费',
    `block_number` BIGINT DEFAULT NULL COMMENT '区块高度',
    `block_hash` VARCHAR(128) DEFAULT NULL COMMENT '区块哈希',
    `block_timestamp` DATETIME(3) DEFAULT NULL COMMENT '区块时间戳',
    `confirmations` INT DEFAULT 0 COMMENT '确认数',
    `status` TINYINT NOT NULL DEFAULT 1 COMMENT '交易状态: 1-待处理, 2-处理中, 3-成功, 4-失败, 5-超时',
    `error_message` TEXT DEFAULT NULL COMMENT '错误信息',
    `retry_count` TINYINT DEFAULT 0 COMMENT '重试次数',
    `last_retry_at` DATETIME(3) DEFAULT NULL COMMENT '最后重试时间',
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
    `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_blockchain_transactions_hash` (`tx_hash`),
    KEY `idx_blockchain_transactions_network` (`network_id`),
    KEY `idx_blockchain_transactions_cert` (`certificate_id`),
    KEY `idx_blockchain_transactions_status` (`status`),
    KEY `idx_blockchain_transactions_block` (`block_number`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='区块链交易表';

-- =====================================================
-- 4.3 区块链证书存证表 (blockchain_certificates)
-- =====================================================
DROP TABLE IF EXISTS `blockchain_certificates`;
CREATE TABLE `blockchain_certificates` (
    `id` VARCHAR(36) NOT NULL COMMENT '记录ID (UUID)',
    `certificate_id` VARCHAR(36) NOT NULL COMMENT '证书ID (唯一)',
    `network_id` VARCHAR(36) NOT NULL COMMENT '区块链网络ID',
    `tx_id` VARCHAR(36) NOT NULL COMMENT '区块链交易ID',
    `cert_hash` VARCHAR(64) NOT NULL COMMENT '证书数据哈希',
    `owner_address` VARCHAR(128) DEFAULT NULL COMMENT '所有者区块链地址',
    `token_id` VARCHAR(64) DEFAULT NULL COMMENT 'NFT Token ID (如果是NFT)',
    `uri` TEXT DEFAULT NULL COMMENT '元数据URI',
    `data` TEXT DEFAULT NULL COMMENT '链上存储的证书数据',
    `is_synced` TINYINT(1) DEFAULT 0 COMMENT '是否已同步: 0-未同步, 1-已同步',
    `synced_at` DATETIME(3) DEFAULT NULL COMMENT '同步时间',
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
    `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_blockchain_certificates_cert` (`certificate_id`),
    KEY `idx_blockchain_certificates_network` (`network_id`),
    KEY `idx_blockchain_certificates_hash` (`cert_hash`),
    KEY `idx_blockchain_certificates_token` (`token_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='区块链证书存证表';

-- =====================================================
-- 4.4 智能合约表 (smart_contracts)
-- =====================================================
DROP TABLE IF EXISTS `smart_contracts`;
CREATE TABLE `smart_contracts` (
    `id` VARCHAR(36) NOT NULL COMMENT '合约ID (UUID)',
    `name` VARCHAR(64) NOT NULL COMMENT '合约名称',
    `code` VARCHAR(32) NOT NULL COMMENT '合约代码 (唯一标识)',
    `network_id` VARCHAR(36) NOT NULL COMMENT '部署网络ID',
    `contract_address` VARCHAR(128) NOT NULL COMMENT '合约地址',
    `abi` JSON DEFAULT NULL COMMENT 'ABI接口定义',
    `bytecode` TEXT DEFAULT NULL COMMENT '字节码',
    `version` VARCHAR(32) DEFAULT NULL COMMENT '合约版本',
    `deployer_address` VARCHAR(128) DEFAULT NULL COMMENT '部署者地址',
    `deployed_at` DATETIME(3) DEFAULT NULL COMMENT '部署时间',
    `status` TINYINT DEFAULT 1 COMMENT '状态: 1-正常, 2-停用, 3-已废弃',
    `is_verified` TINYINT(1) DEFAULT 0 COMMENT '是否已验证源码: 0-否, 1-是',
    `verified_at` DATETIME(3) DEFAULT NULL COMMENT '验证时间',
    `extra_data` JSON DEFAULT NULL COMMENT '扩展数据 (JSON格式)',
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
    `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_smart_contracts_code` (`code`),
    KEY `idx_smart_contracts_network` (`network_id`),
    KEY `idx_smart_contracts_address` (`contract_address`(64))
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='智能合约表';

SET FOREIGN_KEY_CHECKS = 1;