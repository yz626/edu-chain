package redis

import (
	"context"
	"encoding/json"
)

// ===== Hash操作 =====

// HGet 获取hash字段值
func (r *RedisClient) HGet(ctx context.Context, key, field string) (string, error) {
	return r.client.HGet(ctx, key, field).Result()
}

// HSet 设置hash字段值
func (r *RedisClient) HSet(ctx context.Context, key, field string, value interface{}) error {
	var val interface{} = value
	switch v := value.(type) {
	case string, int, int64, float64, bool:
		val = v
	default:
		b, err := json.Marshal(v)
		if err != nil {
			return err
		}
		val = string(b)
	}
	return r.client.HSet(ctx, key, field, val).Err()
}

// HDel 删除hash字段
func (r *RedisClient) HDel(ctx context.Context, key string, fields ...string) error {
	return r.client.HDel(ctx, key, fields...).Err()
}

// HGetAll 获取所有hash字段
func (r *RedisClient) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	return r.client.HGetAll(ctx, key).Result()
}

// HExists 检查hash字段是否存在
func (r *RedisClient) HExists(ctx context.Context, key, field string) (bool, error) {
	return r.client.HExists(ctx, key, field).Result()
}

// HIncrBy 增加hash字段值
func (r *RedisClient) HIncrBy(ctx context.Context, key, field string, incr int64) (int64, error) {
	return r.client.HIncrBy(ctx, key, field, incr).Result()
}

// HIncrByFloat 增加hash字段浮点值
func (r *RedisClient) HIncrByFloat(ctx context.Context, key, field string, incr float64) (float64, error) {
	return r.client.HIncrByFloat(ctx, key, field, incr).Result()
}

// HLen 获取hash字段数量
func (r *RedisClient) HLen(ctx context.Context, key string) (int64, error) {
	return r.client.HLen(ctx, key).Result()
}

// HMGet 批量获取hash字段值
func (r *RedisClient) HMGet(ctx context.Context, key string, fields ...string) ([]interface{}, error) {
	return r.client.HMGet(ctx, key, fields...).Result()
}

// HMSet 批量设置hash字段值
func (r *RedisClient) HMSet(ctx context.Context, key string, values map[string]interface{}) error {
	return r.client.HSet(ctx, key, values).Err()
}

// HKeys 获取hash所有字段
func (r *RedisClient) HKeys(ctx context.Context, key string) ([]string, error) {
	return r.client.HKeys(ctx, key).Result()
}

// HVals 获取hash所有值
func (r *RedisClient) HVals(ctx context.Context, key string) ([]string, error) {
	return r.client.HVals(ctx, key).Result()
}

// HScan 迭代hash字段
func (r *RedisClient) HScan(ctx context.Context, key string, cursor uint64, match string, count int64) ([]string, uint64, error) {
	return r.client.HScan(ctx, key, cursor, match, count).Result()
}
