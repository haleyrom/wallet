package base

import (
	"github.com/haleyrom/wallet/core"
	"github.com/haleyrom/wallet/internal/models"
	"github.com/haleyrom/wallet/internal/params"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"sync"
)

// CreateUser 创建用户
func CreateUser(p *params.BaseParam) error {
	user := models.User{
		Uid:   p.Claims.UserID,
		Name:  p.Claims.Name,
		Email: p.Claims.Email,
	}

	o := core.Orm.New()
	if err := user.IsExistUser(o); err != nil {
		o = o.Begin()
		// 不存在用户，创建
		if err = user.CreateUser(o); err != nil {
			o.Rollback()
			return err
		} else if err = o.Commit().Error; err == nil {
			var wg sync.WaitGroup
			wg.Add(2)

			go func(id uint, o *gorm.DB) {
				defer wg.Done()
				// 创建账本
				ids, _ := models.NewCurrency().GetIdAll(o)
				account := models.Account{
					Uid: id,
				}

				if err := account.CreateAccount(o, ids); err != nil {
					logrus.Errorf("create account uid %d, failure: %v", id, err)
				}

			}(user.ID, core.Orm.New())

			go func(id uint, o *gorm.DB) {
				defer wg.Done()
				_ = CreateDeposit(o, id)
			}(user.ID, core.Orm.New())

			wg.Wait()
		}
	} else if p.Claims.Name != core.DefaultNilString && p.Claims.Email != core.DefaultNilString && (user.Name != p.Claims.Name || user.Email != p.Claims.Email) {
		user.Name, user.Email = p.Claims.Name, p.Claims.Email
		_ = user.UpdateInfo(core.Orm.New())
	}

	p.Uid = user.ID
	return nil
}
