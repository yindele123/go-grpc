package initialize

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"project/userop_web/middlewares"
	"project/userop_web/router"
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
	ApiGroup := r.Group("/up/v1")
	router.InitAddressRouter(ApiGroup)
	router.InitMessageRouter(ApiGroup)
	router.InitUserFavRouter(ApiGroup)
	return r
}
