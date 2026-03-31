package fiscobcosabigen

import (
	"encoding/hex"
	"fmt"
	"os"
	"strings"

	"golang.org/x/crypto/sha3"

	"github.com/yz626/edu-chain/config"
)

// certIDToBytes32 将链下证书 ID 经 keccak256 哈希转为合约 bytes32 键。
//
// 合约以 bytes32 为主键存储证书，而链下使用人类可读的字符串 ID（如 UUID）。
// keccak256 保证：任意长度字符串 → 固定 32 字节，且碰撞概率可忽略不计。
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

// hexToBytes32 将 64 位 hex 字符串（可带 0x 前缀）解码为 [32]byte。
//
// 证书哈希由上层以 hex 字符串传入，合约接口要求 bytes32 类型。
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

// resolveKeyPEMFile 返回账户私钥 PEM 文件的路径。
//
// 优先级：
//  1. account.key_file  — 直接使用文件路径
//  2. account.key       — inline PEM 或 base64，写入临时文件后返回路径
func resolveKeyPEMFile(cfg *config.BlockchainConfig) (string, error) {
	if cfg.Account.KeyFile != "" {
		return strings.TrimSpace(cfg.Account.KeyFile), nil
	}
	if cfg.Account.Key == "" {
		return "", fmt.Errorf("no private key: set account.key or account.key_file in blockchain.yaml")
	}

	raw := strings.TrimSpace(cfg.Account.Key)
	var pemData string
	if strings.HasPrefix(raw, "-----BEGIN") {
		pemData = raw + "\n"
	} else {
		pemData = "-----BEGIN PRIVATE KEY-----\n" + raw + "\n-----END PRIVATE KEY-----\n"
	}

	tmp, err := os.CreateTemp("", "bcos-key-*.pem")
	if err != nil {
		return "", fmt.Errorf("create temp key file: %w", err)
	}
	if _, err := tmp.WriteString(pemData); err != nil {
		_ = tmp.Close()
		return "", fmt.Errorf("write temp key file: %w", err)
	}
	_ = tmp.Close()
	return tmp.Name(), nil
}
