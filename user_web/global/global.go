package global

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/go-redis/redis/v8"
	"project/user_web/config"
)

var (
	Trans        ut.Translator
	ServerConfig *config.ServerConfig = &config.ServerConfig{}
	Rdb *redis.Client
)
