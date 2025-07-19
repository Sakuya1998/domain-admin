package cache

import (
	"context"
	"domain-admin/pkg/config"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client

// 缓存键前缀
const (
	UserCachePrefix     = "user:"
	UserListCachePrefix = "user_list:"
	SessionCachePrefix  = "session:"
)

// 缓存过期时间
const (
	UserCacheExpire     = 30 * time.Minute  // 用户信息缓存30分钟
	UserListCacheExpire = 5 * time.Minute   // 用户列表缓存5分钟
	SessionCacheExpire  = 24 * time.Hour    // 会话缓存24小时
)

func InitCache(cfg config.RedisConfig) {
	redisClient = redis.NewClient(&redis.Options{
		Addr:         cfg.Addr,
		Password:     cfg.Password,
		DB:           cfg.DB,
		ReadTimeout:  parseDuration(cfg.ReadTimeout),
		WriteTimeout: parseDuration(cfg.WriteTimeout),
	})

	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		// Redis连接失败时不panic，而是记录错误并将redisClient设为nil
		fmt.Printf("Warning: Failed to connect to Redis: %v\n", err)
		fmt.Println("Application will continue without Redis cache")
		redisClient = nil
	}
}

func GetRedisClient() *redis.Client {
	return redisClient
}

// Set 设置缓存
func Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return redisClient.Set(ctx, key, data, expiration).Err()
}

// Get 获取缓存
func Get(ctx context.Context, key string, dest interface{}) error {
	data, err := redisClient.Get(ctx, key).Result()
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(data), dest)
}

// Del 删除缓存
func Del(ctx context.Context, keys ...string) error {
	return redisClient.Del(ctx, keys...).Err()
}

// Exists 检查缓存是否存在
func Exists(ctx context.Context, key string) (bool, error) {
	result, err := redisClient.Exists(ctx, key).Result()
	return result > 0, err
}

// SetUserCache 设置用户缓存
func SetUserCache(ctx context.Context, userID uint, user interface{}) error {
	if redisClient == nil {
		return nil // Redis未初始化时静默返回
	}
	key := fmt.Sprintf("%s%d", UserCachePrefix, userID)
	return Set(ctx, key, user, UserCacheExpire)
}

// GetUserCache 获取用户缓存
func GetUserCache(ctx context.Context, userID uint, dest interface{}) error {
	if redisClient == nil {
		return redis.Nil // Redis未初始化时返回未找到
	}
	key := fmt.Sprintf("%s%d", UserCachePrefix, userID)
	return Get(ctx, key, dest)
}

// DelUserCache 删除用户缓存
func DelUserCache(ctx context.Context, userID uint) error {
	if redisClient == nil {
		return nil // Redis未初始化时静默返回
	}
	key := fmt.Sprintf("%s%d", UserCachePrefix, userID)
	return Del(ctx, key)
}

// SetUserListCache 设置用户列表缓存
func SetUserListCache(ctx context.Context, cacheKey string, data interface{}) error {
	if redisClient == nil {
		return nil // Redis未初始化时静默返回
	}
	key := fmt.Sprintf("%s%s", UserListCachePrefix, cacheKey)
	return Set(ctx, key, data, UserListCacheExpire)
}

// GetUserListCache 获取用户列表缓存
func GetUserListCache(ctx context.Context, cacheKey string, dest interface{}) error {
	if redisClient == nil {
		return redis.Nil // Redis未初始化时返回未找到
	}
	key := fmt.Sprintf("%s%s", UserListCachePrefix, cacheKey)
	return Get(ctx, key, dest)
}

// DelUserListCache 删除用户列表缓存
func DelUserListCache(ctx context.Context, pattern string) error {
	if redisClient == nil {
		return nil // Redis未初始化时静默返回
	}
	keys, err := redisClient.Keys(ctx, UserListCachePrefix+pattern).Result()
	if err != nil {
		return err
	}
	if len(keys) > 0 {
		return Del(ctx, keys...)
	}
	return nil
}

// SetSessionCache 设置会话缓存
func SetSessionCache(ctx context.Context, token string, userInfo interface{}) error {
	if redisClient == nil {
		return nil // Redis未初始化时静默返回
	}
	key := fmt.Sprintf("%s%s", SessionCachePrefix, token)
	return Set(ctx, key, userInfo, SessionCacheExpire)
}

// GetSessionCache 获取会话缓存
func GetSessionCache(ctx context.Context, token string, dest interface{}) error {
	if redisClient == nil {
		return redis.Nil // Redis未初始化时返回未找到
	}
	key := fmt.Sprintf("%s%s", SessionCachePrefix, token)
	return Get(ctx, key, dest)
}

// DelSessionCache 删除会话缓存
func DelSessionCache(ctx context.Context, token string) error {
	if redisClient == nil {
		return nil // Redis未初始化时静默返回
	}
	key := fmt.Sprintf("%s%s", SessionCachePrefix, token)
	return Del(ctx, key)
}

func parseDuration(s string) time.Duration {
	d, err := time.ParseDuration(s)
	if err != nil {
		return 0
	}
	return d
}
