package initialize

import (
	"go.uber.org/zap"
)

func NewLog() (*zap.Logger, error) {
	logConfig := zap.NewProductionConfig()
	//todo
	//调试模式下这会报错，找不到./log文件夹
	logConfig.OutputPaths = []string{
		//"./logs/my.log",
		"stderr",
	}
	return logConfig.Build()
}

func InitLogger() {
	logger, _ := NewLog()
	zap.ReplaceGlobals(logger)
}
