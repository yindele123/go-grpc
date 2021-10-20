package initialize

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"project/goods_web/middlewares"
	"project/goods_web/router"
)

func Routers() *gin.Engine {
	r := gin.Default()
	r.Use(middlewares.Cors())
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK,gin.H{
			"code":http.StatusOK,
		})
	})
	//添加链路追踪
	ApiGroup := r.Group("/g/v1")
	router.InitGoodsRouter(ApiGroup)
	router.InitCategoryRouter(ApiGroup)
	router.InitBannerRouter(ApiGroup)
	router.InitBrandRouter(ApiGroup)
	return r
}
