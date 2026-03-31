package fiscobcosabigen

import (
	"context"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

// ----------------------------------------------------------------
// 证书颁发
// ----------------------------------------------------------------

// IssueCertificate 上链颁发单张证书。
func (c *Client) IssueCertificate(ctx context.Context, certID, certHash string) (*TxReceipt, error) {
	if !c.cfg.Enabled {
		return nil, ErrDisabled
	}
	certId32, err := certIDToBytes32(certID)
	if err != nil {
		return nil, err
	}
	certHash32, err := hexToBytes32(certHash)
	if err != nil {
		return nil, err
	}
	c.session.TransactOpts.Context = ctx
	_, receipt, err := c.session.IssueCertificate(certId32, certHash32)
	if err != nil {
		return nil, fmt.Errorf("IssueCertificate: %w", err)
	}
	return toTxReceipt(receipt), nil
}

// IssueCertificateBatch 批量上链颁发证书。
func (c *Client) IssueCertificateBatch(ctx context.Context, items []BatchItem) (*TxReceipt, error) {
	if !c.cfg.Enabled {
		return nil, ErrDisabled
	}
	if len(items) == 0 {
		return nil, fmt.Errorf("batch is empty")
	}
	certIds := make([][32]byte, len(items))
	certHashes := make([][32]byte, len(items))
	for i, item := range items {
		cid, err := certIDToBytes32(item.CertID)
		if err != nil {
			return nil, fmt.Errorf("item[%d] certID: %w", i, err)
		}
		ch, err := hexToBytes32(item.CertHash)
		if err != nil {
			return nil, fmt.Errorf("item[%d] certHash: %w", i, err)
		}
		certIds[i] = cid
		certHashes[i] = ch
	}
	c.session.TransactOpts.Context = ctx
	_, receipt, err := c.session.IssueCertificateBatch(certIds, certHashes)
	if err != nil {
		return nil, fmt.Errorf("IssueCertificateBatch: %w", err)
	}
	return toTxReceipt(receipt), nil
}

// ----------------------------------------------------------------
// 证书撤销与恢复
// ----------------------------------------------------------------

// RevokeCertificate 撤销证书。
func (c *Client) RevokeCertificate(ctx context.Context, certID, reason string) (*TxReceipt, error) {
	if !c.cfg.Enabled {
		return nil, ErrDisabled
	}
	certId32, err := certIDToBytes32(certID)
	if err != nil {
		return nil, err
	}
	c.session.TransactOpts.Context = ctx
	_, receipt, err := c.session.RevokeCertificate(certId32, reason)
	if err != nil {
		return nil, fmt.Errorf("RevokeCertificate: %w", err)
	}
	return toTxReceipt(receipt), nil
}

// RestoreCertificate 恢复已撤销的证书。
func (c *Client) RestoreCertificate(ctx context.Context, certID string) (*TxReceipt, error) {
	if !c.cfg.Enabled {
		return nil, ErrDisabled
	}
	certId32, err := certIDToBytes32(certID)
	if err != nil {
		return nil, err
	}
	c.session.TransactOpts.Context = ctx
	_, receipt, err := c.session.RestoreCertificate(certId32)
	if err != nil {
		return nil, fmt.Errorf("RestoreCertificate: %w", err)
	}
	return toTxReceipt(receipt), nil
}

// ----------------------------------------------------------------
// 查询与验证（call，不消耗 Gas）
// ----------------------------------------------------------------

// CertExists 查询证书是否存在。
func (c *Client) CertExists(ctx context.Context, certID string) (bool, error) {
	if !c.cfg.Enabled {
		return false, ErrDisabled
	}
	certId32, err := certIDToBytes32(certID)
	if err != nil {
		return false, err
	}
	c.session.CallOpts.Context = ctx
	return c.session.CertExists(certId32)
}

// GetCertificate 获取证书链上完整记录。
func (c *Client) GetCertificate(ctx context.Context, certID string) (*CertOnChainRecord, error) {
	if !c.cfg.Enabled {
		return nil, ErrDisabled
	}
	certId32, err := certIDToBytes32(certID)
	if err != nil {
		return nil, err
	}
	c.session.CallOpts.Context = ctx
	out, err := c.session.GetCertificate(certId32)
	if err != nil {
		return nil, fmt.Errorf("GetCertificate: %w", err)
	}
	rec := &CertOnChainRecord{
		CertHash:     hex.EncodeToString(out.CertHash[:]),
		Issuer:       out.Issuer.Hex(),
		IssuedAt:     time.UnixMilli(int64(out.IssuedAt)),
		Revoked:      out.Revoked,
		RevokeReason: out.RevokeReason,
	}
	if out.RevokedAt > 0 {
		rec.RevokedAt = time.UnixMilli(int64(out.RevokedAt))
	}
	return rec, nil
}

// VerifyCertificate 验证证书（哈希匹配且未撤销返回 Valid=true）。
func (c *Client) VerifyCertificate(ctx context.Context, certID, certHash string) (*VerifyResult, error) {
	if !c.cfg.Enabled {
		return nil, ErrDisabled
	}
	certId32, err := certIDToBytes32(certID)
	if err != nil {
		return nil, err
	}
	certHash32, err := hexToBytes32(certHash)
	if err != nil {
		return nil, err
	}
	c.session.CallOpts.Context = ctx
	out, err := c.session.VerifyCertificate(certId32, certHash32)
	if err != nil {
		return nil, fmt.Errorf("VerifyCertificate: %w", err)
	}
	return &VerifyResult{Valid: out.Valid, Revoked: out.Revoked}, nil
}

// ----------------------------------------------------------------
// 颁发机构管理
// ----------------------------------------------------------------

// AddIssuer 添加或重新授权颁发机构（仅合约 owner）。
func (c *Client) AddIssuer(ctx context.Context, address, name string) (*TxReceipt, error) {
	if !c.cfg.Enabled {
		return nil, ErrDisabled
	}
	c.session.TransactOpts.Context = ctx
	_, receipt, err := c.session.AddIssuer(common.HexToAddress(address), name)
	if err != nil {
		return nil, fmt.Errorf("AddIssuer: %w", err)
	}
	return toTxReceipt(receipt), nil
}

// GetIssuerInfo 查询颁发机构链上信息。
func (c *Client) GetIssuerInfo(ctx context.Context, address string) (*IssuerInfo, error) {
	if !c.cfg.Enabled {
		return nil, ErrDisabled
	}
	c.session.CallOpts.Context = ctx
	out, err := c.session.GetIssuerInfo(common.HexToAddress(address))
	if err != nil {
		return nil, fmt.Errorf("GetIssuerInfo: %w", err)
	}
	return &IssuerInfo{
		Address:      address,
		Name:         out.Name,
		Authorized:   out.Authorized,
		AuthorizedAt: time.UnixMilli(int64(out.AuthorizedAt)),
	}, nil
}

// GetStats 获取合约全局统计。
func (c *Client) GetStats(ctx context.Context) (*ContractStats, error) {
	if !c.cfg.Enabled {
		return nil, ErrDisabled
	}
	c.session.CallOpts.Context = ctx
	out, err := c.session.GetStats()
	if err != nil {
		return nil, fmt.Errorf("GetStats: %w", err)
	}
	return &ContractStats{
		TotalIssued:  out.TotalIssued.Uint64(),
		TotalRevoked: out.TotalRevoked.Uint64(),
	}, nil
}
