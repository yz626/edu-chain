package jwts

import "github.com/yz626/edu-chain/pkg/errors"

var (
	ErrUnauthorized     = errors.New(errors.ErrCodeUnauthorized, "未授权")
	ErrTokenExpired     = errors.New(errors.ErrCodeTokenExpired, "令牌已过期")
	ErrTokenInvalid     = errors.New(errors.ErrCodeTokenInvalid, "令牌无效")
	ErrTokenMalformed   = errors.New(errors.ErrCodeTokenMalformed, "令牌格式错误")
	ErrTokenNotValidYet = errors.New(errors.ErrCodeTokenNotValidYet, "令牌还未生效")
	ErrSigningMethod    = errors.New(errors.ErrCodeSigningMethod, "签名方法无效")
)
