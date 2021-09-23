package main

import (
	"fmt"
	"github.com/gin-gonic/gin/binding"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"project/user_web/global"
	"project/user_web/initialize"
	validator2 "project/user_web/validator"
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

	if err := r.Run(fmt.Sprintf(":%d", global.ServerConfig.Port)); err != nil {
		zap.S().Panic("启动失败:", err.Error())
	}
}
