package base

import (
	"github.com/haleyrom/wallet/core"
	"github.com/haleyrom/wallet/internal/models"
	"github.com/haleyrom/wallet/pkg/consul"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

// CreateDeposit 创建钱包地址
func CreateDeposit(o *gorm.DB, uid uint) error {
	chain := models.NewBlockChain()
	if items, _ := chain.GetAll(core.Orm.New()); len(items) > 0 {
		addr := models.DepositAddr{
			Uid: uid,
		}

		for _, item := range items {
			go func(code string) {
				if data, err := consul.GetWalletAddress(code); err == nil && data != nil && len(data.Data.Address) > 0 {
					addr.BlockChainId = uint(item.Id)
					addr.Address = data.Data.Address
					addr.OrderId = data.Data.OrderId
					_ = addr.CreateDepositAddr(o)
				} else {
					logrus.Errorf("RegisterWalletAddr data :%v,failure :%v", data, err)
				}
			}(item.ChainCode)
		}
	}
	return nil
}
