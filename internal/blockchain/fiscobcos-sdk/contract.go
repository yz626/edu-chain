package fiscobcossdk

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/FISCO-BCOS/go-sdk/v3/abi"
	"github.com/FISCO-BCOS/go-sdk/v3/client"
	"github.com/FISCO-BCOS/go-sdk/v3/types"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"

	"github.com/yz626/edu-chain/config"
)

// Contract 封装对 CertificateRegistry 合约的所有调用。
//
// 使用 go-sdk v3 自带的 abi 包进行编解码（github.com/FISCO-BCOS/go-sdk/v3/abi），
// 通过 SDK 的 CreateEncodedTransactionDataV1 + SendEncodedTransaction 发送交易，
// 通过 CallContract 执行只读调用。
type Contract struct {
	sdkClient *client.Client
	cfg       *config.BlockchainConfig
	abi       abi.ABI
	address   common.Address
	abiLoaded bool
}

// NewContract 初始化合约调用器。
func NewContract(sdkClient *client.Client, cfg *config.BlockchainConfig) *Contract {
	return &Contract{
		sdkClient: sdkClient,
		cfg:       cfg,
		address:   common.HexToAddress(cfg.Contract.Address),
	}
}

// loadABI 懒加载并缓存合约 ABI。
func (c *Contract) loadABI() error {
	if c.abiLoaded {
		return nil
	}
	abiJSON, err := loadABIFromFile(c.cfg.Contract.ABIFile)
	if err != nil {
		return err
	}
	parsed, err := abi.JSON(strings.NewReader(abiJSON))
	if err != nil {
		return fmt.Errorf("parse contract abi: %w", err)
	}
	c.abi = parsed
	c.abiLoaded = true
	return nil
}

// ----------------------------------------------------------------
// 证书颁发
// ----------------------------------------------------------------

