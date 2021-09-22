package initialize

import (
	"github.com/gin-gonic/gin"
	"project/user_web/router"
)

func Routers() *gin.Engine {
	r := gin.Default()
	ApiGroup := r.Group("/u/v1")
	router.InitUserRouter(ApiGroup)
	return r
}
