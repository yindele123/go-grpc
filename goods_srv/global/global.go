package global

import (
	"gorm.io/gorm"
	"project/goods_srv/config"
)

var (
	MysqlDb      *gorm.DB
	ServerConfig *config.ServerConfig = &config.ServerConfig{}

	NacosConfig *config.NacosConfig = &config.NacosConfig{}
)
