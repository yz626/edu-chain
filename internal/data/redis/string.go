package redis

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

// ===== 基础缓存操作 =====

// Get 获取缓存
func (r *RedisClient) Get(ctx context.Context, key string) (string, error) {
	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", nil // key不存在
		}
		return "", err
	}
	return val, nil
}

// Set 设置缓存
func (r *RedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	var val string
	switch v := value.(type) {
	case string:
		val = v
	case []byte:
		val = string(v)
	default:
		b, err := json.Marshal(v)
		if err != nil {
			return err
		}
		val = string(b)
	}
	return r.client.Set(ctx, key, val, expiration).Err()
}

// SetNX 设置缓存（仅当key不存在时）
func (r *RedisClient) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error) {
	var val string
	switch v := value.(type) {
	case string:
		val = v
	case []byte:
		val = string(v)
	default:
		b, err := json.Marshal(v)
		if err != nil {
			return false, err
		}
		val = string(b)
	}
	return r.client.SetNX(ctx, key, val, expiration).Result()
}

// Del 删除缓存
func (r *RedisClient) Del(ctx context.Context, keys ...string) error {
	if len(keys) == 0 {
		return nil
	}
	return r.client.Del(ctx, keys...).Err()
}

// Exists 检查key是否存在
func (r *RedisClient) Exists(ctx context.Context, keys ...string) (int64, error) {
	return r.client.Exists(ctx, keys...).Result()
}

// Expire 设置key的过期时间
func (r *RedisClient) Expire(ctx context.Context, key string, expiration time.Duration) (bool, error) {
	return r.client.Expire(ctx, key, expiration).Result()
}

// TTL 获取key的剩余过期时间
func (r *RedisClient) TTL(ctx context.Context, key string) (time.Duration, error) {
	return r.client.TTL(ctx, key).Result()
}

// ===== 复杂数据类型操作 =====

// GetJSON 获取JSON缓存
func (r *RedisClient) GetJSON(ctx context.Context, key string, dest interface{}) error {
	val, err := r.Get(ctx, key)
	if err != nil {
		return err
	}
	if val == "" {
		return nil
	}
	return json.Unmarshal([]byte(val), dest)
}

// SetJSON 设置JSON缓存
func (r *RedisClient) SetJSON(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return r.Set(ctx, key, value, expiration)
}

// Incr 增加数值
func (r *RedisClient) Incr(ctx context.Context, key string) (int64, error) {
	return r.client.Incr(ctx, key).Result()
}

// IncrBy 增加指定数值
func (r *RedisClient) IncrBy(ctx context.Context, key string, increment int64) (int64, error) {
	return r.client.IncrBy(ctx, key, increment).Result()
}

// Decr 减少数值
func (r *RedisClient) Decr(ctx context.Context, key string) (int64, error) {
	return r.client.Decr(ctx, key).Result()
}

// DecrBy 减少指定数值
func (r *RedisClient) DecrBy(ctx context.Context, key string, decrement int64) (int64, error) {
	return r.client.DecrBy(ctx, key, decrement).Result()
}
