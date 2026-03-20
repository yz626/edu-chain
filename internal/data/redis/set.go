package cache

import (
	"context"
)

// ===== Set操作 =====

// SAdd 添加集合成员
func (r *RedisClient) SAdd(ctx context.Context, key string, members ...interface{}) error {
	return r.client.SAdd(ctx, key, members...).Err()
}

// SRem 移除集合成员
func (r *RedisClient) SRem(ctx context.Context, key string, members ...interface{}) error {
	return r.client.SRem(ctx, key, members...).Err()
}

// SMembers 获取集合所有成员
func (r *RedisClient) SMembers(ctx context.Context, key string) ([]string, error) {
	return r.client.SMembers(ctx, key).Result()
}

// SIsMember 检查是否为集合成员
func (r *RedisClient) SIsMember(ctx context.Context, key string, member interface{}) (bool, error) {
	return r.client.SIsMember(ctx, key, member).Result()
}

// SCard 获取集合成员数量
func (r *RedisClient) SCard(ctx context.Context, key string) (int64, error) {
	return r.client.SCard(ctx, key).Result()
}

// SDiff 获取差集
func (r *RedisClient) SDiff(ctx context.Context, keys ...string) ([]string, error) {
	return r.client.SDiff(ctx, keys...).Result()
}

// SUnion 获取并集
func (r *RedisClient) SUnion(ctx context.Context, keys ...string) ([]string, error) {
	return r.client.SUnion(ctx, keys...).Result()
}

// SInter 获取交集
func (r *RedisClient) SInter(ctx context.Context, keys ...string) ([]string, error) {
	return r.client.SInter(ctx, keys...).Result()
}

// SDiffStore 将差集结果存储到新集合
func (r *RedisClient) SDiffStore(ctx context.Context, dest string, keys ...string) (int64, error) {
	return r.client.SDiffStore(ctx, dest, keys...).Result()
}

// SUnionStore 将并集结果存储到新集合
func (r *RedisClient) SUnionStore(ctx context.Context, dest string, keys ...string) (int64, error) {
	return r.client.SUnionStore(ctx, dest, keys...).Result()
}

// SInterStore 将交集结果存储到新集合
func (r *RedisClient) SInterStore(ctx context.Context, dest string, keys ...string) (int64, error) {
	return r.client.SInterStore(ctx, dest, keys...).Result()
}

// SMove 移动集合成员
func (r *RedisClient) SMove(ctx context.Context, src, dest string, member interface{}) (bool, error) {
	return r.client.SMove(ctx, src, dest, member).Result()
}

// SPop 随机弹出成员
func (r *RedisClient) SPop(ctx context.Context, key string) (string, error) {
	return r.client.SPop(ctx, key).Result()
}

// SPopN 随机弹出多个成员
func (r *RedisClient) SPopN(ctx context.Context, key string, count int64) ([]string, error) {
	return r.client.SPopN(ctx, key, count).Result()
}

// SRandMember 随机获取成员
func (r *RedisClient) SRandMember(ctx context.Context, key string) (string, error) {
	return r.client.SRandMember(ctx, key).Result()
}

// SRandMemberN 随机获取多个成员
func (r *RedisClient) SRandMemberN(ctx context.Context, key string, count int64) ([]string, error) {
	return r.client.SRandMemberN(ctx, key, count).Result()
}

// SScan 迭代集合成员
func (r *RedisClient) SScan(ctx context.Context, key string, cursor uint64, match string, count int64) ([]string, uint64, error) {
	return r.client.SScan(ctx, key, cursor, match, count).Result()
}
