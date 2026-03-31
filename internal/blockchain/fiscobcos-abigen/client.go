package fiscobcosabigen

import (
	"fmt"

	"github.com/FISCO-BCOS/go-sdk/v3/client"
	"github.com/ethereum/go-ethereum/common"

	"github.com/yz626/edu-chain/config"
	certificate "github.com/yz626/edu-chain/contracts/go"
)

// Client FISCO BCOS 3.0 abigen 方案区块链客户端。
//
// 使用 abigen 生成的强类型绑定（contracts/go/certificate_registry.go），
// 通过 CertificateRegistrySession 调用合约，编译期检查所有参数类型。
type Client struct {
	cfg     *config.BlockchainConfig                // 区块链配置
	session *certificate.CertificateRegistrySession // 证书注册合约会话
}

// ErrDisabled 区块链未启用时返回此错误。
var ErrDisabled = fmt.Errorf(
	"blockchain is disabled (set blockchain.enabled=true in config/blockchain.yaml)")

// NewClient 根据区块链配置初始化客户端。
func NewClient(cfg *config.BlockchainConfig) (*Client, error) {
	if cfg == nil {
		return nil, fmt.Errorf("blockchain config is nil")
	}
	c := &Client{cfg: cfg}
	if !cfg.Enabled {
		return c, nil
	}

	sdkClient, err := dialNode(cfg)
	if err != nil {
		return nil, fmt.Errorf("dial node: %w", err)
	}

	address := common.HexToAddress(cfg.Contract.Address)
	instance, err := certificate.NewCertificateRegistry(address, sdkClient)
	if err != nil {
		return nil, fmt.Errorf("init contract binding: %w", err)
	}

	c.session = &certificate.CertificateRegistrySession{
		Contract:     instance,
		CallOpts:     *sdkClient.GetCallOpts(),
		TransactOpts: *sdkClient.GetTransactOpts(),
	}
	return c, nil
}

// Enabled 返回区块链是否已启用。
func (c *Client) Enabled() bool { return c.cfg.Enabled }

// newSDKClient 仅暴露给测试使用：直接返回底层 go-sdk client。
// 正常业务代码请使用 NewClient。
func newSDKClient(cfg *config.BlockchainConfig) (*client.Client, error) {
	return dialNode(cfg)
}
