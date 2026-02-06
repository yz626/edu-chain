package middleware

import (
	"log"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
)

var (
	Enforcer *casbin.Enforcer
)

// InitCasbin 初始化 Casbin（通常在 main 或 server 初始化时调用）
func InitCasbin(modelPath, policyPath string) error {
	var err error

	Enforcer, err = casbin.NewEnforcer(modelPath, policyPath)
	if err != nil {
		log.Fatalf("Failed to create casbin enforcer: %v", err)
		return err
	}

	// 启用角色继承
	Enforcer.AddRoleForUser("admin", "SUPER_ADMIN")

	log.Println("Casbin initialized successfully")
	return nil
}

// GetEnforcer 获取 Casbin enforcer 实例
func GetEnforcer() *casbin.Enforcer {
	return Enforcer
}

// Enforce 检查权限
// sub: 角色
// obj: 资源
// act: 操作
func Enforce(sub, obj, act string) (bool, error) {
	return Enforcer.Enforce(sub, obj, act)
}

// AddPolicy 添加策略
func AddPolicy(sub, obj, act string) bool {
	ok, err := Enforcer.AddPolicy(sub, obj, act)
	if err != nil {
		log.Printf("Error adding policy: %v", err)
	}
	return ok
}

// RemovePolicy 删除策略
func RemovePolicy(sub, obj, act string) bool {
	ok, err := Enforcer.RemovePolicy(sub, obj, act)
	if err != nil {
		log.Printf("Error removing policy: %v", err)
	}
	return ok
}

// AddRoleForUser 为用户添加角色
func AddRoleForUser(user, role string) (bool, error) {
	return Enforcer.AddRoleForUser(user, role)
}

// RequirePermission 权限检查中间件
// obj: 资源名称（如 "certificate"、"user"）
// act: 操作（如 "read"、"create"、"update"、"delete"）
func RequirePermission(obj, act string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			// 如果没有设置角色，默认使用 STUDENT（最低权限）
			role = "STUDENT"
		}

		sub := role.(string)
		ok, err := Enforcer.Enforce(sub, obj, act)
		if err != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": "权限检查失败"})
			return
		}

		if !ok {
			c.AbortWithStatusJSON(403, gin.H{
				"error":    "权限不足",
				"role":     sub,
				"resource": obj,
				"action":   act,
			})
			return
		}

		c.Next()
	}
}

// RequireAnyPermission 满足任一权限即可通过的中间件
func RequireAnyPermission(obj string, acts ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			role = "STUDENT"
		}

		sub := role.(string)

		for _, act := range acts {
			ok, err := Enforcer.Enforce(sub, obj, act)
			if err != nil {
				c.AbortWithStatusJSON(500, gin.H{"error": "权限检查失败"})
				return
			}
			if ok {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(403, gin.H{
			"error":    "权限不足",
			"role":     sub,
			"resource": obj,
			"actions":  acts,
		})
	}
}

// RequireAllPermission 需要满足所有权限的中间件
func RequireAllPermission(obj string, acts ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			role = "STUDENT"
		}

		sub := role.(string)

		for _, act := range acts {
			ok, err := Enforcer.Enforce(sub, obj, act)
			if err != nil {
				c.AbortWithStatusJSON(500, gin.H{"error": "权限检查失败"})
				return
			}
			if !ok {
				c.AbortWithStatusJSON(403, gin.H{
					"error":    "权限不足",
					"role":     sub,
					"resource": obj,
					"action":   act,
				})
				return
			}
		}

		c.Next()
	}
}

// RequireSuperAdmin 需要超级管理员权限
func RequireSuperAdmin() gin.HandlerFunc {
	return RequirePermission("*", "*")
}
