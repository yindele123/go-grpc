package initialize

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"project/user_web/global"
	"project/user_web/proto"
	_ "github.com/mbobakov/grpc-consul-resolver" // It's important
)
//负载
func InitSrvConn(){
	consulInfo := global.ServerConfig.ConsulInfo
	userConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", consulInfo.Host, consulInfo.Port, global.ServerConfig.UserSrvInfo.Name),
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		zap.S().Fatal("[InitSrvConn] 连接 【用户服务失败】")
	}

	userSrvClient := proto.NewUserClient(userConn)
	global.UserSrvClient = userSrvClient
}

//普通拉起服务
func InitSrvConn2()  {
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d",global.ServerConfig.ConsulInfo.Host,global.ServerConfig.ConsulInfo.Port)

	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}
	userSrvHost := ""
	userSrvPort := 0
	data, err := client.Agent().ServicesWithFilter(fmt.Sprintf("Service == \"%s\"", global.ServerConfig.UserSrvInfo.Name))
	fmt.Println(global.ServerConfig.UserSrvInfo.Name)
	fmt.Println(err)
	if err != nil {
		panic(err)
	}
	for _, value := range data{
		userSrvHost = value.Address
		userSrvPort = value.Port
		break
	}

	if userSrvHost == ""{
		zap.S().Fatal("[InitSrvConn] 连接 【用户服务失败】")
		return
	}
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", userSrvHost,userSrvPort), grpc.WithInsecure())
	if err != nil {
		zap.S().Errorw("[GetUserList] 链接grpc失败")
		return
	}
	fmt.Println("测试")
	global.UserSrvClient= proto.NewUserClient(conn)
	global.UserGrpcClient=conn
}
