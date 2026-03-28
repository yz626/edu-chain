package fiscobcos

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/yz626/edu-chain/config"
)

// Contract 封装对 CertificateRegistry 合约的所有调用。
// 通过 JSON-RPC over HTTP 与 FISCO BCOS 3.0 节点通信，无 CGO 依赖。
type Contract struct {
	rpc        *rpcClient
	privKeyHex string
	cfg        *config.BlockchainConfig
	abi        *abiWrapper
	address    string // 合约地址（hex，带 0x）
}

// NewContract 初始化合约调用器。
func NewContract(rpc *rpcClient, privKeyHex string, cfg *config.BlockchainConfig) *Contract {
	addr := cfg.Contract.Address
	if !strings.HasPrefix(addr, "0x") && addr != "" {
		addr = "0x" + addr
	}
	return &Contract{
		rpc:        rpc,
		privKeyHex: privKeyHex,
		cfg:        cfg,
		address:    addr,
	}
}

// loadABI 懒加载并缓存 ABI。
func (c *Contract) loadABI() error {
	if c.abi != nil {
		return nil
	}
	abiJSON, err := loadABIFromFile(c.cfg.Contract.ABIFile)
	if err != nil {
		return err
	}
	w, err := newABIWrapper(abiJSON)
	if err != nil {
		return err
	}
	c.abi = w
	return nil
}

// ----------------------------------------------------------------
// 证书颁发
// ----------------------------------------------------------------

func (c *Contract) IssueCertificate(ctx context.Context, req IssueCertRequest) (*TxReceipt, error) {
	if err := c.loadABI(); err != nil {
		return nil, err
	}
	cid, err := uuidToBytes32(req.CertID)
	if err != nil {
		return nil, fmt.Errorf("certID: %w", err)
	}
	ch, err := hexToBytes32(req.CertHash)
	if err != nil {
		return nil, fmt.Errorf("certHash: %w", err)
	}
	data, err := c.abi.pack("issueCertificate", cid, ch)
	if err != nil {
		return nil, err
	}
	return c.sendTx(ctx, data)
}

func (c *Contract) IssueCertificateBatch(ctx context.Context, req BatchIssueCertRequest) (*TxReceipt, error) {
	if err := c.loadABI(); err != nil {
		return nil, err
	}
	certIDs := make([][32]byte, len(req.Items))
	certHashes := make([][32]byte, len(req.Items))
	for i, item := range req.Items {
		cid, err := uuidToBytes32(item.CertID)
		if err != nil {
			return nil, fmt.Errorf("item[%d] certID: %w", i, err)
		}
		ch, err := hexToBytes32(item.CertHash)
		if err != nil {
			return nil, fmt.Errorf("item[%d] certHash: %w", i, err)
		}
		certIDs[i] = cid
		certHashes[i] = ch
	}
	data, err := c.abi.pack("issueCertificateBatch", certIDs, certHashes)
	if err != nil {
		return nil, err
	}
	return c.sendTx(ctx, data)
}

// ----------------------------------------------------------------
// 证书撤销与恢复
// ----------------------------------------------------------------

func (c *Contract) RevokeCertificate(ctx context.Context, req RevokeCertRequest) (*TxReceipt, error) {
	if err := c.loadABI(); err != nil {
		return nil, err
	}
	cid, err := uuidToBytes32(req.CertID)
	if err != nil {
		return nil, fmt.Errorf("certID: %w", err)
	}
	data, err := c.abi.pack("revokeCertificate", cid, req.Reason)
	if err != nil {
		return nil, err
	}
	return c.sendTx(ctx, data)
}

func (c *Contract) RestoreCertificate(ctx context.Context, req RestoreCertRequest) (*TxReceipt, error) {
	if err := c.loadABI(); err != nil {
		return nil, err
	}
	cid, err := uuidToBytes32(req.CertID)
	if err != nil {
		return nil, fmt.Errorf("certID: %w", err)
	}
	data, err := c.abi.pack("restoreCertificate", cid)
	if err != nil {
		return nil, err
	}
	return c.sendTx(ctx, data)
}

// ----------------------------------------------------------------
// 查询与验证（call，不消耗 Gas）
// ----------------------------------------------------------------

