package jwts

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/yz626/edu-chain/pkg/constants"
	"github.com/yz626/edu-chain/pkg/response"
)

// ContextKey JWT上下文键
const (
	ContextKeyUserID   = "user_id"
	ContextKeyUsername = "username"
	ContextKeyEmail    = "email"
	ContextKeyUserType = "user_type"
	ContextKeyStatus   = "status"
	ContextKeyClaims   = "jwt_claims"
)

// AuthMiddleware 创建JWT认证中间件
// 从Authorization header中提取Bearer token并验证
func (j *JWT) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 提取Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, "Authorization header is required")
			c.Abort()
			return
		}

		// 提取Bearer token
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Unauthorized(c, "Invalid authorization header format")
			c.Abort()
			return
		}

		tokenString := parts[1]

		// 解析和验证token
		claims, err := j.ParseToken(tokenString)
		if err != nil {
			response.Unauthorized(c, "Invalid or expired token")
			c.Abort()
			return
		}

		// 检查用户状态
		if claims.Status != 1 { // 1 = 正常状态
			response.Forbidden(c, "User account is disabled or locked")
			c.Abort()
			return
		}

		// 将用户信息存入上下文
		c.Set(ContextKeyUserID, claims.UserID)
		c.Set(ContextKeyUsername, claims.Username)
		c.Set(ContextKeyEmail, claims.Email)
		c.Set(ContextKeyUserType, claims.UserType)
		c.Set(ContextKeyStatus, claims.Status)
		c.Set(ContextKeyClaims, claims)

		c.Next()
	}
}

// OptionalAuthMiddleware 可选的JWT认证中间件
// 如果存在token则解析并设置用户信息，否则继续处理
func (j *JWT) OptionalAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.Next()
			return
		}

		tokenString := parts[1]
		claims, err := j.ParseToken(tokenString)
		if err != nil {
			c.Next()
			return
		}

		// 检查用户状态
		if claims.Status == 1 {
			c.Set(ContextKeyUserID, claims.UserID)
			c.Set(ContextKeyUsername, claims.Username)
			c.Set(ContextKeyEmail, claims.Email)
			c.Set(ContextKeyUserType, claims.UserType)
			c.Set(ContextKeyStatus, claims.Status)
			c.Set(ContextKeyClaims, claims)
		}

		c.Next()
	}
}

// GetUserID 从上下文获取用户ID
func GetUserID(c *gin.Context) string {
	if v, exists := c.Get(ContextKeyUserID); exists {
		return v.(string)
	}
	return ""
}

// GetUsername 从上下文获取用户名
func GetUsername(c *gin.Context) string {
	if v, exists := c.Get(ContextKeyUsername); exists {
		return v.(string)
	}
	return ""
}

// GetEmail 从上下文获取邮箱
func GetEmail(c *gin.Context) string {
	if v, exists := c.Get(ContextKeyEmail); exists {
		return v.(string)
	}
	return ""
}

// GetUserType 从上下文获取用户类型
func GetUserType(c *gin.Context) int32 {
	if v, exists := c.Get(ContextKeyUserType); exists {
		return v.(int32)
	}
	return 0
}

// GetStatus 从上下文获取用户状态
func GetStatus(c *gin.Context) int32 {
	if v, exists := c.Get(ContextKeyStatus); exists {
		return v.(int32)
	}
	return 0
}

// GetClaims 从上下文获取完整Claims
func GetClaims(c *gin.Context) *Claims {
	if v, exists := c.Get(ContextKeyClaims); exists {
		return v.(*Claims)
	}
	return nil
}

// IsAdmin 检查是否是管理员
func IsAdmin(c *gin.Context) bool {
	return GetUserType(c) == constants.UserTypeAdmin
}
