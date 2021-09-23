package initialize

import (
	"github.com/gin-gonic/gin"
	"project/user_web/middlewares"
	"project/user_web/router"
)

func Routers() *gin.Engine {
	r := gin.Default()
	r.Use(middlewares.Cors())
	ApiGroup := r.Group("/u/v1")
	router.InitUserRouter(ApiGroup)
	router.InitBaseRouter(ApiGroup)
	return r
}
