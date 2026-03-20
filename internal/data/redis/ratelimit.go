package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// RateLimit 限流检查
// 使用滑动窗口算法实现限流
// 参数:
//   - ctx: 上下文
//   - key: 限流的资源标识
//   - limit: 时间窗口内允许的最大请求数
//   - window: 时间窗口大小
//
// 返回:
//   - bool: 是否允许通过
//   - error: 错误信息
func (r *RedisClient) RateLimit(ctx context.Context, key string, limit int, window time.Duration) (bool, error) {
	rateKey := PrefixRateLimit + key

	// 使用滑动窗口算法
	now := time.Now().Unix()
	windowStart := now - int64(window.Seconds())

	// 移除窗口外的记录
	r.client.ZRemRangeByScore(ctx, rateKey, "0", fmt.Sprintf("%d", windowStart))

	// 获取当前窗口内的请求数
	count, err := r.client.ZCard(ctx, rateKey).Result()
	if err != nil {
		return false, err
	}

	// 检查是否超过限制
	if count >= int64(limit) {
		return false, nil
	}

	// 添加当前请求
	r.client.ZAdd(ctx, rateKey, redis.Z{Score: float64(now), Member: fmt.Sprintf("%d", now)})
	r.client.Expire(ctx, rateKey, window)

	return true, nil
}

// RateLimitWithPrefix 使用指定前缀的限流检查
func (r *RedisClient) RateLimitWithPrefix(ctx context.Context, prefix, key string, limit int, window time.Duration) (bool, error) {
	rateKey := prefix + key
	return r.rateLimit(ctx, rateKey, limit, window)
}

// rateLimit 限流核心实现
func (r *RedisClient) rateLimit(ctx context.Context, rateKey string, limit int, window time.Duration) (bool, error) {
	now := time.Now().Unix()
	windowStart := now - int64(window.Seconds())

	r.client.ZRemRangeByScore(ctx, rateKey, "0", fmt.Sprintf("%d", windowStart))

	count, err := r.client.ZCard(ctx, rateKey).Result()
	if err != nil {
		return false, err
	}

	if count >= int64(limit) {
		return false, nil
	}

	r.client.ZAdd(ctx, rateKey, redis.Z{Score: float64(now), Member: fmt.Sprintf("%d", now)})
	r.client.Expire(ctx, rateKey, window)

	return true, nil
}

// RateLimitCounter 基于计数器的限流（固定窗口）
// 这种方式比滑动窗口性能更好，但可能产生边界突刺
func (r *RedisClient) RateLimitCounter(ctx context.Context, key string, limit int, window time.Duration) (bool, error) {
	rateKey := PrefixRateLimit + "counter:" + key

	// 使用INCR实现原子计数
	count, err := r.client.Incr(ctx, rateKey).Result()
	if err != nil {
		return false, err
	}

	// 第一次设置过期时间
	if count == 1 {
		r.client.Expire(ctx, rateKey, window)
	}

	return count <= int64(limit), nil
}

// RateLimitInfo 获取限流信息
type RateLimitInfo struct {
	Allowed    bool          `json:"allowed"`
	Limit      int           `json:"limit"`
	Remaining  int           `json:"remaining"`
	ResetTime  time.Time     `json:"reset_time"`
	RetryAfter time.Duration `json:"retry_after"`
}

// GetRateLimitInfo 获取限流详细信息
func (r *RedisClient) GetRateLimitInfo(ctx context.Context, key string, limit int, window time.Duration) (*RateLimitInfo, error) {
	rateKey := PrefixRateLimit + key
	now := time.Now().Unix()
	windowStart := now - int64(window.Seconds())

	r.client.ZRemRangeByScore(ctx, rateKey, "0", fmt.Sprintf("%d", windowStart))

	count, err := r.client.ZCard(ctx, rateKey).Result()
	if err != nil {
		return nil, err
	}

	remaining := limit - int(count)
	if remaining < 0 {
		remaining = 0
	}

	// 获取key的TTL作为重置时间
	ttl, _ := r.client.TTL(ctx, rateKey).Result()

	return &RateLimitInfo{
		Allowed:    count < int64(limit),
		Limit:      limit,
		Remaining:  remaining,
		ResetTime:  time.Now().Add(ttl),
		RetryAfter: ttl,
	}, nil
}

// ClearRateLimit 清除限流记录
func (r *RedisClient) ClearRateLimit(ctx context.Context, key string) error {
	rateKey := PrefixRateLimit + key
	return r.client.Del(ctx, rateKey).Err()
}

// RateLimiter 限流器封装
type RateLimiter struct {
	client *RedisClient
	key    string
	limit  int
	window time.Duration
}

// NewRateLimiter 创建限流器
func NewRateLimiter(client *RedisClient, key string, limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		client: client,
		key:    key,
		limit:  limit,
		window: window,
	}
}

// Allow 检查是否允许请求
func (rl *RateLimiter) Allow(ctx context.Context) (bool, error) {
	return rl.client.RateLimit(ctx, rl.key, rl.limit, rl.window)
}

// GetInfo 获取限流信息
func (rl *RateLimiter) GetInfo(ctx context.Context) (*RateLimitInfo, error) {
	return rl.client.GetRateLimitInfo(ctx, rl.key, rl.limit, rl.window)
}

// Reset 重置限流器
func (rl *RateLimiter) Reset(ctx context.Context) error {
	return rl.client.ClearRateLimit(ctx, rl.key)
}
