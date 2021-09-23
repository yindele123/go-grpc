package router

import (
	"github.com/gin-gonic/gin"
	"project/user_web/api"
	"project/user_web/middlewares"
)

func InitUserRouter(group *gin.RouterGroup) {
	UserRouter := group.Group("user")
	{
		//GetUserList
		UserRouter.GET("/list",middlewares.JWTAuth(),middlewares.IsAdminAuth(), api.GetUserList)
		UserRouter.POST("pwd_login", api.PassWordLogin)
		UserRouter.POST("send_sms", api.SendSms)
	}
}
