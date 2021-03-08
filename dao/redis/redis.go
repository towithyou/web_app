package redis

import (
	"fmt"

	"github.com/towithyou/web_app/settings"

	"github.com/go-redis/redis"
)

var rdb *redis.Client

func Init(cfg *settings.RedisConfig) (err error) {
	rdb = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password, // no password set
		DB:       cfg.Db,       // use default DB
		PoolSize: cfg.PoolSize,
	})

	_, err = rdb.Ping().Result()
	return
}

func Close() {
	_ = rdb.Close()
}
