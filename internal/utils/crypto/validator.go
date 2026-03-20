package crypto

import (
	"regexp"
	"strings"
)

// Password strength levels
const (
	StrengthVeryWeak = iota
	StrengthWeak
	StrengthFair
	StrengthStrong
	StrengthVeryStrong
)

// PasswordValidator 密码验证器
type PasswordValidator struct {
	MinLength      int
	MaxLength      int
	RequireUpper   bool
	RequireLower   bool
	RequireDigit   bool
	RequireSpecial bool
	SpecialChars   string
}

// DefaultPasswordValidator 默认密码验证规则
var DefaultPasswordValidator = PasswordValidator{
	MinLength:      8,
	MaxLength:      32,
	RequireUpper:   true,
	RequireLower:   true,
	RequireDigit:   true,
	RequireSpecial: false,
	SpecialChars:   "!@#$%^&*()-_=+[]{}|;:,.<>?",
}

// ValidatePasswordStrength 验证密码强度
// 返回密码强度等级和错误信息（如果有）
func ValidatePasswordStrength(password string) (int, error) {
	if password == "" {
		return StrengthVeryWeak, ErrPasswordEmpty
	}

	length := len(password)
	if length < DefaultPasswordValidator.MinLength {
		return StrengthVeryWeak, ErrPasswordTooShort
	}
	if length > DefaultPasswordValidator.MaxLength {
		return StrengthVeryWeak, ErrPasswordTooLong
	}

	strength := StrengthVeryWeak
	hasUpper := false
	hasLower := false
	hasDigit := false
	hasSpecial := false

	for _, char := range password {
		switch {
		case char >= 'A' && char <= 'Z':
			hasUpper = true
		case char >= 'a' && char <= 'z':
			hasLower = true
		case char >= '0' && char <= '9':
			hasDigit = true
		case strings.Contains(DefaultPasswordValidator.SpecialChars, string(char)):
			hasSpecial = true
		}
	}

	// 计算强度
	if hasLower {
		strength++
	}
	if hasUpper {
		strength++
	}
	if hasDigit {
		strength++
	}
	if hasSpecial {
		strength += 2
	}
	if length >= 12 {
		strength++
	}
	if length >= 16 {
		strength++
	}

	// 验证必需条件
	if DefaultPasswordValidator.RequireUpper && !hasUpper {
		return strength, ErrPasswordNoUpperCase
	}
	if DefaultPasswordValidator.RequireLower && !hasLower {
		return strength, ErrPasswordNoLowerCase
	}
	if DefaultPasswordValidator.RequireDigit && !hasDigit {
		return strength, ErrPasswordNoDigit
	}
	if DefaultPasswordValidator.RequireSpecial && !hasSpecial {
		return strength, ErrPasswordNoSpecial
	}

	return strength, nil
}

// ValidatePassword 使用默认验证器验证密码
func ValidatePassword(password string) error {
	_, err := ValidatePasswordStrength(password)
	return err
}

// ValidatePasswordWithCustomRules 使用自定义规则验证密码
func ValidatePasswordWithCustomRules(password string, validator PasswordValidator) error {
	if password == "" {
		return ErrPasswordEmpty
	}

	length := len(password)
	if length < validator.MinLength {
		return ErrPasswordTooShort
	}
	if length > validator.MaxLength {
		return ErrPasswordTooLong
	}

	hasUpper := false
	hasLower := false
	hasDigit := false
	hasSpecial := false

	for _, char := range password {
		switch {
		case char >= 'A' && char <= 'Z':
			hasUpper = true
		case char >= 'a' && char <= 'z':
			hasLower = true
		case char >= '0' && char <= '9':
			hasDigit = true
		case strings.Contains(validator.SpecialChars, string(char)):
			hasSpecial = true
		}
	}

	if validator.RequireUpper && !hasUpper {
		return ErrPasswordNoUpperCase
	}
	if validator.RequireLower && !hasLower {
		return ErrPasswordNoLowerCase
	}
	if validator.RequireDigit && !hasDigit {
		return ErrPasswordNoDigit
	}
	if validator.RequireSpecial && !hasSpecial {
		return ErrPasswordNoSpecial
	}

	return nil
}

// GetPasswordStrengthText 获取密码强度文本描述
func GetPasswordStrengthText(strength int) string {
	switch strength {
	case StrengthVeryWeak:
		return "非常弱"
	case StrengthWeak:
		return "弱"
	case StrengthFair:
		return "一般"
	case StrengthStrong:
		return "强"
	case StrengthVeryStrong:
		return "非常强"
	default:
		return "未知"
	}
}

// IsPasswordHashValid 检查哈希值格式是否有效
func IsPasswordHashValid(hash string) bool {
	if hash == "" {
		return false
	}

	// bcrypt 哈希格式: $2a$或$2b$或$2y$ + 22字符盐值 + 31字符哈希
	// 总长度应为 60 字符
	if len(hash) != 60 {
		return false
	}

	// 验证格式
	pattern := regexp.MustCompile(`^\$2[aby]\$\d{2}\$[./A-Za-z0-9]{53}$`)
	return pattern.MatchString(hash)
}
