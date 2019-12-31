package api

import (
	"github.com/gin-gonic/gin"
	"github.com/haleyrom/wallet/core"
	"github.com/haleyrom/wallet/internal/models"
	"github.com/haleyrom/wallet/internal/params"
	"github.com/haleyrom/wallet/internal/resp"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"sync"
)

// Check Check
func Check(c *gin.Context) {
	core.GResp.Success(c, "ok")
	return
}

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
	}

	p.Uid = user.ID
	return nil
}

// UpdatePayPassword 更新支付密码
// @Tags  User 用户
// @Summary 更新支付密码接口
// @Description 更新支付密码
// @Produce json
// @Security ApiKeyAuth
// @Param password formData string true "支付密码"
// @Success 200
// @Router /user/pay/update [post]
func UpdatePayPassword(c *gin.Context) {
	p := &params.UpdatePayPasswordParam{
		Base: core.UserInfoPool.Get().(*params.BaseParam),
	}

	// 绑定参数
	if err := c.ShouldBind(p); err != nil {
		core.GResp.Failure(c, resp.CodeIllegalParam, err)
		return
	}

	user := models.NewUser()
	user.ID = p.Base.Uid
	user.PayPassword = p.Password
	if err := user.UpdatePayPassword(core.Orm.New()); err != nil {
		core.GResp.Failure(c, err)
		return
	}
	core.GResp.Success(c, resp.EmptyData())
	return
}
