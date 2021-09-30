package initialize

import (
	"encoding/json"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"project/user_web/global"
)

func GetEnvInfo(env string)  bool{
	viper.AutomaticEnv()
	return viper.GetBool(env)
}
func InitConfig() {
	debug := GetEnvInfo("MXSHOP_DEBUG")
	configFilePrefix := "config"
	configFileName := fmt.Sprintf("user_web/%s-pro.yaml", configFilePrefix)
	if debug {
		configFileName = fmt.Sprintf("user_web/%s-debug.yaml", configFilePrefix)
	}
	v:=viper.New()
	v.SetConfigFile(configFileName)
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}
	if err := v.Unmarshal(global.NacosConfig); err != nil {
		panic(err)
	}
	//viper的功能 - 动态监控变化
	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		zap.S().Info("配置文件被人修改啦", e.Name)
		_ = v.ReadInConfig()
		_ = v.Unmarshal(global.ServerConfig)
	})

	// 至少一个ServerConfig
	serverConfigs := []constant.ServerConfig{
		{
			IpAddr:      global.NacosConfig.NacosInfo.Host,
			Port:        global.NacosConfig.NacosInfo.Port,
		},
	}
	fmt.Println(global.NacosConfig.NacosInfo)
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

	// 创建动态配置客户端
	configClient, configErr := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": serverConfigs,
		"clientConfig":  clientConfig,
	})
	if configErr!=nil {
		zap.S().Panic("链接nacos失败!", configErr.Error())
	}
	content, contentErr := configClient.GetConfig(vo.ConfigParam{
		DataId: global.NacosConfig.NacosInfo.DataId,
		Group:  global.NacosConfig.NacosInfo.Group})
	if contentErr!=nil {
		zap.S().Panic("获取nacos配置失败!", contentErr.Error())
	}

	fmt.Println(content)

	err := json.Unmarshal([]byte(content), &global.ServerConfig)
	if err != nil {
		zap.S().Fatalf("读取nacos配置失败： %s", err.Error())
	}
}
