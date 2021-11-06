package initialize

import (
	"fmt"
	_ "github.com/mbobakov/grpc-consul-resolver" // It's important
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"project/order_web/global"
	"project/order_web/proto"
)

func NacosClient(host string, port uint64, namespaceId string) (iClient naming_client.INamingClient, err error) {
	// 至少一个ServerConfig
	serverConfigs := []constant.ServerConfig{
		{
			IpAddr: host,
			Port:   port,
		},
	}
	// 创建clientConfig
	clientConfig := constant.ClientConfig{
		NamespaceId:         namespaceId, // 如果需要支持多namespace，我们可以场景多个client,它们有不同的NamespaceId。当namespace是public时，此处填空字符串。
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "/tmp/nacos/log",
		CacheDir:            "/tmp/nacos/cache",
		RotateTime:          "1h",
		MaxAge:              3,
	}
	// 创建服务发现客户端
	configClient, err := clients.CreateNamingClient(map[string]interface{}{
		"serverConfigs": serverConfigs,
		"clientConfig":  clientConfig,
	})
	return configClient, err
}

//负载
func InitSrvConn() {
	consulInfo := global.ServerConfig.ConsulInfo
	conn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", consulInfo.Host, consulInfo.Port, global.ServerConfig.GoodsSrvInfo.Name),
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		zap.S().Fatal("[InitSrvConn] 连接 【goods服务失败】")
		return
	}
	namingClient, err := NacosClient(global.NacosConfig.NacosInfo.Host, global.NacosConfig.NacosInfo.Port, global.NacosConfig.NacosInfo.NamespaceId)
	if err != nil {
		zap.S().Fatal("[InitSrvConn] 连接 【order服务失败】")
		return
	}

	instance, err := namingClient.SelectOneHealthyInstance(vo.SelectOneHealthInstanceParam{
		ServiceName: global.ServerConfig.OrderSrvInfo.Name,
	})
	if err != nil {
		zap.S().Fatal("[InitSrvConn] 连接 【order服务失败】")
		return
	}
	orderConn, err := grpc.Dial(fmt.Sprintf("%s:%d", instance.Ip, instance.Port), grpc.WithInsecure())
	if err != nil {
		zap.S().Fatal("[InitSrvConn] 连接 【order服务失败】")
		return
	}

	invConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", consulInfo.Host, consulInfo.Port, global.ServerConfig.InvSrvInfo.Name),
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		zap.S().Fatal("[InitSrvConn] 连接 【库存服务失败】")
		return
	}

	global.GoodsSrvClient = proto.NewGoodsClient(conn)
	global.OrderSrvClient = proto.NewOrderClient(orderConn)
	global.InvSrvClient = proto.NewInventoryClient(invConn)
}
