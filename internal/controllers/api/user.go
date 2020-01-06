package api

import (
	"fmt"
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
	} else if p.Claims.Name != core.DefaultNilString && p.Claims.Email != core.DefaultNilString && (user.Name != p.Claims.Name || user.Email != p.Claims.Email) {
		user.Name, user.Email = p.Claims.Name, p.Claims.Email
		_ = user.UpdateInfo(core.Orm.New())
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

// ChargeQrCode 收费二维码
// @Tags  User 用户
// @Summary 收费二维码接口
// @Description 收费二维码
// @Produce json
// @Security ApiKeyAuth
// @Param currency_id query string true "币种ID"
// @Param type query string true "收费类型1：不带金额2：带金额"
// @Param symbol query string true "币种标示"
// @Param money query number true "金额"
// @Success 200 {object} resp.ChargeQrCodeResp
// @Router /user/qrcode/change [get]
func ChargeQrCode(c *gin.Context) {
	p := &params.ChargeQrCodeParam{
		Base: core.UserInfoPool.Get().(*params.BaseParam),
	}

	// 绑定参数
	if err := c.ShouldBind(p); err != nil {
		core.GResp.Failure(c, resp.CodeIllegalParam, err)
		return
	}

	qrcode := fmt.Sprintf("code=%d&symbol=%s&type=%d&money=%d&from=%s", p.Base.Uid, p.Symbol, p.Type, p.Money, "changre")
	core.GResp.Success(c, resp.ChargeQrCodeResp{
		UserName: p.Base.Claims.Name,
		Email:    p.Base.Claims.Email,
		Qrcode:   qrcode,
	})
	return
}

// PaymentQrCode 付款二维码
// @Tags  User 用户
// @Summary 付款二维码接口
// @Description 付款二维码
// @Produce json
// @Security ApiKeyAuth
// @Param symbol formData string true "币种标示"
// @Param money formData number true "金额"
// @Param code formData number true "code"
// @Success 200
// @Router /user/qrcode/pay [get]
func PaymentQrCode(c *gin.Context) {
	p := &params.PaymentQrCodeParam{
		Base: core.UserInfoPool.Get().(*params.BaseParam),
	}

	// 绑定参数
	if err := c.ShouldBind(p); err != nil {
		core.GResp.Failure(c, resp.CodeIllegalParam, err)
		return
	}
}

// UserPayInfo 用户付款
// @Tags  User 用户
// @Summary 用户付款接口
// @Description 用户付款
// @Produce json
// @Security ApiKeyAuth
// @Param code formData number true "code标示"
// @Param money formData number true "金额"
// @Param symbol formData string true "币种标识"
// @Param from formData string true "来源"
// @Param order_id formData string false "订单id"
// @Success 200
// @Router /user/pay/info [get]
func UserPayInfo(c *gin.Context) {
	p := &params.UserPayInfoParam{
		Base: core.UserInfoPool.Get().(*params.BaseParam),
	}

	// 绑定参数
	if err := c.ShouldBind(p); err != nil {
		core.GResp.Failure(c, resp.CodeIllegalParam, err)
		return
	}

}

// UserChange 用户付款
// @Tags  User 用户
// @Summary 用户付款接口
// @Description 用户付款
// @Produce json
// @Security ApiKeyAuth
// @Param code formData number true "code标示"
// @Param money formData number true "金额"
// @Param symbol formData string true "币种标识"
// @Param from formData string true "来源"
// @Param order_id formData string false "订单id"
// @Success 200
// @Router /user/pay/change [post]
func UserChange(c *gin.Context) {
	p := &params.UserChangeParam{
		Base: core.UserInfoPool.Get().(*params.BaseParam),
	}

	// 绑定参数
	if err := c.ShouldBind(p); err != nil {
		core.GResp.Failure(c, resp.CodeIllegalParam, err)
		return
	}
}
