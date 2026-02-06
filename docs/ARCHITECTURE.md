# EduChain 区块链学历证书管理系统 - 架构设计

## 1. 项目概述

### 1.1 项目背景
基于区块链技术的学历证书管理系统，实现高校证书的颁发、存储、验证全流程管理，解决证书造假问题，提升学历认证的公信力。

### 1.2 系统角色
| 角色 | 核心诉求 | 主要功能 |
|------|---------|---------|
| **监管方**（教育部/教育厅） | 监管全国高校证书颁发，建立统一标准 | 节点管理、证书管理、合约管理、数据监管、系统管理 |
| **高校**（成员方） | 便捷管理本校证书颁发，保护数据主权 | 身份认证、证书颁发管理、区块链存证、PDF生成与分发、查询统计、系统集成 |
| **验证方**（企业/其他高校/政府） | 快速、可靠验证证书真伪 | 多种验证方式、验证结果展示、批量验证、验证历史管理 |
| **学生** | 方便获取、管理和分享电子证书 | 证书获取、证书管理、授权管理、移动端支持 |
| **第三方开发者** | 通过API集成证书验证功能 | 开发者门户、API服务 |

---

## 2. 技术架构

### 2.1 技术栈选型

| 层级 | 技术选型 | 说明 |
|------|---------|------|
| **后端框架** | Gin 1.x | 高性能HTTP web框架 |
| **数据库** | PostgreSQL 14+ | 关系型数据库，存储核心业务数据 |
| **缓存** | Redis 7.x | 会话缓存、分布式锁、热点数据缓存 |
| **区块链** | Hyperledger Fabric 2.x / Ethereum | 联盟链，支持智能合约 |
| **消息队列** | RabbitMQ 3.x | 异步任务处理、事件通知 |
| **文件存储** | MinIO / 阿里云OSS | PDF证书、附件存储 |
| **PDF生成** | gofpdf / gopdf | 证书PDF生成 |
| **API文档** | Swagger/OpenAPI 3.0 | 接口文档自动生成 |
| **日志** | Zap / Logrus | 结构化日志 |
| **配置** | Viper | 配置管理 |
| **测试** | Testify / Ginkgo | 单元测试、集成测试 |

### 2.2 系统架构图

```
┌─────────────────────────────────────────────────────────────────────┐
│                         客户端层 (Client Layer)                       │
├─────────────┬─────────────┬─────────────┬─────────────┬─────────────┤
│   Web管理端  │  移动端H5   │  微信小程序  │  API调用   │  CLI工具    │
│  (Vue3/React)│ (响应式)     │ (Uni-app)   │ (第三方)   │ (Golang)    │
└─────────────┴─────────────┴─────────────┴─────────────┴─────────────┘
                                  │
                                  ▼
┌─────────────────────────────────────────────────────────────────────┐
│                        API网关层 (Gateway Layer)                      │
├─────────────────────────────────────────────────────────────────────┤
│  认证授权 │  限流熔断  │  负载均衡  │  请求路由  │  日志审计    │
│ (JWT/OAuth2)│ (Redis)    │ (Nginx)    │ (Gin)     │ (Middleware) │
└─────────────────────────────────────────────────────────────────────┘
                                  │
                                  ▼
┌─────────────────────────────────────────────────────────────────────┐
│                        服务层 (Service Layer)                        │
├─────────────┬─────────────┬─────────────┬─────────────┬─────────────┤
│  用户服务    │ 证书服务    │ 验证服务    │ 区块链服务  │ 消息服务    │
│  (User)     │ (Cert)     │ (Verify)   │ (Chain)    │ (Msg)      │
├─────────────┼─────────────┼─────────────┼─────────────┼─────────────┤
│  PDF服务    │ 系统服务    │ 统计服务    │ 存储服务    │ 支付服务    │
│  (PDF)     │ (System)   │ (Stats)    │ (Storage)  │ (Payment)  │
└─────────────┴─────────────┴─────────────┴─────────────┴─────────────┘
                                  │
                                  ▼
┌─────────────────────────────────────────────────────────────────────┐
│                        数据层 (Data Layer)                           │
├─────────────┬─────────────┬─────────────┬─────────────┬─────────────┤
│  PostgreSQL │   Redis     │  RabbitMQ   │    区块链    │  文件存储   │
│  (主数据库)  │  (缓存)      │  (消息队列)  │  (Fabric)   │  (MinIO)   │
└─────────────┴─────────────┴─────────────┴─────────────┴─────────────┘
```

