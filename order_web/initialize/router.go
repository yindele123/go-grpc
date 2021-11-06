package initialize

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"project/order_web/middlewares"
	"project/order_web/router"
)

func Routers() *gin.Engine {
	r := gin.Default()
	r.Use(middlewares.Cors())
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusOK,
		})
	})
	//添加链路追踪
	ApiGroup := r.Group("/o/v1")
	router.InitOrderRouter(ApiGroup)
	router.InitShopCartRouter(ApiGroup)
	return r
}