func (c *Contract) IssueCertificate(ctx context.Context, req IssueCertRequest) (*TxReceipt, error) {
	if err := c.loadABI(); err != nil {
		return nil, err
	}
	cid, err := certIDToBytes32(req.CertID)
	if err != nil {
		return nil, fmt.Errorf("certID: %w", err)
	}
	ch, err := hexToBytes32(req.CertHash)
	if err != nil {
		return nil, fmt.Errorf("certHash: %w", err)
	}
	data, err := c.abi.Pack("issueCertificate", cid, ch)
	if err != nil {
		return nil, fmt.Errorf("pack issueCertificate: %w", err)
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
		cid, err := certIDToBytes32(item.CertID)
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
	data, err := c.abi.Pack("issueCertificateBatch", certIDs, certHashes)
	if err != nil {
		return nil, fmt.Errorf("pack issueCertificateBatch: %w", err)
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
	cid, err := certIDToBytes32(req.CertID)
	if err != nil {
		return nil, fmt.Errorf("certID: %w", err)
	}
	data, err := c.abi.Pack("revokeCertificate", cid, req.Reason)
	if err != nil {
		return nil, fmt.Errorf("pack revokeCertificate: %w", err)
	}
	return c.sendTx(ctx, data)
}

func (c *Contract) RestoreCertificate(ctx context.Context, req RestoreCertRequest) (*TxReceipt, error) {
	if err := c.loadABI(); err != nil {
		return nil, err
	}
	cid, err := certIDToBytes32(req.CertID)
	if err != nil {
		return nil, fmt.Errorf("certID: %w", err)
	}
	data, err := c.abi.Pack("restoreCertificate", cid)
	if err != nil {
		return nil, fmt.Errorf("pack restoreCertificate: %w", err)
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
	cid, err := certIDToBytes32(certID)
	if err != nil {
		return false, err
	}
	data, err := c.abi.Pack("certExists", cid)
	if err != nil {
		return false, fmt.Errorf("pack certExists: %w", err)
	}
	result, err := c.call(ctx, data)
	if err != nil {
		return false, err
	}
	var exists bool
	if err := c.abi.Unpack(&exists, "certExists", result); err != nil {
		return false, fmt.Errorf("unpack certExists: %w", err)
	}
	return exists, nil
}

func (c *Contract) GetCertificate(ctx context.Context, certID string) (*CertOnChainRecord, error) {
	if err := c.loadABI(); err != nil {
		return nil, err
	}
	cid, err := certIDToBytes32(certID)
	if err != nil {
		return nil, err
	}
	data, err := c.abi.Pack("getCertificate", cid)
	if err != nil {
		return nil, fmt.Errorf("pack getCertificate: %w", err)
	}
	result, err := c.call(ctx, data)
	if err != nil {
		return nil, err
	}
	// 合约返回: (bytes32 certHash, address issuer, uint64 issuedAt,
	//            bool revoked, uint64 revokedAt, string revokeReason)
	out := make(map[string]interface{})
	if err := c.abi.UnpackIntoMap(out, "getCertificate", result); err != nil {
		return nil, fmt.Errorf("unpack getCertificate: %w", err)
	}
	certHash32, _ := out["certHash"].([32]byte)
	issuerAddr, _ := out["issuer"].(common.Address)
	issuedAtU64, _ := out["issuedAt"].(uint64)
	revoked, _ := out["revoked"].(bool)
	revokedAtU64, _ := out["revokedAt"].(uint64)
	revokeReason, _ := out["revokeReason"].(string)
	rec := &CertOnChainRecord{
		CertHash:     hex.EncodeToString(certHash32[:]),
		Issuer:       issuerAddr.Hex(),
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
	cid, err := certIDToBytes32(certID)
	if err != nil {
		return nil, err
	}
	ch, err := hexToBytes32(certHash)
	if err != nil {
		return nil, err
	}
	data, err := c.abi.Pack("verifyCertificate", cid, ch)
	if err != nil {
		return nil, fmt.Errorf("pack verifyCertificate: %w", err)
	}
	result, err := c.call(ctx, data)
	if err != nil {
		return nil, err
	}
	out2 := make(map[string]interface{})
	if err := c.abi.UnpackIntoMap(out2, "verifyCertificate", result); err != nil {
		return nil, fmt.Errorf("unpack verifyCertificate: %w", err)
	}
	valid, _ := out2["valid"].(bool)
	revoked2, _ := out2["revoked"].(bool)
	return &VerifyResult{Valid: valid, Revoked: revoked2}, nil
}

func (c *Contract) VerifyCertificateBatch(ctx context.Context, req BatchIssueCertRequest) (*BatchVerifyResult, error) {
	if err := c.loadABI(); err != nil {
		return nil, err
	}
	certIDs := make([][32]byte, len(req.Items))
	certHashes := make([][32]byte, len(req.Items))
	for i, item := range req.Items {
		cid, err := certIDToBytes32(item.CertID)
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
	data, err := c.abi.Pack("verifyCertificateBatch", certIDs, certHashes)
	if err != nil {
		return nil, fmt.Errorf("pack verifyCertificateBatch: %w", err)
	}
	result, err := c.call(ctx, data)
	if err != nil {
		return nil, err
	}
	out3 := make(map[string]interface{})
	if err := c.abi.UnpackIntoMap(out3, "verifyCertificateBatch", result); err != nil {
		return nil, fmt.Errorf("unpack verifyCertificateBatch: %w", err)
	}
	valids, _ := out3["valids"].([]bool)
	revokeds, _ := out3["revokeds"].([]bool)
	batchOut := &BatchVerifyResult{Results: make([]VerifyResult, len(valids))}
	for i := range valids {
		batchOut.Results[i] = VerifyResult{Valid: valids[i], Revoked: revokeds[i]}
	}
	return batchOut, nil
}

// ----------------------------------------------------------------
// 颁发机构管理
// ----------------------------------------------------------------

func (c *Contract) AddIssuer(ctx context.Context, req AddIssuerRequest) (*TxReceipt, error) {
	if err := c.loadABI(); err != nil {
		return nil, err
	}
	data, err := c.abi.Pack("addIssuer", common.HexToAddress(req.Address), req.Name)
	if err != nil {
		return nil, fmt.Errorf("pack addIssuer: %w", err)
	}
	return c.sendTx(ctx, data)
}

func (c *Contract) GetIssuerInfo(ctx context.Context, address string) (*IssuerInfo, error) {
	if err := c.loadABI(); err != nil {
		return nil, err
	}
	data, err := c.abi.Pack("getIssuerInfo", common.HexToAddress(address))
	if err != nil {
		return nil, fmt.Errorf("pack getIssuerInfo: %w", err)
	}
	result, err := c.call(ctx, data)
	if err != nil {
		return nil, err
	}
	out4 := make(map[string]interface{})
	if err := c.abi.UnpackIntoMap(out4, "getIssuerInfo", result); err != nil {
		return nil, fmt.Errorf("unpack getIssuerInfo: %w", err)
	}
	authorized, _ := out4["authorized"].(bool)
	name, _ := out4["name"].(string)
	authorizedAt, _ := out4["authorizedAt"].(uint64)
	return &IssuerInfo{
		Address:      address,
		Authorized:   authorized,
		Name:         name,
		AuthorizedAt: time.Unix(int64(authorizedAt), 0),
	}, nil
}

func (c *Contract) GetStats(ctx context.Context) (*ContractStats, error) {
	if err := c.loadABI(); err != nil {
		return nil, err
	}
	data, err := c.abi.Pack("getStats")
	if err != nil {
		return nil, fmt.Errorf("pack getStats: %w", err)
	}
	result, err := c.call(ctx, data)
	if err != nil {
		return nil, err
	}
	out5 := make(map[string]interface{})
	if err := c.abi.UnpackIntoMap(out5, "getStats", result); err != nil {
		return nil, fmt.Errorf("unpack getStats: %w", err)
	}
	totalIssued, _ := out5["totalIssued"].(*big.Int)
	totalRevoked, _ := out5["totalRevoked"].(*big.Int)
	if totalIssued == nil {
		totalIssued = big.NewInt(0)
	}
	if totalRevoked == nil {
		totalRevoked = big.NewInt(0)
	}
	return &ContractStats{
		TotalIssued:  totalIssued.Uint64(),
		TotalRevoked: totalRevoked.Uint64(),
	}, nil
}

// ----------------------------------------------------------------
// 内部：发交易 / call
// ----------------------------------------------------------------

// sendTx 使用 go-sdk v3 推荐流程发送写操作交易：
// GetBlockNumber → CreateEncodedTransactionDataV1 → CreateEncodedSignature
// → CreateEncodedTransaction → SendEncodedTransaction
func (c *Contract) sendTx(ctx context.Context, input []byte) (*TxReceipt, error) {
	blockNumber, err := c.sdkClient.GetBlockNumber(ctx)
	if err != nil {
		return nil, fmt.Errorf("get block number: %w", err)
	}
	blockLimit := blockNumber + client.BlockLimit

	txData, txDataHash, err := c.sdkClient.CreateEncodedTransactionDataV1(
		&c.address, input, blockLimit, "",
	)
	if err != nil {
		return nil, fmt.Errorf("create tx data: %w", err)
	}
	signature, err := c.sdkClient.CreateEncodedSignature(txDataHash)
	if err != nil {
		return nil, fmt.Errorf("sign tx: %w", err)
	}
	encodedTx, err := c.sdkClient.CreateEncodedTransaction(
		txData, txDataHash, signature, 0, "",
	)
	if err != nil {
		return nil, fmt.Errorf("create encoded tx: %w", err)
	}
	receipt, err := c.sdkClient.SendEncodedTransaction(ctx, encodedTx, false)
	if err != nil {
		return nil, fmt.Errorf("send tx: %w", err)
	}
	return receiptToTxReceipt(receipt), nil
}

// call 发送只读调用（不上链，不消耗 Gas）。
func (c *Contract) call(ctx context.Context, input []byte) ([]byte, error) {
	msg := ethereum.CallMsg{
		To:   &c.address,
		Data: input,
	}
	result, err := c.sdkClient.CallContract(ctx, msg)
	if err != nil {
		return nil, fmt.Errorf("call contract: %w", err)
	}
	return result, nil
}

// receiptToTxReceipt 将 go-sdk v3 types.Receipt 转为项目内部 TxReceipt。
func receiptToTxReceipt(r *types.Receipt) *TxReceipt {
	if r == nil {
		return nil
	}
	return &TxReceipt{
		TxHash:      r.TransactionHash,
		BlockNumber: int64(r.BlockNumber),
		Status:      r.Status,
		Message:     r.Message,
	}
}

// certIDToBytes32 将链下 UUID 经 keccak256 转为合约 bytes32。
func certIDToBytes32(certID string) ([32]byte, error) {
	if certID == "" {
		return [32]byte{}, fmt.Errorf("certID is empty")
	}
	return CertIDToBytes32(certID), nil
}
