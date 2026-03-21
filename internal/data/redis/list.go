package redis

import (
	"context"
	"time"
)

// ===== List操作 =====

// LPush 将元素推入列表左侧
func (r *RedisClient) LPush(ctx context.Context, key string, values ...interface{}) error {
	return r.client.LPush(ctx, key, values...).Err()
}

// RPush 将元素推入列表右侧
func (r *RedisClient) RPush(ctx context.Context, key string, values ...interface{}) error {
	return r.client.RPush(ctx, key, values...).Err()
}

// LPop 从列表左侧弹出元素
func (r *RedisClient) LPop(ctx context.Context, key string) (string, error) {
	return r.client.LPop(ctx, key).Result()
}

// RPop 从列表右侧弹出元素
func (r *RedisClient) RPop(ctx context.Context, key string) (string, error) {
	return r.client.RPop(ctx, key).Result()
}

// LRange 获取列表范围元素
func (r *RedisClient) LRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	return r.client.LRange(ctx, key, start, stop).Result()
}

// LLen 获取列表长度
func (r *RedisClient) LLen(ctx context.Context, key string) (int64, error) {
	return r.client.LLen(ctx, key).Result()
}

// LIndex 获取列表指定位置的元素
func (r *RedisClient) LIndex(ctx context.Context, key string, index int64) (string, error) {
	return r.client.LIndex(ctx, key, index).Result()
}

// LInsert 在列表指定位置插入元素
func (r *RedisClient) LInsert(ctx context.Context, key, op string, pivot, value interface{}) (int64, error) {
	return r.client.LInsert(ctx, key, op, pivot, value).Result()
}

// LRem 从列表中移除元素
func (r *RedisClient) LRem(ctx context.Context, key string, count int64, value interface{}) (int64, error) {
	return r.client.LRem(ctx, key, count, value).Result()
}

// LSet 设置列表指定位置的元素
func (r *RedisClient) LSet(ctx context.Context, key string, index int64, value interface{}) error {
	return r.client.LSet(ctx, key, index, value).Err()
}

// LTrim 修剪列表
func (r *RedisClient) LTrim(ctx context.Context, key string, start, stop int64) error {
	return r.client.LTrim(ctx, key, start, stop).Err()
}

// LPopCount 弹出多个元素
func (r *RedisClient) LPopCount(ctx context.Context, key string, count int64) ([]string, error) {
	var results []string
	for i := int64(0); i < count; i++ {
		val, err := r.client.LPop(ctx, key).Result()
		if err != nil {
			if i > 0 {
				// 如果已经弹出了一些元素，则继续返回
				return results, nil
			}
			return nil, err
		}
		results = append(results, val)
	}
	return results, nil
}

// RPopCount 弹出多个元素（从右侧）
func (r *RedisClient) RPopCount(ctx context.Context, key string, count int64) ([]string, error) {
	var results []string
	for i := int64(0); i < count; i++ {
		val, err := r.client.RPop(ctx, key).Result()
		if err != nil {
			if i > 0 {
				return results, nil
			}
			return nil, err
		}
		results = append(results, val)
	}
	return results, nil
}

// BRPop 阻塞式从列表右侧弹出元素
func (r *RedisClient) BRPop(ctx context.Context, timeout time.Duration, keys ...string) ([]string, error) {
	return r.client.BRPop(ctx, timeout, keys...).Result()
}

// BLPop 阻塞式从列表左侧弹出元素
func (r *RedisClient) BLPop(ctx context.Context, timeout time.Duration, keys ...string) ([]string, error) {
	return r.client.BLPop(ctx, timeout, keys...).Result()
}
