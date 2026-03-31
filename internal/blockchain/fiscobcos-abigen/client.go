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

// NewClient 使用默认连接方式（DialModeConfigFile）初始化客户端。
// 在所有平台上均稳定，推荐使用。
func NewClient(cfg *config.BlockchainConfig) (*Client, error) {
	return NewClientWith(cfg, DialModeConfigFile)
}

// NewClientWith 以指定连接方式初始化客户端。
//
// mode 可选值：
//   - DialModeConfigFile（默认）：通过 INI 配置文件调用 bcos_sdk_create_by_config_file，
//     在 Windows/Linux/macOS 上均稳定。
//   - DialModeContext：通过 client.DialContext 调用 bcos_sdk_create_config，
//     仅适用于 Linux/macOS，在 Windows 上会崩溃。
func NewClientWith(cfg *config.BlockchainConfig, mode DialMode) (*Client, error) {
	if cfg == nil {
		return nil, fmt.Errorf("blockchain config is nil")
	}
	c := &Client{cfg: cfg}
	if !cfg.Enabled {
		return c, nil
	}

	sdkClient, err := dialNodeWith(cfg, mode)
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

// newSDKClient 仅供测试使用：直接返回底层 go-sdk client。
// 正常业务代码请使用 NewClient 或 NewClientWith。
func newSDKClient(cfg *config.BlockchainConfig, mode DialMode) (*client.Client, error) {
	return dialNodeWith(cfg, mode)
}
