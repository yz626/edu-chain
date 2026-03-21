package v1

import (
	"github.com/gin-gonic/gin"
)

// Register 注册路由
func Register(router *gin.Engine) {

	register(router)
}

// register 注册路由
func register(router *gin.Engine) {
	router.GET("/ping", ping)
	router.NoRoute(notFound)
}

// 注册中间件
func initMiddlewares(router *gin.Engine) {

}

// <===== 路由辅助函数 =====>

// 心跳路由
func ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

// 404路由
func notFound(c *gin.Context) {
	c.JSON(404, gin.H{
		"message": "404 Not Found",
	})
}
