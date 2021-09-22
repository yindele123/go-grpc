package main

import (
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
	"project/user_srv/global"
	"project/user_srv/handler"
	"project/user_srv/initialize"
	"project/user_srv/proto"
)

func main() {
	initialize.InitConfig()
	//初始化数据库
	initialize.InitMysql()
	initialize.InitLogger()

	g := grpc.NewServer()
	userServer := handler.UserServer{}
	proto.RegisterUserServer(g, &userServer)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d",global.ServerConfig.Port))
	if err != nil {
		zap.S().Panic("启动失败:", err.Error())
	}
	_ = g.Serve(lis)

}
