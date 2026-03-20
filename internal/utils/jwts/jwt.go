package jwts

import (
	"crypto/rand"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/yz626/edu-chain/config"
)

// Claims JWT Claims结构体
// 携带核心用户信息，避免过于复杂
type Claims struct {
	jwt.RegisteredClaims        // 继承标准Claims（包含Issuer、Subject、Audience、Expiry等）
	UserID               string `json:"user_id"`   // 用户ID
	Username             string `json:"username"`  // 用户名
	Email                string `json:"email"`     // 邮箱
	UserType             int32  `json:"user_type"` // 用户类型
	Status               int32  `json:"status"`    // 用户状态
}

// UserClaims token 用户信息
type UserClaims struct {
	UserID   string `json:"user_id"`   // 用户ID
	Username string `json:"username"`  // 用户名
	Email    string `json:"email"`     // 邮箱
	UserType int32  `json:"user_type"` // 用户类型
	Status   int32  `json:"status"`    // 用户状态
}

// JWT JWT工具类
type JWT struct {
	cfg *config.JWTConfig
}

// New 创建JWT工具实例
func New(cfg *config.JWTConfig) *JWT {
	return &JWT{cfg: cfg}
}

// GenerateToken 生成访问令牌
// 携带必要用户信息：用户ID、用户名、邮箱、用户类型、状态
func (j *JWT) GenerateToken(user *UserClaims) (string, error) {
	now := time.Now()
	expiresAt := now.Add(time.Duration(j.cfg.Expire) * time.Second)

	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt), // 过期时间
			IssuedAt:  jwt.NewNumericDate(now),       // 签发时间
			NotBefore: jwt.NewNumericDate(now),       // 生效时间
			Issuer:    j.cfg.Issuer,                  // 签发者
			ID:        generateTokenID(),             // 唯一标识
		},
		UserID:   user.UserID,
		Username: user.Username,
		Email:    user.Email,
		UserType: user.UserType,
		Status:   user.Status,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.cfg.Secret))
}

// GenerateRefreshToken 生成刷新令牌
func (j *JWT) GenerateRefreshToken(userID string) (string, error) {
	now := time.Now()
	expiresAt := now.Add(time.Duration(j.cfg.RefreshExpire) * time.Second)

	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    j.cfg.Issuer,
			ID:        generateTokenID(),
			Subject:   userID, // 刷新令牌的主体是用户ID
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.cfg.Secret))
}

// ParseToken 解析和验证令牌
func (j *JWT) ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrSigningMethod
		}
		return []byte(j.cfg.Secret), nil
	})

	if err != nil {
		return nil, handleJWTError(err)
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrTokenInvalid
}

// ValidateToken 验证令牌是否有效（仅验证，不解析）
func (j *JWT) ValidateToken(tokenString string) bool {
	_, err := j.ParseToken(tokenString)
	return err == nil
}

// GetUserIDFromToken 从令牌中获取用户ID
func (j *JWT) GetUserIDFromToken(tokenString string) (string, error) {
	claims, err := j.ParseToken(tokenString)
	if err != nil {
		return "", err
	}
	return claims.UserID, nil
}

// handleJWTError 处理JWT错误，转换为项目统一的错误类型
func handleJWTError(jwtErr error) error {
	if jwtErr == nil {
		return nil
	}

	// 使用 errors.Is 来检查错误类型
	if jwtErr == jwt.ErrTokenExpired {
		return ErrTokenExpired
	}
	if jwtErr == jwt.ErrTokenNotValidYet {
		return ErrTokenNotValidYet
	}
	if jwtErr == jwt.ErrTokenMalformed {
		return ErrTokenMalformed
	}

	// 其他错误统一为令牌无效
	return ErrTokenInvalid
}

// generateTokenID 生成唯一的令牌ID
func generateTokenID() string {
	return time.Now().Format("20060102150405.000000") + "-" + randomString(8)
}

// randomString 生成随机字符串
func randomString(length int) string {
	b := make([]byte, length)
	_, _ = rand.Read(b)
	return fmt.Sprintf("%x", b)[:length]
}
