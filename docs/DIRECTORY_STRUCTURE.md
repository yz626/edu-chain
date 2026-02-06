# EduChain 项目目录结构

## 完整目录结构

```
edu-chain/
├── api/                          # API定义层
│   ├── proto/                    # gRPC协议定义
│   │   ├── certificate.proto     # 证书相关协议
│   │   ├── user.proto            # 用户相关协议
│   │   └── verification.proto    # 验证相关协议
│   └── openapi/                  # OpenAPI规范
│       ├── v1.yaml              # v1 API规范
│       └── v2.yaml              # v2 API规范
│
├── cmd/                          # 应用入口
│   ├── server/                   # 主服务
│   │   └── main.go              # 服务启动入口
│   ├── cli/                      # CLI命令行工具
│   │   └── main.go              # CLI入口
│   └── verifier/                 # 独立验证服务
│       └── main.go              # 验证服务入口
│
├── internal/                      # 内部模块（不对外暴露）
│   ├── biz/                      # 业务逻辑层
│   │   ├── user/                # 用户业务
│   │   │   ├── user.go         # 用户业务逻辑
│   │   │   ├── auth.go         # 认证业务
│   │   │   └── permission.go   # 权限业务
│   │   ├── certificate/         # 证书业务
│   │   │   ├── certificate.go  # 证书业务逻辑
│   │   │   ├── template.go     # 模板管理
│   │   │   └── issuer.go       # 颁发管理
│   │   ├── verification/        # 验证业务
│   │   │   ├── verification.go  # 验证逻辑
│   │   │   └── batch.go        # 批量验证
│   │   ├── blockchain/          # 区块链业务
│   │   │   ├── chain.go        # 链上操作
│   │   │   └── transaction.go   # 交易管理
│   │   └── admin/              # 管理业务
│   │       ├── node.go         # 节点管理
│   │       ├── contract.go     # 合约管理
│   │       └── audit.go        # 审计管理
│   │
│   ├── data/                    # 数据访问层
│   │   ├── repository/         # 数据仓储
│   │   │   ├── user_repo.go    # 用户仓储
│   │   │   ├── cert_repo.go    # 证书仓储
│   │   │   └── verify_repo.go  # 验证仓储
│   │   ├── cache/              # 缓存操作
│   │   │   ├── redis.go        # Redis封装
│   │   │   └── session.go      # 会话缓存
│   │   └── db/                  # 数据库操作
│   │       ├── postgres.go      # PostgreSQL连接
│   │       ├── migrate.go       # 数据库迁移
│   │       └── seed.go         # 数据初始化
│   │
│   ├── service/                 # 服务层（对外接口）
│   │   ├── user_service.go     # 用户服务
│   │   ├── certificate_service.go  # 证书服务
│   │   ├── verification_service.go  # 验证服务
│   │   └── admin_service.go    # 管理服务
│   │
│   ├── server/                   # 服务启动配置
│   │   ├── http.go             # HTTP服务
│   │   ├── grpc.go             # gRPC服务
│   │   └── middleware.go      # 服务中间件
│   │
│   ├── blockchain/               # 区块链SDK封装
│   │   ├── fabric/            # Hyperledger Fabric
│   │   │   ├── client.go      # Fabric客户端
│   │   │   ├── chaincode.go   # 链码调用
│   │   │   └── identity.go    # 身份管理
│   │   └── ethereum/          # Ethereum
│   │       ├── client.go      # Ethereum客户端
│   │       ├── contract.go    # 合约调用
│   │       └── wallet.go      # 钱包管理
│   │
│   ├── routes/                   # 路由定义
│   │   ├── router.go           # 路由初始化
│   │   ├── middleware/         # 中间件
│   │   │   ├── cors.go        # 跨域中间件
│   │   │   ├── jwt.go         # JWT认证
│   │   │   ├── rate_limiter.go # 限流
│   │   │   ├── logger.go      # 日志
│   │   │   └── admin.go       # 权限控制
│   │   └── v1/                # API版本v1
│   │       ├── auth.go        # 认证接口
│   │       ├── certificate.go # 证书接口
│   │       ├── verification.go # 验证接口
│   │       ├── user.go        # 用户接口
│   │       └── admin.go       # 管理接口
│   │
│   ├── utils/                    # 工具函数
│   │   ├── crypto/             # 加密工具
│   │   │   ├── hash.go        # 哈希计算
│   │   │   ├── signature.go   # 数字签名
│   │   │   └── encrypt.go     # 加密解密
│   │   ├── pdf/               # PDF工具
│   │   │   ├── generator.go   # PDF生成
│   │   │   ├── template.go   # 模板渲染
│   │   │   └── watermark.go   # 水印添加
│   │   ├── jwt/               # JWT工具
│   │   │   ├── token.go       # Token生成
│   │   │   └── claims.go      # Claims定义
│   │   ├── excel/             # Excel工具
│   │   │   ├── reader.go      # Excel读取
│   │   │   └── writer.go      # Excel写入
│   │   └── time/              # 时间工具
│   │       └── format.go      # 时间格式化
│   │
│   └── config/                   # 配置管理
│       ├── config.go           # 配置加载
│       ├── default.go          # 默认配置
│       └── env.go             # 环境变量
│
├── pkg/                          # 公共包（可被外部引用）
│   ├── response/                # 统一响应
│   │   └── response.go        # 响应封装
│   ├── errors/                  # 错误定义
│   │   ├── errors.go          # 错误类型
│   │   └── codes.go           # 错误码
│   ├── logger/                 # 日志封装
│   │   ├── logger.go          # 日志接口
│   │   └── zap.go             # Zap实现
│   └── validator/              # 参数验证
│       └── validator.go       # 验证器
│
├── scripts/                      # 部署脚本
│   ├── build.sh                # 构建脚本
│   ├── deploy.sh               # 部署脚本
│   └── migrate.sh              # 数据库迁移
│
├── sql/                          # 数据库脚本
│   ├── schema.sql              # 表结构
│   ├── seed.sql               # 初始数据
│   └── migrations/             # 迁移文件
│       ├── 001_init.sql
│       └── 002_add_xxx.sql
│
├── tests/                        # 测试文件
│   ├── unit/                   # 单元测试
│   ├── integration/            # 集成测试
│   └── e2e/                   # 端到端测试
│
├── web/                          # 前端静态文件
│   ├── index.html             # 管理后台
│   ├── css/                   # 样式文件
│   ├── js/                    # 前端脚本
│   └── assets/                # 静态资源
│
├── docs/                        # 文档
│   ├── ARCHITECTURE.md        # 架构设计
│   ├── API.md                 # API文档
│   ├── DATABASE.md            # 数据库设计
│   └── DIRECTORY_STRUCTURE.md  # 目录结构
│
├── .gitignore                  # Git忽略文件
├── .golangci.yml               # Lint配置
├── .dockerignore               # Docker忽略文件
├── Dockerfile                  # Docker构建文件
├── docker-compose.yml          # Docker编排
├── go.mod                      # Go模块文件
├── go.sum                      # Go依赖校验
├── Makefile                    # Make命令
├── README.md                   # 项目说明
└── config.yaml.example         # 配置示例
```

## 目录说明

| 目录/文件 | 说明 |
|----------|------|
| `api/` | API接口定义，使用proto或OpenAPI规范 |
| `cmd/` | 应用入口，每个子目录对应一个可执行程序 |
| `internal/` | 内部模块，Go的可见性规则保护 |
| `pkg/` | 公共包，可以被外部项目引用 |
| `scripts/` | 部署、构建等脚本 |
| `sql/` | 数据库相关脚本 |
| `tests/` | 测试文件 |
| `web/` | 前端静态资源 |
| `docs/` | 项目文档 |

## 模块职责

### internal/biz/ (业务逻辑层)
- 封装核心业务逻辑
- 不依赖外部框架
- 可测试性强

### internal/data/ (数据访问层)
- 封装数据库操作
- 实现仓储模式
- 缓存管理

### internal/service/ (服务层)
- 对外提供接口
- 事务管理
- 参数校验

### internal/routes/ (路由层)
- HTTP路由定义
- 中间件配置
- API版本管理

## 依赖注入建议

推荐使用 Wire 进行依赖注入：

```
internal/
├── di/                        # 依赖注入
│   ├── wire.go              # Wire配置
│   ├── wire_gen.go          # 生成代码
│   └── module.go            # 模块定义