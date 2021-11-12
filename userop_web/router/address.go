package router

import (
	"github.com/gin-gonic/gin"
	"project/userop_web/api/address"
	"project/userop_web/middlewares"
)

func InitAddressRouter(Router *gin.RouterGroup) {
	AddressRouter := Router.Group("address").Use(middlewares.JWTAuth())
	{
		AddressRouter.GET("", address.List)          // 轮播图列表页
		AddressRouter.DELETE("/:id", address.Delete) // 删除轮播图
		AddressRouter.POST("", address.New)          //新建轮播图
		AddressRouter.PUT("/:id", address.Update)    //修改轮播图信息
	}
}
