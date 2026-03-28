package fiscobcossdk

import (
	"time"
)

// ================================================================
// FISCO BCOS 3.0 SDK 方案 — 类型定义
// 基于 github.com/FISCO-BCOS/go-sdk/v3（官方 SDK）
// ================================================================

// IssueCertRequest 上链颁发证书请求
type IssueCertRequest struct {
	CertID   string // 链下 UUID 字符串（内部经 keccak256 转 bytes32）
	CertHash string // sha256(certJSON) hex（不带 0x，32字节）
}

// BatchIssueCertRequest 批量上链请求
type BatchIssueCertRequest struct {
	Items []IssueCertRequest
}

// RevokeCertRequest 撤销证书请求
type RevokeCertRequest struct {
	CertID string
	Reason string
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
	CertHash     string
	Issuer       string
	IssuedAt     time.Time
	Revoked      bool
	RevokedAt    time.Time
	RevokeReason string
}

// VerifyResult 单张证书验证结果
type VerifyResult struct {
	Valid   bool
	Revoked bool
}

// BatchVerifyResult 批量验证结果
type BatchVerifyResult struct {
	Results []VerifyResult
}

// TxReceipt 交易回执
type TxReceipt struct {
	TxHash      string
	BlockNumber int64
	Status      int
	Message     string
}

// IssuerInfo 颁发机构链上信息
type IssuerInfo struct {
	Address      string
	Name         string
	Authorized   bool
	AuthorizedAt time.Time
}

// ContractStats 合约全局统计
type ContractStats struct {
	TotalIssued  uint64
	TotalRevoked uint64
}
