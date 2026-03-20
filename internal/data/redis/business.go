package cache

import (
	"context"
	"fmt"
	"time"
)

// ===== 认证/会话相关 =====

// SetToken 设置Token缓存
func (r *RedisClient) SetToken(ctx context.Context, userID string, token string, expiration time.Duration) error {
	key := PrefixToken + userID
	return r.Set(ctx, key, token, expiration)
}

// GetToken 获取Token缓存
func (r *RedisClient) GetToken(ctx context.Context, userID string) (string, error) {
	key := PrefixToken + userID
	return r.Get(ctx, key)
}

// DeleteToken 删除Token缓存
func (r *RedisClient) DeleteToken(ctx context.Context, userID string) error {
	key := PrefixToken + userID
	return r.Del(ctx, key)
}

// SetSession 设置会话
func (r *RedisClient) SetSession(ctx context.Context, sessionID string, value interface{}, expiration time.Duration) error {
	key := PrefixSession + sessionID
	return r.SetJSON(ctx, key, value, expiration)
}

// GetSession 获取会话
func (r *RedisClient) GetSession(ctx context.Context, sessionID string, dest interface{}) error {
	key := PrefixSession + sessionID
	return r.GetJSON(ctx, key, dest)
}

// DeleteSession 删除会话
func (r *RedisClient) DeleteSession(ctx context.Context, sessionID string) error {
	key := PrefixSession + sessionID
	return r.Del(ctx, key)
}

// SetRefreshToken 设置刷新Token
func (r *RedisClient) SetRefreshToken(ctx context.Context, userID string, token string, expiration time.Duration) error {
	key := PrefixRefresh + userID
	return r.Set(ctx, key, token, expiration)
}

// GetRefreshToken 获取刷新Token
func (r *RedisClient) GetRefreshToken(ctx context.Context, userID string) (string, error) {
	key := PrefixRefresh + userID
	return r.Get(ctx, key)
}

// DeleteRefreshToken 删除刷新Token
func (r *RedisClient) DeleteRefreshToken(ctx context.Context, userID string) error {
	key := PrefixRefresh + userID
	return r.Del(ctx, key)
}

// ===== 用户权限缓存 =====

// CacheUserPermissions 缓存用户权限
func (r *RedisClient) CacheUserPermissions(ctx context.Context, userID string, permissions []string, expiration time.Duration) error {
	key := PrefixUserPerm + userID
	if err := r.client.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("delete key failed: %w", err)
	}
	if err := r.client.SAdd(ctx, key, permissions).Err(); err != nil {
		return fmt.Errorf("add permissions failed: %w", err)
	}
	return r.client.Expire(ctx, key, expiration).Err()
}

// GetUserPermissions 获取用户权限缓存
func (r *RedisClient) GetUserPermissions(ctx context.Context, userID string) ([]string, error) {
	key := PrefixUserPerm + userID
	return r.client.SMembers(ctx, key).Result()
}

// DeleteUserPermissions 删除用户权限缓存
func (r *RedisClient) DeleteUserPermissions(ctx context.Context, userID string) error {
	key := PrefixUserPerm + userID
	return r.Del(ctx, key)
}

// ===== 证书缓存 =====

// CacheCertificate 缓存证书信息
func (r *RedisClient) CacheCertificate(ctx context.Context, certID string, data interface{}, expiration time.Duration) error {
	key := PrefixCert + certID
	return r.SetJSON(ctx, key, data, expiration)
}

// GetCertificate 获取证书缓存
func (r *RedisClient) GetCertificate(ctx context.Context, certID string, dest interface{}) error {
	key := PrefixCert + certID
	return r.GetJSON(ctx, key, dest)
}

// DeleteCertificate 删除证书缓存
func (r *RedisClient) DeleteCertificate(ctx context.Context, certID string) error {
	key := PrefixCert + certID
	return r.Del(ctx, key)
}

