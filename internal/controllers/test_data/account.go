package test_data

import (
	"github.com/gin-gonic/gin"
	"github.com/haleyrom/wallet/core"
	"github.com/haleyrom/wallet/internal/models"
	"github.com/haleyrom/wallet/internal/resp"
)

// TestDataWriteAccount 账本测试数据写入
func TestDataWriteAccount(c *gin.Context) {
	items := make([]models.Account, 0)
	o := core.Orm.New().Begin()

	o.Table(models.GetAccountTable()).Find(&items)

	if len(items) > core.DefaultNilNum {
		account := models.NewAccount()
		detail := make([]models.AccountDetail, 0)
		money := float64(100000)
		for _, item := range items {
			account.Uid = item.Uid
			account.CurrencyId = item.CurrencyId
			if err := account.UpdateBalance(o, core.OperateToUp, money); err == nil {
				detail = append(detail, models.AccountDetail{
					Uid:         item.Uid,
					AccountId:   item.ID,
					Balance:     item.Balance + money,
					LastBalance: item.Balance,
					Income:      money,
					Spend:       float64(core.DefaultNilNum),
					Type:        resp.AccountDetailUp,
					OrderId:     uint(core.DefaultNilNum),
				})
			}
		}

		if len(detail) > core.DefaultNilNum {
			if err := models.NewAccountDetail().CreateAccountDetailAll(o, detail); err == nil {
				o.Rollback()
				core.GResp.Failure(c, resp.CodeUnknow)
				return
			}
		}
	}
	o.Commit()
	core.GResp.Success(c, resp.EmptyData())
	return
}
