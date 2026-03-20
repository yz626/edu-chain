package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/yz626/edu-chain/config"
	"github.com/yz626/edu-chain/pkg/logger"
)

var (
	// DefaultExpiration 默认过期时间
	DefaultExpiration = time.Hour * 24 // 24小时

	DefaultLoggerPrefix = "redis"
)

// RedisClient Redis客户端
type RedisClient struct {
	client *redis.Client
	log    logger.Logger
}

// NewRedisClient 创建Redis客户端
func NewRedisClient(cfg *config.RedisConfig, hepler logger.Logger) (*RedisClient, error) {
	log := hepler.Named(DefaultLoggerPrefix)

	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
		PoolSize: cfg.PoolSize,
	})

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		// log.Error("Redis连接失败", logger.AnyField("error", err))
		return nil, fmt.Errorf("redis connection failed: %w", err)
	}

	log.Info("Redis连接成功", logger.String("addr", fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)))

	return &RedisClient{
		client: client,
		log:    *log,
	}, nil
}

// NewRedisClientFromViper 从Viper创建Redis客户端
func NewRedisClientFromViper(v *config.Config, hepler logger.Logger) (*RedisClient, error) {
	return NewRedisClient(&v.Redis, hepler)
}

// Close 关闭Redis连接
func (r *RedisClient) Close() error {
	r.log.Info("关闭Redis连接")
	return r.client.Close()
}

// Ping 检查Redis连接
func (r *RedisClient) Ping(ctx context.Context) error {
	return r.client.Ping(ctx).Err()
}

// HealthCheck 健康检查
func (r *RedisClient) HealthCheck(ctx context.Context) (bool, string) {
	err := r.Ping(ctx)
	if err != nil {
		return false, err.Error()
	}
	return true, "OK"
}

// GetStats 获取连接池统计信息
func (r *RedisClient) GetStats() map[string]interface{} {
	stats := r.client.PoolStats()
	return map[string]interface{}{
		"hits":        stats.Hits,
		"misses":      stats.Misses,
		"idle_conns":  stats.IdleConns,
		"total_conns": stats.TotalConns,
		"stale_conns": stats.StaleConns,
	}
}

// GetClient 获取原生Redis客户端（用于高级操作）
func (r *RedisClient) GetClient() *redis.Client {
	return r.client
}

// Pipeline 创建流水线
func (r *RedisClient) Pipeline() redis.Pipeliner {
	return r.client.Pipeline()
}

// MGet 批量获取
func (r *RedisClient) MGet(ctx context.Context, keys ...string) ([]string, error) {
	results, err := r.client.MGet(ctx, keys...).Result()
	if err != nil {
		return nil, err
	}
	strResults := make([]string, len(results))
	for i, v := range results {
		if v != nil {
			switch val := v.(type) {
			case string:
				strResults[i] = val
			default:
				return nil, fmt.Errorf("unexpected type %T for key %s", v, keys[i])
			}
		}
	}
	return strResults, nil
}

// MSet 批量设置
func (r *RedisClient) MSet(ctx context.Context, values map[string]interface{}) error {
	ss := make([]interface{}, 0, len(values)*2)
	for k, v := range values {
		ss = append(ss, k, v)
	}
	return r.client.MSet(ctx, ss...).Err()
}
