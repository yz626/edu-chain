package fiscobcos

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/yz626/edu-chain/config"
)

// Client FISCO BCOS 3.0 区块链客户端（纯 Go，无 CGO）。
//
// 通过 JSON-RPC over HTTP 与节点通信，无需 CGO 或 C 语言库。
// 写操作（上链）返回 *TxReceipt，读操作直接返回业务结构体。
// cfg.Enabled == false 时所有方法返回 ErrDisabled。
type Client struct {
	cfg      *config.BlockchainConfig
	rpc      *rpcClient
	contract *Contract
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

	if len(cfg.Nodes) == 0 {
		return nil, fmt.Errorf("no blockchain nodes configured")
	}

	// 加载私钥
	privKeyHex, err := loadPrivateKey(cfg)
	if err != nil {
		return nil, fmt.Errorf("load private key: %w", err)
	}

	// 构建 HTTP 端点（节点地址格式：host:port）
	endpoint := buildEndpoint(cfg.Nodes[0], cfg.TLS.Enabled)
	timeout := time.Duration(cfg.Timeout) * time.Second

	rpcCli := newRPCClient(endpoint, cfg.GroupID, timeout)
	c.rpc = rpcCli
	c.contract = NewContract(rpcCli, privKeyHex, cfg)
	return c, nil
}

// Close 关闭客户端（HTTP 客户端无需显式关闭）。
func (c *Client) Close() {}

// Enabled 返回区块链是否已启用。
func (c *Client) Enabled() bool { return c.cfg.Enabled }

// ----------------------------------------------------------------
// 证书颁发
// ----------------------------------------------------------------

// IssueCertificate 单张证书上链存证。
//   - req.CertID  : 链下 UUID 字符串（内部经 keccak256 转 bytes32）
//   - req.CertHash: sha256(certJSON) hex 字符串（不带 0x，32字节）
func (c *Client) IssueCertificate(ctx context.Context, req IssueCertRequest) (*TxReceipt, error) {
	if !c.cfg.Enabled {
		return nil, ErrDisabled
	}
	ctx, cancel := c.withTimeout(ctx)
	defer cancel()
	return c.contract.IssueCertificate(ctx, req)
}

// IssueCertificateBatch 批量证书上链。
func (c *Client) IssueCertificateBatch(ctx context.Context, req BatchIssueCertRequest) (*TxReceipt, error) {
	if !c.cfg.Enabled {
		return nil, ErrDisabled
	}
	if len(req.Items) == 0 {
		return nil, fmt.Errorf("batch is empty")
	}
	ctx, cancel := c.withTimeout(ctx)
	defer cancel()
	return c.contract.IssueCertificateBatch(ctx, req)
}

// ----------------------------------------------------------------
// 证书撤销与恢复
// ----------------------------------------------------------------

// RevokeCertificate 撤销证书（仅原颁发机构）。
func (c *Client) RevokeCertificate(ctx context.Context, req RevokeCertRequest) (*TxReceipt, error) {
	if !c.cfg.Enabled {
		return nil, ErrDisabled
	}
	ctx, cancel := c.withTimeout(ctx)
	defer cancel()
	return c.contract.RevokeCertificate(ctx, req)
}

// RestoreCertificate 恢复被撤销的证书。
func (c *Client) RestoreCertificate(ctx context.Context, req RestoreCertRequest) (*TxReceipt, error) {
	if !c.cfg.Enabled {
		return nil, ErrDisabled
	}
	ctx, cancel := c.withTimeout(ctx)
	defer cancel()
	return c.contract.RestoreCertificate(ctx, req)
}

// ----------------------------------------------------------------
// 查询与验证（只读，无 Gas）
// ----------------------------------------------------------------

// CertExists 查询证书是否存在。
func (c *Client) CertExists(ctx context.Context, certID string) (bool, error) {
	if !c.cfg.Enabled {
		return false, ErrDisabled
	}
	ctx, cancel := c.withTimeout(ctx)
	defer cancel()
	return c.contract.CertExists(ctx, certID)
}

