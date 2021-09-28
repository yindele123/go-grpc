package initialize

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"project/user_srv/global"
	"project/user_web/utils"
)

func GetEnvInfo(env string) bool {
	viper.AutomaticEnv()
	return viper.GetBool(env)
}
func InitConfig() {
	debug := GetEnvInfo("MXSHOP_DEBUG")
	configFilePrefix := "config"
	configFileName := fmt.Sprintf("%s/%s-pro.yaml", utils.GetCurrentPath(),configFilePrefix)
	if debug {
		configFileName = fmt.Sprintf("user_srv/%s-debug.yaml", configFilePrefix)
	}
	v := viper.New()
	v.SetConfigFile(configFileName)
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}
	if err := v.Unmarshal(global.ServerConfig); err != nil {
		panic(err)
	}
	//viper的功能 - 动态监控变化
	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		zap.S().Info("配置文件被人修改啦", e.Name)
		_ = v.ReadInConfig()
		_ = v.Unmarshal(global.ServerConfig)
	})
}
