package global

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/go-redis/redis/v8"
	"project/order_web/config"
	"project/order_web/proto"
)

var (
	Trans        ut.Translator
	ServerConfig *config.ServerConfig = &config.ServerConfig{}
	NacosConfig  *config.NacosConfig  = &config.NacosConfig{}
	Rdb          *redis.Client

	GoodsSrvClient proto.GoodsClient
	OrderSrvClient proto.OrderClient
	InvSrvClient   proto.InventoryClient
)