func (c *Contract) CertExists(ctx context.Context, certID string) (bool, error) {
	if err := c.loadABI(); err != nil {
		return false, err
	}
	cid, err := uuidToBytes32(certID)
	if err != nil {
		return false, err
	}
	data, err := c.abi.pack("certExists", cid)
	if err != nil {
		return false, err
	}
	result, err := c.rpc.callContract(ctx, c.address, "0x"+hex.EncodeToString(data))
	if err != nil {
		return false, err
	}
	values, err := c.abi.unpack("certExists", hexToBytes(result))
	if err != nil {
		return false, err
	}
	if len(values) < 1 {
		return false, fmt.Errorf("certExists: empty result")
	}
	return values[0].(bool), nil
}

func (c *Contract) GetCertificate(ctx context.Context, certID string) (*CertOnChainRecord, error) {
	if err := c.loadABI(); err != nil {
		return nil, err
	}
	cid, err := uuidToBytes32(certID)
	if err != nil {
		return nil, err
	}
	data, err := c.abi.pack("getCertificate", cid)
	if err != nil {
		return nil, err
	}
	result, err := c.rpc.callContract(ctx, c.address, "0x"+hex.EncodeToString(data))
	if err != nil {
		return nil, err
	}
	values, err := c.abi.unpack("getCertificate", hexToBytes(result))
	if err != nil {
		return nil, err
	}
	if len(values) < 6 {
		return nil, fmt.Errorf("getCertificate: unexpected output count %d", len(values))
	}
	certHash32 := values[0].([32]byte)
	issuer := fmt.Sprintf("%v", values[1])
	issuedAtU64 := values[2].(uint64)
	revoked := values[3].(bool)
	revokedAtU64 := values[4].(uint64)
	revokeReason := values[5].(string)
	rec := &CertOnChainRecord{
		CertHash:     hex.EncodeToString(certHash32[:]),
		Issuer:       issuer,
		IssuedAt:     time.Unix(int64(issuedAtU64), 0),
		Revoked:      revoked,
		RevokeReason: revokeReason,
	}
	if revokedAtU64 > 0 {
		rec.RevokedAt = time.Unix(int64(revokedAtU64), 0)
	}
	return rec, nil
}

func (c *Contract) VerifyCertificate(ctx context.Context, certID, certHash string) (*VerifyResult, error) {
	if err := c.loadABI(); err != nil {
		return nil, err
	}
	cid, err := uuidToBytes32(certID)
	if err != nil {
		return nil, err
	}
	ch, err := hexToBytes32(certHash)
	if err != nil {
		return nil, err
	}
	data, err := c.abi.pack("verifyCertificate", cid, ch)
	if err != nil {
		return nil, err
	}
	result, err := c.rpc.callContract(ctx, c.address, "0x"+hex.EncodeToString(data))
	if err != nil {
		return nil, err
	}
	values, err := c.abi.unpack("verifyCertificate", hexToBytes(result))
	if err != nil {
		return nil, err
	}
	if len(values) < 2 {
		return nil, fmt.Errorf("verifyCertificate: unexpected output count %d", len(values))
	}
	return &VerifyResult{Valid: values[0].(bool), Revoked: values[1].(bool)}, nil
}

func (c *Contract) VerifyCertificateBatch(ctx context.Context, req BatchIssueCertRequest) (*BatchVerifyResult, error) {
	if err := c.loadABI(); err != nil {
		return nil, err
	}
	certIDs := make([][32]byte, len(req.Items))
	certHashes := make([][32]byte, len(req.Items))
	for i, item := range req.Items {
		cid, err := uuidToBytes32(item.CertID)
		if err != nil {
			return nil, fmt.Errorf("item[%d] certID: %w", i, err)
		}
		ch, err := hexToBytes32(item.CertHash)
		if err != nil {
			return nil, fmt.Errorf("item[%d] certHash: %w", i, err)
		}
		certIDs[i] = cid
		certHashes[i] = ch
	}
	data, err := c.abi.pack("verifyCertificateBatch", certIDs, certHashes)
	if err != nil {
		return nil, err
	}
	result, err := c.rpc.callContract(ctx, c.address, "0x"+hex.EncodeToString(data))
	if err != nil {
		return nil, err
	}
	values, err := c.abi.unpack("verifyCertificateBatch", hexToBytes(result))
	if err != nil {
		return nil, err
	}
	if len(values) < 2 {
		return nil, fmt.Errorf("verifyCertificateBatch: unexpected output count %d", len(values))
	}
	valids := values[0].([]bool)
	revokeds := values[1].([]bool)
	out := &BatchVerifyResult{Results: make([]VerifyResult, len(valids))}
	for i := range valids {
		out.Results[i] = VerifyResult{Valid: valids[i], Revoked: revokeds[i]}
	}
	return out, nil
}