// GetCertificate 获取证书链上完整记录。
func (c *Client) GetCertificate(ctx context.Context, certID string) (*CertOnChainRecord, error) {
	if !c.cfg.Enabled {
		return nil, ErrDisabled
	}
	ctx, cancel := c.withTimeout(ctx)
	defer cancel()
	return c.contract.GetCertificate(ctx, certID)
}

// VerifyCertificate 验证证书（哈希匹配且未撤销返回 Valid=true）。
func (c *Client) VerifyCertificate(ctx context.Context, certID, certHash string) (*VerifyResult, error) {
	if !c.cfg.Enabled {
		return nil, ErrDisabled
	}
	ctx, cancel := c.withTimeout(ctx)
	defer cancel()
	return c.contract.VerifyCertificate(ctx, certID, certHash)
}

// VerifyCertificateBatch 批量验证证书。
func (c *Client) VerifyCertificateBatch(ctx context.Context, req BatchIssueCertRequest) (*BatchVerifyResult, error) {
	if !c.cfg.Enabled {
		return nil, ErrDisabled
	}
	ctx, cancel := c.withTimeout(ctx)
	defer cancel()
	return c.contract.VerifyCertificateBatch(ctx, req)
}

// ----------------------------------------------------------------
// 颁发机构管理
// ----------------------------------------------------------------

// AddIssuer 添加或重新授权颁发机构（仅合约 owner）。
func (c *Client) AddIssuer(ctx context.Context, req AddIssuerRequest) (*TxReceipt, error) {
	if !c.cfg.Enabled {
		return nil, ErrDisabled
	}
	ctx, cancel := c.withTimeout(ctx)
	defer cancel()
	return c.contract.AddIssuer(ctx, req)
}

// GetIssuerInfo 查询颁发机构链上信息。
func (c *Client) GetIssuerInfo(ctx context.Context, address string) (*IssuerInfo, error) {
	if !c.cfg.Enabled {
		return nil, ErrDisabled
	}
	ctx, cancel := c.withTimeout(ctx)
	defer cancel()
	return c.contract.GetIssuerInfo(ctx, address)
}

// GetStats 获取合约全局统计。
func (c *Client) GetStats(ctx context.Context) (*ContractStats, error) {
	if !c.cfg.Enabled {
		return nil, ErrDisabled
	}
	ctx, cancel := c.withTimeout(ctx)
	defer cancel()
	return c.contract.GetStats(ctx)
}

// ----------------------------------------------------------------
// 内部辅助
// ----------------------------------------------------------------

// buildEndpoint 构建 HTTP/HTTPS 端点 URL。
func buildEndpoint(nodeAddr string, tlsEnabled bool) string {
	scheme := "http"
	if tlsEnabled {
		scheme = "https"
	}
	if strings.HasPrefix(nodeAddr, "http") {
		return nodeAddr
	}
	return fmt.Sprintf("%s://%s", scheme, nodeAddr)
}

// loadPrivateKey 从配置读取私钥十六进制字符串。
// 优先 cfg.Account.Key（环境变量注入），其次读文件。
func loadPrivateKey(cfg *config.BlockchainConfig) (string, error) {
	if cfg.Account.Key != "" {
		return strings.TrimSpace(cfg.Account.Key), nil
	}
	if cfg.Account.KeyFile == "" {
		return "", fmt.Errorf(
			"no private key: set account.key or account.key_file in blockchain.yaml")
	}
	data, err := os.ReadFile(cfg.Account.KeyFile)
	if err != nil {
		return "", fmt.Errorf("read key file %s: %w", cfg.Account.KeyFile, err)
	}
	return strings.TrimSpace(string(data)), nil
}

// withTimeout 为 context 附加配置超时。
func (c *Client) withTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	timeout := time.Duration(c.cfg.Timeout) * time.Second
	if timeout <= 0 {
		timeout = 30 * time.Second
	}
	return context.WithTimeout(ctx, timeout)
}

// ----------------------------------------------------------------
// 哨兵错误
// ----------------------------------------------------------------

// ErrDisabled 区块链功能未启用时返回此错误。
var ErrDisabled = fmt.Errorf(
	"blockchain is disabled (set blockchain.enabled=true in config/blockchain.yaml)")
