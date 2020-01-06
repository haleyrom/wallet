package router

import (
	"github.com/gin-gonic/gin"
	"github.com/haleyrom/wallet/internal/controllers/admin"
	"github.com/haleyrom/wallet/pkg/middleware"
)

// RegisterAdminRouter 后台路由
func RegisterAdminRouter(r *gin.RouterGroup) {
	v1 := r.Group("admin")
	{
		v1.Use(middleware.HttpInterceptor(middleware.AdminJwt))
		v1.GET("/account/user/list", admin.AccountUserList)
		v1.GET("/account/withdrawal/list", admin.AccountWithdrawalList)
		v1.GET("/account/user/deposit-list", admin.DepositDetailList)
		v1.GET("/account/user/order-list", admin.OrderList)
		v1.GET("/account/withdrawal/detail", admin.AccountWithdrawalDetail)
		v1.POST("/account/withdrawal/customer", admin.WithdrawalDetailCustomer)
		v1.POST("/account/withdrawal/financial", admin.WithdrawalDetailFinancial)
		v1.GET("/account/company/deposit/list", admin.CompanyDepositList)
		v1.GET("/account/company/withdrawal/list", admin.CompanyWithdrawalList)
		v1.GET("/account/company/deposit_addr/list", admin.CompanyDepositAddrList)
		v1.GET("/account/company/withdrawal_addr/list", admin.CompanyWithdrawalAddrList)
		v1.POST("/account/company/deposit_addr/join", admin.JoinCompanyDepositAddr)
		v1.POST("/account/company/withdrawal_addr/join", admin.JoinCompanyWithdrawalAddr)
		v1.POST("/account/company/addr/update", admin.UpdateCompanyAddr)
		v1.POST("/account/company/addr/status", admin.UpdateCompanyAddrStatus)
		v1.POST("/account/company/addr/create", admin.CreateCompanyAddr)
		v1.POST("/currency/quote/create", admin.CreateCurrencyQuote)
		v1.POST("/currency/quote/update", admin.UpdateCurrencyQuote)
		v1.GET("/currency/quote/list", admin.ReadQuotePage)
		v1.GET("/currency/quote_history/list", admin.ReadQuoteHistoryPage)
		v1.POST("/operate/recharge/join", admin.JoinRecharge)
		v1.GET("/operate/recharge/list", admin.ReadRechargePage)
		v1.POST("/operate/recharge/remove", admin.RemoveRecharge)
		v1.POST("/operate/recharge/audio", admin.AudioRecharge)
	}
}
