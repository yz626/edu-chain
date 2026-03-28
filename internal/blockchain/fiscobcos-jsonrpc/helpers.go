package fiscobcos

import (
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"strings"
)

// ----------------------------------------------------------------
// ABI 封装（纯 Go，支持合约所用类型）
// ----------------------------------------------------------------

type abiWrapper struct{ methods map[string]abiMethod }
type abiMethod struct {
	Name    string
	Inputs  []abiParam
	Outputs []abiParam
}
type abiParam struct {
	Name       string     `json:"name"`
	Type       string     `json:"type"`
	Components []abiParam `json:"components,omitempty"`
}
type abiEntry struct {
	Type    string     `json:"type"`
	Name    string     `json:"name"`
	Inputs  []abiParam `json:"inputs"`
	Outputs []abiParam `json:"outputs"`
}

func newABIWrapper(abiJSON string) (*abiWrapper, error) {
	var entries []abiEntry
	if err := json.Unmarshal([]byte(abiJSON), &entries); err != nil {
		return nil, fmt.Errorf("parse abi json: %w", err)
	}
	w := &abiWrapper{methods: make(map[string]abiMethod)}
	for _, e := range entries {
		if e.Type == "function" {
			w.methods[e.Name] = abiMethod{Name: e.Name, Inputs: e.Inputs, Outputs: e.Outputs}
		}
	}
	return w, nil
}

func (w *abiWrapper) pack(method string, args ...interface{}) ([]byte, error) {
	m, ok := w.methods[method]
	if !ok {
		return nil, fmt.Errorf("method %q not found in ABI", method)
	}
	sig := buildMethodSig(m)
	selector := keccak256Hash([]byte(sig))
	encoded, err := abiEncodeArgs(m.Inputs, args)
	if err != nil {
		return nil, fmt.Errorf("abi encode %s: %w", method, err)
	}
	result := make([]byte, 4+len(encoded))
	copy(result[:4], selector[:4])
	copy(result[4:], encoded)
	return result, nil
}

func (w *abiWrapper) unpack(method string, data []byte) ([]interface{}, error) {
	m, ok := w.methods[method]
	if !ok {
		return nil, fmt.Errorf("method %q not found in ABI", method)
	}
	return abiDecodeOutputs(m.Outputs, data)
}

func buildMethodSig(m abiMethod) string {
	params := make([]string, len(m.Inputs))
	for i, p := range m.Inputs {
		params[i] = p.Type
	}
	return fmt.Sprintf("%s(%s)", m.Name, strings.Join(params, ","))
}

// isDynamic 判断是否动态类型。
func isDynamic(t string) bool {
	return t == "string" || t == "bytes" || strings.HasSuffix(t, "[]")
}

// uint256Encode 将 int64 编码为 32 字节大端整数。
func uint256Encode(n int64) []byte {
	b := make([]byte, 32)
	big.NewInt(n).FillBytes(b)
	return b
}

// hexDecodeFixed 解码固定长度的 hex 字节。
func hexDecodeFixed(h string, length int) ([]byte, error) {
	if len(h) != length*2 {
		return nil, fmt.Errorf("expected %d hex chars, got %d", length*2, len(h))
	}
	b := make([]byte, length)
	for i := 0; i < length; i++ {
		var v byte
		fmt.Sscanf(h[i*2:i*2+2], "%02x", &v)
		b[i] = v
	}
	return b, nil
}

// ----------------------------------------------------------------
// ABI 编码
// ----------------------------------------------------------------

func abiEncodeArgs(params []abiParam, args []interface{}) ([]byte, error) {
	if len(params) != len(args) {
		return nil, fmt.Errorf("param count mismatch: expected %d got %d", len(params), len(args))
	}
	if len(params) == 0 {
		return []byte{}, nil
	}
	heads := make([][]byte, len(params))
	var tails [][]byte
	headSize := 32 * len(params)
	tailOffset := headSize
	for i, p := range params {
		if isDynamic(p.Type) {
			heads[i] = uint256Encode(int64(tailOffset))
			tail, err := encodeValue(p.Type, args[i])
			if err != nil {
				return nil, fmt.Errorf("param[%d] %s: %w", i, p.Name, err)
			}
			tails = append(tails, tail)
			tailOffset += len(tail)
		} else {
			enc, err := encodeValue(p.Type, args[i])
			if err != nil {
				return nil, fmt.Errorf("param[%d] %s: %w", i, p.Name, err)
			}
			heads[i] = enc
		}
	}
	var result []byte
	for _, h := range heads {
		result = append(result, h...)
	}
	for _, t := range tails {
		result = append(result, t...)
	}
	return result, nil
}

