package router

import (
	"github.com/gin-gonic/gin"
	"project/user_web/api"
)

func InitUserRouter(group *gin.RouterGroup) {
	UserRouter := group.Group("user")
	{
		//GetUserList
		UserRouter.GET("/list", api.GetUserList)
		UserRouter.POST("pwd_login", api.PassWordLogin)
	}
}
