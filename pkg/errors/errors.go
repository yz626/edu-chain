package errors

import (
	"fmt"
)

// ErrorCode 错误码类型
type ErrorCode string

// 错误码定义
const (
	// 通用错误 (COMMON_xxx)
	ErrCodeUnknown      ErrorCode = "COMMON_000"
	ErrCodeInternal     ErrorCode = "COMMON_001"
	ErrCodeInvalidParam ErrorCode = "COMMON_002"

	// 用户相关错误 (USER_xxx)
	ErrCodeUserNotFound      ErrorCode = "USER_001"
	ErrCodeUserAlreadyExists ErrorCode = "USER_002"
	ErrCodeInvalidPassword   ErrorCode = "USER_003"
	ErrCodeUserDisabled      ErrorCode = "USER_004"

	// 认证相关错误 (AUTH_xxx)
	ErrCodeUnauthorized     ErrorCode = "AUTH_001"
	ErrCodeTokenExpired     ErrorCode = "AUTH_002"
	ErrCodeTokenInvalid     ErrorCode = "AUTH_003"
	ErrCodeTokenMalformed   ErrorCode = "AUTH_004"
	ErrCodeTokenNotValidYet ErrorCode = "AUTH_005"
	ErrCodeSigningMethod    ErrorCode = "AUTH_006"

	// 密码相关错误 (PASSWORD_xxx)
	ErrCodePasswordEmpty       ErrorCode = "PASSWORD_001"
	ErrCodePasswordTooShort    ErrorCode = "PASSWORD_002"
	ErrCodePasswordTooLong     ErrorCode = "PASSWORD_003"
	ErrCodePasswordWeak        ErrorCode = "PASSWORD_004"
	ErrCodePasswordNoUpperCase ErrorCode = "PASSWORD_005"
	ErrCodePasswordNoLowerCase ErrorCode = "PASSWORD_006"
	ErrCodePasswordNoDigit     ErrorCode = "PASSWORD_007"
	ErrCodePasswordNoSpecial   ErrorCode = "PASSWORD_008"

	// 证书相关错误 (CERT_xxx)
	ErrCodeCertNotFound ErrorCode = "CERT_001"
	ErrCodeCertInvalid  ErrorCode = "CERT_002"
	ErrCodeCertExpired  ErrorCode = "CERT_003"
	ErrCodeCertRevoked  ErrorCode = "CERT_004"

	// 区块链相关错误 (CHAIN_xxx)
	ErrCodeChainConnectionFailed  ErrorCode = "CHAIN_001"
	ErrCodeChainTransactionFailed ErrorCode = "CHAIN_002"
	ErrCodeChainContractError     ErrorCode = "CHAIN_003"

	// 权限相关错误 (PERM_xxx)
	ErrCodeForbidden        ErrorCode = "PERM_001"
	ErrCodePermissionDenied ErrorCode = "PERM_002"
)

// Error 错误结构
type Error struct {
	code    ErrorCode `json:"code"`    // 错误码
	message string    `json:"message"` // 错误消息
	err     error     `json:"-"`       // 原始错误
}

// Error 实现 error 接口
func (e *Error) Error() string {
	if e.err != nil {
		return fmt.Sprintf("[%s] %s: %v", e.code, e.message, e.err)
	}
	return fmt.Sprintf("[%s] %s", e.code, e.message)
}

// Unwrap 返回原始错误
func (e *Error) Unwrap() error {
	return e.err
}

// WithError 添加原始错误
func (e *Error) WithError(err error) *Error {
	e.err = err
	return e
}

// New 创建新错误
func New(code ErrorCode, message string) *Error {
	return &Error{
		code:    code,
		message: message,
	}
}

// Newf 创建带格式化的错误
func Newf(code ErrorCode, format string, args ...interface{}) *Error {
	return &Error{
		code:    code,
		message: fmt.Sprintf(format, args...),
	}
}

// FromError 从标准错误转换
func FromError(err error) *Error {
	if err == nil {
		return nil
	}

	// 如果已经是 *Error 类型，直接返回
	if e, ok := err.(*Error); ok {
		return e
	}

	// 默认转换为内部错误
	return &Error{
		code:    ErrCodeInternal,
		message: err.Error(),
		err:     err,
	}
}

// Is 判断错误是否相等
func (e *Error) Is(target error) bool {
	if target == nil {
		return false
	}

	// 检查错误码是否匹配
	if t, ok := target.(*Error); ok {
		return e.code == t.code
	}

	return false
}

// GetCode 获取错误码
func (e *Error) GetCode() ErrorCode {
	return e.code
}

// GetMessage 获取错误消息
func (e *Error) GetMessage() string {
	return e.message
}

// SetMessage 设置错误消息
func (e *Error) SetMessage(msg string) *Error {
	e.message = msg
	return e
}

// 预定义的错误变量
var (
	ErrUnknown      = New(ErrCodeUnknown, "未知错误")
	ErrInternal     = New(ErrCodeInternal, "内部服务器错误")
	ErrInvalidParam = New(ErrCodeInvalidParam, "参数无效")
)
