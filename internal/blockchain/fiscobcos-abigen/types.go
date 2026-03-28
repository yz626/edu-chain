package fiscobcosabigen

import (
	"time"

	"github.com/FISCO-BCOS/go-sdk/v3/types"
)

// ================================================================
// FISCO BCOS 3.0 abigen 方案 — 业务类型定义
// ================================================================

// BatchItem 批量操作单项
type BatchItem struct {
	CertID   string // 链下 UUID
	CertHash string // sha256(certJSON) hex，不带 0x
}

// CertOnChainRecord 链上证书完整记录
type CertOnChainRecord struct {
	CertHash     string    // 证书完整内容的 SHA-256 哈希，用于链下数据完整性校验
	Issuer       string    // 颁发机构区块链地址
	IssuedAt     time.Time // 颁发时间戳（Unix 秒）
	Revoked      bool      // 是否已撤销
	RevokedAt    time.Time // 撤销时间戳（0 = 未撤销）
	RevokeReason string    // 撤销原因简述（详情存链下）
}

// VerifyResult 证书验证结果
type VerifyResult struct {
	Valid   bool // 证书是否有效
	Revoked bool // 证书是否已撤销
}

// TxReceipt 交易回执
type TxReceipt struct {
	TxHash      string // 交易哈希
	BlockNumber int64  // 区块号
	Status      int    // 状态
	Message     string // 消息
}

// IssuerInfo 颁发机构链上信息
type IssuerInfo struct {
	Address      string    // 颁发机构区块链地址
	Name         string    // 机构名称
	Authorized   bool      // 当前是否处于授权状态
	AuthorizedAt time.Time // 最近一次授权时间戳
}

// ContractStats 合约全局统计
type ContractStats struct {
	TotalIssued  uint64 // 总颁发数量
	TotalRevoked uint64 // 总撤销数量
}

// toTxReceipt 将 go-sdk v3 *types.Receipt 转为项目内部 TxReceipt。
func toTxReceipt(r *types.Receipt) *TxReceipt {
	if r == nil {
		return nil
	}
	return &TxReceipt{
		TxHash:      r.TransactionHash,
		BlockNumber: int64(r.BlockNumber),
		Status:      r.Status,
		Message:     r.Message,
	}
}
