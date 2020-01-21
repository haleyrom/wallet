package test_data

import (
	"github.com/gin-gonic/gin"
	"github.com/haleyrom/wallet/core"
	"github.com/haleyrom/wallet/internal/models"
	"github.com/haleyrom/wallet/internal/resp"
	"time"
)

// TestDataWithdrawalIsLocal 测试提笔地址
func TestDataWithdrawalIsLocal(c *gin.Context) {
	o := core.Orm.New()

	items := make([]models.WithdrawalDetail, 0)
	if err := o.Where("address_source = 0").Find(&items).Error; err == nil {
		if len(items) > core.DefaultNilNum {
			deposit := models.NewDepositAddr()
			data := map[string]interface{}{
				"updated_at": time.Now(),
			}

			for _, item := range items {
				// 判断地址是否合法
				deposit.Address = item.Address
				if err := deposit.IsAddress(o); err == nil {
					data["address_source"] = models.WithdrawalAddrLocal
				} else {
					data["address_source"] = models.WithdrawalAddrBack
				}
				o.Table(models.GetWithdrawalAddrTable()).Where("id = ?", item.ID).Update(data)
			}
		}
	}

	core.GResp.Success(c, resp.EmptyData())
	return
}
