package main

import (
	"fmt"
	"github.com/gin-gonic/gin/binding"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/nacos-group/nacos-sdk-go/inner/uuid"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"project/order_web/global"
	"project/order_web/initialize"
	"project/order_web/utils/register"
	validator2 "project/order_web/validator"
	"syscall"
)

func main() {
	initialize.InitConfig()
	initialize.InitRedis()
	if transErr := initialize.InitTrans("zh"); transErr != nil {
		zap.S().Panic("初始化翻译器错误:", transErr.Error())
	}
	initialize.InitLogger()
	//Routers
	//zap.S().Panic("启动失败:")
	r := initialize.Routers()
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("mobile", validator2.ValidateMobile)
		_ = v.RegisterTranslation("mobile", global.Trans, func(ut ut.Translator) error {
			return ut.Add("mobile", "{0} 非法的手机号码!", true) // see universal-translator for details
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("mobile", fe.Field())
			return t
		})
	}

	initialize.InitSrvConn()

	//注册服务

	//ConsulRegister
	uuid, _ := uuid.NewV4()
	serviceId := fmt.Sprintf("%s", uuid)
	var consulRegister register.Register = register.ConsulRegister{
		Host: global.ServerConfig.ConsulInfo.Host,
		Port: global.ServerConfig.ConsulInfo.Port,
	}
	rerr := consulRegister.Register(global.ServerConfig.Host, global.ServerConfig.Port, global.ServerConfig.Name, []string{"xindele", "yindele123", "order-web"}, serviceId)
	if rerr != nil {
		zap.S().Panic("注册服务失败:", rerr.Error())
	}
	go func() {
		_ = r.Run(fmt.Sprintf(":%d", global.ServerConfig.Port))
		/*if err := r.Run(fmt.Sprintf(":%d", global.ServerConfig.Port)); err != nil {
			zap.S().Panic("启动失败:", err.Error())
		}*/
	}()

	//接收终止信号
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	if err := consulRegister.Deregister(serviceId); err != nil {
		zap.S().Info("注销失败:", err.Error())
	} else {
		zap.S().Info("注销成功:")
	}
}
