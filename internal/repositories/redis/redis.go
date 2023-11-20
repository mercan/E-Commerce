package redis

import (
	"github.com/mercan/ecommerce/internal/config"
	"github.com/mercan/ecommerce/internal/utils"
	"github.com/redis/go-redis/v9"
	"log"
)

var client = Connect()

func Connect() *redis.Client {
	ctx, cancel := utils.ContextWithTimeout(10)
	defer cancel()

	opt, err := redis.ParseURL(config.GetRedisConfig().URI)
	if err != nil {
		panic(err)
	}

	rdb := redis.NewClient(opt)

	if _, err := rdb.Ping(ctx).Result(); err != nil {

		panic(err)
	}

	log.Println("Connected to Redis")
	return rdb
}
