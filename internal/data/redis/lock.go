package cache

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// Lock 分布式锁
type Lock struct {
	client  *redis.Client
	key     string
	lockKey string
	ttl     time.Duration
}

// AcquireLock 获取分布式锁
// 参数:
//   - ctx: 上下文
//   - key: 锁的资源标识
//   - ttl: 锁的过期时间
//
// 返回:
//   - *Lock: 锁对象，使用完后需调用 ReleaseLock 释放
//   - error: 获取锁失败时返回错误
func (r *RedisClient) AcquireLock(ctx context.Context, key string, ttl time.Duration) (*Lock, error) {
	lockKey := PrefixLock + key
	ok, err := r.client.SetNX(ctx, lockKey, "1", ttl).Result()
	if err != nil {
		return nil, fmt.Errorf("acquire lock failed: %w", err)
	}
	if !ok {
		return nil, errors.New("lock already held")
	}
	return &Lock{
		client:  r.client,
		key:     key,
		lockKey: lockKey,
		ttl:     ttl,
	}, nil
}

// AcquireLockWithRetry 获取分布式锁（带重试）
// 参数:
//   - ctx: 上下文
//   - key: 锁的资源标识
//   - ttl: 锁的过期时间
//   - retryTimes: 重试次数
//   - retryDelay: 重试间隔
//
// 返回:
//   - *Lock: 锁对象
//   - error: 获取锁失败时返回错误
func (r *RedisClient) AcquireLockWithRetry(ctx context.Context, key string, ttl time.Duration, retryTimes int, retryDelay time.Duration) (*Lock, error) {
	for i := 0; i < retryTimes; i++ {
		lock, err := r.AcquireLock(ctx, key, ttl)
		if err == nil {
			return lock, nil
		}
		if i < retryTimes-1 {
			time.Sleep(retryDelay)
		}
	}
	return nil, fmt.Errorf("failed to acquire lock after %d retries", retryTimes)
}

// ReleaseLock 释放分布式锁
func (l *Lock) ReleaseLock(ctx context.Context) error {
	return l.client.Del(ctx, l.lockKey).Err()
}

// ExtendLock 延长锁的过期时间
func (l *Lock) ExtendLock(ctx context.Context, ttl time.Duration) (bool, error) {
	// 只有锁还存在且值未改变时才延长
	script := redis.NewScript(`if redis.call("get", KEYS[1]) == ARGV[1] then return redis.call("pexpire", KEYS[1], ARGV[2]) else return 0 end`)
	result, err := script.Run(ctx, l.client, []string{l.lockKey}, "1", ttl.Milliseconds()).Int()
	if err != nil {
		return false, err
	}
	return result == 1, nil
}

// GetLockKey 获取锁的key
func (l *Lock) GetLockKey() string {
	return l.lockKey
}

// GetTTL 获取锁的剩余过期时间
func (l *Lock) GetTTL(ctx context.Context) (time.Duration, error) {
	return l.client.TTL(ctx, l.lockKey).Result()
}

// LockOption 锁选项
type LockOption func(*LockOptions)

// LockOptions 锁配置选项
type LockOptions struct {
	Expire     time.Duration // 锁过期时间
	RetryTimes int           // 重试次数
	RetryDelay time.Duration // 重试间隔
}

// WithLockExpire 设置锁过期时间
func WithLockExpire(expire time.Duration) LockOption {
	return func(o *LockOptions) {
		o.Expire = expire
	}
}

// WithLockRetry 设置重试次数和间隔
func WithLockRetry(retryTimes int, retryDelay time.Duration) LockOption {
	return func(o *LockOptions) {
		o.RetryTimes = retryTimes
		o.RetryDelay = retryDelay
	}
}

// AcquireLockWithOptions 使用选项获取分布式锁
func (r *RedisClient) AcquireLockWithOptions(ctx context.Context, key string, opts ...LockOption) (*Lock, error) {
	options := &LockOptions{
		Expire:     10 * time.Second,
		RetryTimes: 3,
		RetryDelay: 100 * time.Millisecond,
	}

	for _, opt := range opts {
		opt(options)
	}

	if options.RetryTimes > 0 {
		return r.AcquireLockWithRetry(ctx, key, options.Expire, options.RetryTimes, options.RetryDelay)
	}
	return r.AcquireLock(ctx, key, options.Expire)
}

// DoWithLock 使用锁执行操作
// 参数:
//   - ctx: 上下文
//   - key: 锁的资源标识
//   - fn: 需要在锁保护下执行的函数
//
// 返回:
//   - error: 函数执行失败时返回错误
func (r *RedisClient) DoWithLock(ctx context.Context, key string, fn func() error) error {
	lock, err := r.AcquireLock(ctx, key, 10*time.Second)
	if err != nil {
		return fmt.Errorf("acquire lock failed: %w", err)
	}
	defer lock.ReleaseLock(ctx)

	return fn()
}

// DoWithLockExpire 使用指定过期时间的锁执行操作
func (r *RedisClient) DoWithLockExpire(ctx context.Context, key string, expire time.Duration, fn func() error) error {
	lock, err := r.AcquireLock(ctx, key, expire)
	if err != nil {
		return fmt.Errorf("acquire lock failed: %w", err)
	}
	defer lock.ReleaseLock(ctx)

	return fn()
}
