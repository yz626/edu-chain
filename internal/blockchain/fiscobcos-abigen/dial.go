package fiscobcosabigen

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/FISCO-BCOS/go-sdk/v3/client"

	"github.com/yz626/edu-chain/config"
)

// dialNode 使用 SDK 配置文件方式连接节点（bcos_sdk_create_by_config_file）。
//
// 背景：client.DialContext 在 Windows 上通过 bcos_sdk_create_config 构建配置，
// 该函数内部会 free 并重新赋值 cert_config 各字段指针，而 bcos_sdk_create 在
// Windows 上会转移这些指针的所有权，导致 defer bcos_sdk_c_config_destroy 时
// double-free，引发 0xc0000005 访问违规崩溃。
//
// 本函数改用 client.Dial（bcos_sdk_create_by_config_file），C SDK 自行解析
// INI 配置文件并管理内存，完全绕过上述问题，在 Windows/Linux/macOS 上均稳定。
func dialNode(cfg *config.BlockchainConfig) (*client.Client, error) {
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
