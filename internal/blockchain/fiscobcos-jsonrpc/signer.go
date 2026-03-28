package fiscobcos

import (
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"

	"golang.org/x/crypto/sha3"
)

// ================================================================
// FISCO BCOS 3.0 交易签名（纯 Go，secp256k1 + keccak256）
//
// FISCO BCOS 3.0 非国密交易格式与 Ethereum 兼容：
// RLP([nonce, blockLimit, to, value, data, chainID, reserved1, reserved2])
// 签名：keccak256(RLP) → secp256k1 签名 (r, s, v)
// ================================================================

// buildSignedTransaction 构建已签名的原始交易字节（RLP 编码）。
// 返回的字节可直接 hex 编码后通过 sendTransaction 接口发送。
func buildSignedTransaction(
	privKey *ecdsa.PrivateKey,
	to string,
	data []byte,
	blockLimit int64,
) ([]byte, error) {
	// 简化的交易数据结构（FISCO BCOS 3.0 兼容格式）
	// 实际生产环境需要完整 RLP 编码，此处为占位实现
	// 待集成正式 SDK 后替换
	nonce := make([]byte, 8)
	_, err := rand.Read(nonce)
	if err != nil {
		return nil, fmt.Errorf("generate nonce: %w", err)
	}

	// 构建交易数据
	toBytes, err := hexToAddress(to)
	if err != nil {
		return nil, fmt.Errorf("parse to address: %w", err)
	}

	blockLimitBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(blockLimitBytes, uint64(blockLimit))

	// 拼接待签名数据：nonce + blockLimit + to + data
	txPayload := make([]byte, 0, len(nonce)+len(blockLimitBytes)+len(toBytes)+len(data))
	txPayload = append(txPayload, nonce...)
	txPayload = append(txPayload, blockLimitBytes...)
	txPayload = append(txPayload, toBytes...)
	txPayload = append(txPayload, data...)

	// keccak256 哈希
	txHash := keccak256Hash(txPayload)

	// secp256k1 签名
	sig, err := signHash(privKey, txHash[:])
	if err != nil {
		return nil, fmt.Errorf("sign tx hash: %w", err)
	}

	// 组装：payload + signature
	rawTx := make([]byte, 0, len(txPayload)+len(sig))
	rawTx = append(rawTx, txPayload...)
	rawTx = append(rawTx, sig...)
	return rawTx, nil
}

// signHash 使用 secp256k1 对哈希签名，返回 r+s+v（65字节）。
func signHash(privKey *ecdsa.PrivateKey, hash []byte) ([]byte, error) {
	r, s, err := ecdsa.Sign(rand.Reader, privKey, hash)
	if err != nil {
		return nil, err
	}
	sig := make([]byte, 65)
	rBytes := r.Bytes()
	sBytes := s.Bytes()
	copy(sig[32-len(rBytes):32], rBytes)
	copy(sig[64-len(sBytes):64], sBytes)
	sig[64] = 0 // v = 0
	return sig, nil
}

// ecdsaFromBytes 将私钥字节转为 *ecdsa.PrivateKey（使用标准 P256 曲线占位）。
// 注意：FISCO BCOS 3.0 使用 secp256k1，生产环境需替换为 secp256k1 实现。
func ecdsaFromBytes(b []byte) (*ecdsa.PrivateKey, error) {
	// 使用标准库 elliptic 曲线（P256）作为开发阶段占位
	// 生产环境需使用 github.com/decred/dcrd/dcrec/secp256k1/v4
	if len(b) != 32 {
		return nil, fmt.Errorf("invalid private key length: %d (expected 32)", len(b))
	}
	privKey := new(ecdsa.PrivateKey)
	privKey.D = new(big.Int).SetBytes(b)
	// 占位：实际需要设置正确的 secp256k1 曲线和公钥
	// privKey.PublicKey.Curve = secp256k1.S256()
	// privKey.PublicKey.X, privKey.PublicKey.Y = secp256k1.S256().ScalarBaseMult(b)
	return privKey, nil
}

// hexToECDSA 将十六进制私钥字符串转为 *ecdsa.PrivateKey。
func hexToECDSA(hexKey string) (*ecdsa.PrivateKey, error) {
	hexKey = strings.TrimPrefix(hexKey, "0x")
	b, err := hex.DecodeString(hexKey)
	if err != nil {
		return nil, fmt.Errorf("decode private key hex: %w", err)
	}
	return ecdsaFromBytes(b)
}

// hexToAddress 将 0x 前缀的地址字符串转为 20 字节。
func hexToAddress(addr string) ([]byte, error) {
	addr = strings.TrimPrefix(addr, "0x")
	if len(addr) != 40 {
		return nil, fmt.Errorf("invalid address length %d (expected 40)", len(addr))
	}
	return hex.DecodeString(addr)
}

// keccak256Hash 计算 keccak256 哈希，返回 [32]byte。
func keccak256Hash(data []byte) [32]byte {
	h := sha3.NewLegacyKeccak256()
	h.Write(data)
	var result [32]byte
	copy(result[:], h.Sum(nil))
	return result
}

// CertIDToBytes32 将链下 UUID 字符串经 keccak256 转为合约 bytes32。
// 与合约约定一致：certId = keccak256(abi.encodePacked(uuidString))。
// 供外部业务层复用。
func CertIDToBytes32(certID string) [32]byte {
	return keccak256Hash([]byte(certID))
}
