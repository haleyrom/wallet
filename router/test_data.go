package router

import (
	"github.com/gin-gonic/gin"
	"github.com/haleyrom/wallet/internal/controllers/test_data"
)

// RegisterTestDataRouter 测试数据路由
func RegisterTestDataRouter(r *gin.RouterGroup) {
	v1 := r.Group("test_data")
	{
		v1.GET("/account/write", test_data.TestDataWriteAccount)
		v1.GET("/withdrawal/addr/is_local", test_data.TestDataWithdrawalIsLocal)
	}
}
