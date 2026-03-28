package redis

import (
	"context"
	"fmt"
	"moxin-zhicheng/internal/config" // 确保你的 config 包能读取到 redis 配置项
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	RedisClient *redis.Client
	ctx         = context.Background()
)

// InitRedis 初始化 Redis 连接
func InitRedis() {
	conf := config.Conf.Redis // 假设你的 config 结构体里有 Redis 字段

	RedisClient = redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%s", conf.Host, conf.Port),
		Password:     conf.Password, // config_dev.yaml 中的密码
		DB:           conf.DB,       // 默认 0
		PoolSize:     20,            // 连接池大小
		MinIdleConns: 5,             // 最小空闲连接数
	})

	// 测试连接状况
	timeoutCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := RedisClient.Ping(timeoutCtx).Result()
	if err != nil {
		panic(fmt.Sprintf("❌ Redis 连接失败: %v", err))
	}

	fmt.Println("🌟 Redis 连接成功！已开启数据预热准备。")
}

func Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return RedisClient.Set(ctx, key, value, expiration).Err()
}

// 封装一个通用的 Get 方法
func Get(ctx context.Context, key string) (string, error) {
	return RedisClient.Get(ctx, key).Result()
}

// CloseRedis 关闭连接
func CloseRedis() {
	if RedisClient != nil {
		_ = RedisClient.Close()
	}
}
