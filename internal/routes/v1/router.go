package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/yz626/edu-chain/pkg/logger"
)

// RouterSet Wire Provider Set
var RouterSet = wire.NewSet(
	NewRouter,
)

func NewRouter(log *logger.Logger) *gin.Engine {
	router := gin.New()

	register(router)
	return router
}

// register 注册路由
func register(router *gin.Engine) {
}

// 注册中间件
func initMiddlewares(router *gin.Engine) {

}
