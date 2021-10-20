package global

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/go-redis/redis/v8"
	"project/goods_web/config"
	"project/goods_web/proto"
)

var (
	Trans        ut.Translator
	ServerConfig *config.ServerConfig = &config.ServerConfig{}
	NacosConfig  *config.NacosConfig  = &config.NacosConfig{}
	Rdb          *redis.Client

	GoodsSrvClient proto.GoodsClient
	BannerSrvClient proto.BannersClient
	CategorySrvClient proto.CategoryClient
	BrandSrvClient proto.BrandsClient

)
