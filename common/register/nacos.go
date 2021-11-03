package register

import (
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"go.uber.org/zap"
	"project/order_srv/global"
)

type NacosRegister struct {
	Host string
	Port uint64
}

func (n NacosRegister) Client() (iClient naming_client.INamingClient, err error) {
	// 至少一个ServerConfig
	serverConfigs := []constant.ServerConfig{
		{
			IpAddr: n.Host,
			Port:   n.Port,
		},
	}
	// 创建clientConfig
	clientConfig := constant.ClientConfig{
		NamespaceId:         global.NacosConfig.NacosInfo.NamespaceId, // 如果需要支持多namespace，我们可以场景多个client,它们有不同的NamespaceId。当namespace是public时，此处填空字符串。
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

func (n NacosRegister) Register(address string, port int, name string, tags interface{}, id string) error {
	configClient, clientErr := n.Client()
	if clientErr != nil {
		zap.S().Panic("链接nacos失败!", clientErr.Error())
		return clientErr
	}
	_, err := configClient.RegisterInstance(vo.RegisterInstanceParam{
		Ip:          address,
		Port:        uint64(port),
		ServiceName: name,
		Weight:      10,
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true,
		Metadata:    tags.(map[string]string),
	})
	if err != nil {
		zap.S().Panic("链接nacos失败!", err.Error())
		return err
	}
	return nil

}

func (n NacosRegister) Deregister(serviceId string) error {
	configClient, clientErr := n.Client()
	if clientErr != nil {
		zap.S().Panic("链接nacos失败!", clientErr.Error())
		return clientErr
	}
	_, err := configClient.DeregisterInstance(vo.DeregisterInstanceParam{
		Ip:          global.ServerConfig.Host,
		Port:        uint64(global.ServerConfig.Port),
		ServiceName: global.ServerConfig.ServiceName,
		Ephemeral:   true,
	})
	if err != nil {
		zap.S().Panic("链接nacos失败!", err.Error())
		return err
	}
	return nil
}
