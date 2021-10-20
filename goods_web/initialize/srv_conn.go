package initialize

import (
	"fmt"
	_ "github.com/mbobakov/grpc-consul-resolver" // It's important
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"project/goods_web/global"
	"project/goods_web/proto"
)
//负载
func InitSrvConn(){
	consulInfo := global.ServerConfig.ConsulInfo
	conn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", consulInfo.Host, consulInfo.Port, global.ServerConfig.GoodsSrvInfo.Name),
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		zap.S().Fatal("[InitSrvConn] 连接 【goods服务失败】")
	}
	global.GoodsSrvClient = proto.NewGoodsClient(conn)
	global.BannerSrvClient=proto.NewBannersClient(conn)
	global.CategorySrvClient=proto.NewCategoryClient(conn)
	global.BrandSrvClient=proto.NewBrandsClient(conn)
}