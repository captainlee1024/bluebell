// Package redis provides ...
package redis

import (
	"fmt"

	"github.com/captainlee1024/bluebell/settings"
	"github.com/go-redis/redis"
)

var rdb *redis.Client

// Init 初始化 redis 连接
func Init(cfg *settings.RedisConfig) (err error) {
	rdb = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d",
			//viper.GetString("redis.host"),
			//viper.GetInt("redis.port"),
			cfg.Host,
			cfg.Port,
		),
		//Password: viper.GetString("redis.password"),
		//DB:       viper.GetInt("redis.db"),
		//PoolSize: viper.GetInt("redis.pool_size"),
		Password: cfg.Password,
		DB:       cfg.DB,
		PoolSize: cfg.PoolSize,
	})

	_, err = rdb.Ping().Result()
	return
}

// Close 关闭 redis 连接
func Close() {
	_ = rdb.Close()
}
