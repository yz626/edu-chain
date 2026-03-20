package crypto

import (
	"crypto/rand"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// Password hashing cost (14-31, higher = more secure but slower)
const (
	DefaultCost = 14
	MinCost     = 4
	MaxCost     = 31
)

// HashPassword 使用 bcrypt 加密密码
// 返回加密后的哈希值
func HashPassword(password string) (string, error) {
	if password == "" {
		return "", ErrPasswordEmpty
	}

	// 使用默认成本生成哈希
	hash, err := bcrypt.GenerateFromPassword([]byte(password), DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

// HashPasswordWithCost 使用指定成本加密密码
// cost 范围: 4-31, 建议使用 14
func HashPasswordWithCost(password string, cost int) (string, error) {
	if password == "" {
		return "", ErrPasswordEmpty
	}

	// 验证成本范围
	if cost < MinCost || cost > MaxCost {
		cost = DefaultCost
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

// CheckPassword 验证密码是否与哈希值匹配
func CheckPassword(password string, hash string) (bool, error) {
	if password == "" || hash == "" {
		return false, ErrPasswordEmpty
	}

	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

// NeedsRehash 检查密码哈希是否需要重新哈希
// 例如当 cost 已经低于当前推荐值时
func NeedsRehash(hash string, minCost int) (bool, error) {
	if hash == "" {
		return false, ErrPasswordEmpty
	}

	// 解析哈希成本
	// 格式: $2a$XX$...
	// XX 是两位数成本
	if len(hash) < 7 {
		return false, errors.New("无效的哈希格式")
	}

	var cost int
	_, err := fmt.Sscanf(hash[4:6], "%d", &cost)
	if err != nil {
		return false, err
	}

	return cost < minCost, nil
}

// GenerateRandomPassword 生成随机密码
// length: 密码长度
// includeSpecial: 是否包含特殊字符
func GenerateRandomPassword(length int, includeSpecial bool) (string, error) {
	if length < 6 || length > 64 {
		length = 16 // 默认长度
	}

	chars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	if includeSpecial {
		chars += "!@#$%^&*()-_=+"
	}

	password := make([]byte, length)
	randBytes := make([]byte, length)
	_, err := rand.Read(randBytes)
	if err != nil {
		return "", err
	}

	for i := range password {
		password[i] = chars[int(randBytes[i])%len(chars)]
	}

	return string(password), nil
}
