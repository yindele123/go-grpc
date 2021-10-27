package global

import (
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"project/inventory_srv/config"
)

var (
	MysqlDb      *gorm.DB
	ServerConfig *config.ServerConfig = &config.ServerConfig{}

	NacosConfig *config.NacosConfig = &config.NacosConfig{}

	Rdb          *redis.Client
)