func encodeValue(t string, v interface{}) ([]byte, error) {
	switch t {
	case "bytes32":
		b, ok := v.([32]byte)
		if !ok {
			return nil, fmt.Errorf("expected [32]byte for bytes32, got %T", v)
		}
		result := make([]byte, 32)
		copy(result, b[:])
		return result, nil
	case "bytes32[]":
		arr, ok := v.([][32]byte)
		if !ok {
			return nil, fmt.Errorf("expected [][32]byte for bytes32[]")
		}
		result := uint256Encode(int64(len(arr)))
		for _, item := range arr {
			result = append(result, item[:]...)
		}
		return result, nil
	case "address":
		s, ok := v.(string)
		if !ok {
			return nil, fmt.Errorf("expected string for address")
		}
		s = strings.TrimPrefix(s, "0x")
		for len(s) < 40 {
			s = "0" + s
		}
		b, err := hexDecodeFixed(s, 20)
		if err != nil {
			return nil, err
		}
		result := make([]byte, 32)
		copy(result[12:], b)
		return result, nil
	case "bool":
		b, ok := v.(bool)
		if !ok {
			return nil, fmt.Errorf("expected bool")
		}
		result := make([]byte, 32)
		if b {
			result[31] = 1
		}
		return result, nil
	case "string":
		s, ok := v.(string)
		if !ok {
			return nil, fmt.Errorf("expected string")
		}
		b := []byte(s)
		result := uint256Encode(int64(len(b)))
		padded := make([]byte, (len(b)+31)/32*32)
		copy(padded, b)
		return append(result, padded...), nil
	case "uint256":
		switch n := v.(type) {
		case *big.Int:
			b := make([]byte, 32)
			n.FillBytes(b)
			return b, nil
		case int64:
			return uint256Encode(n), nil
		case uint64:
			return uint256Encode(int64(n)), nil
		default:
			return nil, fmt.Errorf("unsupported uint256 value type %T", v)
		}
	case "uint64":
		var n uint64
		switch val := v.(type) {
		case uint64:
			n = val
		case int64:
			n = uint64(val)
		default:
			return nil, fmt.Errorf("unsupported uint64 value type %T", v)
		}
		result := make([]byte, 32)
		big.NewInt(0).SetUint64(n).FillBytes(result)
		return result, nil
	default:
		return nil, fmt.Errorf("unsupported abi type %q", t)
	}
}

// ----------------------------------------------------------------
// ABI 解码
// ----------------------------------------------------------------

func abiDecodeOutputs(params []abiParam, data []byte) ([]interface{}, error) {
	if len(params) == 0 {
		return nil, nil
	}
	results := make([]interface{}, len(params))
	for i, p := range params {
		base := i * 32
		if base+32 > len(data) {
			return nil, fmt.Errorf("decode output[%d]: data too short", i)
		}
		if isDynamic(p.Type) {
			off := int(new(big.Int).SetBytes(data[base : base+32]).Int64())
			v, err := decodeValue(p.Type, data, off)
			if err != nil {
				return nil, fmt.Errorf("decode output[%d] %s: %w", i, p.Name, err)
			}
			results[i] = v
		} else {
			v, err := decodeValue(p.Type, data[base:base+32], 0)
			if err != nil {
				return nil, fmt.Errorf("decode output[%d] %s: %w", i, p.Name, err)
			}
			results[i] = v
		}
	}
	return results, nil
}

func decodeValue(t string, data []byte, offset int) (interface{}, error) {
	switch t {
	case "bytes32":
		if len(data) < 32 {
			return nil, fmt.Errorf("bytes32: data too short")
		}
		var b [32]byte
		copy(b[:], data[:32])
		return b, nil
	case "address":
		if len(data) < 32 {
			return nil, fmt.Errorf("address: data too short")
		}
		// 地址在后 20 字节
		return "0x" + fmt.Sprintf("%x", data[12:32]), nil
	case "bool":
		if len(data) < 32 {
			return nil, fmt.Errorf("bool: data too short")
		}
		return data[31] == 1, nil
	case "uint64":
		if len(data) < 32 {
			return nil, fmt.Errorf("uint64: data too short")
		}
		return new(big.Int).SetBytes(data[24:32]).Uint64(), nil
	case "uint256":
		if len(data) < 32 {
			return nil, fmt.Errorf("uint256: data too short")
		}
		return new(big.Int).SetBytes(data[:32]), nil
	case "string":
		if offset+32 > len(data) {
			return nil, fmt.Errorf("string: offset out of bounds")
		}
		length := int(new(big.Int).SetBytes(data[offset : offset+32]).Int64())
		start := offset + 32
		if start+length > len(data) {
			return nil, fmt.Errorf("string: data too short")
		}
		return string(data[start : start+length]), nil
	case "bool[]":
		if offset+32 > len(data) {
			return nil, fmt.Errorf("bool[]: offset out of bounds")
		}
		count := int(new(big.Int).SetBytes(data[offset : offset+32]).Int64())
		result := make([]bool, count)
		for i := 0; i < count; i++ {
			start := offset + 32 + i*32
			if start+32 > len(data) {
				return nil, fmt.Errorf("bool[]: element %d out of bounds", i)
			}
			result[i] = data[start+31] == 1
		}
		return result, nil
	default:
		return nil, fmt.Errorf("unsupported decode type %q", t)
	}
}

// ----------------------------------------------------------------
// 文件工具
// ----------------------------------------------------------------

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
