package api

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/haleyrom/wallet/core"
	"github.com/haleyrom/wallet/internal/models"
	"github.com/haleyrom/wallet/internal/params"
	"github.com/haleyrom/wallet/internal/resp"
	"github.com/haleyrom/wallet/pkg/consul"
	"github.com/haleyrom/wallet/pkg/tools"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"time"
)

// failure_timer 无效时间
const failure_timer = time.Minute * 1

// Check Check
func Check(c *gin.Context) {
	core.GResp.Success(c, "ok")
	return
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

	qrcode := fmt.Sprintf("code=%d&symbol=%s&type=%d&money=%s&from=%s", p.Base.Uid, p.Symbol, p.Type, p.Money, "changre")
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
// @Success 200 {object} resp.ChargeQrCodeResp
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

	o := core.Orm.New()
	// 获取代笔信息
	currency := models.NewCurrency()
	currency.Symbol = p.Symbol
	if err := currency.GetSymbolById(o); err != nil {
		o.Rollback()
		core.GResp.Failure(c, resp.CodeNotCurrency, err)
		return
	}

	// 创建订单
	jsonStr, _ := json.Marshal(c.Request.PostForm)
	order := &models.Order{
		Uid:         p.Base.Uid,
		Context:     string(jsonStr),
		CurrencyId:  currency.ID,
		ExchangeUid: uint(core.DefaultNilNum),
		ExchangeId:  currency.ID,
		Balance:     float64(core.DefaultNilNum),
		Ratio:       float64(core.DefaultNilNum),
		Status:      models.OrderStatusNot,
		Type:        models.OrderTypePayment,
		Form:        models.OrderFormPayment,
		Symbol:      p.Symbol,
	}
	if err := order.CreateOrder(o); err != nil {
		o.Callback()
		core.GResp.Failure(c, errors.Errorf("join order fail:%s", err))
		return
	}

	qrcode := fmt.Sprintf("code=%d&symbol=%s&type=%d&money=%s&from=%s&order_id=%s", p.Base.Uid, p.Symbol, p.Type, p.Money, "payment", order.OrderUuid)

	core.PayChan.MapChan[order.OrderUuid] = make(chan int, 1)
	core.PayChan.MapChan[order.OrderUuid] <- core.DefaultNilNum
	key := int(order.CreatedAt.Add(failure_timer).Unix())
	if len(core.PayChan.MapTime[key]) == core.DefaultNilNum {
		core.PayChan.MapTime[key] = make([]string, 0)
	}
	core.PayChan.MapTime[key] = append(core.PayChan.MapTime[key], order.OrderUuid)

	core.GResp.Success(c, resp.ChargeQrCodeResp{
		UserName: p.Base.Claims.Name,
		Email:    p.Base.Claims.Email,
		Qrcode:   qrcode,
		MinMoney: currency.MinPayMoney,
		OrderId:  order.OrderUuid,
	})
	return

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
// @Success 200 {object} resp.UserPayInfoResp
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

	o := core.Orm.New()
	user := models.NewUser()
	user.ID = p.Code
	if err := user.GetInfo(o); err != nil {
		core.GResp.Failure(c, resp.CodeNotUser, err)
		return
	}

	result, err := consul.GetUserInfo(user.Uid, c.Request.Header.Get(core.HttpHeadToken))
	var data resp.UserInfoResp

	if err != nil {
		o.Rollback()
		core.GResp.Failure(c, resp.CodeNotUser)
		return
	} else if err = mapstructure.Decode(result, &data); err != nil {
		core.GResp.Failure(c, resp.CodeNotUser)
		return
	}

	core.GResp.Success(c, resp.UserPayInfoResp{
		OrderId: p.OrderId,
		Symbol:  p.Symbol,
		Email:   data.Email,
		Money:   p.Money,
	})
	return
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
// @Param pay_password formData string false "支付密码"
// @Success 200
// @Router /user/pay/change [post]
func UserChange(c *gin.Context) {
	var err error
	p := &params.UserChangeParam{
		Base: core.UserInfoPool.Get().(*params.BaseParam),
	}

	// 绑定参数
	if err := c.ShouldBind(p); err != nil {
		core.GResp.Failure(c, resp.CodeIllegalParam, err)
		return
	}

	if p.Base.Uid == p.Code {
		core.GResp.Failure(c, resp.CodeOneselfInto)
		return
	}

	// 验证支付密码
	o := core.Orm.New().Begin()
	user := models.NewUser()
	user.ID = p.Base.Uid
	if err = user.GetInfo(o); err != nil {
		o.Callback()
		core.GResp.Failure(c, resp.CodeNotUser)
		return
	}

	// 获取代笔信息
	currency := models.NewCurrency()
	currency.Symbol = p.Symbol
	if err := currency.GetSymbolById(o); err != nil {
		o.Rollback()
		core.GResp.Failure(c, resp.CodeNotCurrency, err)
		return
	}

	if (p.From == "payment" && p.Money < currency.MinPayMoney) == false {
		if user.PayPassword != tools.Hash256(p.PayPassword, tools.NewPwdSalt(p.Base.Claims.UserID, 1)) {
			o.Callback()
			core.GResp.Failure(c, resp.CodeErrorPayPassword)
			return
		}
	}

	// 获取金额
	account := models.NewAccount()
	account.CurrencyId = currency.ID
	list := make(map[uint]models.Account, 0)
	if list, err = account.GetOrderUidSByCurrencyInfo(o, []uint{p.Base.Uid, p.Code}); err != nil {
		o.Rollback()
		core.GResp.Failure(c, resp.CodeNotAccount, err)
		return
	}

	// 余额不足
	if (list[p.Base.Uid].Balance*100 - list[p.Base.Uid].BlockedBalance*100) < p.Money*100 {
		o.Rollback()
		core.GResp.Failure(c, resp.CodeLessMoney)
		return
	}

	order := models.NewOrder()
	jsonStr, _ := json.Marshal(c.Request.PostForm)

	// 判断订单id
	if len(p.OrderId) > core.DefaultNilNum {
		order.OrderUuid = p.OrderId
		err := order.IsOrderUuid(o)
		if err != nil {
			o.Rollback()
			core.GResp.Failure(c, resp.CodeNotOrderId, err)
			return
		} else if order.Status == models.OrderStatusOk {
			o.Rollback()
			core.GResp.Failure(c, resp.CodeOrderStatusOK)
			return
		} else if order.CreatedAt.Add(failure_timer).Unix() < time.Now().Unix() {
			_ = order.RemoveOrder(o)
			o.Commit()
			core.GResp.Failure(c, resp.CodeFailureQrCode)
			return
		}

		order.Context = string(jsonStr)
		order.Balance, order.ExchangeUid = p.Money, p.Base.Uid
		if err = order.UpdateStatusOk(o); err != nil {
			o.Rollback()
			core.GResp.Failure(c, errors.Errorf("update order fail:%s", err))
			return
		}
	} else {
		currency := models.NewCurrency()
		currency.ID = list[p.Base.Uid].CurrencyId
		_ = currency.IsExistCurrency(o)

		// 创建订单
		order = &models.Order{
			Uid:         p.Base.Uid,
			Context:     string(jsonStr),
			CurrencyId:  list[p.Base.Uid].CurrencyId,
			ExchangeUid: p.Code,
			ExchangeId:  list[p.Code].CurrencyId,
			Balance:     p.Money,
			Ratio:       float64(core.DefaultNilNum),
			Status:      models.OrderStatusOk,
			Type:        models.OrderTypePayment,
			Form:        models.OrderFormPayment,
			Symbol:      currency.Symbol,
		}
		if err := order.CreateOrder(o); err != nil {
			o.Callback()
			core.GResp.Failure(c, errors.Errorf("join order fail:%s", err))
			return
		}
	}

	// 扣费
	if err := AccountOperate(o, list[p.Base.Uid], p.Money, core.OperateToOut, resp.AccountDetailPayment, order.ID); err != nil {
		o.Rollback()
		core.GResp.Failure(c, resp.CodeLessMoney)
		return
	} else if err = AccountOperate(o, list[p.Code], p.Money, core.OperateToUp, resp.AccountDetailGather, order.ID); err != nil {
		o.Rollback()
		core.GResp.Failure(c, err)
		return
	}

	o.Commit()
	core.GResp.Success(c, resp.EmptyData())
	return
}

// UserPayQrCodeStatus 付款码激活状态查询
// @Tags  User 用户
// @Summary 付款码激活状态查询接口
// @Description 付款码激活状态查询 code:200已支付,101138未支付 对接时候返回code为101138时候需要继续请求数据，其他情况下刷新二维码 最长等待时间30秒
// @Produce json
// @Security ApiKeyAuth
// @Param order_id formData string false "订单id"
// @Success 200
// @Router /user/pay/status [get]
func UserPayQrCodeStatus(c *gin.Context) {
	p := &params.UserPayStatusParam{
		Base: core.UserInfoPool.Get().(*params.BaseParam),
	}

	// 绑定参数
	if err := c.ShouldBind(p); err != nil {
		core.GResp.Failure(c, resp.CodeIllegalParam, err)
		return
	}

	fmt.Println(core.PayChan.MapChan, "+++++++")
	if _, ok := core.PayChan.MapChan[p.OrderId]; ok == false {
		core.GResp.Failure(c, resp.CodeFailureQrCode)
		return
	}

	for {
		select {
		case <-core.PayChan.MapChan[p.OrderId]:
			delete(core.PayChan.MapChan, p.OrderId)
			core.GResp.Success(c, resp.EmptyData())
			return
		case <-time.After(100 * time.Millisecond):
			core.GResp.Failure(c, resp.CodeWaitQrCode)
			return
		}
	}
}
