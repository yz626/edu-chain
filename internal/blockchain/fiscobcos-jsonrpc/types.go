package fiscobcos

import (
	"time"
)

// ================================================================
// FISCO BCOS 3.0 区块链交互类型定义
// ================================================================

// IssueCertRequest 上链颁发证书请求
type IssueCertRequest struct {
	// CertID: 链下 UUID 字符串（内部经 keccak256 转 bytes32）
	CertID string
	// CertHash: sha256(certJSON) hex（不带 0x，32字节）
	CertHash string
}

// BatchIssueCertRequest 批量上链请求
type BatchIssueCertRequest struct {
	Items []IssueCertRequest
}

// RevokeCertRequest 撤销证书请求
type RevokeCertRequest struct {
	CertID string // 链下 UUID
	Reason string // 撤销原因简述
}

// RestoreCertRequest 恢复证书请求
type RestoreCertRequest struct {
	CertID string
}

// AddIssuerRequest 添加授权颁发机构请求
type AddIssuerRequest struct {
	Address string // 机构区块链地址
	Name    string // 机构名称
}

// CertOnChainRecord 链上证书完整记录
type CertOnChainRecord struct {
	CertHash     string    // 证书内容哈希（hex，不带 0x）
	Issuer       string    // 颁发机构区块链地址
	IssuedAt     time.Time // 颁发时间
	Revoked      bool      // 是否已撤销
	RevokedAt    time.Time // 撤销时间（零值表示未撤销）
	RevokeReason string    // 撤销原因
}

// VerifyResult 单张证书验证结果
type VerifyResult struct {
	Valid   bool // true = 哈希匹配且未撤销
	Revoked bool // true = 已被撤销
}

// BatchVerifyResult 批量验证结果
type BatchVerifyResult struct {
	Results []VerifyResult
}

// TxReceipt 交易回执（写操作上链结果）
type TxReceipt struct {
	TxHash      string // 交易哈希
	BlockNumber int64  // 所在区块高度
	Status      int    // 0 = 成功，非 0 = 失败
}

// IssuerInfo 颁发机构链上信息
type IssuerInfo struct {
	Address      string    // 机构区块链地址
	Name         string    // 机构名称
	Authorized   bool      // 是否已授权
	AuthorizedAt time.Time // 最近授权时间
}

// ContractStats 合约全局统计
type ContractStats struct {
	TotalIssued  uint64 // 已颁发证书总数
	TotalRevoked uint64 // 已撤销证书净数
}
