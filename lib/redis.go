package lib

import (
	"context"

	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
	"github.com/go-redis/redis/v8"
)

var (
	RedisClient *redis.Client
	ctx         = context.Background()
)

func InitRedis() {
	redisHost := web.AppConfig.DefaultString("redis_host", "localhost")
	redisPort := web.AppConfig.DefaultString("redis_port", "6379")
	redisPassword := web.AppConfig.DefaultString("redis_password", "")
	redisDb := web.AppConfig.DefaultInt("redis_db", 0)

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     redisHost + ":" + redisPort,
		Password: redisPassword,
		DB:       redisDb,
	})

	// 测试连接
	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		logs.Error("Redis连接失败: %v", err)
		panic(err)
	}
}
