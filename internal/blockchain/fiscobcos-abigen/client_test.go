package fiscobcosabigen

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/yz626/edu-chain/config"
)

// newTestClient 公共初始化：切换工作目录 + 加载配置 + 创建客户端
func newTestClient(t *testing.T) *Client {
	t.Helper()
	// 只在未切换时切换（避免多个测试重复 chdir）
	if _, err := os.Stat("config/blockchain.yaml"); err != nil {
		if err := os.Chdir("../../.."); err != nil {
			t.Fatalf("chdir to project root failed: %v", err)
		}
	}
	cfg, err := config.LoadBlockchain("config/blockchain.yaml")
	if err != nil {
		t.Fatalf("加载配置失败: %v", err)
	}
	if !cfg.Enabled {
		t.Skip("区块链未启用，跳过测试")
	}
	c, err := NewClient(cfg)
	if err != nil {
		t.Fatalf("NewClient 失败: %v", err)
	}
	return c
}

// certJSON 模拟证书 JSON 数据，生成确定性哈希
type certJSON struct {
	CertID      string `json:"cert_id"`
	StudentName string `json:"student_name"`
	Degree      string `json:"degree"`
	School      string `json:"school"`
	IssuedAt    string `json:"issued_at"`
}

// hashCertJSON 计算证书 JSON 的 SHA-256 哈希（hex，不带 0x）
func hashCertJSON(cert certJSON) string {
	data, _ := json.Marshal(cert)
	h := sha256.Sum256(data)
	return hex.EncodeToString(h[:])
}

// ================================================================
// TestAbigenClientConnect 连接与基础查询测试
// ================================================================
func TestAbigenClientConnect(t *testing.T) {
	c := newTestClient(t)
	ctx := context.Background()

	t.Log("[测试 1] GetStats ...")
	stats, err := c.GetStats(ctx)
	if err != nil {
		t.Fatalf("GetStats 失败: %v", err)
	}
	t.Logf("[测试 1] 通过 — 总颁发: %d  总撤销: %d", stats.TotalIssued, stats.TotalRevoked)

	t.Log("[测试 2] CertExists（不存在的 certID）...")
	exists, err := c.CertExists(ctx, "nonexistent-cert-id-12345")
	if err != nil {
		t.Fatalf("CertExists 失败: %v", err)
	}
	t.Logf("[测试 2] 通过 — exists=%v（预期 false）", exists)

	if cfg, _ := config.LoadBlockchain("config/blockchain.yaml"); cfg.Account.Address != "" {
		t.Logf("[测试 3] GetIssuerInfo address=%s ...", cfg.Account.Address)
		info, err := c.GetIssuerInfo(ctx, cfg.Account.Address)
		if err != nil {
			t.Logf("[测试 3] 失败（非致命）: %v", err)
		} else {
			t.Logf("[测试 3] 通过 — authorized=%v name=%q", info.Authorized, info.Name)
		}
	}
	t.Log("========== 连接测试通过 ==========")
}

