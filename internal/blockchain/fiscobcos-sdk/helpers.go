package fiscobcossdk

import (
	"fmt"
	"os"

	"golang.org/x/crypto/sha3"
)

// keccak256Hash 计算 keccak256 哈希，返回 [32]byte。
func keccak256Hash(data []byte) [32]byte {
	h := sha3.NewLegacyKeccak256()
	h.Write(data)
	var result [32]byte
	copy(result[:], h.Sum(nil))
	return result
}

// loadABIFromFile 从文件路径读取合约 ABI JSON 字符串。
func loadABIFromFile(path string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("abi_file path is not configured")
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("read abi file %s: %w", path, err)
	}
	if len(data) == 0 {
		return "", fmt.Errorf("abi file %s is empty", path)
	}
	return string(data), nil
}
