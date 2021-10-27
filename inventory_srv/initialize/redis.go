package initialize

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"project/inventory_srv/global"
)

var ctx = context.Background()
func InitRedis() {
	global.Rdb = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", global.ServerConfig.RedisInfo.Host, global.ServerConfig.RedisInfo.Port),
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}
