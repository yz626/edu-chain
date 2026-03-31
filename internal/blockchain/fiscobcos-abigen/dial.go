package fiscobcosabigen

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/FISCO-BCOS/go-sdk/v3/client"

	"github.com/yz626/edu-chain/config"
)

// DialMode 连接方式。
type DialMode int

const (
	// DialModeConfigFile 通过临时 INI 配置文件调用 bcos_sdk_create_by_config_file。
	// 推荐方式，在 Windows/Linux/macOS 上均稳定。
	DialModeConfigFile DialMode = iota

	// DialModeContext 通过 client.DialContext 调用 bcos_sdk_create_config。
	// 在 Linux/macOS 上正常，在 Windows 上因 C SDK 双重释放 cert_config 指针
	// 导致 0xc0000005 访问违规崩溃，仅供在非 Windows 环境或调试时使用。
	DialModeContext
)

// dialNode 使用默认方式（DialModeConfigFile）连接节点。
func dialNode(cfg *config.BlockchainConfig) (*client.Client, error) {
	return dialNodeWith(cfg, DialModeConfigFile)
}

// dialNodeWith 以指定方式连接节点。
//
// DialModeConfigFile（推荐）
//
//	将连接参数写入临时 INI 文件，通过 client.Dial →
//	bcos_sdk_create_by_config_file 建立连接。C SDK 自行管理内存，
//	在 Windows/Linux/macOS 上均无崩溃问题。
//
// DialModeContext
//
//	通过 client.DialContext → bcos_sdk_create_config 在代码中构建配置结构体。
//	在 Windows 上，bcos_sdk_create 会转移 cert_config 各字段指针的所有权，
//	导致 defer bcos_sdk_c_config_destroy 时 double-free，引发访问违规崩溃。
//	该方式保留用于在 Linux/macOS 服务器部署时的可选替代，以及问题复现/对比测试。
func dialNodeWith(cfg *config.BlockchainConfig, mode DialMode) (*client.Client, error) {
	switch mode {
	case DialModeConfigFile:
		return dialByConfigFile(cfg)
	case DialModeContext:
		return dialByContext(cfg)
	default:
		return nil, fmt.Errorf("unknown dial mode: %d", mode)
	}
}

// ----------------------------------------------------------------
// DialModeConfigFile 实现
// ----------------------------------------------------------------

// dialByConfigFile 将连接参数序列化为临时 INI 文件后调用 client.Dial。
func dialByConfigFile(cfg *config.BlockchainConfig) (*client.Client, error) {
	privKeyFile, err := resolveKeyPEMFile(cfg)
	if err != nil {
		return nil, err
	}
	privKey, _, err := client.LoadECPrivateKeyFromPEM(privKeyFile)
	if err != nil {
		return nil, fmt.Errorf("load private key: %w", err)
	}

	configFile, err := writeSDKConfigFile(cfg)
	if err != nil {
		return nil, fmt.Errorf("write sdk config file: %w", err)
	}

	return client.Dial(configFile, cfg.GroupID, privKey)
}

// writeSDKConfigFile 将连接参数序列化为 C SDK 期望的 INI 格式临时文件。
//
// C SDK 使用 INI 格式配置文件，关键节：
//   - [common]  : 线程数、超时、是否禁用 SSL
//   - [cert]    : ssl_type、ca_path、ca_cert、sdk_key、sdk_cert
//   - [peers]   : node.N=host:port
//
// 所有证书路径转换为绝对路径，避免因工作目录不同导致 C SDK 找不到文件。
func writeSDKConfigFile(cfg *config.BlockchainConfig) (string, error) {
	if len(cfg.Nodes) == 0 {
		return "", fmt.Errorf("no blockchain nodes configured")
	}

	absPath := func(p string) string {
		if p == "" || filepath.IsAbs(p) {
			return p
		}
		if abs, err := filepath.Abs(p); err == nil {
			return abs
		}
		return p
	}

	var buf strings.Builder

	// [common]
	buf.WriteString("[common]\n")
	buf.WriteString("    thread_pool_size = 4\n")
	buf.WriteString("    message_timeout_ms = 10000\n")
	if !cfg.TLS.Enabled {
		buf.WriteString("    disable_ssl = true\n")
	}
	buf.WriteString("\n")

	// [cert]
	buf.WriteString("[cert]\n")
	buf.WriteString("    ssl_type = ssl\n")
	if cfg.TLS.Enabled {
		caAbs := absPath(cfg.TLS.CACert)
		fmt.Fprintf(&buf, "    ca_path=%s\n", filepath.Dir(caAbs))
		fmt.Fprintf(&buf, "    ca_cert=%s\n", filepath.Base(caAbs))
		fmt.Fprintf(&buf, "    sdk_key=%s\n", filepath.Base(absPath(cfg.TLS.ClientKey)))
		fmt.Fprintf(&buf, "    sdk_cert=%s\n", filepath.Base(absPath(cfg.TLS.ClientCert)))
	} else {
		buf.WriteString("    ca_path=./\n")
	}
	buf.WriteString("\n")

	// [peers]
	buf.WriteString("[peers]\n")
	for i, node := range cfg.Nodes {
		fmt.Fprintf(&buf, "    node.%d=%s\n", i, node)
	}

	tmp, err := os.CreateTemp("", "bcos-sdk-config-*.ini")
	if err != nil {
		return "", fmt.Errorf("create temp config: %w", err)
	}
	if _, err := tmp.WriteString(buf.String()); err != nil {
		_ = tmp.Close()
		return "", fmt.Errorf("write temp config: %w", err)
	}
	_ = tmp.Close()
	return tmp.Name(), nil
}

// ----------------------------------------------------------------
// DialModeContext 实现
// ----------------------------------------------------------------

// dialByContext 通过 client.DialContext 连接节点（代码构建配置结构体方式）。
//
// 注意：在 Windows 上此方式会因 C SDK 双重释放指针导致崩溃，请勿在 Windows 上使用。
func dialByContext(cfg *config.BlockchainConfig) (*client.Client, error) {
	if len(cfg.Nodes) == 0 {
		return nil, fmt.Errorf("no blockchain nodes configured")
	}

	host, port, err := parseNodeAddr(cfg.Nodes[0])
	if err != nil {
		return nil, err
	}

	keyPEMFile, err := resolveKeyPEMFile(cfg)
	if err != nil {
		return nil, err
	}

	var caFile, tlsCertFile, tlsKeyFile string
	if cfg.TLS.Enabled {
		caFile = cfg.TLS.CACert
		tlsCertFile = cfg.TLS.ClientCert
		tlsKeyFile = cfg.TLS.ClientKey
	}

	// ParseConfigOptions(caFile, tlsKeyFile, tlsCertFile, keyPEMFile, groupID, host, port, isSMCrypto)
	sdkCfg, err := client.ParseConfigOptions(
		caFile, tlsKeyFile, tlsCertFile, keyPEMFile,
		cfg.GroupID, host, port, false,
	)
	if err != nil {
		return nil, fmt.Errorf("parse sdk config: %w", err)
	}
	if !cfg.TLS.Enabled {
		sdkCfg.DisableSsl = true
	}

	return client.DialContext(context.Background(), sdkCfg)
}

// parseNodeAddr 将 "host:port" 拆分为 host 字符串和 port 整数。
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
