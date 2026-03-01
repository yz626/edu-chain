# EduChain 数据库迁移脚本

## 概述

本目录包含 EduChain 区块链学历证书管理系统的数据库表结构定义，使用 MySQL 8.0 语法。

## 文件说明

| 文件 | 模块 | 表数量 | 说明 |
|------|------|--------|------|
| `001_users.sql` | 用户认证 | 9张 | 用户、角色、权限、登录日志等 |
| `002_organizations.sql` | 组织管理 | 3张 | 组织、部门、成员关联 |
| `003_certificates.sql` | 证书管理 | 5张 | 证书类型、模板、主表、批次 |
| `004_blockchain.sql` | 区块链 | 4张 | 网络、交易、存证、合约 |
| `005_verifications.sql` | 验证服务 | 3张 | 验证记录、请求、报告 |
| `006_audit_logs.sql` | 审计日志 | 3张 | 审计、操作、API日志 |
| `007_system.sql` | 系统管理 | 5张 | 配置、字典、文件、通知、任务 |
| `008_statistics.sql` | 统计分析 | 1张 | 每日统计 |

## 特性

- **无外键约束**：关联关系由应用层维护
- **UUID主键**：使用 UUID 作为主键
- **JSON字段**：使用 JSON 类型存储扩展数据
- **完整索引**：针对常见查询场景优化
- **字段注释**：每个字段都有详细注释
- **乐观锁**：version 字段支持并发控制

## 使用方法

### 1. 创建数据库

```sql
CREATE DATABASE edu_chain DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

### 2. 执行迁移脚本 (按序号顺序)

```bash
# 执行所有脚本
mysql -u root -p edu_chain < 001_users.sql
mysql -u root -p edu_chain < 002_organizations.sql
mysql -u root -p edu_chain < 003_certificates.sql
mysql -u root -p edu_chain < 004_blockchain.sql
mysql -u root -p edu_chain < 005_verifications.sql
mysql -u root -p edu_chain < 006_audit_logs.sql
mysql -u root -p edu_chain < 007_system.sql
mysql -u root -p edu_chain < 008_statistics.sql
```

### 3. 推荐引擎设置

所有表使用 `InnoDB` 引擎，支持事务和外键(如果需要)。

### 4. 字符集设置

```sql
SET NAMES utf8mb4;
SET CHARACTER SET utf8mb4;
```

## 表统计

- **总计**: 33张表
- **用户认证模块**: 9张表
- **组织管理模块**: 3张表
- **证书管理模块**: 5张表
- **区块链模块**: 4张表
- **验证服务模块**: 3张表
- **审计日志模块**: 3张表
- **系统管理模块**: 5张表
- **统计分析模块**: 1张表