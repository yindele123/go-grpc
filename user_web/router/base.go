package router

import (
	"github.com/gin-gonic/gin"
	"project/user_web/api"
)

func InitBaseRouter(group *gin.RouterGroup) {
	UserRouter := group.Group("base")
	{
		UserRouter.GET("/captcha", api.GetCaptcha)
		UserRouter.POST("send_sms", api.SendSms)
	}
}
