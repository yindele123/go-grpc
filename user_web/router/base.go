package router

import (
	"github.com/gin-gonic/gin"
	"project/user_web/api"
)

func InitBaseRouter(group *gin.RouterGroup) {
	UserRouter := group.Group("base")
	{
		//GetUserList
		UserRouter.GET("/captcha", api.GetCaptcha)
	}
}