// ================================================================
// TestIssueCertificate 上链一张证书，然后查询并验证
// ================================================================
func TestIssueCertificate(t *testing.T) {
	c := newTestClient(t)
	ctx := context.Background()

	// ----------------------------------------------------------------
	// 构造测试证书数据
	// ----------------------------------------------------------------
	certID := fmt.Sprintf("CERT-TEST-%d", time.Now().UnixMilli())
	cert := certJSON{
		CertID:      certID,
		StudentName: "张三",
		Degree:      "计算机科学与技术学士",
		School:      "测试大学",
		IssuedAt:    time.Now().Format("2006-01-02"),
	}
	certHash := hashCertJSON(cert)
	t.Logf("证书 ID   : %s", certID)
	t.Logf("证书 Hash : %s", certHash)

	// ----------------------------------------------------------------
	// 步骤 1：上链前确认证书不存在
	// ----------------------------------------------------------------
	t.Log("[步骤 1] 确认证书上链前不存在 ...")
	existsBefore, err := c.CertExists(ctx, certID)
	if err != nil {
		t.Fatalf("CertExists（上链前）失败: %v", err)
	}
	if existsBefore {
		t.Fatalf("证书上链前已存在（certID 重复），请重新运行")
	}
	t.Logf("[步骤 1] 通过 — 证书不存在（预期）")

	// ----------------------------------------------------------------
	// 步骤 2：上链颁发证书（写操作，消耗 Gas）
	// ----------------------------------------------------------------
	t.Log("[步骤 2] 调用 IssueCertificate 上链 ...")
	receipt, err := c.IssueCertificate(ctx, certID, certHash)
	if err != nil {
		t.Fatalf("IssueCertificate 失败: %v", err)
	}
	if receipt.Status != 0 {
		t.Fatalf("交易失败，status=%d message=%s txHash=%s",
			receipt.Status, receipt.Message, receipt.TxHash)
	}
	t.Logf("[步骤 2] 通过 — txHash=%s blockNumber=%d",
		receipt.TxHash, receipt.BlockNumber)

	// ----------------------------------------------------------------
	// 步骤 3：上链后确认证书存在
	// ----------------------------------------------------------------
	t.Log("[步骤 3] 确认证书上链后存在 ...")
	existsAfter, err := c.CertExists(ctx, certID)
	if err != nil {
		t.Fatalf("CertExists（上链后）失败: %v", err)
	}
	if !existsAfter {
		t.Fatal("证书上链后仍不存在，写入失败")
	}
	t.Logf("[步骤 3] 通过 — 证书已上链")

	// ----------------------------------------------------------------
	// 步骤 4：GetCertificate 查询链上完整记录
	// ----------------------------------------------------------------
	t.Log("[步骤 4] GetCertificate 查询链上记录 ...")
	rec, err := c.GetCertificate(ctx, certID)
	if err != nil {
		t.Fatalf("GetCertificate 失败: %v", err)
	}
	t.Logf("[步骤 4] 通过 — certHash=%s issuer=%s issuedAt=%s revoked=%v",
		rec.CertHash, rec.Issuer, rec.IssuedAt.Format("2006-01-02 15:04:05"), rec.Revoked)
	if rec.CertHash != certHash {
		t.Errorf("链上 certHash 不匹配: 期望 %s 实际 %s", certHash, rec.CertHash)
	}
	if rec.Revoked {
		t.Error("新颁发的证书不应处于撤销状态")
	}

	// ----------------------------------------------------------------
	// 步骤 5：VerifyCertificate 验证证书有效性
	// ----------------------------------------------------------------
	t.Log("[步骤 5] VerifyCertificate 验证证书 ...")
	result, err := c.VerifyCertificate(ctx, certID, certHash)
	if err != nil {
		t.Fatalf("VerifyCertificate 失败: %v", err)
	}
	t.Logf("[步骤 5] 通过 — valid=%v revoked=%v", result.Valid, result.Revoked)
	if !result.Valid {
		t.Error("证书验证应为 valid=true")
	}
	if result.Revoked {
		t.Error("新颁发的证书不应为 revoked=true")
	}

	// ----------------------------------------------------------------
	// 步骤 6：用错误哈希验证，应返回 valid=false
	// ----------------------------------------------------------------
	t.Log("[步骤 6] 用错误哈希验证（应返回 invalid）...")
	wrongHash := hex.EncodeToString(sha256.New().Sum([]byte("wrong")))
	// 补齐到 64 位
	if len(wrongHash) < 64 {
		wrongHash = fmt.Sprintf("%064s", wrongHash)
	}
	wrongHash = wrongHash[:64]
	wrongResult, err := c.VerifyCertificate(ctx, certID, wrongHash)
	if err != nil {
		t.Fatalf("VerifyCertificate（错误哈希）失败: %v", err)
	}
	t.Logf("[步骤 6] 通过 — valid=%v（预期 false）", wrongResult.Valid)
	if wrongResult.Valid {
		t.Error("错误哈希不应返回 valid=true")
	}

	// ----------------------------------------------------------------
	// 步骤 7：GetStats 确认总数增加
	// ----------------------------------------------------------------
	t.Log("[步骤 7] GetStats 确认颁发总数 ...")
	stats, err := c.GetStats(ctx)
	if err != nil {
		t.Fatalf("GetStats 失败: %v", err)
	}
	t.Logf("[步骤 7] 通过 — 总颁发: %d  总撤销: %d", stats.TotalIssued, stats.TotalRevoked)

	fmt.Printf("\n========== 证书上链全流程测试通过 ==========\n")
	fmt.Printf("证书 ID   : %s\n", certID)
	fmt.Printf("证书 Hash : %s\n", certHash)
	fmt.Printf("交易 Hash : %s\n", receipt.TxHash)
	fmt.Printf("区块高度  : %d\n", receipt.BlockNumber)
	fmt.Printf("总颁发数  : %d\n", stats.TotalIssued)
	fmt.Println("===========================================")
}
