package base

import (
	"github.com/gin-gonic/gin"
	"github.com/haleyrom/wallet/core"
	"github.com/haleyrom/wallet/internal/models"
	"github.com/haleyrom/wallet/internal/params"
	"github.com/haleyrom/wallet/internal/resp"
	"github.com/haleyrom/wallet/pkg/consul"
	"github.com/jinzhu/gorm"
	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
	"sync"
)

// CreateUser 创建用户
func CreateUser(c *gin.Context, p *params.BaseParam) error {
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
	} else {
		go UpdateUserInfo(c, p, user)
	}

	p.Uid = user.ID
	return nil
}

// UpdateUserInfo 根据用户信息
func UpdateUserInfo(c *gin.Context, p *params.BaseParam, user models.User) {
	data, err := GetConsulUserInfo(c, p.Claims.UserID)
	if err == nil {
		if p.Claims.Name != core.DefaultNilString && p.Claims.Email != core.DefaultNilString && (data.Nickname != p.Claims.Name || data.Email != p.Claims.Email) {
			user.Name, user.Email = data.Nickname, data.Email
			_ = user.UpdateInfo(core.Orm.New())
		}
	}
}

// GetConsulUserInfo 获取consul用户信息
func GetConsulUserInfo(c *gin.Context, uid string) (resp.UserInfoResp, error) {
	var data resp.UserInfoResp
	// 查询用户创建用户
	result, err := consul.GetUserInfo(uid, c.Request.Header.Get(core.HttpHeadToken))

	if err != nil {
		return data, resp.CodeNotUser
	} else if err = mapstructure.Decode(result, &data); err != nil {
		return data, resp.CodeNotUser
	}
	return data, nil
}