### 2.3 微服务模块划分

```
edu-chain/
├── api/                    # API定义层（protobuf/OpenAPI）
│   ├── proto/             # gRPC协议定义
│   └── openapi/           # OpenAPI规范
├── cmd/                    # 应用入口
│   ├── server/            # 主服务启动
│   ├── cli/               # CLI命令行工具
│   └── verifier/          # 独立验证服务
├── internal/               # 内部模块
│   ├── biz/               # 业务逻辑层
│   │   ├── user/         # 用户业务
│   │   ├── certificate/  # 证书业务
│   │   ├── verification/ # 验证业务
│   │   ├── blockchain/   # 区块链业务
│   │   └── admin/        # 管理业务
│   ├── data/              # 数据访问层
│   │   ├── repository/   # 数据仓储
│   │   ├── cache/        # 缓存操作
│   │   └── db/           # 数据库操作
│   ├── service/           # 服务层（对外接口）
│   │   ├── user.go
│   │   ├── certificate.go
│   │   └── verification.go
│   ├── server/            # 服务启动配置
│   │   ├── http.go       # HTTP服务
│   │   └── grpc.go       # gRPC服务
│   ├── blockchain/        # 区块链SDK封装
│   │   ├── fabric/       # Fabric链码调用
│   │   └── ethereum/     # Ethereum合约调用
│   ├── routes/            # 路由定义
│   │   ├── router.go
│   │   ├── middleware.go
│   │   └── v1/           # API版本
│   ├── utils/             # 工具函数
│   │   ├── crypto/       # 加密工具
│   │   ├── pdf/          # PDF生成
│   │   └── jwt/          # JWT工具
│   └── config/            # 配置管理
├── pkg/                    # 公共包
│   ├── logger/            # 日志封装
│   ├── errors/            # 错误定义
│   └── response/          # 响应封装
├── scripts/               # 部署脚本
├── sql/                   # 数据库脚本
├── tests/                 # 测试文件
├── docs/                  # 文档
└── web/                   # 前端静态文件
```

---

## 3. 核心功能模块设计

### 3.1 用户认证模块 (User Auth)

```
功能点：
├── 登录认证
│   ├── 数字证书/UKey登录
│   ├── 统一身份认证对接
│   └── JWT Token管理
├── 权限管理
│   ├── RBAC角色权限模型
│   ├── 监管方：超级管理员、审计员、操作员
│   ├── 高校：校管理员、院系管理员、操作员
│   └── 验证方：HR、部门经理、审计员
└── 操作审计
    ├── 登录日志
    └── 操作日志
```

### 3.2 证书管理模块 (Certificate Management)

```
功能点：
├── 证书类型
│   ├── 学位证书（学士、硕士、博士）
│   ├── 成绩单
│   └── 在读证明
├── 证书模板
│   ├── 多模板支持
│   ├── 拖拽式模板编辑器
│   └── 校徽、签名等元素管理
├── 证书颁发
│   ├── 学位授予名单导入（Excel/CSV）
│   ├── 批量生成证书
│   ├── 单个证书补充颁发
│   └── 证书重发/补办
├── 证书生命周期
│   ├── 草稿 → 已生成 → 已签名 → 已上链 → 已分发
│   └── 证书状态跟踪
└── 数据同步
    ├── 教务系统对接（API/SFTP/数据库）
    └── 自动/手动同步
```

