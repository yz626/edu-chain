package fiscobcossdk

import (
	"context"
	"encoding/hex"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/FISCO-BCOS/go-sdk/v3/client"

	"github.com/yz626/edu-chain/config"
)

// Client FISCO BCOS 3.0 SDK 方案区块链客户端。
//
// 基于官方 go-sdk v3（github.com/FISCO-BCOS/go-sdk/v3），
// 使用 client.DialContext + ParseConfigOptions 连接节点。
// 需要 64 位 CGO 环境（MinGW-w64 / Linux GCC）。
//
// cfg.Enabled == false 时所有方法返回 ErrDisabled。
type Client struct {
	cfg       *config.BlockchainConfig
	sdkClient *client.Client
	contract  *Contract
}

// NewClient 根据区块链配置初始化 SDK 客户端。
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
		return nil, fmt.Errorf("dial fisco bcos node %s:%d: %w",
			sdkCfg.Host, sdkCfg.Port, err)
	}

	c.sdkClient = sdkClient
	c.contract = NewContract(sdkClient, cfg)
	return c, nil
}

// Close 关闭 SDK 连接。
func (c *Client) Close() {
	if c.sdkClient != nil {
		c.sdkClient.Close()
	}
}

// Enabled 返回区块链是否已启用。
func (c *Client) Enabled() bool { return c.cfg.Enabled }

// ----------------------------------------------------------------
// 证书颁发
// ----------------------------------------------------------------

// IssueCertificate 颁发证书。
func (c *Client) IssueCertificate(ctx context.Context, req IssueCertRequest) (*TxReceipt, error) {
	if !c.cfg.Enabled {
		return nil, ErrDisabled
	}
	ctx, cancel := c.withTimeout(ctx)
	defer cancel()
	return c.contract.IssueCertificate(ctx, req)
}

// IssueCertificateBatch 批量颁发证书。
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

// RevokeCertificate 撤销证书。
func (c *Client) RevokeCertificate(ctx context.Context, req RevokeCertRequest) (*TxReceipt, error) {
	if !c.cfg.Enabled {
		return nil, ErrDisabled
	}
	ctx, cancel := c.withTimeout(ctx)
	defer cancel()
	return c.contract.RevokeCertificate(ctx, req)
}

// RestoreCertificate 恢复证书。
func (c *Client) RestoreCertificate(ctx context.Context, req RestoreCertRequest) (*TxReceipt, error) {
	if !c.cfg.Enabled {
		return nil, ErrDisabled
	}
	ctx, cancel := c.withTimeout(ctx)
	defer cancel()
	return c.contract.RestoreCertificate(ctx, req)
}

// ----------------------------------------------------------------
// 查询与验证
// ----------------------------------------------------------------

func (c *Client) CertExists(ctx context.Context, certID string) (bool, error) {
	if !c.cfg.Enabled {
		return false, ErrDisabled
	}
	ctx, cancel := c.withTimeout(ctx)
	defer cancel()
	return c.contract.CertExists(ctx, certID)
}

func (c *Client) GetCertificate(ctx context.Context, certID string) (*CertOnChainRecord, error) {
	if !c.cfg.Enabled {
		return nil, ErrDisabled
	}
	ctx, cancel := c.withTimeout(ctx)
	defer cancel()
	return c.contract.GetCertificate(ctx, certID)
}

func (c *Client) VerifyCertificate(ctx context.Context, certID, certHash string) (*VerifyResult, error) {
	if !c.cfg.Enabled {
		return nil, ErrDisabled
	}
	ctx, cancel := c.withTimeout(ctx)
	defer cancel()
	return c.contract.VerifyCertificate(ctx, certID, certHash)
}

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

func (c *Client) AddIssuer(ctx context.Context, req AddIssuerRequest) (*TxReceipt, error) {
	if !c.cfg.Enabled {
		return nil, ErrDisabled
	}
	ctx, cancel := c.withTimeout(ctx)
	defer cancel()
	return c.contract.AddIssuer(ctx, req)
}