// CacheCertificateMetadata 缓存证书元数据
func (r *RedisClient) CacheCertificateMetadata(ctx context.Context, certID string, metadata map[string]interface{}, expiration time.Duration) error {
	key := PrefixCertMeta + certID
	for k, v := range metadata {
		if err := r.client.HSet(ctx, key, k, v).Err(); err != nil {
			return fmt.Errorf("hset field %s failed: %w", k, err)
		}
	}
	return r.client.Expire(ctx, key, expiration).Err()
}

// GetCertificateMetadata 获取证书元数据
func (r *RedisClient) GetCertificateMetadata(ctx context.Context, certID string) (map[string]string, error) {
	key := PrefixCertMeta + certID
	return r.HGetAll(ctx, key)
}

// ===== 验证缓存 =====

// CacheVerificationResult 缓存验证结果
func (r *RedisClient) CacheVerificationResult(ctx context.Context, verifyKey string, result interface{}, expiration time.Duration) error {
	key := PrefixVerify + verifyKey
	return r.SetJSON(ctx, key, result, expiration)
}

// GetVerificationResult 获取验证结果
func (r *RedisClient) GetVerificationResult(ctx context.Context, verifyKey string, dest interface{}) error {
	key := PrefixVerify + verifyKey
	return r.GetJSON(ctx, key, dest)
}

// ===== 系统配置缓存 =====

// CacheSysConfig 缓存系统配置
func (r *RedisClient) CacheSysConfig(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	cacheKey := PrefixSysConfig + key
	return r.Set(ctx, cacheKey, value, expiration)
}

// GetSysConfig 获取系统配置缓存
func (r *RedisClient) GetSysConfig(ctx context.Context, key string) (string, error) {
	cacheKey := PrefixSysConfig + key
	return r.Get(ctx, cacheKey)
}

// DeleteSysConfig 删除系统配置缓存
func (r *RedisClient) DeleteSysConfig(ctx context.Context, key string) error {
	cacheKey := PrefixSysConfig + key
	return r.Del(ctx, cacheKey)
}

// ===== 区块链交易缓存 =====

// CacheTransaction 缓存交易信息
func (r *RedisClient) CacheTransaction(ctx context.Context, txHash string, data interface{}, expiration time.Duration) error {
	key := PrefixTx + txHash
	return r.SetJSON(ctx, key, data, expiration)
}

// GetTransaction 获取交易缓存
func (r *RedisClient) GetTransaction(ctx context.Context, txHash string, dest interface{}) error {
	key := PrefixTx + txHash
	return r.GetJSON(ctx, key, dest)
}

// ===== 用户缓存 =====

// CacheUser 缓存用户信息
func (r *RedisClient) CacheUser(ctx context.Context, userID string, data interface{}, expiration time.Duration) error {
	key := PrefixUser + userID
	return r.SetJSON(ctx, key, data, expiration)
}

// GetUser 获取用户缓存
func (r *RedisClient) GetUser(ctx context.Context, userID string, dest interface{}) error {
	key := PrefixUser + userID
	return r.GetJSON(ctx, key, dest)
}

// DeleteUser 删除用户缓存
func (r *RedisClient) DeleteUser(ctx context.Context, userID string) error {
	key := PrefixUser + userID
	return r.Del(ctx, key)
}

// ===== 组织缓存 =====

// CacheOrganization 缓存组织信息
func (r *RedisClient) CacheOrganization(ctx context.Context, orgID string, data interface{}, expiration time.Duration) error {
	key := PrefixOrg + orgID
	return r.SetJSON(ctx, key, data, expiration)
}

// GetOrganization 获取组织缓存
func (r *RedisClient) GetOrganization(ctx context.Context, orgID string, dest interface{}) error {
	key := PrefixOrg + orgID
	return r.GetJSON(ctx, key, dest)
}

// DeleteOrganization 删除组织缓存
func (r *RedisClient) DeleteOrganization(ctx context.Context, orgID string) error {
	key := PrefixOrg + orgID
	return r.Del(ctx, key)
}