### 3.3 区块链存证模块 (Blockchain)

```
功能点：
├── 哈希计算
│   ├── SHA-256/SHA-3算法
│   └── 证书数据哈希
├── 数字签名
│   ├── 高校私钥签名
│   └── 硬件加密模块支持（HSM/UKey）
├── 区块链交互
│   ├── 上链交易
│   ├── 交易状态监控
│   ├── 重试机制（网络异常自动重试）
│   └── 离线签名支持
└── 智能合约
    ├── 证书存证合约
    ├── 验证合约
    └── 证书撤销合约
```

### 3.4 PDF生成模块 (PDF Generation)

```
功能点：
├── 证书生成
│   ├── PDF模板渲染
│   ├── 防伪水印添加
│   ├── 二维码集成（交易哈希）
│   └── 文件加密（密码保护）
├── 批量处理
│   ├── 批量生成
│   ├── 批量下载/打包
    └── 邮件自动发送
└── 分享功能
    ├── 验证链接生成
    └── 带时效的分享链接
```

### 3.5 验证服务模块 (Verification)

```
功能点：
├── 验证方式
│   ├── 文件上传验证（PDF/图片）
│   ├── 手动输入验证（交易哈希/学生ID）
│   ├── 扫描二维码
│   └── API批量验证
├── 验证结果
│   ├── 有效/警告/无效状态
│   ├── 颁发信息展示
│   ├── 区块链确认数
│   └── 验证报告生成
├── 批量验证
│   ├── Excel导入
│   ├── 进度跟踪
│   └── 结果导出
└── 验证历史
    ├── 验证记录查询
    ├── 验证备注
    └── 收藏管理
```

---

## 4. 数据库设计

### 4.1 核心表结构

```sql
-- 用户与权限
users                    -- 用户表
roles                    -- 角色表
permissions              -- 权限表
user_roles               -- 用户角色关联
organizations            -- 组织（高校/机构）
organization_users       -- 组织用户关联

-- 证书管理
certificates             -- 证书主表
certificate_templates    -- 证书模板
certificate_types        -- 证书类型
certificate_applications -- 证书申请表
certificate_records      -- 证书颁发记录

-- 区块链
blockchain_transactions  -- 区块链交易
blockchain_nodes         -- 区块链节点
certificate_hashes        -- 证书哈希记录

-- 验证服务
verification_logs        -- 验证日志
verification_requests    -- 验证请求
verification_reports     -- 验证报告

-- 系统管理
audit_logs               -- 审计日志
system_configs           -- 系统配置
data_backups             -- 数据备份

-- 文件存储
file_records             -- 文件记录
```

### 4.2 E-R关系图

```
users ──< user_roles >── roles ──< role_permissions >── permissions
                                                      │
organizations ──< organization_users >── users ───┘
     │
     └──< certificates >── certificate_records
                │
                └──< verification_logs >
```

---

## 5. API接口设计

### 5.1 API版本规划

| 版本 | 路径 | 说明 |
|------|------|------|
| v1 | `/api/v1/` | 第一版API，稳定版本 |
| v2 | `/api/v2/` | 第二版API，新功能 |
| internal | `/api/internal/` | 内部服务调用 |

### 5.2 主要API分组

```
认证模块
├── POST   /api/v1/auth/login
├── POST   /api/v1/auth/logout
├── POST   /api/v1/auth/refresh
└── GET    /api/v1/auth/info

证书模块
├── GET    /api/v1/certificates
├── POST   /api/v1/certificates
├── GET    /api/v1/certificates/:id
├── PUT    /api/v1/certificates/:id
├── DELETE /api/v1/certificates/:id
├── POST   /api/v1/certificates/batch
└── POST   /api/v1/certificates/import

验证模块
├── POST   /api/v1/verify
├── POST   /api/v1/verify/batch
├── GET    /api/v1/verify/history
└── GET    /api/v1/verify/:id

用户模块
├── GET    /api/v1/users/profile
├── PUT    /api/v1/users/profile
└── GET    /api/v1/users/verifications

管理模块
├── GET    /api/v1/admin/nodes
├── POST   /api/v1/admin/nodes
├── GET    /api/v1/admin/contracts
├── POST   /api/v1/admin/contracts
├── GET    /api/v1/admin/audit-logs
└── GET    /api/v1/admin/statistics
```

