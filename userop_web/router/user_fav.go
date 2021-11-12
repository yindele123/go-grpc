package router

import (
	"github.com/gin-gonic/gin"
	"project/userop_web/api/userFav"
	"project/userop_web/middlewares"
)

func InitUserFavRouter(Router *gin.RouterGroup) {
	UserFavRouter := Router.Group("userfavs")
	{
		UserFavRouter.DELETE("/:id", middlewares.JWTAuth(), userFav.Delete) // 删除收藏记录
		UserFavRouter.GET("/:id", middlewares.JWTAuth(), userFav.Detail)    // 获取收藏记录
		UserFavRouter.POST("", middlewares.JWTAuth(), userFav.New)          //新建收藏记录
		UserFavRouter.GET("", middlewares.JWTAuth(), userFav.List)          //获取当前用户的收藏
	}
}
