package fiscobcos

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// ================================================================
// FISCO BCOS 3.0 JSON-RPC 客户端（纯 Go，无 CGO 依赖）
//
// FISCO BCOS 3.0 节点默认在 20200 端口提供 JSON-RPC over HTTP 服务。
// 文档：https://fisco-bcos-doc.readthedocs.io/zh-cn/latest/docs/develop/api.html
// ================================================================

// rpcRequest JSON-RPC 2.0 请求结构
type rpcRequest struct {
	JSONRPC string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
	ID      int           `json:"id"`
}

// rpcResponse JSON-RPC 2.0 响应结构
type rpcResponse struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      int             `json:"id"`
	Result  json.RawMessage `json:"result"`
	Error   *rpcError       `json:"error"`
}

// rpcError JSON-RPC 错误
type rpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *rpcError) Error() string {
	return fmt.Sprintf("rpc error %d: %s", e.Code, e.Message)
}

// rpcClient 封装对 FISCO BCOS 3.0 节点的 JSON-RPC 调用
type rpcClient struct {
	endpoint string       // http://host:port
	groupID  string       // group0
	httpClient *http.Client
	idCounter  int
}

// newRPCClient 创建 JSON-RPC 客户端
func newRPCClient(endpoint, groupID string, timeout time.Duration) *rpcClient {
	if timeout <= 0 {
		timeout = 30 * time.Second
	}
	return &rpcClient{
		endpoint: strings.TrimRight(endpoint, "/"),
		groupID:  groupID,
		httpClient: &http.Client{Timeout: timeout},
	}
}

// call 执行 JSON-RPC 调用
func (c *rpcClient) call(ctx context.Context, method string, params []interface{}, result interface{}) error {
	c.idCounter++
	req := rpcRequest{
		JSONRPC: "2.0",
		Method:  method,
		Params:  params,
		ID:      c.idCounter,
	}

	body, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("marshal rpc request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.endpoint, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("create http request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("http post to %s: %w", c.endpoint, err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response body: %w", err)
	}

	var rpcResp rpcResponse
	if err := json.Unmarshal(respBody, &rpcResp); err != nil {
		return fmt.Errorf("unmarshal rpc response: %w", err)
	}
	if rpcResp.Error != nil {
		return rpcResp.Error
	}
	if result != nil && rpcResp.Result != nil {
		if err := json.Unmarshal(rpcResp.Result, result); err != nil {
			return fmt.Errorf("unmarshal rpc result: %w", err)
		}
	}
	return nil
}

// getBlockNumber 获取当前区块高度
func (c *rpcClient) getBlockNumber(ctx context.Context) (int64, error) {
	var result string
	if err := c.call(ctx, "getBlockNumber", []interface{}{c.groupID}, &result); err != nil {
		return 0, fmt.Errorf("getBlockNumber: %w", err)
	}
	var num int64
	fmt.Sscanf(result, "%d", &num)
	return num, nil
}

// sendRawTransaction 发送已签名的原始交易
func (c *rpcClient) sendRawTransaction(ctx context.Context, rawTx string) (*txReceiptRaw, error) {
	var result txReceiptRaw
	if err := c.call(ctx, "sendTransaction",
		[]interface{}{c.groupID, rawTx, false}, &result); err != nil {
		return nil, fmt.Errorf("sendTransaction: %w", err)
	}
	return &result, nil
}

// callContract 只读合约调用（call，不上链）
func (c *rpcClient) callContract(ctx context.Context, to, data string) (string, error) {
	params := map[string]string{
		"to":   to,
		"data": data,
	}
	var result string
	if err := c.call(ctx, "call",
		[]interface{}{c.groupID, params}, &result); err != nil {
		return "", fmt.Errorf("call: %w", err)
	}
	return result, nil
}

// txReceiptRaw JSON-RPC 返回的交易回执原始格式
type txReceiptRaw struct {
	TransactionHash string `json:"transactionHash"`
	BlockNumber     string `json:"blockNumber"`
	Status          string `json:"status"`
	Output          string `json:"output"`
	Message         string `json:"message"`
}