---

## 6. 部署架构

### 6.1 物理部署图

```
┌─────────────────────────────────────────────────────────────────────┐
│                          生产环境架构                                 │
├─────────────────────────────────────────────────────────────────────┤
│                                                                       │
│  ┌─────────────┐    ┌─────────────┐    ┌─────────────┐              │
│  │   Nginx     │    │   Nginx     │    │   Nginx     │              │
│  │   (LB)      │    │   (LB)      │    │   (LB)      │              │
│  └──────┬──────┘    └──────┬──────┘    └──────┬──────┘              │
│         │                 │                 │                       │
│         └────────────┬────┴─────────────────┘                       │
│                      ▼                                              │
│         ┌─────────────────────────────┐                              │
│         │      应用服务器集群          │                              │
│         │  ┌─────┐ ┌─────┐ ┌─────┐  │                              │
│         │  │App1 │ │App2 │ │App3 │  │                              │
│         │  └─────┘ └─────┘ └─────┘  │                              │
│         └─────────────────────────────┘                              │
│                      │                                              │
│    ┌─────────────────┼─────────────────┐                            │
│    ▼                 ▼                 ▼                              │
│ ┌────────┐      ┌────────┐       ┌────────┐                          │
│ │PostgreSQL│    │  Redis  │       │RabbitMQ│                          │
│ │(主从复制)│    │(集群)   │       │(集群)   │                          │
│ └────────┘      └────────┘       └────────┘                          │
│                      │                                              │
│                      ▼                                              │
│              ┌──────────────┐                                        │
│              │  区块链网络   │                                        │
│              │ (Fabric排序+ │                                        │
│              │  peer节点)   │                                        │
│              └──────────────┘                                        │
│                                                                       │
└─────────────────────────────────────────────────────────────────────┘
```

### 6.2 Docker Compose部署

```yaml
version: '3.8'
services:
  nginx:
    image: nginx:latest
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf
    depends_on:
      - app

  app:
    build: .
    ports:
      - "8080:8080"
    environment:
      - DATABASE_URL=postgres://user:pass@postgres:5432/edu_chain
      - REDIS_ADDR=redis:6379
      - RABBITMQ_URL=amqp://guest:guest@rabbitmq:5672/
    depends_on:
      - postgres
      - redis
      - rabbitmq

  postgres:
    image: postgres:14
    environment:
      POSTGRES_DB: edu_chain
      POSTGRES_USER: user
      POSTGRES_PASSWORD: pass
    volumes:
      - postgres_data:/var/lib/postgresql/data

  redis:
    image: redis:7-alpine
    volumes:
      - redis_data:/data

  rabbitmq:
    image: rabbitmq:3-management
    ports:
      - "15672:15672"
    volumes:
      - rabbitmq_data:/var/lib/rabbitmq

volumes:
  postgres_data:
  redis_data:
  rabbitmq_data:
```

---

## 7. 安全设计

### 7.1 安全措施

```
├── 传输安全
│   └── HTTPS/TLS 1.3
├── 认证授权
│   ├── JWT Token（RS256签名）
│   ├── RBAC权限控制
│   └── API密钥管理
├── 数据安全
│   ├── 数据库加密（透明加密）
│   ├── 敏感数据脱敏
│   └── 备份加密
├── 区块链安全
│   ├── 私钥安全存储（HSM）
│   ├── 多重签名机制
│   └── 智能合约审计
└── 应用安全
    ├── SQL注入防护
    ├── XSS防护
    ├── CSRF防护
    ├── 限流熔断
    └── 安全审计日志
```