func (c *Client) GetIssuerInfo(ctx context.Context, address string) (*IssuerInfo, error) {
	if !c.cfg.Enabled {
		return nil, ErrDisabled
	}
	ctx, cancel := c.withTimeout(ctx)
	defer cancel()
	return c.contract.GetIssuerInfo(ctx, address)
}

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

// buildSDKConfig 将项目 BlockchainConfig 转为 go-sdk v3 client.Config。
// 使用 ParseConfigOptions(caFile, privKeyHex, certFile, keyFile,
//
//	groupId, host, port, isSMCrypto)
func buildSDKConfig(cfg *config.BlockchainConfig) (*client.Config, error) {
	if len(cfg.Nodes) == 0 {
		return nil, fmt.Errorf("no blockchain nodes configured")
	}
	privKeyHex, err := loadPrivateKey(cfg)
	if err != nil {
		return nil, fmt.Errorf("load private key: %w", err)
	}
	host, port, err := parseNodeAddr(cfg.Nodes[0])
	if err != nil {
		return nil, fmt.Errorf("parse node address %q: %w", cfg.Nodes[0], err)
	}

	var caFile, certFile, keyFile string
	if cfg.TLS.Enabled {
		caFile = cfg.TLS.CACert
		certFile = cfg.TLS.ClientCert
		keyFile = cfg.TLS.ClientKey
	}

	// 将私钥 hex 转为字节切片（go-sdk v3 ParseConfigOptions 接受 hex 字符串作为 key 参数）
	sdkCfg, err := client.ParseConfigOptions(
		caFile,      // CA 证书路径
		privKeyHex,  // 私钥十六进制（key 参数）
		certFile,    // 客户端证书路径
		keyFile,     // 客户端私钥文件路径
		cfg.GroupID, // 群组 ID
		host,        // 节点主机
		port,        // 节点端口
		false,       // isSMCrypto：非国密
	)
	if err != nil {
		return nil, fmt.Errorf("parse sdk config: %w", err)
	}
	if !cfg.TLS.Enabled {
		sdkCfg.DisableSsl = true
	}
	return sdkCfg, nil
}

// loadPrivateKey 读取私钥十六进制字符串。
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

// parseNodeAddr 解析 "host:port" 格式节点地址。
func parseNodeAddr(addr string) (host string, port int, err error) {
	parts := strings.SplitN(addr, ":", 2)
	if len(parts) != 2 {
		return "", 0, fmt.Errorf("expected host:port format")
	}
	port, err = strconv.Atoi(parts[1])
	if err != nil {
		return "", 0, fmt.Errorf("invalid port %q: %w", parts[1], err)
	}
	return parts[0], port, nil
}

// CertIDToBytes32 将链下 UUID 经 keccak256 转为合约 bytes32。
// 与合约约定一致：certId = keccak256(abi.encodePacked(uuidString))。
func CertIDToBytes32(certID string) [32]byte {
	return keccak256Hash([]byte(certID))
}

// hexToBytes32 将 hex 字符串（可带 0x，32字节）转为 [32]byte。
func hexToBytes32(hexStr string) ([32]byte, error) {
	hexStr = strings.TrimPrefix(hexStr, "0x")
	if len(hexStr) != 64 {
		return [32]byte{}, fmt.Errorf(
			"invalid bytes32 hex length %d (expected 64)", len(hexStr))
	}
	b, err := hex.DecodeString(hexStr)
	if err != nil {
		return [32]byte{}, fmt.Errorf("decode hex: %w", err)
	}
	var result [32]byte
	copy(result[:], b)
	return result, nil
}

// ErrDisabled 区块链功能未启用时返回此错误。
var ErrDisabled = fmt.Errorf(
	"blockchain is disabled (set blockchain.enabled=true in config/blockchain.yaml)")

// withTimeout 附加配置超时。
func (c *Client) withTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	timeout := time.Duration(c.cfg.Timeout) * time.Second
	if timeout <= 0 {
		timeout = 30 * time.Second
	}
	return context.WithTimeout(ctx, timeout)
}
