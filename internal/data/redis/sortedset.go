package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// ===== Sorted Set操作 =====

// ZAdd 添加有序集合成员
func (r *RedisClient) ZAdd(ctx context.Context, key string, members ...redis.Z) error {
	return r.client.ZAdd(ctx, key, members...).Err()
}

// ZRem 移除有序集合成员
func (r *RedisClient) ZRem(ctx context.Context, key string, members ...interface{}) error {
	return r.client.ZRem(ctx, key, members...).Err()
}

// ZCard 获取有序集合成员数量
func (r *RedisClient) ZCard(ctx context.Context, key string) (int64, error) {
	return r.client.ZCard(ctx, key).Result()
}

// ZCount 获取指定分数范围的成员数量
func (r *RedisClient) ZCount(ctx context.Context, key string, min, max string) (int64, error) {
	return r.client.ZCount(ctx, key, min, max).Result()
}

// ZRange 获取指定索引范围的成员
func (r *RedisClient) ZRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	return r.client.ZRange(ctx, key, start, stop).Result()
}

// ZRangeWithScores 获取指定索引范围的成员及分数
func (r *RedisClient) ZRangeWithScores(ctx context.Context, key string, start, stop int64) ([]redis.Z, error) {
	return r.client.ZRangeWithScores(ctx, key, start, stop).Result()
}

// ZRangeByScore 按分数范围获取成员
func (r *RedisClient) ZRangeByScore(ctx context.Context, key string, opt *redis.ZRangeBy) ([]string, error) {
	return r.client.ZRangeByScore(ctx, key, opt).Result()
}

// ZRangeByScoreWithScores 按分数范围获取成员及分数
func (r *RedisClient) ZRangeByScoreWithScores(ctx context.Context, key string, opt *redis.ZRangeBy) ([]redis.Z, error) {
	return r.client.ZRangeByScoreWithScores(ctx, key, opt).Result()
}

// ZRevRange 获取指定索引范围的成员（从大到小）
func (r *RedisClient) ZRevRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	return r.client.ZRevRange(ctx, key, start, stop).Result()
}

// ZRevRangeWithScores 获取指定索引范围的成员及分数（从大到小）
func (r *RedisClient) ZRevRangeWithScores(ctx context.Context, key string, start, stop int64) ([]redis.Z, error) {
	return r.client.ZRevRangeWithScores(ctx, key, start, stop).Result()
}

// ZRevRangeByScore 按分数范围获取成员（从大到小）
func (r *RedisClient) ZRevRangeByScore(ctx context.Context, key string, opt *redis.ZRangeBy) ([]string, error) {
	return r.client.ZRevRangeByScore(ctx, key, opt).Result()
}

// ZRank 获取成员排名（从小到大）
func (r *RedisClient) ZRank(ctx context.Context, key, member string) (int64, error) {
	return r.client.ZRank(ctx, key, member).Result()
}

// ZRevRank 获取成员排名（从大到小）
func (r *RedisClient) ZRevRank(ctx context.Context, key, member string) (int64, error) {
	return r.client.ZRevRank(ctx, key, member).Result()
}

// ZScore 获取成员分数
func (r *RedisClient) ZScore(ctx context.Context, key, member string) (float64, error) {
	return r.client.ZScore(ctx, key, member).Result()
}

// ZIncrBy 增加成员分数
func (r *RedisClient) ZIncrBy(ctx context.Context, key string, increment float64, member string) (float64, error) {
	return r.client.ZIncrBy(ctx, key, increment, member).Result()
}

// ZRemRangeByRank 移除指定排名范围的成员
func (r *RedisClient) ZRemRangeByRank(ctx context.Context, key string, start, stop int64) (int64, error) {
	return r.client.ZRemRangeByRank(ctx, key, start, stop).Result()
}

// ZRemRangeByScore 移除指定分数范围的成员
func (r *RedisClient) ZRemRangeByScore(ctx context.Context, key string, min, max string) (int64, error) {
	return r.client.ZRemRangeByScore(ctx, key, min, max).Result()
}

// ZUnionStore 并集存储
func (r *RedisClient) ZUnionStore(ctx context.Context, dest string, store *redis.ZStore) (int64, error) {
	return r.client.ZUnionStore(ctx, dest, store).Result()
}

// ZInterStore 交集存储
func (r *RedisClient) ZInterStore(ctx context.Context, dest string, store *redis.ZStore) (int64, error) {
	return r.client.ZInterStore(ctx, dest, store).Result()
}

// ZScan 迭代有序集合
func (r *RedisClient) ZScan(ctx context.Context, key string, cursor uint64, match string, count int64) ([]string, uint64, error) {
	return r.client.ZScan(ctx, key, cursor, match, count).Result()
}

// ZPopMax 弹出分数最大的成员
func (r *RedisClient) ZPopMax(ctx context.Context, key string, count ...int64) ([]redis.Z, error) {
	return r.client.ZPopMax(ctx, key, count...).Result()
}

// ZPopMin 弹出分数最小的成员
func (r *RedisClient) ZPopMin(ctx context.Context, key string, count ...int64) ([]redis.Z, error) {
	return r.client.ZPopMin(ctx, key, count...).Result()
}
