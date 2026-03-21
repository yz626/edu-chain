package redis

// 缓存Key前缀管理常量
// 用于区分不同业务模块的缓存key，避免key冲突
const (
	// ========== 认证相关 ==========
	// PrefixToken Token缓存前缀
	// 用法: token:{userID} -> 存储用户JWT Token
	PrefixToken = "token:"

	// PrefixSession 会话缓存前缀
	// 用法: session:{sessionID} -> 存储用户会话信息
	PrefixSession = "session:"

	// PrefixRefresh 刷新Token前缀
	// 用法: refresh:{userID} -> 存储刷新Token
	PrefixRefresh = "refresh:"

	// PrefixAuth 认证相关前缀（通用）
	PrefixAuth = "auth:"

	// ========== 用户相关 ==========
	// PrefixUser 用户缓存前缀
	// 用法: user:{userID} -> 存储用户基本信息
	PrefixUser = "user:"

	// PrefixUserPerm 用户权限缓存前缀
	// 用法: user:perm:{userID} -> 存储用户权限列表
	PrefixUserPerm = "user:perm:"

	// PrefixUserProfile 用户资料缓存前缀
	// 用法: user:profile:{userID} -> 存储用户详细信息
	PrefixUserProfile = "user:profile:"

	// ========== 证书相关 ==========
	// PrefixCert 证书缓存前缀
	// 用法: cert:{certID} -> 存储证书信息
	PrefixCert = "cert:"

	// PrefixCertMeta 证书元数据缓存前缀
	// 用法: cert:meta:{certID} -> 存储证书元数据（Hash结构）
	PrefixCertMeta = "cert:meta:"

	// PrefixCertTemplate 证书模板缓存前缀
	// 用法: cert:template:{templateID} -> 存储证书模板
	PrefixCertTemplate = "cert:template:"

	// PrefixCertBatch 证书批次缓存前缀
	// 用法: cert:batch:{batchID} -> 存储证书批次信息
	PrefixCertBatch = "cert:batch:"

	// ========== 验证相关 ==========
	// PrefixVerify 验证结果缓存前缀
	// 用法: verify:{verifyKey} -> 存储证书验证结果
	PrefixVerify = "verify:"

	// PrefixVerifyHistory 验证历史缓存前缀
	// 用法: verify:history:{certID} -> 存储验证历史记录
	PrefixVerifyHistory = "verify:history:"

	// ========== 组织相关 ==========
	// PrefixOrg 组织缓存前缀
	// 用法: org:{orgID} -> 存储组织信息
	PrefixOrg = "org:"

	// PrefixOrgUser 组织用户关联缓存前缀
	// 用法: org:user:{orgID} -> 存储组织下的用户列表
	PrefixOrgUser = "org:user:"

	// ========== 系统相关 ==========
	// PrefixSysConfig 系统配置缓存前缀
	// 用法: sys:config:{key} -> 存储系统配置
	PrefixSysConfig = "sys:config:"

	// PrefixRateLimit 限流缓存前缀
	// 用法: ratelimit:{key} -> 存储限流计数器
	PrefixRateLimit = "ratelimit:"

	// PrefixLock 分布式锁前缀
	// 用法: lock:{resource} -> 存储锁标记
	PrefixLock = "lock:"

	// PrefixSequence 序列号缓存前缀
	// 用法: sequence:{name} -> 存储分布式序列号
	PrefixSequence = "sequence:"

	// ========== 区块链相关 ==========
	// PrefixChain 区块链相关缓存前缀
	PrefixChain = "chain:"

	// PrefixTx 交易缓存前缀
	// 用法: chain:tx:{txHash} -> 存储交易信息
	PrefixTx = "chain:tx:"

	// PrefixBlock 区块缓存前缀
	// 用法: chain:block:{blockNumber} -> 存储区块信息
	PrefixBlock = "chain:block:"

	// PrefixChainConfig 链配置缓存前缀
	// 用法: chain:config:{networkID} -> 存储链配置
	PrefixChainConfig = "chain:config:"

	// ========== 审计相关 ==========
	// PrefixAudit 审计日志缓存前缀
	// 用法: audit:{logID} -> 存储审计日志
	PrefixAudit = "audit:"

	// ========== 临时缓存 ==========
	// PrefixTemp 临时缓存前缀（用于短期缓存）
	// 用法: temp:{key} -> 临时数据
	PrefixTemp = "temp:"
)
