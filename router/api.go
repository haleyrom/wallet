package router

import (
	"github.com/gin-gonic/gin"
	"github.com/haleyrom/wallet/internal/controllers/api"
	"github.com/haleyrom/wallet/pkg/middleware"
)

// RegisterApiRouter 注册api文件
func RegisterApiRouter(r *gin.RouterGroup) {
	r.Use(middleware.HttpInterceptor(middleware.Jwt))
	r.GET("/check", api.Check)
	r.POST("/user/pay/update", api.UpdatePayPassword)
	r.POST("/user/pay/set-paypwd", api.SetPayPassWordHandler)
	r.POST("/user/pay/reset-paypwd", api.ReSetPayWordHandler)
	r.POST("/user/pay/send-email_util", api.SendEmailPayHandler)
	r.GET("/user/pay/is_init", api.IsSetPassWord)

	r.GET("/account/info", api.AccountInfo)
	r.GET("/account/tfor/info", api.AccountTFORInfo)
	r.GET("/account/detail", api.AccountDetail)
	r.POST("/account/transfer", api.AccountTransfer)
	r.POST("/account/change", api.AccountChange)
	r.POST("/account/share/bonus", api.AccountShareBonus)
	r.POST("/account/withdrawal", api.AccountWithdrawal)

	r.GET("/currency/list", api.ReadCurrencyList)
	r.GET("/currency/transfer_list", api.ReadCurrencyTransferList)
	r.POST("/currency/update", api.UpdateCurrency)
	r.POST("/currency/add", api.AddCurrency)
	r.POST("/currency/remove", api.RemoveCurrency)
	r.POST("/currency/status", api.UpdateCurrencyStatus)
	r.GET("/currency/quote", api.CurrencyQuote)

	r.GET("/chain/list", api.ReadListBlockChain)
	r.POST("/chain/update", api.UpdateBlockChain)
	r.POST("/chain/add", api.AddBlockChain)
	r.POST("/chain/remove", api.RemoveBlockChain)
	r.GET("/chain/symbol", api.ReadListSymbolBlockChain)

	r.GET("/coin/list", api.ReadCoinList)
	r.GET("/coin/info", api.ReadCoinInfo)
	r.POST("/coin/update", api.UpdateCoin)
	r.POST("/coin/add", api.AddCoin)
	r.POST("/coin/status", api.UpdateCoinStatus)
	r.POST("/coin/remove", api.RemoveCoin)
	r.GET("/coin/deposit", api.ReadCoinDepositInfo)

	r.GET("/deposit/list", api.ReadDepositAddList)
	r.POST("/deposit/join", api.JoinDepositDetail)
	r.GET("/deposit/info", api.ReadDepositAddr)
	r.POST("/deposit/top_up", api.TopUpDeposit)
	r.GET("/deposit/detail", api.ReadDepositDetail)

	r.GET("/withdrawal/detail", api.ReadWithdrawalDetail)
	r.GET("/withdrawal/addr/list", api.ReadWithdrawalAddrList)
	r.POST("/withdrawal/addr/add", api.CreateWithdrawalAddr)
	r.POST("/withdrawal/addr/update", api.UpdateWithdrawalAddr)
	r.POST("/withdrawal/addr/remove", api.RemoveWithdrawalAddr)
	r.POST("/withdrawal/callback", api.WithdrawalCallback)
}