---

## 8. 性能优化

### 8.1 优化策略

```
├── 数据库优化
│   ├── 读写分离
│   ├── 分库分表
│   ├── 索引优化
│   └── 连接池管理
├── 缓存优化
│   ├── 多级缓存架构
│   ├── 缓存预热
│   └── 缓存淘汰策略
├── 异步处理
│   ├── 消息队列削峰
│   ├── 任务队列处理
│   └── 定时任务调度
└── 区块链优化
    ├── 批量上链
    └── 链下数据存储
```

---

## 9. 监控运维

### 9.1 监控指标

```
├── 应用监控
│   ├── 请求量（QPS）
│   ├── 响应时间（P50/P95/P99）
│   └── 错误率
├── 系统监控
│   ├── CPU/内存/磁盘
│   ├── 网络带宽
│   └── 容器健康状态
├── 业务监控
│   ├── 证书颁发量
│   ├── 验证请求量
│   └── 上链成功率
└── 区块链监控
    ├── 节点状态
    ├── 交易确认时间
    └── 区块高度
```

### 9.2 日志收集

```
├── 应用日志
│   ├── 结构化日志（JSON格式）
│   └── 日志级别：DEBUG/INFO/WARN/ERROR
├── 日志存储
│   ├── ELK Stack（Elasticsearch + Logstash + Kibana）
│   └── Loki + Grafana
└── 日志分析
    ├── 错误告警
    ├── 性能分析
    └── 安全审计
```

---

## 10. 开发规范

### 10.1 代码规范

- **Go编码规范**：遵循Effective Go和Go Code Review Comments
- **API设计**：遵循RESTful最佳实践
- **文档**：使用Swagger生成API文档
- **测试**：单元测试覆盖率 > 80%

### 10.2 Git规范

- **分支策略**：Git Flow
  - `main`：主分支（生产环境）
  - `develop`：开发分支
  - `feature/*`：功能分支
  - `release/*`：发布分支
  - `hotfix/*`：紧急修复分支

- **提交规范**：
  ```
  feat: 新功能
  fix: 修复bug
  docs: 文档更新
  style: 代码格式
  refactor: 重构
  test: 测试
  chore: 构建/工具
  ```

---

## 11. 项目里程碑

| 阶段 | 时间 | 交付物 |
|------|------|--------|
| Phase 1 | 1-2周 | 基础框架搭建、用户认证模块 |
| Phase 2 | 3-4周 | 证书管理核心功能、区块链集成 |
| Phase 3 | 5-6周 | 验证服务、PDF生成 |
| Phase 4 | 7-8周 | 批量处理、系统集成 |
| Phase 5 | 9-10周 | 性能优化、安全加固 |
| Phase 6 | 11-12周 | 测试、上线部署 |

---

## 12. 风险评估

| 风险 | 影响 | 应对措施 |
|------|------|---------|
| 区块链性能瓶颈 | 高 | 链下缓存、批量处理 |
| 数据一致性 | 中 | 最终一致性、补偿机制 |
| 私钥泄露 | 高 | HSM硬件加密、多重签名 |
| 系统可用性 | 高 | 多活部署、自动故障转移 |
| 合规性 | 中 | 隐私保护、数据脱敏 |

---

## 13. 附录

### 13.1 参考文档
- Hyperledger Fabric官方文档
- Gin框架文档
- PostgreSQL最佳实践
- 区块链密码学标准

### 13.2 术语表
- **CA**：证书颁发机构（Certificate Authority）
- **PKI**：公钥基础设施（Public Key Infrastructure）
- **HSM**：硬件安全模块（Hardware Security Module）
- **Merkle Tree**：默克尔树（用于数据验证）