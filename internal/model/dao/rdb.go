package dao

import (
	"DigitalCurrency/internal/config"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
)

var Rdb *redis.Client

func InitRedis(conf *config.RedisConf) {
	options := redis.Options{
		Addr: fmt.Sprintf(
			"%s:%d",
			conf.Host,
			conf.Port), // Redis地址
		DB:          conf.Database,                   // Redis库
		PoolSize:    5,                               // Redis连接池大小
		MaxRetries:  3,                               // 最大重试次数
		IdleTimeout: time.Second * time.Duration(10), // 空闲链接超时时间
		Password:    conf.Password,
	}
	if viper.GetString("redis_passwd") != "" {
		options.Password = viper.GetString("redis_passwd")
	}
	Rdb = redis.NewClient(&options)
}
