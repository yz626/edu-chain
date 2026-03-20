package jwts

import (
	stdErrors "errors"
	"testing"

	errPkg "github.com/yz626/edu-chain/pkg/errors"
)

// TestErrorVariables 测试错误变量定义
func TestErrorVariables(t *testing.T) {
	// 测试 ErrUnauthorized
	if ErrUnauthorized == nil {
		t.Error("ErrUnauthorized 不应为空")
	}
	if ErrUnauthorized.GetCode() != errPkg.ErrCodeUnauthorized {
		t.Errorf("ErrUnauthorized 错误码不匹配: 期望 %s, 实际 %s", errPkg.ErrCodeUnauthorized, ErrUnauthorized.GetCode())
	}
	if ErrUnauthorized.GetMessage() == "" {
		t.Error("ErrUnauthorized 消息不应为空")
	}

	// 测试 ErrTokenExpired
	if ErrTokenExpired == nil {
		t.Error("ErrTokenExpired 不应为空")
	}
	if ErrTokenExpired.GetCode() != errPkg.ErrCodeTokenExpired {
		t.Errorf("ErrTokenExpired 错误码不匹配: 期望 %s, 实际 %s", errPkg.ErrCodeTokenExpired, ErrTokenExpired.GetCode())
	}
	if ErrTokenExpired.GetMessage() == "" {
		t.Error("ErrTokenExpired 消息不应为空")
	}

	// 测试 ErrTokenInvalid
	if ErrTokenInvalid == nil {
		t.Error("ErrTokenInvalid 不应为空")
	}
	if ErrTokenInvalid.GetCode() != errPkg.ErrCodeTokenInvalid {
		t.Errorf("ErrTokenInvalid 错误码不匹配: 期望 %s, 实际 %s", errPkg.ErrCodeTokenInvalid, ErrTokenInvalid.GetCode())
	}
	if ErrTokenInvalid.GetMessage() == "" {
		t.Error("ErrTokenInvalid 消息不应为空")
	}

	// 测试 ErrTokenMalformed
	if ErrTokenMalformed == nil {
		t.Error("ErrTokenMalformed 不应为空")
	}
	if ErrTokenMalformed.GetCode() != errPkg.ErrCodeTokenMalformed {
		t.Errorf("ErrTokenMalformed 错误码不匹配: 期望 %s, 实际 %s", errPkg.ErrCodeTokenMalformed, ErrTokenMalformed.GetCode())
	}
	if ErrTokenMalformed.GetMessage() == "" {
		t.Error("ErrTokenMalformed 消息不应为空")
	}

	// 测试 ErrTokenNotValidYet
	if ErrTokenNotValidYet == nil {
		t.Error("ErrTokenNotValidYet 不应为空")
	}
	if ErrTokenNotValidYet.GetCode() != errPkg.ErrCodeTokenNotValidYet {
		t.Errorf("ErrTokenNotValidYet 错误码不匹配: 期望 %s, 实际 %s", errPkg.ErrCodeTokenNotValidYet, ErrTokenNotValidYet.GetCode())
	}
	if ErrTokenNotValidYet.GetMessage() == "" {
		t.Error("ErrTokenNotValidYet 消息不应为空")
	}

	// 测试 ErrSigningMethod
	if ErrSigningMethod == nil {
		t.Error("ErrSigningMethod 不应为空")
	}
	if ErrSigningMethod.GetCode() != errPkg.ErrCodeSigningMethod {
		t.Errorf("ErrSigningMethod 错误码不匹配: 期望 %s, 实际 %s", errPkg.ErrCodeSigningMethod, ErrSigningMethod.GetCode())
	}
	if ErrSigningMethod.GetMessage() == "" {
		t.Error("ErrSigningMethod 消息不应为空")
	}
}

// TestErrorMessages 测试错误消息内容
func TestErrorMessages(t *testing.T) {
	tests := []struct {
		err      *errPkg.Error
		expected string
	}{
		{
			err:      ErrUnauthorized,
			expected: "未授权",
		},
		{
			err:      ErrTokenExpired,
			expected: "令牌已过期",
		},
		{
			err:      ErrTokenInvalid,
			expected: "令牌无效",
		},
		{
			err:      ErrTokenMalformed,
			expected: "令牌格式错误",
		},
		{
			err:      ErrTokenNotValidYet,
			expected: "令牌还未生效",
		},
		{
			err:      ErrSigningMethod,
			expected: "签名方法无效",
		},
	}

	for _, tt := range tests {
		if tt.err.GetMessage() != tt.expected {
			t.Errorf("错误消息不匹配: 期望 '%s', 实际 '%s'", tt.expected, tt.err.GetMessage())
		}
	}
}

