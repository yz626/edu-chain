package crypto

import "github.com/yz626/edu-chain/pkg/errors"

// Password validation errors - 使用统一的错误定义
var (
	ErrPasswordEmpty       = errors.New(errors.ErrCodePasswordEmpty, "密码不能为空")
	ErrPasswordTooShort    = errors.New(errors.ErrCodePasswordTooShort, "密码长度至少为8个字符")
	ErrPasswordTooLong     = errors.New(errors.ErrCodePasswordTooLong, "密码长度不能超过32个字符")
	ErrPasswordNoUpperCase = errors.New(errors.ErrCodePasswordNoUpperCase, "密码必须包含至少一个大写字母")
	ErrPasswordNoLowerCase = errors.New(errors.ErrCodePasswordNoLowerCase, "密码必须包含至少一个小写字母")
	ErrPasswordNoDigit     = errors.New(errors.ErrCodePasswordNoDigit, "密码必须包含至少一个数字")
	ErrPasswordNoSpecial   = errors.New(errors.ErrCodePasswordNoSpecial, "密码必须包含至少一个特殊字符")
)