func (c *Contract) AddIssuer(ctx context.Context, req AddIssuerRequest) (*TxReceipt, error) {
	if err := c.loadABI(); err != nil {
		return nil, err
	}
	addr := req.Address
	if !strings.HasPrefix(addr, "0x") {
		addr = "0x" + addr
	}
	data, err := c.abi.pack("addIssuer", addr, req.Name)
	if err != nil {
		return nil, err
	}
	return c.sendTx(ctx, data)
}

func (c *Contract) GetIssuerInfo(ctx context.Context, address string) (*IssuerInfo, error) {
	if err := c.loadABI(); err != nil {
		return nil, err
	}
	addr := address
	if !strings.HasPrefix(addr, "0x") {
		addr = "0x" + addr
	}
	data, err := c.abi.pack("getIssuerInfo", addr)
	if err != nil {
		return nil, err
	}
	result, err := c.rpc.callContract(ctx, c.address, "0x"+hex.EncodeToString(data))
	if err != nil {
		return nil, err
	}
	values, err := c.abi.unpack("getIssuerInfo", hexToBytes(result))
	if err != nil {
		return nil, err
	}
	if len(values) < 3 {
		return nil, fmt.Errorf("getIssuerInfo: unexpected output count %d", len(values))
	}
	return &IssuerInfo{
		Address:      address,
		Authorized:   values[0].(bool),
		Name:         values[1].(string),
		AuthorizedAt: time.Unix(int64(values[2].(uint64)), 0),
	}, nil
}

func (c *Contract) GetStats(ctx context.Context) (*ContractStats, error) {
	if err := c.loadABI(); err != nil {
		return nil, err
	}
	data, err := c.abi.pack("getStats")
	if err != nil {
		return nil, err
	}
	result, err := c.rpc.callContract(ctx, c.address, "0x"+hex.EncodeToString(data))
	if err != nil {
		return nil, err
	}
	values, err := c.abi.unpack("getStats", hexToBytes(result))
	if err != nil {
		return nil, err
	}
	if len(values) < 2 {
		return nil, fmt.Errorf("getStats: unexpected output count %d", len(values))
	}
	return &ContractStats{
		TotalIssued:  values[0].(*big.Int).Uint64(),
		TotalRevoked: values[1].(*big.Int).Uint64(),
	}, nil
}

// ----------------------------------------------------------------
// 内部：发交易
// ----------------------------------------------------------------

// sendTx 构建、签名并发送交易。
func (c *Contract) sendTx(ctx context.Context, input []byte) (*TxReceipt, error) {
	blockNum, err := c.rpc.getBlockNumber(ctx)
	if err != nil {
		return nil, fmt.Errorf("get block number: %w", err)
	}
	privKey, err := hexToECDSA(c.privKeyHex)
	if err != nil {
		return nil, fmt.Errorf("load private key: %w", err)
	}
	rawTx, err := buildSignedTransaction(privKey, c.address, input, blockNum+500)
	if err != nil {
		return nil, fmt.Errorf("build signed tx: %w", err)
	}
	receiptRaw, err := c.rpc.sendRawTransaction(ctx, "0x"+hex.EncodeToString(rawTx))
	if err != nil {
		return nil, fmt.Errorf("send raw transaction: %w", err)
	}
	return &TxReceipt{
		TxHash: receiptRaw.TransactionHash,
		Status: parseHexInt(receiptRaw.Status),
	}, nil
}

// ----------------------------------------------------------------
// 工具函数
// ----------------------------------------------------------------

func uuidToBytes32(certID string) ([32]byte, error) {
	if certID == "" {
		return [32]byte{}, fmt.Errorf("certID is empty")
	}
	return CertIDToBytes32(certID), nil
}

func hexToBytes32(hexStr string) ([32]byte, error) {
	hexStr = strings.TrimPrefix(hexStr, "0x")
	if len(hexStr) != 64 {
		return [32]byte{}, fmt.Errorf("invalid bytes32 hex length %d (expected 64)", len(hexStr))
	}
	b, err := hex.DecodeString(hexStr)
	if err != nil {
		return [32]byte{}, fmt.Errorf("decode hex: %w", err)
	}
	var result [32]byte
	copy(result[:], b)
	return result, nil
}

func hexToBytes(hexStr string) []byte {
	hexStr = strings.TrimPrefix(hexStr, "0x")
	b, _ := hex.DecodeString(hexStr)
	return b
}

func parseHexInt(hexStr string) int {
	hexStr = strings.TrimPrefix(hexStr, "0x")
	var n int64
	fmt.Sscanf(hexStr, "%x", &n)
	return int(n)
}