// TestErrorCodes 测试错误码
func TestErrorCodes(t *testing.T) {
	tests := []struct {
		err      *errPkg.Error
		expected errPkg.ErrorCode
	}{
		{
			err:      ErrUnauthorized,
			expected: errPkg.ErrCodeUnauthorized,
		},
		{
			err:      ErrTokenExpired,
			expected: errPkg.ErrCodeTokenExpired,
		},
		{
			err:      ErrTokenInvalid,
			expected: errPkg.ErrCodeTokenInvalid,
		},
		{
			err:      ErrTokenMalformed,
			expected: errPkg.ErrCodeTokenMalformed,
		},
		{
			err:      ErrTokenNotValidYet,
			expected: errPkg.ErrCodeTokenNotValidYet,
		},
		{
			err:      ErrSigningMethod,
			expected: errPkg.ErrCodeSigningMethod,
		},
	}

	for _, tt := range tests {
		if tt.err.GetCode() != tt.expected {
			t.Errorf("错误码不匹配: 期望 '%s', 实际 '%s'", tt.expected, tt.err.GetCode())
		}
	}
}

// TestErrorIs 测试错误Is方法
func TestErrorIs(t *testing.T) {
	// 测试相同的错误
	if !ErrUnauthorized.Is(ErrUnauthorized) {
		t.Error("ErrUnauthorized 应该等于自身")
	}

	// 测试不同的错误（应该返回false）
	if ErrUnauthorized.Is(ErrTokenExpired) {
		t.Error("ErrUnauthorized 不应该等于 ErrTokenExpired")
	}
}

// TestErrorString 测试错误String方法
func TestErrorString(t *testing.T) {
	tests := []struct {
		err      *errPkg.Error
		checkSub string
	}{
		{
			err:      ErrUnauthorized,
			checkSub: "未授权",
		},
		{
			err:      ErrTokenExpired,
			checkSub: "令牌已过期",
		},
		{
			err:      ErrTokenInvalid,
			checkSub: "令牌无效",
		},
		{
			err:      ErrTokenMalformed,
			checkSub: "令牌格式错误",
		},
		{
			err:      ErrTokenNotValidYet,
			checkSub: "令牌还未生效",
		},
		{
			err:      ErrSigningMethod,
			checkSub: "签名方法无效",
		},
	}

	for _, tt := range tests {
		errStr := tt.err.Error()
		if errStr == "" {
			t.Errorf("错误 '%s' 的Error()方法返回空字符串", tt.checkSub)
		}
		// 错误字符串应该包含错误码或消息
		if len(errStr) < 5 {
			t.Errorf("错误 '%s' 的Error()方法返回字符串太短: %s", tt.checkSub, errStr)
		}
	}
}

// TestErrorUnwrap 测试错误Unwrap方法
func TestErrorUnwrap(t *testing.T) {
	// 预定义的错误没有包装原始错误，所以Unwrap返回nil
	if ErrUnauthorized.Unwrap() != nil {
		t.Error("ErrUnauthorized 的 Unwrap() 应该返回 nil")
	}

	// 测试带有原始错误的错误
	originalErr := stdErrors.New("original error")
	wrappedErr := ErrUnauthorized.WithError(originalErr)
	if wrappedErr.Unwrap() == nil {
		t.Error("带原始错误的 Unwrap() 不应该返回 nil")
	}
}

// TestErrorWithError 测试 WithError 方法
func TestErrorWithError(t *testing.T) {
	originalErr := stdErrors.New("original error")
	wrappedErr := ErrUnauthorized.WithError(originalErr)

	if wrappedErr.Unwrap() != originalErr {
		t.Error("WithError 应该设置原始错误")
	}

	// 测试错误消息
	if wrappedErr.GetMessage() != ErrUnauthorized.GetMessage() {
		t.Error("WithError 不应该改变错误消息")
	}
}

// TestErrorNew 测试创建自定义错误
func TestErrorNew(t *testing.T) {
	customErr := errPkg.New(errPkg.ErrCodeUnauthorized, "自定义未授权消息")

	if customErr.GetCode() != errPkg.ErrCodeUnauthorized {
		t.Error("自定义错误的错误码应该匹配")
	}

	if customErr.GetMessage() != "自定义未授权消息" {
		t.Error("自定义错误的消息应该匹配")
	}
}

// TestErrorNewf 测试带格式的错误创建
func TestErrorNewf(t *testing.T) {
	customErr := errPkg.Newf(errPkg.ErrCodeUnauthorized, "用户 %s 未授权", "testuser")

	if customErr.GetCode() != errPkg.ErrCodeUnauthorized {
		t.Error("带格式错误的错误码应该匹配")
	}

	expectedMsg := "用户 testuser 未授权"
	if customErr.GetMessage() != expectedMsg {
		t.Errorf("带格式错误的消息应该匹配: 期望 '%s', 实际 '%s'", expectedMsg, customErr.GetMessage())
	}
}

// TestErrorFromError 测试从标准错误转换
func TestErrorFromError(t *testing.T) {
	// 测试从nil转换
	result := errPkg.FromError(nil)
	if result != nil {
		t.Error("nil错误转换应该返回nil")
	}

	// 测试从自定义错误转换
	customErr := errPkg.New(errPkg.ErrCodeUnauthorized, "test")
	result = errPkg.FromError(customErr)
	if result.GetCode() != customErr.GetCode() {
		t.Error("从自定义错误转换应该保留错误码")
	}

	// 测试从标准错误转换
	stdErr := stdErrors.New("standard error")
	result = errPkg.FromError(stdErr)
	// 标准错误会转换为内部错误
	if result.GetCode() != errPkg.ErrCodeInternal {
		t.Error("从标准错误转换应该得到内部错误")
	}
}
