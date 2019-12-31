package router

import (
	"github.com/gin-gonic/gin"
	_ "github.com/haleyrom/wallet/docs"
	"github.com/haleyrom/wallet/pkg/middleware"
	"github.com/spf13/viper"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// InitRouter 初始化路由
func InitRouter() *gin.Engine {
	r := gin.Default()
	registerSwagger(r)
	r.Use(middleware.HttpBindGResp(), middleware.HttpCors())
	registerRouter(r)
	return r
}

// registerSwagger 注册swagger
func registerSwagger(r *gin.Engine) {
	url := ginSwagger.URL("/api/v1/wallet/swagger/doc.json")
	r.Group("api").Group("v1").Group("wallet").GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
}

// registerRouter 注册路由
func registerRouter(r *gin.Engine) {
	v1 := r.Group("api").Group("v1").Group("wallet")
	RegisterApiRouter(v1)
	RegisterAdminRouter(v1)

	if viper.GetString("runmode") == "debug" {
		RegisterTestDataRouter(v1)
	}
}
