package fiscobcosabigen

import (
	"context"
	"encoding/hex"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/FISCO-BCOS/go-sdk/v3/abi/bind"
	"github.com/FISCO-BCOS/go-sdk/v3/client"
	"github.com/ethereum/go-ethereum/common"
	"golang.org/x/crypto/sha3"

	certificate "github.com/yz626/edu-chain/contracts/go"
	"github.com/yz626/edu-chain/config"
)

// Client FISCO BCOS 3.0 abigen 方案区块链客户端。
//
// 使用 abigen 生成的强类型绑定（contracts/go/certificate_registry.go），
// 通过 CertificateRegistrySession 调用合约，编译期检查所有参数类型。
type Client struct {
	cfg     *config.BlockchainConfig
	session *certificate.CertificateRegistrySession
}

// NewClient 根据区块链配置初始化客户端。
func NewClient(cfg *config.BlockchainConfig) (*Client, error) {
	if cfg == nil {
		return nil, fmt.Errorf("blockchain config is nil")
	}
	c := &Client{cfg: cfg}
	if !cfg.Enabled {
		return c, nil
	}

	sdkCfg, err := buildSDKConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("build sdk config: %w", err)
	}
	sdkClient, err := client.DialContext(context.Background(), sdkCfg)
	if err != nil {
		return nil, fmt.Errorf("dial node: %w", err)
	}

	address := common.HexToAddress(cfg.Contract.Address)
	instance, err := certificate.NewCertificateRegistry(address, sdkClient)
	if err != nil {
		return nil, fmt.Errorf("init contract binding: %w", err)
	}

	// CertificateRegistrySession 预设 call/transact opts，调用更简洁
	c.session = &certificate.CertificateRegistrySession{
		Contract:     instance,
		CallOpts:     *sdkClient.GetCallOpts(),
		TransactOpts: *sdkClient.GetTransactOpts(),
	}
	return c, nil
}

// Enabled 返回区块链是否已启用。
func (c *Client) Enabled() bool { return c.cfg.Enabled }

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
	// 强类型调用：参数类型错误编译期报错
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

// RestoreCertificate 恢复证书。
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
	// 返回匿名结构体（abigen 生成）
	out, err := c.session.GetCertificate(certId32)
	if err != nil {
		return nil, fmt.Errorf("GetCertificate: %w", err)
	}
	rec := &CertOnChainRecord{
		CertHash:     hex.EncodeToString(out.CertHash[:]),
		Issuer:       out.Issuer.Hex(),
		IssuedAt:     time.Unix(int64(out.IssuedAt), 0),
		Revoked:      out.Revoked,
		RevokeReason: out.RevokeReason,
	}
	if out.RevokedAt > 0 {
		rec.RevokedAt = time.Unix(int64(out.RevokedAt), 0)
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
		AuthorizedAt: time.Unix(int64(out.AuthorizedAt), 0),
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

// ----------------------------------------------------------------
// 内部辅助
// ----------------------------------------------------------------

func buildSDKConfig(cfg *config.BlockchainConfig) (*client.Config, error) {
	if len(cfg.Nodes) == 0 {
		return nil, fmt.Errorf("no blockchain nodes configured")
	}
	privKeyHex, err := loadPrivateKey(cfg)
	if err != nil {
		return nil, err
	}
	host, port, err := parseNodeAddr(cfg.Nodes[0])
	if err != nil {
		return nil, err
	}
	var caFile, certFile, keyFile string
	if cfg.TLS.Enabled {
		caFile = cfg.TLS.CACert
		certFile = cfg.TLS.ClientCert
		keyFile = cfg.TLS.ClientKey
	}
	sdkCfg, err := client.ParseConfigOptions(
		caFile, privKeyHex, certFile, keyFile,
		cfg.GroupID, host, port, false,
	)
	if err != nil {
		return nil, fmt.Errorf("parse sdk config: %w", err)
	}
	if !cfg.TLS.Enabled {
		sdkCfg.DisableSsl = true
	}
	return sdkCfg, nil
}

func loadPrivateKey(cfg *config.BlockchainConfig) (string, error) {
	if cfg.Account.Key != "" {
		return strings.TrimSpace(cfg.Account.Key), nil
	}
	if cfg.Account.KeyFile == "" {
		return "", fmt.Errorf("no private key: set account.key or account.key_file")
	}
	data, err := os.ReadFile(cfg.Account.KeyFile)
	if err != nil {
		return "", fmt.Errorf("read key file: %w", err)
	}
	return strings.TrimSpace(string(data)), nil
}

func parseNodeAddr(addr string) (string, int, error) {
	parts := strings.SplitN(addr, ":", 2)
	if len(parts) != 2 {
		return "", 0, fmt.Errorf("expected host:port, got %q", addr)
	}
	port, err := strconv.Atoi(parts[1])
	if err != nil {
		return "", 0, fmt.Errorf("invalid port: %w", err)
	}
	return parts[0], port, nil
}

// certIDToBytes32 将链下 UUID 经 keccak256 转为合约 bytes32。
func certIDToBytes32(certID string) ([32]byte, error) {
	if certID == "" {
		return [32]byte{}, fmt.Errorf("certID is empty")
	}
	h := sha3.NewLegacyKeccak256()
	h.Write([]byte(certID))
	var result [32]byte
	copy(result[:], h.Sum(nil))
	return result, nil
}

// hexToBytes32 将 hex 字符串（可带 0x）转为 [32]byte。
func hexToBytes32(s string) ([32]byte, error) {
	s = strings.TrimPrefix(s, "0x")
	if len(s) != 64 {
		return [32]byte{}, fmt.Errorf("invalid bytes32 hex length %d (expected 64)", len(s))
	}
	b, err := hex.DecodeString(s)
	if err != nil {
		return [32]byte{}, err
	}
	var result [32]byte
	copy(result[:], b)
	return result, nil
}

// ErrDisabled 区块链未启用时返回此错误。
var ErrDisabled = fmt.Errorf(
	"blockchain is disabled (set blockchain.enabled=true in config/blockchain.yaml)")

// withTimeout 为 context 附加配置超时（供外部调用方使用）。
func (c *Client) withTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	timeout := time.Duration(c.cfg.Timeout) * time.Second
	if timeout <= 0 {
		timeout = 30 * time.Second
	}
	return context.WithTimeout(ctx, timeout)
}

// 确保 withTimeout 不被编译器报未使用
var _ = (*Client).withTimeout

// 需要 bind 包（避免 import 未使用报错）
var _ bind.CallOpts
