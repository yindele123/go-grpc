package global

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/go-redis/redis/v8"
	"google.golang.org/grpc"
	"project/user_web/config"
	"project/user_web/proto"
)

var (
	Trans        ut.Translator
	ServerConfig *config.ServerConfig = &config.ServerConfig{}
	Rdb          *redis.Client

	UserSrvClient proto.UserClient

	UserGrpcClient *grpc.ClientConn
)
