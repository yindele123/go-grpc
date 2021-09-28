package initialize

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"project/user_web/middlewares"
	"project/user_web/router"
)

func Routers() *gin.Engine {
	r := gin.Default()
	r.Use(middlewares.Cors())
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK,gin.H{
			"code":http.StatusOK,
		})
	})
	ApiGroup := r.Group("/u/v1")
	router.InitUserRouter(ApiGroup)
	router.InitBaseRouter(ApiGroup)
	return r
}
